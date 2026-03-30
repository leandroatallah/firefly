# User Story: Create Clean Boilerplate Repository

**Status:** Backlog
**ID:** USER_STORY_001

## Story Description

**As a** game developer,
**I want** a clean, minimal Ebitengine project boilerplate,
**So that** I can quickly start new game projects without manually stripping game-specific code.

## Acceptance Criteria

1.  **Engine Integrity:** The `internal/engine/` directory must remain completely intact and functional, ensuring that the core engine logic is ready for any new project.
2.  **Platformer Foundation:** `internal/game/scenes/phases/` must be preserved as it provides the core platformer logic for future projects.
3.  **Game Stubbing:** All other non-essential files in `internal/game/` must be replaced with minimal stubs (e.g., an empty Menu scene) that satisfy the engine's requirements and provide a starting point for new content.
4.  **Asset Sanitization:** 
    - `assets/` subdirectories (audio, images, particles, sequences, tilemap) must be cleared of game-specific content but maintain their structure via `.gitkeep` files.
    - Debugging assets (`empty.png`, `grid-32.png`) and all `assets/fonts/` must be kept to support initial development.
    - `assets/lang/` files must be stripped of Growbel-specific keys, keeping only engine/stub keys (e.g., generic UI labels).
5.  **Module Refactoring:** 
    - The `go.mod` module name must be changed to `github.com/boilerplate/ebiten-template`.
    - All internal import paths in `.go` files must be updated to match the new module name.
6.  **Metadata Update:** `README.md` and `TODO.md` must be rewritten to reflect a generic boilerplate project, removing Growbel-specific instructions and goals.
7.  **Functional Verification:** Running `go run main.go` must open a window displaying the minimal stub scene without errors, confirming the boilerplate is ready for development.

## Behavioral Edge Cases

- **Asset Path Resolution:** Ensure that the stub scene doesn't attempt to load assets that were removed during sanitization, which would cause a panic on startup.
- **I18n Fallbacks:** If a stub scene uses a translation key that was accidentally stripped, the `I18nManager` should return the key itself instead of crashing, as per engine design.
- **Empty Phases:** If `internal/game/scenes/phases/` is kept but no levels exist in `assets/tilemap/`, the engine should handle this gracefully (e.g., showing a "No levels found" message or staying in the stub menu).
- **Import Path Consistency:** A single missed import path update should be caught during the "Functional Verification" step (it will fail to compile).
- **Missing Audio:** Ensure the `AudioPlayer` handles empty audio directories without error, using `.gitkeep` files to maintain structure.

## Sub-tasks (Estimated)

- [ ] Rename Go module and update all import paths.
- [ ] Clear game-specific assets and add `.gitkeep` files.
- [ ] Stub `internal/game` scenes and UI while preserving the `phases` logic.
- [ ] Update `internal/game/app/setup.go` to initialize the minimal stub environment.
- [ ] Refactor `assets/lang/` for generic use.
- [ ] Update documentation (`README.md`, `TODO.md`).
- [ ] Validate the boilerplate by running the application.
