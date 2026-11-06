# BETA Project: Architecture, Ideology & Comprehensive Repository Status

> **Context Document for AI Models & Developers**
> This document provides a complete architectural overview, coding standards, language & library stacks, in-depth documentation of the active **Photography CRM** application, and an exhaustive file-by-file catalog of every single source file in the `BETA` codebase.

---

## 1. Overall Project Status & Ideology

### 1.1 Architecture & Multi-Module Workspace
The `BETA` project is organized as a high-performance Go multi-module workspace (`go.work`) adopting **Clean Architecture** and modern **Hypermedia-Driven UI** principles. Rather than maintaining heavy Single Page Application (SPA) frontend frameworks that require complex build chains, `BETA` leverages **Server-Side Rendered (SSR) Go templates with A-H/Templ**, **HTMX** for selective DOM swaps, **Vanilla JavaScript** for targeted micro-interactions, and pure CSS for sleek obsidian-dark/enterprise-slate themes.

The workspace unifies multiple distinct applications and shared packages:
- **`cmd/runner`**: Multi-service runner and supervisor that can boot all micro-services concurrently.
- **`cmd/seed`**: CLI database seeding utility for populating roles and initial users into any target SQLite database.
- **`shared`**: Cross-cutting cryptographic security, session tokens, and password hashing.
- **`projects/core`**: Domain data layer, SQLite connection pooling (`DBPool`), Sqlc-generated typesafe queries, repository patterns, and authentication middleware.
- **`components`**: Reusable Templ UI component library (buttons, cards, badges, dialogs, charts, spotlight search).
- **`projects/photography_crm`**: Flagship production SaaS CRM tailored for professional photographers and studios.
- **`projects/cms`**: Content management service sharing the unified core authentication and schema.
- **`projects/superadmin`**: Multi-tenant administrative portal for managing users, organizations, and global settings.
- **`projects/restaurant_pos`**: Placeholder for upcoming point-of-sale module.

### 1.2 Technology Stack & Modules Used
| Layer | Technologies & Libraries |
|---|---|
| **Core Language** | Go 1.26+ (`go.work` multi-module structure) |
| **HTTP Routing** | `github.com/go-chi/chi/v5` (modular sub-routing, clean middleware chaining) |
| **Templating Engine** | `github.com/a-h/templ` (type-safe, pre-compiled Go HTML templates) |
| **Database & Storage** | SQLite 3 via `modernc.org/sqlite` (pure Go, CGO-free, WAL mode enabled) |
| **Type-Safe SQL** | `sqlc` (compiles raw SQL into idiomatic Go structs and queries) |
| **Security & Crypto** | Argon2id & `golang.org/x/crypto/bcrypt`, HMAC-SHA256 session cookie encryption |
| **Frontend Hypermedia** | HTMX (`htmx.min.js`), SortableJS (`sortable.min.js`), FullCalendar (`fullcalendar.min.js`) |
| **Design & Icons** | Custom Obsidian Dark & Enterprise Slate CSS design tokens, Lucide SVG icons |
| **Configuration** | `github.com/joho/godotenv` with fallback environment defaults |
| **API Testing** | Bruno collection runner (`.yml` requests) and standard Go `testing` package |

### 1.3 Coding Style & Engineering Principles
1. **Zero-CGO Pure Go Deployment**: Database drivers use `modernc.org/sqlite` to ensure single-binary portable builds without external C compiler dependencies.
2. **Separation of Concerns**: Each service cleanly delineates:
   - `routes/router.go`: Middlewares and top-level routing assembly.
   - `routes/frontend.go`: Page endpoints and HTML template rendering.
   - `routes/backend.go`: JSON REST APIs and mutation endpoints.
   - `handlers/`: Request parsing, service calls, and view model preparation.
   - `ui/pages/` and `ui/layouts/`: Pure `.templ` presentation code.
3. **Compiled Static Assets**: All CSS, JavaScript, and SVG assets are bundled into Go binaries using `//go:embed` via `embed.go`, eliminating missing file issues in production.
4. **Safe Concurrency & WAL Mode**: SQLite connection pools (`DBPool`) enforce Write-Ahead Logging (`PRAGMA journal_mode = WAL;`) and synchronous normal mode for high-throughput reads and concurrent writes.
5. **Server & Port Safety Rules**: Workspace execution rules forbid running unmonitored background daemons; ports remain clean when tests or builds complete.


### 1.4 Architecture Evolution & Middleware Specification
The project evolved from an initial single-service HttpRouter/SQLx setup to a modular, multi-service Go workspace (`go.work`):
- **Routing**: Migrated from `HttpRouter` to `go-chi/chi/v5` for standard `http.Handler` interoperability, idiomatic route grouping (`r.Group`), and clean middleware stacking.
- **Database Access**: Replaced raw SQLx with `modernc.org/sqlite` (pure Go, CGO-free, WAL mode) and `sqlc` for compile-time verified database models and queries.
- **Shared Architecture**: Extracted cross-cutting concerns into `projects/core` (auth service, repository pattern, DBPool) and `shared/security` (Argon2id hashing, HMAC sessions).

#### Middleware Convention
All HTTP middlewares adhere strictly to standard Go `http.Handler` wrapping:
```go
package routes

import "net/http"

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
```
Available in `projects/photography_crm/routes/middleware.go` and `projects/core/middleware/middleware.go`, this pattern allows standard middleware functions (`DatabaseRoutingMiddleware`, `RequireAuth`, `RequireAuthPage`, `RequireRole`) to be used both with `Chain(handler, m1, m2)` and directly in Chi routers via `r.Use(m)`.


---

## 2. In-Depth Focus: Photography CRM (`projects/photography_crm`)

### 2.1 Mission & Domain Purpose
The **Photography CRM** is the centerpiece application of the `BETA` platform. It provides end-to-end studio management for commercial, portrait, and wedding photographers. It unifies the entire client lifecycle—from initial lead capture and follow-ups, to contract signing, photoshoot scheduling, deliverable tracking, invoice generation, and financial reporting.

### 2.2 Functional Modules & Feature Breakdown
- **Leads & Pipeline Management (`/leads`, `/leads/details`)**:
  - Kanban pipeline visualization and table views for client inquiries.
  - Drag-and-drop lead stages powered by SortableJS and HTMX.
  - Lead score tracking, source attribution, and conversion metrics.
- **Inbox & Communication (`/inbox`)**:
  - Unified messaging interface for client email communications, notes, and direct threads.
- **Follow-Ups Engine (`/followups`)**:
  - Automated scheduling and reminder queues to ensure prompt prospect replies.
- **Client Directory (`/clients`, `/clients/details`)**:
  - Comprehensive client relationship records, booking history, family/event notes, and contact details.
- **Projects & Shoots (`/projects`, `/projects/details`)**:
  - Milestone-based shoot project tracking, crew assignments, gear checklists, and location details.
- **Sessions & Bookings (`/sessions`)**:
  - Individual photoshoot session tracking, slot booking, and time logs.
- **Interactive Studio Calendar (`/calendar`)**:
  - Integrated FullCalendar.js interface displaying booked sessions, due deliverables, and team schedules.
- **Contracts & E-Signatures (`/contracts`)**:
  - Digital photography contracts, release forms, and legal template management.
- **Deliverables & Proofing (`/deliverables`)**:
  - Digital gallery deliverable pipelines, proofing review statuses, retouching queues, and download limits.
- **Billing & Invoicing (`/billing`, `/orders`, `/orders/details`)**:
  - Order calculation, split deposit invoicing, tax calculation, and payment status tracking.
- **Studio Analytics & Reports (`/reports`)**:
  - Revenue trajectory, conversion funnel charts, average shoot value, and seasonal performance KPIs.
- **Comprehensive Studio Settings (`/settings/*`)**:
  - Deep customization including Business Settings, Crew & Assistant management, Deliverable Presets, Forms, Integrations, Invoice Templates, Packages, Quotation Templates, Sequence ID auto-numbering, Services, and Multi-stage Workflows.

