---
trigger: always_on
---

# AI Coding & Development Standard

**Status:** Adopted\
**Coding philosophy:** Option A --- Strict simplicity\
**Applies to:** All new and modified code across all projects

------------------------------------------------------------------------

## 1. Core principle

Prioritize the simplest implementation that correctly meets the
requirement.

Write code that a beginner developer can read, follow, and maintain.
Do not introduce complexity merely because a pattern is common in larger
software projects.

The goal is not to avoid good engineering. The goal is to use only the
engineering structure and complexity that the feature actually needs.

**Default rule:** If a straightforward solution is sufficient, use it.

## 2. Code organization

Clean separation is encouraged, but it must remain proportional to the
feature.

-   Use separate files when they make responsibilities easier to find
    and understand.
-   Use separate folders when they represent meaningful modules or
    distinct areas of responsibility.
-   Use small, clearly named functions for distinct operations.
-   Keep related code together when splitting it would make the flow
    harder to follow.
-   Do not create a new file, folder, package, interface, or layer for
    every small operation.
-   Avoid both giant files and excessive fragmentation.
-   Follow the existing project structure unless there is a clear reason
    to change it.

A good structure is easy to navigate without requiring the developer to
jump through many files to understand one simple operation.

## 3. General coding rules

Prefer familiar, direct, readable code:

-   Use ordinary control flow statements (`if`, `else`, basic loops).
-   Use simple functions with clear names and focused responsibilities.
-   Use straightforward data structures.
-   Keep control flow visible and easy to trace.
-   Handle errors explicitly and immediately.
-   Use comments to explain non-obvious decisions, not to narrate every
    line.
-   Prefer clear variable names over clever or abbreviated names.

Avoid introducing these unless a specific requirement genuinely calls
for them:

-   Overly generic abstractions, dependency-injection frameworks, or reflection.
-   Complex concurrency, threading, or asynchronous synchronization unless strictly required.
-   Custom frameworks or home-grown generic helper systems.
-   Deep chains of wrappers, callbacks, or indirection.
-   Unnecessary design patterns and speculative generalized architecture.

Do not rewrite simple code into a more abstract form just to make it
look sophisticated.

## 4. Functions and files

-   A function should have one understandable purpose.
-   Split a function when doing so makes the code easier to read, test,
    or reuse.
-   Do not split every few lines into separate functions without a
    reason.
-   Keep functions close to the code and module that use them.
-   Use separate files for genuinely separate responsibilities.
-   Avoid creating single-purpose one-file packages or folders that add
    navigation overhead without providing meaningful separation.

The preferred balance is **organized, not over-engineered**.

## 5. APIs and handlers

For ordinary API operations:

1.  Read the request.
2.  Decode/parse the input.
3.  Perform clear, explicit validation.
4.  Execute the required operation.
5.  Return an appropriate HTTP status and response.

Keep handlers straightforward. Do not add DTO mappers, validator layers,
service layers, repository interfaces, or other indirection for a simple
CRUD endpoint unless the feature has a concrete need for them.

Use existing response helpers and project conventions when they make the
handler clearer. Do not create duplicate helpers for the same purpose.

## 6. CRUD and database access

Follow the project's established database setup, libraries, and query tooling.

-   Write clear queries with meaningful names.
-   Use parameterized queries; never build queries by concatenating
    untrusted input.
-   Keep query files organized by feature or responsibility where
    appropriate.
-   Do not manually edit generated files (e.g., ORM or SQL generators).
-   Do not introduce another database access library or query pattern
    without approval.
-   Do not add repository or service layers around simple queries unless
    they provide a real benefit.

Keep database operations understandable and consistent with the existing
codebase.

## 7. Validation and error handling

Keep ordinary validation explicit and simple.

For example, basic required-field checks should generally use direct
conditions rather than a large validation framework.

However:

-   Do not remove validation that protects data integrity.
-   Do not ignore errors.
-   Do not hide errors behind complicated helper chains.
-   Return appropriate HTTP status codes.
-   Give users clear, safe error messages.
-   Avoid exposing internal errors, database details, secrets, or sensitive
    information in API responses.

If a feature requires complex validation, explain why it is needed and
isolate it in a clearly named function or file.

## 8. Security and multi-tenancy

Simplicity must never weaken security.

Preserve and correctly implement:

-   Authentication and session verification.
-   Authorization and role/permission checks.
-   Tenant isolation on every tenant-owned resource.
-   Password hashing and secure credential handling.
-   Secure cookie settings and session expiration.
-   Parameterized queries / protection against injection.
-   Input validation and data-integrity constraints.
-   Safe file-upload and file-access rules where applicable.

Never assume that a user may access a record just because they know its
ID. Ownership and authorization must always be checked.

Security-sensitive complexity may be necessary. Keep it centralized,
explain its purpose, and avoid spreading it throughout ordinary feature
code.

## 9. Project technology and conventions

Follow the established project stack and conventions:

-   Adhere to the language, framework, routing, templating, and database
    libraries already in place in the repository.
-   Do not replace an established technology or add a new dependency
    without first explaining the reason and receiving approval.
-   Follow existing project layouts and module conventions. Shared code
    belongs in a shared location only when it is genuinely used by multiple
    parts of the project.

## 10. Handling unavoidable complexity

Strict simplicity does not mean pretending every problem is simple.

When advanced logic is genuinely necessary:

1.  Explain the problem in plain English.
2.  Explain why a simpler approach is insufficient.
3.  Keep the complex logic in a small, clearly named function or focused
    file.
4.  Keep ordinary handlers and CRUD code simple.
5.  Document important assumptions and edge cases.
6.  Avoid making unrelated parts of the application depend on the
    complex implementation.

A reviewer should be able to understand what the complex part
is responsible for, even if the internal implementation requires
specialized domain knowledge.

## 11. AI implementation workflow

Before writing code, AI must:

1.  Understand the requested behavior.
2.  Inspect the relevant existing files and conventions.
3.  Identify the smallest set of files that need changes.
4.  Check whether an existing helper or component already solves the
    problem.
5.  Identify security, authorization, and tenant-isolation requirements.

While implementing:

1.  Prefer the simplest correct solution.
2.  Keep changes limited to the requested feature.
3.  Avoid unrelated refactoring.
4.  Avoid unnecessary dependencies and abstractions.
5.  Follow existing naming, routing, database, and response conventions.

After implementing:

1.  Summarize what changed and which files changed.
2.  Explain the implementation in beginner-friendly language.
3.  Explain any unavoidable advanced concepts and why they were needed.
4.  State which tests or checks were run and their results.
5.  Clearly disclose anything not tested or not completed.

Do not claim that code works or tests pass unless that has actually been
verified.

## 12. Before adding complexity

AI must ask itself:

-   Can this be implemented with a simple function?
-   Can an ordinary `if` statement handle this condition?
-   Is an additional file or folder genuinely useful?
-   Is this interface, abstraction, or dependency actually required?
-   Does this layer make the code easier to understand, or merely add
    indirection?
-   Can the complex part be isolated without complicating ordinary code?
-   Am I changing anything outside the requested scope?

If the simple solution is sufficient, use it.

------------------------------------------------------------------------

## Final rule

**Build with clean organization and the least complexity necessary.
Keep ordinary code beginner-readable. Separate files, folders, and
functions when they improve clarity, but do not overdo it. Preserve
security and correctness. Explain and isolate complexity that cannot
reasonably be avoided.**
