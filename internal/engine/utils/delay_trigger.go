package utils

// DelayTrigger is a reusable trigger with countdown and single-fire behavior.
type DelayTrigger struct {
	enabled   bool
	triggered bool
	delay     int
}

// Enable activates the trigger with a delay in frames.
func (t *DelayTrigger) Enable(delay int) {
	t.enabled = true
	t.triggered = false
	t.delay = delay
}

// Update decrements the delay if enabled. Should be called every frame.
func (t *DelayTrigger) Update() {
	if t.enabled && !t.triggered && t.delay > 0 {
		t.delay--
	}
}

// Trigger checks if action should occur. Returns true only once when ready.
func (t *DelayTrigger) Trigger() bool {
	if t.enabled && !t.triggered && t.delay == 0 {
		t.triggered = true
		return true
	}
	return false
}

// Reset clears the trigger state for reuse.
func (t *DelayTrigger) Reset() {
	t.enabled = false
	t.triggered = false
	t.delay = 0
}

// IsEnabled returns true if the trigger is enabled.
func (t *DelayTrigger) IsEnabled() bool {
	return t.enabled
}

// IsTriggered returns true if the trigger has already fired.
func (t *DelayTrigger) IsTriggered() bool {
	return t.triggered
}

// IsReady returns true if the trigger is enabled and the delay has finished.
func (t *DelayTrigger) IsReady() bool {
	return t.enabled && t.delay == 0
}
