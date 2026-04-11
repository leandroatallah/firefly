# US-033 — Combat Particle Definitions

**Branch:** `033-combat-particle-definitions`
**Bounded Context:** Assets (`assets/particles/`)

## Story

As a game developer using this boilerplate,
I want predefined particle effects for combat actions,
so that muzzle flash, impact, and despawn effects work without creating custom particle configs.

## Context

The current `assets/particles/vfx.json` contains only jump and landing particle definitions. This story adds combat-specific particle definitions using pixel-based particles with 1-bit (black and white) colors.

## Acceptance Criteria

- **AC1** — `muzzle_flash` particle type defined in `vfx.json` using pixel-based config with white flash effect.
- **AC2** — `bullet_impact` particle type defined in `vfx.json` using pixel-based config with white spark effect.
- **AC3** — `bullet_despawn` particle type defined in `vfx.json` using pixel-based config with white fade-out effect.
- **AC4** — All particle configs use pixel-based particle system (no image assets required).
- **AC5** — Particle colors limited to black (#000000) and white (#FFFFFF) for 1-bit aesthetic.
- **AC6** — Particle effects have appropriate duration and behavior for combat feedback.
- **AC7** — JSON validation passes for updated `vfx.json`.

## Notes

- File path: `assets/particles/vfx.json`
- Uses existing pixel-based particle system
- 1-bit color palette: black and white only
- Should be implemented FIRST as foundation for other VFX stories
