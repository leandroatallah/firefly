package event

import "testing"

type testEvt struct{ k string }

func (e testEvt) Type() string { return e.k }

func TestEventSubscribeAndPublish(t *testing.T) {
	m := NewManager()

	count := 0
	m.Subscribe("X", func(e Event) {
		count++
		if e.Type() != "X" {
			t.Fatalf("wrong type")
		}
	})

	m.Publish(testEvt{"X"})
	m.Publish(testEvt{"Y"})
	m.Publish(testEvt{"X"})

	if count != 2 {
		t.Fatalf("expected 2 invocations, got %d", count)
	}
}

func TestGenericEvent(t *testing.T) {
	e := GenericEvent{EventType: "A", Payload: map[string]interface{}{"v": 1}}
	if e.Type() != "A" || e.Payload["v"].(int) != 1 {
		t.Fatalf("unexpected generic event")
	}
}

func TestUnsubscribe(t *testing.T) {
	m := NewManager()

	count := 0
	unsubscribe := m.Subscribe("X", func(e Event) {
		count++
	})

	m.Publish(testEvt{"X"})
	if count != 1 {
		t.Fatalf("expected 1 invocation before unsubscribe, got %d", count)
	}

	unsubscribe()

	m.Publish(testEvt{"X"})
	if count != 1 {
		t.Fatalf("expected 1 invocation after unsubscribe, got %d", count)
	}
}

func TestUnsubscribeMultipleListeners(t *testing.T) {
	m := NewManager()

	count1, count2 := 0, 0
	unsub1 := m.Subscribe("X", func(e Event) { count1++ })
	unsub2 := m.Subscribe("X", func(e Event) { count2++ })

	m.Publish(testEvt{"X"})
	if count1 != 1 || count2 != 1 {
		t.Fatalf("expected both listeners to be called, got %d, %d", count1, count2)
	}

	unsub1()

	m.Publish(testEvt{"X"})
	if count1 != 1 || count2 != 2 {
		t.Fatalf("expected only second listener after unsub1, got %d, %d", count1, count2)
	}

	unsub2()

	m.Publish(testEvt{"X"})
	if count1 != 1 || count2 != 2 {
		t.Fatalf("expected no listeners after unsub2, got %d, %d", count1, count2)
	}
}

func TestUnsubscribeIdempotent(t *testing.T) {
	m := NewManager()

	count := 0
	unsubscribe := m.Subscribe("X", func(e Event) { count++ })

	m.Publish(testEvt{"X"})

	unsubscribe()
	unsubscribe() // Should not panic or cause errors
	unsubscribe()

	m.Publish(testEvt{"X"})
	if count != 1 {
		t.Fatalf("expected 1 invocation, got %d", count)
	}
}

func TestUnsubscribeDifferentEventTypes(t *testing.T) {
	m := NewManager()

	countX, countY := 0, 0
	unsubX := m.Subscribe("X", func(e Event) { countX++ })
	m.Subscribe("Y", func(e Event) { countY++ })

	m.Publish(testEvt{"X"})
	m.Publish(testEvt{"Y"})

	if countX != 1 || countY != 1 {
		t.Fatalf("expected both events handled, got %d, %d", countX, countY)
	}

	unsubX()

	m.Publish(testEvt{"X"})
	m.Publish(testEvt{"Y"})

	if countX != 1 || countY != 2 {
		t.Fatalf("expected only Y handled after unsubX, got %d, %d", countX, countY)
	}
}
