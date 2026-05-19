// Command domain-docs reads a domaindocs.Index from stdin as JSON and
// writes output/docs/index.html plus one output/docs/<slug>.html per entry,
// under the current working directory.
package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/boilerplate/ebiten-template/scripts/gen/baseasset"
	"github.com/boilerplate/ebiten-template/scripts/gen/domaindocs"
	"github.com/boilerplate/ebiten-template/scripts/gen/render"
)

func main() {
	var idx domaindocs.Index
	if err := json.NewDecoder(os.Stdin).Decode(&idx); err != nil {
		log.Fatalf("decode: %v", err)
	}
	r, err := render.New(baseasset.HTML)
	if err != nil {
		log.Fatalf("renderer: %v", err)
	}
	outDir := filepath.Join("output", "docs")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		log.Fatalf("mkdir: %v", err)
	}

	// Write index.html
	indexPath := filepath.Join(outDir, "index.html")
	indexFile, err := os.Create(indexPath)
	if err != nil {
		log.Fatalf("create index: %v", err)
	}
	if err := domaindocs.RenderIndex(indexFile, r, idx); err != nil {
		indexFile.Close()
		log.Fatalf("render index: %v", err)
	}
	indexFile.Close()

	// Write per-entry pages
	for _, entry := range idx.Entries {
		entryPath := filepath.Join(outDir, entry.Slug+".html")
		entryFile, err := os.Create(entryPath)
		if err != nil {
			log.Fatalf("create entry %s: %v", entry.Slug, err)
		}
		if err := domaindocs.RenderEntry(entryFile, r, idx, entry); err != nil {
			entryFile.Close()
			log.Fatalf("render entry %s: %v", entry.Slug, err)
		}
		entryFile.Close()
	}
}
