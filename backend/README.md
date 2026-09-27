# Shaolin Academy Portal - Go Backend Engine

Welcome to the backend engine for the **Shaolin Academy** school management system. This Go-based backend provides a highly optimized, fully tested, and secure API supporting a React/TypeScript/TailwindCSS frontend. It manages martial arts curriculum schedules, class bookings, token-based transactions, dynamic online checkouts, student grading tracks, lay disciple certifications, parent guides, and student forums.

---

## 🛠️ Technology Stack

1. **Language**: Go (1.22+) using the standard library `net/http` multiplexer with modern wildcard path parameters.
2. **Database**: SQLite (relational schema with strict foreign key constraints and `CHECK` validation limits).
3. **Authentication**: JSON Web Tokens (JWT) signed using a secure key, combined with `golang.org/x/crypto/bcrypt` for secure hashing.
4. **Telemetry & Dashboard**: Ephemeral in-memory analytics engine with query sparks, health logs, and transaction rollback recovery gauges.

---

## 📂 Codebase Directory Layout

```
backend/
├── cmd/
│   ├── api/             # Production API Server Entrypoint (listening on Port 8081)
│   │   ├── main.go      # Wire repositories, handlers, and middlewares
│   │   └── *_test.go    # Endpoint-level API integration test suites
│   ├── dbmigrate/       # Standalone database migration runner CLI
│   └── explorer/        # Developer Visual Console Entrypoint (listening on Port 8080)
├── internal/
│   ├── api/             # Shared HTTP Handlers & Middlewares (Auth, Roles, Logs, CORS)
│   ├── domain/          # Core definitions and abstract repository interfaces (Clean Arch)
│   └── repository/
│       └── sqlite/      # SQLite implementation of domain repository interfaces
└── migrations/          # SQL database migration scripts (0001 to 0021)
```

---

## ⚡ Key Features

*   **Charitable Household Links**: Allows student accounts to group under a shared family household to manage charitable memberships, waivers, and discounts.
*   **Token Transaction Ledger**: Track purchased token packages, deduct tokens upon class bookings, and credit them back on eligible cancellations.
*   **Waitlisting & Attendance Modes**: Blocks physical class bookings once capacity is reached (waitlisting the student) while letting online live-stream bookings pass capacity limits.
*   **Dynamic Checkout Calculations**: Validates shopping carts, processes promo codes (percentage or flat discounts), enforces shipping thresholds, supports non-profit donations, and handles credit card surcharges (2.4%).
*   **Belt Exams & Certification Gating**: Stores student stance times (Mabu) and grades, preventing access to restricted kung fu digital media files until the student passes prerequisite levels.
*   **Role-Based Thread Security**: Protects forum message board topics so students can read/reply on open boards but are prevented from posting on admin announcement boards.

---

## 🚀 Getting Started

### Prerequisites

Ensure you have [Go](https://go.dev/doc/install) installed.

### Setup and Running

1.  **Clone the Repository** and navigate to the backend directory:
    ```bash
    cd backend
    ```

2.  **Run All Unit Tests**:
    Verify database migrations, repositories, constraints, and REST controllers:
    ```bash
    go test ./...
    ```

3.  **Start the Developer Explorer & Dashboard**:
    This starts the developer visual console which seeds the database automatically:
    ```bash
    go run cmd/explorer/main.go
    ```
    *   **Dashboard URL**: `http://localhost:8080`
    *   *Features*: Visual database telemetry, query sparkline logs, interactive SQL terminal, backoffice CRUD records editor, and transaction scenarios lab.

4.  **Start the Production API Server**:
    This initializes the production server with auto-migrations:
    ```bash
    go run cmd/api/main.go
    ```
    *   **API URL**: `http://localhost:8081`

---

## 📋 REST API Reference

### 🔐 Auth & Student Profile
*   `POST /api/auth/register` - Create a student profile
*   `POST /api/auth/login` - Authenticate credentials and receive a Bearer JWT
*   `GET /api/auth/me` - Fetch profile metadata (`AuthRequired`)

### 👥 Household & Leaderboard
*   `GET /api/users/family` - Retrieve linked family members (`AuthRequired`)
*   `POST /api/users/family` - Create a new family household (`AuthRequired`)
*   `POST /api/users/family/members` - Link a user by email to the family (`AuthRequired`)
*   `GET /api/users/memberships` - Query charitable memberships (`AuthRequired`)
*   `POST /api/users/memberships` - Issue charitable membership receipts (`AdminOnly`)
*   `GET /api/users/leaderboard/weekly` - Fetch weekly point ranking spectrum
*   `GET /api/users/leaderboard/overall` - Fetch cumulative student points leaderboard

### 📅 Terms & Schedule Configuration
*   `GET /api/terms` - List school curriculum terms
*   `POST /api/terms` - Create a new term (`AdminOnly`)
*   `PUT /api/terms/{id}` - Modify term details (`AdminOnly`)
*   `DELETE /api/terms/{id}` - Terminate a term (`AdminOnly`)
*   `GET /api/terms/{term_id}/breaks` - Retrieve breaks during a term
*   `GET /api/classes` - Search and list classes catalog
*   `GET /api/occurrences` - Query class occurrences range

### 🎟️ Bookings & Token Transactions
*   `GET /api/tokens/balance?term_id={id}` - Check remaining tokens for a term (`AuthRequired`)
*   `GET /api/tokens/packages` - Fetch available token pricing packages
*   `GET /api/tokens/transactions` - View audit transactions history ledger (`AuthRequired`)
*   `POST /api/tokens/purchase` - Buy packages and update balance (`AuthRequired`)
*   `POST /api/bookings` - Confirm a class booking (checks capacity & token limits) (`AuthRequired`)
*   `DELETE /api/bookings/{id}` - Cancel a booking and process eligible token refunds (`AuthRequired`)

### 🛒 E-Commerce Online Store
*   `GET /api/store/products` - Catalog of products grouped by categories
*   `GET /api/cart` - Retrieve shopping cart items (`AuthRequired`)
*   `POST /api/cart` - Add products or registrations to cart (`AuthRequired`)
*   `DELETE /api/cart/{id}` - Delete item from cart (`AuthRequired`)
*   `POST /api/store/checkout` - Unified order checkout (calculates tax, discounts, surcharges) (`AuthRequired`)
*   `GET /api/store/orders/history` - User billing order history (`AuthRequired`)

### 🏆 Grading Tracks & Exams
*   `GET /api/grading/tracks` - Fetch complete grading syllabus track structure
*   `GET /api/grading/exams/my-history` - Student's exam attempts and grade status (`AuthRequired`)
*   `POST /api/grading/exams` - Evaluate student marks and issue certificates (`AdminOnly`)

### 💬 Forums, Media & Guides
*   `GET /api/forums/boards` - List message boards
*   `GET /api/forums/topics/{id}` - List threads under board
*   `GET /api/forums/topics/{id}/posts` - Read replies under topic thread
*   `POST /api/forums/topics` - Create new topic thread (validates roles) (`AuthRequired`)
*   `POST /api/forums/posts` - Add a reply post (`AuthRequired`)
*   `GET /api/media` - Fetch qualifications-gated video portfolio (`AuthRequired`)
*   `GET /api/guides` - Retrieve markdown onboarding manuals (Parent/Student)
