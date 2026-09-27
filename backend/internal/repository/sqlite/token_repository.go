package sqlite

import (
	"database/sql"
	"time"

	"shaolin/backend/internal/domain"
)

type SQLiteTokenRepository struct {
	db *sql.DB
}

func NewSQLiteTokenRepository(db *sql.DB) *SQLiteTokenRepository {
	return &SQLiteTokenRepository{db: db}
}

func (r *SQLiteTokenRepository) CreateTransaction(tx *domain.TokenTransaction) error {
	query := `INSERT INTO token_transactions (id, user_id, term_id, tokens_added, price_paid_cents, transaction_type, checkout_id, created_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	createdAt := tx.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now()
	}
	_, err := r.db.Exec(query, tx.ID, tx.UserID, tx.TermID, tx.TokensAdded, tx.PricePaidCents, tx.TransactionType, tx.CheckoutID, createdAt)
	return err
}

func (r *SQLiteTokenRepository) GetBalance(userID, termID string) (*domain.UserTermTokens, error) {
	query := `SELECT user_id, term_id, tokens_remaining FROM user_term_tokens WHERE user_id = ? AND term_id = ?`
	row := r.db.QueryRow(query, userID, termID)

	var bal domain.UserTermTokens
	err := row.Scan(&bal.UserID, &bal.TermID, &bal.TokensRemaining)
	if err != nil {
		return nil, err
	}
	return &bal, nil
}

func (r *SQLiteTokenRepository) UpdateBalance(userID, termID string, remaining int) error {
	query := `INSERT INTO user_term_tokens (user_id, term_id, tokens_remaining) VALUES (?, ?, ?)
	          ON CONFLICT(user_id, term_id) DO UPDATE SET tokens_remaining = excluded.tokens_remaining`
	_, err := r.db.Exec(query, userID, termID, remaining)
	return err
}

func (r *SQLiteTokenRepository) GetTransactionsByUserID(userID string) ([]*domain.TokenTransaction, error) {
	query := `SELECT id, user_id, term_id, tokens_added, price_paid_cents, transaction_type, checkout_id, created_at 
	          FROM token_transactions WHERE user_id = ? ORDER BY created_at DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.TokenTransaction
	for rows.Next() {
		var tx domain.TokenTransaction
		var createdAt string
		if err := rows.Scan(&tx.ID, &tx.UserID, &tx.TermID, &tx.TokensAdded, &tx.PricePaidCents, &tx.TransactionType, &tx.CheckoutID, &createdAt); err != nil {
			return nil, err
		}
		t, err := parseTime(createdAt)
		if err == nil {
			tx.CreatedAt = t
		}
		list = append(list, &tx)
	}
	return list, nil
}

func (r *SQLiteTokenRepository) GetTransactionsByTermID(termID string) ([]*domain.TokenTransaction, error) {
	query := `SELECT id, user_id, term_id, tokens_added, price_paid_cents, transaction_type, checkout_id, created_at 
	          FROM token_transactions WHERE term_id = ? ORDER BY created_at DESC`
	rows, err := r.db.Query(query, termID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.TokenTransaction
	for rows.Next() {
		var tx domain.TokenTransaction
		var createdAt string
		if err := rows.Scan(&tx.ID, &tx.UserID, &tx.TermID, &tx.TokensAdded, &tx.PricePaidCents, &tx.TransactionType, &tx.CheckoutID, &createdAt); err != nil {
			return nil, err
		}
		t, err := parseTime(createdAt)
		if err == nil {
			tx.CreatedAt = t
		}
		list = append(list, &tx)
	}
	return list, nil
}

func (r *SQLiteTokenRepository) GetPackages() ([]*domain.TokenPackage, error) {
	query := `SELECT id, tokens_count, terms_duration, adult_price_cents, child_price_cents, senior_price_cents, is_active FROM token_packages ORDER BY tokens_count ASC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.TokenPackage
	for rows.Next() {
		var p domain.TokenPackage
		var isActive int
		if err := rows.Scan(&p.ID, &p.TokensCount, &p.TermsDuration, &p.AdultPriceCents, &p.ChildPriceCents, &p.SeniorPriceCents, &isActive); err != nil {
			return nil, err
		}
		p.IsActive = (isActive == 1)
		list = append(list, &p)
	}
	return list, nil
}

func (r *SQLiteTokenRepository) GetPackageByID(id string) (*domain.TokenPackage, error) {
	query := `SELECT id, tokens_count, terms_duration, adult_price_cents, child_price_cents, senior_price_cents, is_active FROM token_packages WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var p domain.TokenPackage
	var isActive int
	err := row.Scan(&p.ID, &p.TokensCount, &p.TermsDuration, &p.AdultPriceCents, &p.ChildPriceCents, &p.SeniorPriceCents, &isActive)
	if err != nil {
		return nil, err
	}
	p.IsActive = (isActive == 1)
	return &p, nil
}

func (r *SQLiteTokenRepository) GetPingPongPackages() ([]*domain.PingPongPackage, error) {
	query := `SELECT id, tokens_count, price_cents, is_active FROM ping_pong_packages ORDER BY tokens_count ASC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.PingPongPackage
	for rows.Next() {
		var p domain.PingPongPackage
		var isActive int
		if err := rows.Scan(&p.ID, &p.TokensCount, &p.PriceCents, &isActive); err != nil {
			return nil, err
		}
		p.IsActive = (isActive == 1)
		list = append(list, &p)
	}
	return list, nil
}