### 2.3 UX/UI Design & Interaction Architecture
- **Obsidian Dark & Clean Slate Themes**: Implements high-contrast dark tones (`#08080a`, `#131318`) alongside enterprise light mode tokens.
- **Spotlight Search (`<cmd-k>`)**: Universal quick search component dynamically filtering leads, projects, clients, and settings.
- **Kanban & Drag-and-Drop**: Interactive board layouts (`ui/layouts/kanban_layout.templ`) with optimistic UI updates.
- **Modular CSS Isolation**: Each feature page maintains its own stylesheet in `ui/static/css/pages/` preventing unintended CSS cascades.

### 2.4 Photography CRM Endpoints & Routing Architecture
The application boots via `projects/photography_crm/main.go` on port `5500` (or `8082` via the runner):
```go
// projects/photography_crm/routes/router.go
func SetupRouter(pool *utils.DBPool, cfg utils.Config) *chi.Mux {
    r := chi.NewRouter()
    RegisterFrontendRoutes(r, cfg)
    RegisterBackendRoutes(r, pool, cfg)
    return r
}
```

---

## 3. Exhaustive File-by-File Repository Catalog

Below is the complete, exhaustive catalog of every file across the `BETA` project. Each entry contains the **absolute file path** and a concise **3 to 5 line description** detailing its role and implementation within the codebase.

### 1. Root Workspace & Environment Configuration (6 files)

#### `/home/desktop/especa/beta/.env`
- **Relative Path**: `.env`
- **Purpose**: Root environment configuration declaring global ports, database paths, and session keys.

#### `/home/desktop/especa/beta/.gitignore`
- **Relative Path**: `.gitignore`
- **Purpose**: Specifies files and patterns ignored by Git, including compiled binaries, database journals, and temporary build outputs.

#### `/home/desktop/especa/beta/_temp`
- **Relative Path**: `_temp`
- **Purpose**: Scratchpad and temporary developer notes detailing database schemas, task lists, and command snippets.

#### `/home/desktop/especa/beta/generate_status.py`
- **Relative Path**: `generate_status.py`
- **Purpose**: Project file generate_status.py delivering essential application configuration or source logic.

#### `/home/desktop/especa/beta/go.work`
- **Relative Path**: `go.work`
- **Purpose**: Defines the Go 1.26 multi-module workspace integrating runner, components, core, shared, and applications into an unified build structure.

#### `/home/desktop/especa/beta/go.work.sum`
- **Relative Path**: `go.work.sum`
- **Purpose**: Cryptographic checksum verification lockfile for module dependencies across the entire Go workspace.

---

### 2. Multi-Service Runner (`cmd/runner`) (3 files)

#### `/home/desktop/especa/beta/cmd/runner/go.mod`
- **Relative Path**: `cmd/runner/go.mod`
- **Purpose**: Go module manifest declaring dependencies for the cmd/runner multi-service supervisor.

#### `/home/desktop/especa/beta/cmd/runner/go.sum`
- **Relative Path**: `cmd/runner/go.sum`
- **Purpose**: Checksum lockfile ensuring reproducible builds for the cmd/runner supervisor service.

#### `/home/desktop/especa/beta/cmd/runner/main.go`
- **Relative Path**: `cmd/runner/main.go`
- **Purpose**: Unified supervisor application that concurrently boots CMS (:8081), Photography CRM (:8082), and Superadmin (:8083) with graceful shutdown.

---

### 3. Database Seed CLI (`cmd/seed`) (2 files)

#### `/home/desktop/especa/beta/cmd/seed/go.mod`
- **Relative Path**: `cmd/seed/go.mod`
- **Purpose**: Go module manifest declaring dependencies for the standalone CLI database seeding tool.

#### `/home/desktop/especa/beta/cmd/seed/main.go`
- **Relative Path**: `cmd/seed/main.go`
- **Purpose**: CLI utility tool to seed database roles and initial users into SQLite databases via flags (`-db`, `-role`, `-first`, `-last`, `-email`, `-password`).

---

### 4. Shared Reusable UI Component Library (`components`) (34 files)

#### `/home/desktop/especa/beta/components/avatar.templ`
- **Relative Path**: `components/avatar.templ`
- **Purpose**: Templ markup definition providing design tokens and slots for the reusable avatar component.

#### `/home/desktop/especa/beta/components/avatar_templ.go`
- **Relative Path**: `components/avatar_templ.go`
- **Purpose**: Compiled Go template code generated by templ for the reusable avatar UI component.

#### `/home/desktop/especa/beta/components/badge.templ`
- **Relative Path**: `components/badge.templ`
- **Purpose**: Templ markup definition providing design tokens and slots for the reusable badge component.

#### `/home/desktop/especa/beta/components/badge_templ.go`
- **Relative Path**: `components/badge_templ.go`
- **Purpose**: Compiled Go template code generated by templ for the reusable badge UI component.

#### `/home/desktop/especa/beta/components/baseui/lifecycle.js`
- **Relative Path**: `components/baseui/lifecycle.js`
- **Purpose**: Client-side component lifecycle orchestrator managing dynamic DOM mounts and state bindings.

#### `/home/desktop/especa/beta/components/bubble.templ`
- **Relative Path**: `components/bubble.templ`
- **Purpose**: Templ markup definition providing design tokens and slots for the reusable bubble component.

#### `/home/desktop/especa/beta/components/bubble_templ.go`
- **Relative Path**: `components/bubble_templ.go`
- **Purpose**: Compiled Go template code generated by templ for the reusable bubble UI component.

#### `/home/desktop/especa/beta/components/button.templ`
- **Relative Path**: `components/button.templ`
- **Purpose**: Templ markup definition providing design tokens and slots for the reusable button component.

#### `/home/desktop/especa/beta/components/button_templ.go`
- **Relative Path**: `components/button_templ.go`
- **Purpose**: Compiled Go template code generated by templ for the reusable button UI component.

#### `/home/desktop/especa/beta/components/card.templ`
- **Relative Path**: `components/card.templ`
- **Purpose**: Templ markup definition providing design tokens and slots for the reusable card component.

#### `/home/desktop/especa/beta/components/card_templ.go`
- **Relative Path**: `components/card_templ.go`
- **Purpose**: Compiled Go template code generated by templ for the reusable card UI component.

#### `/home/desktop/especa/beta/components/chart/chart.js`
- **Relative Path**: `components/chart/chart.js`
- **Purpose**: Lightweight chart rendering engine providing interactive SVG chart data visualization.

#### `/home/desktop/especa/beta/components/chart/chart.templ`
- **Relative Path**: `components/chart/chart.templ`
- **Purpose**: Templ markup definition providing design tokens and slots for the reusable chart component.

#### `/home/desktop/especa/beta/components/chart/chart_templ.go`
- **Relative Path**: `components/chart/chart_templ.go`
- **Purpose**: Compiled Go template code generated by templ for the reusable chart UI component.

#### `/home/desktop/especa/beta/components/chart/scale.go`
- **Relative Path**: `components/chart/scale.go`
- **Purpose**: Mathematical scale calculation functions converting data ranges into SVG coordinate points.

#### `/home/desktop/especa/beta/components/dialog.templ`
- **Relative Path**: `components/dialog.templ`
- **Purpose**: Templ markup definition providing design tokens and slots for the reusable dialog component.

#### `/home/desktop/especa/beta/components/dialog_templ.go`
- **Relative Path**: `components/dialog_templ.go`
- **Purpose**: Compiled Go template code generated by templ for the reusable dialog UI component.

#### `/home/desktop/especa/beta/components/go.mod`
- **Relative Path**: `components/go.mod`
- **Purpose**: Go module definition for the shared, reusable Templ UI component library.

#### `/home/desktop/especa/beta/components/go.sum`
- **Relative Path**: `components/go.sum`
- **Purpose**: Checksum lockfile verifying dependencies for the components module.

#### `/home/desktop/especa/beta/components/input.templ`
- **Relative Path**: `components/input.templ`
- **Purpose**: Templ markup definition providing design tokens and slots for the reusable input component.

