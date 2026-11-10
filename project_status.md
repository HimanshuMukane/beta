# BETA Project: Architecture, Ideology & Comprehensive Repository Status

> **Context Document for AI Models & Developers**
> This document provides a complete architectural overview, coding standards, language & library stacks, in-depth documentation of all active platform applications (**CMS**, **Photography CRM**, **Superadmin**, **Especa Finance**, **Especa CRM**, and **HRMS**), and an exhaustive file-by-file catalog of every single source file in the `BETA` codebase.

---

## 1. Overall Project Status & Ideology

### 1.1 Architecture & Multi-Module Workspace
The `BETA` project is organized as a high-performance Go multi-module workspace (`go.work`) adopting **Clean Architecture** and modern **Hypermedia-Driven UI** principles. Rather than maintaining heavy Single Page Application (SPA) frontend frameworks that require complex build chains, `BETA` leverages **Server-Side Rendered (SSR) Go templates with A-H/Templ**, **HTMX** for selective DOM swaps, **Vanilla JavaScript** for targeted micro-interactions, and pure CSS for sleek obsidian-dark/enterprise-slate themes.

The workspace unifies multiple distinct applications and shared packages:
- **`cmd/runner`**: Multi-service runner and supervisor that can boot all platform micro-services concurrently with graceful shutdown.
- **`projects/core`**: Domain data layer, database connection pooling (`DBPool`), Sqlc-generated typesafe queries, repository patterns, and cross-cutting authentication/session middleware.
- **`ui`**: Centralized, reusable Templ UI component library (buttons, cards, badges, dialogs, charts, tables, messages, spotlight search).
- **`projects/superadmin`**: Multi-tenant administrative portal for managing users, organizations, subscriptions, and platform-wide configurations.
- **`projects/cms`**: Content management service sharing the unified core authentication and schema.
- **`projects/photography_crm`**: Flagship production SaaS CRM tailored for professional photographers and studios.
- **`projects/especa_crm`**: General customer relationship management module tracking leads, deals, contacts, and quotes.
- **`projects/especa_finance`**: Double-entry accounting, financial ledgers, invoicing, chart of accounts, and project financial analytics.
- **`projects/hrms`**: Human Resource Management System handling employees, departments, leave management, and attendance.
- **`tests`**: Centralized integration, unit, and end-to-end test suites verifying security, session revocation, database pools, and business flows.

### 1.2 Technology Stack & Modules Used
| Layer | Technologies & Libraries |
|---|---|
| **Core Language** | Go 1.26+ (`go.work` multi-module structure) |
| **HTTP Routing** | `github.com/go-chi/chi/v5` (modular sub-routing, clean middleware chaining) |
| **Templating Engine** | `github.com/a-h/templ` (type-safe, pre-compiled Go HTML templates) |
| **Database & Storage** | MySQL 8.0 (Docker containerized) & SQLite 3 fallback via `modernc.org/sqlite` |
| **Type-Safe SQL** | `sqlc` (compiles raw SQL into idiomatic Go structs and queries) |
| **Security & Crypto** | Argon2id & `golang.org/x/crypto/bcrypt`, HMAC-SHA256 session cookie encryption, CSRF protection |
| **Frontend Hypermedia** | HTMX (`htmx.min.js`), SortableJS (`sortable.min.js`), FullCalendar (`fullcalendar.min.js`) |
| **Design & Icons** | Custom Obsidian Dark & Enterprise Slate CSS design tokens, Lucide SVG icons |
| **Configuration** | `github.com/joho/godotenv` with fallback environment defaults |
| **API Testing** | Bruno collection runner (`.yml` requests) and standard Go `testing` package |

### 1.3 Coding Style & Engineering Principles
1. **Strict Simplicity**: Prioritize the simplest implementation that correctly meets the requirement. Keep code readable for any engineer.
2. **Separation of Concerns**: Each service cleanly delineates:
   - `routes/router.go` or `routes/routes.go`: Middlewares and top-level routing assembly.
   - `routes/frontend.go`: Page endpoints and HTML template rendering.
   - `routes/backend.go`: JSON REST APIs and mutation endpoints.
   - `handlers/`: Request parsing, service calls, and view model preparation.
   - `ui/pages/` and `ui/layouts/`: Pure `.templ` presentation code.
3. **Compiled Static Assets**: All CSS, JavaScript, and SVG assets are bundled into Go binaries using `//go:embed` via `embed.go`, eliminating missing file issues in production.
4. **Safe Connection Pooling (`DBPool`)**: Enforces connection limits, ping verification, and multi-tenant routing with dynamic environment swapping support.
5. **Server & Port Safety Rules**: Workspace execution rules forbid running unmonitored background daemons; ports remain clean when tests or builds complete.

### 1.4 Port & Service Allocation Map
The supervisor `cmd/runner/main.go` assigns dedicated ports for each platform microservice:
| Service | Project Directory | Default Port | Session Cookie Name | Database Name |
|---|---|---|---|---|
| **CMS** | `projects/cms` | `:8081` | `cms_session` | `cms_db` |
| **Photography CRM** | `projects/photography_crm` | `:8082` | `photography_session` | `photography_crm_db` |
| **Superadmin** | `projects/superadmin` | `:8083` | `superadmin_session` | `superadmin_db` |
| **Finance** | `projects/especa_finance` | `:8084` | `finance_session` | `finance_db` |
| **Especa CRM** | `projects/especa_crm` | `:8086` | `crm_session` | `especa_crm_db` |
| **HRMS** | `projects/hrms` | `:8087` | `hrms_session` | `hrms_db` |

---

## 2. Platform Services Overview

### 2.1 Photography CRM (`projects/photography_crm`)
End-to-end studio management for commercial, portrait, and wedding photographers, powered by a modern React SPA (`frontend/apps/photography_crm`) and a direct-SQL Go Chi backend (`:8082`):
- **Dual-Portal Architecture**: Distinct role-based interfaces for studio leadership and field crew:
  - **Admin Studio Portal**: Full management of Leads, Projects, Sessions, Orders, and Billings; 8 operational & financial reports; 12 customizable studio settings panels.
  - **Team Member / User Portal**: Streamlined portal for photographers, cinematographers, and editors (`/home` agenda, `/shoots` shoot details & timings, `/tasks` Kanban production matrix, `/deliverables` culling and retouching tracking).
- **Leads & Pipeline Management (`/leads`)**: Kanban drag-and-drop pipeline and table views with customizable stages, source attribution, and one-click lead-to-client conversion.
- **Projects & Shoots (`/projects`)**: Milestone tracking, multi-day shoot segments, crew assignments, session scheduling, and shoot expenditure tracking.
- **Sessions & Bookings (`/sessions`)**: Photoshoot session timetable, shoot timing slots, crew allocation, and session status drawer.
- **Merchandise & Orders (`/orders`)**: Client orders for photobooks, frames, canvas wraps, and raw footage USBs.
- **Billing & Financials (`/billings`)**: Double-entry financial tracking for invoices, client payments, operational expenses, cancellations, and refunds.
- **Studio Configuration (`/settings/*`)**: Workflows, master data, crew roster, sequence ID generators, service presets, package bundles, order presets, deliverable presets, contract templates, and feedback surveys.
- **Telemetry & Error Logging**: Centralized error reporting (`POST /api/logs/frontend`) connected to React Error Boundaries.

### 2.2 Especa Finance (`projects/especa_finance`)
Comprehensive financial and accounting management:
- **Chart of Accounts (`/accounts`)**: Asset, liability, equity, revenue, and expense account hierarchy.
- **General Ledger & Journal Entries (`/journal`)**: Double-entry bookkeeping with strict balance validation.
- **Invoicing & Billing (`/invoices`)**: Invoice issuance, payment tracking, and balance tracking.
- **Project Financials (`/project-financials`)**: Revenue, cost allocations, and margin calculation per project.
- **Subscriptions (`/subscriptions`)**: Recurring billing contracts and renewal schedules.

### 2.3 Especa CRM (`projects/especa_crm`)
General-purpose CRM tailored for client and deal pipelines:
- **Deals & Opportunities (`/deals`)**: Value tracking, pipeline stages, probability forecasting.
- **Lead Capture (`/leads`)**: Prospect acquisition, qualification scoring, and assignment.
- **Quotes & Proposals (`/quotes`)**: Quote generation and client acceptance flows.

### 2.4 HRMS (`projects/hrms`)
Human Resource Management System:
- **Employee Directory (`/employees`)**: Profile data, job titles, department assignments, salary tiers.
- **Department Structure (`/departments`)**: Organizational units and managerial hierarchies.
- **Attendance & Time Logs (`/attendance`)**: Clock-in/out records, timesheets, and overtime tracking.
- **Leave Management (`/leaves`)**: Leave request submissions, balance tracking, and manager approvals.

### 2.5 Superadmin (`projects/superadmin`)
Platform-wide administration:
- **Tenant Management**: Organization provisioning, schema isolation, and domain routing.
- **User & Role Administration**: Global permission assignments, role definitions, and access revocation.
- **Platform Analytics**: Cross-tenant health metrics and audit logs.

### 2.6 CMS (`projects/cms`)
Headless and rendered content management:
- **Article & Page Publishing**: Draft, review, and publication lifecycles.
- **Developer API**: Headless REST endpoints providing content feeds to external web frontends.

---

## 3. Exhaustive File-by-File Repository Catalog

Below is the complete, exhaustive catalog of every file across the `BETA` project. Each entry contains the **relative file path**, **absolute path**, and a concise **role and implementation description**.


### 1. Root Workspace & Environment Configuration (6 files)

#### `/home/desktop/especa/beta/.env.example`
- **Relative Path**: `.env.example`
- **Purpose**: Example environment configuration template with default ports, credentials, and connection strings.

#### `/home/desktop/especa/beta/.gitignore`
- **Relative Path**: `.gitignore`
- **Purpose**: Specifies files and patterns ignored by Git, including compiled binaries, database journals, logs, and sensitive environment files.

#### `/home/desktop/especa/beta/docker-compose.yml`
- **Relative Path**: `docker-compose.yml`
- **Purpose**: Docker Compose orchestration file configuring MySQL 8.0 container and phpMyAdmin dashboard.

#### `/home/desktop/especa/beta/go.work`
- **Relative Path**: `go.work`
- **Purpose**: Defines the Go multi-module workspace unifying runner, ui, core, cms, superadmin, photography_crm, especa_crm, especa_finance, hrms, and tests.

#### `/home/desktop/especa/beta/go.work.sum`
- **Relative Path**: `go.work.sum`
- **Purpose**: Cryptographic checksum verification lockfile for module dependencies across the entire Go workspace.

#### `/home/desktop/especa/beta/project_status.md`
- **Relative Path**: `project_status.md`
- **Purpose**: Context document and exhaustive architecture status catalog for the entire BETA repository.


### 2. Container & Database Initialization (`docker`) (1 files)

#### `/home/desktop/especa/beta/docker/volumes/mysql/initdb/01_schema.sql`
- **Relative Path**: `docker/volumes/mysql/initdb/01_schema.sql`
- **Purpose**: Initial database DDL provisioning multi-tenant platform schemas, identity tables, and service databases.


### 3. Multi-Service Runner (`cmd/runner`) (3 files)

#### `/home/desktop/especa/beta/cmd/runner/go.mod`
- **Relative Path**: `cmd/runner/go.mod`
- **Purpose**: Go module manifest declaring dependencies for the cmd/runner multi-service supervisor.

#### `/home/desktop/especa/beta/cmd/runner/go.sum`
- **Relative Path**: `cmd/runner/go.sum`
- **Purpose**: Checksum lockfile ensuring reproducible builds for the cmd/runner supervisor service.

#### `/home/desktop/especa/beta/cmd/runner/main.go`
- **Relative Path**: `cmd/runner/main.go`
- **Purpose**: Unified multi-service supervisor application booting CMS (:8081), Photography CRM (:8082), Superadmin (:8083), Finance (:8084), Especa CRM (:8086), and HRMS (:8087) with graceful shutdown.


### 4. Shared UI Component Library (`ui`) (276 files)

#### `/home/desktop/especa/beta/ui/baseui/lifecycle.js`
- **Relative Path**: `ui/baseui/lifecycle.js`
- **Purpose**: Client-side component lifecycle orchestrator managing dynamic DOM mounts and state bindings.

#### `/home/desktop/especa/beta/ui/components/buttons/button.templ`
- **Relative Path**: `ui/components/buttons/button.templ`
- **Purpose**: Templ component rendering reusable UI element for button.templ in the buttons category.

#### `/home/desktop/especa/beta/ui/components/buttons/button_sm.templ`
- **Relative Path**: `ui/components/buttons/button_sm.templ`
- **Purpose**: Templ component rendering reusable UI element for button_sm.templ in the buttons category.

#### `/home/desktop/especa/beta/ui/components/buttons/button_sm_templ.go`
- **Relative Path**: `ui/components/buttons/button_sm_templ.go`
- **Purpose**: Compiled Go code generated by templ for button_sm.templ.

#### `/home/desktop/especa/beta/ui/components/buttons/button_templ.go`
- **Relative Path**: `ui/components/buttons/button_templ.go`
- **Purpose**: Compiled Go code generated by templ for button.templ.

#### `/home/desktop/especa/beta/ui/components/buttons/button_types.go`
- **Relative Path**: `ui/components/buttons/button_types.go`
- **Purpose**: Go types and view model definitions for the buttons component library.

#### `/home/desktop/especa/beta/ui/components/buttons/shadcn_button.templ`
- **Relative Path**: `ui/components/buttons/shadcn_button.templ`
- **Purpose**: Templ component rendering reusable UI element for shadcn_button.templ in the buttons category.

#### `/home/desktop/especa/beta/ui/components/buttons/shadcn_button_templ.go`
- **Relative Path**: `ui/components/buttons/shadcn_button_templ.go`
- **Purpose**: Compiled Go code generated by templ for shadcn_button.templ.

#### `/home/desktop/especa/beta/ui/components/cards/card.templ`
- **Relative Path**: `ui/components/cards/card.templ`
- **Purpose**: Templ component rendering reusable UI element for card.templ in the cards category.

#### `/home/desktop/especa/beta/ui/components/cards/card_action.templ`
- **Relative Path**: `ui/components/cards/card_action.templ`
- **Purpose**: Templ component rendering reusable UI element for card_action.templ in the cards category.

#### `/home/desktop/especa/beta/ui/components/cards/card_action_templ.go`
- **Relative Path**: `ui/components/cards/card_action_templ.go`
- **Purpose**: Compiled Go code generated by templ for card_action.templ.

#### `/home/desktop/especa/beta/ui/components/cards/card_content.templ`
- **Relative Path**: `ui/components/cards/card_content.templ`
- **Purpose**: Templ component rendering reusable UI element for card_content.templ in the cards category.

#### `/home/desktop/especa/beta/ui/components/cards/card_content_templ.go`
- **Relative Path**: `ui/components/cards/card_content_templ.go`
- **Purpose**: Compiled Go code generated by templ for card_content.templ.

#### `/home/desktop/especa/beta/ui/components/cards/card_description.templ`
- **Relative Path**: `ui/components/cards/card_description.templ`
- **Purpose**: Templ component rendering reusable UI element for card_description.templ in the cards category.

#### `/home/desktop/especa/beta/ui/components/cards/card_description_templ.go`
- **Relative Path**: `ui/components/cards/card_description_templ.go`
- **Purpose**: Compiled Go code generated by templ for card_description.templ.

#### `/home/desktop/especa/beta/ui/components/cards/card_footer.templ`
- **Relative Path**: `ui/components/cards/card_footer.templ`
- **Purpose**: Templ component rendering reusable UI element for card_footer.templ in the cards category.

#### `/home/desktop/especa/beta/ui/components/cards/card_footer_templ.go`
- **Relative Path**: `ui/components/cards/card_footer_templ.go`
- **Purpose**: Compiled Go code generated by templ for card_footer.templ.

#### `/home/desktop/especa/beta/ui/components/cards/card_header.templ`
- **Relative Path**: `ui/components/cards/card_header.templ`
- **Purpose**: Templ component rendering reusable UI element for card_header.templ in the cards category.

#### `/home/desktop/especa/beta/ui/components/cards/card_header_templ.go`
- **Relative Path**: `ui/components/cards/card_header_templ.go`
- **Purpose**: Compiled Go code generated by templ for card_header.templ.

#### `/home/desktop/especa/beta/ui/components/cards/card_templ.go`
- **Relative Path**: `ui/components/cards/card_templ.go`
- **Purpose**: Compiled Go code generated by templ for card.templ.

#### `/home/desktop/especa/beta/ui/components/cards/card_title.templ`
- **Relative Path**: `ui/components/cards/card_title.templ`
- **Purpose**: Templ component rendering reusable UI element for card_title.templ in the cards category.

#### `/home/desktop/especa/beta/ui/components/cards/card_title_templ.go`
- **Relative Path**: `ui/components/cards/card_title_templ.go`
- **Purpose**: Compiled Go code generated by templ for card_title.templ.

#### `/home/desktop/especa/beta/ui/components/cards/card_types.go`
- **Relative Path**: `ui/components/cards/card_types.go`
- **Purpose**: Go types and view model definitions for the cards component library.

#### `/home/desktop/especa/beta/ui/components/chart/area.templ`
- **Relative Path**: `ui/components/chart/area.templ`
- **Purpose**: Templ component rendering reusable UI element for area.templ in the chart category.

#### `/home/desktop/especa/beta/ui/components/chart/area_chart.templ`
- **Relative Path**: `ui/components/chart/area_chart.templ`
- **Purpose**: Templ component rendering reusable UI element for area_chart.templ in the chart category.

#### `/home/desktop/especa/beta/ui/components/chart/area_chart_templ.go`
- **Relative Path**: `ui/components/chart/area_chart_templ.go`
- **Purpose**: Compiled Go code generated by templ for area_chart.templ.

#### `/home/desktop/especa/beta/ui/components/chart/area_templ.go`
- **Relative Path**: `ui/components/chart/area_templ.go`
- **Purpose**: Compiled Go code generated by templ for area.templ.

#### `/home/desktop/especa/beta/ui/components/chart/bar.templ`
- **Relative Path**: `ui/components/chart/bar.templ`
- **Purpose**: Templ component rendering reusable UI element for bar.templ in the chart category.

