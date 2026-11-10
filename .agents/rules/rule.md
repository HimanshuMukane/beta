---
trigger: always_on
---

# Execution Rules

## Server and Port Management
- **Never start or run the application server / backend daemon** (e.g., `go run`, `npm run dev`, `python app.py`, etc.) unless explicitly instructed by the user.
- **Never launch browser subagents or automated browser testing** that requires running the dev server.
- **Always keep ports free**: Do not leave processes listening on localhost ports (e.g., port 3000, 8080).
- Only run non-server development build/compilation checks (such as `templ generate` or `go build`) when validating code syntax or updating generated templates.

- no comments

- No migrations, No models, use DB querries directly

## Syntax & Naming Rules
- **Explicit Variable Declarations Always**: Use explicit typed variable declarations (`var <name> <Type> = ...` or `var <name> <Type>; <name> = ...`) instead of short variable declarations (`:=`).
- **Explicit Loops Always**: For slices and collections, use explicit index-based loops with pre-declared iteration variables (e.g. `var index int; for index = 0; index < len(items); index++`). For map iterations, pre-declare the key and value variables. Avoid `[]rune(str)` for iteration or indexing.
- **Proper Semantic Naming Always**: Use full, domain-accurate semantic variable and parameter names (e.g., `responseWriter`, `request`, `database`, `tenantID`, `scanErr`, `queryErr`, `insertErr`, `requestPayload`), never single-letter or abbreviated names (such as `w`, `r`, `db`, `err`, `tid`, `req`).
