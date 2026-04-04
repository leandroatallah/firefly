package gamesetup

import (
	"io/fs"
	"strings"
)

func collectSpeechBleeps(assets fs.FS) []string {
	paths := collectAudioInDir(assets, "assets/audio/bleeps")
	if len(paths) > 0 {
		return paths
	}
	return collectAudioWithPrefix(assets, "assets/audio", "bleep")
}

func collectAudioInDir(assets fs.FS, dir string) []string {
	entries, err := fs.ReadDir(assets, dir)
	if err != nil {
		return nil
	}
	paths := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".ogg") && !strings.HasSuffix(name, ".wav") && !strings.HasSuffix(name, ".mp3") {
			continue
		}
		paths = append(paths, dir+"/"+name)
	}
	return paths
}

func collectAudioWithPrefix(assets fs.FS, dir, prefix string) []string {
	entries, err := fs.ReadDir(assets, dir)
	if err != nil {
		return nil
	}
	paths := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasPrefix(name, prefix) {
			continue
		}
		if !strings.HasSuffix(name, ".ogg") && !strings.HasSuffix(name, ".wav") && !strings.HasSuffix(name, ".mp3") {
			continue
		}
		paths = append(paths, dir+"/"+name)
	}
	return paths
}
