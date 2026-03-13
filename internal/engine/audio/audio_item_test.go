package audio

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestAudioItem(t *testing.T) {
	item := &AudioItem{
		name: "test",
		data: []byte{1, 2, 3},
	}

	if item.Name() != "test" {
		t.Errorf("expected name test, got %s", item.Name())
	}
	if len(item.Data()) != 3 {
		t.Errorf("expected data length 3, got %d", len(item.Data()))
	}
}

func TestLoad(t *testing.T) {
	am := getTestAudioManager()

	// Create a temp file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.wav")
	err := os.WriteFile(tmpFile, []byte("dummy data"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	item, err := am.Load(tmpFile)
	if err != nil {
		t.Fatalf("failed to load audio: %v", err)
	}

	if item.Name() != tmpFile {
		t.Errorf("expected name %s, got %s", tmpFile, item.Name())
	}

	// Test non-existent file
	_, err = am.Load("non-existent.wav")
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}

func TestLoadFromFS(t *testing.T) {
	am := getTestAudioManager()

	mockFS := fstest.MapFS{
		"test.wav": &fstest.MapFile{Data: []byte("dummy data")},
	}

	item, err := am.LoadFromFS(mockFS, "test.wav")
	if err != nil {
		t.Fatalf("failed to load from FS: %v", err)
	}

	if item.Name() != "test.wav" {
		t.Errorf("expected name test.wav, got %s", item.Name())
	}

	// Test non-existent file
	_, err = am.LoadFromFS(mockFS, "non-existent.wav")
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}

type errorFS struct{}

func (errorFS) Open(name string) (fs.File, error) {
	return nil, fs.ErrNotExist
}

func TestLoadFromFSError(t *testing.T) {
	am := getTestAudioManager()
	_, err := am.LoadFromFS(errorFS{}, "test.wav")
	if err == nil {
		t.Error("expected error")
	}
}