#### `/home/desktop/especa/beta/components/input_templ.go`
- **Relative Path**: `components/input_templ.go`
- **Purpose**: Compiled Go template code generated by templ for the reusable input UI component.

#### `/home/desktop/especa/beta/components/lucide_icons.templ`
- **Relative Path**: `components/lucide_icons.templ`
- **Purpose**: Templ markup definition providing design tokens and slots for the reusable lucide_icons component.

#### `/home/desktop/especa/beta/components/lucide_icons_templ.go`
- **Relative Path**: `components/lucide_icons_templ.go`
- **Purpose**: Compiled Go template code generated by templ for the reusable lucide_icons UI component.

#### `/home/desktop/especa/beta/components/marker.templ`
- **Relative Path**: `components/marker.templ`
- **Purpose**: Templ markup definition providing design tokens and slots for the reusable marker component.

#### `/home/desktop/especa/beta/components/marker_templ.go`
- **Relative Path**: `components/marker_templ.go`
- **Purpose**: Compiled Go template code generated by templ for the reusable marker UI component.

#### `/home/desktop/especa/beta/components/message.templ`
- **Relative Path**: `components/message.templ`
- **Purpose**: Templ markup definition providing design tokens and slots for the reusable message component.

#### `/home/desktop/especa/beta/components/message_templ.go`
- **Relative Path**: `components/message_templ.go`
- **Purpose**: Compiled Go template code generated by templ for the reusable message UI component.

#### `/home/desktop/especa/beta/components/spotlight_search.templ`
- **Relative Path**: `components/spotlight_search.templ`
- **Purpose**: Templ markup definition providing design tokens and slots for the reusable spotlight_search component.

#### `/home/desktop/especa/beta/components/spotlight_search_templ.go`
- **Relative Path**: `components/spotlight_search_templ.go`
- **Purpose**: Compiled Go template code generated by templ for the reusable spotlight_search UI component.

#### `/home/desktop/especa/beta/components/table.templ`
- **Relative Path**: `components/table.templ`
- **Purpose**: Templ markup definition providing design tokens and slots for the reusable table component.

#### `/home/desktop/especa/beta/components/table_templ.go`
- **Relative Path**: `components/table_templ.go`
- **Purpose**: Compiled Go template code generated by templ for the reusable table UI component.

#### `/home/desktop/especa/beta/components/tabs.templ`
- **Relative Path**: `components/tabs.templ`
- **Purpose**: Templ markup definition providing design tokens and slots for the reusable tabs component.

#### `/home/desktop/especa/beta/components/tabs_templ.go`
- **Relative Path**: `components/tabs_templ.go`
- **Purpose**: Compiled Go template code generated by templ for the reusable tabs UI component.

#### `/home/desktop/especa/beta/components/utils/utils.go`
- **Relative Path**: `components/utils/utils.go`
- **Purpose**: Utility helpers for HTML class concatenation, conditional styles, and string transformations.

---

### 5. Shared Cross-Cutting Security Module (`shared`) (5 files)

#### `/home/desktop/especa/beta/shared/go.mod`
- **Relative Path**: `shared/go.mod`
- **Purpose**: Go module definition for the shared cross-cutting security and cryptography utilities.

#### `/home/desktop/especa/beta/shared/go.sum`
- **Relative Path**: `shared/go.sum`
- **Purpose**: Checksum lockfile for the shared utilities module.

#### `/home/desktop/especa/beta/shared/security/password.go`
- **Relative Path**: `shared/security/password.go`
- **Purpose**: Provides Argon2id and bcrypt cryptographic password hashing, salting, and constant-time verification.

#### `/home/desktop/especa/beta/shared/security/security_test.go`
- **Relative Path**: `shared/security/security_test.go`
- **Purpose**: Unit tests verifying cryptographic integrity of password hashing and session signature validation.

#### `/home/desktop/especa/beta/shared/security/session.go`
- **Relative Path**: `shared/security/session.go`
- **Purpose**: Provides HMAC-SHA256 authenticated session cookie generation, encryption, and validation.

---

### 6. Core Business, Database & Auth Layer (`projects/core`) (23 files)

#### `/home/desktop/especa/beta/projects/core/db/queries/queries.sql`
- **Relative Path**: `projects/core/db/queries/queries.sql`
- **Purpose**: SQL query definitions consumed by sqlc to generate type-safe Go database functions.

#### `/home/desktop/especa/beta/projects/core/db/schema/schema.go`
- **Relative Path**: `projects/core/db/schema/schema.go`
- **Purpose**: Initializes SQLite schema, performs automatic table migrations, and validates database connectivity.

#### `/home/desktop/especa/beta/projects/core/db/schema/schema.sql`
- **Relative Path**: `projects/core/db/schema/schema.sql`
- **Purpose**: Master SQLite DDL schema declaring tables, constraints, foreign keys, and indexes for users and roles.

#### `/home/desktop/especa/beta/projects/core/db/sqlc/db.go`
- **Relative Path**: `projects/core/db/sqlc/db.go`
- **Purpose**: Generated DBTX interface and database connection handler created by sqlc.

#### `/home/desktop/especa/beta/projects/core/db/sqlc/models.go`
- **Relative Path**: `projects/core/db/sqlc/models.go`
- **Purpose**: Generated Go structs representing table rows (User, Role) with strict typing.

#### `/home/desktop/especa/beta/projects/core/db/sqlc/queries.sql.go`
- **Relative Path**: `projects/core/db/sqlc/queries.sql.go`
- **Purpose**: SQL query definitions consumed by sqlc to generate type-safe Go database functions.

#### `/home/desktop/especa/beta/projects/core/go.mod`
- **Relative Path**: `projects/core/go.mod`
- **Purpose**: Go module definition for core shared backend logic, database abstraction, and auth services.

#### `/home/desktop/especa/beta/projects/core/go.sum`
- **Relative Path**: `projects/core/go.sum`
- **Purpose**: Checksum lockfile for dependencies used in the core backend package.

#### `/home/desktop/especa/beta/projects/core/handlers/handlers.go`
- **Relative Path**: `projects/core/handlers/handlers.go`
- **Purpose**: Shared HTTP controller handlers implementing core REST API endpoints and error handling.

#### `/home/desktop/especa/beta/projects/core/middleware/middleware.go`
- **Relative Path**: `projects/core/middleware/middleware.go`
- **Purpose**: HTTP middleware providing session verification, cookie extraction, and route guard authentication.

#### `/home/desktop/especa/beta/projects/core/models/models.go`
- **Relative Path**: `projects/core/models/models.go`
- **Purpose**: Data transfer objects (DTOs) and request/response models shared across microservices.

#### `/home/desktop/especa/beta/projects/core/repositories/repositories.go`
- **Relative Path**: `projects/core/repositories/repositories.go`
- **Purpose**: Repository layer encapsulating database access logic behind clean interface abstractions.

#### `/home/desktop/especa/beta/projects/core/services/auth_service.go`
- **Relative Path**: `projects/core/services/auth_service.go`
- **Purpose**: Authentication service containing business logic for logins, signups, credential checks, and role enforcement.

#### `/home/desktop/especa/beta/projects/core/sqlc.yaml`
- **Relative Path**: `projects/core/sqlc.yaml`
- **Purpose**: Sqlc configuration file configuring SQLite code generation for type-safe database queries.

#### `/home/desktop/especa/beta/projects/core/tests/auth_service_test.go`
- **Relative Path**: `projects/core/tests/auth_service_test.go`
- **Purpose**: Automated unit and integration test suite asserting correctness of core auth_service.

#### `/home/desktop/especa/beta/projects/core/tests/password_test.go`
- **Relative Path**: `projects/core/tests/password_test.go`
- **Purpose**: Automated unit and integration test suite asserting correctness of core password.

#### `/home/desktop/especa/beta/projects/core/tests/repository_test.go`
- **Relative Path**: `projects/core/tests/repository_test.go`
- **Purpose**: Automated unit and integration test suite asserting correctness of core repository.

