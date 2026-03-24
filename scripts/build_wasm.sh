#!/bin/bash
# scripts/build_wasm.sh

# Exit on error
set -e

# Project root
ROOT_DIR=$(pwd)
ZIP_NAME="growbel-wasm.zip"
OUTPUT_WASM="game.wasm"
WASM_EXEC_JS="wasm_exec.js"
INDEX_HTML="index.html"

echo "Building WASM binary..."
GOOS=js GOARCH=wasm go build -o "$OUTPUT_WASM" main.go

# Ensure wasm_exec.js is up-to-date with the current Go version
echo "Updating wasm_exec.js from Go root..."
WASM_EXEC_PATH="$(go env GOROOT)/lib/wasm/wasm_exec.js"
if [ ! -f "$WASM_EXEC_PATH" ]; then
    WASM_EXEC_PATH="$(go env GOROOT)/misc/wasm/wasm_exec.js"
fi

if [ -f "$WASM_EXEC_PATH" ]; then
    cp "$WASM_EXEC_PATH" .
else
    echo "Warning: wasm_exec.js not found in GOROOT, using existing one if present."
fi

echo "Creating zip file..."
# -j flag to junk paths (don't include directories if they were built in a subfolder, 
# though here they are at the root)
# But here they are all in root, so zip is fine.
zip -r "$ZIP_NAME" "$OUTPUT_WASM" "$INDEX_HTML" "$WASM_EXEC_JS"

# Clean up the intermediate wasm file if desired? 
# Usually, keeping it is fine for testing.
# rm "$OUTPUT_WASM"

echo "Done! Zip file created at $ROOT_DIR/$ZIP_NAME"