#### `/home/desktop/especa/beta/ui/components/chart/bar_chart.templ`
- **Relative Path**: `ui/components/chart/bar_chart.templ`
- **Purpose**: Templ component rendering reusable UI element for bar_chart.templ in the chart category.

#### `/home/desktop/especa/beta/ui/components/chart/bar_chart_templ.go`
- **Relative Path**: `ui/components/chart/bar_chart_templ.go`
- **Purpose**: Compiled Go code generated by templ for bar_chart.templ.

#### `/home/desktop/especa/beta/ui/components/chart/bar_templ.go`
- **Relative Path**: `ui/components/chart/bar_templ.go`
- **Purpose**: Compiled Go code generated by templ for bar.templ.

#### `/home/desktop/especa/beta/ui/components/chart/cartesian_grid.templ`
- **Relative Path**: `ui/components/chart/cartesian_grid.templ`
- **Purpose**: Templ component rendering reusable UI element for cartesian_grid.templ in the chart category.

#### `/home/desktop/especa/beta/ui/components/chart/cartesian_grid_templ.go`
- **Relative Path**: `ui/components/chart/cartesian_grid_templ.go`
- **Purpose**: Compiled Go code generated by templ for cartesian_grid.templ.

#### `/home/desktop/especa/beta/ui/components/chart/cell.templ`
- **Relative Path**: `ui/components/chart/cell.templ`
- **Purpose**: Templ component rendering reusable UI element for cell.templ in the chart category.

#### `/home/desktop/especa/beta/ui/components/chart/cell_templ.go`
- **Relative Path**: `ui/components/chart/cell_templ.go`
- **Purpose**: Compiled Go code generated by templ for cell.templ.

#### `/home/desktop/especa/beta/ui/components/chart/chart.js`
- **Relative Path**: `ui/components/chart/chart.js`
- **Purpose**: Client-side JavaScript managing client interactions, DOM operations, or library integrations.

#### `/home/desktop/especa/beta/ui/components/chart/chart_output.templ`
- **Relative Path**: `ui/components/chart/chart_output.templ`
- **Purpose**: Templ component rendering reusable UI element for chart_output.templ in the chart category.

#### `/home/desktop/especa/beta/ui/components/chart/chart_output_templ.go`
- **Relative Path**: `ui/components/chart/chart_output_templ.go`
- **Purpose**: Compiled Go code generated by templ for chart_output.templ.

#### `/home/desktop/especa/beta/ui/components/chart/container.templ`
- **Relative Path**: `ui/components/chart/container.templ`
- **Purpose**: Templ component rendering reusable UI element for container.templ in the chart category.

#### `/home/desktop/especa/beta/ui/components/chart/container_templ.go`
- **Relative Path**: `ui/components/chart/container_templ.go`
- **Purpose**: Compiled Go code generated by templ for container.templ.

#### `/home/desktop/especa/beta/ui/components/chart/defs.templ`
- **Relative Path**: `ui/components/chart/defs.templ`
- **Purpose**: Templ component rendering reusable UI element for defs.templ in the chart category.

#### `/home/desktop/especa/beta/ui/components/chart/defs_templ.go`
- **Relative Path**: `ui/components/chart/defs_templ.go`
- **Purpose**: Compiled Go code generated by templ for defs.templ.

#### `/home/desktop/especa/beta/ui/components/chart/label.templ`
- **Relative Path**: `ui/components/chart/label.templ`
- **Purpose**: Templ component rendering reusable UI element for label.templ in the chart category.

#### `/home/desktop/especa/beta/ui/components/chart/label_list.templ`
- **Relative Path**: `ui/components/chart/label_list.templ`
- **Purpose**: Templ component rendering reusable UI element for label_list.templ in the chart category.

#### `/home/desktop/especa/beta/ui/components/chart/label_list_templ.go`
- **Relative Path**: `ui/components/chart/label_list_templ.go`
- **Purpose**: Compiled Go code generated by templ for label_list.templ.

#### `/home/desktop/especa/beta/ui/components/chart/label_templ.go`
- **Relative Path**: `ui/components/chart/label_templ.go`
- **Purpose**: Compiled Go code generated by templ for label.templ.

#### `/home/desktop/especa/beta/ui/components/chart/legend.templ`
- **Relative Path**: `ui/components/chart/legend.templ`
- **Purpose**: Templ component rendering reusable UI element for legend.templ in the chart category.

#### `/home/desktop/especa/beta/ui/components/chart/legend_content.templ`
- **Relative Path**: `ui/components/chart/legend_content.templ`
- **Purpose**: Templ component rendering reusable UI element for legend_content.templ in the chart category.

#### `/home/desktop/especa/beta/ui/components/chart/legend_content_templ.go`
- **Relative Path**: `ui/components/chart/legend_content_templ.go`
- **Purpose**: Compiled Go code generated by templ for legend_content.templ.

#### `/home/desktop/especa/beta/ui/components/chart/legend_templ.go`
- **Relative Path**: `ui/components/chart/legend_templ.go`
- **Purpose**: Compiled Go code generated by templ for legend.templ.

#### `/home/desktop/especa/beta/ui/components/chart/line.templ`
- **Relative Path**: `ui/components/chart/line.templ`
- **Purpose**: Templ component rendering reusable UI element for line.templ in the chart category.

#### `/home/desktop/especa/beta/ui/components/chart/line_chart.templ`
- **Relative Path**: `ui/components/chart/line_chart.templ`
- **Purpose**: Templ component rendering reusable UI element for line_chart.templ in the chart category.

#### `/home/desktop/especa/beta/ui/components/chart/line_chart_templ.go`
- **Relative Path**: `ui/components/chart/line_chart_templ.go`
- **Purpose**: Compiled Go code generated by templ for line_chart.templ.

#### `/home/desktop/especa/beta/ui/components/chart/line_templ.go`
- **Relative Path**: `ui/components/chart/line_templ.go`
- **Purpose**: Compiled Go code generated by templ for line.templ.

#### `/home/desktop/especa/beta/ui/components/chart/pie.templ`
- **Relative Path**: `ui/components/chart/pie.templ`
- **Purpose**: Templ component rendering reusable UI element for pie.templ in the chart category.

#### `/home/desktop/especa/beta/ui/components/chart/pie_chart.templ`
- **Relative Path**: `ui/components/chart/pie_chart.templ`
- **Purpose**: Templ component rendering reusable UI element for pie_chart.templ in the chart category.

#### `/home/desktop/especa/beta/ui/components/chart/pie_chart_templ.go`
- **Relative Path**: `ui/components/chart/pie_chart_templ.go`
- **Purpose**: Compiled Go code generated by templ for pie_chart.templ.

#### `/home/desktop/especa/beta/ui/components/chart/pie_templ.go`
- **Relative Path**: `ui/components/chart/pie_templ.go`
- **Purpose**: Compiled Go code generated by templ for pie.templ.

#### `/home/desktop/especa/beta/ui/components/chart/polar_angle_axis.templ`
- **Relative Path**: `ui/components/chart/polar_angle_axis.templ`
- **Purpose**: Templ component rendering reusable UI element for polar_angle_axis.templ in the chart category.

#### `/home/desktop/especa/beta/ui/components/chart/polar_angle_axis_templ.go`
- **Relative Path**: `ui/components/chart/polar_angle_axis_templ.go`
- **Purpose**: Compiled Go code generated by templ for polar_angle_axis.templ.

#### `/home/desktop/especa/beta/ui/components/chart/polar_grid.templ`
- **Relative Path**: `ui/components/chart/polar_grid.templ`
- **Purpose**: Templ component rendering reusable UI element for polar_grid.templ in the chart category.

#### `/home/desktop/especa/beta/ui/components/chart/polar_grid_templ.go`
- **Relative Path**: `ui/components/chart/polar_grid_templ.go`
- **Purpose**: Compiled Go code generated by templ for polar_grid.templ.

#### `/home/desktop/especa/beta/ui/components/chart/polar_radius_axis.templ`
- **Relative Path**: `ui/components/chart/polar_radius_axis.templ`
- **Purpose**: Templ component rendering reusable UI element for polar_radius_axis.templ in the chart category.

#### `/home/desktop/especa/beta/ui/components/chart/polar_radius_axis_templ.go`
- **Relative Path**: `ui/components/chart/polar_radius_axis_templ.go`
- **Purpose**: Compiled Go code generated by templ for polar_radius_axis.templ.

#### `/home/desktop/especa/beta/ui/components/chart/radar.templ`
- **Relative Path**: `ui/components/chart/radar.templ`
- **Purpose**: Templ component rendering reusable UI element for radar.templ in the chart category.

#### `/home/desktop/especa/beta/ui/components/chart/radar_chart.templ`
- **Relative Path**: `ui/components/chart/radar_chart.templ`
- **Purpose**: Templ component rendering reusable UI element for radar_chart.templ in the chart category.

#### `/home/desktop/especa/beta/ui/components/chart/radar_chart_templ.go`
- **Relative Path**: `ui/components/chart/radar_chart_templ.go`
- **Purpose**: Compiled Go code generated by templ for radar_chart.templ.

#### `/home/desktop/especa/beta/ui/components/chart/radar_templ.go`
- **Relative Path**: `ui/components/chart/radar_templ.go`
- **Purpose**: Compiled Go code generated by templ for radar.templ.

#### `/home/desktop/especa/beta/ui/components/chart/radial_bar.templ`
- **Relative Path**: `ui/components/chart/radial_bar.templ`
- **Purpose**: Templ component rendering reusable UI element for radial_bar.templ in the chart category.

#### `/home/desktop/especa/beta/ui/components/chart/radial_bar_chart.templ`
- **Relative Path**: `ui/components/chart/radial_bar_chart.templ`
- **Purpose**: Templ component rendering reusable UI element for radial_bar_chart.templ in the chart category.

#### `/home/desktop/especa/beta/ui/components/chart/radial_bar_chart_templ.go`
- **Relative Path**: `ui/components/chart/radial_bar_chart_templ.go`
- **Purpose**: Compiled Go code generated by templ for radial_bar_chart.templ.

#### `/home/desktop/especa/beta/ui/components/chart/radial_bar_templ.go`
- **Relative Path**: `ui/components/chart/radial_bar_templ.go`
- **Purpose**: Compiled Go code generated by templ for radial_bar.templ.

#### `/home/desktop/especa/beta/ui/components/chart/responsive_container.templ`
- **Relative Path**: `ui/components/chart/responsive_container.templ`
- **Purpose**: Templ component rendering reusable UI element for responsive_container.templ in the chart category.

#### `/home/desktop/especa/beta/ui/components/chart/responsive_container_templ.go`
- **Relative Path**: `ui/components/chart/responsive_container_templ.go`
- **Purpose**: Compiled Go code generated by templ for responsive_container.templ.

#### `/home/desktop/especa/beta/ui/components/chart/scale.go`
- **Relative Path**: `ui/components/chart/scale.go`
- **Purpose**: Go types and view model definitions for the chart component library.

#### `/home/desktop/especa/beta/ui/components/chart/style.templ`
- **Relative Path**: `ui/components/chart/style.templ`
- **Purpose**: Templ component rendering reusable UI element for style.templ in the chart category.

#### `/home/desktop/especa/beta/ui/components/chart/style_templ.go`
- **Relative Path**: `ui/components/chart/style_templ.go`
- **Purpose**: Compiled Go code generated by templ for style.templ.

#### `/home/desktop/especa/beta/ui/components/chart/tooltip.templ`
- **Relative Path**: `ui/components/chart/tooltip.templ`
- **Purpose**: Templ component rendering reusable UI element for tooltip.templ in the chart category.

#### `/home/desktop/especa/beta/ui/components/chart/tooltip_templ.go`
- **Relative Path**: `ui/components/chart/tooltip_templ.go`
- **Purpose**: Compiled Go code generated by templ for tooltip.templ.

#### `/home/desktop/especa/beta/ui/components/chart/types.go`
- **Relative Path**: `ui/components/chart/types.go`
- **Purpose**: Go types and view model definitions for the chart component library.

#### `/home/desktop/especa/beta/ui/components/chart/x_axis.templ`
- **Relative Path**: `ui/components/chart/x_axis.templ`
- **Purpose**: Templ component rendering reusable UI element for x_axis.templ in the chart category.

#### `/home/desktop/especa/beta/ui/components/chart/x_axis_templ.go`
- **Relative Path**: `ui/components/chart/x_axis_templ.go`
- **Purpose**: Compiled Go code generated by templ for x_axis.templ.

#### `/home/desktop/especa/beta/ui/components/chart/y_axis.templ`
- **Relative Path**: `ui/components/chart/y_axis.templ`
- **Purpose**: Templ component rendering reusable UI element for y_axis.templ in the chart category.

#### `/home/desktop/especa/beta/ui/components/chart/y_axis_templ.go`
- **Relative Path**: `ui/components/chart/y_axis_templ.go`
- **Purpose**: Compiled Go code generated by templ for y_axis.templ.

#### `/home/desktop/especa/beta/ui/components/dialogs/dialog.templ`
- **Relative Path**: `ui/components/dialogs/dialog.templ`
- **Purpose**: Templ component rendering reusable UI element for dialog.templ in the dialogs category.

#### `/home/desktop/especa/beta/ui/components/dialogs/dialog_body.templ`
- **Relative Path**: `ui/components/dialogs/dialog_body.templ`
- **Purpose**: Templ component rendering reusable UI element for dialog_body.templ in the dialogs category.

#### `/home/desktop/especa/beta/ui/components/dialogs/dialog_body_templ.go`
- **Relative Path**: `ui/components/dialogs/dialog_body_templ.go`
- **Purpose**: Compiled Go code generated by templ for dialog_body.templ.

#### `/home/desktop/especa/beta/ui/components/dialogs/dialog_content.templ`
- **Relative Path**: `ui/components/dialogs/dialog_content.templ`
- **Purpose**: Templ component rendering reusable UI element for dialog_content.templ in the dialogs category.

#### `/home/desktop/especa/beta/ui/components/dialogs/dialog_content_templ.go`
- **Relative Path**: `ui/components/dialogs/dialog_content_templ.go`
- **Purpose**: Compiled Go code generated by templ for dialog_content.templ.

#### `/home/desktop/especa/beta/ui/components/dialogs/dialog_description.templ`
- **Relative Path**: `ui/components/dialogs/dialog_description.templ`
- **Purpose**: Templ component rendering reusable UI element for dialog_description.templ in the dialogs category.

#### `/home/desktop/especa/beta/ui/components/dialogs/dialog_description_templ.go`
- **Relative Path**: `ui/components/dialogs/dialog_description_templ.go`
- **Purpose**: Compiled Go code generated by templ for dialog_description.templ.

#### `/home/desktop/especa/beta/ui/components/dialogs/dialog_footer.templ`
- **Relative Path**: `ui/components/dialogs/dialog_footer.templ`
- **Purpose**: Templ component rendering reusable UI element for dialog_footer.templ in the dialogs category.

#### `/home/desktop/especa/beta/ui/components/dialogs/dialog_footer_templ.go`
- **Relative Path**: `ui/components/dialogs/dialog_footer_templ.go`
- **Purpose**: Compiled Go code generated by templ for dialog_footer.templ.

#### `/home/desktop/especa/beta/ui/components/dialogs/dialog_header.templ`
- **Relative Path**: `ui/components/dialogs/dialog_header.templ`
- **Purpose**: Templ component rendering reusable UI element for dialog_header.templ in the dialogs category.

#### `/home/desktop/especa/beta/ui/components/dialogs/dialog_header_templ.go`
- **Relative Path**: `ui/components/dialogs/dialog_header_templ.go`
- **Purpose**: Compiled Go code generated by templ for dialog_header.templ.

#### `/home/desktop/especa/beta/ui/components/dialogs/dialog_templ.go`
- **Relative Path**: `ui/components/dialogs/dialog_templ.go`
- **Purpose**: Compiled Go code generated by templ for dialog.templ.

#### `/home/desktop/especa/beta/ui/components/dialogs/dialog_title.templ`
- **Relative Path**: `ui/components/dialogs/dialog_title.templ`
- **Purpose**: Templ component rendering reusable UI element for dialog_title.templ in the dialogs category.

#### `/home/desktop/especa/beta/ui/components/dialogs/dialog_title_templ.go`
- **Relative Path**: `ui/components/dialogs/dialog_title_templ.go`
- **Purpose**: Compiled Go code generated by templ for dialog_title.templ.

#### `/home/desktop/especa/beta/ui/components/feedback/avatar.templ`
- **Relative Path**: `ui/components/feedback/avatar.templ`
- **Purpose**: Templ component rendering reusable UI element for avatar.templ in the feedback category.

#### `/home/desktop/especa/beta/ui/components/feedback/avatar_templ.go`
- **Relative Path**: `ui/components/feedback/avatar_templ.go`
- **Purpose**: Compiled Go code generated by templ for avatar.templ.

#### `/home/desktop/especa/beta/ui/components/feedback/badge.templ`
- **Relative Path**: `ui/components/feedback/badge.templ`
- **Purpose**: Templ component rendering reusable UI element for badge.templ in the feedback category.

#### `/home/desktop/especa/beta/ui/components/feedback/badge_templ.go`
- **Relative Path**: `ui/components/feedback/badge_templ.go`
- **Purpose**: Compiled Go code generated by templ for badge.templ.

#### `/home/desktop/especa/beta/ui/components/feedback/empty_state.templ`
- **Relative Path**: `ui/components/feedback/empty_state.templ`
- **Purpose**: Templ component rendering reusable UI element for empty_state.templ in the feedback category.

#### `/home/desktop/especa/beta/ui/components/feedback/empty_state_templ.go`
- **Relative Path**: `ui/components/feedback/empty_state_templ.go`
- **Purpose**: Compiled Go code generated by templ for empty_state.templ.

#### `/home/desktop/especa/beta/ui/components/feedback/marker.templ`
- **Relative Path**: `ui/components/feedback/marker.templ`
- **Purpose**: Templ component rendering reusable UI element for marker.templ in the feedback category.

#### `/home/desktop/especa/beta/ui/components/feedback/marker_content.templ`
- **Relative Path**: `ui/components/feedback/marker_content.templ`
- **Purpose**: Templ component rendering reusable UI element for marker_content.templ in the feedback category.

