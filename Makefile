.PHONY: build-wasm clean help

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build-wasm  Build the WASM binary and create a zip for deployment"
	@echo "  clean       Remove build artifacts"
	@echo "  help        Show this help message"

build-wasm:
	@bash scripts/build_wasm.sh

clean:
	rm -f game.wasm growbel-wasm.zip
