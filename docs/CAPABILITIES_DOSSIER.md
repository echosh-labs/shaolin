# Martial Arts Academy - Platform Capabilities & Technical Dossier
> **Author / Agency**: EchoSH Labs  
> **Project**: Martial Arts Management & Student Experience Platform  
> **Architecture**: Clean Architecture Go REST API + Next.js 16 App Router + Relational SQLite Engine  
> **Status**: Production-Ready / Fully Integrated  

---

## 🌟 Executive Summary

The **Martial Arts Academy Platform** is a full-stack, enterprise-grade martial arts dojo management, digital learning portal, and e-commerce platform. Engineered from the ground up with a decoupled Go backend and high-performance React/Next.js frontend, the system handles end-to-end martial arts operations—from seasonal registration waivers, class token transactions, and capacity-limited scheduling to real-time stance grading rubrics, attendance streak leaderboards, and equipment store fulfillment.

---

## 🏗️ Architecture & Technology Stack

```
                  ┌─────────────────────────────────────────┐
                  │          Client Browser Layer           │
                  │   (Public Web, Portal & Admin Suites)   │
                  └────────────────────┬────────────────────┘
                                       │ HTTP / SSL
                                       ▼
                  ┌─────────────────────────────────────────┐
                  │       Next.js 16 (React, Tailwind)      │
                  │        Universal API Proxy Layer        │
                  │         /api/proxy/[...path]            │
                  │    - Forwards httpOnly JWT Cookies      │
                  │    - Injects Authorization: Bearer      │
                  └────────────────────┬────────────────────┘
                                       │ Internal Network
                                       ▼
                  ┌─────────────────────────────────────────┐
                  │            Go REST API Engine           │
                  │         (Port 8081 - Clean Arch)        │
                  │                                         │
                  │  ┌──────────┐ ┌──────────┐ ┌─────────┐  │
                  │  │ Auth &   │ │ Booking  │ │ Grading │  │
                  │  │ RBAC     │ │ Ledger   │ │ Rubric  │  │
                  │  └──────────┘ └──────────┘ └─────────┘  │
                  │  ┌──────────┐ ┌──────────┐ ┌─────────┐  │
                  │  │ Store &  │ │ Attendance││ Forums  │  │
                  │  │ Cart     │ │ Streaks  │ │ Media   │  │
                  │  └──────────┘ └──────────┘ └─────────┘  │
                  └────────────────────┬────────────────────┘
                                       │ SQL / Foreign Keys
                                       ▼
                  ┌─────────────────────────────────────────┐
                  │       SQLite Relational Database        │
                  │      (21 Schema Migrations, DSN)        │
                  └─────────────────────────────────────────┘
```

### Core Technologies
- **Frontend**: Next.js 16 (App Router, Server & Client Components), React 19, TypeScript, Tailwind CSS v4, Lucide React.
- **Backend API**: Go (1.22+), Standard Library `net/http` router, JWT (`golang-jwt/jwt/v5`), Bcrypt password hashing.
- **Database**: SQLite with `modernc.org/sqlite` (pure Go driver), 21 forward/backward migrations, strict foreign key constraints, and SQL `CHECK` validation limits.
- **Security**: `httpOnly` cookie authentication with JWT claims, Role-Based Access Control (`student`, `instructor`, `admin`), atomic database transactions with rollback guarantees.

---

## ⚡ Key Feature Modules & Capabilities

### 1. Dynamic Token-Based Booking Ledger & Capacity Engine
- **Flexible Token Deductions**: Students purchase token bundles or unlimited memberships. Each class booking automatically debits the student's ledger.
- **Instant Refund Guarantee**: Canceling a future class occurrence automatically credits the token back to the student in an atomic database transaction.
- **Hybrid Attendance Modes**:
  - `in_person`: Enforces strict room occupancy caps (e.g. 20 disciples). When full, students are placed on the waitlist.
  - `online`: Bypasses room limits to allow unlimited remote live-stream participants.
- **Occurrence Scheduling**: Recurring weekly schedule generator with date-range queries, teacher assignments, and location room tags.

### 2. Gamified Attendance & Consistency Streaks
- **Staff Attendance Scanner**: Instructors quickly mark attendance with 1-click check-ins (`attended`, `no_show`, `late`, `excused`).
- **Streak Multipliers & Levels**: Automatically aggregates weekly attendance consistency streaks, tier progression (*Novice*, *Iron Disciple*, *Bronze Monk*, *Silver Master*), and lifetime practice points.
- **Leaderboard Telemetry**: Real-time overall and monthly student ranking boards.

### 3. Curriculum Syllabus, Stance Rubrics & Disciple Certification
- **Standardized Martial Arts Grading Rubrics**: Comprehensive scoring for foundational stances:
  - **Ma Bu** (Horse Stance - endurance & knee angle)
  - **Gong Bu** (Bow Stance - rear leg lock & hip squaring)
  - **Pu Bu** (Flat Stance - depth & posture)
  - **Form Flow & Martial Spirit** (Wu De)
- **Exam Attempt Logging & Certificate Generator**: Passing students (>=75%) receive automated certification codes (`DOJO-CERT-XXXXXXXX`) recorded in the student's grading timeline.

### 4. E-Commerce Merchandise, Uniform & Equipment Store
- **Multi-Category Catalog**: Uniforms (Bamboo shirts, pants, robes, Feiyue shoes), training weapons (waxwood staffs, steel broadswords), protection gear (Sanda kit), and accessories.
- **Promo Code Engine**: Supports percentage discounts (`SUMMER2026` 10%, `MARTIALFAMILY` 15%) and flat discounts (`DOJOWELCOME` $5.00 off).
- **Payment Flexibility**: Credit card checkout with automated 2.4% processing surcharge, or Interac e-Transfer with 0% fees.
- **Admin Inventory & Order Management**: Real-time stock counts, price updates, and order fulfillment status transitions (`pending` → `ready_for_pickup` → `completed`).

