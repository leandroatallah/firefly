package items

import (
	"testing"
)

// mockConcreteItem implements Item for testing
type mockConcreteItem struct {
	*BaseItem
}

func TestNewItemFactory(t *testing.T) {
	itemMap := ItemMap[*mockConcreteItem]{
		"test": func(x, y int, id string) *mockConcreteItem {
			return &mockConcreteItem{}
		},
	}

	factory := NewItemFactory(itemMap)
	if factory == nil {
		t.Fatal("NewItemFactory() returned nil")
	}
	if factory.itemMap == nil {
		t.Error("NewItemFactory() did not initialize itemMap")
	}
}

func TestItemFactory_Create_Success(t *testing.T) {
	created := false
	itemMap := ItemMap[*mockConcreteItem]{
		"test": func(x, y int, id string) *mockConcreteItem {
			created = true
			if x != 100 {
				t.Errorf("Create() x = %d, want 100", x)
			}
			if y != 200 {
				t.Errorf("Create() y = %d, want 200", y)
			}
			if id != "test-id" {
				t.Errorf("Create() id = %q, want %q", id, "test-id")
			}
			return &mockConcreteItem{}
		},
	}

	factory := NewItemFactory(itemMap)
	_, err := factory.Create("test", 100, 200, "test-id")

	if err != nil {
		t.Errorf("Create() returned error: %v", err)
	}
	if !created {
		t.Error("Create() did not call item constructor")
	}
}

func TestItemFactory_Create_UnknownType(t *testing.T) {
	itemMap := ItemMap[*mockConcreteItem]{
		"test": func(x, y int, id string) *mockConcreteItem {
			return &mockConcreteItem{}
		},
	}

	factory := NewItemFactory(itemMap)
	item, err := factory.Create("unknown_type", 0, 0, "id")

	if err == nil {
		t.Error("Create() with unknown type should return error")
	}
	if item != nil {
		t.Error("Create() with unknown type should return nil item")
	}
}

func TestItemFactory_Create_MultipleTypes(t *testing.T) {
	itemMap := ItemMap[*mockConcreteItem]{
		"type1": func(x, y int, id string) *mockConcreteItem {
			return &mockConcreteItem{}
		},
		"type2": func(x, y int, id string) *mockConcreteItem {
			return &mockConcreteItem{}
		},
	}

	factory := NewItemFactory(itemMap)

	_, err := factory.Create("type1", 0, 0, "id1")
	if err != nil {
		t.Errorf("Create() type1 returned error: %v", err)
	}

	_, err = factory.Create("type2", 0, 0, "id2")
	if err != nil {
		t.Errorf("Create() type2 returned error: %v", err)
	}
}

func TestItemFactory_Create_EmptyMap(t *testing.T) {
	itemMap := ItemMap[*mockConcreteItem]{}

	factory := NewItemFactory(itemMap)
	item, err := factory.Create("any", 0, 0, "id")

	if err == nil {
		t.Error("Create() with empty map should return error")
	}
	if item != nil {
		t.Error("Create() with empty map should return nil item")
	}
}