#### `/home/desktop/especa/beta/ui/components/feedback/marker_content_templ.go`
- **Relative Path**: `ui/components/feedback/marker_content_templ.go`
- **Purpose**: Compiled Go code generated by templ for marker_content.templ.

#### `/home/desktop/especa/beta/ui/components/feedback/marker_templ.go`
- **Relative Path**: `ui/components/feedback/marker_templ.go`
- **Purpose**: Compiled Go code generated by templ for marker.templ.

#### `/home/desktop/especa/beta/ui/components/feedback/marker_types.go`
- **Relative Path**: `ui/components/feedback/marker_types.go`
- **Purpose**: Go types and view model definitions for the feedback component library.

#### `/home/desktop/especa/beta/ui/components/feedback/page_header.templ`
- **Relative Path**: `ui/components/feedback/page_header.templ`
- **Purpose**: Templ component rendering reusable UI element for page_header.templ in the feedback category.

#### `/home/desktop/especa/beta/ui/components/feedback/page_header_templ.go`
- **Relative Path**: `ui/components/feedback/page_header_templ.go`
- **Purpose**: Compiled Go code generated by templ for page_header.templ.

#### `/home/desktop/especa/beta/ui/components/feedback/stat_card.templ`
- **Relative Path**: `ui/components/feedback/stat_card.templ`
- **Purpose**: Templ component rendering reusable UI element for stat_card.templ in the feedback category.

#### `/home/desktop/especa/beta/ui/components/feedback/stat_card_templ.go`
- **Relative Path**: `ui/components/feedback/stat_card_templ.go`
- **Purpose**: Compiled Go code generated by templ for stat_card.templ.

#### `/home/desktop/especa/beta/ui/components/feedback/status_pill.templ`
- **Relative Path**: `ui/components/feedback/status_pill.templ`
- **Purpose**: Templ component rendering reusable UI element for status_pill.templ in the feedback category.

#### `/home/desktop/especa/beta/ui/components/feedback/status_pill_templ.go`
- **Relative Path**: `ui/components/feedback/status_pill_templ.go`
- **Purpose**: Compiled Go code generated by templ for status_pill.templ.

#### `/home/desktop/especa/beta/ui/components/forms/input.templ`
- **Relative Path**: `ui/components/forms/input.templ`
- **Purpose**: Templ component rendering reusable UI element for input.templ in the forms category.

#### `/home/desktop/especa/beta/ui/components/forms/input_templ.go`
- **Relative Path**: `ui/components/forms/input_templ.go`
- **Purpose**: Compiled Go code generated by templ for input.templ.

#### `/home/desktop/especa/beta/ui/components/icons/icon_arrow_right.templ`
- **Relative Path**: `ui/components/icons/icon_arrow_right.templ`
- **Purpose**: Templ markup component rendering Lucide arrow_right SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_arrow_right_templ.go`
- **Relative Path**: `ui/components/icons/icon_arrow_right_templ.go`
- **Purpose**: Compiled Go template code for the Lucide arrow_right SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_bar_chart3.templ`
- **Relative Path**: `ui/components/icons/icon_bar_chart3.templ`
- **Purpose**: Templ markup component rendering Lucide bar_chart3 SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_bar_chart3_templ.go`
- **Relative Path**: `ui/components/icons/icon_bar_chart3_templ.go`
- **Purpose**: Compiled Go template code for the Lucide bar_chart3 SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_bell.templ`
- **Relative Path**: `ui/components/icons/icon_bell.templ`
- **Purpose**: Templ markup component rendering Lucide bell SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_bell_templ.go`
- **Relative Path**: `ui/components/icons/icon_bell_templ.go`
- **Purpose**: Compiled Go template code for the Lucide bell SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_brush.templ`
- **Relative Path**: `ui/components/icons/icon_brush.templ`
- **Purpose**: Templ markup component rendering Lucide brush SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_brush_templ.go`
- **Relative Path**: `ui/components/icons/icon_brush_templ.go`
- **Purpose**: Compiled Go template code for the Lucide brush SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_calendar.templ`
- **Relative Path**: `ui/components/icons/icon_calendar.templ`
- **Purpose**: Templ markup component rendering Lucide calendar SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_calendar_templ.go`
- **Relative Path**: `ui/components/icons/icon_calendar_templ.go`
- **Purpose**: Compiled Go template code for the Lucide calendar SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_camera.templ`
- **Relative Path**: `ui/components/icons/icon_camera.templ`
- **Purpose**: Templ markup component rendering Lucide camera SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_camera_templ.go`
- **Relative Path**: `ui/components/icons/icon_camera_templ.go`
- **Purpose**: Compiled Go template code for the Lucide camera SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_check.templ`
- **Relative Path**: `ui/components/icons/icon_check.templ`
- **Purpose**: Templ markup component rendering Lucide check SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_check_templ.go`
- **Relative Path**: `ui/components/icons/icon_check_templ.go`
- **Purpose**: Compiled Go template code for the Lucide check SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_chevron_down.templ`
- **Relative Path**: `ui/components/icons/icon_chevron_down.templ`
- **Purpose**: Templ markup component rendering Lucide chevron_down SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_chevron_down_templ.go`
- **Relative Path**: `ui/components/icons/icon_chevron_down_templ.go`
- **Purpose**: Compiled Go template code for the Lucide chevron_down SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_chevron_left.templ`
- **Relative Path**: `ui/components/icons/icon_chevron_left.templ`
- **Purpose**: Templ markup component rendering Lucide chevron_left SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_chevron_left_templ.go`
- **Relative Path**: `ui/components/icons/icon_chevron_left_templ.go`
- **Purpose**: Compiled Go template code for the Lucide chevron_left SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_chevron_right.templ`
- **Relative Path**: `ui/components/icons/icon_chevron_right.templ`
- **Purpose**: Templ markup component rendering Lucide chevron_right SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_chevron_right_templ.go`
- **Relative Path**: `ui/components/icons/icon_chevron_right_templ.go`
- **Purpose**: Compiled Go template code for the Lucide chevron_right SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_chevrons_up_down.templ`
- **Relative Path**: `ui/components/icons/icon_chevrons_up_down.templ`
- **Purpose**: Templ markup component rendering Lucide chevrons_up_down SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_chevrons_up_down_templ.go`
- **Relative Path**: `ui/components/icons/icon_chevrons_up_down_templ.go`
- **Purpose**: Compiled Go template code for the Lucide chevrons_up_down SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_circle_dot.templ`
- **Relative Path**: `ui/components/icons/icon_circle_dot.templ`
- **Purpose**: Templ markup component rendering Lucide circle_dot SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_circle_dot_templ.go`
- **Relative Path**: `ui/components/icons/icon_circle_dot_templ.go`
- **Purpose**: Compiled Go template code for the Lucide circle_dot SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_clock.templ`
- **Relative Path**: `ui/components/icons/icon_clock.templ`
- **Purpose**: Templ markup component rendering Lucide clock SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_clock_templ.go`
- **Relative Path**: `ui/components/icons/icon_clock_templ.go`
- **Purpose**: Compiled Go template code for the Lucide clock SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_compass.templ`
- **Relative Path**: `ui/components/icons/icon_compass.templ`
- **Purpose**: Templ markup component rendering Lucide compass SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_compass_templ.go`
- **Relative Path**: `ui/components/icons/icon_compass_templ.go`
- **Purpose**: Compiled Go template code for the Lucide compass SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_copy.templ`
- **Relative Path**: `ui/components/icons/icon_copy.templ`
- **Purpose**: Templ markup component rendering Lucide copy SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_copy_templ.go`
- **Relative Path**: `ui/components/icons/icon_copy_templ.go`
- **Purpose**: Compiled Go template code for the Lucide copy SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_credit_card.templ`
- **Relative Path**: `ui/components/icons/icon_credit_card.templ`
- **Purpose**: Templ markup component rendering Lucide credit_card SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_credit_card_templ.go`
- **Relative Path**: `ui/components/icons/icon_credit_card_templ.go`
- **Purpose**: Compiled Go template code for the Lucide credit_card SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_dashboard.templ`
- **Relative Path**: `ui/components/icons/icon_dashboard.templ`
- **Purpose**: Templ markup component rendering Lucide dashboard SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_dashboard_templ.go`
- **Relative Path**: `ui/components/icons/icon_dashboard_templ.go`
- **Purpose**: Compiled Go template code for the Lucide dashboard SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_dollar_sign.templ`
- **Relative Path**: `ui/components/icons/icon_dollar_sign.templ`
- **Purpose**: Templ markup component rendering Lucide dollar_sign SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_dollar_sign_templ.go`
- **Relative Path**: `ui/components/icons/icon_dollar_sign_templ.go`
- **Purpose**: Compiled Go template code for the Lucide dollar_sign SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_file_text.templ`
- **Relative Path**: `ui/components/icons/icon_file_text.templ`
- **Purpose**: Templ markup component rendering Lucide file_text SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_file_text_templ.go`
- **Relative Path**: `ui/components/icons/icon_file_text_templ.go`
- **Purpose**: Compiled Go template code for the Lucide file_text SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_filter.templ`
- **Relative Path**: `ui/components/icons/icon_filter.templ`
- **Purpose**: Templ markup component rendering Lucide filter SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_filter_templ.go`
- **Relative Path**: `ui/components/icons/icon_filter_templ.go`
- **Purpose**: Compiled Go template code for the Lucide filter SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_folder.templ`
- **Relative Path**: `ui/components/icons/icon_folder.templ`
- **Purpose**: Templ markup component rendering Lucide folder SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_folder_templ.go`
- **Relative Path**: `ui/components/icons/icon_folder_templ.go`
- **Purpose**: Compiled Go template code for the Lucide folder SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_kanban.templ`
- **Relative Path**: `ui/components/icons/icon_kanban.templ`
- **Purpose**: Templ markup component rendering Lucide kanban SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_kanban_templ.go`
- **Relative Path**: `ui/components/icons/icon_kanban_templ.go`
- **Purpose**: Compiled Go template code for the Lucide kanban SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_layout.templ`
- **Relative Path**: `ui/components/icons/icon_layout.templ`
- **Purpose**: Templ markup component rendering Lucide layout SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_layout_templ.go`
- **Relative Path**: `ui/components/icons/icon_layout_templ.go`
- **Purpose**: Compiled Go template code for the Lucide layout SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_list.templ`
- **Relative Path**: `ui/components/icons/icon_list.templ`
- **Purpose**: Templ markup component rendering Lucide list SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_list_templ.go`
- **Relative Path**: `ui/components/icons/icon_list_templ.go`
- **Purpose**: Compiled Go template code for the Lucide list SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_map_pin.templ`
- **Relative Path**: `ui/components/icons/icon_map_pin.templ`
- **Purpose**: Templ markup component rendering Lucide map_pin SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_map_pin_templ.go`
- **Relative Path**: `ui/components/icons/icon_map_pin_templ.go`
- **Purpose**: Compiled Go template code for the Lucide map_pin SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_megaphone.templ`
- **Relative Path**: `ui/components/icons/icon_megaphone.templ`
- **Purpose**: Templ markup component rendering Lucide megaphone SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_megaphone_templ.go`
- **Relative Path**: `ui/components/icons/icon_megaphone_templ.go`
- **Purpose**: Compiled Go template code for the Lucide megaphone SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_message_square.templ`
- **Relative Path**: `ui/components/icons/icon_message_square.templ`
- **Purpose**: Templ markup component rendering Lucide message_square SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_message_square_templ.go`
- **Relative Path**: `ui/components/icons/icon_message_square_templ.go`
- **Purpose**: Compiled Go template code for the Lucide message_square SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_more_horizontal.templ`
- **Relative Path**: `ui/components/icons/icon_more_horizontal.templ`
- **Purpose**: Templ markup component rendering Lucide more_horizontal SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_more_horizontal_templ.go`
- **Relative Path**: `ui/components/icons/icon_more_horizontal_templ.go`
- **Purpose**: Compiled Go template code for the Lucide more_horizontal SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_package.templ`
- **Relative Path**: `ui/components/icons/icon_package.templ`
- **Purpose**: Templ markup component rendering Lucide package SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_package_templ.go`
- **Relative Path**: `ui/components/icons/icon_package_templ.go`
- **Purpose**: Compiled Go template code for the Lucide package SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_plus.templ`
- **Relative Path**: `ui/components/icons/icon_plus.templ`
- **Purpose**: Templ markup component rendering Lucide plus SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_plus_templ.go`
- **Relative Path**: `ui/components/icons/icon_plus_templ.go`
- **Purpose**: Compiled Go template code for the Lucide plus SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_receipt.templ`
- **Relative Path**: `ui/components/icons/icon_receipt.templ`
- **Purpose**: Templ markup component rendering Lucide receipt SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_receipt_templ.go`
- **Relative Path**: `ui/components/icons/icon_receipt_templ.go`
- **Purpose**: Compiled Go template code for the Lucide receipt SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_search.templ`
- **Relative Path**: `ui/components/icons/icon_search.templ`
- **Purpose**: Templ markup component rendering Lucide search SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_search_templ.go`
- **Relative Path**: `ui/components/icons/icon_search_templ.go`
- **Purpose**: Compiled Go template code for the Lucide search SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_settings.templ`
- **Relative Path**: `ui/components/icons/icon_settings.templ`
- **Purpose**: Templ markup component rendering Lucide settings SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_settings_templ.go`
- **Relative Path**: `ui/components/icons/icon_settings_templ.go`
- **Purpose**: Compiled Go template code for the Lucide settings SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_sidebar.templ`
- **Relative Path**: `ui/components/icons/icon_sidebar.templ`
- **Purpose**: Templ markup component rendering Lucide sidebar SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_sidebar_templ.go`
- **Relative Path**: `ui/components/icons/icon_sidebar_templ.go`
- **Purpose**: Compiled Go template code for the Lucide sidebar SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_sliders_horizontal.templ`
- **Relative Path**: `ui/components/icons/icon_sliders_horizontal.templ`
- **Purpose**: Templ markup component rendering Lucide sliders_horizontal SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_sliders_horizontal_templ.go`
- **Relative Path**: `ui/components/icons/icon_sliders_horizontal_templ.go`
- **Purpose**: Compiled Go template code for the Lucide sliders_horizontal SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_sparkles.templ`
- **Relative Path**: `ui/components/icons/icon_sparkles.templ`
- **Purpose**: Templ markup component rendering Lucide sparkles SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_sparkles_templ.go`
- **Relative Path**: `ui/components/icons/icon_sparkles_templ.go`
- **Purpose**: Compiled Go template code for the Lucide sparkles SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_star.templ`
- **Relative Path**: `ui/components/icons/icon_star.templ`
- **Purpose**: Templ markup component rendering Lucide star SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_star_templ.go`
- **Relative Path**: `ui/components/icons/icon_star_templ.go`
- **Purpose**: Compiled Go template code for the Lucide star SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_sun.templ`
- **Relative Path**: `ui/components/icons/icon_sun.templ`
- **Purpose**: Templ markup component rendering Lucide sun SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_sun_templ.go`
- **Relative Path**: `ui/components/icons/icon_sun_templ.go`
- **Purpose**: Compiled Go template code for the Lucide sun SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_table.templ`
- **Relative Path**: `ui/components/icons/icon_table.templ`
- **Purpose**: Templ markup component rendering Lucide table SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_table_templ.go`
- **Relative Path**: `ui/components/icons/icon_table_templ.go`
- **Purpose**: Compiled Go template code for the Lucide table SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_tag.templ`
- **Relative Path**: `ui/components/icons/icon_tag.templ`
- **Purpose**: Templ markup component rendering Lucide tag SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_tag_templ.go`
- **Relative Path**: `ui/components/icons/icon_tag_templ.go`
- **Purpose**: Compiled Go template code for the Lucide tag SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_target.templ`
- **Relative Path**: `ui/components/icons/icon_target.templ`
- **Purpose**: Templ markup component rendering Lucide target SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_target_templ.go`
- **Relative Path**: `ui/components/icons/icon_target_templ.go`
- **Purpose**: Compiled Go template code for the Lucide target SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_trash2.templ`
- **Relative Path**: `ui/components/icons/icon_trash2.templ`
- **Purpose**: Templ markup component rendering Lucide trash2 SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_trash2_templ.go`
- **Relative Path**: `ui/components/icons/icon_trash2_templ.go`
- **Purpose**: Compiled Go template code for the Lucide trash2 SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_trending_up.templ`
- **Relative Path**: `ui/components/icons/icon_trending_up.templ`
- **Purpose**: Templ markup component rendering Lucide trending_up SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_trending_up_templ.go`
- **Relative Path**: `ui/components/icons/icon_trending_up_templ.go`
- **Purpose**: Compiled Go template code for the Lucide trending_up SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_user.templ`
- **Relative Path**: `ui/components/icons/icon_user.templ`
- **Purpose**: Templ markup component rendering Lucide user SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_user_plus.templ`
- **Relative Path**: `ui/components/icons/icon_user_plus.templ`
- **Purpose**: Templ markup component rendering Lucide user_plus SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_user_plus_templ.go`
- **Relative Path**: `ui/components/icons/icon_user_plus_templ.go`
- **Purpose**: Compiled Go template code for the Lucide user_plus SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_user_templ.go`
- **Relative Path**: `ui/components/icons/icon_user_templ.go`
- **Purpose**: Compiled Go template code for the Lucide user SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_users.templ`
- **Relative Path**: `ui/components/icons/icon_users.templ`
- **Purpose**: Templ markup component rendering Lucide users SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_users_templ.go`
- **Relative Path**: `ui/components/icons/icon_users_templ.go`
- **Purpose**: Compiled Go template code for the Lucide users SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_x.templ`
- **Relative Path**: `ui/components/icons/icon_x.templ`
- **Purpose**: Templ markup component rendering Lucide x SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_x_templ.go`
- **Relative Path**: `ui/components/icons/icon_x_templ.go`
- **Purpose**: Compiled Go template code for the Lucide x SVG icon component.

#### `/home/desktop/especa/beta/ui/components/icons/icon_zap.templ`
- **Relative Path**: `ui/components/icons/icon_zap.templ`
- **Purpose**: Templ markup component rendering Lucide zap SVG icon with customizable size and CSS classes.

#### `/home/desktop/especa/beta/ui/components/icons/icon_zap_templ.go`
- **Relative Path**: `ui/components/icons/icon_zap_templ.go`
- **Purpose**: Compiled Go template code for the Lucide zap SVG icon component.

#### `/home/desktop/especa/beta/ui/components/messages/bubble.templ`
- **Relative Path**: `ui/components/messages/bubble.templ`
- **Purpose**: Templ component rendering reusable UI element for bubble.templ in the messages category.

#### `/home/desktop/especa/beta/ui/components/messages/bubble_content.templ`
- **Relative Path**: `ui/components/messages/bubble_content.templ`
- **Purpose**: Templ component rendering reusable UI element for bubble_content.templ in the messages category.

#### `/home/desktop/especa/beta/ui/components/messages/bubble_content_templ.go`
- **Relative Path**: `ui/components/messages/bubble_content_templ.go`
- **Purpose**: Compiled Go code generated by templ for bubble_content.templ.

#### `/home/desktop/especa/beta/ui/components/messages/bubble_group.templ`
- **Relative Path**: `ui/components/messages/bubble_group.templ`
- **Purpose**: Templ component rendering reusable UI element for bubble_group.templ in the messages category.

#### `/home/desktop/especa/beta/ui/components/messages/bubble_group_templ.go`
- **Relative Path**: `ui/components/messages/bubble_group_templ.go`
- **Purpose**: Compiled Go code generated by templ for bubble_group.templ.

#### `/home/desktop/especa/beta/ui/components/messages/bubble_reactions.templ`
- **Relative Path**: `ui/components/messages/bubble_reactions.templ`
- **Purpose**: Templ component rendering reusable UI element for bubble_reactions.templ in the messages category.

#### `/home/desktop/especa/beta/ui/components/messages/bubble_reactions_templ.go`
- **Relative Path**: `ui/components/messages/bubble_reactions_templ.go`
- **Purpose**: Compiled Go code generated by templ for bubble_reactions.templ.

#### `/home/desktop/especa/beta/ui/components/messages/bubble_templ.go`
- **Relative Path**: `ui/components/messages/bubble_templ.go`
- **Purpose**: Compiled Go code generated by templ for bubble.templ.

#### `/home/desktop/especa/beta/ui/components/messages/bubble_types.go`
- **Relative Path**: `ui/components/messages/bubble_types.go`
- **Purpose**: Go types and view model definitions for the messages component library.

#### `/home/desktop/especa/beta/ui/components/messages/message.templ`
- **Relative Path**: `ui/components/messages/message.templ`
- **Purpose**: Templ component rendering reusable UI element for message.templ in the messages category.

#### `/home/desktop/especa/beta/ui/components/messages/message_avatar_slot.templ`
- **Relative Path**: `ui/components/messages/message_avatar_slot.templ`
- **Purpose**: Templ component rendering reusable UI element for message_avatar_slot.templ in the messages category.

#### `/home/desktop/especa/beta/ui/components/messages/message_avatar_slot_templ.go`
- **Relative Path**: `ui/components/messages/message_avatar_slot_templ.go`
- **Purpose**: Compiled Go code generated by templ for message_avatar_slot.templ.

#### `/home/desktop/especa/beta/ui/components/messages/message_content.templ`
- **Relative Path**: `ui/components/messages/message_content.templ`
- **Purpose**: Templ component rendering reusable UI element for message_content.templ in the messages category.

#### `/home/desktop/especa/beta/ui/components/messages/message_content_templ.go`
- **Relative Path**: `ui/components/messages/message_content_templ.go`
- **Purpose**: Compiled Go code generated by templ for message_content.templ.

#### `/home/desktop/especa/beta/ui/components/messages/message_footer.templ`
- **Relative Path**: `ui/components/messages/message_footer.templ`
- **Purpose**: Templ component rendering reusable UI element for message_footer.templ in the messages category.

#### `/home/desktop/especa/beta/ui/components/messages/message_footer_templ.go`
- **Relative Path**: `ui/components/messages/message_footer_templ.go`
- **Purpose**: Compiled Go code generated by templ for message_footer.templ.

#### `/home/desktop/especa/beta/ui/components/messages/message_group.templ`
- **Relative Path**: `ui/components/messages/message_group.templ`
- **Purpose**: Templ component rendering reusable UI element for message_group.templ in the messages category.

#### `/home/desktop/especa/beta/ui/components/messages/message_group_templ.go`
- **Relative Path**: `ui/components/messages/message_group_templ.go`
- **Purpose**: Compiled Go code generated by templ for message_group.templ.

#### `/home/desktop/especa/beta/ui/components/messages/message_header.templ`
- **Relative Path**: `ui/components/messages/message_header.templ`
- **Purpose**: Templ component rendering reusable UI element for message_header.templ in the messages category.

#### `/home/desktop/especa/beta/ui/components/messages/message_header_templ.go`
- **Relative Path**: `ui/components/messages/message_header_templ.go`
- **Purpose**: Compiled Go code generated by templ for message_header.templ.

#### `/home/desktop/especa/beta/ui/components/messages/message_templ.go`
- **Relative Path**: `ui/components/messages/message_templ.go`
- **Purpose**: Compiled Go code generated by templ for message.templ.

#### `/home/desktop/especa/beta/ui/components/messages/message_types.go`
- **Relative Path**: `ui/components/messages/message_types.go`
- **Purpose**: Go types and view model definitions for the messages component library.

#### `/home/desktop/especa/beta/ui/components/navigation/pagination.templ`
- **Relative Path**: `ui/components/navigation/pagination.templ`
- **Purpose**: Templ component rendering reusable UI element for pagination.templ in the navigation category.

#### `/home/desktop/especa/beta/ui/components/navigation/pagination_templ.go`
- **Relative Path**: `ui/components/navigation/pagination_templ.go`
- **Purpose**: Compiled Go code generated by templ for pagination.templ.

#### `/home/desktop/especa/beta/ui/components/navigation/sidebar.templ`
- **Relative Path**: `ui/components/navigation/sidebar.templ`
- **Purpose**: Templ component rendering reusable UI element for sidebar.templ in the navigation category.

#### `/home/desktop/especa/beta/ui/components/navigation/sidebar_templ.go`
- **Relative Path**: `ui/components/navigation/sidebar_templ.go`
- **Purpose**: Compiled Go code generated by templ for sidebar.templ.

#### `/home/desktop/especa/beta/ui/components/navigation/sidebar_types.go`
- **Relative Path**: `ui/components/navigation/sidebar_types.go`
- **Purpose**: Go types and view model definitions for the navigation component library.

#### `/home/desktop/especa/beta/ui/components/navigation/spotlight_search.templ`
- **Relative Path**: `ui/components/navigation/spotlight_search.templ`
- **Purpose**: Templ component rendering reusable UI element for spotlight_search.templ in the navigation category.

#### `/home/desktop/especa/beta/ui/components/navigation/spotlight_search_templ.go`
- **Relative Path**: `ui/components/navigation/spotlight_search_templ.go`
- **Purpose**: Compiled Go code generated by templ for spotlight_search.templ.

#### `/home/desktop/especa/beta/ui/components/navigation/tab_trigger.templ`
- **Relative Path**: `ui/components/navigation/tab_trigger.templ`
- **Purpose**: Templ component rendering reusable UI element for tab_trigger.templ in the navigation category.

#### `/home/desktop/especa/beta/ui/components/navigation/tab_trigger_templ.go`
- **Relative Path**: `ui/components/navigation/tab_trigger_templ.go`
- **Purpose**: Compiled Go code generated by templ for tab_trigger.templ.

#### `/home/desktop/especa/beta/ui/components/navigation/tabs.templ`
- **Relative Path**: `ui/components/navigation/tabs.templ`
- **Purpose**: Templ component rendering reusable UI element for tabs.templ in the navigation category.

#### `/home/desktop/especa/beta/ui/components/navigation/tabs_templ.go`
- **Relative Path**: `ui/components/navigation/tabs_templ.go`
- **Purpose**: Compiled Go code generated by templ for tabs.templ.

#### `/home/desktop/especa/beta/ui/components/tables/table.templ`
- **Relative Path**: `ui/components/tables/table.templ`
- **Purpose**: Templ component rendering reusable UI element for table.templ in the tables category.

#### `/home/desktop/especa/beta/ui/components/tables/table_body.templ`
- **Relative Path**: `ui/components/tables/table_body.templ`
- **Purpose**: Templ component rendering reusable UI element for table_body.templ in the tables category.

#### `/home/desktop/especa/beta/ui/components/tables/table_body_templ.go`
- **Relative Path**: `ui/components/tables/table_body_templ.go`
- **Purpose**: Compiled Go code generated by templ for table_body.templ.

#### `/home/desktop/especa/beta/ui/components/tables/table_cell.templ`
- **Relative Path**: `ui/components/tables/table_cell.templ`
- **Purpose**: Templ component rendering reusable UI element for table_cell.templ in the tables category.

#### `/home/desktop/especa/beta/ui/components/tables/table_cell_templ.go`
- **Relative Path**: `ui/components/tables/table_cell_templ.go`
- **Purpose**: Compiled Go code generated by templ for table_cell.templ.

#### `/home/desktop/especa/beta/ui/components/tables/table_footer.templ`
- **Relative Path**: `ui/components/tables/table_footer.templ`
- **Purpose**: Templ component rendering reusable UI element for table_footer.templ in the tables category.

#### `/home/desktop/especa/beta/ui/components/tables/table_footer_templ.go`
- **Relative Path**: `ui/components/tables/table_footer_templ.go`
- **Purpose**: Compiled Go code generated by templ for table_footer.templ.

#### `/home/desktop/especa/beta/ui/components/tables/table_head.templ`
- **Relative Path**: `ui/components/tables/table_head.templ`
- **Purpose**: Templ component rendering reusable UI element for table_head.templ in the tables category.

#### `/home/desktop/especa/beta/ui/components/tables/table_head_templ.go`
- **Relative Path**: `ui/components/tables/table_head_templ.go`
- **Purpose**: Compiled Go code generated by templ for table_head.templ.

#### `/home/desktop/especa/beta/ui/components/tables/table_header_cell.templ`
- **Relative Path**: `ui/components/tables/table_header_cell.templ`
- **Purpose**: Templ component rendering reusable UI element for table_header_cell.templ in the tables category.

#### `/home/desktop/especa/beta/ui/components/tables/table_header_cell_templ.go`
- **Relative Path**: `ui/components/tables/table_header_cell_templ.go`
- **Purpose**: Compiled Go code generated by templ for table_header_cell.templ.

#### `/home/desktop/especa/beta/ui/components/tables/table_row.templ`
- **Relative Path**: `ui/components/tables/table_row.templ`
- **Purpose**: Templ component rendering reusable UI element for table_row.templ in the tables category.

#### `/home/desktop/especa/beta/ui/components/tables/table_row_templ.go`
- **Relative Path**: `ui/components/tables/table_row_templ.go`
- **Purpose**: Compiled Go code generated by templ for table_row.templ.

#### `/home/desktop/especa/beta/ui/components/tables/table_templ.go`
- **Relative Path**: `ui/components/tables/table_templ.go`
- **Purpose**: Compiled Go code generated by templ for table.templ.

#### `/home/desktop/especa/beta/ui/components/tables/table_types.go`
- **Relative Path**: `ui/components/tables/table_types.go`
- **Purpose**: Go types and view model definitions for the tables component library.

#### `/home/desktop/especa/beta/ui/go.mod`
- **Relative Path**: `ui/go.mod`
- **Purpose**: Go module manifest declaring dependencies and module path for ui.

#### `/home/desktop/especa/beta/ui/go.sum`
- **Relative Path**: `ui/go.sum`
- **Purpose**: Checksum lockfile ensuring reproducible builds and dependency integrity for ui.

#### `/home/desktop/especa/beta/ui/layouts/auth_layout.templ`
- **Relative Path**: `ui/layouts/auth_layout.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for auth_layout.templ.

