---
name: gen-diff-report
description: Generate an HTML diff report from unstaged git changes. Reads git diff, builds the JSON input, runs make gen-diff, and reports the output path.
---

# Generate Diff Report

Inspect unstaged changes, produce the JSON input for `make gen-diff`, run the generator, and report the output path.

## When to Use

User asks to:
- "generate a diff report for my changes"
- "create an HTML report of the unstaged changes"
- "document what I changed"

## Workflow

**1. Read unstaged changes**

```bash
git diff
git status
```

If there are no unstaged changes, stop and tell the user.

**2. Build the JSON input**

Map the diff output to this exact structure:

```json
{
  "title": "<short descriptive title, e.g. 'Fix collision detection'>",
  "explanation": "<2–3 paragraph plain text summary of what changed and why, paragraphs separated by \\n\\n>",
  "files": [
    {
      "path": "<file path from diff header>",
      "hunks": [
        {
          "header": "<hunk header, e.g. '@@ -10,5 +10,12 @@'>",
          "lines": [
            {"kind": "add", "text": "<line without leading +>"},
            {"kind": "del", "text": "<line without leading ->"},
            {"kind": "ctx", "text": "<line without leading space>"}
          ]
        }
      ]
    }
  ]
}
```

Rules:
- `title` — concise, lowercase, describes the intent (used as the output filename)
- `explanation` — plain prose only, no markdown, no bullet points; paragraphs split on `\n\n`
- `kind` — `"add"` for `+` lines, `"del"` for `-` lines, `"ctx"` for unchanged context lines
- `text` — strip the leading `+`, `-`, or space; do NOT HTML-escape (the renderer does it)
- Skip diff metadata lines (`diff --git`, `index`, `---`, `+++`)

**3. Save the JSON**

Write to `output/tmp/input-<slug>.json` where `<slug>` matches the title slug.

**4. Run the generator**

```bash
make gen-diff INPUT=output/tmp/input-<slug>.json
```

**5. Report the output**

Tell the user:
- The generated file path: `output/tmp/<slug>.html`
- How to preview: `make serve` then open `http://localhost:8080/tmp/`
