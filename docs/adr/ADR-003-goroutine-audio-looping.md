# ADR-003 — Goroutine-Based Audio Looping

## Status
Accepted

## Context
Ebitengine's `audio.Player` does not natively support seamless looping. The built-in approach requires wrapping the stream in an `audio.InfiniteLoop` reader, which ties the loop to the stream decoder and makes it difficult to interrupt (e.g. for fade-out or pause). A clean pause/resume/fade-out cycle needs the ability to stop the loop at any point without restarting the decoder.

## Decision
When `PlayMusic` is called with `loop = true`, a goroutine is spawned that polls `player.IsPlaying()` every 100 ms. When the player stops (track ended) and is not paused or fading, the goroutine calls `player.Rewind()` and `player.Play()` to restart. The goroutine exits when a `context.CancelFunc` stored in `fadeCancel` is triggered, enabling fade-out and stop operations to cleanly terminate the loop.

## Consequences
- Looping is interruptible at any time via context cancellation.
- Pause and fade-out integrate naturally: the goroutine checks `paused` and `fadeCancel` flags before rewinding.
- There is a ~100 ms polling interval, which introduces a small gap at the loop point. This is acceptable for background music but not for sample-accurate loops.
- Each looping track owns one goroutine for its lifetime; goroutine leaks are prevented by always cancelling via `fadeCancel` before starting a new loop on the same track.
