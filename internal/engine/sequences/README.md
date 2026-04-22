# Sequences

The `sequences` package runs **scripted event chains** — cutscenes, intros, dialogues, scripted actor movements, and any other frame-by-frame choreography that doesn't fit into a per-frame `Update`. Sequences are defined as JSON, parsed into typed commands, and driven by a `SequencePlayer`.

## Core Types

- `Sequence` — an ordered list of `sequences.Command` (contract in `contracts/sequences/`) plus per-command blocking flags and sequence-level flags (`Interruptible`, `OneTime`, `BlockPlayerMovement`).
- `SequencePlayer` — the executor. Holds the active sequence, tracks the current command index, and keeps a background-command list for non-blocking commands that should keep ticking while the timeline advances.
- `CommandData` / `SequenceData` — JSON parse wrappers. `CommandData.ToCommand()` discriminates on the `"command"` field and returns the concrete `Command` implementation.

## Command Families

Commands live in separate files so they can be read in isolation:

| File | Commands |
|---|---|
| `commands.go` | `DialogueCommand`, `DelayCommand`, `EventCommand` |
| `commands_actor.go` | Actor movement, following, speed overrides |
| `commands_camera.go` | Camera zoom / move / reset / shake, vignette |
| `commands_music.go` | Background music play / stop / fade |
| `commands_vfx.go` | Floating / overhead / screen text, particle bursts |
| `commands_sequence.go` | Nested `call_sequence` (chained execution) |

Every command implements `Init(appContext)` (wiring) and `Update() bool` (returns `true` when finished).

## Blocking Model

Each command has a `BlockSequence` flag (`block_sequence` in JSON):

- **Blocking** (default, or `true`): the player waits for `Update()` to return `true` before advancing to the next command.
- **Non-blocking** (`false`): the command is kicked off, then moved to `backgroundCommands` so the next command starts immediately. Background commands keep ticking until they report done.

This means a single sequence can e.g. start a screen shake + music fade in parallel, then block on a dialogue line.

## One-Shot & Interruptible

- `OneTime: true` — the `SequencePlayer` records the sequence's path in `consumedOneTimeSequences` on first `Play` and silently ignores future calls.
- `Interruptible: true` — another `Play` call can replace an in-flight sequence. Non-interruptible sequences are sticky until their timeline completes.

## Usage

```go
player := sequences.NewSequencePlayer(appContext)

// Play directly from a JSON file in the FS (assets or embed.FS).
player.PlaySequence("sequences/intro.json")

// Or build a Sequence programmatically and call Play.
seq, _ := sequences.NewSequenceFromFS(assets, "sequences/intro.json")
player.Play(seq)

// Tick the player every frame.
player.Update()
```

## Testing

Every command family has a matching `*_test.go` file with table-driven coverage. Tests drive `Update()` directly and assert on side effects through mock app contexts defined in `mocks_test.go` / `sequences_test.go`. No Ebitengine window is required.

## Example

See `example_usage.go` for a runnable reference implementation that wires a `SequencePlayer` into a minimal app context.
