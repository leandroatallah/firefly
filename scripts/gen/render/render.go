// Package render provides a shared HTML template renderer built on
// html/template.
package render

import (
	"errors"
	"html/template"
	"io"
	"regexp"
	"strings"
)

// Renderer renders child templates against a shared base template.
type Renderer struct {
	base *template.Template
}

// New parses the supplied base.html bytes and returns a *Renderer.
// Returns an error containing "empty" if baseHTML is nil or zero-length.
func New(baseHTML []byte) (*Renderer, error) {
	if len(baseHTML) == 0 {
		return nil, errors.New("base template is empty")
	}
	tmpl, err := template.New("base").Parse(string(baseHTML))
	if err != nil {
		return nil, err
	}
	return &Renderer{base: tmpl}, nil
}

// Render clones the base template, parses childTmpl into the clone, then
// executes the "base" template writing the result to w. Cloning per call
// ensures block overrides do not accumulate across calls.
func (r *Renderer) Render(w io.Writer, childTmpl string, data any) error {
	clone, err := r.base.Clone()
	if err != nil {
		return err
	}
	_, err = clone.Parse(childTmpl)
	if err != nil {
		return err
	}
	return clone.ExecuteTemplate(w, "base", data)
}

var nonAlnum = regexp.MustCompile(`[^a-z0-9]+`)

// Slugify converts s into a filename-safe slug: lowercase, runs of
// non-[a-z0-9] characters replaced with "-", with leading/trailing "-" trimmed.
func Slugify(s string) string {
	s = strings.ToLower(s)
	s = nonAlnum.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}