#### `/home/desktop/especa/beta/ui/layouts/auth_layout_templ.go`
- **Relative Path**: `ui/layouts/auth_layout_templ.go`
- **Purpose**: Compiled Go template code generated by templ for auth_layout.templ.

#### `/home/desktop/especa/beta/ui/layouts/dashboard_layout.templ`
- **Relative Path**: `ui/layouts/dashboard_layout.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for dashboard_layout.templ.

#### `/home/desktop/especa/beta/ui/layouts/dashboard_layout_templ.go`
- **Relative Path**: `ui/layouts/dashboard_layout_templ.go`
- **Purpose**: Compiled Go template code generated by templ for dashboard_layout.templ.

#### `/home/desktop/especa/beta/ui/layouts/kanban_card.templ`
- **Relative Path**: `ui/layouts/kanban_card.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for kanban_card.templ.

#### `/home/desktop/especa/beta/ui/layouts/kanban_card_templ.go`
- **Relative Path**: `ui/layouts/kanban_card_templ.go`
- **Purpose**: Compiled Go template code generated by templ for kanban_card.templ.

#### `/home/desktop/especa/beta/ui/layouts/kanban_column.templ`
- **Relative Path**: `ui/layouts/kanban_column.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for kanban_column.templ.

#### `/home/desktop/especa/beta/ui/layouts/kanban_column_templ.go`
- **Relative Path**: `ui/layouts/kanban_column_templ.go`
- **Purpose**: Compiled Go template code generated by templ for kanban_column.templ.

#### `/home/desktop/especa/beta/ui/layouts/kanban_layout.templ`
- **Relative Path**: `ui/layouts/kanban_layout.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for kanban_layout.templ.

#### `/home/desktop/especa/beta/ui/layouts/kanban_layout_templ.go`
- **Relative Path**: `ui/layouts/kanban_layout_templ.go`
- **Purpose**: Compiled Go template code generated by templ for kanban_layout.templ.

