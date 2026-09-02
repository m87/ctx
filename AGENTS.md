# AGENTS.md

## Role

Work as a maintainer of `ctx`, a time-tracking application with a Go CLI/server and an optional Angular UI. Make the smallest complete change that solves the task, preserve existing public behavior unless the task explicitly changes it, and leave the repository in a verifiable state.

Treat executable configuration as the source of truth. In particular, prefer `go.mod`, `ui/package.json`, `ui/package-lock.json`, `ui/angular.json`, and `.github/workflows/` over version numbers or commands in prose documentation.

## Commands

Run commands from the repository root unless a command starts with `cd ui`.

### Prerequisites and setup

- Go version: `1.26.3` from `go.mod`.
- UI runtime: Node.js `^22.22.3`, `^24.15.0`, or `>=26.0.0`; npm `>=11.10.0`. The pinned package manager is npm `11.12.1`.
- Install Go modules with `go mod download`.
- Install UI dependencies reproducibly with `cd ui && npm ci`. Do not use `npm install` for routine setup.

### Backend and CLI

```sh
# Format changed Go files.
gofmt -w path/to/changed.go path/to/changed_test.go

# Run all Go tests and build all packages.
go test ./...
go build ./...

# Run a focused package or test while iterating.
go test ./core
go test ./core -run '^TestIntervalSplit$'
```

CI rejects any Go file reported by `gofmt -l .`. There is no separate lint command configured; do not invent one as a required check.

### UI

```sh
cd ui
npm run format:check
npm test -- --watch=false
npm run build -- --configuration production
```

`npm test` uses Vitest through Angular's unit-test builder. Always disable watch mode in non-interactive runs.

### Embedded all-in-one build

Build the UI before preparing embedded assets:

```sh
cd ui
npm ci
npm run build -- --configuration production
cd ..
sh ./scripts/prepare-spa-assets.sh
go build -tags allinone ./...
```

The preparation script copies `ui/dist/ctx/browser` to `server/spa/dist`. Both locations are generated output and must not be edited by hand.

### Full CI-equivalent validation

For cross-cutting or release-sensitive work, reproduce `.github/workflows/go.yml` in this order:

1. `cd ui && npm ci`
2. `npm run format:check`
3. `npm run build -- --configuration production`
4. From the root, `sh ./scripts/prepare-spa-assets.sh`
5. `test -z "$(gofmt -l .)"`
6. `go build -v ./...`
7. `go build -tags allinone -v ./...`
8. `go test -v ./...`

Also run `cd ui && npm test -- --watch=false` when UI code changes; the current CI workflow does not run UI unit tests for you.

## Project structure

- `main.go` is the binary entry point.
- `cmd/` defines Cobra commands, output formats, local/remote dispatch, and CLI integration tests.
- `core/` contains domain models, repository interfaces, business rules, integrity checks/repairs, synchronization, and in-memory test doubles.
- `storage/` implements the repository interfaces with GORM and SQLite. Storage tests use an in-memory database.
- `migration/` migrates data from the legacy `nod` schema into the current model.
- `bootstrap/` opens storage and wires repositories into managers.
- `server/` contains HTTP handlers and route registration. APIs are mounted under `/api` and also at legacy root-level paths for existing clients.
- `server/spa_enabled.go` and `server/spa_disabled.go` select embedded-SPA behavior with the `allinone` build tag.
- `ui/src/api/` contains typed HTTP services plus TanStack Query query/mutation definitions.
- `ui/src/app/` contains standalone Angular application components and state.
- `ui/libs/ui/` contains the local Spartan/Helm UI primitives addressed through TypeScript path aliases.
- `scripts/prepare-spa-assets.sh` stages a production UI bundle for Go embedding.
- `.github/workflows/` is the authoritative CI and release automation.

The normal dependency flow is:

```text
CLI or HTTP handler -> bootstrap -> core manager -> repository interface -> storage implementation
Angular component -> query/mutation -> API service -> /api endpoint
```

Keep business rules in `core/`. Handlers and CLI commands should validate transport-specific input, call a manager, and translate the result; they should not duplicate domain or persistence logic. Keep GORM-specific entities and queries in `storage/`.

## Code style and patterns

### Go

- Let `gofmt` determine formatting and import grouping.
- Return errors instead of logging and continuing. Preserve typed domain errors when callers need to map them to CLI or HTTP behavior.
- Validate required input, propagate repository errors, distinguish a missing entity, and then delegate persistence. Follow this existing shape:

```go
workspace, err := m.WorkspaceRepository.GetById(context.WorkspaceId)
if err != nil {
	return "", err
}
if workspace == nil {
	return "", &WorkspaceNotFoundError{WorkspaceId: context.WorkspaceId}
}
return m.ContextRepository.Save(context)
```