#### `/home/desktop/especa/beta/projects/core/tests/session_test.go`
- **Relative Path**: `projects/core/tests/session_test.go`
- **Purpose**: Automated unit and integration test suite asserting correctness of core session.

#### `/home/desktop/especa/beta/projects/core/utils/config.go`
- **Relative Path**: `projects/core/utils/config.go`
- **Purpose**: Configuration parser loading environment variables into structured Config structs with defaults.

#### `/home/desktop/especa/beta/projects/core/utils/db.go`
- **Relative Path**: `projects/core/utils/db.go`
- **Purpose**: Database connection pool manager setting up SQLite in WAL mode with connection limits and test hooks.

#### `/home/desktop/especa/beta/projects/core/utils/response.go`
- **Relative Path**: `projects/core/utils/response.go`
- **Purpose**: Standardized HTTP JSON response utilities providing uniform success and error payloads.

#### `/home/desktop/especa/beta/projects/core/utils/user_context.go`
- **Relative Path**: `projects/core/utils/user_context.go`
- **Purpose**: Context utilities for injecting and retrieving authenticated user claims from HTTP request context.

#### `/home/desktop/especa/beta/projects/core/validation/validation.go`
- **Relative Path**: `projects/core/validation/validation.go`
- **Purpose**: Input validation logic ensuring email syntax, password complexity, and required field safety.

---

### 7. Photography CRM (`projects/photography_crm`) (141 files)

#### `/home/desktop/especa/beta/projects/photography_crm/.env`
- **Relative Path**: `projects/photography_crm/.env`
- **Purpose**: Active local configuration specifying SQLite DB path, port 5500, session secrets, and test database flags.

#### `/home/desktop/especa/beta/projects/photography_crm/.env.example`
- **Relative Path**: `projects/photography_crm/.env.example`
- **Purpose**: Template environment file outlining configuration parameters needed to run Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/data/photography_crm.db`
- **Relative Path**: `projects/photography_crm/data/photography_crm.db`
- **Purpose**: Local SQLite database / WAL journal storage file preserving photography CRM data (photography_crm.db).

#### `/home/desktop/especa/beta/projects/photography_crm/data/photography_crm.db-shm`
- **Relative Path**: `projects/photography_crm/data/photography_crm.db-shm`
- **Purpose**: Local SQLite database / WAL journal storage file preserving photography CRM data (photography_crm.db-shm).

#### `/home/desktop/especa/beta/projects/photography_crm/data/photography_crm.db-wal`
- **Relative Path**: `projects/photography_crm/data/photography_crm.db-wal`
- **Purpose**: Local SQLite database / WAL journal storage file preserving photography CRM data (photography_crm.db-wal).

#### `/home/desktop/especa/beta/projects/photography_crm/data/photography_crm_test.db`
- **Relative Path**: `projects/photography_crm/data/photography_crm_test.db`
- **Purpose**: Local SQLite database / WAL journal storage file preserving photography CRM data (photography_crm_test.db).

#### `/home/desktop/especa/beta/projects/photography_crm/data/photography_crm_test.db-shm`
- **Relative Path**: `projects/photography_crm/data/photography_crm_test.db-shm`
- **Purpose**: Local SQLite database / WAL journal storage file preserving photography CRM data (photography_crm_test.db-shm).

#### `/home/desktop/especa/beta/projects/photography_crm/data/photography_crm_test.db-wal`
- **Relative Path**: `projects/photography_crm/data/photography_crm_test.db-wal`
- **Purpose**: Local SQLite database / WAL journal storage file preserving photography CRM data (photography_crm_test.db-wal).

#### `/home/desktop/especa/beta/projects/photography_crm/go.mod`
- **Relative Path**: `projects/photography_crm/go.mod`
- **Purpose**: Go module file defining Photography CRM dependencies, including Chi v5, templ, and local core modules.

#### `/home/desktop/especa/beta/projects/photography_crm/go.sum`
- **Relative Path**: `projects/photography_crm/go.sum`
- **Purpose**: Cryptographic dependency checksum lockfile ensuring reproducible builds for Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/handlers/auth_handler.go`
- **Relative Path**: `projects/photography_crm/handlers/auth_handler.go`
- **Purpose**: HTTP controllers managing user authentication, login redirects, session creation, and sign-out in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/handlers/crm_handlers.go`
- **Relative Path**: `projects/photography_crm/handlers/crm_handlers.go`
- **Purpose**: HTTP handler suite rendering CRM views: leads, clients, orders, projects, sessions, calendar, and settings.

#### `/home/desktop/especa/beta/projects/photography_crm/handlers/home_handler.go`
- **Relative Path**: `projects/photography_crm/handlers/home_handler.go`
- **Purpose**: Root dashboard handler loading pipeline summaries, active session metrics, and recent inquiries.

#### `/home/desktop/especa/beta/projects/photography_crm/main.go`
- **Relative Path**: `projects/photography_crm/main.go`
- **Purpose**: Standalone HTTP entrypoint loading configs, establishing SQLite DBPool, setting up Chi router, and serving Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/routes/backend.go`
- **Relative Path**: `projects/photography_crm/routes/backend.go`
- **Purpose**: API router declaring JSON REST endpoints for roles, authentication, and background CRM data actions.

#### `/home/desktop/especa/beta/projects/photography_crm/routes/frontend.go`
- **Relative Path**: `projects/photography_crm/routes/frontend.go`
- **Purpose**: Frontend router binding URLs (/leads, /projects, /calendar, /settings, etc.) to templated view handlers.

#### `/home/desktop/especa/beta/projects/photography_crm/routes/router.go`
- **Relative Path**: `projects/photography_crm/routes/router.go`
- **Purpose**: Main Chi router constructor attaching logging, session recovery, static file serving, and subrouters.

#### `/home/desktop/especa/beta/projects/photography_crm/tests/auth_api_test.go`
- **Relative Path**: `projects/photography_crm/tests/auth_api_test.go`
- **Purpose**: Automated test suite asserting correct behavior of Photography CRM auth_api.

#### `/home/desktop/especa/beta/projects/photography_crm/tests/auth_middleware_test.go`
- **Relative Path**: `projects/photography_crm/tests/auth_middleware_test.go`
- **Purpose**: Automated test suite asserting correct behavior of Photography CRM auth_middleware.

#### `/home/desktop/especa/beta/projects/photography_crm/tests/bruno/Photography_crm/opencollection.yml`
- **Relative Path**: `projects/photography_crm/tests/bruno/Photography_crm/opencollection.yml`
- **Purpose**: Bruno API collection test file verifying CRM endpoint contracts (opencollection.yml).

#### `/home/desktop/especa/beta/projects/photography_crm/tests/frontend_auth_test.go`
- **Relative Path**: `projects/photography_crm/tests/frontend_auth_test.go`
- **Purpose**: Automated test suite asserting correct behavior of Photography CRM frontend_auth.

#### `/home/desktop/especa/beta/projects/photography_crm/tests/middleware_test.go`
- **Relative Path**: `projects/photography_crm/tests/middleware_test.go`
- **Purpose**: Automated test suite asserting correct behavior of Photography CRM middleware.