#### `/home/desktop/especa/beta/ui/layouts/settings_layout.templ`
- **Relative Path**: `ui/layouts/settings_layout.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for settings_layout.templ.

#### `/home/desktop/especa/beta/ui/layouts/settings_layout_templ.go`
- **Relative Path**: `ui/layouts/settings_layout_templ.go`
- **Purpose**: Compiled Go template code generated by templ for settings_layout.templ.

#### `/home/desktop/especa/beta/ui/utils/utils.go`
- **Relative Path**: `ui/utils/utils.go`
- **Purpose**: Utility helpers for templ markup generation, class toggling, and rendering abstractions.


### 5. Core Architecture & Platform Services (`projects/core`) (27 files)

#### `/home/desktop/especa/beta/projects/core/app.go`
- **Relative Path**: `projects/core/app.go`
- **Purpose**: Core application definitions and runtime primitives.

#### `/home/desktop/especa/beta/projects/core/db/migrator.go`
- **Relative Path**: `projects/core/db/migrator.go`
- **Purpose**: Database migration runner applying versioned SQL migration files to active database pools.

#### `/home/desktop/especa/beta/projects/core/db/queries/queries.sql`
- **Relative Path**: `projects/core/db/queries/queries.sql`
- **Purpose**: SQL query definitions consumed by sqlc to generate type-safe Go database functions.

#### `/home/desktop/especa/beta/projects/core/db/schema/schema.go`
- **Relative Path**: `projects/core/db/schema/schema.go`
- **Purpose**: Schema migration executor and table initialization utilities for core services.

#### `/home/desktop/especa/beta/projects/core/db/schema/schema.sql`
- **Relative Path**: `projects/core/db/schema/schema.sql`
- **Purpose**: Master DDL schema declaring tables, constraints, foreign keys, and indexes for platform users and roles.

#### `/home/desktop/especa/beta/projects/core/db/sqlc/db.go`
- **Relative Path**: `projects/core/db/sqlc/db.go`
- **Purpose**: Generated DBTX interface and database connection handler created by sqlc.

#### `/home/desktop/especa/beta/projects/core/db/sqlc/models.go`
- **Relative Path**: `projects/core/db/sqlc/models.go`
- **Purpose**: Generated Go structs representing core table rows with strict typing.

#### `/home/desktop/especa/beta/projects/core/db/sqlc/queries.sql.go`
- **Relative Path**: `projects/core/db/sqlc/queries.sql.go`
- **Purpose**: Generated type-safe database query functions created by sqlc.

#### `/home/desktop/especa/beta/projects/core/go.mod`
- **Relative Path**: `projects/core/go.mod`
- **Purpose**: Go module manifest declaring dependencies and module path for projects/core.

#### `/home/desktop/especa/beta/projects/core/go.sum`
- **Relative Path**: `projects/core/go.sum`
- **Purpose**: Checksum lockfile ensuring reproducible builds and dependency integrity for projects/core.

#### `/home/desktop/especa/beta/projects/core/handlers/handlers.go`
- **Relative Path**: `projects/core/handlers/handlers.go`
- **Purpose**: Shared HTTP controller handlers implementing core REST API endpoints and error handling.

#### `/home/desktop/especa/beta/projects/core/middleware/csrf.go`
- **Relative Path**: `projects/core/middleware/csrf.go`
- **Purpose**: Cross-Site Request Forgery (CSRF) protection middleware and token validation.

#### `/home/desktop/especa/beta/projects/core/middleware/middleware.go`
- **Relative Path**: `projects/core/middleware/middleware.go`
- **Purpose**: HTTP middleware providing session verification, cookie extraction, and route guard authentication.

#### `/home/desktop/especa/beta/projects/core/middleware/tenant_middleware.go`
- **Relative Path**: `projects/core/middleware/tenant_middleware.go`
- **Purpose**: Multi-tenant HTTP middleware enforcing tenant isolation and injecting tenant context into requests.

#### `/home/desktop/especa/beta/projects/core/models/models.go`
- **Relative Path**: `projects/core/models/models.go`
- **Purpose**: Core data models, user and session entities, and request/response DTOs.

#### `/home/desktop/especa/beta/projects/core/repositories/registration_repository.go`
- **Relative Path**: `projects/core/repositories/registration_repository.go`
- **Purpose**: Repository handling user and company registration transactions across platform databases.

#### `/home/desktop/especa/beta/projects/core/repositories/repositories.go`
- **Relative Path**: `projects/core/repositories/repositories.go`
- **Purpose**: Repository layer encapsulating database access logic behind clean interface abstractions.

#### `/home/desktop/especa/beta/projects/core/security/password.go`
- **Relative Path**: `projects/core/security/password.go`
- **Purpose**: Provides Argon2id and bcrypt cryptographic password hashing, salting, and constant-time verification.

#### `/home/desktop/especa/beta/projects/core/security/session.go`
- **Relative Path**: `projects/core/security/session.go`
- **Purpose**: Provides HMAC-SHA256 authenticated session cookie generation, encryption, and validation.

#### `/home/desktop/especa/beta/projects/core/services/auth_service.go`
- **Relative Path**: `projects/core/services/auth_service.go`
- **Purpose**: Authentication service containing business logic for logins, signups, credential checks, and role enforcement.

#### `/home/desktop/especa/beta/projects/core/sqlc.yaml`
- **Relative Path**: `projects/core/sqlc.yaml`
- **Purpose**: Sqlc configuration file configuring code generation for type-safe database queries.

#### `/home/desktop/especa/beta/projects/core/utils/config.go`
- **Relative Path**: `projects/core/utils/config.go`
- **Purpose**: Loads environment variables into structured Config structs with fallback defaults.

#### `/home/desktop/especa/beta/projects/core/utils/db.go`
- **Relative Path**: `projects/core/utils/db.go`
- **Purpose**: Database connection pool manager setting up MySQL/SQLite connection pools (DBPool) with WAL mode and health checks.

#### `/home/desktop/especa/beta/projects/core/utils/error_logger.go`
- **Relative Path**: `projects/core/utils/error_logger.go`
- **Purpose**: Structured error logging and diagnostics utility for API handlers.

#### `/home/desktop/especa/beta/projects/core/utils/money.go`
- **Relative Path**: `projects/core/utils/money.go`
- **Purpose**: Type-safe precision monetary arithmetic and currency formatting utilities.

#### `/home/desktop/especa/beta/projects/core/utils/response.go`
- **Relative Path**: `projects/core/utils/response.go`
- **Purpose**: Standardized HTTP JSON response utilities providing uniform success and error payloads.

#### `/home/desktop/especa/beta/projects/core/utils/user_context.go`
- **Relative Path**: `projects/core/utils/user_context.go`
- **Purpose**: Context utilities for injecting and retrieving authenticated user claims and tenant IDs from HTTP request context.


### 6. Superadmin Service (`projects/superadmin`) (42 files)

#### `/home/desktop/especa/beta/projects/superadmin/go.mod`
- **Relative Path**: `projects/superadmin/go.mod`
- **Purpose**: Go module manifest declaring dependencies and module path for projects/superadmin.

#### `/home/desktop/especa/beta/projects/superadmin/go.sum`
- **Relative Path**: `projects/superadmin/go.sum`
- **Purpose**: Checksum lockfile ensuring reproducible builds and dependency integrity for projects/superadmin.

#### `/home/desktop/especa/beta/projects/superadmin/handlers/auth_handler.go`
- **Relative Path**: `projects/superadmin/handlers/auth_handler.go`
- **Purpose**: HTTP controller handler implementing request decoding, validation, service coordination, and responses.

#### `/home/desktop/especa/beta/projects/superadmin/handlers/home_handler.go`
- **Relative Path**: `projects/superadmin/handlers/home_handler.go`
- **Purpose**: HTTP controller handler implementing request decoding, validation, service coordination, and responses.

#### `/home/desktop/especa/beta/projects/superadmin/handlers/tenant_cms_features_handler.go`
- **Relative Path**: `projects/superadmin/handlers/tenant_cms_features_handler.go`
- **Purpose**: HTTP controller handler implementing request decoding, validation, service coordination, and responses.

#### `/home/desktop/especa/beta/projects/superadmin/main.go`
- **Relative Path**: `projects/superadmin/main.go`
- **Purpose**: Main entrypoint initializing database connections, service dependencies, HTTP routers, and listeners for projects/superadmin.

#### `/home/desktop/especa/beta/projects/superadmin/models/models.go`
- **Relative Path**: `projects/superadmin/models/models.go`
- **Purpose**: Domain models, database entities, and Data Transfer Objects (DTOs).

#### `/home/desktop/especa/beta/projects/superadmin/repositories/platform_admin_repository.go`
- **Relative Path**: `projects/superadmin/repositories/platform_admin_repository.go`
- **Purpose**: Data access repository implementing database queries and CRUD operations for domain models.

#### `/home/desktop/especa/beta/projects/superadmin/routes/backend.go`
- **Relative Path**: `projects/superadmin/routes/backend.go`
- **Purpose**: HTTP router configuration and endpoint routing definitions for projects/superadmin/routes.

#### `/home/desktop/especa/beta/projects/superadmin/routes/frontend.go`
- **Relative Path**: `projects/superadmin/routes/frontend.go`
- **Purpose**: HTTP router configuration and endpoint routing definitions for projects/superadmin/routes.

#### `/home/desktop/especa/beta/projects/superadmin/routes/router.go`
- **Relative Path**: `projects/superadmin/routes/router.go`
- **Purpose**: HTTP router configuration and endpoint routing definitions for projects/superadmin/routes.

#### `/home/desktop/especa/beta/projects/superadmin/services/platform_admin_service.go`
- **Relative Path**: `projects/superadmin/services/platform_admin_service.go`
- **Purpose**: Business logic service layer handling core workflows, validation rules, and operations.

#### `/home/desktop/especa/beta/projects/superadmin/ui/pages/dashboard.templ`
- **Relative Path**: `projects/superadmin/ui/pages/dashboard.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for dashboard.templ.

#### `/home/desktop/especa/beta/projects/superadmin/ui/pages/dashboard_templ.go`
- **Relative Path**: `projects/superadmin/ui/pages/dashboard_templ.go`
- **Purpose**: Compiled Go template code generated by templ for dashboard.templ.

#### `/home/desktop/especa/beta/projects/superadmin/ui/pages/forgot_password.templ`
- **Relative Path**: `projects/superadmin/ui/pages/forgot_password.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for forgot_password.templ.

#### `/home/desktop/especa/beta/projects/superadmin/ui/pages/forgot_password_templ.go`
- **Relative Path**: `projects/superadmin/ui/pages/forgot_password_templ.go`
- **Purpose**: Compiled Go template code generated by templ for forgot_password.templ.

#### `/home/desktop/especa/beta/projects/superadmin/ui/pages/internal_error.templ`
- **Relative Path**: `projects/superadmin/ui/pages/internal_error.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for internal_error.templ.

#### `/home/desktop/especa/beta/projects/superadmin/ui/pages/internal_error_templ.go`
- **Relative Path**: `projects/superadmin/ui/pages/internal_error_templ.go`
- **Purpose**: Compiled Go template code generated by templ for internal_error.templ.

#### `/home/desktop/especa/beta/projects/superadmin/ui/pages/not_found.templ`
- **Relative Path**: `projects/superadmin/ui/pages/not_found.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for not_found.templ.

#### `/home/desktop/especa/beta/projects/superadmin/ui/pages/not_found_templ.go`
- **Relative Path**: `projects/superadmin/ui/pages/not_found_templ.go`
- **Purpose**: Compiled Go template code generated by templ for not_found.templ.

#### `/home/desktop/especa/beta/projects/superadmin/ui/pages/signin.templ`
- **Relative Path**: `projects/superadmin/ui/pages/signin.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for signin.templ.

#### `/home/desktop/especa/beta/projects/superadmin/ui/pages/signin_templ.go`
- **Relative Path**: `projects/superadmin/ui/pages/signin_templ.go`
- **Purpose**: Compiled Go template code generated by templ for signin.templ.

#### `/home/desktop/especa/beta/projects/superadmin/ui/pages/signup.templ`
- **Relative Path**: `projects/superadmin/ui/pages/signup.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for signup.templ.

#### `/home/desktop/especa/beta/projects/superadmin/ui/pages/signup_templ.go`
- **Relative Path**: `projects/superadmin/ui/pages/signup_templ.go`
- **Purpose**: Compiled Go template code generated by templ for signup.templ.

#### `/home/desktop/especa/beta/projects/superadmin/ui/static/css/auth.css`
- **Relative Path**: `projects/superadmin/ui/static/css/auth.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for auth.css.

#### `/home/desktop/especa/beta/projects/superadmin/ui/static/css/crm_shared.css`
- **Relative Path**: `projects/superadmin/ui/static/css/crm_shared.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for crm_shared.css.

#### `/home/desktop/especa/beta/projects/superadmin/ui/static/css/dashboard.css`
- **Relative Path**: `projects/superadmin/ui/static/css/dashboard.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for dashboard.css.

#### `/home/desktop/especa/beta/projects/superadmin/ui/static/css/globals.css`
- **Relative Path**: `projects/superadmin/ui/static/css/globals.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for globals.css.

#### `/home/desktop/especa/beta/projects/superadmin/ui/static/css/pages/billing.css`
- **Relative Path**: `projects/superadmin/ui/static/css/pages/billing.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for billing.css.

#### `/home/desktop/especa/beta/projects/superadmin/ui/static/css/pages/calendar.css`
- **Relative Path**: `projects/superadmin/ui/static/css/pages/calendar.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for calendar.css.

#### `/home/desktop/especa/beta/projects/superadmin/ui/static/css/pages/leads.css`
- **Relative Path**: `projects/superadmin/ui/static/css/pages/leads.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for leads.css.

#### `/home/desktop/especa/beta/projects/superadmin/ui/static/css/pages/projects.css`
- **Relative Path**: `projects/superadmin/ui/static/css/pages/projects.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for projects.css.

#### `/home/desktop/especa/beta/projects/superadmin/ui/static/css/pages/sessions.css`
- **Relative Path**: `projects/superadmin/ui/static/css/pages/sessions.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for sessions.css.

#### `/home/desktop/especa/beta/projects/superadmin/ui/static/css/pages/settings.css`
- **Relative Path**: `projects/superadmin/ui/static/css/pages/settings.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for settings.css.

#### `/home/desktop/especa/beta/projects/superadmin/ui/static/css/pages/settings_extra.css`
- **Relative Path**: `projects/superadmin/ui/static/css/pages/settings_extra.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for settings_extra.css.

#### `/home/desktop/especa/beta/projects/superadmin/ui/static/css/style.css`
- **Relative Path**: `projects/superadmin/ui/static/css/style.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for style.css.

#### `/home/desktop/especa/beta/projects/superadmin/ui/static/embed.go`
- **Relative Path**: `projects/superadmin/ui/static/embed.go`
- **Purpose**: Go embed directives packaging static CSS, JavaScript, and assets into the application binary.

#### `/home/desktop/especa/beta/projects/superadmin/ui/static/js/auth.js`
- **Relative Path**: `projects/superadmin/ui/static/js/auth.js`
- **Purpose**: Client-side JavaScript managing client interactions, DOM operations, or library integrations.

#### `/home/desktop/especa/beta/projects/superadmin/ui/static/js/fullcalendar.min.js`
- **Relative Path**: `projects/superadmin/ui/static/js/fullcalendar.min.js`
- **Purpose**: Client-side JavaScript managing client interactions, DOM operations, or library integrations.

#### `/home/desktop/especa/beta/projects/superadmin/ui/static/js/htmx.min.js`
- **Relative Path**: `projects/superadmin/ui/static/js/htmx.min.js`
- **Purpose**: Client-side JavaScript managing client interactions, DOM operations, or library integrations.

#### `/home/desktop/especa/beta/projects/superadmin/ui/static/js/kanban.js`
- **Relative Path**: `projects/superadmin/ui/static/js/kanban.js`
- **Purpose**: Client-side JavaScript managing client interactions, DOM operations, or library integrations.

#### `/home/desktop/especa/beta/projects/superadmin/ui/static/js/sortable.min.js`
- **Relative Path**: `projects/superadmin/ui/static/js/sortable.min.js`
- **Purpose**: Client-side JavaScript managing client interactions, DOM operations, or library integrations.


### 7. CMS Service (`projects/cms`) (41 files)

#### `/home/desktop/especa/beta/projects/cms/go.mod`
- **Relative Path**: `projects/cms/go.mod`
- **Purpose**: Go module manifest declaring dependencies and module path for projects/cms.

#### `/home/desktop/especa/beta/projects/cms/go.sum`
- **Relative Path**: `projects/cms/go.sum`
- **Purpose**: Checksum lockfile ensuring reproducible builds and dependency integrity for projects/cms.

#### `/home/desktop/especa/beta/projects/cms/handlers/auth_handler.go`
- **Relative Path**: `projects/cms/handlers/auth_handler.go`
- **Purpose**: HTTP controller handler implementing request decoding, validation, service coordination, and responses.

#### `/home/desktop/especa/beta/projects/cms/handlers/developer_api_handler.go`
- **Relative Path**: `projects/cms/handlers/developer_api_handler.go`
- **Purpose**: HTTP controller handler implementing request decoding, validation, service coordination, and responses.

#### `/home/desktop/especa/beta/projects/cms/handlers/home_handler.go`
- **Relative Path**: `projects/cms/handlers/home_handler.go`
- **Purpose**: HTTP controller handler implementing request decoding, validation, service coordination, and responses.

#### `/home/desktop/especa/beta/projects/cms/main.go`
- **Relative Path**: `projects/cms/main.go`
- **Purpose**: Main entrypoint initializing database connections, service dependencies, HTTP routers, and listeners for projects/cms.

#### `/home/desktop/especa/beta/projects/cms/models/models.go`
- **Relative Path**: `projects/cms/models/models.go`
- **Purpose**: Domain models, database entities, and Data Transfer Objects (DTOs).

#### `/home/desktop/especa/beta/projects/cms/repositories/cms_repository.go`
- **Relative Path**: `projects/cms/repositories/cms_repository.go`
- **Purpose**: Data access repository implementing database queries and CRUD operations for domain models.

#### `/home/desktop/especa/beta/projects/cms/routes/backend.go`
- **Relative Path**: `projects/cms/routes/backend.go`
- **Purpose**: HTTP router configuration and endpoint routing definitions for projects/cms/routes.

#### `/home/desktop/especa/beta/projects/cms/routes/frontend.go`
- **Relative Path**: `projects/cms/routes/frontend.go`
- **Purpose**: HTTP router configuration and endpoint routing definitions for projects/cms/routes.

#### `/home/desktop/especa/beta/projects/cms/routes/router.go`
- **Relative Path**: `projects/cms/routes/router.go`
- **Purpose**: HTTP router configuration and endpoint routing definitions for projects/cms/routes.

#### `/home/desktop/especa/beta/projects/cms/ui/pages/dashboard.templ`
- **Relative Path**: `projects/cms/ui/pages/dashboard.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for dashboard.templ.

#### `/home/desktop/especa/beta/projects/cms/ui/pages/dashboard_templ.go`
- **Relative Path**: `projects/cms/ui/pages/dashboard_templ.go`
- **Purpose**: Compiled Go template code generated by templ for dashboard.templ.

#### `/home/desktop/especa/beta/projects/cms/ui/pages/forgot_password.templ`
- **Relative Path**: `projects/cms/ui/pages/forgot_password.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for forgot_password.templ.

#### `/home/desktop/especa/beta/projects/cms/ui/pages/forgot_password_templ.go`
- **Relative Path**: `projects/cms/ui/pages/forgot_password_templ.go`
- **Purpose**: Compiled Go template code generated by templ for forgot_password.templ.

#### `/home/desktop/especa/beta/projects/cms/ui/pages/internal_error.templ`
- **Relative Path**: `projects/cms/ui/pages/internal_error.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for internal_error.templ.

#### `/home/desktop/especa/beta/projects/cms/ui/pages/internal_error_templ.go`
- **Relative Path**: `projects/cms/ui/pages/internal_error_templ.go`
- **Purpose**: Compiled Go template code generated by templ for internal_error.templ.

#### `/home/desktop/especa/beta/projects/cms/ui/pages/not_found.templ`
- **Relative Path**: `projects/cms/ui/pages/not_found.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for not_found.templ.

#### `/home/desktop/especa/beta/projects/cms/ui/pages/not_found_templ.go`
- **Relative Path**: `projects/cms/ui/pages/not_found_templ.go`
- **Purpose**: Compiled Go template code generated by templ for not_found.templ.

#### `/home/desktop/especa/beta/projects/cms/ui/pages/signin.templ`
- **Relative Path**: `projects/cms/ui/pages/signin.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for signin.templ.

#### `/home/desktop/especa/beta/projects/cms/ui/pages/signin_templ.go`
- **Relative Path**: `projects/cms/ui/pages/signin_templ.go`
- **Purpose**: Compiled Go template code generated by templ for signin.templ.

#### `/home/desktop/especa/beta/projects/cms/ui/pages/signup.templ`
- **Relative Path**: `projects/cms/ui/pages/signup.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for signup.templ.

#### `/home/desktop/especa/beta/projects/cms/ui/pages/signup_templ.go`
- **Relative Path**: `projects/cms/ui/pages/signup_templ.go`
- **Purpose**: Compiled Go template code generated by templ for signup.templ.

#### `/home/desktop/especa/beta/projects/cms/ui/static/css/auth.css`
- **Relative Path**: `projects/cms/ui/static/css/auth.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for auth.css.

#### `/home/desktop/especa/beta/projects/cms/ui/static/css/crm_shared.css`
- **Relative Path**: `projects/cms/ui/static/css/crm_shared.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for crm_shared.css.

#### `/home/desktop/especa/beta/projects/cms/ui/static/css/dashboard.css`
- **Relative Path**: `projects/cms/ui/static/css/dashboard.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for dashboard.css.

#### `/home/desktop/especa/beta/projects/cms/ui/static/css/globals.css`
- **Relative Path**: `projects/cms/ui/static/css/globals.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for globals.css.

#### `/home/desktop/especa/beta/projects/cms/ui/static/css/pages/billing.css`
- **Relative Path**: `projects/cms/ui/static/css/pages/billing.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for billing.css.

#### `/home/desktop/especa/beta/projects/cms/ui/static/css/pages/calendar.css`
- **Relative Path**: `projects/cms/ui/static/css/pages/calendar.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for calendar.css.

#### `/home/desktop/especa/beta/projects/cms/ui/static/css/pages/leads.css`
- **Relative Path**: `projects/cms/ui/static/css/pages/leads.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for leads.css.

#### `/home/desktop/especa/beta/projects/cms/ui/static/css/pages/projects.css`
- **Relative Path**: `projects/cms/ui/static/css/pages/projects.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for projects.css.

#### `/home/desktop/especa/beta/projects/cms/ui/static/css/pages/sessions.css`
- **Relative Path**: `projects/cms/ui/static/css/pages/sessions.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for sessions.css.

