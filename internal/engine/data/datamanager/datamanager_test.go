package datamanager

import (
	"testing"
)

func TestNewDataManager(t *testing.T) {
	mgr := NewDataManager()

	if mgr == nil {
		t.Fatal("expected NewDataManager to return non-nil manager")
	}
	if mgr.data == nil {
		t.Error("expected data map to be initialized")
	}
	if len(mgr.data) != 0 {
		t.Errorf("expected empty data map, got %d items", len(mgr.data))
	}
}

func TestAddAndGet(t *testing.T) {
	mgr := NewDataManager()

	testData := []byte("test content")
	mgr.Add("test_key", testData)

	retrieved := mgr.Get("test_key")
	if retrieved == nil {
		t.Fatal("expected Get to return non-nil data")
	}
	if string(retrieved) != "test content" {
		t.Errorf("expected 'test content', got '%s'", string(retrieved))
	}
}

func TestGetNotFound(t *testing.T) {
	mgr := NewDataManager()

	retrieved := mgr.Get("nonexistent_key")
	if retrieved != nil {
		t.Errorf("expected nil for nonexistent key, got %v", retrieved)
	}
}

func TestAddOverwrite(t *testing.T) {
	mgr := NewDataManager()

	mgr.Add("key", []byte("original"))
	mgr.Add("key", []byte("overwritten"))

	retrieved := mgr.Get("key")
	if string(retrieved) != "overwritten" {
		t.Errorf("expected 'overwritten', got '%s'", string(retrieved))
	}
}

func TestMultipleEntries(t *testing.T) {
	mgr := NewDataManager()

	mgr.Add("key1", []byte("value1"))
	mgr.Add("key2", []byte("value2"))
	mgr.Add("key3", []byte("value3"))

	if string(mgr.Get("key1")) != "value1" {
		t.Errorf("expected 'value1' for key1")
	}
	if string(mgr.Get("key2")) != "value2" {
		t.Errorf("expected 'value2' for key2")
	}
	if string(mgr.Get("key3")) != "value3" {
		t.Errorf("expected 'value3' for key3")
	}
}
