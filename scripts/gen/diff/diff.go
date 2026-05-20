// Package diff defines the data model for diff reports and renders them
// through the shared render.Renderer.
package diff

import (
	"errors"
	"io"
	"strings"

	"github.com/boilerplate/ebiten-template/scripts/gen/render"
)

// Line is a single diff line with a kind tag for CSS styling.
type Line struct {
	Kind string // "add" | "del" | "ctx"
	Text string
}

// Hunk groups contiguous diff lines under a hunk header.
type Hunk struct {
	Header string
	Lines  []Line
}

// FileDiff is the per-file collection of hunks.
type FileDiff struct {
	Path  string
	Hunks []Hunk
}

// Report is the top-level diff report payload.
type Report struct {
	Title       string
	Explanation string
	Files       []FileDiff
}

// reportData is the template data struct passed to the renderer.
type reportData struct {
	Title      string
	Paragraphs []string
	Files      []FileDiff
}

const diffTmpl = `{{define "title"}}{{.Title}}{{end}}
{{define "content"}}
{{- range .Paragraphs}}<p>{{.}}</p>{{end}}
{{range .Files}}<h2>{{.Path}}</h2>
{{range .Hunks}}<pre class="hunk"><span class="hunk-header">{{.Header}}</span>
{{range .Lines}}<span class="{{.Kind}}">{{.Text}}</span>
{{end}}</pre>{{end}}{{end}}
{{end}}`

// paragraphs splits text on blank lines and returns non-empty paragraphs.
func paragraphs(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, "\n\n")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

// Render writes the rendered HTML diff report to w.
func Render(w io.Writer, r *render.Renderer, rep Report) error {
	if rep.Title == "" {
		return errors.New("diff: report title must not be empty")
	}
	data := reportData{
		Title:      rep.Title,
		Paragraphs: paragraphs(rep.Explanation),
		Files:      rep.Files,
	}
	return r.Render(w, diffTmpl, data)
}
