package sqlite

import (
	"database/sql"

	"shaolin/backend/internal/domain"
)

type SQLiteStoreRepository struct {
	db *sql.DB
}

func NewSQLiteStoreRepository(db *sql.DB) *SQLiteStoreRepository {
	return &SQLiteStoreRepository{db: db}
}

func (r *SQLiteStoreRepository) GetCategories() ([]*domain.StoreCategory, error) {
	query := `SELECT id, name, display_order, created_at FROM store_categories ORDER BY display_order`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.StoreCategory
	for rows.Next() {
		var c domain.StoreCategory
		var createdAt string
		if err := rows.Scan(&c.ID, &c.Name, &c.DisplayOrder, &createdAt); err != nil {
			return nil, err
		}
		if t, err := parseTime(createdAt); err == nil {
			c.CreatedAt = t
		}
		list = append(list, &c)
	}
	return list, nil
}

func (r *SQLiteStoreRepository) GetProductsByCategoryID(categoryID string) ([]*domain.StoreProduct, error) {
	query := `SELECT id, category_id, name, price_cents, stock_quantity, description, is_active, created_at FROM store_products WHERE category_id = ? AND is_active = 1 ORDER BY name`
	rows, err := r.db.Query(query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.StoreProduct
	for rows.Next() {
		var p domain.StoreProduct
		var desc sql.NullString
		var isActive int
		var createdAt string
		if err := rows.Scan(&p.ID, &p.CategoryID, &p.Name, &p.PriceCents, &p.StockQuantity, &desc, &isActive, &createdAt); err != nil {
			return nil, err
		}
		if desc.Valid {
			p.Description = &desc.String
		}
		p.IsActive = (isActive == 1)
		if t, err := parseTime(createdAt); err == nil {
			p.CreatedAt = t
		}
		list = append(list, &p)
	}
	return list, nil
}

func (r *SQLiteStoreRepository) GetProductByID(id string) (*domain.StoreProduct, error) {
	query := `SELECT id, category_id, name, price_cents, stock_quantity, description, is_active, created_at FROM store_products WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var p domain.StoreProduct
	var desc sql.NullString
	var isActive int
	var createdAt string
	err := row.Scan(&p.ID, &p.CategoryID, &p.Name, &p.PriceCents, &p.StockQuantity, &desc, &isActive, &createdAt)
	if err != nil {
		return nil, err
	}
	if desc.Valid {
		p.Description = &desc.String
	}
	p.IsActive = (isActive == 1)
	if t, err := parseTime(createdAt); err == nil {
		p.CreatedAt = t
	}
	return &p, nil
}

func (r *SQLiteStoreRepository) UpdateProductStock(id string, quantity int) error {
	query := `UPDATE store_products SET stock_quantity = ? WHERE id = ?`
	_, err := r.db.Exec(query, quantity, id)
	return err
}

func (r *SQLiteStoreRepository) CreateOrder(order *domain.StoreOrder, items []*domain.StoreOrderItem) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	queryOrder := `INSERT INTO store_orders (id, user_id, pre_tax_cents, tax_cents, total_cents, payment_status, checkout_id, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = tx.Exec(queryOrder, order.ID, order.UserID, order.PreTaxCents, order.TaxCents, order.TotalCents, order.PaymentStatus, order.CheckoutID, order.CreatedAt)
	if err != nil {
		return err
	}

	queryItem := `INSERT INTO store_order_items (id, order_id, product_id, quantity, price_paid_cents) VALUES (?, ?, ?, ?, ?)`
	for _, item := range items {
		_, err = tx.Exec(queryItem, item.ID, item.OrderID, item.ProductID, item.Quantity, item.PricePaidCents)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *SQLiteStoreRepository) GetOrderByID(id string) (*domain.StoreOrder, error) {
	query := `SELECT id, user_id, pre_tax_cents, tax_cents, total_cents, payment_status, checkout_id, created_at FROM store_orders WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var o domain.StoreOrder
	var checkoutID sql.NullString
	var createdAt string
	err := row.Scan(&o.ID, &o.UserID, &o.PreTaxCents, &o.TaxCents, &o.TotalCents, &o.PaymentStatus, &checkoutID, &createdAt)
	if err != nil {
		return nil, err
	}
	if checkoutID.Valid {
		o.CheckoutID = &checkoutID.String
	}
	if t, err := parseTime(createdAt); err == nil {
		o.CreatedAt = t
	}
	return &o, nil
}

func (r *SQLiteStoreRepository) GetOrderItems(orderID string) ([]*domain.StoreOrderItem, error) {
	query := `SELECT id, order_id, product_id, quantity, price_paid_cents FROM store_order_items WHERE order_id = ?`
	rows, err := r.db.Query(query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.StoreOrderItem
	for rows.Next() {
		var item domain.StoreOrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.PricePaidCents); err != nil {
			return nil, err
		}
		list = append(list, &item)
	}
	return list, nil
}

func (r *SQLiteStoreRepository) GetPromoCodeByID(id string) (*domain.PromoCode, error) {
	query := `SELECT id, discount_type, discount_value, min_order_value_cents, is_active, expires_at, created_at FROM promo_codes WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var p domain.PromoCode
	var expiresAt sql.NullString
	var isActive int
	var createdAt string
	err := row.Scan(&p.ID, &p.DiscountType, &p.DiscountValue, &p.MinOrderValueCents, &isActive, &expiresAt, &createdAt)
	if err != nil {
		return nil, err
	}
	p.IsActive = (isActive == 1)
	if expiresAt.Valid {
		p.ExpiresAt = &expiresAt.String
	}
	if t, err := parseTime(createdAt); err == nil {
		p.CreatedAt = t
	}
	return &p, nil
}

func (r *SQLiteStoreRepository) CreateCheckout(c *domain.Checkout) error {
	query := `INSERT INTO checkouts (id, user_id, promo_code_id, donation_cents, payment_method, items_total_cents, discounts_total_cents, shipping_fee_cents, pre_tax_cents, tax_cents, credit_surcharge_cents, order_total_cents, payment_status, order_number, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, c.ID, c.UserID, c.PromoCodeID, c.DonationCents, c.PaymentMethod, c.ItemsTotalCents, c.DiscountsTotalCents, c.ShippingFeeCents, c.PreTaxCents, c.TaxCents, c.CreditSurchargeCents, c.OrderTotalCents, c.PaymentStatus, c.OrderNumber, c.CreatedAt)
	return err
}

func (r *SQLiteStoreRepository) GetCheckoutByID(id string) (*domain.Checkout, error) {
	query := `SELECT id, user_id, promo_code_id, donation_cents, payment_method, items_total_cents, discounts_total_cents, shipping_fee_cents, pre_tax_cents, tax_cents, credit_surcharge_cents, order_total_cents, payment_status, order_number, created_at FROM checkouts WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var c domain.Checkout
	var promoCodeID sql.NullString
	var orderNumber sql.NullInt64
	var createdAt string
	err := row.Scan(&c.ID, &c.UserID, &promoCodeID, &c.DonationCents, &c.PaymentMethod, &c.ItemsTotalCents, &c.DiscountsTotalCents, &c.ShippingFeeCents, &c.PreTaxCents, &c.TaxCents, &c.CreditSurchargeCents, &c.OrderTotalCents, &c.PaymentStatus, &orderNumber, &createdAt)
	if err != nil {
		return nil, err
	}
	if promoCodeID.Valid {
		c.PromoCodeID = &promoCodeID.String
	}
	if orderNumber.Valid {
		val := int(orderNumber.Int64)
		c.OrderNumber = &val
	}
	if t, err := parseTime(createdAt); err == nil {
		c.CreatedAt = t
	}
	return &c, nil
}

func (r *SQLiteStoreRepository) GetCheckoutsByUserID(userID string) ([]*domain.Checkout, error) {
	query := `SELECT id, user_id, promo_code_id, donation_cents, payment_method, items_total_cents, discounts_total_cents, shipping_fee_cents, pre_tax_cents, tax_cents, credit_surcharge_cents, order_total_cents, payment_status, order_number, created_at FROM checkouts WHERE user_id = ? ORDER BY created_at DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.Checkout
	for rows.Next() {
		var c domain.Checkout
		var promoCodeID sql.NullString
		var orderNumber sql.NullInt64
		var createdAt string
		err := rows.Scan(&c.ID, &c.UserID, &promoCodeID, &c.DonationCents, &c.PaymentMethod, &c.ItemsTotalCents, &c.DiscountsTotalCents, &c.ShippingFeeCents, &c.PreTaxCents, &c.TaxCents, &c.CreditSurchargeCents, &c.OrderTotalCents, &c.PaymentStatus, &orderNumber, &createdAt)
		if err != nil {
			return nil, err
		}
		if promoCodeID.Valid {
			c.PromoCodeID = &promoCodeID.String
		}
		if orderNumber.Valid {
			val := int(orderNumber.Int64)
			c.OrderNumber = &val
		}
		if t, err := parseTime(createdAt); err == nil {
			c.CreatedAt = t
		}
		list = append(list, &c)
	}
	return list, nil
}

func (r *SQLiteStoreRepository) GetCartItemsByUserID(userID string) ([]*domain.CartItem, error) {
	query := `SELECT id, user_id, product_id, registration_id, quantity, created_at FROM cart_items WHERE user_id = ? ORDER BY created_at ASC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.CartItem
	for rows.Next() {
		var ci domain.CartItem
		var productID, registrationID sql.NullString
		var createdAt string
		if err := rows.Scan(&ci.ID, &ci.UserID, &productID, &registrationID, &ci.Quantity, &createdAt); err != nil {
			return nil, err
		}
		if productID.Valid {
			ci.ProductID = &productID.String
		}
		if registrationID.Valid {
			ci.RegistrationID = &registrationID.String
		}
		if t, err := parseTime(createdAt); err == nil {
			ci.CreatedAt = t
		}
		list = append(list, &ci)
	}
	return list, nil
}

func (r *SQLiteStoreRepository) SaveCartItem(item *domain.CartItem) error {
	query := `INSERT OR REPLACE INTO cart_items (id, user_id, product_id, registration_id, quantity, created_at) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, item.ID, item.UserID, item.ProductID, item.RegistrationID, item.Quantity, item.CreatedAt)
	return err
}

func (r *SQLiteStoreRepository) DeleteCartItem(id string) error {
	query := `DELETE FROM cart_items WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *SQLiteStoreRepository) ClearCart(userID string) error {
	query := `DELETE FROM cart_items WHERE user_id = ?`
	_, err := r.db.Exec(query, userID)
	return err
}