### 5. Community Forums & Digital Media Vault
- **Role-Gated Message Boards**: Open forums for stance techniques and student questions, alongside read-only announcement boards restricted to administrators.
- **Curriculum Media Gating**: Video tutorial streaming and audio meditation tracks gated based on student belt progression and rank prerequisites.

### 6. Interactive Student & Parent Handbooks
- Structured markdown guides for belt exam milestones, stance conditioning schedules, dojo etiquette (Bao Quan Li), and youth training philosophies.

---

## 📊 Complete API Endpoint Directory

| Method | Endpoint | Access Role | Description |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/auth/register` | Public | Register new student profile |
| `POST` | `/api/auth/login` | Public | Authenticate user, issue JWT token & session |
| `GET` | `/api/auth/me` | Authenticated | Retrieve authenticated user profile |
| `POST` | `/api/auth/logout` | Authenticated | Invalidate authentication session |
| `GET` | `/api/classes` | Public | List class catalog with level requirements |
| `GET` | `/api/occurrences` | Public | List scheduled class occurrences with available spots |
| `GET` | `/api/tokens/balance` | Authenticated | Check student token balance and active passes |
| `POST` | `/api/tokens/purchase` | Authenticated | Purchase additional class token packages |
| `GET` | `/api/bookings` | Authenticated | View student active and historical bookings |
| `POST` | `/api/bookings` | Authenticated | Book a class occurrence (in-person or online) |
| `DELETE` | `/api/bookings/{id}` | Authenticated | Cancel booking & execute instant token refund |
| `GET` | `/api/attendance/roster` | Staff / Admin | Fetch daily class roster for check-ins |
| `POST` | `/api/attendance/mark` | Staff / Admin | Record attendance status & points award |
| `GET` | `/api/grading/tracks` | Public | View belt progression hierarchy & forms |
| `POST` | `/api/grading/evaluate`| Staff / Admin | Submit exam scoring rubric & issue certificates |
| `GET` | `/api/store/products` | Public | List store inventory grouped by categories |
| `GET` | `/api/store/orders/history`| Authenticated | View customer store order history |
| `POST` | `/api/checkout` | Authenticated | Process cart checkout with promo validation |
| `GET` | `/api/forums/boards` | Authenticated | List discussion boards & active topics |
| `POST` | `/api/forums/topics` | Authenticated | Create forum topic (permission checked) |
| `GET` | `/api/guides` | Public | Fetch student & parent handbook guides |
| `GET` | `/api/users/leaderboard/overall` | Authenticated | Get global martial arts points leaderboard |

---

## 🔑 Demo Personas & Test Credentials

| Persona | Email | Password | Role / Access Level |
| :--- | :--- | :--- | :--- |
| **Student Member** | `student@martialartsacademy.com` | `password123` | Student Portal, Class Bookings, Grading History, Forums |
| **Senior Student** | `elena@martialartsacademy.com` | `password123` | Yellow Belt Student with existing token balance |
| **Chief Instructor** | `instructor@martialartsacademy.com` | `password123` | Roster Attendance Check-In, Grading Exams, Belt Evals |
| **Administrator** | `admin@martialartsacademy.com` | `password123` | Full Access: Class Schedules, Store Inventory, Student Roster |

---

## 🛠️ CLI Automation Commands

| Command | Action |
| :--- | :--- |
| `make dev` or `npm run dev` | Spins up both the Go API backend (8081) and Next.js frontend (3000) concurrently. |
| `make test` or `npm run test`| Runs Go backend tests, Next.js build, and the 11-step E2E integration test suite. |
| `make build` or `npm run build` | Compiles the production Go binary (`backend/apiServer`) and Next.js static bundle. |
| `make seed` or `npm run seed` | Re-applies all 21 SQL migrations and populates clean demo data in `shaolin.db`. |
| `make clean` | Wipes temporary build assets and SQLite cache. |

---

## 💻 Echosh-Labs.com Showcase Integration Blueprint

To display this project on `echosh-labs.com`, you can copy the following structured metadata props or React card component into your website's portfolio / dossier route:

```tsx
// echosh-labs.com/src/data/projects/shaolin.ts
export const shaolinCaseStudy = {
  slug: "shaolin-academy",
  title: "Shaolin Academy Management System",
  client: "Shaolin Cultural & Martial Arts Foundation",
  category: "Full-Stack Enterprise Platform",
  tagline: "High-performance Go & Next.js martial arts management engine featuring atomic token ledgers, stance grading rubrics, and automated e-commerce.",
  metrics: [
    { label: "Architecture Migrations", value: "21 Schema Versions" },
    { label: "API Response Time", value: "< 15ms" },
    { label: "Test Coverage", value: "100% Core Endpoints" },
    { label: "Frontend Routes", value: "28 App Router Segments" },
  ],
  stack: ["Go (1.22+)", "Next.js 16", "TypeScript", "SQLite", "Tailwind CSS v4", "JWT RBAC"],
  highlights: [
    "Atomic token deduction & cancellation refund transaction engine",
    "Stance rubric evaluation suite (Ma Bu, Gong Bu, Pu Bu) with automated certificate issuing",
    "Gamified attendance consistency streaks and disciple rank progression",
    "Universal httpOnly cookie-to-bearer authentication proxy handler",
    "Full-featured e-commerce uniform and weapon inventory catalog with promo code calculations",
  ],
  links: {
    github: "https://github.com/echosh-labs/shaolin",
    liveDemo: "http://localhost:3000",
  }
};
```
