# ADR-007: FP16 Scale Factor is 16, Not 65536

**Status:** Accepted  
**Date:** 2026-04-06  
**Context:** Projectile system implementation (US-028)

## Problem

The naming "fp16" suggests a 16-bit fixed-point system with 65536 units per pixel (2^16), but the actual implementation uses a scale factor of 16.

## Investigation

During projectile implementation, we discovered:

1. `internal/engine/utils/fp16/fp16.go` defines `const scale = 16`
2. `fp16.From16(value)` converts by dividing by 16, not 65536
3. `fp16.To16(value)` converts by multiplying by 16, not 65536
4. Player position x16=768 = 48 pixels (768/16), not 0.01 pixels (768/65536)

## Decision

**The coordinate system uses scale factor 16:**
- 1 pixel = 16 units
- Position shifts use `<<4` (multiply by 16), not `<<16`
- Velocity for 6 pixels/frame = 96 units/frame (6 * 16)
- Bounds checks use `w<<4` for width in pixels

## Consequences

### Positive
- Smaller integer values (easier to debug)
- Less risk of overflow with 32-bit integers
- Sufficient precision for pixel-perfect platformer physics

### Negative
- Misleading name "fp16" suggests 16-bit fixed-point (65536 scale)
- Easy to confuse with actual 16.16 fixed-point arithmetic

## Implementation Notes

When working with positions and velocities:

```go
// Correct: scale factor 16
x16 := 768              // 48 pixels
velocity := 96          // 6 pixels/frame
bounds := width << 4    // Convert pixels to fp16

// Wrong: assuming scale 65536
x16 := 3145728          // Would be 48 pixels in true fp16
velocity := 393216      // Would be 6 pixels/frame in true fp16
bounds := width << 16   // Wrong shift amount
```

## Related

- ADR-001: FP16 Fixed-Point Arithmetic (describes the system but doesn't clarify scale)
- US-028: Dual weapon setup (where this was discovered)
