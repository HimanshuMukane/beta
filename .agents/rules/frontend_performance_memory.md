# Frontend Performance & Memory Management Standard

## 1. Code Splitting & Lazy Loading
- Never bundle all pages and heavy views into a single entry file.
- Use `React.lazy()` with dynamic `import()` for all individual page components, modal-heavy subviews, and complex widgets (e.g. data tables, Kanban boards, chart engines).
- Wrap lazy components in a `Suspense` boundary with a clean, lightweight fallback loader.

## 2. Memory Hygiene & Garbage Collection
- **View-scoped state:** Keep feature-specific and table data in the component's local state so unmounting automatically frees the memory heap. Do not store ephemeral page data in permanent global state.
- **Subscription & listener cleanup:** Every `useEffect` that attaches DOM event listeners, intervals, timeouts, or WebSocket subscriptions must return a cleanup function.
- **Request abortion:** Long-running or heavy fetch requests should use `AbortController` and abort on component unmount to avoid background memory retention.
- **Cache limits:** If using client-side query caching, keep cache durations bounded to prevent unbounded memory growth in long sessions.