#### `/home/desktop/especa/beta/projects/photography_crm/tests/signup_api_test.go`
- **Relative Path**: `projects/photography_crm/tests/signup_api_test.go`
- **Purpose**: Automated test suite asserting correct behavior of Photography CRM signup_api.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/components/common.templ`
- **Relative Path**: `projects/photography_crm/ui/components/common.templ`
- **Purpose**: Templ component definitions for reusable UI fragments, search modals, and status badges in CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/components/common_templ.go`
- **Relative Path**: `projects/photography_crm/ui/components/common_templ.go`
- **Purpose**: Compiled Go template code for reusable CRM subcomponents like search bars and status pills.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/layouts/auth_layout.templ`
- **Relative Path**: `projects/photography_crm/ui/layouts/auth_layout.templ`
- **Purpose**: Templ markup defining the auth_layout layout structure, navigation sidebar, topbar, and slots.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/layouts/auth_layout_templ.go`
- **Relative Path**: `projects/photography_crm/ui/layouts/auth_layout_templ.go`
- **Purpose**: Compiled Go template code for the auth_layout UI layout shell in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/layouts/dashboard_layout.templ`
- **Relative Path**: `projects/photography_crm/ui/layouts/dashboard_layout.templ`
- **Purpose**: Templ markup defining the dashboard_layout layout structure, navigation sidebar, topbar, and slots.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/layouts/dashboard_layout_templ.go`
- **Relative Path**: `projects/photography_crm/ui/layouts/dashboard_layout_templ.go`
- **Purpose**: Compiled Go template code for the dashboard_layout UI layout shell in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/layouts/kanban_layout.templ`
- **Relative Path**: `projects/photography_crm/ui/layouts/kanban_layout.templ`
- **Purpose**: Templ markup defining the kanban_layout layout structure, navigation sidebar, topbar, and slots.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/layouts/kanban_layout_templ.go`
- **Relative Path**: `projects/photography_crm/ui/layouts/kanban_layout_templ.go`
- **Purpose**: Compiled Go template code for the kanban_layout UI layout shell in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/layouts/settings_layout.templ`
- **Relative Path**: `projects/photography_crm/ui/layouts/settings_layout.templ`
- **Purpose**: Templ markup defining the settings_layout layout structure, navigation sidebar, topbar, and slots.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/layouts/settings_layout_templ.go`
- **Relative Path**: `projects/photography_crm/ui/layouts/settings_layout_templ.go`
- **Purpose**: Compiled Go template code for the settings_layout UI layout shell in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/billing/billing_overview.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/billing/billing_overview.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for billing/billing_overview.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/billing/billing_overview_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/billing/billing_overview_templ.go`
- **Purpose**: Compiled Go template implementation for the billing/billing_overview view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/billing/cancelled.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/billing/cancelled.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for billing/cancelled.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/billing/cancelled_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/billing/cancelled_templ.go`
- **Purpose**: Compiled Go template implementation for the billing/cancelled view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/billing/columns_dropdown.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/billing/columns_dropdown.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for billing/columns_dropdown.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/billing/columns_dropdown_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/billing/columns_dropdown_templ.go`
- **Purpose**: Compiled Go template implementation for the billing/columns_dropdown view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/billing/edit_invoice_modal.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/billing/edit_invoice_modal.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for billing/edit_invoice_modal.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/billing/edit_invoice_modal_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/billing/edit_invoice_modal_templ.go`
- **Purpose**: Compiled Go template implementation for the billing/edit_invoice_modal view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/billing/expenses.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/billing/expenses.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for billing/expenses.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/billing/expenses_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/billing/expenses_templ.go`
- **Purpose**: Compiled Go template implementation for the billing/expenses view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/billing/invoices.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/billing/invoices.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for billing/invoices.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/billing/invoices_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/billing/invoices_templ.go`
- **Purpose**: Compiled Go template implementation for the billing/invoices view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/billing/payments.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/billing/payments.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for billing/payments.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/billing/payments_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/billing/payments_templ.go`
- **Purpose**: Compiled Go template implementation for the billing/payments view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/billing/refunds.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/billing/refunds.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for billing/refunds.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/billing/refunds_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/billing/refunds_templ.go`
- **Purpose**: Compiled Go template implementation for the billing/refunds view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/billing/subnav.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/billing/subnav.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for billing/subnav.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/billing/subnav_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/billing/subnav_templ.go`
- **Purpose**: Compiled Go template implementation for the billing/subnav view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/calendar/calendar.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/calendar/calendar.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for calendar/calendar.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/calendar/calendar_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/calendar/calendar_templ.go`
- **Purpose**: Compiled Go template implementation for the calendar/calendar view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/clients/client_details.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/clients/client_details.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for clients/client_details.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/clients/client_details_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/clients/client_details_templ.go`
- **Purpose**: Compiled Go template implementation for the clients/client_details view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/clients/clients.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/clients/clients.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for clients/clients.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/clients/clients_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/clients/clients_templ.go`
- **Purpose**: Compiled Go template implementation for the clients/clients view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/contracts/contracts.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/contracts/contracts.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for contracts/contracts.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/contracts/contracts_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/contracts/contracts_templ.go`
- **Purpose**: Compiled Go template implementation for the contracts/contracts view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/dashboard.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/dashboard.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for pages/dashboard.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/dashboard/dashboard.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/dashboard/dashboard.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for dashboard/dashboard.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/dashboard/dashboard_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/dashboard/dashboard_templ.go`
- **Purpose**: Compiled Go template implementation for the dashboard/dashboard view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/dashboard_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/dashboard_templ.go`
- **Purpose**: Compiled Go template implementation for the pages/dashboard view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/deliverables/deliverables.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/deliverables/deliverables.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for deliverables/deliverables.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/deliverables/deliverables_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/deliverables/deliverables_templ.go`
- **Purpose**: Compiled Go template implementation for the deliverables/deliverables view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/errors.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/errors.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for pages/errors.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/errors_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/errors_templ.go`
- **Purpose**: Compiled Go template implementation for the pages/errors view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/followups/followups.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/followups/followups.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for followups/followups.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/followups/followups_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/followups/followups_templ.go`
- **Purpose**: Compiled Go template implementation for the followups/followups view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/forgot_password.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/forgot_password.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for pages/forgot_password.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/forgot_password_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/forgot_password_templ.go`
- **Purpose**: Compiled Go template implementation for the pages/forgot_password view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/inbox/inbox.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/inbox/inbox.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for inbox/inbox.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/inbox/inbox_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/inbox/inbox_templ.go`
- **Purpose**: Compiled Go template implementation for the inbox/inbox view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/leads/lead_details.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/leads/lead_details.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for leads/lead_details.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/leads/lead_details_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/leads/lead_details_templ.go`
- **Purpose**: Compiled Go template implementation for the leads/lead_details view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/leads/leads.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/leads/leads.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for leads/leads.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/leads/leads_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/leads/leads_templ.go`
- **Purpose**: Compiled Go template implementation for the leads/leads view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/orders/order_details.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/orders/order_details.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for orders/order_details.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/orders/order_details_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/orders/order_details_templ.go`
- **Purpose**: Compiled Go template implementation for the orders/order_details view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/orders/orders.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/orders/orders.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for orders/orders.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/orders/orders_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/orders/orders_templ.go`
- **Purpose**: Compiled Go template implementation for the orders/orders view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/projects/project_details.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/projects/project_details.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for projects/project_details.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/projects/project_details_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/projects/project_details_templ.go`
- **Purpose**: Compiled Go template implementation for the projects/project_details view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/projects/projects.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/projects/projects.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for projects/projects.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/projects/projects_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/projects/projects_templ.go`
- **Purpose**: Compiled Go template implementation for the projects/projects view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/reports/reports.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/reports/reports.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for reports/reports.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/reports/reports_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/reports/reports_templ.go`
- **Purpose**: Compiled Go template implementation for the reports/reports view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/sessions/sessions.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/sessions/sessions.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for sessions/sessions.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/sessions/sessions_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/sessions/sessions_templ.go`
- **Purpose**: Compiled Go template implementation for the sessions/sessions view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/settings/business_settings.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/settings/business_settings.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for settings/business_settings.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/settings/business_settings_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/settings/business_settings_templ.go`
- **Purpose**: Compiled Go template implementation for the settings/business_settings view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/settings/contract_templates.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/settings/contract_templates.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for settings/contractates.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/settings/contract_templates_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/settings/contract_templates_templ.go`
- **Purpose**: Compiled Go template implementation for the settings/contractates view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/settings/crew.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/settings/crew.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for settings/crew.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/settings/crew_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/settings/crew_templ.go`
- **Purpose**: Compiled Go template implementation for the settings/crew view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/settings/deliverable_presets.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/settings/deliverable_presets.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for settings/deliverable_presets.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/settings/deliverable_presets_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/settings/deliverable_presets_templ.go`
- **Purpose**: Compiled Go template implementation for the settings/deliverable_presets view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/settings/feedback_templates.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/settings/feedback_templates.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for settings/feedbackates.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/settings/feedback_templates_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/settings/feedback_templates_templ.go`
- **Purpose**: Compiled Go template implementation for the settings/feedbackates view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/settings/forms.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/settings/forms.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for settings/forms.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/settings/forms_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/settings/forms_templ.go`
- **Purpose**: Compiled Go template implementation for the settings/forms view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/settings/integrations.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/settings/integrations.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for settings/integrations.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/settings/integrations_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/settings/integrations_templ.go`
- **Purpose**: Compiled Go template implementation for the settings/integrations view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/settings/invoice_templates.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/settings/invoice_templates.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for settings/invoiceates.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/settings/invoice_templates_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/settings/invoice_templates_templ.go`
- **Purpose**: Compiled Go template implementation for the settings/invoiceates view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/settings/order_presets.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/settings/order_presets.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for settings/order_presets.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/settings/order_presets_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/settings/order_presets_templ.go`
- **Purpose**: Compiled Go template implementation for the settings/order_presets view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/settings/packages.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/settings/packages.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for settings/packages.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/settings/packages_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/settings/packages_templ.go`
- **Purpose**: Compiled Go template implementation for the settings/packages view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/settings/quotation_templates.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/settings/quotation_templates.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for settings/quotationates.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/settings/quotation_templates_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/settings/quotation_templates_templ.go`
- **Purpose**: Compiled Go template implementation for the settings/quotationates view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/settings/sequence_ids.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/settings/sequence_ids.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for settings/sequence_ids.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/settings/sequence_ids_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/settings/sequence_ids_templ.go`
- **Purpose**: Compiled Go template implementation for the settings/sequence_ids view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/settings/services.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/settings/services.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for settings/services.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/settings/services_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/settings/services_templ.go`
- **Purpose**: Compiled Go template implementation for the settings/services view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/settings/settings.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/settings/settings.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for settings/settings.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/settings/settings_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/settings/settings_templ.go`
- **Purpose**: Compiled Go template implementation for the settings/settings view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/settings/workflow.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/settings/workflow.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for settings/workflow.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/settings/workflow_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/settings/workflow_templ.go`
- **Purpose**: Compiled Go template implementation for the settings/workflow view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/signin.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/signin.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for pages/signin.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/signin_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/signin_templ.go`
- **Purpose**: Compiled Go template implementation for the pages/signin view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/signup.templ`
- **Relative Path**: `projects/photography_crm/ui/pages/signup.templ`
- **Purpose**: Templ component rendering the user interface, tables, cards, or forms for pages/signup.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/pages/signup_templ.go`
- **Relative Path**: `projects/photography_crm/ui/pages/signup_templ.go`
- **Purpose**: Compiled Go template implementation for the pages/signup view in Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/static/css/auth.css`
- **Relative Path**: `projects/photography_crm/ui/static/css/auth.css`
- **Purpose**: CSS stylesheet declaring custom styles, responsive layouts, and design tokens for auth.css.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/static/css/crm_shared.css`
- **Relative Path**: `projects/photography_crm/ui/static/css/crm_shared.css`
- **Purpose**: CSS stylesheet declaring custom styles, responsive layouts, and design tokens for crm_shared.css.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/static/css/dashboard.css`
- **Relative Path**: `projects/photography_crm/ui/static/css/dashboard.css`
- **Purpose**: CSS stylesheet declaring custom styles, responsive layouts, and design tokens for dashboard.css.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/static/css/globals.css`
- **Relative Path**: `projects/photography_crm/ui/static/css/globals.css`
- **Purpose**: CSS stylesheet declaring custom styles, responsive layouts, and design tokens for globals.css.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/static/css/pages/billing.css`
- **Relative Path**: `projects/photography_crm/ui/static/css/pages/billing.css`
- **Purpose**: CSS stylesheet declaring custom styles, responsive layouts, and design tokens for billing.css.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/static/css/pages/calendar.css`
- **Relative Path**: `projects/photography_crm/ui/static/css/pages/calendar.css`
- **Purpose**: CSS stylesheet declaring custom styles, responsive layouts, and design tokens for calendar.css.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/static/css/pages/leads.css`
- **Relative Path**: `projects/photography_crm/ui/static/css/pages/leads.css`
- **Purpose**: CSS stylesheet declaring custom styles, responsive layouts, and design tokens for leads.css.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/static/css/pages/projects.css`
- **Relative Path**: `projects/photography_crm/ui/static/css/pages/projects.css`
- **Purpose**: CSS stylesheet declaring custom styles, responsive layouts, and design tokens for projects.css.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/static/css/pages/sessions.css`
- **Relative Path**: `projects/photography_crm/ui/static/css/pages/sessions.css`
- **Purpose**: CSS stylesheet declaring custom styles, responsive layouts, and design tokens for sessions.css.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/static/css/pages/settings.css`
- **Relative Path**: `projects/photography_crm/ui/static/css/pages/settings.css`
- **Purpose**: CSS stylesheet declaring custom styles, responsive layouts, and design tokens for settings.css.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/static/css/pages/settings_extra.css`
- **Relative Path**: `projects/photography_crm/ui/static/css/pages/settings_extra.css`
- **Purpose**: CSS stylesheet declaring custom styles, responsive layouts, and design tokens for settings_extra.css.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/static/css/style.css`
- **Relative Path**: `projects/photography_crm/ui/static/css/style.css`
- **Purpose**: CSS stylesheet declaring custom styles, responsive layouts, and design tokens for style.css.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/static/embed.go`
- **Relative Path**: `projects/photography_crm/ui/static/embed.go`
- **Purpose**: Embeds static CSS stylesheets, JavaScript libraries, and assets into the Photography CRM binary.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/static/js/auth.js`
- **Relative Path**: `projects/photography_crm/ui/static/js/auth.js`
- **Purpose**: Client-side script or third-party vendor library (HTMX, Sortable, FullCalendar) powering auth.js.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/static/js/fullcalendar.min.js`
- **Relative Path**: `projects/photography_crm/ui/static/js/fullcalendar.min.js`
- **Purpose**: Client-side script or third-party vendor library (HTMX, Sortable, FullCalendar) powering fullcalendar.min.js.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/static/js/htmx.min.js`
- **Relative Path**: `projects/photography_crm/ui/static/js/htmx.min.js`
- **Purpose**: Client-side script or third-party vendor library (HTMX, Sortable, FullCalendar) powering htmx.min.js.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/static/js/kanban.js`
- **Relative Path**: `projects/photography_crm/ui/static/js/kanban.js`
- **Purpose**: Client-side script or third-party vendor library (HTMX, Sortable, FullCalendar) powering kanban.js.

#### `/home/desktop/especa/beta/projects/photography_crm/ui/static/js/sortable.min.js`
- **Relative Path**: `projects/photography_crm/ui/static/js/sortable.min.js`
- **Purpose**: Client-side script or third-party vendor library (HTMX, Sortable, FullCalendar) powering sortable.min.js.

---

### 8. CMS Service (`projects/cms`) (30 files)

#### `/home/desktop/especa/beta/projects/cms/.env.example`
- **Relative Path**: `projects/cms/.env.example`
- **Purpose**: Example environment configuration specifying ports and settings required for CMS.

#### `/home/desktop/especa/beta/projects/cms/go.mod`
- **Relative Path**: `projects/cms/go.mod`
- **Purpose**: Go module definition specifying dependencies for CMS.

#### `/home/desktop/especa/beta/projects/cms/go.sum`
- **Relative Path**: `projects/cms/go.sum`
- **Purpose**: Checksum lockfile locking exact dependency versions for CMS.

#### `/home/desktop/especa/beta/projects/cms/handlers/auth_handler.go`
- **Relative Path**: `projects/cms/handlers/auth_handler.go`
- **Purpose**: HTTP controllers processing user login, signup, session creation, and logout for CMS.

#### `/home/desktop/especa/beta/projects/cms/handlers/home_handler.go`
- **Relative Path**: `projects/cms/handlers/home_handler.go`
- **Purpose**: HTTP handler rendering the primary dashboard view and root index page for CMS.

#### `/home/desktop/especa/beta/projects/cms/main.go`
- **Relative Path**: `projects/cms/main.go`
- **Purpose**: Standalone HTTP server entrypoint initializing database pool, Chi router, and serving CMS.

#### `/home/desktop/especa/beta/projects/cms/routes/backend.go`
- **Relative Path**: `projects/cms/routes/backend.go`
- **Purpose**: Chi route definitions for backend REST API endpoints and data endpoints in CMS.

#### `/home/desktop/especa/beta/projects/cms/routes/frontend.go`
- **Relative Path**: `projects/cms/routes/frontend.go`
- **Purpose**: Chi route definitions mapping frontend HTML pages and static assets for CMS.

#### `/home/desktop/especa/beta/projects/cms/routes/router.go`
- **Relative Path**: `projects/cms/routes/router.go`
- **Purpose**: Top-level Chi router constructor configuring middlewares, frontend views, and API subrouters for CMS.

#### `/home/desktop/especa/beta/projects/cms/tests/auth_api_test.go`
- **Relative Path**: `projects/cms/tests/auth_api_test.go`
- **Purpose**: Integration test suite testing auth_api for the CMS service.

#### `/home/desktop/especa/beta/projects/cms/tests/auth_middleware_test.go`
- **Relative Path**: `projects/cms/tests/auth_middleware_test.go`
- **Purpose**: Integration test suite testing auth_middleware for the CMS service.

#### `/home/desktop/especa/beta/projects/cms/tests/frontend_auth_test.go`
- **Relative Path**: `projects/cms/tests/frontend_auth_test.go`
- **Purpose**: Integration test suite testing frontend_auth for the CMS service.

#### `/home/desktop/especa/beta/projects/cms/tests/middleware_test.go`
- **Relative Path**: `projects/cms/tests/middleware_test.go`
- **Purpose**: Integration test suite testing middleware for the CMS service.

#### `/home/desktop/especa/beta/projects/cms/tests/signup_api_test.go`
- **Relative Path**: `projects/cms/tests/signup_api_test.go`
- **Purpose**: Integration test suite testing signup_api for the CMS service.

#### `/home/desktop/especa/beta/projects/cms/ui/layouts/auth_layout.templ`
- **Relative Path**: `projects/cms/ui/layouts/auth_layout.templ`
- **Purpose**: Templ markup component defining the base HTML skeleton and navigation layout for CMS.

#### `/home/desktop/especa/beta/projects/cms/ui/layouts/auth_layout_templ.go`
- **Relative Path**: `projects/cms/ui/layouts/auth_layout_templ.go`
- **Purpose**: Compiled Go template code for the auth_layout layout in CMS.

#### `/home/desktop/especa/beta/projects/cms/ui/layouts/dashboard_layout.templ`
- **Relative Path**: `projects/cms/ui/layouts/dashboard_layout.templ`
- **Purpose**: Templ markup component defining the base HTML skeleton and navigation layout for CMS.

#### `/home/desktop/especa/beta/projects/cms/ui/layouts/dashboard_layout_templ.go`
- **Relative Path**: `projects/cms/ui/layouts/dashboard_layout_templ.go`
- **Purpose**: Compiled Go template code for the dashboard_layout layout in CMS.

#### `/home/desktop/especa/beta/projects/cms/ui/pages/auth.templ`
- **Relative Path**: `projects/cms/ui/pages/auth.templ`
- **Purpose**: Templ page component providing UI rendering and layout for auth in CMS.

#### `/home/desktop/especa/beta/projects/cms/ui/pages/auth_templ.go`
- **Relative Path**: `projects/cms/ui/pages/auth_templ.go`
- **Purpose**: Compiled Go code for auth page in CMS.

#### `/home/desktop/especa/beta/projects/cms/ui/pages/dashboard.templ`
- **Relative Path**: `projects/cms/ui/pages/dashboard.templ`
- **Purpose**: Templ page component providing UI rendering and layout for dashboard in CMS.

#### `/home/desktop/especa/beta/projects/cms/ui/pages/dashboard_templ.go`
- **Relative Path**: `projects/cms/ui/pages/dashboard_templ.go`
- **Purpose**: Compiled Go code for dashboard page in CMS.

#### `/home/desktop/especa/beta/projects/cms/ui/pages/errors.templ`
- **Relative Path**: `projects/cms/ui/pages/errors.templ`
- **Purpose**: Templ page component providing UI rendering and layout for errors in CMS.

#### `/home/desktop/especa/beta/projects/cms/ui/pages/errors_templ.go`
- **Relative Path**: `projects/cms/ui/pages/errors_templ.go`
- **Purpose**: Compiled Go code for errors page in CMS.

#### `/home/desktop/especa/beta/projects/cms/ui/static/css/auth.css`
- **Relative Path**: `projects/cms/ui/static/css/auth.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme rules for CMS (auth.css).

