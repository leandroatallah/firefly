# Mock Specifications — 029 Weapon Item Collectible

## Shared Mocks (Already Available)

Located in `internal/engine/mocks/`:

### MockInventory
```go
// From internal/engine/mocks/combat.go
type MockInventory struct {
    AddWeaponFunc    func(weapon combat.Weapon)
    ActiveWeaponFunc func() combat.Weapon
    SwitchNextFunc   func()
    SwitchPrevFunc   func()
    SwitchToFunc     func(index int)
    HasAmmoFunc      func(weaponID string) bool
    ConsumeAmmoFunc  func(weaponID string, amount int)
    SetAmmoFunc      func(weaponID string, amount int)
    UpdateFunc       func()
}
```

### MockWeapon
```go
// From internal/engine/mocks/combat.go
type MockWeapon struct {
    IDFunc          func() string
    FireFunc        func(x16, y16 int, faceDir animation.FacingDirectionEnum, direction body.ShootDirection)
    CanFireFunc     func() bool
    UpdateFunc      func()
    CooldownFunc    func() int
    SetCooldownFunc func(frames int)
}
```

### MockActor
```go
// From internal/engine/mocks/actors.go
type MockActor struct {
    Id             string
    Pos            image.Rectangle
    SpeedVal       int
    MaxSpeedVal    int
    HealthVal      int
    MaxHealthVal   int
    // ... other fields
}
```

## Package-Local Mocks (for `item_weapon_cannon_test.go`)

### MockBodiesSpace
Implements `body.BodiesSpace` interface for Update() calls.

```go
type MockBodiesSpace struct {
    AddFunc    func(body contracts.Body)
    UpdateFunc func()
    // ... other methods as needed
}

func (m *MockBodiesSpace) Add(body contracts.Body) {
    if m.AddFunc != nil {
        m.AddFunc(body)
    }
}

func (m *MockBodiesSpace) Update() {
    if m.UpdateFunc != nil {
        m.UpdateFunc()
    }
}
```

### MockAppContext
Wrapper to inject mocked ActorManager.

```go
type MockAppContext struct {
    ActorManagerVal *actors.Manager
}

func (m *MockAppContext) ActorManager() *actors.Manager {
    return m.ActorManagerVal
}
```

## Test Setup Pattern

```go
func setupTestContext(t *testing.T) (*app.AppContext, *mocks.MockInventory, *mocks.MockActor) {
    // Create mocked player with inventory
    mockInventory := &mocks.MockInventory{}
    mockPlayer := &mocks.MockActor{Id: "player"}
    
    // Create app context with mocked actor manager
    ctx := &app.AppContext{
        ActorManager: actors.NewManager(),
        Assets:       os.DirFS("."),
    }
    ctx.ActorManager.RegisterPrimary(mockPlayer)
    
    return ctx, mockInventory, mockPlayer
}
```

## Integration Notes

- **Shared mocks** are used across multiple test packages
- **Package-local mocks** are specific to `item_weapon_cannon_test.go`
- All mocks follow the callback pattern (no `_ = variable` anti-pattern)
- Mocks are formatted with `gofmt` after generation
