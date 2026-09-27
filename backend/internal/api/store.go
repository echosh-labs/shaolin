package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"shaolin/backend/internal/domain"
)

type StoreHandler struct {
	repo domain.StoreRepository
}

func NewStoreHandler(repo domain.StoreRepository) *StoreHandler {
	return &StoreHandler{repo: repo}
}

type CheckoutRequest struct {
	PromoCodeID   *string `json:"promo_code_id,omitempty"`
	DonationCents int     `json:"donation_cents"`
	PaymentMethod string  `json:"payment_method"` // 'credit_card', 'e_transfer', 'cash', 'debit'
}

type CartRequest struct {
	ProductID      *string `json:"product_id,omitempty"`
	RegistrationID *string `json:"registration_id,omitempty"`
	Quantity       int     `json:"quantity"`
}

// ListProducts returns the catalog of products grouped by category or filtered by category_id
func (h *StoreHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	catID := r.URL.Query().Get("category_id")
	if catID != "" {
		list, err := h.repo.GetProductsByCategoryID(catID)
		if err != nil {
			log.Printf("Failed to get products by category: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)
		return
	}

	categories, err := h.repo.GetCategories()
	if err != nil {
		log.Printf("Failed to list categories: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	type CategoryResponse struct {
		Category *domain.StoreCategory  `json:"category"`
		Products []*domain.StoreProduct `json:"products"`
	}

	var resp []CategoryResponse
	for _, cat := range categories {
		products, err := h.repo.GetProductsByCategoryID(cat.ID)
		if err != nil {
			log.Printf("Failed to get products for category %s: %v", cat.ID, err)
			continue
		}
		resp = append(resp, CategoryResponse{
			Category: cat,
			Products: products,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// GetCart retrieves active shopping cart items for the user
func (h *StoreHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	claims, err := GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	items, err := h.repo.GetCartItemsByUserID(claims.UserID)
	if err != nil {
		log.Printf("Failed to get cart items: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(items)
}

// AddToCart adds a product or term registration to the user's shopping cart
func (h *StoreHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	claims, err := GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if (req.ProductID == nil && req.RegistrationID == nil) || (req.ProductID != nil && req.RegistrationID != nil) {
		http.Error(w, "Exactly one of ProductID or RegistrationID must be provided", http.StatusBadRequest)
		return
	}

	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	// Verify product exists and is active/has stock if ProductID is provided
	if req.ProductID != nil {
		p, err := h.repo.GetProductByID(*req.ProductID)
		if err != nil {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		if !p.IsActive {
			http.Error(w, "Product is inactive", http.StatusBadRequest)
			return
		}
		if p.StockQuantity < req.Quantity {
			http.Error(w, "Insufficient product stock availability", http.StatusConflict)
			return
		}
	}

	// Create unique ID derived from user/product or random uuid
	cartItemID := uuid.New().String()
	// Let's check if the user already has this item in their cart to increment quantity
	existingCartItems, err := h.repo.GetCartItemsByUserID(claims.UserID)
	if err == nil {
		for _, ci := range existingCartItems {
			if req.ProductID != nil && ci.ProductID != nil && *ci.ProductID == *req.ProductID {
				cartItemID = ci.ID
				req.Quantity += ci.Quantity
				break
			}
			if req.RegistrationID != nil && ci.RegistrationID != nil && *ci.RegistrationID == *req.RegistrationID {
				cartItemID = ci.ID
				req.Quantity += ci.Quantity
				break
			}
		}
	}

	item := &domain.CartItem{
		ID:             cartItemID,
		UserID:         claims.UserID,
		ProductID:      req.ProductID,
		RegistrationID: req.RegistrationID,
		Quantity:       req.Quantity,
		CreatedAt:      time.Now(),
	}

	if err := h.repo.SaveCartItem(item); err != nil {
		log.Printf("Failed to save cart item: %v", err)
		http.Error(w, "Failed to save cart item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
}

// RemoveFromCart deletes a specific cart item by ID
func (h *StoreHandler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Cart item ID is required", http.StatusBadRequest)
		return
	}

	if err := h.repo.DeleteCartItem(id); err != nil {
		log.Printf("Failed to delete cart item: %v", err)
		http.Error(w, "Failed to delete cart item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Checkout checks out the active shopping cart items, applying taxes and surcharges
func (h *StoreHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	claims, err := GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// 1. Fetch active cart items
	cartItems, err := h.repo.GetCartItemsByUserID(claims.UserID)
	if err != nil || len(cartItems) == 0 {
		http.Error(w, "Shopping cart is empty", http.StatusBadRequest)
		return
	}

	// 2. Validate products and calculate items subtotal
	itemsSubtotal := 0
	var orderItems []*domain.StoreOrderItem
	type stockUpdate struct {
		productID string
		newStock  int
	}
	var stocksToUpdate []stockUpdate

	for _, ci := range cartItems {
		if ci.ProductID != nil {
			p, err := h.repo.GetProductByID(*ci.ProductID)
			if err != nil {
				http.Error(w, "Product not found: "+*ci.ProductID, http.StatusNotFound)
				return
			}
			if p.StockQuantity < ci.Quantity {
				http.Error(w, "Product "+p.Name+" has insufficient stock", http.StatusConflict)
				return
			}
			itemsSubtotal += p.PriceCents * ci.Quantity
			orderItems = append(orderItems, &domain.StoreOrderItem{
				ID:             uuid.New().String(),
				ProductID:      p.ID,
				Quantity:       ci.Quantity,
				PricePaidCents: p.PriceCents,
			})
			stocksToUpdate = append(stocksToUpdate, stockUpdate{
				productID: p.ID,
				newStock:  p.StockQuantity - ci.Quantity,
			})
		}
		// Registration item pricing can be handled separately if integrated.
		// For general store inventory cart checkouts, we calculate products subtotal.
	}

	// 3. Apply Promo Code
	discountCents := 0
	if req.PromoCodeID != nil && *req.PromoCodeID != "" {
		promo, err := h.repo.GetPromoCodeByID(*req.PromoCodeID)
		if err == nil && promo.IsActive {
			if itemsSubtotal >= promo.MinOrderValueCents {
				if promo.DiscountType == "percent" {
					discountCents = (itemsSubtotal * promo.DiscountValue) / 100
				} else if promo.DiscountType == "flat" {
					discountCents = promo.DiscountValue
				}
				if discountCents > itemsSubtotal {
					discountCents = itemsSubtotal
				}
			}
		}
	}

	// 4. Calculate Tax (13% HST on taxable subtotal: itemsSubtotal - discountCents)
	taxableSubtotal := itemsSubtotal - discountCents
	if taxableSubtotal < 0 {
		taxableSubtotal = 0
	}
	taxCents := (taxableSubtotal * 13) / 100

	// 5. Pre-tax cents (subtotal - discount + donation + shipping)
	// (Note: donation and shipping are non-taxable on top of products)
	shippingFeeCents := 0
	if taxableSubtotal > 0 && taxableSubtotal < 10000 { // Under $100.00 pays $10 shipping
		shippingFeeCents = 1000
	}
	preTaxCents := taxableSubtotal + req.DonationCents + shippingFeeCents

	// 6. Credit Card Surcharge (2.4% on items total + tax)
	surchargeCents := 0
	if req.PaymentMethod == "credit_card" {
		surchargeCents = ((preTaxCents + taxCents) * 24) / 1000
	}

	grandTotalCents := preTaxCents + taxCents + surchargeCents

	// 7. Write Checkout Record
	checkoutID := uuid.New().String()
	orderNum := int(time.Now().UnixNano() % 1000000)
	checkout := &domain.Checkout{
		ID:                   checkoutID,
		UserID:               claims.UserID,
		PromoCodeID:          req.PromoCodeID,
		DonationCents:        req.DonationCents,
		PaymentMethod:        req.PaymentMethod,
		ItemsTotalCents:      itemsSubtotal,
		DiscountsTotalCents:  discountCents,
		ShippingFeeCents:     shippingFeeCents,
		PreTaxCents:          preTaxCents,
		TaxCents:             taxCents,
		CreditSurchargeCents: surchargeCents,
		OrderTotalCents:      grandTotalCents,
		PaymentStatus:        "paid", // simulated checkout resolves instantly to paid
		OrderNumber:          &orderNum,
		CreatedAt:            time.Now(),
	}

	if err := h.repo.CreateCheckout(checkout); err != nil {
		log.Printf("Failed to create checkout record: %v", err)
		http.Error(w, "Checkout failed", http.StatusInternalServerError)
		return
	}

	// 8. Create Store Order and Order Items
	if len(orderItems) > 0 {
		orderID := uuid.New().String()
		order := &domain.StoreOrder{
			ID:            orderID,
			UserID:        &claims.UserID,
			PreTaxCents:   taxableSubtotal,
			TaxCents:      taxCents,
			TotalCents:    taxableSubtotal + taxCents,
			PaymentStatus: "paid",
			CheckoutID:    &checkoutID,
			CreatedAt:     time.Now(),
		}

		for _, item := range orderItems {
			item.OrderID = orderID
		}

		if err := h.repo.CreateOrder(order, orderItems); err != nil {
			log.Printf("Failed to create store order: %v", err)
			// Non-blocking for checkout itself but logged
		} else {
			// Update stocks on success
			for _, su := range stocksToUpdate {
				h.repo.UpdateProductStock(su.productID, su.newStock)
			}
		}
	}

	// 9. Clear the user's cart
	h.repo.ClearCart(claims.UserID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(checkout)
}

// GetOrderHistory returns all previous checkouts/invoices for the student
func (h *StoreHandler) GetOrderHistory(w http.ResponseWriter, r *http.Request) {
	claims, err := GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	list, err := h.repo.GetCheckoutsByUserID(claims.UserID)
	if err != nil {
		log.Printf("Failed to fetch order history: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(list)
}
