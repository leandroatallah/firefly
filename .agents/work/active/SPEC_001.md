# Technical Specification: SPEC_001 - Create Clean Boilerplate Repository

## Overview
This specification details the process of transforming the Growbel game repository into a clean, minimal Ebitengine boilerplate. The goal is to strip game-specific logic and assets while preserving the engine core and the platformer foundation.

## 1. Module Rename Strategy
All references to the current module name must be updated to the new template name.

- **New Module Name:** `github.com/boilerplate/ebiten-template`
- **Steps:**
  1.  Update `go.mod`: `module github.com/boilerplate/ebiten-template`
  2.  Recursively find and replace all instances of `github.com/leandroatallah/firefly` (the current module name) with `github.com/boilerplate/ebiten-template` in all `.go` files.
  3.  Update the `index.html` (if it contains module references for WASM) to reflect any name changes if necessary.

## 2. Asset Sanitization
The `assets/` directory must be cleared of game-specific content but maintain its structure for the engine.

- **Keep (Do NOT delete):**
  - All files in `assets/fonts/` (`monogram.ttf`, `pressstart2p.ttf`, `tiny5.ttf`).
  - `assets/images/empty.png` and `assets/images/grid-32.png` (essential for debugging and placeholder use).
  - All `.gitkeep` files in empty subdirectories.
- **Delete (Except subdirectories and .gitkeep):**
  - `assets/audio/`: Delete all `.ogg` files.
  - `assets/images/`: Delete all `.png` files (except the ones kept).
  - `assets/particles/`: Delete all `.json` files.
  - `assets/sequences/`: Delete all `.json` files.
  - `assets/tilemap/`: Delete all `.tmj`, `.tsx`, and `.png` files.
- **Note:** After deletion, ensure each empty directory contains a `.gitkeep` file to preserve the directory structure in Git.

## 3. I18n Management
The translation files must be minimal.

- **Files:** `assets/lang/en.json`, `assets/lang/pt-br.json`.
- **Strategy:**
  - Remove all Growbel-specific keys (e.g., character names, story text).
  - Retain only generic keys needed for the stub menu:
    ```json
    {
      "menu.start": "Start Game",
      "menu.exit": "Exit",
      "ui.loading": "Loading..."
    }
    ```

## 4. Minimal Stub Scene Structure
The `internal/game/` logic must be simplified to provide a clean starting point.

- **Preserve:**
  - `internal/game/scenes/phases/`: Preserve this directory as it contains the platformer logic.
- **Stub/Clean:**
  - `internal/game/scenes/`: Remove all non-essential scene implementations (Intro, Summary, Credits, Story, etc.).
  - Create/Update a minimal `internal/game/scenes/menu/` that just displays a "Start Game" option using the `menu.start` i18n key.
  - Simplify `internal/game/scenes/factory.go` (or whichever file uses `scenestypes`) to only register:
    - `SceneMenu`
    - `ScenePhases`
- **Phases Configuration:**
  - The function `GetPhases()` (located in `internal/game/app/` in the same package as `Setup`) should be updated to return an empty slice of `phases.Phase` (or a single example phase referencing a non-existent/placeholder tilemap that the engine handles gracefully).

## 5. Main Entry Point and Setup
Update the application entry point to initialize the minimal environment.

- **`main.go`:** Ensure imports are updated to the new module name.
- **`internal/game/app/setup.go`:**
  - Update `ebiten.SetWindowTitle("Ebitengine Boilerplate")`.
  - Remove Growbel-specific initializations (e.g., custom speech bubble typing sounds if they were deleted from `assets/audio`).
  - Update the initial navigation to start at `scenestypes.SceneMenu`.
  - Ensure `appContext` is populated with minimal dependencies.
- **`internal/game/app/config.go`:**
  - Keep default resolutions (256x240) and physics settings as they are a good baseline.
  - Update any Growbel-specific flag defaults if they exist.

## 6. Metadata and Documentation
Update repository metadata to reflect its new purpose.

- **`README.md`:** Replace content with a description of the boilerplate, installation instructions, and how to start a new project.
- **`TODO.md`:** Clear or replace with boilerplate-specific tasks (e.g., "Implement your first scene").
- **`growbel.gif`:** Delete or replace with a generic placeholder.

## 7. Verification Steps
The following must be verified after implementation:

1.  **Compilation:** `go build ./...` must succeed without errors.
2.  **Linting/Imports:** All imports must correctly point to `github.com/boilerplate/ebiten-template`.
3.  **Runtime:** `go run main.go` must:
    - Open a window titled "Ebitengine Boilerplate".
    - Show the minimal Menu scene.
    - Not crash due to missing assets (I18n keys or images).
4.  **Asset Structure:** Verify that `assets/` subdirectories exist but are mostly empty (except for `.gitkeep`, fonts, and debug images).
5.  **Phase Logic:** Ensure navigating to `ScenePhases` from the menu (if implemented) doesn't crash the engine even if no phases are defined (it should handle the empty case).
