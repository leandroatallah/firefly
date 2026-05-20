# User Story 060 — HTML Report Generator

## Story

**As a** developer working on the Firefly engine,  
**I want** a Go-based HTML report generator with a shared base template,  
**So that** I can produce human-readable outputs (code diffs with explanations, domain documentation, implementation plans) that are easy to navigate and visually consistent.

## Context

Agents and developers need to produce structured, readable outputs beyond plain text. These outputs fall into two categories:

- **Permanent docs** — domain documentation, architecture overviews. Navigable via sidebar.
- **Temporary reports** — code diff explanations, visual implementation plans. Standalone, no sidebar nav required.

All outputs share the same visual identity (CSS, layout) but are generated independently by small Go programs under `scripts/gen/`.

## Acceptance Criteria

**AC-1: Shared base template**  
Given a base template exists at `scripts/gen/base.html`,  
when any generator renders a page,  
then the output includes the shared CSS, header, and footer.

**AC-2: Named template blocks**  
Given the base template defines `{{block "title"}}`, `{{block "content"}}`, and `{{block "sidebar"}}` blocks,  
when a generator provides only `content` and `title`,  
then the sidebar block renders empty (no broken layout).

**AC-3: Diff report generator**  
Given a structured diff input (file path, hunks, explanation text),  
when `scripts/gen/diff.go` is run,  
then it renders an HTML page with syntax-highlighted diff hunks and prose explanation, written to `output/tmp/`.

**AC-4: Domain documentation generator**  
Given a list of domain entries (name, description, key types),  
when `scripts/gen/domain-docs.go` is run,  
then it renders a navigable HTML page with a sidebar index, written to `output/docs/`.

**AC-5: Output directory separation**  
Given the generator runs,  
when the output type is temporary (diff, plan),  
then the file is written under `output/tmp/`;  
when the output type is permanent (docs),  
then the file is written under `output/docs/`.

**AC-6: No external dependencies**  
Given the generated HTML files,  
when opened in a browser without a server,  
then all styles and layout render correctly (CSS inlined or co-located, no CDN links).

**AC-7: Standalone pages work without a server**  
Given relative links between permanent doc pages,  
when navigating the sidebar,  
then links resolve correctly from a flat `output/docs/` directory.

## Out of Scope

- JavaScript interactivity.
- Automatic sidebar regeneration across all pages (manual or script-driven is acceptable).
- CI integration or automatic report publishing.

## Domain Language

- **Report** — a single generated HTML file.
- **Base template** — the shared `base.html` defining layout and styles.
- **Diff report** — a temporary report showing code changes with explanation.
- **Domain doc** — a permanent report documenting a bounded context or package.
- **Generator** — a Go program under `scripts/gen/` that renders a specific report type.
