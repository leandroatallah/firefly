# Audio Engine

The `audio` package provides a robust, goroutine-based audio management system built on top of Ebitengine's `audio` package. It handles multi-format decoding, volume management, and sophisticated fading mechanics.

## Table of Contents

- [Core Architecture](#core-architecture)
- [Sophisticated Mechanics](#sophisticated-mechanics)
  - [Manual Music Looping](#manual-music-looping)
  - [Asynchronous Fade System](#asynchronous-fade-system)
- [Architectural Decisions](#architectural-decisions)
- [Non-Obvious Edge Cases](#non-obvious-edge-cases)
- [Agent Quick-Reference](#agent-quick-reference)

---

## Core Architecture

The system operates using a centralized `AudioManager` that manages a single `audio.Context` (sampled at 44100Hz).

1.  **AudioItem**: A simple container for raw audio bytes, used during the loading phase to decouple file I/O from player initialization.
2.  **AudioManager**: The primary orchestrator. It maintains a registry of `audio.Player` instances and manages global volume state.
3.  **Loader**: A utility module that scans the filesystem/embed FS for supported formats (`.mp3`, `.ogg`, `.wav`) and populates the manager.

---

## Sophisticated Mechanics

### Manual Music Looping

Instead of relying on format-specific looping (like OGG metadata), the engine implements a **Goroutine-Managed Loop**.

- **Mechanism**: When `PlayMusic` is called with `loop: true`, a dedicated goroutine is spawned. It monitors the player's `IsPlaying()` status every 100ms. When the track ends, it automatically rewinds and restarts the player.
- **Interruption**: The loop goroutine is "fade-aware." If a fade-out starts for that specific track, the goroutine detects the entry in `fadeCancel` and terminates itself to prevent a "restart-during-fade" conflict.

### Asynchronous Fade System

Fading is implemented using Go's concurrency primitives rather than per-frame updates in the main game loop.

- **Context-Controlled**: Each fade operation creates a `context.WithCancel`. The cancellation function is stored in the `fadeCancel` map.
- **Interruption Logic**: If you start playing a song that is currently fading out, or start a new fade-out on a track already fading, the previous operation is **immediately cancelled** via its `context.CancelFunc` to ensure volume state consistency.
- **Ticker-Based**: Fades use a `time.Ticker` (100ms resolution) to interpolate volume linearly. This keeps the logic independent of the game's UPS (Updates Per Second).

---

## Architectural Decisions

- **Single Context**: All audio shares one `audio.Context`. This is an Ebitengine requirement to avoid resource exhaustion and synchronization issues.
- **Sample Rate (44100)**: Hardcoded for consistency across all decoded formats.
- **Lazy Player Initialization**: Players are created once during the `Add` phase and reused. This minimizes runtime allocation but means all audio assets are kept in memory as decoded seekers or raw buffers.

---

## Non-Obvious Edge Cases

### Volume Restoration after Fade
When `FadeOut` or `FadeOutAll` completes, the affected players have their volume **restored** to the global `AudioManager.volume` level before being paused. 
- **Why**: This ensures that the next time `Play()` is called on that player, it doesn't start at zero volume, which is a common bug in simpler fade implementations.

### The `_all` Fade Key
The `fadeCancel` map uses a special reserved key `"_all"`. 
- `FadeOutAll` sets this key. 
- Individual `PlayMusic` calls check for this key to ensure that a global fade-out-all doesn't get "overruled" by a single track trying to restart unless explicitly intended.

### Decoding Errors
The `Add` method logs errors but does not panic if a file is corrupted. This allows the game to continue even if a specific asset fails to load, though the `audioPlayers` map will simply lack that entry.

---

## Agent Quick-Reference

### Implementation Checklist

- **Supported Formats**: `.mp3`, `.ogg`, `.wav`. Anything else will be ignored by `LoadAudioAssetsFromFS`.
- **Concurrency**: `AudioManager` is generally thread-safe for playback, but goroutines are heavily used for fades and loops.
- **Global Toggle**: `config.NoSound` completely disables playback and forces `NewAudioManager` to initialize with 0.0 volume.

### Common Pitfalls

- **Looping SFX**: `PlaySound` does **not** support looping. Use `PlayMusic` for any track that needs to loop, even if it's technically a sound effect.
- **Pathing**: The `Loader` expects a specific directory structure: `assets/audio`. Files outside this or in subdirectories are currently ignored by the auto-loader.
- **Memory**: Since everything is decoded and kept in memory, very large uncompressed WAV files can impact the heap. Prefer compressed OGG for long background tracks.
