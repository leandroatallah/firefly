package engine_test

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestEngineLayerHasNoKitOrGameDependencies asserts that no package under
// internal/engine/... transitively imports any package under
// internal/kit/..., internal/game/..., or internal/engine/combat/....
//
// The first two prefixes enforce the three-layer architecture (engine must
// not depend on genre-specific or game-specific code).
//
// The third prefix encodes the migration goal of story 049: the legacy
// `internal/engine/combat` sub-tree is being relocated to
// `internal/kit/combat`. After the move, no package under internal/engine
// must transitively pull in the (deleted) engine/combat path. Today this
// assertion fails because:
//   - internal/engine/app/context.go imports engine/combat/projectile
//   - internal/engine/entity/actors/character.go imports engine/combat
//   - internal/engine/entity/actors/builder/configure_enemy_weapon.go
//     imports engine/combat/weapon
//
// All three call sites must be refactored (per SPEC §5) before the test
// goes green.
//
// The test shells out to `go list -deps` and runs from the module root
// (discovered via `go env GOMOD`) so the result is independent of the
// test's working directory.
func TestEngineLayerHasNoKitOrGameDependencies(t *testing.T) {
	const modulePath = "github.com/boilerplate/ebiten-template"

	moduleRoot := findModuleRoot(t)

	forbiddenPrefixes := []struct {
		name   string
		prefix string
	}{
		{name: "kit", prefix: modulePath + "/internal/kit"},
		{name: "game", prefix: modulePath + "/internal/game"},
		{name: "engine_combat_legacy", prefix: modulePath + "/internal/engine/combat"},
	}

	cmd := exec.Command("go", "list", "-deps", "-f", "{{.ImportPath}}", "./internal/engine/...")
	cmd.Dir = moduleRoot
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("go list failed: %v\nstderr: %s", err, stderr.String())
	}

	deps := strings.Split(strings.TrimSpace(stdout.String()), "\n")

	for _, fp := range forbiddenPrefixes {
		t.Run(fp.name, func(t *testing.T) {
			var violations []string
			for _, dep := range deps {
				dep = strings.TrimSpace(dep)
				if dep == "" {
					continue
				}
				// Skip self — engine/combat itself is in the engine tree
				// today, and this loop also enumerates it as a "consumer".
				// We only care about dependencies declared by OTHER engine
				// packages. Once engine/combat is deleted, this filter is
				// a no-op.
				if dep == fp.prefix || strings.HasPrefix(dep, fp.prefix+"/") {
					// We still want to flag it as a violation if some
					// non-combat engine package imports it, so we check
					// the full deps list below by other means. For
					// simplicity, list every dependency match here — the
					// presence of engine/combat/* in the transitive
					// closure of `./internal/engine/...` is itself the
					// violation we want to catch (it means at least one
					// non-combat engine package pulls it in, OR the
					// combat package still lives under engine and is
					// reachable, both of which are forbidden post-migration).
					violations = append(violations, dep)
				}
			}
			if len(violations) > 0 {
				t.Errorf("engine layer transitively depends on %s packages (forbidden):\n  %s",
					fp.name, strings.Join(violations, "\n  "))
			}
		})
	}
}

// findModuleRoot returns the directory containing go.mod for the current
// module by invoking `go env GOMOD`.
func findModuleRoot(t *testing.T) string {
	t.Helper()
	cmd := exec.Command("go", "env", "GOMOD")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("go env GOMOD failed: %v", err)
	}
	gomod := strings.TrimSpace(string(out))
	if gomod == "" || gomod == "/dev/null" {
		t.Fatalf("could not locate go.mod (GOMOD=%q)", gomod)
	}
	return filepath.Dir(gomod)
}
