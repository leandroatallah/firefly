package events

const (
	// CharacterDiedEventType is the event type for when a character dies.
	CharacterDiedEventType = "CharacterDied"
)

// CharacterDiedEvent is dispatched when a character dies.
type CharacterDiedEvent struct{}

func (e *CharacterDiedEvent) Type() string {
	return CharacterDiedEventType
}
