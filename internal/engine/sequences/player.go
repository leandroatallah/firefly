package sequences

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/sequences"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/debug"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// commandName returns the concrete command type name without the package
// prefix (e.g. "*sequences.MoveActorCommand" -> "MoveActorCommand"). Used for
// the command_init debug channel so logged names track struct renames.
func commandName(cmd sequences.Command) string {
	name := fmt.Sprintf("%T", cmd)
	if i := strings.LastIndex(name, "."); i >= 0 {
		name = name[i+1:]
	}
	return name
}

// SequencePlayer manages the execution of a sequence.
type SequencePlayer struct {
	app.AppContextHolder

	currentSequence     sequences.Sequence
	currentSequencePath string
	currentCommandIndex int
	isPlaying           bool

	hasActiveCommands bool
	lastBlockingIndex int
	blockingEnded     bool
	blockedByParent   bool
	debugActive       bool
	debugPaused       bool

	// updating is true while the player is driving command Init/Update calls.
	// A command may synchronously publish an event whose handler calls back into
	// Play (e.g. one sequence chaining to another). Mutating player state from
	// inside that stack corrupts the in-flight advance, so such calls are queued
	// in pendingPlays and drained once the current update/advance unwinds.
	updating     bool
	pendingPlays []sequences.Sequence

	backgroundCommands       []sequences.Command
	consumedOneTimeSequences map[string]struct{}
}

// NewSequencePlayer creates a new player.
func NewSequencePlayer(appContext *app.AppContext) *SequencePlayer {
	ctx := app.AppContextHolder{}
	ctx.SetAppContext(appContext)
	return &SequencePlayer{
		AppContextHolder:         ctx,
		consumedOneTimeSequences: make(map[string]struct{}),
	}
}

// PlaySequence loads and plays a sequence from a JSON file.
func (p *SequencePlayer) PlaySequence(filePath string) {
	debug.Watch("sequence_play", "", filePath)
	sequence, err := NewSequenceFromFS(p.AppContext().Assets, filePath)
	if err != nil {
		return
	}
	p.Play(sequence)
}

// Play starts executing a sequence.
//
// If called re-entrantly (from within a command's Init/Update, e.g. an event
// command whose handler chains to another sequence), the request is queued and
// applied after the current update/advance unwinds, so it never mutates player
// state mid-stack.
func (p *SequencePlayer) Play(sequence sequences.Sequence) {
	if p.updating {
		p.pendingPlays = append(p.pendingPlays, sequence)
		return
	}
	p.play(sequence)
	p.drainPending()
}

// drainPending applies any sequences queued by re-entrant Play calls. play may
// itself queue further sequences (a chain), so this loops until the queue is
// empty. Each queued sequence interrupts the previous one, matching the
// synchronous semantics of back-to-back Play calls.
func (p *SequencePlayer) drainPending() {
	for len(p.pendingPlays) > 0 {
		next := p.pendingPlays[0]
		p.pendingPlays = p.pendingPlays[1:]
		p.play(next)
	}
}

