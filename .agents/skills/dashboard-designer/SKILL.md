---
name: dashboard-designer
description: 'Use this skill when designing a data dashboard—choosing KPIs, structuring layout, applying visual hierarchy, or deciding which BI tool to use. Trigger phrases: ''design a dashboard'', ''build a KPI dashboard'', ''what should my dashboard show'', ''help me layout a dashboard'', ''dashboard for monitoring''. Not for building chart code from scratch (use chart-builder), writing SQL queries (use sql-analyst), or designing marketing landing pages.'
version: 1.2.0
author: community
tags:
  - data
  - dashboard
  - visualization
  - business-intelligence
  - ui-ux
  - lucide-icons
license: MIT
keywords:
  - design dashboard
  - KPI dashboard
  - BI tool
  - dashboard layout
  - dashboard
  - designer
  - dashboard designer
  - lucide icons
---

# Dashboard Designer

## Overview
This skill provides a structured approach to dashboard design—from defining the right KPIs and audience through information architecture, visual hierarchy, color coding, design systems, and BI tool selection. A well-designed dashboard surfaces the right information to the right person at the right time. This skill combines data visualization principles with modern UI/UX thinking to create dashboards that drive decisions, not just display data.

## When to Use
- You're building a new dashboard and need to decide what to show and how to arrange it
- An existing dashboard is cluttered, ignored, or confusing to users
- You need to choose between BI tools (Tableau, Power BI, Looker, Grafana)
- You want to design a KPI framework before building the technical implementation
- You're presenting data to an executive audience and need to structure the narrative

## When NOT to Use
- Writing SQL queries to power the dashboard (use sql-analyst skill)
- Building chart code in Python/JavaScript (use chart-builder skill)
- Creating a static report or slide deck (use a presentation skill)
- Designing the underlying data model or warehouse (use a data engineering skill)

## Quick Reference
| Design Principle | Rule |
|-----------------|------|
| Visual hierarchy | Most important KPI top-left; details bottom-right |
| White space | 15–20% of dashboard area; prevents cognitive overload |
| Color coding | Max 4–5 colors; use red/green for status, grey for context |
| Chart density | Max 6–8 charts per screen; prefer fewer, larger charts |
| KPI count | Executive: 3–5 KPIs. Operational: 8–12. Analyst: unlimited |
| Text | Titles as insights ("Revenue Up 12%"), not labels ("Revenue") |
| Refresh rate | Real-time for ops, daily for management, weekly/monthly for strategy |
| Interactivity | Filters at top; drill-down by click; default to most common view |
| Mobile | Stack charts vertically for mobile; limit to top 3 KPIs |
| Accessibility | Colorblind-safe palette; min contrast ratio 4.5:1 |
| Typography | Tabular figures (`tabular-nums`) for numeric data alignment |
| Grid system | 12-column layout; 8px baseline spacing system |
| Icon system | Lucide Icons (24×24 viewBox, stroke-width: 1.5px–2px, `stroke="currentColor"`) |

---

## Design System & Theme Specifications

Derived from modern productivity & workspace dashboards (ClickUp/Linear/Notion-inspired): high-contrast obsidian dark mode, crisp clean light mode, saturated status pill badges, and compact avatar/metadata tags.

### Color Tokens

#### Dark Mode (Deep Obsidian / Pitch-Black Minimal)
- **App Canvas / Page Background**: `#08080a` or `#0b0b0e` (true deep dark tone)
- **Primary Sidebar**: `#0f0f13` (distinct subtle offset with `1px solid rgba(255, 255, 255, 0.06)`)
- **Card / Surface Background**: `#131318` (elevated items: `#181820`, hover: `#1f1f2a`)
- **Borders & Dividers**: `rgba(255, 255, 255, 0.07)` to `rgba(255, 255, 255, 0.12)`
- **Text Hierarchy**:
  - Primary (Headings, Values): `#f4f4f6` (95% brightness)
  - Secondary (Labels, Metadata): `#9898a6`
  - Muted / Placeholder: `#5a5a68`
- **Focus / Accent Ring**: `#6366f1` (Indigo/Violet) or `#3b82f6` (Blue)

