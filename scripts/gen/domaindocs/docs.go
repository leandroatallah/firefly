// Package domaindocs defines the data model for the domain documentation
// site (index + per-entry pages) and renders them through the shared
// render.Renderer.
package domaindocs

import (
	"io"

	"github.com/boilerplate/ebiten-template/scripts/gen/render"
)

// DomainEntry is a single domain page entry.
type DomainEntry struct {
	Name        string
	Slug        string
	Description string
	KeyTypes    []string
}

// Index is the full set of domain entries plus a page title.
type Index struct {
	Title   string
	Entries []DomainEntry
}

const indexTmpl = `{{define "title"}}{{.Idx.Title}}{{end}}
{{define "sidebar"}}
<nav><ul>
{{range .Idx.Entries}}<li><a href="{{.Slug}}.html">{{.Name}}</a></li>
{{end}}</ul></nav>
{{end}}
{{define "content"}}<h2>Index</h2><ul>
{{range .Idx.Entries}}<li><a href="{{.Slug}}.html">{{.Name}}</a></li>
{{end}}</ul>{{end}}`

const entryTmpl = `{{define "title"}}{{.Entry.Name}}{{end}}
{{define "sidebar"}}
<nav><ul>
{{range .Idx.Entries}}<li><a href="{{.Slug}}.html">{{.Name}}</a></li>
{{end}}</ul></nav>
{{end}}
{{define "content"}}<h1>{{.Entry.Name}}</h1><p>{{.Entry.Description}}</p>
<h3>Key Types</h3><ul>{{range .Entry.KeyTypes}}<li><code>{{.}}</code></li>{{end}}</ul>
{{end}}`

// RenderIndex writes the rendered HTML for the docs index page.
func RenderIndex(w io.Writer, r *render.Renderer, idx Index) error {
	data := struct{ Idx Index }{Idx: idx}
	return r.Render(w, indexTmpl, data)
}

// RenderEntry writes the rendered HTML for a single domain entry page.
func RenderEntry(w io.Writer, r *render.Renderer, idx Index, entry DomainEntry) error {
	data := struct {
		Idx   Index
		Entry DomainEntry
	}{Idx: idx, Entry: entry}
	return r.Render(w, entryTmpl, data)
}
