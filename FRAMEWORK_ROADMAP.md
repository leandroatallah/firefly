# Firefly Framework Roadmap: Building a Practical 2D Game Starter Kit

## Executive Summary

This document outlines the plan to evolve Firefly from a game-specific implementation into a **personal-use mini-framework/starter kit** for quickly building 2D games with Go/Ebiten. It's a middle ground between a boilerplate and a full framework - providing reusable systems while keeping flexibility for rapid iteration. The architecture already has excellent separation between `internal/engine` (framework code) and `internal/game` (game-specific code), providing a solid foundation.

**Philosophy**: Practical over perfect. Focus on features that help you quickly start new games, not production-ready standards.

## Current State Analysis

### ✅ Strengths

1. **Clean Architecture**: Excellent separation between engine and game-specific code
2. **Modular Design**: Well-organized packages (physics, entities, scenes, rendering, audio, input)
3. **Interface-Driven**: Good use of contracts/interfaces for decoupling
4. **Core Systems Present**:
   - Scene management with transitions
   - Physics system (collision detection, movement models)
   - Entity system (actors, items with state management)
   - Audio system (sound/music management)
   - Rendering pipeline (camera, sprites, tilemaps)
   - Input handling
   - UI system (HUD, dialogue/speech bubbles)
   - Sequence/cutscene system
   - Phase/level management
   - Asset loading (embedded FS support)

5. **Testing Foundation**: One test file exists (`body_movable_test.go`), coverage script available

### ⚠️ Areas for Improvement

1. **Error Handling**: Uses `log.Fatal` in some paths (works but not ideal for debugging)
2. **Testing**: Minimal coverage (1 test file) - enough for critical bugs but could improve confidence
3. **Documentation**: Basic READMEs - could use better inline docs for quick reference
4. **Configuration**: Hard-coded defaults - would benefit from easier customization
5. **Event System**: No event bus - direct coupling works but limits flexibility
6. **Resource Management**: Basic asset loading - works but could be more efficient
7. **Save System**: Not implemented - would be useful for games
8. **Quick Start**: Could use template/project generator for new games
9. **Development Tools**: Basic debug tools - could add more utilities

---

## Roadmap: Practical Improvements for Rapid Game Development

### 🎯 **PHASE 1: Quality of Life Improvements (Weeks 1-2)**

**Goal**: Fix annoying issues and make the codebase easier to work with.