#### Light Mode (Clean Warm-Slate / Crisp Enterprise)
- **App Canvas / Page Background**: `#f8f9fb` or `#f3f4f6`
- **Primary Sidebar**: `#ffffff` (bordered by `1px solid #e5e7eb`)
- **Card / Surface Background**: `#ffffff` (hover: `#f9fafb`, active: `#f3f4f6`)
- **Borders & Dividers**: `#e5e7eb` or `#e2e8f0`
- **Text Hierarchy**:
  - Primary (Headings, Values): `#0f172a` (Slate 900)
  - Secondary (Labels, Metadata): `#475569` (Slate 600)
  - Muted / Placeholder: `#94a3b8` (Slate 400)
- **Card Elevation**: `0 1px 3px 0 rgba(0, 0, 0, 0.05), 0 1px 2px -1px rgba(0, 0, 0, 0.05)`

### Saturated Status & Category Pill Badges
Use distinct rounded pills (`border-radius: 9999px` or `rounded-full`, 11px–12px uppercase font, bold tracking, subtle icon + text):
- **To Do / Upcoming**: 
  - Dark: BG `#26262e`, Text `#d4d4d8`, Border `1px dashed rgba(255,255,255,0.2)`
  - Light: BG `#f1f5f9`, Text `#475569`, Border `1px dashed #cbd5e1`
- **In Progress (Active Blue)**:
  - Dark: BG `#1d4ed8`, Text `#ffffff` (or BG `rgba(59, 130, 246, 0.15)`, Text `#60a5fa`)
  - Light: BG `#0284c7`, Text `#ffffff`
- **Bugs / Critical Alert (Crimson)**:
  - Dark: BG `#991b1b`, Text `#fecaca`
  - Light: BG `#dc2626`, Text `#ffffff`
- **Testing / Review (Violet / Purple)**:
  - Dark: BG `#581c87`, Text `#e9d5ff`
  - Light: BG `#7c3aed`, Text `#ffffff`
- **Tag / Category (e.g. `feature`, `crm`)**:
  - Pill badge with a small tag icon: `#0284c7` (blue badge with white text) or light slate with rounded corners.

### UI Element Specs
- **Column / Group Headers**: Section title in uppercase/bold (`12px`, tracking `0.05em`) followed by count badge (e.g., `IN PROGRESS [2]`).
- **Stacked Avatar Badges**: Circular (20×20px or 24×24px) overlapping avatars with initials or user profile photos, bordered by card surface color.
- **Top Command Bar / Search**: Compact centered search input with keyboard shortcut indicator (e.g., `Search ⌘K`) and AI / quick action trigger.

---

## Lucide Icon System Standards

