package engine_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// TestEnginePhysicsSkillDirectoryAbsent asserts that the legacy
// internal/engine/physics/skill directory has been removed. Its concrete skill
// implementations are migrated to internal/kit/skills as part of story 050.
//
// The path is computed relative to this test file's location so the test is
// independent of the caller's working directory.
func TestEnginePhysicsSkillDirectoryAbsent(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed; cannot determine test file location")
	}

	skillDir := filepath.Join(filepath.Dir(thisFile), "physics", "skill")

	_, err := os.Stat(skillDir)
	if err == nil {
		t.Fatalf("expected %s to be absent, but it still exists", skillDir)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("unexpected error stat-ing %s: %v (want fs.ErrNotExist)", skillDir, err)
	}
}