#### 1.1 Error Handling Improvements
- **Priority**: MEDIUM (fix as you encounter issues)
- **Tasks**:
  - Replace `log.Fatal` calls with proper error returns where it makes sense
  - Keep `log.Fatal` for truly unrecoverable errors (that's fine for personal use)
  - Add context to error messages for easier debugging
  - Focus on critical paths that you actually use

#### 1.2 Testing for Confidence
- **Priority**: LOW (test as needed, not comprehensive)
- **Tasks**:
  - Add tests for critical physics bugs (collision, movement)
  - Test complex state transitions (actor states, scene transitions)
  - Add tests when you fix bugs (regression tests)
  - Don't aim for high coverage - just enough to catch breaking changes
  - Keep tests simple and fast

#### 1.3 Better Inline Documentation
- **Priority**: MEDIUM (as you work on code)
- **Tasks**:
  - Add package-level comments explaining what each package does
  - Document non-obvious functions and complex logic
  - Add examples in comments where helpful
  - Focus on "why" not "what" - the code should be readable

#### 1.4 Configuration System
- **Priority**: HIGH (makes starting new games easier)
- **Tasks**:
  - Support configuration from JSON files (simple, human-readable)
  - Keep hard-coded defaults (they're fine)
  - Allow overriding via code (for programmatic setup)
  - Environment variable support (optional, useful for deployment)
  - Keep it simple - don't over-engineer

#### 1.5 Fix Technical Debt from TODO.md
- **Priority**: HIGH (fix bugs as you encounter them)
- **Tasks**:
  - Fix horizontal collision gap when walking left
  - Fix sprite animation for jumping
  - Refactor shape naming if it causes confusion
  - Reduce scene duplication where it makes sense

**Deliverables**:
- ✅ Configuration system for easier customization
- ✅ Critical bugs fixed
- ✅ Better error messages for debugging
- ✅ Enough tests to prevent regressions

---

### 🔧 **PHASE 2: Useful Features (Weeks 3-6)**

**Goal**: Add features that make game development faster and more enjoyable.

#### 2.1 Simple Event System
- **Priority**: MEDIUM (only if you find direct coupling annoying)
- **Tasks**:
  - Simple event bus with basic pub/sub
  - Common events: `ActorSpawned`, `ActorDestroyed`, `Collision`, `SceneTransition`
  - Keep it lightweight - no need for priorities or filtering
  - Add only when you actually need it

#### 2.2 Better Asset Management
- **Priority**: MEDIUM (if memory becomes an issue)
- **Tasks**:
  - Asset hot-reloading for development (super useful!)
  - Better asset caching (reference counting if needed)
  - Asset validation (catch missing assets early)
  - Keep it simple - don't over-engineer

#### 2.3 Enhanced Input System
- **Priority**: MEDIUM (if you want gamepad support)
- **Tasks**:
  - Gamepad/controller support (if you use controllers)
  - Input action mapping (abstract actions like "jump", "move_left")
  - Key rebinding (useful for different games)
  - Keep keyboard support simple - it works fine as-is

#### 2.4 Project Template Generator
- **Priority**: HIGH (makes starting new games super fast)
- **Tasks**:
  - Simple script/tool to generate new game project
  - Template with basic structure: scenes, player, enemies
  - Copy assets structure
  - Minimal config setup
  - This is probably the most useful feature!

**Deliverables**:
- ✅ Simple event system (if needed)
- ✅ Asset hot-reloading for development
- ✅ Gamepad support (if needed)
- ✅ Project template generator

---

### 🎨 **PHASE 3: Nice-to-Have Features (Ongoing)**

**Goal**: Add features as you need them for specific games. Don't build features you won't use.

#### 3.1 Save System
- **Priority**: HIGH (when you need it for a game)
- **Tasks**:
  - Simple JSON-based save files (human-readable, easy to debug)
  - Save/load game state (player position, inventory, etc.)
  - Multiple save slots
  - Don't worry about versioning/encryption unless you need it
  - Keep it simple - just serialize game state to JSON

#### 3.2 Particle System
- **Priority**: LOW (only if you want visual effects)
- **Tasks**:
  - Simple particle emitter (explosions, trails, smoke)
  - Particle pools for performance
  - Basic integration with physics
  - Keep it minimal - don't build a full particle editor

#### 3.3 Collision System Improvements
- **Priority**: MEDIUM (only if current system becomes limiting)
- **Tasks**:
  - Collision layers (if you need them)
  - Collision types (solid, trigger) - simple flags
  - Better collision debugging visualization
  - Only add what you actually need

#### 3.4 Animation Improvements
- **Priority**: LOW (only if current animation is limiting)
- **Tasks**:
  - Animation easing (ease-in/out) - nice but not essential
  - Better animation state management
  - Sprite sheet tools (if you find yourself doing this manually often)
  - Only add if it saves you time

**Philosophy**: Build features when you need them, not "just in case". Each feature should solve a real problem you're facing.

---

### 📚 **PHASE 4: Developer Experience (Ongoing)**

**Goal**: Make it easier for you to work with the codebase and start new games.

#### 4.1 Documentation (Just Enough)
- **Priority**: MEDIUM (document as you work)
- **Tasks**:
  - Better README with quick start guide
  - Document common patterns you use
  - Inline comments for complex code
  - Architecture overview (high-level, not detailed)
  - Troubleshooting notes (things that tripped you up)
  - Keep it practical - don't write docs you won't read

#### 4.2 Example Game Templates
- **Priority**: HIGH (super useful for starting new games)
- **Tasks**:
  - Platformer template (current game)
  - Top-down template (if you make one)
  - Minimal template (just the basics)
  - Each template is a starting point for a new game
  - Well-commented so you remember how things work

#### 4.3 Development Tools
- **Priority**: MEDIUM (build when they save time)
- **Tasks**:
  - Simple project generator script (bash/go script)
  - Debug visualization (F1 already works, maybe expand it)
  - Asset validation script (catch missing assets)
  - Build script (optional, if you find yourself doing repetitive builds)
  - Keep tools simple - bash scripts are fine

#### 4.4 Testing Utilities
- **Priority**: LOW (only if you write many tests)
- **Tasks**:
  - Test helpers for common scenarios (only if needed)
  - Keep it minimal - Go's testing package is usually enough

**Philosophy**: Build tools that save you time. If a task is annoying, automate it. Otherwise, keep it manual.

---

### 🚀 **PHASE 5: Polish & Optimization (As Needed)**

**Goal**: Make the codebase fast and maintainable for your use. Skip anything unnecessary.

#### 5.1 Performance Optimization
- **Priority**: MEDIUM (only when performance is an issue)
- **Tasks**:
  - Profile when games feel slow
  - Optimize hot paths (rendering, physics)
  - Memory optimization (if you see memory issues)
  - Don't optimize prematurely - measure first

#### 5.2 Code Organization
- **Priority**: LOW (refactor when code becomes hard to navigate)
- **Tasks**:
  - Refactor when you find yourself confused
  - Extract common patterns
  - Keep package structure logical
  - Don't over-engineer - simple is better

#### 5.3 Git Workflow (Optional)
- **Priority**: LOW (if you work on multiple games/branches)
- **Tasks**:
  - Simple branching strategy (if needed)
  - Tags for releases (if you version your games)
  - Keep it simple - you're the only developer

**Philosophy**: Only do this if it helps. Don't add complexity for its own sake.

---

## Technical Debt from TODO.md

Prioritize these as you encounter issues:

1. **Fix collision gaps** (horizontal collision when walking left) - HIGH (it's a bug)
2. **Sprite animation improvements** (jumping states) - MEDIUM (visual polish)
3. **Refactor `shape` to `Shape`** - LOW (only if naming causes confusion)
4. **Create shape package** - LOW (only if physics package gets too large)
5. **Reduce scene duplication** - MEDIUM (if it becomes hard to maintain)
6. **Add camera movement to sequences** - MEDIUM (if you need it)

**Approach**: Fix bugs as you find them. Refactor when code becomes hard to work with. Don't refactor just for the sake of it.

---

## Key Decisions (Simplified)

### Decision 1: Keep `internal/` Structure
**Answer**: Yes, keep it. It works fine for personal use. No need to make it public.

### Decision 2: Breaking Changes
**Answer**: Do what makes sense. You're the only user, so breaking changes are fine. Just document major changes so future-you remembers.

### Decision 3: Dependency Management
**Answer**: Pin to Ebiten versions that work. Upgrade when needed. Don't overthink it.

### Decision 4: Configuration Format
**Answer**: JSON files are fine. Simple and human-readable. Optional Go code overrides.

---

## Success Criteria (Simplified)

### Phase 1 (Foundation)
- [ ] Configuration system works
- [ ] Critical bugs fixed
- [ ] Code is easy to navigate and modify
- [ ] Enough tests to catch regressions

### Phase 2 (Features)
- [ ] Project template generator works
- [ ] Asset hot-reloading in development
- [ ] Features you actually use are implemented

### Phase 3 (Nice-to-Haves)
- [ ] Save system works (when you need it)
- [ ] Other features added as needed

### Phase 4 (DX)
- [ ] Starting a new game is fast (<10 minutes)
- [ ] Code is documented enough for future-you
- [ ] Development workflow is smooth

### Phase 5 (Polish)
- [ ] Performance is acceptable
- [ ] Code is maintainable
- [ ] You can quickly iterate on games

---

## Risk Assessment (Minimal)

### Low Risk (Personal Use)
1. **Breaking Changes**: Fine - you're the only user
2. **Performance**: Profile when it matters, optimize hot paths
3. **Documentation**: Document as you go, enough for future-you
4. **Scope Creep**: Build features as you need them, not "just in case"
5. **Adoption**: Not relevant - it's for personal use

---

## Next Immediate Steps (This Week)

1. **Fix Critical Bugs**: Horizontal collision gap, sprite animation for jumping
2. **Configuration System**: JSON config support (makes starting new games easier)
3. **Better Documentation**: Add inline comments where code is unclear
4. **Project Template**: Script to generate new game project (most valuable!)

---

## Resources Needed (Realistic)

### Development Time
- **Phase 1**: ~20-40 hours (work on it when you have time)
- **Phase 2**: ~40-60 hours (as you build games)
- **Phase 3**: Ongoing (build features as needed)
- **Phase 4**: Ongoing (improve DX as you work)
- **Phase 5**: As needed (optimize when performance matters)

**Total**: Ongoing project, not a fixed timeline. Work on it when it helps.

### External Tools/Libraries
- Testing: Go's built-in testing (usually enough)
- Documentation: Comments in code, README
- Performance: `go test -bench`, `pprof` (when needed)
- CI/CD: Optional (not necessary for personal use)

---

## Conclusion

Firefly is already a solid starter kit for building 2D games. The clean architecture and modular design make it easy to iterate on games. The focus should be on:

1. **Making it convenient** (configuration, project templates)
2. **Fixing bugs** (as you encounter them)
3. **Adding features** (when you need them)
4. **Keeping it simple** (don't over-engineer)

This is an ongoing project - improve it as you build games. Don't try to build everything upfront. Focus on features that save you time when starting new games.

**Remember**: It's a starter kit, not a production framework. Practical over perfect.

---

**Last Updated**: 2025-01-27  
**Next Review**: As you work on it

