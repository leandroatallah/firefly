// Command diff reads a diff.Report from stdin as JSON and writes the
// rendered HTML to output/tmp/<slug>.html under the current working directory.
package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/boilerplate/ebiten-template/scripts/gen/baseasset"
	"github.com/boilerplate/ebiten-template/scripts/gen/diff"
	"github.com/boilerplate/ebiten-template/scripts/gen/render"
)

func main() {
	var rep diff.Report
	if err := json.NewDecoder(os.Stdin).Decode(&rep); err != nil {
		log.Fatalf("decode: %v", err)
	}
	r, err := render.New(baseasset.HTML)
	if err != nil {
		log.Fatalf("renderer: %v", err)
	}
	slug := render.Slugify(rep.Title)
	outDir := filepath.Join("output", "tmp")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		log.Fatalf("mkdir: %v", err)
	}
	outPath := filepath.Join(outDir, slug+".html")
	f, err := os.Create(outPath)
	if err != nil {
		log.Fatalf("create: %v", err)
	}
	defer f.Close()
	if err := diff.Render(f, r, rep); err != nil {
		log.Fatalf("render: %v", err)
	}
}
