<div align="center">

# 📋 On Duty (`on-duty`)
### Enterprise Duty-Leave & Attendance Workflow System for Amrita CSE Department
**Production-grade duty-leave management with compile-time type-safe SQL, versioned migrations, hierarchical role-based approvals, and automated attendance reconciliation.**

[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![Gin Framework](https://img.shields.io/badge/Framework-Gin%20Web%20Engine-008ECF?style=for-the-badge&logo=gin)](https://gin-gonic.com/)
[![PostgreSQL](https://img.shields.io/badge/Database-PostgreSQL%2016-336791?style=for-the-badge&logo=postgresql)](https://www.postgresql.org/)
[![sqlc](https://img.shields.io/badge/ORM-sqlc%20(Compile--Time%20Type--Safe)-4B5563?style=for-the-badge)](https://sqlc.dev/)
[![goose](https://img.shields.io/badge/Migrations-goose%20v3-F59E0B?style=for-the-badge)](https://github.com/pressly/goose)
[![React](https://img.shields.io/badge/Frontend-React%2018%20%2B%20Vite%20%2B%20Tailwind-61DAFB?style=for-the-badge&logo=react)](https://react.dev/)
[![License: Apache-2.0](https://img.shields.io/badge/License-Apache_2.0-blue.svg?style=for-the-badge)](https://opensource.org/licenses/Apache-2.0)

<br/>

```
   ____          ____        _         
  / __ \____    / __ \__  __/ /___  __ 
 / / / / __ \  / / / / / / / __/ / / / 
/ /_/ / / / / / /_/ / /_/ / /_/ /_/ /  
\____/_/ /_/ /_____/\__,_/\__/\__, /   
                             /____/    
 Amrita CSE On-Duty & Attendance Workflow Engine
```

<p align="center">
  <i>Automating academic leaves, conference representations, hackathon delegations, and Anokha techfest participation with zero manual paper-chasing.</i>
</p>

</div>

---

## 🏛️ System Overview

**On Duty** is an enterprise-grade academic duty-leave and attendance workflow system engineered specifically for the Department of Computer Science & Engineering at **Amrita Vishwa Vidyapeetham**.

In large engineering departments, hundreds of students participate in national hackathons (Smart India Hackathon), technical conferences, coding summits, sports meets, and national techfests like **Anokha**. Traditional paper OD forms frequently lead to lost applications, delayed faculty approvals, attendance discrepancies, and administrative overhead.

**On Duty** solves this with:
- 🛡️ **Role-Based Access Control (RBAC)** across 5 tiers: `student`, `faculty_advisor`, `club_coordinator`, `hod`, and `dept_admin`.
- ⚡ **Strict Hierarchical State Machine**:
  `Student Submission` ➔ `Faculty Advisor Review` ➔ `Event / Club Coordinator Validation` ➔ `HOD Final Sign-off` ➔ `Automated Attendance Adjustment`.
- 🔒 **Type-Safe SQL with `sqlc`**: Zero dynamic SQL queries. 100% of queries are verified at compile-time with type-safe Go structs and pgx/v5 drivers.
- 📜 **Versioned Schema Migrations with `goose`**: Automated version-controlled database migrations with complete rollback (`+goose Up` / `+goose Down`) support.
- 📊 **Department Analytics**: Real-time tracking of semester quota caps (maximum 20 OD days per semester), category utilization, and attendance impact.

---

## 📐 Approval Workflow State Machine

```
   [ Student ]
        |
        v  (Submit Application: Category, Dates, Reason, Proof)
 [ PENDING_ADVISOR ]  <------------------------+
        |                                      |
        +-- (Reject) ------------> [ REJECTED ]|
        |                                      |
        v  (Approve)                           |
   Requires Coord?                             |
      /        \                               |
    (Yes)      (No)                            |
    /            \                             |
   v              \                            |
[ PENDING_COORD ]  \                           |
   |                \                          |
   +-- (Reject) -----\-------------------------+
   |                  \                        
   v  (Approve)        v                       
       [ PENDING_HOD ]                         
              |                                
              +-- (Reject) --------------------+
              |                                
              v  (Approve)                     
        [ APPROVED ]                           
              |                                
              v                                
   (Auto Attendance Credit)                    
```

---

## 🛠️ Tech Stack & Architecture

| Layer | Technologies | Key Architectural Patterns |
| :--- | :--- | :--- |
| **Backend API** | Go 1.22+, Gin Web Framework | Clean Architecture, Middleware pipelines, Graceful shutdown |
| **Database** | PostgreSQL 16, pgx/v5 | Foreign key cascades, check constraints, transaction isolation |
| **SQL Generation** | `sqlc` v2 | Compile-time query verification, zero runtime reflection |
| **Migrations** | `goose` v3 | Versioned SQL migration files with idempotent rollbacks |
| **Security & Auth**| JWT (HS256), bcrypt (cost 10) | Stateless tokens, Argon2/bcrypt password hashing, RBAC claims |
| **Frontend** | React 18, Vite, Tailwind CSS | Persona switcher, responsive timelines, status cards |
| **DevOps** | Docker, Docker Compose, GitHub Actions | Multi-stage scratch build, CI test pipelines |

---

## 🚀 Quickstart Guide

### 1. Run with Docker Compose (Recommended)
```bash
git clone https://github.com/miriyaladhanwinn/on-duty.git
cd on-duty

# Start PostgreSQL 16, Redis, and Gin REST API
docker-compose up --build
```

The API will be live at `http://localhost:8080`.

### 2. Run Go Backend Locally
```bash
# Verify unit & integration test suite
go test -v ./...

# Build and start server
go run cmd/api/main.go
```

Output:
```text
===================================================
   ____          ____        _         
  / __ \____    / __ \__  __/ /___  __ 
 / / / / __ \  / / / / / / / __/ / / / 
/ /_/ / / / / / /_/ / /_/ / /_/ /_/ /  
\____/_/ /_/ /_____/\__,_/\__/\__, /   
                             /____/    
 Amrita CSE On-Duty & Attendance Workflow Engine
===================================================
[+] Amrita CSE On-Duty REST API server listening on :8080
[+] Swagger/API Docs: http://localhost:8080/api/v1/health
```

### 3. Run Frontend UI
```bash
cd frontend
npm install
npm run dev
```

---

## 📡 REST API Reference

### Authentication & Users
| Method | Endpoint | Access | Description |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/auth/login` | Public | Authenticate user, returns 24h JWT token |
| `POST` | `/api/v1/auth/register`| Public | Register new student or faculty account |
| `GET`  | `/api/v1/auth/me` | Authenticated | Fetch current user claims and active role |

### On-Duty Applications & Workflow
| Method | Endpoint | Access | Description |
| :--- | :--- | :--- | :--- |
| `GET`  | `/api/v1/categories` | Public | List OD categories and semester day quotas |
| `POST` | `/api/v1/od/requests` | `student` | Submit new duty-leave request |
| `GET`  | `/api/v1/od/requests` | Authenticated | List applications (own for students, queue for faculty) |
| `GET`  | `/api/v1/od/requests/:id` | Authenticated | Retrieve request details with full approval timeline |
| `POST` | `/api/v1/od/requests/:id/action` | `faculty`, `coord`, `hod` | Submit `APPROVE` or `REJECT` decision with remarks |
| `GET`  | `/api/v1/od/stats` | Authenticated | View department summary statistics |

---

## 🧪 Comprehensive Automated Tests

```bash
go test -v ./...
```

```text
=== RUN   TestHealthEndpoint
--- PASS: TestHealthEndpoint (0.06s)
=== RUN   TestLoginAndGetCategories
--- PASS: TestLoginAndGetCategories (0.11s)
PASS
ok      github.com/miriyaladhanwinn/on-duty/internal/api        0.933s
=== RUN   TestTokenGenerationAndValidation
--- PASS: TestTokenGenerationAndValidation (0.00s)
=== RUN   TestPasswordHashingAndVerification
--- PASS: TestPasswordHashingAndVerification (0.19s)
PASS
ok      github.com/miriyaladhanwinn/on-duty/internal/auth       0.681s
=== RUN   TestCompleteApprovalWorkflow
--- PASS: TestCompleteApprovalWorkflow (0.06s)
=== RUN   TestCategoryDaysQuotaEnforcement
--- PASS: TestCategoryDaysQuotaEnforcement (0.06s)
PASS
ok      github.com/miriyaladhanwinn/on-duty/internal/service    0.580s
```

---

## 📄 License

This project is licensed under the [Apache-2.0 License](LICENSE).
Copyright © 2026 [MIRIYALA DHANWINN](https://github.com/miriyaladhanwinn).
Amrita Vishwa Vidyapeetham, Department of Computer Science & Engineering.