- Preserve the established public model names such as `Id`, `ContextId`, and `WorkspaceId`; do not introduce a compatibility-breaking rename as incidental cleanup. Local variables commonly use idiomatic names such as `contextID`.
- Store and compare persisted timestamps in UTC. Use explicit time zones at user/API boundaries.
- Use `writeJSON` and `writeError` for HTTP responses so error bodies retain the `{code, description}` contract.
- When adding or changing a server resource, consider all affected surfaces: `core`, `storage`, `server`, remote CLI code, and UI API clients.

### TypeScript and Angular

- TypeScript and Angular template checking are strict. Do not bypass them with `any`, unchecked non-null assertions, or disabled diagnostics unless the reason is documented and unavoidable.
- Follow `.editorconfig` and Prettier: two-space indentation, single quotes in TypeScript, 100-character print width, and trailing commas where Prettier emits them.
- Use standalone Angular components and the existing `inject(...)`, signals, and inline template/style patterns.
- Keep raw HTTP and DTO types in `ui/src/api/<resource>/<resource>.service.ts`.
- Keep TanStack Query keys/options in `*.queries.ts`, mutations and cache follow-up behavior in `*.mutations.ts`, and central cache invalidation in `CacheService`.
- Prefer derived state over duplicated component state. When state must be shared across features, follow the existing signal service or NGXS patterns already used by that feature.
- Reuse primitives from `ui/libs/ui/` and existing shared components before creating a new one.

An API query should retain stable, resource-scoped keys and gate requests that lack required identifiers:

```ts
detail(contextId: string) {
  return {
    queryKey: contextQueryKeys.detail(contextId),
    queryFn: () => lastValueFrom(this.contextService.getContext(contextId)),
    enabled: contextId.length > 0,
  };
}
```

## Testing

- Add or update tests for every behavior change. Put tests next to the code as `*_test.go` or `*.spec.ts`.
- Test business rules in `core/` with repository mocks and fixed UTC timestamps.
- Test storage behavior with `CreateTestInMemoryStorage`; never point tests at a developer's `ctx.db`.
- Test handlers with `net/http/httptest`, checking status codes and decoded response bodies.
- Test Cobra wiring and local/remote parity in `cmd/` when CLI behavior changes.
- Use Go subtests for related cases and `testify/require` when later assertions depend on setup succeeding. Existing tests also use `testify/assert` for independent comparisons.
- Use Vitest's `describe`, `it`, and `expect` style for UI unit tests. Cover success, empty/error state, and cache invalidation where relevant.
- For API route changes, verify both `/api/...` and the corresponding legacy root route unless the user explicitly approves a compatibility break.
- Run the narrowest useful test while iterating, then the complete relevant suite before handoff. Run both Go and UI suites for changes to shared API contracts.

## Documentation and compatibility

- Update `README.md` when setup, configuration, commands, modes, or user-visible behavior changes.
- Update `CHANGELOG.md` only when the task includes a release-facing change or explicitly requests a changelog entry; follow its existing format.
- Preserve CLI output formats (`text`, `json`, `yaml`, and `shell`) and remote/local command parity.
- Preserve the `/api` routes and legacy root routes unless removal is explicitly in scope.
- Treat database upgrades as data-safety work. Cover both fresh schema creation and legacy migration paths.

## Git workflow

- Inspect `git status --short` before and after editing. Existing changes belong to the user; do not overwrite, revert, or reformat unrelated files.
- Keep patches scoped to the task. Avoid opportunistic refactors mixed with behavior changes.
- Do not commit, amend, rebase, force-push, tag, publish, or deploy unless the user explicitly asks.
- Do not edit `go.sum` or `ui/package-lock.json` manually; update them only through the relevant package manager when dependency changes are approved.
- In the handoff, state what changed, which validation commands ran, and any failures or checks that could not run.

## Boundaries

### Always

- Read the nearest implementation and tests before changing a pattern.
- Format changed files and run tests proportional to the change.
- Keep domain, storage, HTTP, CLI, and UI contracts synchronized when a shared model changes.
- Use an isolated temporary SQLite database for manual CLI/server checks. The ignored root `ctx.db` may contain developer data.

### Ask first, unless explicitly required by the task

- Add, remove, or upgrade dependencies.
- Change the database schema, migration semantics, or destructive data behavior.
- Break a public CLI flag/output, HTTP endpoint, JSON shape, or legacy route.
- Modify release automation, package publication, container publication, or deployment settings.
- Perform a broad generated-UI replacement or a repository-wide refactor.

### Never

- Read, write, commit, or print secrets, credentials, API keys, personal config files, or real user databases.
- Edit generated/cache/dependency directories: `ui/node_modules/`, `ui/.angular/`, `ui/dist/`, or generated files under `server/spa/dist/`.
- Remove, skip, or weaken a failing test merely to make a check pass.
- Silently discard unrelated work or use destructive Git commands.
- Run a release, push a container, deploy, or push Git changes without explicit authorization.
