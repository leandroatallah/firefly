# ADR-001 — FP16 Fixed-Point Arithmetic for Positions

## Status
Accepted

## Context
Game physics requires sub-pixel precision for smooth movement and accurate collision detection. Using `float64` introduces floating-point rounding errors that accumulate over time and produce non-deterministic collision results. Integer arithmetic is deterministic and fast, but whole-pixel resolution is too coarse for smooth platformer movement.

## Decision
All body positions are stored as `x16`/`y16` integers scaled by a factor of 16 (one pixel = 16 units). The `fp16` package provides two helpers: `To16(v int) int` (multiply by 16) and `From16(v int) int` (divide by 16). `Body.SetPosition` converts to x16 on write; `Body.Position()` converts back to pixels on read.

## Consequences
- Collision and movement math is fully deterministic and integer-only.
- Sub-pixel movement is representable in 1/16-pixel increments.
- All position comparisons and physics calculations must use x16 values internally; callers that work in pixel space use `SetPosition`/`Position()` and remain unaware of the scaling.
- Positions passed directly as x16 (e.g. `SetPosition16`) bypass conversion and must be handled carefully to avoid scale mismatches.
