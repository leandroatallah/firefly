package diff_test

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/boilerplate/ebiten-template/scripts/gen/diff"
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

// T-D1: diff renders file path as h2.
func TestRender_FilePathAsH2(t *testing.T) {
	r := newRenderer(t)
	rep := diff.Report{
		Title: "Test",
		Files: []diff.FileDiff{
			{
				Path: "a.go",
				Hunks: []diff.Hunk{
					{Header: "@@ -1 +1 @@", Lines: []diff.Line{{Kind: "ctx", Text: "x"}}},
				},
			},
		},
	}
	var buf bytes.Buffer
	if err := diff.Render(&buf, r, rep); err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(buf.String(), "<h2>a.go</h2>") {
		t.Fatalf("expected <h2>a.go</h2> in output; got:\n%s", buf.String())
	}
}

// T-D2: diff line kinds get correct CSS class.
func TestRender_LineKindCSSClasses(t *testing.T) {
	r := newRenderer(t)
	rep := diff.Report{
		Title: "Kinds",
		Files: []diff.FileDiff{
			{
				Path: "a.go",
				Hunks: []diff.Hunk{
					{
						Header: "@@",
						Lines: []diff.Line{
							{Kind: "add", Text: "added"},
							{Kind: "del", Text: "deleted"},
							{Kind: "ctx", Text: "context"},
						},
					},
				},
			},
		},
	}
	var buf bytes.Buffer
	if err := diff.Render(&buf, r, rep); err != nil {
		t.Fatalf("Render: %v", err)
	}
	out := buf.String()
	for _, needle := range []string{`class="add"`, `class="del"`, `class="ctx"`} {
		if got := strings.Count(out, needle); got != 1 {
			t.Errorf("expected exactly 1 occurrence of %s, got %d; output:\n%s",
				needle, got, out)
		}
	}
}

// T-D3: diff escapes HTML in line text.
func TestRender_EscapesHTMLInLineText(t *testing.T) {
	r := newRenderer(t)
	payload := "<script>x</script>"
	rep := diff.Report{
		Title: "Escape",
		Files: []diff.FileDiff{
			{
				Path: "a.go",
				Hunks: []diff.Hunk{
					{Header: "@@", Lines: []diff.Line{{Kind: "ctx", Text: payload}}},
				},
			},
		},
	}
	var buf bytes.Buffer
	if err := diff.Render(&buf, r, rep); err != nil {
		t.Fatalf("Render: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "&lt;script&gt;") {
		t.Errorf("expected escaped &lt;script&gt; in output, got:\n%s", out)
	}
	if strings.Contains(out, "<script>") {
		t.Errorf("output must not contain literal <script>; got:\n%s", out)
	}
}

// T-D4: cmd/diff writes to output/tmp/<slug>.html under cwd.
func TestCmdDiff_WritesToOutputTmp(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping cmd test in -short mode")
	}
	repoRoot := findRepoRoot(t)
	tmp := t.TempDir()

	rep := diff.Report{
		Title: "My Diff",
		Files: []diff.FileDiff{
			{
				Path: "a.go",
				Hunks: []diff.Hunk{
					{Header: "@@", Lines: []diff.Line{{Kind: "ctx", Text: "x"}}},
				},
			},
		},
	}
	body, err := json.Marshal(rep)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	bin := buildCmd(t, repoRoot, "diff")
	if err := runWithStdin(bin, tmp, bytes.NewReader(body)); err != nil {
		t.Fatalf("cmd/diff failed: %v", err)
	}

	want := filepath.Join(tmp, "output", "tmp", "my-diff.html")
	if _, err := os.Stat(want); err != nil {
		t.Fatalf("expected output file at %s, stat err: %v", want, err)
	}
}

// buildCmd compiles scripts/gen/cmd/<name> into a temp binary using the
// repository's module context, returning the binary path.
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

// runWithStdin runs bin with cwd=workDir and pipes stdin into it.
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

// findRepoRoot walks upward from the test's cwd to locate go.mod.
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