#### `/home/desktop/especa/beta/projects/cms/ui/static/css/dashboard.css`
- **Relative Path**: `projects/cms/ui/static/css/dashboard.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme rules for CMS (dashboard.css).

#### `/home/desktop/especa/beta/projects/cms/ui/static/css/globals.css`
- **Relative Path**: `projects/cms/ui/static/css/globals.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme rules for CMS (globals.css).

#### `/home/desktop/especa/beta/projects/cms/ui/static/css/style.css`
- **Relative Path**: `projects/cms/ui/static/css/style.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme rules for CMS (style.css).

#### `/home/desktop/especa/beta/projects/cms/ui/static/embed.go`
- **Relative Path**: `projects/cms/ui/static/embed.go`
- **Purpose**: Go embed directive embedding static CSS and JavaScript files into the CMS binary.

#### `/home/desktop/especa/beta/projects/cms/ui/static/js/auth.js`
- **Relative Path**: `projects/cms/ui/static/js/auth.js`
- **Purpose**: Client-side JavaScript managing client interactions and form validations for CMS.

---

### 9. Superadmin Service (`projects/superadmin`) (30 files)

#### `/home/desktop/especa/beta/projects/superadmin/.env.example`
- **Relative Path**: `projects/superadmin/.env.example`
- **Purpose**: Example environment configuration specifying ports and settings required for SUPERADMIN.

#### `/home/desktop/especa/beta/projects/superadmin/go.mod`
- **Relative Path**: `projects/superadmin/go.mod`
- **Purpose**: Go module definition specifying dependencies for SUPERADMIN.

#### `/home/desktop/especa/beta/projects/superadmin/go.sum`
- **Relative Path**: `projects/superadmin/go.sum`
- **Purpose**: Checksum lockfile locking exact dependency versions for SUPERADMIN.

#### `/home/desktop/especa/beta/projects/superadmin/handlers/auth_handler.go`
- **Relative Path**: `projects/superadmin/handlers/auth_handler.go`
- **Purpose**: HTTP controllers processing user login, signup, session creation, and logout for SUPERADMIN.

#### `/home/desktop/especa/beta/projects/superadmin/handlers/home_handler.go`
- **Relative Path**: `projects/superadmin/handlers/home_handler.go`
- **Purpose**: HTTP handler rendering the primary dashboard view and root index page for SUPERADMIN.

#### `/home/desktop/especa/beta/projects/superadmin/main.go`
- **Relative Path**: `projects/superadmin/main.go`
- **Purpose**: Standalone HTTP server entrypoint initializing database pool, Chi router, and serving SUPERADMIN.

#### `/home/desktop/especa/beta/projects/superadmin/routes/backend.go`
- **Relative Path**: `projects/superadmin/routes/backend.go`
- **Purpose**: Chi route definitions for backend REST API endpoints and data endpoints in SUPERADMIN.

#### `/home/desktop/especa/beta/projects/superadmin/routes/frontend.go`
- **Relative Path**: `projects/superadmin/routes/frontend.go`
- **Purpose**: Chi route definitions mapping frontend HTML pages and static assets for SUPERADMIN.

#### `/home/desktop/especa/beta/projects/superadmin/routes/router.go`
- **Relative Path**: `projects/superadmin/routes/router.go`
- **Purpose**: Top-level Chi router constructor configuring middlewares, frontend views, and API subrouters for SUPERADMIN.

#### `/home/desktop/especa/beta/projects/superadmin/tests/auth_api_test.go`
- **Relative Path**: `projects/superadmin/tests/auth_api_test.go`
- **Purpose**: Integration test suite testing auth_api for the SUPERADMIN service.

#### `/home/desktop/especa/beta/projects/superadmin/tests/auth_middleware_test.go`
- **Relative Path**: `projects/superadmin/tests/auth_middleware_test.go`
- **Purpose**: Integration test suite testing auth_middleware for the SUPERADMIN service.

#### `/home/desktop/especa/beta/projects/superadmin/tests/frontend_auth_test.go`
- **Relative Path**: `projects/superadmin/tests/frontend_auth_test.go`
- **Purpose**: Integration test suite testing frontend_auth for the SUPERADMIN service.

#### `/home/desktop/especa/beta/projects/superadmin/tests/middleware_test.go`
- **Relative Path**: `projects/superadmin/tests/middleware_test.go`
- **Purpose**: Integration test suite testing middleware for the SUPERADMIN service.

#### `/home/desktop/especa/beta/projects/superadmin/tests/signup_api_test.go`
- **Relative Path**: `projects/superadmin/tests/signup_api_test.go`
- **Purpose**: Integration test suite testing signup_api for the SUPERADMIN service.

#### `/home/desktop/especa/beta/projects/superadmin/ui/layouts/auth_layout.templ`
- **Relative Path**: `projects/superadmin/ui/layouts/auth_layout.templ`
- **Purpose**: Templ markup component defining the base HTML skeleton and navigation layout for SUPERADMIN.

#### `/home/desktop/especa/beta/projects/superadmin/ui/layouts/auth_layout_templ.go`
- **Relative Path**: `projects/superadmin/ui/layouts/auth_layout_templ.go`
- **Purpose**: Compiled Go template code for the auth_layout layout in SUPERADMIN.

#### `/home/desktop/especa/beta/projects/superadmin/ui/layouts/dashboard_layout.templ`
- **Relative Path**: `projects/superadmin/ui/layouts/dashboard_layout.templ`
- **Purpose**: Templ markup component defining the base HTML skeleton and navigation layout for SUPERADMIN.

#### `/home/desktop/especa/beta/projects/superadmin/ui/layouts/dashboard_layout_templ.go`
- **Relative Path**: `projects/superadmin/ui/layouts/dashboard_layout_templ.go`
- **Purpose**: Compiled Go template code for the dashboard_layout layout in SUPERADMIN.

#### `/home/desktop/especa/beta/projects/superadmin/ui/pages/auth.templ`
- **Relative Path**: `projects/superadmin/ui/pages/auth.templ`
- **Purpose**: Templ page component providing UI rendering and layout for auth in SUPERADMIN.

#### `/home/desktop/especa/beta/projects/superadmin/ui/pages/auth_templ.go`
- **Relative Path**: `projects/superadmin/ui/pages/auth_templ.go`
- **Purpose**: Compiled Go code for auth page in SUPERADMIN.

#### `/home/desktop/especa/beta/projects/superadmin/ui/pages/dashboard.templ`
- **Relative Path**: `projects/superadmin/ui/pages/dashboard.templ`
- **Purpose**: Templ page component providing UI rendering and layout for dashboard in SUPERADMIN.

#### `/home/desktop/especa/beta/projects/superadmin/ui/pages/dashboard_templ.go`
- **Relative Path**: `projects/superadmin/ui/pages/dashboard_templ.go`
- **Purpose**: Compiled Go code for dashboard page in SUPERADMIN.

#### `/home/desktop/especa/beta/projects/superadmin/ui/pages/errors.templ`
- **Relative Path**: `projects/superadmin/ui/pages/errors.templ`
- **Purpose**: Templ page component providing UI rendering and layout for errors in SUPERADMIN.

#### `/home/desktop/especa/beta/projects/superadmin/ui/pages/errors_templ.go`
- **Relative Path**: `projects/superadmin/ui/pages/errors_templ.go`
- **Purpose**: Compiled Go code for errors page in SUPERADMIN.

#### `/home/desktop/especa/beta/projects/superadmin/ui/static/css/auth.css`
- **Relative Path**: `projects/superadmin/ui/static/css/auth.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme rules for SUPERADMIN (auth.css).

