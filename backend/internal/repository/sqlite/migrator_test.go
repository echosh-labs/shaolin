package sqlite_test

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
	"shaolin/backend/internal/repository/sqlite"
)

func TestMigrate(t *testing.T) {
	// Open in-memory SQLite database
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open sqlite database: %v", err)
	}
	defer db.Close()

	// Run migrations
	err = sqlite.Migrate(db)
	if err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	// Verify schema_migrations table contains version 1
	var version int
	err = db.QueryRow("SELECT version FROM schema_migrations WHERE version = 1").Scan(&version)
	if err != nil {
		t.Errorf("failed to find version 1 in migrations table: %v", err)
	}

	// Test user insertion
	_, err = db.Exec(`
		INSERT INTO users (id, email, password_hash, first_name, last_name, date_of_birth, role)
		VALUES ('u-123', 'test@martialartsacademy.com', 'passwordhash123', 'Marcus', 'Vance', '1980-05-15', 'admin')
	`)
	if err != nil {
		t.Errorf("failed to insert user: %v", err)
	}

	// Verify email uniqueness constraint
	_, err = db.Exec(`
		INSERT INTO users (id, email, password_hash, first_name, last_name, date_of_birth, role)
		VALUES ('u-456', 'test@martialartsacademy.com', 'anotherhash', 'Student', 'Mulan', '1995-10-10', 'student')
	`)
	if err == nil {
		t.Error("expected duplicate email insertion to fail, but it succeeded")
	}

	// Test term insertion
	_, err = db.Exec(`
		INSERT INTO terms (id, name, start_date, end_date, is_active)
		VALUES ('t-summer-2026', 'Summer 2026', '2026-06-01', '2026-08-31', 1)
	`)
	if err != nil {
		t.Errorf("failed to insert term: %v", err)
	}

	// Test user_term_tokens foreign key constraints
	_, err = db.Exec(`
		INSERT INTO user_term_tokens (user_id, term_id, tokens_remaining)
		VALUES ('u-123', 't-summer-2026', 14)
	`)
	if err != nil {
		t.Errorf("failed to insert user term tokens: %v", err)
	}

	// Test user format check constraint (invalid date_of_birth)
	_, err = db.Exec(`
		INSERT INTO users (id, email, password_hash, first_name, last_name, date_of_birth, role)
		VALUES ('u-bad-dob', 'bad-dob@martialartsacademy.com', 'pwd_hash', 'Bad', 'Dob', '1980/05/15', 'student')
	`)
	if err == nil {
		t.Error("expected check constraint on date_of_birth format to fail, but it succeeded")
	}

	// Test foreign key constraint enforcement
	_, err = db.Exec(`
		INSERT INTO term_breaks (id, term_id, start_date, end_date, notes)
		VALUES ('tb-fake', 'nonexistent-term', '2026-06-01', '2026-06-07', 'Should fail')
	`)
	if err == nil {
		t.Error("expected foreign key constraint violation to fail, but it succeeded")
	}

	// Query values back to verify
	var email, firstName string
	err = db.QueryRow("SELECT email, first_name FROM users WHERE id = 'u-123'").Scan(&email, &firstName)
	if err != nil {
		t.Errorf("failed to query inserted user: %v", err)
	}
	if email != "test@martialartsacademy.com" || firstName != "Marcus" {
		t.Errorf("queried user values do not match: got %s, %s", email, firstName)
	}

	// Verify schema_migrations contains version 2
	var version2 int
	err = db.QueryRow("SELECT version FROM schema_migrations WHERE version = 2").Scan(&version2)
	if err != nil {
		t.Errorf("failed to find version 2 in migrations table: %v", err)
	}

	// Setup prerequisites for bookings test
	_, err = db.Exec(`
		INSERT INTO event_types (id, name, bg_color) VALUES (10, 'Test Event', '#000');
		INSERT INTO classes (id, term_id, event_type_id, name, day_of_week, start_time, end_time) 
		VALUES ('c-test', 't-summer-2026', 10, 'Test Class', 1, '10:00', '11:00');
		INSERT INTO class_occurrences (id, class_id, date, status) 
		VALUES ('o-test', 'c-test', '2026-06-22', 'scheduled');
	`)
	if err != nil {
		t.Fatalf("failed to insert booking test prerequisites: %v", err)
	}

	// Test valid booking with attendance_mode 'live_stream'
	_, err = db.Exec(`
		INSERT INTO bookings (id, user_id, occurrence_id, status, attendance_mode)
		VALUES ('b-valid-stream', 'u-123', 'o-test', 'confirmed', 'live_stream')
	`)
	if err != nil {
		t.Errorf("failed to insert valid booking: %v", err)
	}

	// Test invalid booking with invalid attendance_mode
	_, err = db.Exec(`
		INSERT INTO bookings (id, user_id, occurrence_id, status, attendance_mode)
		VALUES ('b-invalid-mode', 'u-123', 'o-test', 'confirmed', 'zoom_call')
	`)
	if err == nil {
		t.Error("expected invalid attendance_mode constraint to fail, but it succeeded")
	}

	// Verify schema_migrations contains version 3
	var version3 int
	err = db.QueryRow("SELECT version FROM schema_migrations WHERE version = 3").Scan(&version3)
	if err != nil {
		t.Errorf("failed to find version 3 in migrations table: %v", err)
	}

	// Test valid paid monthly subscription ($30/mo)
	_, err = db.Exec(`
		INSERT INTO user_subscriptions (id, user_id, term_id, start_date, end_date, status, payment_status, price_paid_cents)
		VALUES ('sub-paid-1', 'u-123', NULL, '2026-06-01', '2026-07-01', 'active', 'paid', 3000)
	`)
	if err != nil {
		t.Errorf("failed to insert valid subscription: %v", err)
	}

	// Test invalid subscription (negative price)
	_, err = db.Exec(`
		INSERT INTO user_subscriptions (id, user_id, term_id, start_date, end_date, status, payment_status, price_paid_cents)
		VALUES ('sub-bad-price', 'u-123', NULL, '2026-06-01', '2026-07-01', 'active', 'paid', -500)
	`)
	if err == nil {
		t.Error("expected check constraint on subscription price_paid_cents to fail, but it succeeded")
	}

	// Test invalid subscription (invalid payment status)
	_, err = db.Exec(`
		INSERT INTO user_subscriptions (id, user_id, term_id, start_date, end_date, status, payment_status, price_paid_cents)
		VALUES ('sub-bad-status', 'u-123', NULL, '2026-06-01', '2026-07-01', 'active', 'crypto_payment', 3000)
	`)
	if err == nil {
		t.Error("expected check constraint on subscription payment_status to fail, but it succeeded")
	}

	// Verify schema_migrations contains version 10
	var version10 int
	err = db.QueryRow("SELECT version FROM schema_migrations WHERE version = 10").Scan(&version10)
	if err != nil {
		t.Errorf("failed to find version 10 in migrations table: %v", err)
	}

	// Test locations inserting
	_, err = db.Exec(`
		INSERT INTO locations (id, name, address, is_active)
		VALUES ('loc-scarborough', 'Scarborough', '123 Lawrence Ave E, Toronto', 1)
	`)
	if err != nil {
		t.Errorf("failed to insert valid location: %v", err)
	}

	// Test invalid location active constraint
	_, err = db.Exec(`
		INSERT INTO locations (id, name, address, is_active)
		VALUES ('loc-bad', 'Bad Location', 'Address', 5)
	`)
	if err == nil {
		t.Error("expected check constraint on location is_active to fail, but it succeeded")
	}

	// Test token package pricing inserting
	_, err = db.Exec(`
		INSERT INTO token_packages (id, tokens_count, adult_price_cents, child_price_cents, senior_price_cents, is_active)
		VALUES ('pkg-custom', 20, 40000, 35000, 35000, 1)
	`)
	if err != nil {
		t.Errorf("failed to insert valid token package: %v", err)
	}

	// Test invalid token package (negative child price)
	_, err = db.Exec(`
		INSERT INTO token_packages (id, tokens_count, adult_price_cents, child_price_cents, senior_price_cents, is_active)
		VALUES ('pkg-bad-price', 20, 40000, -100, 35000, 1)
	`)
	if err == nil {
		t.Error("expected check constraint on token package price to fail, but it succeeded")
	}

	// Verify Registration FAQs count
	var faqCount int
	err = db.QueryRow("SELECT COUNT(*) FROM registration_faqs").Scan(&faqCount)
	if err != nil {
		t.Errorf("failed to query registration_faqs: %v", err)
	}
	if faqCount != 12 {
		t.Errorf("expected 12 registration FAQs seeded, got %d", faqCount)
	}

	// Test Family creation and mapping
	_, err = db.Exec(`
		INSERT INTO families (id, name) VALUES ('fam-lee', 'Lee Household')
	`)
	if err != nil {
		t.Errorf("failed to insert family: %v", err)
	}

	_, err = db.Exec(`
		INSERT INTO family_members (family_id, user_id) VALUES ('fam-lee', 'u-123')
	`)
	if err != nil {
		t.Errorf("failed to link user to family: %v", err)
	}

	// Test User Memberships constraints
	_, err = db.Exec(`
		INSERT INTO user_memberships (id, user_id, calendar_year, amount_cents, payment_date, receipt_issued)
		VALUES ('memb-2026', 'u-123', 2026, 10000, '2026-01-15', 1)
	`)
	if err != nil {
		t.Errorf("failed to insert valid user membership: %v", err)
	}

	// Verify unique constraint on user_id + calendar_year
	_, err = db.Exec(`
		INSERT INTO user_memberships (id, user_id, calendar_year, amount_cents, payment_date, receipt_issued)
		VALUES ('memb-2026-dup', 'u-123', 2026, 15000, '2026-02-20', 0)
	`)
	if err == nil {
		t.Error("expected unique constraint on (user_id, calendar_year) to fail, but it succeeded")
	}

	// Test Discount Rules constraint
	_, err = db.Exec(`
		INSERT INTO discount_rules (id, name, discount_type, trigger_value, discount_percentage)
		VALUES ('disc-bad-type', 'Bad Type', 'senior', 1, 10)
	`)
	if err == nil {
		t.Error("expected CHECK constraint on discount_type to fail, but it succeeded")
	}

	_, err = db.Exec(`
		INSERT INTO discount_rules (id, name, discount_type, trigger_value, discount_percentage)
		VALUES ('disc-bad-percentage', 'Bad Percentage', 'family', 2, 150)
	`)
	if err == nil {
		t.Error("expected CHECK constraint on discount_percentage to fail, but it succeeded")
	}

	// Test valid student registration with uniform option, ping pong, membership option, and fee breakdown
	_, err = db.Exec(`
		INSERT INTO student_term_registrations (
			id, user_id, term_id, token_package_id, membership_option_id,
			uniform_package_id, uniform_ordered, uniform_size, shoe_size, 
			ping_pong_package_id, ping_pong_club_joined,
			waiver_id, waiver_signer_name, waiver_signed_date, waiver_signed_ip,
			family_discount_applied_percent, returning_discount_applied_percent, 
			class_tokens_fee_cents, ping_pong_fee_cents, uniform_fee_cents, tax_cents, membership_fee_cents, total_fee_cents,
			payment_status
		) VALUES (
			'reg-1', 'u-123', 't-summer-2026', 'pkg-custom', 'memb-opt-standard',
			'uni-pkg-tshirt-high', 1, 'M', 42,
			'pp-pkg-1', 1,
			'waiver-class-registration', 'Justin Wood', '2026-06-20', '127.0.0.1',
			3, 5,
			32200, 2000, 5000, 5096, 2500, 47796,
			'paid'
		)
	`)
	if err != nil {
		t.Errorf("failed to insert valid student registration: %v", err)
	}

	// Test unique constraint (one registration per student per term)
	_, err = db.Exec(`
		INSERT INTO student_term_registrations (
			id, user_id, term_id, token_package_id, membership_option_id,
			waiver_id, waiver_signer_name, waiver_signed_date, waiver_signed_ip,
			uniform_ordered, class_tokens_fee_cents, tax_cents, total_fee_cents, payment_status
		) VALUES (
			'reg-dup', 'u-123', 't-summer-2026', 'pkg-custom', 'memb-opt-standard',
			'waiver-class-registration', 'Justin Wood', '2026-06-20', '127.0.0.1',
			0, 32200, 4186, 36386, 'pending'
		)
	`)
	if err == nil {
		t.Error("expected duplicate term registration to fail, but it succeeded")
	}

	// Seed t-fall-2026 so that constraint testing on other terms doesn't fail on foreign keys
	_, err = db.Exec(`
		INSERT INTO terms (id, name, start_date, end_date, is_active)
		VALUES ('t-fall-2026', 'Fall 2026', '2026-09-07', '2026-12-20', 0)
	`)
	if err != nil {
		t.Errorf("failed to insert second term: %v", err)
	}

	// Test invalid uniform size constraint failure
	_, err = db.Exec(`
		INSERT INTO student_term_registrations (
			id, user_id, term_id, token_package_id, membership_option_id,
			uniform_package_id, uniform_ordered, uniform_size, shoe_size, 
			waiver_id, waiver_signer_name, waiver_signed_date, waiver_signed_ip,
			class_tokens_fee_cents, tax_cents, total_fee_cents, payment_status
		) VALUES (
			'reg-bad-size', 'u-123', 't-fall-2026', 'pkg-custom', 'memb-opt-standard',
			'uni-pkg-tshirt-low', 1, 'XXS', 42, 
			'waiver-class-registration', 'Test Signer', '2026-06-20', '127.0.0.1',
			32200, 4186, 36386, 'pending'
		)
	`)
	if err == nil {
		t.Error("expected CHECK constraint on uniform_size to fail, but it succeeded")
	}

	// Test invalid shoe size constraint failure
	_, err = db.Exec(`
		INSERT INTO student_term_registrations (
			id, user_id, term_id, token_package_id, membership_option_id,
			uniform_package_id, uniform_ordered, uniform_size, shoe_size, 
			waiver_id, waiver_signer_name, waiver_signed_date, waiver_signed_ip,
			class_tokens_fee_cents, tax_cents, total_fee_cents, payment_status
		) VALUES (
			'reg-bad-shoes', 'u-123', 't-fall-2026', 'pkg-custom', 'memb-opt-standard',
			'uni-pkg-tshirt-low', 1, 'M', 55, 
			'waiver-class-registration', 'Test Signer', '2026-06-20', '127.0.0.1',
			32200, 4186, 36386, 'pending'
		)
	`)
	if err == nil {
		t.Error("expected CHECK constraint on shoe_size to fail, but it succeeded")
	}

	// Test Nullable token_package_id ("None" tokens option selection)
	_, err = db.Exec(`
		INSERT INTO student_term_registrations (
			id, user_id, term_id, token_package_id, membership_option_id,
			uniform_ordered, waiver_id, waiver_signer_name, waiver_signed_date, waiver_signed_ip,
			class_tokens_fee_cents, tax_cents, total_fee_cents, payment_status
		) VALUES (
			'reg-no-tokens', 'u-123', 't-fall-2026', NULL, 'memb-opt-odsp',
			0, 'waiver-class-registration', 'Test Signer', '2026-06-20', '127.0.0.1',
			0, 0, 0, 'pending'
		)
	`)
	if err != nil {
		t.Errorf("failed to insert student registration with NULL token package: %v", err)
	}

	// Test inserting token packages with size < 14 and terms_duration (should succeed)
	_, err = db.Exec(`
		INSERT INTO token_packages (id, tokens_count, terms_duration, adult_price_cents, child_price_cents, senior_price_cents, is_active)
		VALUES ('pkg-small-test', 2, 1, 5000, 4400, 4400, 1)
	`)
	if err != nil {
		t.Errorf("failed to insert token package with tokens_count < 14: %v", err)
	}

	// Test invalid token package (tokens_count = 0 should fail)
	_, err = db.Exec(`
		INSERT INTO token_packages (id, tokens_count, terms_duration, adult_price_cents, child_price_cents, senior_price_cents, is_active)
		VALUES ('pkg-invalid-zero', 0, 1, 1000, 1000, 1000, 1)
	`)
	if err == nil {
		t.Error("expected CHECK constraint on tokens_count = 0 to fail, but it succeeded")
	}

	// Test invalid token package (tokens_count = -5 should fail)
	_, err = db.Exec(`
		INSERT INTO token_packages (id, tokens_count, terms_duration, adult_price_cents, child_price_cents, senior_price_cents, is_active)
		VALUES ('pkg-invalid-neg', -5, 1, 1000, 1000, 1000, 1)
	`)
	if err == nil {
		t.Error("expected CHECK constraint on tokens_count = -5 to fail, but it succeeded")
	}

	// Test invalid token package (terms_duration = 0 should fail)
	_, err = db.Exec(`
		INSERT INTO token_packages (id, tokens_count, terms_duration, adult_price_cents, child_price_cents, senior_price_cents, is_active)
		VALUES ('pkg-invalid-duration', 10, 0, 1000, 1000, 1000, 1)
	`)
	if err == nil {
		t.Error("expected CHECK constraint on terms_duration = 0 to fail, but it succeeded")
	}

	// === Phase 2 Shopping Cart Support Test ===

	// 1. Verify promo codes table seeding and functionality
	var promoCount int
	err = db.QueryRow("SELECT COUNT(*) FROM promo_codes WHERE is_active = 1").Scan(&promoCount)
	if err != nil {
		t.Errorf("failed to query promo_codes: %v", err)
	}
	if promoCount < 3 {
		t.Errorf("expected at least 3 active promo codes seeded, got %d", promoCount)
	}

	// 2. Test invalid promo code constraint (invalid discount_type)
	_, err = db.Exec(`
		INSERT INTO promo_codes (id, discount_type, discount_value, min_order_value_cents, is_active)
		VALUES ('BADCODE', 'invalid_type', 10, 0, 1)
	`)
	if err == nil {
		t.Error("expected CHECK constraint on promo_codes discount_type to fail, but it succeeded")
	}

	// 3. Add items to shopping cart
	// Cart item representing a store product
	_, err = db.Exec(`
		INSERT INTO cart_items (id, user_id, product_id, registration_id, quantity)
		VALUES ('cart-item-prod', 'u-123', 'prod-hightop-shoes', NULL, 2)
	`)
	if err != nil {
		t.Errorf("failed to insert valid store product cart item: %v", err)
	}

	// Cart item representing a term class registration
	_, err = db.Exec(`
		INSERT INTO cart_items (id, user_id, product_id, registration_id, quantity)
		VALUES ('cart-item-reg', 'u-123', NULL, 'reg-no-tokens', 1)
	`)
	if err != nil {
		t.Errorf("failed to insert valid registration cart item: %v", err)
	}

	// Cart item check constraint violation: both product and registration specified
	_, err = db.Exec(`
		INSERT INTO cart_items (id, user_id, product_id, registration_id, quantity)
		VALUES ('cart-item-bad', 'u-123', 'prod-hightop-shoes', 'reg-no-tokens', 1)
	`)
	if err == nil {
		t.Error("expected CHECK constraint on cart_items (both product and registration) to fail, but it succeeded")
	}

	// Cart item check constraint violation: neither product nor registration specified
	_, err = db.Exec(`
		INSERT INTO cart_items (id, user_id, product_id, registration_id, quantity)
		VALUES ('cart-item-empty', 'u-123', NULL, NULL, 1)
	`)
	if err == nil {
		t.Error("expected CHECK constraint on cart_items (neither product nor registration) to fail, but it succeeded")
	}

	// Cart item unique constraint violation: duplicate product for user
	_, err = db.Exec(`
		INSERT INTO cart_items (id, user_id, product_id, registration_id, quantity)
		VALUES ('cart-item-dup', 'u-123', 'prod-hightop-shoes', NULL, 1)
	`)
	if err == nil {
		t.Error("expected UNIQUE constraint on (user_id, product_id) to fail, but it succeeded")
	}

	// 4. Test checkout process and calculations
	_, err = db.Exec(`
		INSERT INTO checkouts (
			id, user_id, registration_id, promo_code_id, donation_cents, payment_method,
			items_total_cents, discounts_total_cents, shipping_fee_cents, pre_tax_cents,
			tax_cents, credit_surcharge_cents, order_total_cents, payment_status
		) VALUES (
			'chk-test-1', 'u-123', 'reg-no-tokens', 'SUMMER2026', 1000, 'credit_card',
			7000, 700, 500, 7800,
			1014, 211, 9025, 'pending'
		)
	`)
	if err != nil {
		t.Errorf("failed to insert valid checkout transaction: %v", err)
	}

	// Checkout constraint validation: invalid payment method
	_, err = db.Exec(`
		INSERT INTO checkouts (
			id, user_id, registration_id, promo_code_id, donation_cents, payment_method,
			items_total_cents, discounts_total_cents, shipping_fee_cents, pre_tax_cents,
			tax_cents, credit_surcharge_cents, order_total_cents, payment_status
		) VALUES (
			'chk-bad-pay', 'u-123', NULL, NULL, 0, 'bitcoin',
			5000, 0, 0, 5000,
			650, 0, 5650, 'pending'
		)
	`)
	if err == nil {
		t.Error("expected CHECK constraint on payment_method to fail, but it succeeded")
	}

	// 5. Test store_orders mapping to checkout
	_, err = db.Exec(`
		INSERT INTO store_orders (id, user_id, pre_tax_cents, tax_cents, total_cents, payment_status, checkout_id)
		VALUES ('ord-test-checkout', 'u-123', 7000, 910, 7910, 'pending', 'chk-test-1')
	`)
	if err != nil {
		t.Errorf("failed to link store order to checkout: %v", err)
	}

	// === Phase 2 Order History Test ===

	// 1. Update checkout to set order_number and verify unique index constraint
	_, err = db.Exec(`
		UPDATE checkouts SET order_number = 10805 WHERE id = 'chk-test-1'
	`)
	if err != nil {
		t.Errorf("failed to set order_number on checkout: %v", err)
	}

	// Create another checkout with duplicate order_number (should fail unique index)
	_, err = db.Exec(`
		INSERT INTO checkouts (
			id, user_id, registration_id, promo_code_id, donation_cents, payment_method,
			items_total_cents, discounts_total_cents, shipping_fee_cents, pre_tax_cents,
			tax_cents, credit_surcharge_cents, order_total_cents, payment_status, order_number
		) VALUES (
			'chk-test-dup-number', 'u-123', NULL, NULL, 0, 'cash',
			1000, 0, 0, 1000,
			130, 0, 1130, 'pending', 10805
		)
	`)
	if err == nil {
		t.Error("expected duplicate order_number to fail unique index constraint, but it succeeded")
	}

	// 2. Link registration to checkout using checkout_id and verify it links correctly
	_, err = db.Exec(`
		UPDATE student_term_registrations SET checkout_id = 'chk-test-1' WHERE id = 'reg-no-tokens'
	`)
	if err != nil {
		t.Errorf("failed to set checkout_id on student_term_registrations: %v", err)
	}

	// Verify that we can query checkout and join both store orders and registrations
	var regCount, storeCount int
	err = db.QueryRow(`
		SELECT 
			(SELECT COUNT(*) FROM student_term_registrations WHERE checkout_id = 'chk-test-1') AS reg_count,
			(SELECT COUNT(*) FROM store_orders WHERE checkout_id = 'chk-test-1') AS store_count
	`).Scan(&regCount, &storeCount)
	if err != nil {
		t.Errorf("failed to query linked objects for checkout: %v", err)
	}
	if regCount != 1 || storeCount != 1 {
		t.Errorf("expected 1 registration and 1 store order linked to checkout 'chk-test-1', got reg: %d, store: %d", regCount, storeCount)
	}

	// 3. Link token transaction to checkout and verify it links correctly
	_, err = db.Exec(`
		INSERT INTO token_transactions (id, user_id, term_id, tokens_added, price_paid_cents, transaction_type, checkout_id)
		VALUES ('tx-test-cart', 'u-123', 't-summer-2026', 14, 28000, 'purchase', 'chk-test-1')
	`)
	if err != nil {
		t.Errorf("failed to insert token_transaction linked to checkout: %v", err)
	}

	var txCount int
	err = db.QueryRow("SELECT COUNT(*) FROM token_transactions WHERE checkout_id = 'chk-test-1'").Scan(&txCount)
	if err != nil {
		t.Errorf("failed to query linked token transactions: %v", err)
	}
	if txCount != 1 {
		t.Errorf("expected 1 token transaction linked to checkout 'chk-test-1', got %d", txCount)
	}

	// === Phase 2 Grading System Test ===

	// 1. Verify tracks are seeded
	var trackCount int
	err = db.QueryRow("SELECT COUNT(*) FROM grading_tracks").Scan(&trackCount)
	if err != nil {
		t.Errorf("failed to query grading_tracks: %v", err)
	}
	if trackCount != 3 {
		t.Errorf("expected 3 grading tracks, got %d", trackCount)
	}

	// 2. Test inserting invalid level (level_number = 6 should fail check constraint)
	_, err = db.Exec(`
		INSERT INTO grading_levels (id, track_id, level_number, name, min_training_months, mabu_level, mabu_duration_seconds)
		VALUES ('kf-level-6', 'kung_fu', 6, 'Kung Fu Level 6', 72, 4, 180)
	`)
	if err == nil {
		t.Error("expected CHECK constraint on level_number = 6 to fail, but it succeeded")
	}

	// 3. Test inserting a student grading exam (attempt)
	_, err = db.Exec(`
		INSERT INTO student_grading_exams (
			id, user_id, level_id, exam_date, examiner_id,
			mabu_duration_achieved_seconds, flexibility_percent_achieved,
			technical_score, effectiveness_score, knowledge_score, status, notes
		) VALUES (
			'exam-test-1', 'u-123', 'kf-level-1', '2026-06-20', 'u-123',
			130, 75, 80, 85, 90, 'passed', 'Excellent performance on Wububquan'
		)
	`)
	if err != nil {
		t.Errorf("failed to insert student grading exam: %v", err)
	}

	// 4. Test duplicate exam unique constraint failure (user, level, date)
	_, err = db.Exec(`
		INSERT INTO student_grading_exams (
			id, user_id, level_id, exam_date, examiner_id, status
		) VALUES (
			'exam-test-dup', 'u-123', 'kf-level-1', '2026-06-20', 'u-123', 'passed'
		)
	`)
	if err == nil {
		t.Error("expected unique constraint on (user_id, level_id, exam_date) to fail, but it succeeded")
	}

	// 5. Test invalid status constraint failure
	_, err = db.Exec(`
		INSERT INTO student_grading_exams (
			id, user_id, level_id, exam_date, status
		) VALUES (
			'exam-test-status', 'u-123', 'kf-level-1', '2026-06-21', 'gold_star'
		)
	`)
	if err == nil {
		t.Error("expected CHECK constraint on exam status to fail, but it succeeded")
	}

	// 6. Test caching passed grade in student_grades
	_, err = db.Exec(`
		INSERT INTO student_grades (user_id, level_id, passed_at, certificate_number)
		VALUES ('u-123', 'kf-level-1', '2026-06-20', 'CERT-KF1-0001')
	`)
	if err != nil {
		t.Errorf("failed to record passed grade: %v", err)
	}

	// === Phase 3 Student/Parent Guides Test ===

	// 1. Verify we can insert a valid guide
	_, err = db.Exec(`
		INSERT INTO student_parent_guides (id, title, summary, content_markdown, target_audience, category, display_order, last_updated_by)
		VALUES ('g-test-1', 'Test Guide', 'A guide summary', '# MD Content', 'student', 'classes_resources', 1, 'u-123')
	`)
	if err != nil {
		t.Errorf("failed to insert valid student_parent_guide: %v", err)
	}

	// 2. Verify target_audience CHECK constraint ('invalid_audience' should fail)
	_, err = db.Exec(`
		INSERT INTO student_parent_guides (id, title, summary, content_markdown, target_audience, category, display_order)
		VALUES ('g-test-2', 'Bad Audience Guide', 'summary', '# content', 'invalid_audience', 'classes_resources', 2)
	`)
	if err == nil {
		t.Error("expected CHECK constraint on target_audience to fail, but it succeeded")
	}

	// 3. Verify category CHECK constraint ('invalid_category' should fail)
	_, err = db.Exec(`
		INSERT INTO student_parent_guides (id, title, summary, content_markdown, target_audience, category, display_order)
		VALUES ('g-test-3', 'Bad Category Guide', 'summary', '# content', 'student', 'invalid_category', 3)
	`)
	if err == nil {
		t.Error("expected CHECK constraint on category to fail, but it succeeded")
	}
}

