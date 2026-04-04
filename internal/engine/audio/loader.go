package audio

import (
	"io/fs"
	"log"
	"strings"
)

// LoadAudioAssetsFromFS is a helper function to load all audio files from an fs.FS.
func LoadAudioAssetsFromFS(assets fs.FS, am *AudioManager) {
	// WalkDir ensures recursive loading, including assets/audio/bleeps
	dir := "assets/audio"
	if err := fs.WalkDir(assets, dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".ogg") && !strings.HasSuffix(path, ".wav") && !strings.HasSuffix(path, ".mp3") {
			return nil
		}
		data, err := fs.ReadFile(assets, path)
		if err != nil {
			log.Printf("failed to read embedded file %s: %v", path, err)
			return nil
		}
		am.Add(path, data)
		return nil
	}); err != nil {
		// Handle case where audio directory doesn't exist (e.g., in tests)
		if !strings.Contains(err.Error(), "file does not exist") {
			log.Fatalf("error reading embedded audio dir: %v", err)
		}
		log.Printf("audio directory not found: %s", dir)
	}
}
