package render_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/boilerplate/ebiten-template/scripts/gen/render"
)

const validBase = `<!doctype html>
<html><head>
  <title>{{block "title" .}}Firefly Report{{end}}</title>
  <style>/* css */</style>
</head><body>
  <aside>{{block "sidebar" .}}{{end}}</aside>
  <main>{{block "content" .}}{{end}}</main>
</body></html>`

// T-R1: render.New rejects empty base.
func TestNew_NilBase_ReturnsError(t *testing.T) {
	r, err := render.New(nil)
	if err == nil {
		t.Fatalf("expected error for nil base, got nil")
	}
	if r != nil {
		t.Fatalf("expected nil renderer on error, got %#v", r)
	}
	if !strings.Contains(strings.ToLower(err.Error()), "empty") {
		t.Fatalf("expected error message to contain 'empty', got %q", err.Error())
	}
}

// T-R2: render.New parses valid base.
func TestNew_ValidBase_ReturnsRenderer(t *testing.T) {
	r, err := render.New([]byte(validBase))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if r == nil {
		t.Fatalf("expected non-nil renderer")
	}
}

// T-R3: child omitting sidebar block renders empty <aside>, no panic.
func TestRender_MissingSidebarBlock_RendersEmptyAside(t *testing.T) {
	r, err := render.New([]byte(validBase))
	if err != nil {
		t.Fatalf("setup: New: %v", err)
	}
	child := `{{define "title"}}Hello{{end}}{{define "content"}}<p>body</p>{{end}}`
	var buf bytes.Buffer
	if err := r.Render(&buf, child, nil); err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "<aside></aside>") && !strings.Contains(out, "<aside>\n  </aside>") && !strings.Contains(out, "<aside>  </aside>") {
		t.Fatalf("expected empty <aside> section in output, got:\n%s", out)
	}
}

// T-R4: rendered HTML has no external resources.
func TestRender_NoExternalResources(t *testing.T) {
	r, err := render.New([]byte(validBase))
	if err != nil {
		t.Fatalf("setup: New: %v", err)
	}
	child := `{{define "title"}}T{{end}}{{define "content"}}<p>c</p>{{end}}{{define "sidebar"}}<nav/>{{end}}`
	var buf bytes.Buffer
	if err := r.Render(&buf, child, nil); err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	out := buf.String()
	forbidden := []string{"http://", "https://", "<script", "cdn."}
	for _, needle := range forbidden {
		if strings.Contains(out, needle) {
			t.Errorf("rendered output must not contain %q; output:\n%s", needle, out)
		}
	}
	// Sanity: ensure the renderer actually produced something.
	if len(out) == 0 {
		t.Fatalf("expected non-empty output")
	}
}

// Slugify cases.
func TestSlugify(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"My Diff!", "my-diff"},
		{"a/b c", "a-b-c"},
		{"", ""},
	}
	for _, c := range cases {
		c := c
		t.Run(c.in, func(t *testing.T) {
			got := render.Slugify(c.in)
			if got != c.want {
				t.Fatalf("Slugify(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}
