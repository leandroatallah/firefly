# NOTES — 060 HTML Report Generator

## Investigation

- Repo already uses `scripts/` for Go dev tooling (`kanban.go`, `next-id.go`, `sdd-dashboard.go`). These are run via `go run scripts/<name>.go` and are not part of the engine/kit/game three-layer architecture, so the constitution's dependency rules do not constrain `scripts/gen/`.
- Existing `scripts/` files are flat `.go` files in the `scripts/` package. For multi-file generators with templates + helpers + cmds, a subdirectory tree is cleaner. Need to verify the repo's module layout supports `scripts/gen/...` as importable subpackages — if `scripts/` is excluded from the module (e.g., via build tag or separate go.mod), fall back to `internal/tools/gen/` with the same structure. Implementer should confirm during Green phase.

## Design Rationale

- **`html/template` over `text/template`**: gives automatic context-aware escaping. Critical for AC-3 (diff line text safety) and to keep generators dependency-free (AC-6).
- **Inline CSS in `base.html`**: simplest way to honor AC-6 (no CDN, no separate stylesheet to track). The 240px sidebar / 960px content layout is small enough (<1 KB) that inlining has no real cost.
- **`render.Renderer` clones base per call**: prevents mutating the parsed base template when overriding blocks; allows parallel rendering without locks.
- **`Slug` field provided by caller, not derived**: domain entries may have ambiguous names ("i18n" → "i18n" vs "i-18-n"). Letting the caller specify the slug keeps URLs stable across regenerations.
- **JSON stdin for cmds**: cmds stay dumb wrappers; the rendering libraries are unit-testable without subprocess plumbing.

## Risks

- **Module boundary uncertainty**: if `scripts/gen/render` cannot be imported by `scripts/gen/diff` (Go module rules around `scripts/`), the implementer must move the tree to `internal/tools/gen/`. This is a layout question, not a behavior change — tests stay identical.
- **`html/template` block override gotcha**: defining a `{{block}}` in the base sets a default, and a child template parsed into the same set overrides it. The renderer must `Clone()` the base before parsing each child, otherwise consecutive `Render` calls will leak block definitions from one report into another. Covered indirectly by T-R3 if tests are run in sequence; consider adding a sequential render test if regressions appear.
- **HTML escaping of explanation prose**: AC-3 says "prose explanation" — current spec treats it as plain text split into `<p>` paragraphs. If markdown is desired later, swap to a parser; for now plain-text + paragraph splitting keeps the no-dependency rule intact.

## Out of Scope

- JavaScript-driven sidebar collapse/expand (story explicitly excludes JS).
- Cross-page search.
- Auto-discovery of which domain entries exist (input is explicit `Index{Entries:...}`).
- Pretty diff colors beyond add/del/ctx (e.g., word-level intra-line diff).
- Watch mode / incremental rebuild.

## Open Questions (defer to implementer)

1. Final module path for the package tree (`scripts/gen/...` vs `internal/tools/gen/...`).
2. Whether `cmd/diff` and `cmd/domain-docs` should also accept a `-o <path>` flag for explicit output location, in addition to defaulting to `output/tmp/` and `output/docs/`. Not required by ACs; leave for a follow-up if needed.
