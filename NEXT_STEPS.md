# Firefly Starter Kit: Immediate Next Steps

This document provides actionable next steps to improve Firefly as a personal-use mini-framework/starter kit for rapidly building 2D games. For the complete roadmap, see [FRAMEWORK_ROADMAP.md](./FRAMEWORK_ROADMAP.md).

## Quick Assessment

### Current Status: ⭐⭐⭐⭐☆ (4/5)
- **Architecture**: Excellent (clean separation, modular design)
- **Functionality**: Great (core systems work well)
- **Usability**: Good (easy to iterate on games)
- **Testing**: Basic (enough for confidence, but minimal)
- **Documentation**: Basic (READMEs, could use more inline docs)

### Starter Kit Maturity Score: ~70%
- Core systems: ✅ 90% (systems exist and work)
- Testing: ⚠️ 20% (enough to catch regressions)
- Documentation: ⚠️ 40% (works, but could be clearer)
- Developer Experience: ✅ 60% (good, but can improve)
- Quick Start: ⚠️ 40% (needs project template generator)

---

## Immediate Actions (This Week)

### 1. High: Fix Critical Bugs ✅ HIGH PRIORITY
**Impact**: Fixes annoying bugs that affect gameplay  
**Effort**: 1-2 days

**Tasks**:
- [ ] Fix horizontal collision gap when walking left
- [ ] Fix sprite animation for jumping
- [ ] Test the fixes thoroughly
- [ ] Add a simple regression test if needed

**Files to check**:
```bash
# Check collision code
grep -r "collision\|Collision" internal/engine/physics/

# Check animation states
grep -r "Jumping\|Falling" internal/engine/entity/actors/
```

### 2. High: Configuration System ✅ HIGH PRIORITY
**Impact**: Makes starting new games much easier  
**Effort**: 2-3 days