#### `/home/desktop/especa/beta/projects/cms/ui/static/css/pages/settings.css`
- **Relative Path**: `projects/cms/ui/static/css/pages/settings.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for settings.css.

#### `/home/desktop/especa/beta/projects/cms/ui/static/css/pages/settings_extra.css`
- **Relative Path**: `projects/cms/ui/static/css/pages/settings_extra.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for settings_extra.css.

#### `/home/desktop/especa/beta/projects/cms/ui/static/css/style.css`
- **Relative Path**: `projects/cms/ui/static/css/style.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for style.css.

#### `/home/desktop/especa/beta/projects/cms/ui/static/embed.go`
- **Relative Path**: `projects/cms/ui/static/embed.go`
- **Purpose**: Go embed directives packaging static CSS, JavaScript, and assets into the application binary.

#### `/home/desktop/especa/beta/projects/cms/ui/static/js/auth.js`
- **Relative Path**: `projects/cms/ui/static/js/auth.js`
- **Purpose**: Client-side JavaScript managing client interactions, DOM operations, or library integrations.

#### `/home/desktop/especa/beta/projects/cms/ui/static/js/fullcalendar.min.js`
- **Relative Path**: `projects/cms/ui/static/js/fullcalendar.min.js`
- **Purpose**: Client-side JavaScript managing client interactions, DOM operations, or library integrations.

#### `/home/desktop/especa/beta/projects/cms/ui/static/js/htmx.min.js`
- **Relative Path**: `projects/cms/ui/static/js/htmx.min.js`
- **Purpose**: Client-side JavaScript managing client interactions, DOM operations, or library integrations.

#### `/home/desktop/especa/beta/projects/cms/ui/static/js/kanban.js`
- **Relative Path**: `projects/cms/ui/static/js/kanban.js`
- **Purpose**: Client-side JavaScript managing client interactions, DOM operations, or library integrations.

#### `/home/desktop/especa/beta/projects/cms/ui/static/js/sortable.min.js`
- **Relative Path**: `projects/cms/ui/static/js/sortable.min.js`
- **Purpose**: Client-side JavaScript managing client interactions, DOM operations, or library integrations.


### 8. Photography CRM (`projects/photography_crm`) (28 files)

#### `/home/desktop/especa/beta/projects/photography_crm/.env`
- **Relative Path**: `projects/photography_crm/.env`
- **Purpose**: Local environment variables configuration file specifying ports and secrets.

#### `/home/desktop/especa/beta/projects/photography_crm/.env.example`
- **Relative Path**: `projects/photography_crm/.env.example`
- **Purpose**: Template environment variables documentation file.

#### `/home/desktop/especa/beta/projects/photography_crm/go.mod`
- **Relative Path**: `projects/photography_crm/go.mod`
- **Purpose**: Go module manifest declaring dependencies and module path for projects/photography_crm.

#### `/home/desktop/especa/beta/projects/photography_crm/go.sum`
- **Relative Path**: `projects/photography_crm/go.sum`
- **Purpose**: Checksum lockfile ensuring reproducible builds and dependency integrity for projects/photography_crm.

#### `/home/desktop/especa/beta/projects/photography_crm/main.go`
- **Relative Path**: `projects/photography_crm/main.go`
- **Purpose**: Main entrypoint booting the photography CRM microservice on port 8082 via core.Run.

#### `/home/desktop/especa/beta/projects/photography_crm/handlers/auth_handler.go`
- **Relative Path**: `projects/photography_crm/handlers/auth_handler.go`
- **Purpose**: HTTP controller handler implementing studio session authentication and signin/signup flows.

#### `/home/desktop/especa/beta/projects/photography_crm/handlers/billing_handler.go`
- **Relative Path**: `projects/photography_crm/handlers/billing_handler.go`
- **Purpose**: Direct MySQL handlers for financial payments, invoices, expenses, refunds, and cancellations.

#### `/home/desktop/especa/beta/projects/photography_crm/handlers/crew_handler.go`
- **Relative Path**: `projects/photography_crm/handlers/crew_handler.go`
- **Purpose**: Direct MySQL handlers managing studio crew roster, photographers, and assignments.

#### `/home/desktop/especa/beta/projects/photography_crm/handlers/deliverables_handler.go`
- **Relative Path**: `projects/photography_crm/handlers/deliverables_handler.go`
- **Purpose**: Handlers managing deliverable groups and milestone templates with batch update support.

#### `/home/desktop/especa/beta/projects/photography_crm/handlers/home_handler.go`
- **Relative Path**: `projects/photography_crm/handlers/home_handler.go`
- **Purpose**: HTTP handler for dashboard and home operational views.

#### `/home/desktop/especa/beta/projects/photography_crm/handlers/leads_handler.go`
- **Relative Path**: `projects/photography_crm/handlers/leads_handler.go`
- **Purpose**: Direct MySQL CRUD controller for inquiries, leads pipeline, Kanban stage updates, and lead-to-client conversion.

#### `/home/desktop/especa/beta/projects/photography_crm/handlers/order_presets_handler.go`
- **Relative Path**: `projects/photography_crm/handlers/order_presets_handler.go`
- **Purpose**: Handlers for merchandise product presets (albums, frames, canvas) with batch update support.

#### `/home/desktop/especa/beta/projects/photography_crm/handlers/orders_handler.go`
- **Relative Path**: `projects/photography_crm/handlers/orders_handler.go`
- **Purpose**: Direct MySQL CRUD controller for merchandise and print orders and manufacturing stages.

#### `/home/desktop/especa/beta/projects/photography_crm/handlers/packages_handler.go`
- **Relative Path**: `projects/photography_crm/handlers/packages_handler.go`
- **Purpose**: Handlers managing packaged service bundles and pricing tiers.

#### `/home/desktop/especa/beta/projects/photography_crm/handlers/projects_handler.go`
- **Relative Path**: `projects/photography_crm/handlers/projects_handler.go`
- **Purpose**: Direct MySQL CRUD controller for studio photoshoot projects, multi-day segments, milestones, and expenses.

#### `/home/desktop/especa/beta/projects/photography_crm/handlers/registration_handler.go`
- **Relative Path**: `projects/photography_crm/handlers/registration_handler.go`
- **Purpose**: HTTP handler submitting tenant registration onboarding requests for administrative review.

#### `/home/desktop/especa/beta/projects/photography_crm/handlers/roles_handler.go`
- **Relative Path**: `projects/photography_crm/handlers/roles_handler.go`
- **Purpose**: Handlers managing studio roles (photographer, cinematographer, editor, assistant).

#### `/home/desktop/especa/beta/projects/photography_crm/handlers/sequence_id_handler.go`
- **Relative Path**: `projects/photography_crm/handlers/sequence_id_handler.go`
- **Purpose**: Handlers retrieving and saving custom document sequence ID numbering and prefixes.

#### `/home/desktop/especa/beta/projects/photography_crm/handlers/services_handler.go`
- **Relative Path**: `projects/photography_crm/handlers/services_handler.go`
- **Purpose**: Handlers managing photography service categories and pricing with batch update support.

#### `/home/desktop/especa/beta/projects/photography_crm/handlers/sessions_handler.go`
- **Relative Path**: `projects/photography_crm/handlers/sessions_handler.go`
- **Purpose**: Direct MySQL CRUD controller for photoshoot session scheduling, timings, location, and crew allocation.

#### `/home/desktop/especa/beta/projects/photography_crm/handlers/settings_handler.go`
- **Relative Path**: `projects/photography_crm/handlers/settings_handler.go`
- **Purpose**: Atomic JSON key-value store handlers for studio settings, templates, and presets.

#### `/home/desktop/especa/beta/projects/photography_crm/handlers/workflow_handler.go`
- **Relative Path**: `projects/photography_crm/handlers/workflow_handler.go`
- **Purpose**: Handlers managing production workflow pipelines, stage color coding, and batch ordering.

#### `/home/desktop/especa/beta/projects/photography_crm/handlers/workspace_handler.go`
- **Relative Path**: `projects/photography_crm/handlers/workspace_handler.go`
- **Purpose**: HTTP handler listing user memberships and accessible studio workspaces.

#### `/home/desktop/especa/beta/projects/photography_crm/repositories/crm_repository.go`
- **Relative Path**: `projects/photography_crm/repositories/crm_repository.go`
- **Purpose**: Data access repository implementing database queries and CRUD operations for domain models.

#### `/home/desktop/especa/beta/projects/photography_crm/routes/backend.go`
- **Relative Path**: `projects/photography_crm/routes/backend.go`
- **Purpose**: HTTP router configuration and endpoint routing definitions for REST APIs and mutation routes.

#### `/home/desktop/especa/beta/projects/photography_crm/routes/frontend.go`
- **Relative Path**: `projects/photography_crm/routes/frontend.go`
- **Purpose**: HTTP router configuration serving the compiled React SPA from frontend/apps/photography_crm/dist.

#### `/home/desktop/especa/beta/projects/photography_crm/routes/router.go`
- **Relative Path**: `projects/photography_crm/routes/router.go`
- **Purpose**: Top-level HTTP middleware assembly and router constructor for Photography CRM.

#### `/home/desktop/especa/beta/projects/photography_crm/services/crm_service.go`
- **Relative Path**: `projects/photography_crm/services/crm_service.go`
- **Purpose**: Business logic service layer handling core workflows, validation rules, and operations.


### 9. Especa CRM (`projects/especa_crm`) (39 files)

#### `/home/desktop/especa/beta/projects/especa_crm/db/db.go`
- **Relative Path**: `projects/especa_crm/db/db.go`
- **Purpose**: Source and configuration file delivering essential logic for db.go.

#### `/home/desktop/especa/beta/projects/especa_crm/db/migrations/000001_init.sql`
- **Relative Path**: `projects/especa_crm/db/migrations/000001_init.sql`
- **Purpose**: Database migration or schema DDL script declaring table structure and indexes for 000001_init.sql.

#### `/home/desktop/especa/beta/projects/especa_crm/go.mod`
- **Relative Path**: `projects/especa_crm/go.mod`
- **Purpose**: Go module manifest declaring dependencies and module path for projects/especa_crm.

#### `/home/desktop/especa/beta/projects/especa_crm/go.sum`
- **Relative Path**: `projects/especa_crm/go.sum`
- **Purpose**: Checksum lockfile ensuring reproducible builds and dependency integrity for projects/especa_crm.

#### `/home/desktop/especa/beta/projects/especa_crm/handlers/auth_handlers.go`
- **Relative Path**: `projects/especa_crm/handlers/auth_handlers.go`
- **Purpose**: HTTP controller handler implementing request decoding, validation, service coordination, and responses.

#### `/home/desktop/especa/beta/projects/especa_crm/handlers/handlers.go`
- **Relative Path**: `projects/especa_crm/handlers/handlers.go`
- **Purpose**: HTTP controller handler implementing request decoding, validation, service coordination, and responses.

#### `/home/desktop/especa/beta/projects/especa_crm/main.go`
- **Relative Path**: `projects/especa_crm/main.go`
- **Purpose**: Main entrypoint initializing database connections, service dependencies, HTTP routers, and listeners for projects/especa_crm.

#### `/home/desktop/especa/beta/projects/especa_crm/models/models.go`
- **Relative Path**: `projects/especa_crm/models/models.go`
- **Purpose**: Domain models, database entities, and Data Transfer Objects (DTOs).

#### `/home/desktop/especa/beta/projects/especa_crm/routes/routes.go`
- **Relative Path**: `projects/especa_crm/routes/routes.go`
- **Purpose**: HTTP router configuration and endpoint routing definitions for projects/especa_crm/routes.

#### `/home/desktop/especa/beta/projects/especa_crm/ui/layouts/layout.templ`
- **Relative Path**: `projects/especa_crm/ui/layouts/layout.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for layout.templ.

#### `/home/desktop/especa/beta/projects/especa_crm/ui/layouts/layout_templ.go`
- **Relative Path**: `projects/especa_crm/ui/layouts/layout_templ.go`
- **Purpose**: Compiled Go template code generated by templ for layout.templ.

#### `/home/desktop/especa/beta/projects/especa_crm/ui/pages/dashboard.templ`
- **Relative Path**: `projects/especa_crm/ui/pages/dashboard.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for dashboard.templ.

#### `/home/desktop/especa/beta/projects/especa_crm/ui/pages/dashboard_templ.go`
- **Relative Path**: `projects/especa_crm/ui/pages/dashboard_templ.go`
- **Purpose**: Compiled Go template code generated by templ for dashboard.templ.

#### `/home/desktop/especa/beta/projects/especa_crm/ui/pages/deals.templ`
- **Relative Path**: `projects/especa_crm/ui/pages/deals.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for deals.templ.

#### `/home/desktop/especa/beta/projects/especa_crm/ui/pages/deals_templ.go`
- **Relative Path**: `projects/especa_crm/ui/pages/deals_templ.go`
- **Purpose**: Compiled Go template code generated by templ for deals.templ.

#### `/home/desktop/especa/beta/projects/especa_crm/ui/pages/leads.templ`
- **Relative Path**: `projects/especa_crm/ui/pages/leads.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for leads.templ.

#### `/home/desktop/especa/beta/projects/especa_crm/ui/pages/leads_templ.go`
- **Relative Path**: `projects/especa_crm/ui/pages/leads_templ.go`
- **Purpose**: Compiled Go template code generated by templ for leads.templ.

#### `/home/desktop/especa/beta/projects/especa_crm/ui/pages/projects.templ`
- **Relative Path**: `projects/especa_crm/ui/pages/projects.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for projects.templ.

#### `/home/desktop/especa/beta/projects/especa_crm/ui/pages/projects_templ.go`
- **Relative Path**: `projects/especa_crm/ui/pages/projects_templ.go`
- **Purpose**: Compiled Go template code generated by templ for projects.templ.

#### `/home/desktop/especa/beta/projects/especa_crm/ui/pages/quotes_contacts.templ`
- **Relative Path**: `projects/especa_crm/ui/pages/quotes_contacts.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for quotes_contacts.templ.

#### `/home/desktop/especa/beta/projects/especa_crm/ui/pages/quotes_contacts_templ.go`
- **Relative Path**: `projects/especa_crm/ui/pages/quotes_contacts_templ.go`
- **Purpose**: Compiled Go template code generated by templ for quotes_contacts.templ.

#### `/home/desktop/especa/beta/projects/especa_crm/ui/static/css/auth.css`
- **Relative Path**: `projects/especa_crm/ui/static/css/auth.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for auth.css.

#### `/home/desktop/especa/beta/projects/especa_crm/ui/static/css/crm_shared.css`
- **Relative Path**: `projects/especa_crm/ui/static/css/crm_shared.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for crm_shared.css.

#### `/home/desktop/especa/beta/projects/especa_crm/ui/static/css/dashboard.css`
- **Relative Path**: `projects/especa_crm/ui/static/css/dashboard.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for dashboard.css.

#### `/home/desktop/especa/beta/projects/especa_crm/ui/static/css/globals.css`
- **Relative Path**: `projects/especa_crm/ui/static/css/globals.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for globals.css.

#### `/home/desktop/especa/beta/projects/especa_crm/ui/static/css/pages/billing.css`
- **Relative Path**: `projects/especa_crm/ui/static/css/pages/billing.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for billing.css.

#### `/home/desktop/especa/beta/projects/especa_crm/ui/static/css/pages/calendar.css`
- **Relative Path**: `projects/especa_crm/ui/static/css/pages/calendar.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for calendar.css.

#### `/home/desktop/especa/beta/projects/especa_crm/ui/static/css/pages/leads.css`
- **Relative Path**: `projects/especa_crm/ui/static/css/pages/leads.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for leads.css.

#### `/home/desktop/especa/beta/projects/especa_crm/ui/static/css/pages/projects.css`
- **Relative Path**: `projects/especa_crm/ui/static/css/pages/projects.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for projects.css.

#### `/home/desktop/especa/beta/projects/especa_crm/ui/static/css/pages/sessions.css`
- **Relative Path**: `projects/especa_crm/ui/static/css/pages/sessions.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for sessions.css.

#### `/home/desktop/especa/beta/projects/especa_crm/ui/static/css/pages/settings.css`
- **Relative Path**: `projects/especa_crm/ui/static/css/pages/settings.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for settings.css.

#### `/home/desktop/especa/beta/projects/especa_crm/ui/static/css/pages/settings_extra.css`
- **Relative Path**: `projects/especa_crm/ui/static/css/pages/settings_extra.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for settings_extra.css.

