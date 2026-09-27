# 🥋 Shaolin Platform — Martial Arts Dojo Management & Portal System
> A full-stack school management system, belt grading syllabus, tokenized class booking engine, and community hub for traditional martial arts academies.

---

## 🌟 The Vibe
Shaolin Platform bridges ancient traditional martial arts lineage with modern full-stack web architecture. It provides an all-in-one portal where students book classes using attendance tokens, track belt syllabus criteria and lay disciple requirements, purchase uniforms, and discuss techniques in the forum—while dojo administrators manage rosters, attendance roll-calls, store inventory, and grading exams.

---

## 🚀 60-Second Quickstart

```bash
# 1. Clone the repository
git clone https://github.com/echosh-labs/shaolin.git
cd shaolin

# 2. Seed the SQLite database with pristine demo records
make seed

# 3. Boot both the Go API Server and Next.js Frontend
make dev
```

Once running:
- **Student & Public Portal**: Visit [http://localhost:3000](http://localhost:3000)
- **Go REST API & Health**: Visit [http://localhost:8081/health](http://localhost:8081/health)

---

## 🔑 Demo Login Accounts

The database comes pre-seeded with three tiers of authenticated users (all passwords are `password123`):

| Role | Email Address | Password | Capabilities |
| :--- | :--- | :--- | :--- |
| **Head Master / Admin** | `admin@martialartsacademy.com` | `password123` | Full administrative control: student grading, class schedules, store inventory, and attendance records. |
| **Chief Instructor** | `instructor@martialartsacademy.com` | `password123` | Mark attendance, evaluate belt milestones, and review student bookings. |
| **Enrolled Student** | `student@martialartsacademy.com` | `password123` | Book training sessions, view belt grading syllabus, purchase uniforms, and participate in community forums. |

---

## 🗺️ Interactive Tour & Web Routes

### Public Pages
- [`/`](http://localhost:3000/): Dojo landing page with lineage overview, programs, and schedule previews.
- [`/classes`](http://localhost:3000/classes): Class descriptions (Kung Fu, Qigong, Sanshou sparring, weapons).
- [`/schedule`](http://localhost:3000/schedule): Live weekly timetable with filterable time slots.
- [`/tuition`](http://localhost:3000/tuition): Membership tiers, drop-in token packages, and unlimited terms.
- [`/store`](http://localhost:3000/store): Online dojo gear store (sashes, uniforms, training shoes).
- [`/guides`](http://localhost:3000/guides): Martial etiquette guide, uniform care, and parent handbook.

### Student Portal (`/portal/*`)
- [`/portal`](http://localhost:3000/portal): Student dashboard showing current rank, attendance streak, and next grading date.
- [`/portal/bookings`](http://localhost:3000/portal/bookings): Interactive calendar to reserve spots in upcoming training sessions using tokens.
- [`/portal/grading`](http://localhost:3000/portal/grading): Detailed curriculum syllabus checklist required for the next belt rank.
- [`/portal/tokens`](http://localhost:3000/portal/tokens): Token balance manager and top-up checkout.
- [`/portal/community`](http://localhost:3000/portal/community): Academy message board and video technique discussions.

### Admin Console (`/admin/*`)
- [`/admin`](http://localhost:3000/admin): Executive dashboard with revenue sparks, active headcount, and daily check-ins.
- [`/admin/attendance`](http://localhost:3000/admin/attendance): Real-time tablet-friendly dojo check-in and attendance taker.
- [`/admin/grading`](http://localhost:3000/admin/grading): Candidate evaluation portal for awarding belt advancements.
- [`/admin/students`](http://localhost:3000/admin/students): Student roster manager and waiver compliance status.
- [`/admin/schedules`](http://localhost:3000/admin/schedules): Class calendar and instructor assignment editor.
- [`/admin/store`](http://localhost:3000/admin/store): Inventory manager, restocking alerts, and order fulfillment.

---

## 🏗️ Technical Architecture

```
shaolin/
├── backend/                  # High-performance Go Backend (Port 8081)
│   ├── cmd/
│   │   ├── api/              # Production REST API entrypoint & test suites
│   │   ├── seed/             # Database migration and sample fixture seeder
│   │   └── dbmigrate/        # Standalone schema migration tool
│   ├── internal/
│   │   ├── api/              # HTTP route handlers (auth, bookings, store, grading)
│   │   ├── domain/           # Core domain entity interfaces
│   │   └── repository/       # SQLite persistence engines with strict foreign keys
│   └── migrations/           # 21 atomic up/down SQL schema migrations
├── frontend/                 # Next.js 15 (App Router + TailwindCSS) (Port 3000)
│   ├── src/app/(admin)/      # Administrative management console
│   ├── src/app/(portal)/     # Authenticated student portal
│   ├── src/app/(public)/     # Public marketing site and online checkout
│   └── src/context/          # Auth and shopping cart state providers
├── scripts/                  # Lifecycle scripts (dev.sh, seed.sh, build.sh, test.sh)
└── Makefile                  # Unified automation interface
```

### Backend Engine Highlights
- **Standard Library Routing**: Uses modern Go standard library `net/http` path routing with wildcard matching (`GET /api/v1/classes/{id}`).
- **SQLite Engine**: Powered by `modernc.org/sqlite` (pure Go, zero CGO required) with WAL journal mode and foreign key constraints enabled.
- **Authentication**: Stateless HMAC-SHA256 JWT tokens with bcrypt password hashing (`cost 10`).

---

## ⚙️ Available Makefile Commands

| Command | Action |
| :--- | :--- |
| `make dev` | Starts Go backend (Port 8081) and Next.js frontend (Port 3000) concurrently. |
| `make seed` | Re-runs database migrations and re-seeds clean demo fixtures into `shaolin.db`. |
| `make test` | Executes full Go backend test suite, frontend linting, and typechecks. |
| `make build` | Compiles Go binary (`backend/bin/apiServer`) and builds Next.js production bundle. |
| `make clean` | Removes compiled binaries, build artifacts, and test database files. |

---

## 📄 License

This software is dual-licensed:
- **Open Source Edition**: Governed by the [GNU Affero General Public License v3.0 (AGPL-3.0)](LICENSE.md) for individual, educational, and open-source usage.
- **Commercial & Enterprise Edition**: Requires a commercial license from echoSH labs for proprietary integration, corporate deployment, or advanced enterprise features. See [COMMERCIAL.md](COMMERCIAL.md) or visit [echosh-labs.com](https://echosh-labs.com).

For commercial licensing inquiries, contact [justin@echosh-labs.com](mailto:justin@echosh-labs.com).
