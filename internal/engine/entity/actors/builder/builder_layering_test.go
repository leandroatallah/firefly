package builder_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestBuilderDoesNotImportKit asserts that the builder source files do not
// import any package under internal/kit/. This guards the dependency inversion
// introduced in story 050: the builder (engine layer) must not depend on kit.
func TestBuilderDoesNotImportKit(t *testing.T) {
	builderDir := "."
	entries, err := os.ReadDir(builderDir)
	if err != nil {
		t.Fatalf("failed to read builder directory: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}
		if strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}

		path := filepath.Join(builderDir, entry.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("failed to read %s: %v", path, err)
		}

		if strings.Contains(string(content), `"internal/kit/`) {
			t.Errorf("%s imports internal/kit/... (forbidden for engine layer)", entry.Name())
		}
	}
}