#### `/home/desktop/especa/beta/projects/superadmin/ui/static/css/dashboard.css`
- **Relative Path**: `projects/superadmin/ui/static/css/dashboard.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme rules for SUPERADMIN (dashboard.css).

#### `/home/desktop/especa/beta/projects/superadmin/ui/static/css/globals.css`
- **Relative Path**: `projects/superadmin/ui/static/css/globals.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme rules for SUPERADMIN (globals.css).

#### `/home/desktop/especa/beta/projects/superadmin/ui/static/css/style.css`
- **Relative Path**: `projects/superadmin/ui/static/css/style.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme rules for SUPERADMIN (style.css).

#### `/home/desktop/especa/beta/projects/superadmin/ui/static/embed.go`
- **Relative Path**: `projects/superadmin/ui/static/embed.go`
- **Purpose**: Go embed directive embedding static CSS and JavaScript files into the SUPERADMIN binary.

#### `/home/desktop/especa/beta/projects/superadmin/ui/static/js/auth.js`
- **Relative Path**: `projects/superadmin/ui/static/js/auth.js`
- **Purpose**: Client-side JavaScript managing client interactions and form validations for SUPERADMIN.

---

### 10. Testing & API Collections (`bruno`) (4 files)

#### `/home/desktop/especa/beta/bruno/API.txt`
- **Relative Path**: `bruno/API.txt`
- **Purpose**: Documentation and notes describing Bruno API collections, authentication headers, and test execution.

#### `/home/desktop/especa/beta/bruno/photography/create roles.yml`
- **Relative Path**: `bruno/photography/create roles.yml`
- **Purpose**: Bruno API request definition to seed or create user roles in the Photography CRM authentication system.

#### `/home/desktop/especa/beta/bruno/photography/opencollection.yml`
- **Relative Path**: `bruno/photography/opencollection.yml`
- **Purpose**: Bruno collection configuration defining collection-wide variables and authentication tokens for Photography CRM.

#### `/home/desktop/especa/beta/bruno/photography/signup.yml`
- **Relative Path**: `bruno/photography/signup.yml`
- **Purpose**: Bruno API request definition for automated user registration and credential testing.

---

### 11. Agent Rules, Skills & Editor Config (`.agents` & `.vscode`) (4 files)

#### `/home/desktop/especa/beta/.agents/rules/rule.md`
- **Relative Path**: `.agents/rules/rule.md`
- **Purpose**: Defines workspace-wide developer rules, forbidding daemon server backgrounding and requiring ports to stay clean.

#### `/home/desktop/especa/beta/.agents/skills/dashboard-designer/SKILL.md`
- **Relative Path**: `.agents/skills/dashboard-designer/SKILL.md`
- **Purpose**: Skill reference detailing dashboard UX/UI layout rules, visual hierarchy, color tokens, and Lucide icon guidelines.

#### `/home/desktop/especa/beta/.vscode/launch.json`
- **Relative Path**: `.vscode/launch.json`
- **Purpose**: VS Code and Cursor runtime debug configuration launching Delve for Photography CRM and the multi-service runner.

#### `/home/desktop/especa/beta/.vscode/settings.json`
- **Relative Path**: `.vscode/settings.json`
- **Purpose**: Workspace editor settings configuring the Go language server (gopls) and formatting rules.

---
