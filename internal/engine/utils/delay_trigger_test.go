package utils

import "testing"

func TestDelayTrigger_EnableAndTrigger(t *testing.T) {
	var trigger DelayTrigger

	// Initially disabled
	if trigger.IsEnabled() {
		t.Fatal("expected trigger to be disabled initially")
	}

	// Enable with delay
	trigger.Enable(3)

	if !trigger.IsEnabled() {
		t.Fatal("expected trigger to be enabled after Enable()")
	}
	if trigger.IsTriggered() {
		t.Fatal("expected trigger to not be triggered yet")
	}
	if trigger.IsReady() {
		t.Fatal("expected trigger to not be ready while delay > 0")
	}

	// Update should decrement delay
	trigger.Update()
	if trigger.delay != 2 {
		t.Fatalf("expected delay=2, got %d", trigger.delay)
	}

	trigger.Update()
	if trigger.delay != 1 {
		t.Fatalf("expected delay=1, got %d", trigger.delay)
	}

	trigger.Update()
	if trigger.delay != 0 {
		t.Fatalf("expected delay=0, got %d", trigger.delay)
	}

	if !trigger.IsReady() {
		t.Fatal("expected trigger to be ready when delay==0")
	}

	// Trigger should return true only once
	if !trigger.Trigger() {
		t.Fatal("expected Trigger() to return true when ready")
	}
	if !trigger.IsTriggered() {
		t.Fatal("expected trigger to be marked as triggered")
	}

	// Subsequent Trigger() calls should return false
	if trigger.Trigger() {
		t.Fatal("expected Trigger() to return false after already triggered")
	}
}

func TestDelayTrigger_EnableWithZeroDelay(t *testing.T) {
	var trigger DelayTrigger

	trigger.Enable(0)

	if !trigger.IsReady() {
		t.Fatal("expected trigger to be ready immediately with delay=0")
	}

	if !trigger.Trigger() {
		t.Fatal("expected Trigger() to return true with delay=0")
	}
}

func TestDelayTrigger_UpdateDoesNotGoNegative(t *testing.T) {
	var trigger DelayTrigger

	trigger.Enable(1)
	trigger.Update() // delay = 0
	trigger.Update() // should stay at 0

	if trigger.delay != 0 {
		t.Fatalf("expected delay to stay at 0, got %d", trigger.delay)
	}
}

func TestDelayTrigger_UpdateAfterTriggered(t *testing.T) {
	var trigger DelayTrigger

	trigger.Enable(1)
	trigger.Update() // delay = 0
	trigger.Trigger()
	trigger.Update() // should not change anything

	if !trigger.IsTriggered() {
		t.Fatal("expected trigger to remain triggered")
	}
}

func TestDelayTrigger_Reset(t *testing.T) {
	var trigger DelayTrigger

	trigger.Enable(2)
	trigger.Update() // delay = 1
	trigger.Update() // delay = 0
	trigger.Trigger()

	if !trigger.IsTriggered() {
		t.Fatal("expected trigger to be triggered before reset")
	}

	trigger.Reset()

	if trigger.IsEnabled() {
		t.Fatal("expected trigger to be disabled after reset")
	}
	if trigger.IsTriggered() {
		t.Fatal("expected trigger to not be triggered after reset")
	}
	if trigger.IsReady() {
		t.Fatal("expected trigger to not be ready after reset")
	}
}

func TestDelayTrigger_MultipleUpdatesBeforeTrigger(t *testing.T) {
	var trigger DelayTrigger

	trigger.Enable(5)

	// Update more times than delay
	for i := 0; i < 10; i++ {
		trigger.Update()
	}

	if !trigger.IsReady() {
		t.Fatal("expected trigger to be ready after enough updates")
	}

	if !trigger.Trigger() {
		t.Fatal("expected Trigger() to return true")
	}
}
