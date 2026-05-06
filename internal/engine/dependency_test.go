package engine_test

import (
	"os/exec"
	"regexp"
	"runtime"
	"testing"
)

// TestEngineDoesNotDependOnKitOrGame is the layer-rule safety net for
// story 052-kit-ui-split. It must continue to pass after the dialogue
// orchestrator relocates from internal/engine/ui/speech to
// internal/kit/ui/speech, proving that no engine package was naively
// re-pointed at kit during the move.
//
// The constitution forbids any package under internal/engine/... from
// transitively importing internal/kit/... or internal/game/...
func TestEngineDoesNotDependOnKitOrGame(t *testing.T) {
	// Locate module root relative to this test file.
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not determine test file path")
	}
	// filename is .../internal/engine/dependency_test.go; strip that suffix to get module root.
	moduleRoot := filename[:len(filename)-len("internal/engine/dependency_test.go")]
	cmd := exec.Command("go", "list", "-deps", "github.com/boilerplate/ebiten-template/internal/engine/...")
	cmd.Dir = moduleRoot
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("go list: %v", err)
	}
	// Match only our project's kit/game packages, not ebiten's internal/gamepad etc.
	re := regexp.MustCompile(`boilerplate/ebiten-template/internal/(kit|game)`)
	if loc := re.FindString(string(out)); loc != "" {
		t.Fatalf("engine has forbidden dep: %s", loc)
	}
}
