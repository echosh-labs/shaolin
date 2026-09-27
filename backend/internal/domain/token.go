package domain

import "time"

type TokenTransaction struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	TermID          string    `json:"term_id"`
	TokensAdded     int       `json:"tokens_added"`
	PricePaidCents  int       `json:"price_paid_cents"`
	TransactionType string    `json:"transaction_type"` // 'purchase', 'refund', 'admin_adjustment', 'referral_bonus'
	CheckoutID      *string   `json:"checkout_id,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

type UserTermTokens struct {
	UserID          string `json:"user_id"`
	TermID          string `json:"term_id"`
	TokensRemaining int    `json:"tokens_remaining"`
}

type TokenPackage struct {
	ID             string `json:"id"`
	TokensCount    int    `json:"tokens_count"`
	TermsDuration  int    `json:"terms_duration"`
	AdultPriceCents int    `json:"adult_price_cents"`
	ChildPriceCents int    `json:"child_price_cents"`
	SeniorPriceCents int   `json:"senior_price_cents"`
	IsActive       bool   `json:"is_active"`
}

type PingPongPackage struct {
	ID          string `json:"id"`
	TokensCount int    `json:"tokens_count"`
	PriceCents  int    `json:"price_cents"`
	IsActive    bool   `json:"is_active"`
}

type TokenRepository interface {
	CreateTransaction(tx *TokenTransaction) error
	GetBalance(userID, termID string) (*UserTermTokens, error)
	UpdateBalance(userID, termID string, remaining int) error
	GetTransactionsByUserID(userID string) ([]*TokenTransaction, error)
	GetTransactionsByTermID(termID string) ([]*TokenTransaction, error)

	GetPackages() ([]*TokenPackage, error)
	GetPackageByID(id string) (*TokenPackage, error)
	GetPingPongPackages() ([]*PingPongPackage, error)
}
