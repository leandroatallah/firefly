---
name: fixed-point-physics
description: Testing strategies for fixed-point arithmetic and physics collision edge cases.
---

# Fixed-Point Arithmetic & Physics Testing

Validate positions using `fp16` conversions:

```go
x := fp16.From16(rawX)
y := fp16.From16(rawY)
```

Always use `fp16.To16()` when asserting `x16`/`y16` values in tests.

## Collision Edge Cases to Cover

- One pixel before collision
- Partial overlap
- Full overlap
- Multiple collidables in one space
- Fast movement (skipping over thin walls)
