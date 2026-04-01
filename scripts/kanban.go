//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func cards(dir string) string {
	entries, _ := os.ReadDir(dir)
	var b strings.Builder
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		f, err := os.Open(filepath.Join(dir, e.Name(), "USER_STORY.md"))
		if err != nil {
			continue
		}
		scanner := bufio.NewScanner(f)
		scanner.Scan()
		f.Close()
		title := strings.TrimPrefix(scanner.Text(), "# ")
		fmt.Fprintf(&b, `<div class="card">%s</div>`+"\n", title)
	}
	return b.String()
}

func main() {
	cols := []struct{ id, label string }{
		{"backlog", "Backlog"},
		{"active", "Active"},
		{"done", "Done"},
	}

	var cols_html strings.Builder
	for _, c := range cols {
		fmt.Fprintf(&cols_html, `<div class="col"><h2>%s</h2>%s</div>`+"\n",
			c.label, cards(filepath.Join(".agents/work", c.id)))
	}

	html := `<!doctype html><html><head><meta charset="utf-8"><title>Kanban</title>
<style>
body{font-family:sans-serif;display:flex;gap:1rem;padding:1rem;background:#f4f4f4}
.col{flex:1;background:#e2e8f0;border-radius:8px;padding:.75rem}
h2{margin:0 0 .75rem;font-size:1rem;text-transform:uppercase;letter-spacing:.05em}
.card{background:#fff;border-radius:6px;padding:.6rem .8rem;margin-bottom:.5rem;
      box-shadow:0 1px 3px rgba(0,0,0,.1);font-size:.875rem}
</style></head><body>` + cols_html.String() + `</body></html>`

	os.WriteFile("kanban.html", []byte(html), 0644)
	fmt.Println("kanban.html written.")
}
