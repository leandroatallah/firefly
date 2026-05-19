# HTML Report Generators

The `scripts/gen/` toolkit provides two generators for producing styled HTML reports from structured Go data:

1. **Diff Report Generator** — code diffs with explanations
2. **Domain Docs Generator** — navigable domain documentation site

Both generators use a shared base template (`scripts/gen/base.html`) with inlined CSS, producing offline-friendly HTML that works without a server.

## Quick Start

### Build the Commands

```bash
go build -o scripts/gen/cmd/diff/diff ./scripts/gen/cmd/diff
go build -o scripts/gen/cmd/domain-docs/domain-docs ./scripts/gen/cmd/domain-docs
```

### Diff Report

Generate an HTML page showing code changes with explanation:

```bash
cat << 'EOF' | ./scripts/gen/cmd/diff/diff
{
  "title": "Add Physics Collision",
  "explanation": "Implements HasCollision check for space bodies. Validates body position against obstacle bounds.\n\nUses fixed-point comparison for precision.",
  "files": [
    {
      "path": "internal/engine/space/body.go",
      "hunks": [
        {
          "header": "@@ -10,5 +10,12 @@",
          "lines": [
            {"kind": "ctx", "text": "func (b *Body) Position() fp16.Vec {"},
            {"kind": "add", "text": "func (b *Body) HasCollision(obs []Obstacle) bool {"},
            {"kind": "add", "text": "  for _, o := range obs {"},
            {"kind": "add", "text": "    if b.pos.X >= o.Min.X && b.pos.X <= o.Max.X {"},
            {"kind": "add", "text": "      return true"},
            {"kind": "add", "text": "    }"},
            {"kind": "add", "text": "  }"},
            {"kind": "add", "text": "  return false"},
            {"kind": "add", "text": "}"},
            {"kind": "ctx", "text": "func (b *Body) Update() {"}
          ]
        }
      ]
    }
  ]
}
EOF
```

Output: `output/tmp/add-physics-collision.html`

### Domain Docs Generator

Generate a navigable documentation site:

```bash
cat << 'EOF' | ./scripts/gen/cmd/domain-docs/domain-docs
{
  "title": "Firefly Engine Domain Model",
  "entries": [
    {
      "name": "Physics Engine",
      "slug": "physics",
      "description": "Handles collision detection and body movement using fixed-point arithmetic.",
      "keytypes": ["Body", "Space", "Obstacle"]
    },
    {
      "name": "Scene System",
      "slug": "scene",
      "description": "Manages actor lifecycle and state transitions across game phases.",
      "keytypes": ["Scene", "Actor", "PhaseState"]
    }
  ]
}
EOF
```

Output:
- `output/docs/index.html` — navigable index page with sidebar
- `output/docs/physics.html` — Physics Engine page
- `output/docs/scene.html` — Scene System page

Open `output/docs/index.html` in a browser (no server required).

## Input Schemas

### Diff Report

```go
type Report struct {
  Title       string     // Page title (slugified for output filename)
  Explanation string     // Plain text with \n\n for paragraph breaks
  Files       []FileDiff
}

type FileDiff struct {
  Path  string
  Hunks []Hunk
}

type Hunk struct {
  Header string   // e.g. "@@ -10,5 +10,12 @@"
  Lines  []Line
}

type Line struct {
  Kind string // "add" | "del" | "ctx"
  Text string // Raw line text (HTML auto-escaped)
}
```

**Key points:**
- `Title` must not be empty; non-alphanumeric chars are converted to `-` in filename
- `Explanation` paragraphs are split on blank lines (`\n\n`)
- `Line.Text` is HTML-escaped automatically; no need to pre-escape
- `Kind` values: `"add"` (green), `"del"` (red), `"ctx"` (context, unchanged)

### Domain Docs Index

```go
type Index struct {
  Title   string
  Entries []DomainEntry
}

type DomainEntry struct {
  Name        string   // Display name
  Slug        string   // URL-safe slug (a-z, 0-9, -)
  Description string   // Short summary (rendered as plain text)
  KeyTypes    []string // List of important type names
}
```

**Key points:**
- `Slug` appears in filenames and sidebar links; should be lowercase and URL-safe
- `Description` is not HTML-escaped (future markdown support possible; keep plain for now)
- `KeyTypes` renders in `<code>` tags in a "Key Types" section
- Each entry gets its own page at `output/docs/<slug>.html`

## Usage Patterns

### From Go Code

Use the library packages directly instead of shelling out:

```go
import (
  "os"
  "github.com/boilerplate/ebiten-template/scripts/gen/baseasset"
  "github.com/boilerplate/ebiten-template/scripts/gen/diff"
  "github.com/boilerplate/ebiten-template/scripts/gen/render"
)

func main() {
  r, _ := render.New(baseasset.HTML)
  report := diff.Report{
    Title:       "Fix Physics Bug",
    Explanation: "Corrects collision edge case.",
    Files: []diff.FileDiff{
      // ... populate from git diff, AST analysis, etc.
    },
  }
  
  f, _ := os.Create("output/tmp/fix-physics-bug.html")
  defer f.Close()
  diff.Render(f, r, report)
}
```

### From External Tools

Agents and scripts can generate JSON and pipe to the commands:

```bash
# From a bash script
my-diff-generator | ./scripts/gen/cmd/diff/diff

# From Python
python3 gen_domain_docs.py | go run scripts/gen/cmd/domain-docs
```

### Output Directories

- **Temporary reports** (diffs, plans): `output/tmp/<slug>.html`
  - One file per report
  - No sidebar navigation
  - Safe to delete/regenerate

- **Permanent docs** (domain, architecture): `output/docs/`
  - `index.html` — main entry point with sidebar
  - `<slug>.html` — one page per domain entry
  - Sidebar links all entries
  - Pages are relative-linked (works offline)

## Customization

### Styling

Edit `scripts/gen/base.html` `<style>` section to customize:
- Layout grid (sidebar width, content max-width)
- Colors, fonts, spacing
- Code block styling

Changes apply to all future renders.

### Adding New Generator Types

To create a new generator (e.g., architecture diagrams):

1. Create package under `scripts/gen/<newtype>/`
2. Define data model and `Render(io.Writer, *render.Renderer, Data) error`
3. Create cmd under `scripts/gen/cmd/<newtype>/main.go` that reads JSON stdin
4. Add tests under `scripts/gen/<newtype>/<newtype>_test.go`

All generators share the same base template and render helper.

## Troubleshooting

**File not created / "no such file or directory"**
- Generator creates `output/` directories automatically
- Ensure repo root is the current working directory when running commands
- Check file permissions on parent directory

**Sidebar links broken**
- Domain docs slug must match filename (case-sensitive)
- Use only `[a-z0-9-]` in slugs
- Relative links work only when opening `output/docs/index.html` directly (not from a subdir)

**HTML escaped incorrectly**
- Diff line text is auto-escaped; do not pre-escape `<` to `&lt;`
- Domain doc descriptions are plain text; use unicode or HTML entities if needed

**External CSS/JS not loading**
- By design: all CSS is inlined in `base.html`
- No `<script>` tags or external resources allowed
- Reports work offline in any browser
