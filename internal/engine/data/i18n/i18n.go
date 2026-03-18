package i18n

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
)

type I18nManager struct {
	translations map[string]string
	assets       fs.FS
}

func NewI18nManager(assets fs.FS) *I18nManager {
	return &I18nManager{
		translations: make(map[string]string),
		assets:       assets,
	}
}

func (m *I18nManager) Load(langCode string) error {
	path := fmt.Sprintf("assets/lang/%s.json", langCode)
	f, err := m.assets.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open language file %s: %w", path, err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("failed to read language file %s: %w", path, err)
	}

	var translations map[string]string
	if err := json.Unmarshal(data, &translations); err != nil {
		return fmt.Errorf("failed to unmarshal language file %s: %w", path, err)
	}

	m.translations = translations
	return nil
}

func (m *I18nManager) T(key string, args ...any) string {
	val, ok := m.translations[key]
	if !ok {
		return key
	}
	if len(args) > 0 {
		return fmt.Sprintf(val, args...)
	}
	return val
}
