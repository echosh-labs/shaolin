package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"shaolin/backend/internal/api"
	"shaolin/backend/internal/domain"
	"shaolin/backend/internal/repository/sqlite"
)

func TestStoreAndCartAPI(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userRepo := sqlite.NewSQLiteUserRepository(db)
	storeRepo := sqlite.NewSQLiteStoreRepository(db)
	storeHandler := api.NewStoreHandler(storeRepo)

	// Create test user
	student := &domain.User{
		ID:           "u-student1",
		Email:        "student1@shaolin.com",
		PasswordHash: "pwd",
		FirstName:    "Store",
		LastName:     "Customer",
		DateOfBirth:  "1990-01-01",
		Role:         "student",
		CurrentRank:  "White Belt",
		CreatedAt:    time.Now(),
	}
	userRepo.Create(student)
	token, _ := api.GenerateToken(student)

	// Mux Registration
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/store/products", storeHandler.ListProducts)
	mux.Handle("GET /api/cart", api.AuthMiddleware(http.HandlerFunc(storeHandler.GetCart)))
	mux.Handle("POST /api/cart", api.AuthMiddleware(http.HandlerFunc(storeHandler.AddToCart)))
	mux.Handle("DELETE /api/cart/{id}", api.AuthMiddleware(http.HandlerFunc(storeHandler.RemoveFromCart)))
	mux.Handle("POST /api/store/checkout", api.AuthMiddleware(http.HandlerFunc(storeHandler.Checkout)))
	mux.Handle("GET /api/store/orders/history", api.AuthMiddleware(http.HandlerFunc(storeHandler.GetOrderHistory)))

	// --- 1. Test ListProducts ---
	reqCatalog := httptest.NewRequest("GET", "/api/store/products", nil)
	wCatalog := httptest.NewRecorder()
	mux.ServeHTTP(wCatalog, reqCatalog)

	if wCatalog.Code != http.StatusOK {
		t.Errorf("expected 200 OK for catalog, got %d", wCatalog.Code)
	}

	type CategoryResponse struct {
		Category *domain.StoreCategory  `json:"category"`
		Products []*domain.StoreProduct `json:"products"`
	}
	var catalog []CategoryResponse
	json.NewDecoder(wCatalog.Body).Decode(&catalog)

	if len(catalog) == 0 {
		t.Fatal("expected populated catalog")
	}

	// Verify categories are listed and find target test products
	var prodShoesID string
	var prodStaffID string
	for _, cat := range catalog {
		for _, p := range cat.Products {
			if p.ID == "prod-lowcut-shoes" {
				prodShoesID = p.ID
			}
			if p.ID == "prod-cudgel-staff" {
				prodStaffID = p.ID
			}
		}
	}
	if prodShoesID == "" || prodStaffID == "" {
		t.Fatal("lowcut shoes and cudgel staff target products not found in seeded db")
	}

	// --- 2. Test GetCart (Initially empty) ---
	reqCartEmpty := httptest.NewRequest("GET", "/api/cart", nil)
	reqCartEmpty.Header.Set("Authorization", "Bearer "+token)
	wCartEmpty := httptest.NewRecorder()
	mux.ServeHTTP(wCartEmpty, reqCartEmpty)

	if wCartEmpty.Code != http.StatusOK {
		t.Errorf("expected 200 OK for empty cart, got %d", wCartEmpty.Code)
	}
	var cartItems []domain.CartItem
	json.NewDecoder(wCartEmpty.Body).Decode(&cartItems)
	if len(cartItems) != 0 {
		t.Errorf("expected empty cart list, got %d items", len(cartItems))
	}

	// --- 3. Test AddToCart (Product 1: lowcut shoes) ---
	addReq1 := api.CartRequest{
		ProductID: &prodShoesID,
		Quantity:  1,
	}
	body1, _ := json.Marshal(addReq1)
	reqAdd1 := httptest.NewRequest("POST", "/api/cart", bytes.NewReader(body1))
	reqAdd1.Header.Set("Authorization", "Bearer "+token)
	wAdd1 := httptest.NewRecorder()
	mux.ServeHTTP(wAdd1, reqAdd1)

	if wAdd1.Code != http.StatusOK {
		t.Errorf("expected 200 OK for adding product 1 to cart, got %d", wAdd1.Code)
	}
	var addedItem domain.CartItem
	json.NewDecoder(wAdd1.Body).Decode(&addedItem)
	if addedItem.ID == "" || addedItem.ProductID == nil || *addedItem.ProductID != prodShoesID {
		t.Errorf("invalid cart item response: %+v", addedItem)
	}

	// --- 4. Test AddToCart (Product 2: staff) ---
	addReq2 := api.CartRequest{
		ProductID: &prodStaffID,
		Quantity:  1,
	}
	body2, _ := json.Marshal(addReq2)
	reqAdd2 := httptest.NewRequest("POST", "/api/cart", bytes.NewReader(body2))
	reqAdd2.Header.Set("Authorization", "Bearer "+token)
	wAdd2 := httptest.NewRecorder()
	mux.ServeHTTP(wAdd2, reqAdd2)

	if wAdd2.Code != http.StatusOK {
		t.Errorf("expected 200 OK for adding product 2 to cart, got %d", wAdd2.Code)
	}

	// --- 5. Verify GetCart contents ---
	reqCart2 := httptest.NewRequest("GET", "/api/cart", nil)
	reqCart2.Header.Set("Authorization", "Bearer "+token)
	wCart2 := httptest.NewRecorder()
	mux.ServeHTTP(wCart2, reqCart2)

	json.NewDecoder(wCart2.Body).Decode(&cartItems)
	if len(cartItems) != 2 {
		t.Errorf("expected 2 items in cart, got %d", len(cartItems))
	}

	// --- 6. Test RemoveFromCart ---
	// Delete first item (shoes)
	reqDel := httptest.NewRequest("DELETE", "/api/cart/"+addedItem.ID, nil)
	reqDel.Header.Set("Authorization", "Bearer "+token)
	wDel := httptest.NewRecorder()
	mux.ServeHTTP(wDel, reqDel)

	if wDel.Code != http.StatusNoContent {
		t.Errorf("expected 240 No Content for item deletion, got %d", wDel.Code)
	}

	// Verify cart now only contains staff (1 item)
	wCart3 := httptest.NewRecorder()
	mux.ServeHTTP(wCart3, reqCart2)
	json.NewDecoder(wCart3.Body).Decode(&cartItems)
	if len(cartItems) != 1 || cartItems[0].ProductID == nil || *cartItems[0].ProductID != prodStaffID {
		t.Errorf("expected cart to contain only staff, got: %+v", cartItems)
	}

	// Re-add shoes to test checkout math with multiple products
	reqAddRe := httptest.NewRequest("POST", "/api/cart", bytes.NewReader(body1))
	reqAddRe.Header.Set("Authorization", "Bearer "+token)
	wAddRe := httptest.NewRecorder()
	mux.ServeHTTP(wAddRe, reqAddRe)

	// Fetch initial stocks
	pShoesBefore, _ := storeRepo.GetProductByID(prodShoesID)
	pStaffBefore, _ := storeRepo.GetProductByID(prodStaffID)

	// --- 7. Test Checkout with Promo Code, CC Surcharges & Donations ---
	promoCode := "SUMMER2026"
	chkReq := api.CheckoutRequest{
		PromoCodeID:   &promoCode,
		DonationCents: 2000, // $20.00 donation
		PaymentMethod: "credit_card",
	}
	chkBody, _ := json.Marshal(chkReq)
	reqChk := httptest.NewRequest("POST", "/api/store/checkout", bytes.NewReader(chkBody))
	reqChk.Header.Set("Authorization", "Bearer "+token)
	wChk := httptest.NewRecorder()
	mux.ServeHTTP(wChk, reqChk)

	if wChk.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for checkout, got %d: %s", wChk.Code, wChk.Body.String())
	}

	var chkResp domain.Checkout
	json.NewDecoder(wChk.Body).Decode(&chkResp)

	// Validate Math:
	// Shoes ($25.00) + Staff ($25.00) = $50.00
	// Promo 'SUMMER2026' = 10% off -> $5.00 discount -> Pretax Subtotal: $45.00
	// Shipping Fee: $10.00 (under $100)
	// Donation: $20.00
	// Pre-tax Total: 4500 + 1000 + 2000 = $75.00 -> 7500 cents
	// Tax: 13% on products only ($45.00 * 0.13 = $5.85) -> 585 cents
	// Surcharge: 2.4% on CC order total (pre-tax + tax = 7500 + 585 = 8085 cents. 8085 * 0.024 = 194.04 cents) -> 194 cents
	// Total Order: 7500 + 585 + 194 = 8279 cents
	if chkResp.ItemsTotalCents != 5000 {
		t.Errorf("expected items total 5000, got %d", chkResp.ItemsTotalCents)
	}
	if chkResp.DiscountsTotalCents != 5000*10/100 {
		t.Errorf("expected discount 500, got %d", chkResp.DiscountsTotalCents)
	}
	if chkResp.ShippingFeeCents != 1000 {
		t.Errorf("expected shipping 1000, got %d", chkResp.ShippingFeeCents)
	}
	if chkResp.DonationCents != 2000 {
		t.Errorf("expected donation 2000, got %d", chkResp.DonationCents)
	}
	if chkResp.PreTaxCents != 7500 {
		t.Errorf("expected pre-tax cents 7500, got %d", chkResp.PreTaxCents)
	}
	if chkResp.TaxCents != 585 {
		t.Errorf("expected tax cents 585, got %d", chkResp.TaxCents)
	}
	if chkResp.CreditSurchargeCents != 194 {
		t.Errorf("expected credit surcharge 194, got %d", chkResp.CreditSurchargeCents)
	}
	if chkResp.OrderTotalCents != 8279 {
		t.Errorf("expected grand total 8279, got %d", chkResp.OrderTotalCents)
	}

	// --- 8. Verify Cart Cleared ---
	wCart4 := httptest.NewRecorder()
	mux.ServeHTTP(wCart4, reqCartEmpty)
	json.NewDecoder(wCart4.Body).Decode(&cartItems)
	if len(cartItems) != 0 {
		t.Errorf("expected cart to be cleared, got %d items", len(cartItems))
	}

	// --- 9. Verify Product Stocks Decremented ---
	pShoesAfter, _ := storeRepo.GetProductByID(prodShoesID)
	pStaffAfter, _ := storeRepo.GetProductByID(prodStaffID)
	if pShoesAfter.StockQuantity != pShoesBefore.StockQuantity-1 {
		t.Errorf("expected shoes stock to decrement to %d, got %d", pShoesBefore.StockQuantity-1, pShoesAfter.StockQuantity)
	}
	if pStaffAfter.StockQuantity != pStaffBefore.StockQuantity-1 {
		t.Errorf("expected staff stock to decrement to %d, got %d", pStaffBefore.StockQuantity-1, pStaffAfter.StockQuantity)
	}

	// --- 10. Verify Order History Retrieval ---
	reqHist := httptest.NewRequest("GET", "/api/store/orders/history", nil)
	reqHist.Header.Set("Authorization", "Bearer "+token)
	wHist := httptest.NewRecorder()
	mux.ServeHTTP(wHist, reqHist)

	if wHist.Code != http.StatusOK {
		t.Errorf("expected 200 OK for order history, got %d", wHist.Code)
	}
	var orderHistory []domain.Checkout
	json.NewDecoder(wHist.Body).Decode(&orderHistory)
	if len(orderHistory) != 1 || orderHistory[0].ID != chkResp.ID {
		t.Errorf("expected 1 history record matching checkout, got %+v", orderHistory)
	}
}
