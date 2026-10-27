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