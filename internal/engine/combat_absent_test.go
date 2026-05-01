package engine_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// TestEngineCombatDirectoryAbsent asserts that the legacy
// internal/engine/combat directory has been removed. Its concrete combat
// implementations are migrated to internal/kit/combat as part of story 049.
//
// The path is computed relative to this test file's location so the test is
// independent of the caller's working directory.
func TestEngineCombatDirectoryAbsent(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed; cannot determine test file location")
	}

	combatDir := filepath.Join(filepath.Dir(thisFile), "combat")

	_, err := os.Stat(combatDir)
	if err == nil {
		t.Fatalf("expected %s to be absent, but it still exists", combatDir)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("unexpected error stat-ing %s: %v (want fs.ErrNotExist)", combatDir, err)
	}
}