// play performs the actual sequence start. It must only be called when not
// already driving commands (callers: Play and drainPending), so that the
// advanceToNextCommand below — which can trigger re-entrant Play calls — runs
// under the updating guard rather than corrupting an outer advance.
func (p *SequencePlayer) play(sequence sequences.Sequence) {
	sequencePath := sequence.GetPath()

	// Check if this is a one-time sequence that has already been consumed
	if sequence.OneTime() {
		if _, consumed := p.consumedOneTimeSequences[sequencePath]; consumed {
			return
		}
		// Mark as consumed immediately to prevent re-entry during playback
		p.consumedOneTimeSequences[sequencePath] = struct{}{}
	}

	// If same sequence is already playing, don't restart it
	if p.hasActiveCommands && p.currentSequencePath == sequencePath {
		return
	}

	// If a different sequence is requested and current is non-interruptible, skip
	if p.hasActiveCommands && !p.currentSequence.Interruptible() {
		return
	}

	// Stop current sequence if playing
	if p.hasActiveCommands {
		p.Stop()
	}

	p.currentSequence = sequence
	p.currentSequencePath = sequencePath
	p.currentCommandIndex = -1
	p.hasActiveCommands = true
	p.blockingEnded = false
	p.lastBlockingIndex = -1
	p.backgroundCommands = nil

	if seq, ok := sequence.(*Sequence); ok && len(seq.blockSequenceFlags) == len(seq.commands) {
		for i, flag := range seq.blockSequenceFlags {
			if flag {
				p.lastBlockingIndex = i
			}
		}
	} else {
		p.lastBlockingIndex = len(sequence.Commands()) - 1
	}

	p.isPlaying = p.lastBlockingIndex >= 0

	if config.Get().SequenceDebug {
		p.debugActive = true
		p.debugPaused = true
		debug.Log("sequence_debug", "started %q — Enter: next command | Esc: end sequence", sequencePath)
	}

	if seq, ok := sequence.(*Sequence); ok && seq.BlockPlayerMovement && p.lastBlockingIndex >= 0 && !p.blockedByParent {
		if player, found := p.AppContext().ActorManager.GetPlayer(); found {
			player.BlockMovement()
		}
	}

	debug.Log("command_init", "=== START seq=%s (%d commands) ===",
		p.currentSequencePath, len(sequence.Commands()))
	// In debug mode this loads + Init's the first blocking command and leaves
	// debugPaused=true (set in advanceToNextCommand). The command then runs on
	// the next Enter via the normal Update() path.
	p.updating = true
	p.advanceToNextCommand()
	p.updating = false
}

// IsPlaying returns true if a sequence is currently being played.
func (p *SequencePlayer) IsPlaying() bool {
	debug.Watch("sequence_isPlaying", "", p.isPlaying)
	return p.isPlaying
}

func (p *SequencePlayer) IsOver() bool {
	if p.currentSequence == nil {
		return true
	}
	return p.currentCommandIndex >= len(p.currentSequence.Commands())
}

func (p *SequencePlayer) IsDebugPaused() bool {
	return p.debugPaused
}

// Update should be called every frame. It updates the current command.
//
// Command Init/Update runs under the updating guard so that any sequence chained
// from inside a command (via Play) is queued, then applied via drainPending once
// this update unwinds.
func (p *SequencePlayer) Update() {
	if p.updating {
		return
	}
	p.updating = true
	p.update()
	p.updating = false
	p.drainPending()
}

func (p *SequencePlayer) update() {
	if !p.hasActiveCommands {
		return
	}

	ctx := p.AppContext()

	// Update fade overlay
	if ctx != nil && ctx.FadeOverlay != nil {
		ctx.FadeOverlay.Update()
	}

	if ctx != nil && ctx.SolidColorOverlay != nil {
		ctx.SolidColorOverlay.Update()
	}

	if p.debugPaused {
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			// Exit the debugger but keep the sequence running. The currently
			// staged command resumes via the normal Update() path and the rest
			// of the sequence plays out without pausing.
			p.debugActive = false
			p.debugPaused = false
			debug.Log("sequence_debug", "debugger ended — sequence resumes %q", p.currentSequencePath)
			return
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			// Unpause and let the already-Init'd blocking command run via the
			// normal Update() path below. When it completes,
			// advanceToNextCommand re-pauses before the next command. This lets
			// time-based commands (e.g. DelayCommand) actually tick their
			// frames instead of being skipped.
			p.debugPaused = false
		}
		return
	}

	// Update the currently blocking command (if any)
	if p.currentCommandIndex >= 0 && p.currentCommandIndex < len(p.currentSequence.Commands()) {
		currentCommand := p.currentSequence.Commands()[p.currentCommandIndex]
		if currentCommand.Update() {
			p.advanceToNextCommand()
		}
	}

	// Update background (non-blocking) commands
	if len(p.backgroundCommands) > 0 {
		remaining := p.backgroundCommands[:0]
		for _, cmd := range p.backgroundCommands {
			if !cmd.Update() {
				remaining = append(remaining, cmd)
			}
		}
		p.backgroundCommands = remaining
	}

	// End entire sequence only when no more commands and no more background commands
	if p.currentCommandIndex >= len(p.currentSequence.Commands()) && len(p.backgroundCommands) == 0 {
		p.endSequence()
		return
	}
}