Always use **Lucide Icons** ([lucide.dev](https://lucide.dev)) for visual consistency across all dashboard navigation, KPI metrics, buttons, and status indicators.

### Implementation Guidelines
- **SVG ViewBox**: Always `0 0 24 24`.
- **Stroke Width**: `1.5px` (for secondary elements/large displays) to `2px` (standard UI buttons & badges).
- **Color inheritance**: Always use `stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round"`.
- **Standard Sizes**:
  - `14×14px` (Class: `w-3.5 h-3.5`): Inside compact status pills, inline tags, table rows, and delta indicators.
  - `16×16px` (Class: `w-4 h-4`): Standard button icons, dropdown arrows, filter controls.
  - `20×20px` (Class: `w-5 h-5`): Navigation sidebar links, card header actions, search bar icon.
  - `24×24px` (Class: `w-6 h-6`): Primary Level-1 KPI metric highlight icons (inside a subtle rounded background square).

### Recommended Lucide Icon Mappings
| UI Element | Recommended Lucide Icon | Usage |
|---|---|---|
| **KPI: Growth / Uptrend** | `<TrendingUp />` / `trending-up` | Used next to positive delta values (`+14%`) |
| **KPI: Downtrend / Drop** | `<TrendingDown />` / `trending-down` | Used next to negative delta values (`-8%`) |
| **KPI: Target / Goal** | `<Target />` / `target` | Alongside quota / goal progression indicators |
| **KPI: System / Health** | `<Activity />` / `activity` | Real-time monitoring and heartbeat metrics |
| **Global Date Filter** | `<Calendar />` / `calendar` | Date range picker pill in top header bar |
| **Global Filters** | `<Filter />` / `filter` | Pinned filter trigger with badge count |
| **Search Bar** | `<Search />` / `search` | Left-aligned inside `⌘K` search box |
| **Views / Switcher** | `<LayoutGrid />`, `<Table />`, `<Kanban />`, `<Calendar />` | Switch between Board, Table, List, and Gantt |
| **Status: In Progress** | `<Loader2 />` or `<CircleDot />` | Inside active progress pills |
| **Status: Completed** | `<CheckCircle2 />` / `check-circle-2` | Success indicator |
| **Status: Bug / Alert** | `<AlertCircle />` / `alert-circle` | Error/at-risk highlight |
| **Tags / Labels** | `<Tag />` / `tag` | Next to category tag pills |
| **Assigned / Users** | `<Users />` / `<User />` | Assignee filters or member summaries |
| **Settings / More** | `<MoreHorizontal />` / `<Settings />` | Card action menus and configuration |

---

## Instructions

### Step 1: Define the Audience and Purpose
Answer these questions before designing anything:
1. **Who uses this dashboard?** Executive, team manager, analyst, engineer, customer?
2. **What decision does it inform?** ("Should we increase marketing spend in the West region?")
3. **How often will they look at it?** Real-time (ops), daily, weekly, monthly?
4. **What action can they take?** If no action → reconsider whether a dashboard is the right format.

**Dashboard archetypes:**
| Type | Audience | Refresh | KPI Count | Example |
|------|----------|---------|-----------|---------|
| Strategic | C-suite, VPs | Weekly/Monthly | 3–5 | Company health scorecard |
| Operational | Team leads, managers | Daily | 8–15 | Sales pipeline dashboard |
| Analytical | Data analysts | On-demand | Unlimited | Cohort analysis explorer |
| Real-time monitoring | Engineers, support | Live | 10–20 | System uptime & error rates |

### Step 2: Select the Right KPIs
A KPI must be:
- **Actionable**: the viewer can change it based on what they see
- **Timely**: available at the frequency the dashboard refreshes
- **Comparable**: has a target, benchmark, or prior period to compare against
- **Unambiguous**: has one clear definition shared by all stakeholders

**KPI selection process:**
1. List 20 metrics the audience cares about
2. Filter to metrics the audience can act on (eliminate vanity metrics)
3. Group remaining into 3–5 themes (e.g., Revenue, Customers, Operations, Quality)
4. Pick 1–3 primary metrics per theme
5. Define calculation, data source, and owner for each

**KPI tiers:**
- **Level 1 (headline)**: 3–5 KPIs shown large at the top (Total Revenue, MAU, NPS) paired with Lucide micro-icons, sparklines, or delta pills.
- **Level 2 (supporting)**: 6–10 charts providing context (by region, by product, trend).
- **Level 3 (detail)**: Tables, board lists, or drill-downs for investigation.

### Step 3: Information Architecture (Layout Patterns)

**F-pattern layout** (most common): Users scan top-left → top-right → down the left.
```
┌─────────────────────────────────────────────────────┐
│  KPI 1        KPI 2        KPI 3        KPI 4       │  ← Headline KPIs
├───────────────────────┬─────────────────────────────┤
│                       │                             │
│  Primary chart        │  Secondary chart            │  ← Primary context
│  (e.g. revenue trend) │  (e.g. by region)           │
│                       │                             │
├───────────────────────┴─────────────────────────────┤
│  Detail chart 1    │  Detail chart 2  │  Table      │  ← Supporting detail
└────────────────────┴──────────────────┴─────────────┘
```

**Z-pattern layout** (executive dashboards): Fewer elements, high-level narrative flow.

**Grid layout** (operational dashboards): Equal-sized tiles arranged in rows based on a 12-column grid.

### Step 4: Visual Hierarchy Principles
1. **Size**: Larger = more important. Make the most critical KPI card the biggest element.
2. **Position**: Top-left corner gets most attention. Put the hero metric there.
3. **Color**: Use deep obsidian surfaces in dark mode; stark white cards on slate in light mode. Use saturated color pills only for state highlights.
4. **Typography**:
   - Hero metric values: 24–32pt bold with `font-variant-numeric: tabular-nums`.
   - Labels / descriptions: 10–12pt regular or muted.
   - Section / card titles: 14–16pt semi-bold.
5. **White space**: Empty space is not wasted—it guides the eye and reduces fatigue (15–20% allocation).

### Step 5: Color Coding System
Establish a consistent color language:
- 🟢 **Green**: good, on track, above target
- 🔴 **Red**: bad, alert, below target
- 🟡 **Yellow/Orange**: caution, approaching threshold
- 🔵 **Blue**: primary data, neutral information
- ⬜ **Grey**: secondary data, prior period, context


### Step 8: UI/UX Rules and Component Standards

#### KPI Card Specification
- **Card container**: Subtle border (`1px solid var(--border)`), 8px–12px radius, dark/light surface tokens.
- **Card header**: Metric title (12–14px, medium weight, muted), paired with a relevant Lucide icon (e.g. `TrendingUp` or `DollarSign`).
- **Hero value**: Large bold typography (28–32px), always using tabular figures (`tabular-nums`).
- **Context indicator / Delta pill**:
  - Saturated badge (e.g., `+12.4% vs LW` with `TrendingUp` icon).
- **Micro-visualization**: Embedded 7-day or 30-day sparkline or horizontal target bar underneath the value.

#### Grid & Spatial Hierarchy
- **12-column grid layout**:
  - 4 cards across: `col-span-3` each.
  - 3 cards across: `col-span-4` each.
  - Split view: `col-span-8` (primary chart) + `col-span-4` (secondary chart or feed).
- **8px spacing baseline**: Use uniform increments (8px, 16px, 24px) for padding, grid gaps, and margins.

#### Interactivity & Filter UX
- **Sticky / pinned filter bar**: Date ranges (`Calendar` icon) and critical dimensions stay fixed during page scroll.
- **Active filter feedback**: Display filter pills with clear badge counts (`Filters (3)`) and a one-click "Reset All" option.
- **Chart drill-downs & tooltips**:
  - Show hover states with `cursor: pointer` on interactive bars/points.
  - Tooltips should render contextual values, dates, and deltas—not raw chart coordinates.

#### Edge States & Reliability
- **Loading states**: Use layout-matching skeleton screens; avoid page-blocking full-screen spinners.
- **Error containment**: Isolate data load failures to the specific card affected with a localized retry button.
- **Empty states**: Never display bare empty boxes; show a lightweight illustration with a clear message and action ("No data for selected date range. Try widening filters.").

### Step 9: Dashboard Quality Checklist
Before publishing, verify:
- [ ] Every chart has a title stating the insight (not just the metric name)
- [ ] Every KPI has a comparison value (vs. target, vs. last period)
- [ ] Tabular numbers (`tabular-nums`) are active for numerical alignment
- [ ] Lucide icons used consistently (24×24 viewBox, stroke-width 1.5–2px, `stroke="currentColor"`)
- [ ] Colors match dark/light theme tokens; saturated badges reserved for status pills
- [ ] Date range is clearly labeled and user-controlled with Lucide `Calendar` trigger
- [ ] Filters are sticky or at the top and default to the most common view
- [ ] Skeleton loaders and empty states are styled and tested
- [ ] All numbers are formatted with appropriate precision (revenue: $1.2M not $1,234,567.89)
- [ ] Data source and refresh date are visible in a footer
- [ ] Mobile layout tested (stacked cards, priority KPIs visible without scroll)

## Examples

### Example 1: Design a Sales Performance Dashboard
**Input:** "I need a dashboard for our VP of Sales to review weekly. We care about revenue, pipeline health, and team performance."

**Audience analysis:** VP of Sales — executive level — weekly review — strategic decisions about team resourcing and focus.

**KPI selection (Level 1 — headline):**
1. **Bookings (ARR)** — new closed revenue this week vs target + Lucide `TrendingUp` sparkline
2. **Pipeline Coverage** — pipeline value / quota (target: 3× coverage) + gauge bar
3. **Win Rate** — % of deals closed-won (target: 25%) + delta pill vs last quarter
4. **Deals at Risk** — open deals past expected close date + Lucide `AlertCircle` alert badge

**Layout design:**
```
┌────────────────────────────────────────────────────────────────┐
│ FILTER: [📅 Week ▼]  [🌍 Region ▼]  [👤 Rep ▼]   Data as of Mon│
├──────────────┬──────────────┬──────────────┬──────────────────┤
│ Bookings     │ Pipeline     │ Win Rate     │ Deals at Risk    │
│ $1.2M        │ 3.4× quota   │ 28%          │ 12 deals         │
│ ▲ 15% vs LW  │ ✅ on track  │ ▲ 3pt vs LW  │ 🔴 $340K at risk │
│ [ ▂▃▅▆▇ spark]│ [████░░ prog]│ [▲ badge]    │ [warning pill]   │
├──────────────┴──────────────┴──────────────┴──────────────────┤
│ Bookings Trend (12 weeks) ← primary chart, 60% width          │
│                                                                │
├─────────────────────────────┬──────────────────────────────────┤
│ Revenue by Rep (bar chart)  │ Pipeline Stage Funnel            │
├─────────────────────────────┴──────────────────────────────────┤
│ Deal Table: Name | Stage | Value | Expected Close | Days Late  │
└────────────────────────────────────────────────────────────────┘
```

**Color logic:**
- Green: metric above target
- Red: metric below target or "at risk" items
- Blue: trend line (primary data)
- Grey: prior period comparison line

**Tool recommendation:** Salesforce Einstein Analytics (if using Salesforce CRM) or Tableau with Salesforce connector. For smaller teams, HubSpot built-in reporting or Metabase connected to CRM database.

---

### Example 2: Design an Engineering Metrics Dashboard
**Input:** "Engineering team wants a dashboard to track system health, deployment frequency, and incident response."

**Audience:** Engineering team leads, on-call engineers — daily + real-time use.

**KPI selection (DORA Metrics framework):**
1. **Deployment Frequency** — deploys/week (target: daily) + Lucide `Rocket` icon
2. **Lead Time for Changes** — PR merge → production (target: < 1 day) + Lucide `Clock` icon
3. **Change Failure Rate** — % of deploys causing incidents (target: < 5%) + Lucide `AlertTriangle` icon
4. **MTTR (Mean Time to Recovery)** — avg incident duration (target: < 1 hour) + Lucide `Activity` icon

**Supporting charts:**
- Error rate by service (last 24h line chart)
- P95 API latency by endpoint (heatmap)
- Open incidents table with severity and SLA status
- Recent deployments log

**Tool recommendation:** Grafana (connects directly to Prometheus, DataDog, or CloudWatch; native alerting; real-time refresh).

**Color coding:** Red = SLA breach, Yellow = degraded, Green = healthy, Grey = no data.

## Best Practices
- Design for the question, not the data—start with "what decision does this enable?" not "what data do I have?"
- One dashboard, one audience—don't try to serve executives and analysts with the same view
- KPIs need targets: a number without a benchmark is meaningless
- "Alert" color (red) loses impact if overused; reserve it for genuine problems
- Show the trend, not just the snapshot—a current value without history lacks context
- Label data directly on charts when possible rather than relying on legend lookups
- Build the mobile version after the desktop version, not as an afterthought
- Isolate failures: one failing widget must never crash the entire dashboard interface
- Use Lucide SVG icons consistently across dark and light modes with semantic stroke coloring

## Common Mistakes
- Too many KPIs: 20+ metrics creates analysis paralysis, not insight
- Vanity metrics: page views, social followers, total registered users (not actionable)
- Missing comparison context: showing $1.2M revenue without target or prior period
- Inconsistent date filters: one chart shows "last 30 days" while another shows "this month"
- Using pie charts for more than 5 segments—switch to a ranked bar chart
- Auto-refreshing a dashboard that users need to study calmly—adds cognitive load
- Forgetting mobile: 40%+ of dashboard views happen on tablets or phones
- Neglecting empty/error states: displaying raw JSON errors or empty white voids when queries fail
- Mixing disparate icon packs with differing stroke-weights, visual grids, or corner radii

## Tips & Tricks
- Run a "5-second test": show the dashboard to someone for 5 seconds, then ask "what's the most important thing?"—if they can't answer, the hierarchy isn't clear
- Use grey as your default color for all secondary data; use color only to highlight what matters
- Always include a "last updated" timestamp so users know if data is stale
- Dashboard titles should answer the question the dashboard addresses: "Sales Performance — Week of Jan 15" not "Sales Dashboard"
- For executive dashboards, start with a written "dashboard narrative" document before building—executives often want a briefing, not a tool to explore
