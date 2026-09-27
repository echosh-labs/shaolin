package domain

import "time"

type StoreCategory struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	DisplayOrder int       `json:"display_order"`
	CreatedAt    time.Time `json:"created_at"`
}

type StoreProduct struct {
	ID            string    `json:"id"`
	CategoryID    string    `json:"category_id"`
	Name          string    `json:"name"`
	PriceCents    int       `json:"price_cents"`
	StockQuantity int       `json:"stock_quantity"`
	Description   *string   `json:"description,omitempty"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
}

type StoreOrder struct {
	ID            string    `json:"id"`
	UserID        *string   `json:"user_id,omitempty"`
	PreTaxCents   int       `json:"pre_tax_cents"`
	TaxCents      int       `json:"tax_cents"`
	TotalCents    int       `json:"total_cents"`
	PaymentStatus string    `json:"payment_status"` // 'pending', 'paid', 'refunded', 'cancelled'
	CheckoutID    *string   `json:"checkout_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

type StoreOrderItem struct {
	ID            string `json:"id"`
	OrderID       string `json:"order_id"`
	ProductID     string `json:"product_id"`
	Quantity      int    `json:"quantity"`
	PricePaidCents int    `json:"price_paid_cents"`
}

type PromoCode struct {
	ID                 string    `json:"id"`
	DiscountType       string    `json:"discount_type"` // 'percent', 'flat'
	DiscountValue      int       `json:"discount_value"`
	MinOrderValueCents int       `json:"min_order_value_cents"`
	IsActive           bool      `json:"is_active"`
	ExpiresAt          *string   `json:"expires_at,omitempty"` // YYYY-MM-DD
	CreatedAt          time.Time `json:"created_at"`
}

type Checkout struct {
	ID                   string    `json:"id"`
	UserID               string    `json:"user_id"`
	PromoCodeID          *string   `json:"promo_code_id,omitempty"`
	DonationCents        int       `json:"donation_cents"`
	PaymentMethod        string    `json:"payment_method"` // 'credit_card', 'e_transfer', 'cash', 'debit'
	ItemsTotalCents      int       `json:"items_total_cents"`
	DiscountsTotalCents  int       `json:"discounts_total_cents"`
	ShippingFeeCents     int       `json:"shipping_fee_cents"`
	PreTaxCents          int       `json:"pre_tax_cents"`
	TaxCents             int       `json:"tax_cents"`
	CreditSurchargeCents int       `json:"credit_surcharge_cents"`
	OrderTotalCents      int       `json:"order_total_cents"`
	PaymentStatus        string    `json:"payment_status"` // 'pending', 'paid', 'cancelled'
	OrderNumber          *int      `json:"order_number,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
}

type CartItem struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	ProductID      *string   `json:"product_id,omitempty"`
	RegistrationID *string   `json:"registration_id,omitempty"`
	Quantity       int       `json:"quantity"`
	CreatedAt      time.Time `json:"created_at"`
}

type StoreRepository interface {
	GetCategories() ([]*StoreCategory, error)
	GetProductsByCategoryID(categoryID string) ([]*StoreProduct, error)
	GetProductByID(id string) (*StoreProduct, error)
	UpdateProductStock(id string, quantity int) error

	CreateOrder(order *StoreOrder, items []*StoreOrderItem) error
	GetOrderByID(id string) (*StoreOrder, error)
	GetOrderItems(orderID string) ([]*StoreOrderItem, error)

	GetPromoCodeByID(id string) (*PromoCode, error)

	CreateCheckout(checkout *Checkout) error
	GetCheckoutByID(id string) (*Checkout, error)
	GetCheckoutsByUserID(userID string) ([]*Checkout, error)

	GetCartItemsByUserID(userID string) ([]*CartItem, error)
	SaveCartItem(item *CartItem) error
	DeleteCartItem(id string) error
	ClearCart(userID string) error
}