#### `/home/desktop/especa/beta/projects/especa_crm/ui/static/css/style.css`
- **Relative Path**: `projects/especa_crm/ui/static/css/style.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for style.css.

#### `/home/desktop/especa/beta/projects/especa_crm/ui/static/embed.go`
- **Relative Path**: `projects/especa_crm/ui/static/embed.go`
- **Purpose**: Go embed directives packaging static CSS, JavaScript, and assets into the application binary.

#### `/home/desktop/especa/beta/projects/especa_crm/ui/static/js/auth.js`
- **Relative Path**: `projects/especa_crm/ui/static/js/auth.js`
- **Purpose**: Client-side JavaScript managing client interactions, DOM operations, or library integrations.

#### `/home/desktop/especa/beta/projects/especa_crm/ui/static/js/fullcalendar.min.js`
- **Relative Path**: `projects/especa_crm/ui/static/js/fullcalendar.min.js`
- **Purpose**: Client-side JavaScript managing client interactions, DOM operations, or library integrations.

#### `/home/desktop/especa/beta/projects/especa_crm/ui/static/js/htmx.min.js`
- **Relative Path**: `projects/especa_crm/ui/static/js/htmx.min.js`
- **Purpose**: Client-side JavaScript managing client interactions, DOM operations, or library integrations.

#### `/home/desktop/especa/beta/projects/especa_crm/ui/static/js/kanban.js`
- **Relative Path**: `projects/especa_crm/ui/static/js/kanban.js`
- **Purpose**: Client-side JavaScript managing client interactions, DOM operations, or library integrations.

#### `/home/desktop/especa/beta/projects/especa_crm/ui/static/js/sortable.min.js`
- **Relative Path**: `projects/especa_crm/ui/static/js/sortable.min.js`
- **Purpose**: Client-side JavaScript managing client interactions, DOM operations, or library integrations.


### 10. Especa Finance (`projects/especa_finance`) (88 files)

#### `/home/desktop/especa/beta/projects/especa_finance/db/db.go`
- **Relative Path**: `projects/especa_finance/db/db.go`
- **Purpose**: Source and configuration file delivering essential logic for db.go.

#### `/home/desktop/especa/beta/projects/especa_finance/db/migrations/000001_init.sql`
- **Relative Path**: `projects/especa_finance/db/migrations/000001_init.sql`
- **Purpose**: Database migration or schema DDL script declaring table structure and indexes for 000001_init.sql.

#### `/home/desktop/especa/beta/projects/especa_finance/db/migrations/000002_especa_finance_plan.sql`
- **Relative Path**: `projects/especa_finance/db/migrations/000002_especa_finance_plan.sql`
- **Purpose**: Database migration or schema DDL script declaring table structure and indexes for 000002_especa_finance_plan.sql.

#### `/home/desktop/especa/beta/projects/especa_finance/go.mod`
- **Relative Path**: `projects/especa_finance/go.mod`
- **Purpose**: Go module manifest declaring dependencies and module path for projects/especa_finance.

#### `/home/desktop/especa/beta/projects/especa_finance/go.sum`
- **Relative Path**: `projects/especa_finance/go.sum`
- **Purpose**: Checksum lockfile ensuring reproducible builds and dependency integrity for projects/especa_finance.

#### `/home/desktop/especa/beta/projects/especa_finance/handlers/auth_handlers.go`
- **Relative Path**: `projects/especa_finance/handlers/auth_handlers.go`
- **Purpose**: HTTP controller handler implementing request decoding, validation, service coordination, and responses.

#### `/home/desktop/especa/beta/projects/especa_finance/handlers/handlers.go`
- **Relative Path**: `projects/especa_finance/handlers/handlers.go`
- **Purpose**: HTTP controller handler implementing request decoding, validation, service coordination, and responses.

#### `/home/desktop/especa/beta/projects/especa_finance/main.go`
- **Relative Path**: `projects/especa_finance/main.go`
- **Purpose**: Main entrypoint initializing database connections, service dependencies, HTTP routers, and listeners for projects/especa_finance.

#### `/home/desktop/especa/beta/projects/especa_finance/models/models.go`
- **Relative Path**: `projects/especa_finance/models/models.go`
- **Purpose**: Domain models, database entities, and Data Transfer Objects (DTOs).

#### `/home/desktop/especa/beta/projects/especa_finance/routes/routes.go`
- **Relative Path**: `projects/especa_finance/routes/routes.go`
- **Purpose**: HTTP router configuration and endpoint routing definitions for projects/especa_finance/routes.

#### `/home/desktop/especa/beta/projects/especa_finance/services/accounting.go`
- **Relative Path**: `projects/especa_finance/services/accounting.go`
- **Purpose**: Business logic service layer handling core workflows, validation rules, and operations.

#### `/home/desktop/especa/beta/projects/especa_finance/services/financial_calculator.go`
- **Relative Path**: `projects/especa_finance/services/financial_calculator.go`
- **Purpose**: Business logic service layer handling core workflows, validation rules, and operations.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/layouts/layout.templ`
- **Relative Path**: `projects/especa_finance/ui/layouts/layout.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for layout.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/layouts/layout_templ.go`
- **Relative Path**: `projects/especa_finance/ui/layouts/layout_templ.go`
- **Purpose**: Compiled Go template code generated by templ for layout.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/analytics_data.go`
- **Relative Path**: `projects/especa_finance/ui/pages/analytics_data.go`
- **Purpose**: Source and configuration file delivering essential logic for analytics_data.go.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/analytics_view.templ`
- **Relative Path**: `projects/especa_finance/ui/pages/analytics_view.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for analytics_view.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/analytics_view_templ.go`
- **Relative Path**: `projects/especa_finance/ui/pages/analytics_view_templ.go`
- **Purpose**: Compiled Go template code generated by templ for analytics_view.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/bill_create.templ`
- **Relative Path**: `projects/especa_finance/ui/pages/bill_create.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for bill_create.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/bill_create_templ.go`
- **Relative Path**: `projects/especa_finance/ui/pages/bill_create_templ.go`
- **Purpose**: Compiled Go template code generated by templ for bill_create.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/bills_list.templ`
- **Relative Path**: `projects/especa_finance/ui/pages/bills_list.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for bills_list.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/bills_list_templ.go`
- **Relative Path**: `projects/especa_finance/ui/pages/bills_list_templ.go`
- **Purpose**: Compiled Go template code generated by templ for bills_list.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/chart_of_accounts.templ`
- **Relative Path**: `projects/especa_finance/ui/pages/chart_of_accounts.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for chart_of_accounts.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/chart_of_accounts_templ.go`
- **Relative Path**: `projects/especa_finance/ui/pages/chart_of_accounts_templ.go`
- **Purpose**: Compiled Go template code generated by templ for chart_of_accounts.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/client_create.templ`
- **Relative Path**: `projects/especa_finance/ui/pages/client_create.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for client_create.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/client_create_templ.go`
- **Relative Path**: `projects/especa_finance/ui/pages/client_create_templ.go`
- **Purpose**: Compiled Go template code generated by templ for client_create.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/clients_list.templ`
- **Relative Path**: `projects/especa_finance/ui/pages/clients_list.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for clients_list.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/clients_list_templ.go`
- **Relative Path**: `projects/especa_finance/ui/pages/clients_list_templ.go`
- **Purpose**: Compiled Go template code generated by templ for clients_list.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/employee_create.templ`
- **Relative Path**: `projects/especa_finance/ui/pages/employee_create.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for employee_create.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/employee_create_templ.go`
- **Relative Path**: `projects/especa_finance/ui/pages/employee_create_templ.go`
- **Purpose**: Compiled Go template code generated by templ for employee_create.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/employee_detail.templ`
- **Relative Path**: `projects/especa_finance/ui/pages/employee_detail.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for employee_detail.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/employee_detail_templ.go`
- **Relative Path**: `projects/especa_finance/ui/pages/employee_detail_templ.go`
- **Purpose**: Compiled Go template code generated by templ for employee_detail.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/employee_templates_list.templ`
- **Relative Path**: `projects/especa_finance/ui/pages/employee_templates_list.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for employee_templates_list.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/employee_templates_list_templ.go`
- **Relative Path**: `projects/especa_finance/ui/pages/employee_templates_list_templ.go`
- **Purpose**: Compiled Go template code generated by templ for employee_templates_list.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/employees_list.templ`
- **Relative Path**: `projects/especa_finance/ui/pages/employees_list.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for employees_list.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/employees_list_templ.go`
- **Relative Path**: `projects/especa_finance/ui/pages/employees_list_templ.go`
- **Purpose**: Compiled Go template code generated by templ for employees_list.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/financial_statements.templ`
- **Relative Path**: `projects/especa_finance/ui/pages/financial_statements.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for financial_statements.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/financial_statements_templ.go`
- **Relative Path**: `projects/especa_finance/ui/pages/financial_statements_templ.go`
- **Purpose**: Compiled Go template code generated by templ for financial_statements.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/invoice_create.templ`
- **Relative Path**: `projects/especa_finance/ui/pages/invoice_create.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for invoice_create.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/invoice_create_templ.go`
- **Relative Path**: `projects/especa_finance/ui/pages/invoice_create_templ.go`
- **Purpose**: Compiled Go template code generated by templ for invoice_create.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/invoices_list.templ`
- **Relative Path**: `projects/especa_finance/ui/pages/invoices_list.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for invoices_list.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/invoices_list_templ.go`
- **Relative Path**: `projects/especa_finance/ui/pages/invoices_list_templ.go`
- **Purpose**: Compiled Go template code generated by templ for invoices_list.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/journal_entries.templ`
- **Relative Path**: `projects/especa_finance/ui/pages/journal_entries.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for journal_entries.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/journal_entries_templ.go`
- **Relative Path**: `projects/especa_finance/ui/pages/journal_entries_templ.go`
- **Purpose**: Compiled Go template code generated by templ for journal_entries.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/journal_entry_create.templ`
- **Relative Path**: `projects/especa_finance/ui/pages/journal_entry_create.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for journal_entry_create.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/journal_entry_create_templ.go`
- **Relative Path**: `projects/especa_finance/ui/pages/journal_entry_create_templ.go`
- **Purpose**: Compiled Go template code generated by templ for journal_entry_create.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/overview.templ`
- **Relative Path**: `projects/especa_finance/ui/pages/overview.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for overview.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/overview_templ.go`
- **Relative Path**: `projects/especa_finance/ui/pages/overview_templ.go`
- **Purpose**: Compiled Go template code generated by templ for overview.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/payment_create.templ`
- **Relative Path**: `projects/especa_finance/ui/pages/payment_create.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for payment_create.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/payment_create_templ.go`
- **Relative Path**: `projects/especa_finance/ui/pages/payment_create_templ.go`
- **Purpose**: Compiled Go template code generated by templ for payment_create.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/payments_list.templ`
- **Relative Path**: `projects/especa_finance/ui/pages/payments_list.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for payments_list.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/payments_list_templ.go`
- **Relative Path**: `projects/especa_finance/ui/pages/payments_list_templ.go`
- **Purpose**: Compiled Go template code generated by templ for payments_list.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/payroll_register.templ`
- **Relative Path**: `projects/especa_finance/ui/pages/payroll_register.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for payroll_register.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/payroll_register_templ.go`
- **Relative Path**: `projects/especa_finance/ui/pages/payroll_register_templ.go`
- **Purpose**: Compiled Go template code generated by templ for payroll_register.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/product_create.templ`
- **Relative Path**: `projects/especa_finance/ui/pages/product_create.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for product_create.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/product_create_templ.go`
- **Relative Path**: `projects/especa_finance/ui/pages/product_create_templ.go`
- **Purpose**: Compiled Go template code generated by templ for product_create.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/product_detail.templ`
- **Relative Path**: `projects/especa_finance/ui/pages/product_detail.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for product_detail.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/product_detail_templ.go`
- **Relative Path**: `projects/especa_finance/ui/pages/product_detail_templ.go`
- **Purpose**: Compiled Go template code generated by templ for product_detail.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/product_templates_list.templ`
- **Relative Path**: `projects/especa_finance/ui/pages/product_templates_list.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for product_templates_list.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/product_templates_list_templ.go`
- **Relative Path**: `projects/especa_finance/ui/pages/product_templates_list_templ.go`
- **Purpose**: Compiled Go template code generated by templ for product_templates_list.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/products_list.templ`
- **Relative Path**: `projects/especa_finance/ui/pages/products_list.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for products_list.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/products_list_templ.go`
- **Relative Path**: `projects/especa_finance/ui/pages/products_list_templ.go`
- **Purpose**: Compiled Go template code generated by templ for products_list.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/project_financials.templ`
- **Relative Path**: `projects/especa_finance/ui/pages/project_financials.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for project_financials.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/project_financials_templ.go`
- **Relative Path**: `projects/especa_finance/ui/pages/project_financials_templ.go`
- **Purpose**: Compiled Go template code generated by templ for project_financials.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/subscriptions.templ`
- **Relative Path**: `projects/especa_finance/ui/pages/subscriptions.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for subscriptions.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/subscriptions_templ.go`
- **Relative Path**: `projects/especa_finance/ui/pages/subscriptions_templ.go`
- **Purpose**: Compiled Go template code generated by templ for subscriptions.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/tags_and_tax_settings.templ`
- **Relative Path**: `projects/especa_finance/ui/pages/tags_and_tax_settings.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for tags_and_tax_settings.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/tags_and_tax_settings_templ.go`
- **Relative Path**: `projects/especa_finance/ui/pages/tags_and_tax_settings_templ.go`
- **Purpose**: Compiled Go template code generated by templ for tags_and_tax_settings.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/template_create.templ`
- **Relative Path**: `projects/especa_finance/ui/pages/template_create.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for template_create.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/pages/template_create_templ.go`
- **Relative Path**: `projects/especa_finance/ui/pages/template_create_templ.go`
- **Purpose**: Compiled Go template code generated by templ for template_create.templ.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/static/css/auth.css`
- **Relative Path**: `projects/especa_finance/ui/static/css/auth.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for auth.css.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/static/css/crm_shared.css`
- **Relative Path**: `projects/especa_finance/ui/static/css/crm_shared.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for crm_shared.css.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/static/css/dashboard.css`
- **Relative Path**: `projects/especa_finance/ui/static/css/dashboard.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for dashboard.css.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/static/css/finance.css`
- **Relative Path**: `projects/especa_finance/ui/static/css/finance.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for finance.css.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/static/css/globals.css`
- **Relative Path**: `projects/especa_finance/ui/static/css/globals.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for globals.css.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/static/css/pages/billing.css`
- **Relative Path**: `projects/especa_finance/ui/static/css/pages/billing.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for billing.css.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/static/css/pages/calendar.css`
- **Relative Path**: `projects/especa_finance/ui/static/css/pages/calendar.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for calendar.css.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/static/css/pages/leads.css`
- **Relative Path**: `projects/especa_finance/ui/static/css/pages/leads.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for leads.css.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/static/css/pages/projects.css`
- **Relative Path**: `projects/especa_finance/ui/static/css/pages/projects.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for projects.css.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/static/css/pages/sessions.css`
- **Relative Path**: `projects/especa_finance/ui/static/css/pages/sessions.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for sessions.css.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/static/css/pages/settings.css`
- **Relative Path**: `projects/especa_finance/ui/static/css/pages/settings.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for settings.css.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/static/css/pages/settings_extra.css`
- **Relative Path**: `projects/especa_finance/ui/static/css/pages/settings_extra.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for settings_extra.css.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/static/css/style.css`
- **Relative Path**: `projects/especa_finance/ui/static/css/style.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for style.css.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/static/embed.go`
- **Relative Path**: `projects/especa_finance/ui/static/embed.go`
- **Purpose**: Go embed directives packaging static CSS, JavaScript, and assets into the application binary.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/static/js/auth.js`
- **Relative Path**: `projects/especa_finance/ui/static/js/auth.js`
- **Purpose**: Client-side JavaScript managing client interactions, DOM operations, or library integrations.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/static/js/fullcalendar.min.js`
- **Relative Path**: `projects/especa_finance/ui/static/js/fullcalendar.min.js`
- **Purpose**: Client-side JavaScript managing client interactions, DOM operations, or library integrations.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/static/js/htmx.min.js`
- **Relative Path**: `projects/especa_finance/ui/static/js/htmx.min.js`
- **Purpose**: Client-side JavaScript managing client interactions, DOM operations, or library integrations.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/static/js/kanban.js`
- **Relative Path**: `projects/especa_finance/ui/static/js/kanban.js`
- **Purpose**: Client-side JavaScript managing client interactions, DOM operations, or library integrations.

#### `/home/desktop/especa/beta/projects/especa_finance/ui/static/js/sortable.min.js`
- **Relative Path**: `projects/especa_finance/ui/static/js/sortable.min.js`
- **Purpose**: Client-side JavaScript managing client interactions, DOM operations, or library integrations.


### 11. HRMS Service (`projects/hrms`) (45 files)

#### `/home/desktop/especa/beta/projects/hrms/db/store.go`
- **Relative Path**: `projects/hrms/db/store.go`
- **Purpose**: Source and configuration file delivering essential logic for store.go.

#### `/home/desktop/especa/beta/projects/hrms/go.mod`
- **Relative Path**: `projects/hrms/go.mod`
- **Purpose**: Go module manifest declaring dependencies and module path for projects/hrms.

#### `/home/desktop/especa/beta/projects/hrms/go.sum`
- **Relative Path**: `projects/hrms/go.sum`
- **Purpose**: Checksum lockfile ensuring reproducible builds and dependency integrity for projects/hrms.

#### `/home/desktop/especa/beta/projects/hrms/handlers/handlers.go`
- **Relative Path**: `projects/hrms/handlers/handlers.go`
- **Purpose**: HTTP controller handler implementing request decoding, validation, service coordination, and responses.

#### `/home/desktop/especa/beta/projects/hrms/main.go`
- **Relative Path**: `projects/hrms/main.go`
- **Purpose**: Main entrypoint initializing database connections, service dependencies, HTTP routers, and listeners for projects/hrms.

#### `/home/desktop/especa/beta/projects/hrms/models/models.go`
- **Relative Path**: `projects/hrms/models/models.go`
- **Purpose**: Domain models, database entities, and Data Transfer Objects (DTOs).

#### `/home/desktop/especa/beta/projects/hrms/routes/routes.go`
- **Relative Path**: `projects/hrms/routes/routes.go`
- **Purpose**: HTTP router configuration and endpoint routing definitions for projects/hrms/routes.

#### `/home/desktop/especa/beta/projects/hrms/ui/layouts/layout.templ`
- **Relative Path**: `projects/hrms/ui/layouts/layout.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for layout.templ.

#### `/home/desktop/especa/beta/projects/hrms/ui/layouts/layout_templ.go`
- **Relative Path**: `projects/hrms/ui/layouts/layout_templ.go`
- **Purpose**: Compiled Go template code generated by templ for layout.templ.

#### `/home/desktop/especa/beta/projects/hrms/ui/pages/dashboard.templ`
- **Relative Path**: `projects/hrms/ui/pages/dashboard.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for dashboard.templ.

#### `/home/desktop/especa/beta/projects/hrms/ui/pages/dashboard_templ.go`
- **Relative Path**: `projects/hrms/ui/pages/dashboard_templ.go`
- **Purpose**: Compiled Go template code generated by templ for dashboard.templ.

#### `/home/desktop/especa/beta/projects/hrms/ui/pages/deductions_manager.templ`
- **Relative Path**: `projects/hrms/ui/pages/deductions_manager.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for deductions_manager.templ.

#### `/home/desktop/especa/beta/projects/hrms/ui/pages/deductions_manager_templ.go`
- **Relative Path**: `projects/hrms/ui/pages/deductions_manager_templ.go`
- **Purpose**: Compiled Go template code generated by templ for deductions_manager.templ.