**Tasks**:
- [ ] Create JSON config file format (simple structure)
- [ ] Load config from `config.json` in project root
- [ ] Support overriding via code (for programmatic setup)
- [ ] Keep hard-coded defaults (they're fine as fallback)
- [ ] Document config options in README

**Example config.json**:
```json
{
  "screen_width": 320,
  "screen_height": 180,
  "physics": {
    "jump_force": 6,
    "gravity": 6
  }
}
```

### 3. High: Project Template Generator ✅ HIGH PRIORITY
**Impact**: Most valuable feature - starts new games in minutes  
**Effort**: 2-3 days

**Tasks**:
- [ ] Create simple script (`scripts/new-game.sh` or Go script)
- [ ] Template structure:
  - [ ] Basic scene setup
  - [ ] Player actor
  - [ ] Basic config.json
  - [ ] Assets folder structure
  - [ ] README template
- [ ] Copy current game as reference template
- [ ] Make it easy to customize

**Usage**:
```bash
./scripts/new-game.sh my-new-game
# Creates new game with basic structure
```

### 4. Medium: Better Inline Documentation ✅ MEDIUM PRIORITY
**Impact**: Makes code easier to understand later  
**Effort**: Ongoing (as you work)

**Tasks**:
- [ ] Add package-level comments (what each package does)
- [ ] Document complex functions (especially physics)
- [ ] Add examples in comments where helpful
- [ ] Focus on "why" not "what"
- [ ] Don't over-document - code should be readable

---

## Short-term Goals (Next 2 Weeks)

### Week 1 Focus: Fix & Improve
1. ✅ Fix critical bugs (collision, animation)
2. ✅ Configuration system (JSON support)
3. ✅ Better inline documentation (as you work)
4. ✅ Project template generator (start)

### Week 2 Focus: Developer Experience
1. ✅ Complete project template generator
2. ✅ Test template with a new game
3. ✅ Asset hot-reloading (development mode)
4. ⚠️ Save system (if you need it for a game)

---

## Technical Debt to Address (From TODO.md)

### Priority 1: Critical Bugs (Fix These)
1. **Fix horizontal collision gap** (walking left)
   - Impact: Game-breaking bug
   - Fix when you encounter it in gameplay

2. **Fix sprite animation for jumping**
   - Impact: Visual bug
   - Fix when you notice it's wrong

### Priority 2: Nice to Have (Fix When Needed)
1. **Refactor `shape` to `Shape`** - Only if naming causes confusion
2. **Create `shape` package** - Only if physics package gets too large
3. **Reduce scene duplication** - Only if it becomes hard to maintain
4. **Add camera movement to sequences** - Only if you need it for a game

---

## Quick Wins (Can Do Now)

### 1. Add Package Documentation (30 minutes)
Add package comments as you work on each package. Don't do it all at once - just add comments when code is unclear.

### 2. Fix log.Fatal Where It Matters (1 hour)
Find `log.Fatal` calls that happen in normal operation (not startup errors):
```bash
grep -rn "log.Fatal" internal/engine/ | grep -v "init\|setup\|New"
```
Replace with error returns where it makes sense. Keep `log.Fatal` for truly unrecoverable errors.

### 3. Create Project Template (2-3 hours)
Copy your current game structure, remove game-specific content, create a minimal template:
- Basic scene
- Player actor
- Config.json
- README

### 4. Add Asset Hot-Reload Script (1 hour)
Simple script that watches assets and restarts game (development only):
```bash
# scripts/dev-reload.sh (example)
while inotifywait -e modify assets/; do
  killall firefly
  go run main.go &
done
```

---

## Success Metrics (Realistic)

### Phase 1 (Weeks 1-2) Targets:
- [ ] Critical bugs fixed (collision, animation)
- [ ] Configuration system works (JSON support)
- [ ] Project template generator functional
- [ ] Code is easier to navigate (better docs)
- [ ] Can start new game in <10 minutes

### Quick Health Check:
```bash
# Test coverage (not aiming for high, just check)
go test ./internal/engine/... -cover

# Count test files (should grow as you fix bugs)
find internal/engine -name "*_test.go" | wc -l

# Check for obvious bugs (grep for common issues)
grep -r "TODO\|FIXME\|XXX" internal/engine/
```

---

## Decision Points (Already Made - Keep It Simple)

### 1. Package Structure
**Decision**: Keep `internal/` - it works fine for personal use. No need to make it public.

### 2. Versioning Strategy
**Decision**: Version games, not the framework. Use git tags if helpful, but don't overthink it.

### 3. Breaking Changes
**Decision**: Do what makes sense. You're the only user, so breaking changes are fine. Just note major changes so future-you remembers.

### 4. Testing Strategy
**Decision**: Test when you fix bugs (regression tests). Don't aim for high coverage - just enough confidence.

### 5. Documentation Strategy
**Decision**: Document as you work. Add comments when code is unclear. Keep it practical.

---

## Resources & Tools (Keep It Simple)

### Testing Tools
- Go's built-in `testing` package (usually enough)
- `testify` (optional, if you find yourself repeating assertions)

### Documentation Tools
- Inline comments (most important)
- README (update as needed)
- `godoc` (if you want to generate docs, optional)

### Performance Tools
- `go test -bench` - When performance matters
- `pprof` - When profiling is needed
- Don't use tools you don't need

---

## Getting Started Right Now

### Step 1: Fix Critical Bugs (1-2 hours)
```bash
# Find collision-related code
grep -rn "collision\|Collision" internal/engine/physics/

# Test the bug in game, then fix it
# Add a simple regression test
```

### Step 2: Create Project Template (2-3 hours)
```bash
# Create a clean template from current game
mkdir templates/platformer
# Copy basic structure, remove game-specific content
# Create script to generate new game from template
```

### Step 3: Configuration System (2-3 hours)
- Create `config.json` structure
- Add config loading code
- Test with current game

### Step 4: Better Docs (Ongoing)
- As you work, add comments where code is unclear
- Update README when you add features

---

## Philosophy: Practical Over Perfect

Remember:
- **Build features when you need them**, not "just in case"
- **Fix bugs as you find them**, not all at once
- **Document as you work**, not comprehensively upfront
- **Test what matters**, not everything
- **Keep it simple** - complexity slows you down
- **Iterate fast** - that's the whole point

This is a starter kit to help you build games quickly. Don't turn it into a production framework unless that's what you want.

---

**Next Review**: After completing Week 1 tasks

**Quick Command Reference**:
```bash
# Run all tests
go test ./internal/engine/... -v

# Test coverage
go test ./internal/engine/... -coverprofile=coverage.out
go tool cover -func=coverage.out

# Find TODO items
grep -rn "TODO\|FIXME" internal/engine/

# Check for log.Fatal
grep -rn "log.Fatal\|log.Panic" internal/engine/

# Count test files
find internal/engine -name "*_test.go" | wc -l
```

