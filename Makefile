.PHONY: build-wasm clean sync-skills sync-agents setup help

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build-wasm   Build the WASM binary and create a zip for deployment"
	@echo "  sync-skills  Sync skill files to all AI tools"
	@echo "  sync-agents  Sync agent files to all AI tools"
	@echo "  setup        Install git hooks"
	@echo "  clean        Remove build artifacts"
	@echo "  help         Show this help message"

build-wasm:
	@bash scripts/build_wasm.sh

sync-skills:
	@bash scripts/sync-skills.sh

sync-agents:
	@bash scripts/sync-agents.sh

setup:
	@cp scripts/hooks/pre-commit .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "Git hooks installed."

clean:
	rm -f game.wasm growbel-wasm.zip
