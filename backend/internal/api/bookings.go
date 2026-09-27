package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"shaolin/backend/internal/domain"
)

type BookingHandler struct {
	bookingRepo domain.BookingRepository
	tokenRepo   domain.TokenRepository
}

func NewBookingHandler(bookingRepo domain.BookingRepository, tokenRepo domain.TokenRepository) *BookingHandler {
	return &BookingHandler{
		bookingRepo: bookingRepo,
		tokenRepo:   tokenRepo,
	}
}

type PurchaseRequest struct {
	PackageID string `json:"package_id"`
	TermID    string `json:"term_id"`
}

type BookingRequest struct {
	OccurrenceID   string `json:"occurrence_id"`
	AttendanceMode string `json:"attendance_mode"` // 'in_person', 'live_stream'
}

func (h *BookingHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	claims, err := GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	termID := r.URL.Query().Get("term_id")
	if termID == "" {
		termID = "term-summer-2026"
	}

	bal, err := h.tokenRepo.GetBalance(claims.UserID, termID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Return zero balance if no ledger entry exists
			bal = &domain.UserTermTokens{
				UserID:          claims.UserID,
				TermID:          termID,
				TokensRemaining: 0,
			}
		} else {
			log.Printf("Failed to get balance: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bal)
}

func (h *BookingHandler) ListPackages(w http.ResponseWriter, r *http.Request) {
	list, err := h.tokenRepo.GetPackages()
	if err != nil {
		log.Printf("Failed to list token packages: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(list)
}

func (h *BookingHandler) ListPingPongPackages(w http.ResponseWriter, r *http.Request) {
	list, err := h.tokenRepo.GetPingPongPackages()
	if err != nil {
		log.Printf("Failed to list ping pong packages: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(list)
}

func (h *BookingHandler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	claims, err := GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	list, err := h.tokenRepo.GetTransactionsByUserID(claims.UserID)
	if err != nil {
		log.Printf("Failed to list token transactions: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(list)
}

func (h *BookingHandler) PurchaseTokens(w http.ResponseWriter, r *http.Request) {
	claims, err := GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req PurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.TermID == "" {
		req.TermID = "term-summer-2026"
	}

	pkg, err := h.tokenRepo.GetPackageByID(req.PackageID)
	if err != nil {
		http.Error(w, "Token package not found", http.StatusNotFound)
		return
	}

	// Calculate price (we assume adult rate for generic API sim, or check DOB/age)
	pricePaid := pkg.AdultPriceCents

	// 1. Log transaction
	tx := &domain.TokenTransaction{
		ID:              uuid.New().String(),
		UserID:          claims.UserID,
		TermID:          req.TermID,
		TokensAdded:     pkg.TokensCount,
		PricePaidCents:  pricePaid,
		TransactionType: "purchase",
		CreatedAt:       time.Now(),
	}

	if err := h.tokenRepo.CreateTransaction(tx); err != nil {
		log.Printf("Failed to create token transaction: %v", err)
		http.Error(w, "Transaction failed", http.StatusInternalServerError)
		return
	}

	// 2. Update user's token balance
	bal, err := h.tokenRepo.GetBalance(claims.UserID, req.TermID)
	tokensCountBefore := 0
	if err == nil && bal != nil {
		tokensCountBefore = bal.TokensRemaining
	}

	tokensRemainingNew := tokensCountBefore + pkg.TokensCount
	if err := h.tokenRepo.UpdateBalance(claims.UserID, req.TermID, tokensRemainingNew); err != nil {
		log.Printf("Failed to update tokens balance: %v", err)
		http.Error(w, "Failed to update balance", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"transaction_id":   tx.ID,
		"tokens_added":     pkg.TokensCount,
		"tokens_remaining": tokensRemainingNew,
	})
}

func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	claims, err := GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req BookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.OccurrenceID == "" || (req.AttendanceMode != "in_person" && req.AttendanceMode != "live_stream") {
		http.Error(w, "OccurrenceID is required and AttendanceMode must be 'in_person' or 'live_stream'", http.StatusBadRequest)
		return
	}

	// 1. Get class occurrence details
	occ, err := h.bookingRepo.GetOccurrenceByID(req.OccurrenceID)
	if err != nil {
		http.Error(w, "Class occurrence not found", http.StatusNotFound)
		return
	}
	if occ.Status == "cancelled" {
		http.Error(w, "Cannot book a cancelled class occurrence", http.StatusBadRequest)
		return
	}

	// 2. Get class template to verify capacity and term info
	class, err := h.bookingRepo.GetClassByID(occ.ClassID)
	if err != nil {
		http.Error(w, "Base class config not found", http.StatusNotFound)
		return
	}

	// Only enforce capacity limits for in_person attendance
	if req.AttendanceMode == "in_person" && occ.BookedCount >= class.Capacity {
		http.Error(w, "Class occurrence is full", http.StatusConflict)
		return
	}

	// 3. Check user's token balance for this term
	bal, err := h.tokenRepo.GetBalance(claims.UserID, class.TermID)
	if err != nil || bal == nil || bal.TokensRemaining < 1 {
		http.Error(w, "Insufficient tokens to book this class", http.StatusBadRequest)
		return
	}

	// 4. Perform booking operations (deduct token, increment booked_count, create booking)
	booking := &domain.Booking{
		ID:             uuid.New().String(),
		UserID:         claims.UserID,
		OccurrenceID:   req.OccurrenceID,
		Status:         "confirmed",
		AttendanceMode: req.AttendanceMode,
		BookedAt:       time.Now(),
	}

	// Start database transaction simulation (we perform queries sequentially)
	// Deduct token
	if err := h.tokenRepo.UpdateBalance(claims.UserID, class.TermID, bal.TokensRemaining-1); err != nil {
		log.Printf("Failed to deduct token: %v", err)
		http.Error(w, "Booking failed", http.StatusInternalServerError)
		return
	}

	// Increment booked count
	occ.BookedCount++
	if err := h.bookingRepo.UpdateOccurrence(occ); err != nil {
		// Rollback balance (simple manual query rollback)
		h.tokenRepo.UpdateBalance(claims.UserID, class.TermID, bal.TokensRemaining)
		log.Printf("Failed to increment occurrence booking count: %v", err)
		http.Error(w, "Booking failed", http.StatusInternalServerError)
		return
	}

	// Write booking record
	if err := h.bookingRepo.CreateBooking(booking); err != nil {
		// Rollback balance & booked count
		h.tokenRepo.UpdateBalance(claims.UserID, class.TermID, bal.TokensRemaining)
		occ.BookedCount--
		h.bookingRepo.UpdateOccurrence(occ)
		log.Printf("Failed to save booking: %v", err)
		http.Error(w, "Booking failed", http.StatusInternalServerError)
		return
	}

	// Write negative token transaction record for audit history
	h.tokenRepo.CreateTransaction(&domain.TokenTransaction{
		ID:              uuid.New().String(),
		UserID:          claims.UserID,
		TermID:          class.TermID,
		TokensAdded:     -1,
		PricePaidCents:  0,
		TransactionType: "class_booking",
		CreatedAt:       time.Now(),
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(booking)
}

func (h *BookingHandler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	claims, err := GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Booking ID is required", http.StatusBadRequest)
		return
	}

	// 1. Get booking
	booking, err := h.bookingRepo.GetBookingByID(id)
	if err != nil {
		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	}
	if booking.UserID != claims.UserID {
		http.Error(w, "Forbidden: you do not own this booking", http.StatusForbidden)
		return
	}
	if booking.Status == "cancelled" {
		http.Error(w, "Booking is already cancelled", http.StatusBadRequest)
		return
	}

	// 2. Get occurrence details
	occ, err := h.bookingRepo.GetOccurrenceByID(booking.OccurrenceID)
	if err != nil {
		http.Error(w, "Linked class occurrence not found", http.StatusNotFound)
		return
	}

	// 3. Determine if refund is eligible (must be cancelled prior to occurrence date/time)
	// (Check if cancellation is before occurrence date)
	nowStr := time.Now().Format("2006-01-02")
	refundEligible := occ.Date >= nowStr

	// 4. Update stats: decrement booked count, update booking, and refund tokens if eligible
	booking.Status = "cancelled"
	now := time.Now()
	booking.CancelledAt = &now

	if err := h.bookingRepo.UpdateBooking(booking); err != nil {
		log.Printf("Failed to cancel booking: %v", err)
		http.Error(w, "Cancellation failed", http.StatusInternalServerError)
		return
	}

	if occ.BookedCount > 0 {
		occ.BookedCount--
		h.bookingRepo.UpdateOccurrence(occ)
	}

	refunded := false
	if refundEligible {
		// Find term from class
		class, err := h.bookingRepo.GetClassByID(occ.ClassID)
		if err == nil && class != nil {
			bal, err := h.tokenRepo.GetBalance(claims.UserID, class.TermID)
			if err == nil && bal != nil {
				h.tokenRepo.UpdateBalance(claims.UserID, class.TermID, bal.TokensRemaining+1)
				h.tokenRepo.CreateTransaction(&domain.TokenTransaction{
					ID:              uuid.New().String(),
					UserID:          claims.UserID,
					TermID:          class.TermID,
					TokensAdded:     1,
					PricePaidCents:  0,
					TransactionType: "refund",
					CreatedAt:       time.Now(),
				})
				refunded = true
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"booking_id":      booking.ID,
		"status":          "cancelled",
		"refund_eligible": refundEligible,
		"tokens_refunded": refunded,
	})
}
