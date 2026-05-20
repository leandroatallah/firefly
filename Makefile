.PHONY: build-wasm clean sync-agents sync-skills setup help dashboard kanban build-gen gen-diff gen-docs serve

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  dashboard    Show real-time SDD pipeline dashboard"
	@echo "  kanban       Generate kanban.html from stories"
	@echo "  build-wasm   Build the WASM binary and create a zip for deployment"
	@echo "  build-gen    Build HTML report generators (diff, domain-docs)"
	@echo "  gen-diff     Generate diff report from JSON input (INPUT=path/to/input.json)"
	@echo "  gen-docs     Generate domain docs from JSON input (INPUT=path/to/input.json)"
	@echo "  serve        Start local HTTP server to preview reports (http://localhost:8080)"
	@echo "  sync-skills  Create/update skill symlinks for all AI tools"
	@echo "  sync-agents  Sync agent files to all AI tools"
	@echo "  setup        Install git hooks and create skill symlinks"
	@echo "  clean        Remove build artifacts"
	@echo "  help         Show this help message"

dashboard:
	@bash scripts/story.sh dashboard

kanban:
	@go run scripts/kanban.go
	@echo "kanban.html generated."

build-wasm:
	@bash scripts/build_wasm.sh

sync-skills:
	@bash scripts/setup-skill-symlinks.sh

sync-agents:
	@bash scripts/sync-agents.sh

setup:
	@bash scripts/setup-skill-symlinks.sh
	@lefthook install 2>/dev/null || $(HOME)/go/bin/lefthook install
	@echo "Git hooks installed."

clean:
	rm -f game.wasm growbel-wasm.zip

build-gen:
	@go build -o scripts/gen/cmd/diff/diff ./scripts/gen/cmd/diff
	@go build -o scripts/gen/cmd/domain-docs/domain-docs ./scripts/gen/cmd/domain-docs
	@echo "Generators built."

gen-diff:
	@if [ ! -f scripts/gen/cmd/diff/diff ]; then $(MAKE) build-gen; fi
	@if [ -z "$(INPUT)" ]; then echo "Usage: make gen-diff INPUT=path/to/input.json"; exit 1; fi
	@cat $(INPUT) | ./scripts/gen/cmd/diff/diff
	@echo "Diff report generated."

gen-docs:
	@if [ ! -f scripts/gen/cmd/domain-docs/domain-docs ]; then $(MAKE) build-gen; fi
	@if [ -z "$(INPUT)" ]; then echo "Usage: make gen-docs INPUT=path/to/input.json"; exit 1; fi
	@cat $(INPUT) | ./scripts/gen/cmd/domain-docs/domain-docs
	@echo "Domain docs generated."

serve:
	@go run scripts/serve.go -port=:8080
