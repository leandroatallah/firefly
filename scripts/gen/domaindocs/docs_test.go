package domaindocs_test

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/boilerplate/ebiten-template/scripts/gen/domaindocs"
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

func newRenderer(t *testing.T) *render.Renderer {
	t.Helper()
	r, err := render.New([]byte(validBase))
	if err != nil {
		t.Fatalf("render.New: %v", err)
	}
	return r
}

func sampleIndex() domaindocs.Index {
	return domaindocs.Index{
		Title: "Domains",
		Entries: []domaindocs.DomainEntry{
			{Name: "Physics", Slug: "physics", Description: "p", KeyTypes: []string{"Body"}},
			{Name: "Scene", Slug: "scene", Description: "s", KeyTypes: []string{"Scene"}},
		},
	}
}

// T-X1: index renders sidebar with all entries.
func TestRenderIndex_SidebarContainsAllEntries(t *testing.T) {
	r := newRenderer(t)
	var buf bytes.Buffer
	if err := domaindocs.RenderIndex(&buf, r, sampleIndex()); err != nil {
		t.Fatalf("RenderIndex: %v", err)
	}
	out := buf.String()
	for _, needle := range []string{`href="physics.html"`, `href="scene.html"`} {
		if !strings.Contains(out, needle) {
			t.Errorf("expected sidebar to contain %s; output:\n%s", needle, out)
		}
	}
}

// T-X2: entry page uses relative links matching ^[a-z0-9-]+\.html$.
func TestRenderEntry_SidebarLinksAreRelative(t *testing.T) {
	r := newRenderer(t)
	idx := sampleIndex()
	var buf bytes.Buffer
	if err := domaindocs.RenderEntry(&buf, r, idx, idx.Entries[0]); err != nil {
		t.Fatalf("RenderEntry: %v", err)
	}
	out := buf.String()

	hrefRE := regexp.MustCompile(`href="([^"]+)"`)
	matches := hrefRE.FindAllStringSubmatch(out, -1)
	if len(matches) == 0 {
		t.Fatalf("expected at least one href in sidebar; output:\n%s", out)
	}
	valid := regexp.MustCompile(`^[a-z0-9-]+\.html$`)
	for _, m := range matches {
		href := m[1]
		if strings.HasPrefix(href, "/") {
			t.Errorf("href %q must not be absolute", href)
		}
		if strings.Contains(href, "://") {
			t.Errorf("href %q must not contain scheme", href)
		}
		if !valid.MatchString(href) {
			t.Errorf("href %q does not match ^[a-z0-9-]+\\.html$", href)
		}
	}
}

// T-X3: cmd/domain-docs writes index.html plus one file per entry.
func TestCmdDomainDocs_WritesOutputDocs(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping cmd test in -short mode")
	}
	repoRoot := findRepoRoot(t)
	tmp := t.TempDir()

	body, err := json.Marshal(sampleIndex())
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	bin := buildCmd(t, repoRoot, "domain-docs")
	if err := runWithStdin(bin, tmp, bytes.NewReader(body)); err != nil {
		t.Fatalf("cmd/domain-docs failed: %v", err)
	}

	for _, rel := range []string{
		filepath.Join("output", "docs", "index.html"),
		filepath.Join("output", "docs", "physics.html"),
		filepath.Join("output", "docs", "scene.html"),
	} {
		want := filepath.Join(tmp, rel)
		if _, err := os.Stat(want); err != nil {
			t.Errorf("expected output file at %s, stat err: %v", want, err)
		}
	}
}

func buildCmd(t *testing.T, repoRoot, name string) string {
	t.Helper()
	binName := name
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}
	bin := filepath.Join(t.TempDir(), binName)
	pkg := "./scripts/gen/cmd/" + name
	build := exec.Command("go", "build", "-o", bin, pkg)
	build.Dir = repoRoot
	var stderr bytes.Buffer
	build.Stderr = &stderr
	if err := build.Run(); err != nil {
		t.Fatalf("go build %s: %v\nstderr:\n%s", pkg, err, stderr.String())
	}
	return bin
}

func runWithStdin(bin, workDir string, stdin io.Reader) error {
	cmd := exec.Command(bin)
	cmd.Dir = workDir
	cmd.Stdin = stdin
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return &runErr{err: err, stderr: stderr.String()}
	}
	return nil
}

type runErr struct {
	err    error
	stderr string
}

func (e *runErr) Error() string {
	return e.err.Error() + "\nstderr:\n" + e.stderr
}

func findRepoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("could not locate go.mod above %s", dir)
		}
		dir = parent
	}
}
