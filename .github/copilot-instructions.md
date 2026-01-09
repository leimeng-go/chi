# Copilot Instructions for this Repository

Overview: This repo is the core `chi` HTTP router (Go). It provides a small, fast radix-trie router, middleware composition, and route introspection. It’s stdlib-only (net/http) with zero external deps.

Core Types & Files:
  - Router API and entrypoint: [chi.go](chi.go), [mux.go](mux.go)
  - Radix routing tree and walking: [tree.go](tree.go)
  - Request routing context: [context.go](context.go)
  - Middleware chaining: [chain.go](chain.go)
  - Go 1.23 `Request.Pattern` support: [pattern.go](pattern.go) and fallback [pattern_fallback.go](pattern_fallback.go)
  - Optional/stdlib middlewares: [middleware/](middleware)
  - Examples of real usage: [_examples/]( _examples )
  - Experimental generic radix tree (non-HTTP): [rtree.go](rtree.go)

Big Picture Architecture:
  - `Mux` composes a middleware chain with a radix routing tree (`node`) to route requests. See [mux.go](mux.go) and [tree.go](tree.go).
  - URL params and matched pattern flow via `Context` stored at key `RouteCtxKey`. Handlers access with `chi.URLParam()` or standard `r.PathValue()` (Go ≥1.22). See [context.go](context.go) and [mux.go](mux.go#L176-L196).
  - Route discovery/docs: routers expose `Routes()`; tree walking via `Walk()` in [tree.go](tree.go#L348-L416). External tools like `docgen` consume this.

Routing Patterns (must start with `/`):
  - Named params: `/users/{userID}` → `chi.URLParam(r, "userID")`
  - Regexp params: `/date/{yyyy:\d\d\d\d}` (RE2, `/` never matches)
  - Catch-all: `/files/*` (matches the rest of the path)
  - Wildcard must be last. Adjacent params/wildcards are validated in [tree.go](tree.go#L463-L548).

Composition & Mounting:
  - Define middlewares before routes; calling `Use()` after routes panic-protects (see [mux.go](mux.go#L54-L74) and [mux.go](mux.go#L89-L100)).
  - Inline middlewares with `With()`, group with `Group()`, nest routers with `Route()`/`Mount()`. `Mount()` wires a wildcard and passes control to subrouter; don’t mount two handlers on the exact same pattern (it panics). See [mux.go](mux.go#L129-L220).
  - Customize `NotFound()` and `MethodNotAllowed()`; values cascade to subrouters unless overridden. See [mux.go](mux.go#L101-L128).

Stdlib Integrations:
  - Go ≥1.22: `mux.routeHTTP` populates `http.Request` path values via `r.SetPathValue()`; use `r.PathValue("userID")` in handlers. See [mux.go](mux.go#L176-L196).
  - Go ≥1.23: when available, `r.Pattern` is set to the matched route (see [pattern.go](pattern.go) and [pattern_fallback.go](pattern_fallback.go)).

Developer Workflows:
  - Tests (race + verbose): `make test` (runs router and middleware tests). Targets: `make test-router` and `make test-middleware` (see [Makefile](Makefile)).
  - Run examples:
    - Hello world: `go run _examples/hello-world/main.go`
    - REST demo: `go run _examples/rest/main.go` (routes printed in `_examples/rest/routes.{json,md}`)
  - Module: [go.mod](go.mod) sets `module github.com/go-chi/chi/v5` and `go 1.22`. Keep compatibility with the most recent four Go major versions.

Conventions & Gotchas:
  - Always register all `Use()` middlewares before adding any routes to a given `Mux`; the first route seals the chain (`updateRouteHandler`).
  - Patterns are validated: wildcard must be last; duplicate param keys in a pattern panic (see [tree.go](tree.go#L550-L590)).
  - Custom HTTP methods supported via `RegisterMethod("PURGE")` etc. (see [tree.go](tree.go#L60-L101)). Then use `Router.Method()`/`MethodFunc()`.
  - To compute the final matched pattern after a request, call `RouteContext(r.Context()).RoutePattern()` after `next.ServeHTTP` in middleware (see doc in [context.go](context.go#L72-L120)).

Middlewares:
  - Provided middlewares are plain `net/http` compatible (e.g., `Logger`, `Recoverer`, `RequestID`, `Throttle`, etc.). See [middleware/](middleware) and top-level list in [README.md](README.md).
  - Compose endpoint-specific middleware with `r.With(mw...).Get(...)` or route groups with `r.Group(...)`.

Non-HTTP Radix Tree (optional):
  - [rtree.go](rtree.go) implements a separate, generic radix tree (`RTree`) with Chinese comments; it’s not wired into `Mux`. Use independently or in tests.

Pointers for AI Agents:
  - When adding routes/middleware, prefer touching [mux.go](mux.go) APIs; avoid modifying [tree.go](tree.go) unless changing routing semantics.
  - For path param or pattern behavior, add tests in `*_test.go` next to [tree.go](tree.go) or create new ones mirroring existing style (see `pattern_test.go`, `tree_test.go`).
  - Keep stdlib-only policy: do not add external dependencies without explicit approval.

Useful References:
  - Project overview and examples: [README.md](README.md), [_examples/]( _examples )
  - Contribution workflow: [CONTRIBUTING.md](CONTRIBUTING.md)