// advanceToNextCommand moves to the next command in the queue and initializes it.
func (p *SequencePlayer) advanceToNextCommand() {
	// Keep advancing until we hit a blocking command or the end
	for {
		p.currentCommandIndex++
		if p.currentCommandIndex >= len(p.currentSequence.Commands()) {
			// No more commands in pipeline; background may still be running
			if !p.blockingEnded {
				debug.Log("sequence_player", "status", "END")
				p.endBlockingPhase()
			}
			return
		}

		if !p.blockingEnded && p.currentCommandIndex > p.lastBlockingIndex {
			p.endBlockingPhase()
		}

		nextCommand := p.currentSequence.Commands()[p.currentCommandIndex]

		// Check if this command is blocking or not
		isBlocking := true
		if seq, ok := p.currentSequence.(*Sequence); ok && p.currentCommandIndex < len(seq.blockSequenceFlags) {
			isBlocking = seq.blockSequenceFlags[p.currentCommandIndex]
		}

		debug.Log("command_init", "[%d/%d] %s blocking=%v seq=%s",
			p.currentCommandIndex+1, len(p.currentSequence.Commands()),
			commandName(nextCommand), isBlocking, p.currentSequencePath)
		nextCommand.Init(p.AppContext())
		if isBlocking {
			if p.debugActive {
				p.debugPaused = true
			}
			return
		}

		// Non-blocking: run in background and immediately advance to try next
		p.backgroundCommands = append(p.backgroundCommands, nextCommand)
		// Loop to look for the next blocking command (or finish)
	}
}

func (p *SequencePlayer) endSequence() {
	debug.Log("command_init", "=== END seq=%s ===", p.currentSequencePath)
	// endBlockingPhase must run before currentSequence is cleared: it reads the
	// concrete *Sequence to decide whether to unblock player movement. Nilling
	// first would make that type assertion fail and silently skip the unblock,
	// leaving the player stuck.
	if !p.blockingEnded {
		p.endBlockingPhase()
	}
	p.hasActiveCommands = false
	p.currentSequence = nil
	p.currentSequencePath = ""
}

// Stop cleanly stops the current sequence.
func (p *SequencePlayer) Stop() {
	if !p.hasActiveCommands {
		return
	}

	debug.Log("command_init", "=== END (stopped) seq=%s ===", p.currentSequencePath)

	// Unblock player if needed
	if p.currentSequence != nil && p.currentSequencePath != "" {
		if seq, ok := p.currentSequence.(*Sequence); ok && seq.BlockPlayerMovement && !p.blockedByParent {
			if player, found := p.AppContext().ActorManager.GetPlayer(); found {
				player.UnblockMovement()
			}
		}
	}

	p.hasActiveCommands = false
	p.currentSequence = nil
	p.currentSequencePath = ""
	p.blockingEnded = true
	p.isPlaying = false
	p.backgroundCommands = nil
}

func (p *SequencePlayer) endBlockingPhase() {
	p.blockingEnded = true
	p.isPlaying = false

	if seq, ok := p.currentSequence.(*Sequence); ok && seq.BlockPlayerMovement && !p.blockedByParent {
		if player, found := p.AppContext().ActorManager.GetPlayer(); found {
			player.UnblockMovement()
		}
	}
}

// Draw renders debug overlay when sequence is paused in debug mode.
func (p *SequencePlayer) Draw(screen *ebiten.Image) {
	if !p.debugPaused {
		return
	}

	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()
	overlay := ebiten.NewImage(w, h)
	overlay.Fill(color.RGBA{0, 0, 0, 180})

	opts := &ebiten.DrawImageOptions{}
	screen.DrawImage(overlay, opts)
}

func (p *SequencePlayer) DrawOver(screen *ebiten.Image) {
	ctx := p.AppContext()
	if ctx != nil && ctx.FadeOverlay != nil {
		ctx.FadeOverlay.Draw(screen)
	}
	if ctx != nil && ctx.SolidColorOverlay != nil {
		ctx.SolidColorOverlay.Draw(screen)
	}
}