#### `/home/desktop/especa/beta/projects/hrms/ui/pages/employment_types.templ`
- **Relative Path**: `projects/hrms/ui/pages/employment_types.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for employment_types.templ.

#### `/home/desktop/especa/beta/projects/hrms/ui/pages/employment_types_templ.go`
- **Relative Path**: `projects/hrms/ui/pages/employment_types_templ.go`
- **Purpose**: Compiled Go template code generated by templ for employment_types.templ.

#### `/home/desktop/especa/beta/projects/hrms/ui/pages/financials_overview.templ`
- **Relative Path**: `projects/hrms/ui/pages/financials_overview.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for financials_overview.templ.

#### `/home/desktop/especa/beta/projects/hrms/ui/pages/financials_overview_templ.go`
- **Relative Path**: `projects/hrms/ui/pages/financials_overview_templ.go`
- **Purpose**: Compiled Go template code generated by templ for financials_overview.templ.

#### `/home/desktop/especa/beta/projects/hrms/ui/pages/payment_create.templ`
- **Relative Path**: `projects/hrms/ui/pages/payment_create.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for payment_create.templ.

#### `/home/desktop/especa/beta/projects/hrms/ui/pages/payment_create_templ.go`
- **Relative Path**: `projects/hrms/ui/pages/payment_create_templ.go`
- **Purpose**: Compiled Go template code generated by templ for payment_create.templ.

#### `/home/desktop/especa/beta/projects/hrms/ui/pages/worker_create.templ`
- **Relative Path**: `projects/hrms/ui/pages/worker_create.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for worker_create.templ.

#### `/home/desktop/especa/beta/projects/hrms/ui/pages/worker_create_templ.go`
- **Relative Path**: `projects/hrms/ui/pages/worker_create_templ.go`
- **Purpose**: Compiled Go template code generated by templ for worker_create.templ.

#### `/home/desktop/especa/beta/projects/hrms/ui/pages/worker_detail.templ`
- **Relative Path**: `projects/hrms/ui/pages/worker_detail.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for worker_detail.templ.

#### `/home/desktop/especa/beta/projects/hrms/ui/pages/worker_detail_templ.go`
- **Relative Path**: `projects/hrms/ui/pages/worker_detail_templ.go`
- **Purpose**: Compiled Go template code generated by templ for worker_detail.templ.

#### `/home/desktop/especa/beta/projects/hrms/ui/pages/workers_list.templ`
- **Relative Path**: `projects/hrms/ui/pages/workers_list.templ`
- **Purpose**: Templ presentation template component providing UI rendering and layout for workers_list.templ.

#### `/home/desktop/especa/beta/projects/hrms/ui/pages/workers_list_templ.go`
- **Relative Path**: `projects/hrms/ui/pages/workers_list_templ.go`
- **Purpose**: Compiled Go template code generated by templ for workers_list.templ.

#### `/home/desktop/especa/beta/projects/hrms/ui/static/css/auth.css`
- **Relative Path**: `projects/hrms/ui/static/css/auth.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for auth.css.

#### `/home/desktop/especa/beta/projects/hrms/ui/static/css/crm_shared.css`
- **Relative Path**: `projects/hrms/ui/static/css/crm_shared.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for crm_shared.css.

#### `/home/desktop/especa/beta/projects/hrms/ui/static/css/dashboard.css`
- **Relative Path**: `projects/hrms/ui/static/css/dashboard.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for dashboard.css.

#### `/home/desktop/especa/beta/projects/hrms/ui/static/css/finance.css`
- **Relative Path**: `projects/hrms/ui/static/css/finance.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for finance.css.

#### `/home/desktop/especa/beta/projects/hrms/ui/static/css/globals.css`
- **Relative Path**: `projects/hrms/ui/static/css/globals.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for globals.css.

#### `/home/desktop/especa/beta/projects/hrms/ui/static/css/hrms.css`
- **Relative Path**: `projects/hrms/ui/static/css/hrms.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for hrms.css.

#### `/home/desktop/especa/beta/projects/hrms/ui/static/css/pages/billing.css`
- **Relative Path**: `projects/hrms/ui/static/css/pages/billing.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for billing.css.

#### `/home/desktop/especa/beta/projects/hrms/ui/static/css/pages/calendar.css`
- **Relative Path**: `projects/hrms/ui/static/css/pages/calendar.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for calendar.css.

#### `/home/desktop/especa/beta/projects/hrms/ui/static/css/pages/leads.css`
- **Relative Path**: `projects/hrms/ui/static/css/pages/leads.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for leads.css.

#### `/home/desktop/especa/beta/projects/hrms/ui/static/css/pages/projects.css`
- **Relative Path**: `projects/hrms/ui/static/css/pages/projects.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for projects.css.

#### `/home/desktop/especa/beta/projects/hrms/ui/static/css/pages/sessions.css`
- **Relative Path**: `projects/hrms/ui/static/css/pages/sessions.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for sessions.css.

#### `/home/desktop/especa/beta/projects/hrms/ui/static/css/pages/settings.css`
- **Relative Path**: `projects/hrms/ui/static/css/pages/settings.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for settings.css.

#### `/home/desktop/especa/beta/projects/hrms/ui/static/css/pages/settings_extra.css`
- **Relative Path**: `projects/hrms/ui/static/css/pages/settings_extra.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for settings_extra.css.

#### `/home/desktop/especa/beta/projects/hrms/ui/static/css/style.css`
- **Relative Path**: `projects/hrms/ui/static/css/style.css`
- **Purpose**: Modular CSS stylesheet delivering styling and theme design tokens for style.css.

#### `/home/desktop/especa/beta/projects/hrms/ui/static/embed.go`
- **Relative Path**: `projects/hrms/ui/static/embed.go`
- **Purpose**: Go embed directives packaging static CSS, JavaScript, and assets into the application binary.

#### `/home/desktop/especa/beta/projects/hrms/ui/static/js/auth.js`
- **Relative Path**: `projects/hrms/ui/static/js/auth.js`
- **Purpose**: Client-side JavaScript managing client interactions, DOM operations, or library integrations.

#### `/home/desktop/especa/beta/projects/hrms/ui/static/js/fullcalendar.min.js`
- **Relative Path**: `projects/hrms/ui/static/js/fullcalendar.min.js`
- **Purpose**: Client-side JavaScript managing client interactions, DOM operations, or library integrations.

#### `/home/desktop/especa/beta/projects/hrms/ui/static/js/htmx.min.js`
- **Relative Path**: `projects/hrms/ui/static/js/htmx.min.js`
- **Purpose**: Client-side JavaScript managing client interactions, DOM operations, or library integrations.

#### `/home/desktop/especa/beta/projects/hrms/ui/static/js/kanban.js`
- **Relative Path**: `projects/hrms/ui/static/js/kanban.js`
- **Purpose**: Client-side JavaScript managing client interactions, DOM operations, or library integrations.

#### `/home/desktop/especa/beta/projects/hrms/ui/static/js/sortable.min.js`
- **Relative Path**: `projects/hrms/ui/static/js/sortable.min.js`
- **Purpose**: Client-side JavaScript managing client interactions, DOM operations, or library integrations.


### 12. Centralized Integration Test Suites (`tests`) (35 files)

#### `/home/desktop/especa/beta/tests/cms/auth_api_test.go`
- **Relative Path**: `tests/cms/auth_api_test.go`
- **Purpose**: Automated integration and unit test suite validating functionality and contracts for auth_api_test.go.

#### `/home/desktop/especa/beta/tests/cms/auth_middleware_test.go`
- **Relative Path**: `tests/cms/auth_middleware_test.go`
- **Purpose**: HTTP middleware handling tenant isolation, authentication, session validation, or security headers.

#### `/home/desktop/especa/beta/tests/cms/frontend_auth_test.go`
- **Relative Path**: `tests/cms/frontend_auth_test.go`
- **Purpose**: Automated integration and unit test suite validating functionality and contracts for frontend_auth_test.go.

#### `/home/desktop/especa/beta/tests/cms/headless_api_test.go`
- **Relative Path**: `tests/cms/headless_api_test.go`
- **Purpose**: Automated integration and unit test suite validating functionality and contracts for headless_api_test.go.

#### `/home/desktop/especa/beta/tests/cms/middleware_test.go`
- **Relative Path**: `tests/cms/middleware_test.go`
- **Purpose**: HTTP middleware handling tenant isolation, authentication, session validation, or security headers.

#### `/home/desktop/especa/beta/tests/cms/signup_api_test.go`
- **Relative Path**: `tests/cms/signup_api_test.go`
- **Purpose**: Automated integration and unit test suite validating functionality and contracts for signup_api_test.go.

#### `/home/desktop/especa/beta/tests/cms/test_helpers.go`
- **Relative Path**: `tests/cms/test_helpers.go`
- **Purpose**: Automated integration and unit test suite validating functionality and contracts for test_helpers.go.

#### `/home/desktop/especa/beta/tests/core/auth_service_test.go`
- **Relative Path**: `tests/core/auth_service_test.go`
- **Purpose**: Automated integration and unit test suite validating functionality and contracts for auth_service_test.go.

#### `/home/desktop/especa/beta/tests/core/csrf_test.go`
- **Relative Path**: `tests/core/csrf_test.go`
- **Purpose**: Automated integration and unit test suite validating functionality and contracts for csrf_test.go.

#### `/home/desktop/especa/beta/tests/core/db_test.go`
- **Relative Path**: `tests/core/db_test.go`
- **Purpose**: Automated integration and unit test suite validating functionality and contracts for db_test.go.

#### `/home/desktop/especa/beta/tests/core/migrator_test.go`
- **Relative Path**: `tests/core/migrator_test.go`
- **Purpose**: Automated integration and unit test suite validating functionality and contracts for migrator_test.go.

#### `/home/desktop/especa/beta/tests/core/money_test.go`
- **Relative Path**: `tests/core/money_test.go`
- **Purpose**: Automated integration and unit test suite validating functionality and contracts for money_test.go.

#### `/home/desktop/especa/beta/tests/core/password_test.go`
- **Relative Path**: `tests/core/password_test.go`
- **Purpose**: Automated integration and unit test suite validating functionality and contracts for password_test.go.

#### `/home/desktop/especa/beta/tests/core/repository_test.go`
- **Relative Path**: `tests/core/repository_test.go`
- **Purpose**: Automated integration and unit test suite validating functionality and contracts for repository_test.go.

#### `/home/desktop/especa/beta/tests/core/security_test.go`
- **Relative Path**: `tests/core/security_test.go`
- **Purpose**: Automated integration and unit test suite validating functionality and contracts for security_test.go.

#### `/home/desktop/especa/beta/tests/core/session_revocation_test.go`
- **Relative Path**: `tests/core/session_revocation_test.go`
- **Purpose**: Automated integration and unit test suite validating functionality and contracts for session_revocation_test.go.

#### `/home/desktop/especa/beta/tests/core/session_test.go`
- **Relative Path**: `tests/core/session_test.go`
- **Purpose**: Automated integration and unit test suite validating functionality and contracts for session_test.go.

#### `/home/desktop/especa/beta/tests/core/tenant_isolation_test.go`
- **Relative Path**: `tests/core/tenant_isolation_test.go`
- **Purpose**: Automated integration and unit test suite validating functionality and contracts for tenant_isolation_test.go.

#### `/home/desktop/especa/beta/tests/core/test_helpers.go`
- **Relative Path**: `tests/core/test_helpers.go`
- **Purpose**: Automated integration and unit test suite validating functionality and contracts for test_helpers.go.

#### `/home/desktop/especa/beta/tests/finance/finance_flow_test.go`
- **Relative Path**: `tests/finance/finance_flow_test.go`
- **Purpose**: Automated integration and unit test suite validating functionality and contracts for finance_flow_test.go.

#### `/home/desktop/especa/beta/tests/go.mod`
- **Relative Path**: `tests/go.mod`
- **Purpose**: Go module manifest declaring dependencies and module path for tests.

#### `/home/desktop/especa/beta/tests/go.sum`
- **Relative Path**: `tests/go.sum`
- **Purpose**: Checksum lockfile ensuring reproducible builds and dependency integrity for tests.

#### `/home/desktop/especa/beta/tests/photography_crm/auth_api_test.go`
- **Relative Path**: `tests/photography_crm/auth_api_test.go`
- **Purpose**: Automated integration and unit test suite validating functionality and contracts for auth_api_test.go.

#### `/home/desktop/especa/beta/tests/photography_crm/auth_middleware_test.go`
- **Relative Path**: `tests/photography_crm/auth_middleware_test.go`
- **Purpose**: HTTP middleware handling tenant isolation, authentication, session validation, or security headers.

#### `/home/desktop/especa/beta/tests/photography_crm/bruno/Photography_crm/opencollection.yml`
- **Relative Path**: `tests/photography_crm/bruno/Photography_crm/opencollection.yml`
- **Purpose**: Automated integration and unit test suite validating functionality and contracts for opencollection.yml.

#### `/home/desktop/especa/beta/tests/photography_crm/frontend_auth_test.go`
- **Relative Path**: `tests/photography_crm/frontend_auth_test.go`
- **Purpose**: Automated integration and unit test suite validating functionality and contracts for frontend_auth_test.go.

#### `/home/desktop/especa/beta/tests/photography_crm/middleware_test.go`
- **Relative Path**: `tests/photography_crm/middleware_test.go`
- **Purpose**: HTTP middleware handling tenant isolation, authentication, session validation, or security headers.

#### `/home/desktop/especa/beta/tests/photography_crm/signup_api_test.go`
- **Relative Path**: `tests/photography_crm/signup_api_test.go`
- **Purpose**: Automated integration and unit test suite validating functionality and contracts for signup_api_test.go.

#### `/home/desktop/especa/beta/tests/photography_crm/test_helpers.go`
- **Relative Path**: `tests/photography_crm/test_helpers.go`
- **Purpose**: Automated integration and unit test suite validating functionality and contracts for test_helpers.go.

#### `/home/desktop/especa/beta/tests/superadmin/auth_api_test.go`
- **Relative Path**: `tests/superadmin/auth_api_test.go`
- **Purpose**: Automated integration and unit test suite validating functionality and contracts for auth_api_test.go.

#### `/home/desktop/especa/beta/tests/superadmin/auth_middleware_test.go`
- **Relative Path**: `tests/superadmin/auth_middleware_test.go`
- **Purpose**: HTTP middleware handling tenant isolation, authentication, session validation, or security headers.

#### `/home/desktop/especa/beta/tests/superadmin/frontend_auth_test.go`
- **Relative Path**: `tests/superadmin/frontend_auth_test.go`
- **Purpose**: Automated integration and unit test suite validating functionality and contracts for frontend_auth_test.go.

#### `/home/desktop/especa/beta/tests/superadmin/middleware_test.go`
- **Relative Path**: `tests/superadmin/middleware_test.go`
- **Purpose**: HTTP middleware handling tenant isolation, authentication, session validation, or security headers.

#### `/home/desktop/especa/beta/tests/superadmin/signup_api_test.go`
- **Relative Path**: `tests/superadmin/signup_api_test.go`
- **Purpose**: Automated integration and unit test suite validating functionality and contracts for signup_api_test.go.

#### `/home/desktop/especa/beta/tests/superadmin/test_helpers.go`
- **Relative Path**: `tests/superadmin/test_helpers.go`
- **Purpose**: Automated integration and unit test suite validating functionality and contracts for test_helpers.go.


### 13. API Testing Collections & Bruno (`bruno`) (4 files)

#### `/home/desktop/especa/beta/bruno/API.txt`
- **Relative Path**: `bruno/API.txt`
- **Purpose**: Documentation and usage notes for Bruno API test requests and collections.

#### `/home/desktop/especa/beta/bruno/photography/create roles.yml`
- **Relative Path**: `bruno/photography/create roles.yml`
- **Purpose**: Bruno request definition to create and verify role records in photography CRM.

#### `/home/desktop/especa/beta/bruno/photography/opencollection.yml`
- **Relative Path**: `bruno/photography/opencollection.yml`
- **Purpose**: Bruno collection settings and environment variables configuration.

#### `/home/desktop/especa/beta/bruno/photography/signup.yml`
- **Relative Path**: `bruno/photography/signup.yml`
- **Purpose**: Bruno request definition asserting user signup endpoint functionality.


### 14. CLI Tools & Migration Scripts (`cli`) (5 files)

#### `/home/desktop/especa/beta/cli/finance/migrate/finance_migrate`
- **Relative Path**: `cli/finance/migrate/finance_migrate`
- **Purpose**: Compiled executable binary for finance migration utility.

#### `/home/desktop/especa/beta/cli/finance/migrate/go.mod`
- **Relative Path**: `cli/finance/migrate/go.mod`
- **Purpose**: Go module manifest for finance migration CLI.

#### `/home/desktop/especa/beta/cli/finance/migrate/main.go`
- **Relative Path**: `cli/finance/migrate/main.go`
- **Purpose**: Dedicated database migration runner script applying financial schema migrations.

#### `/home/desktop/especa/beta/cli/seed/go.mod`
- **Relative Path**: `cli/seed/go.mod`
- **Purpose**: Go module manifest for seed CLI.

#### `/home/desktop/especa/beta/cli/seed/main.go`
- **Relative Path**: `cli/seed/main.go`
- **Purpose**: Database seeding CLI tool populating initial roles and administrative credentials.


### 15. Agent Configuration & Skills (`.agents`, `.vscode`) (5 files)

#### `/home/desktop/especa/beta/.agents/rules/coding_standard.md`
- **Relative Path**: `.agents/rules/coding_standard.md`
- **Purpose**: Workspace coding standard enforcing strict simplicity and clear maintainability.

#### `/home/desktop/especa/beta/.agents/rules/rule.md`
- **Relative Path**: `.agents/rules/rule.md`
- **Purpose**: Workspace operational rules governing server lifecycle, port preservation, and test policies.

#### `/home/desktop/especa/beta/.agents/skills/dashboard-designer/SKILL.md`
- **Relative Path**: `.agents/skills/dashboard-designer/SKILL.md`
- **Purpose**: Skill guidance specifying UX/UI layout, visual hierarchy, and dashboard design tokens.

#### `/home/desktop/especa/beta/.vscode/launch.json`
- **Relative Path**: `.vscode/launch.json`
- **Purpose**: VS Code / Cursor debug configurations for Go Delve debugger sessions.

#### `/home/desktop/especa/beta/.vscode/settings.json`
- **Relative Path**: `.vscode/settings.json`
- **Purpose**: Workspace editor settings configuring gopls and formatting rules.
