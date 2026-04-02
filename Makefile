.PHONY: build-wasm clean sync-agents sync-skills setup help

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build-wasm   Build the WASM binary and create a zip for deployment"
	@echo "  sync-skills  Create/update skill symlinks for all AI tools"
	@echo "  sync-agents  Sync agent files to all AI tools"
	@echo "  setup        Install git hooks and create skill symlinks"
	@echo "  clean        Remove build artifacts"
	@echo "  help         Show this help message"

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
