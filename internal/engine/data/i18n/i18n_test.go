package i18n

import (
	"io"
	"io/fs"
	"testing"
	"testing/fstest"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		fs      fstest.MapFS
		lang    string
		wantErr bool
		errMsg  string
	}{
		{
			name: "AC1: Load success with valid JSON",
			fs: fstest.MapFS{
				"assets/lang/en.json": &fstest.MapFile{
					Data: []byte(`{"greeting":"Hello","farewell":"Goodbye"}`),
				},
			},
			lang:    "en",
			wantErr: false,
		},
		{
			name:    "AC2: Load file not found",
			fs:      fstest.MapFS{},
			lang:    "xx",
			wantErr: true,
			errMsg:  "failed to open language file",
		},
		{
			name: "AC3: Load invalid JSON",
			fs: fstest.MapFS{
				"assets/lang/bad.json": &fstest.MapFile{
					Data: []byte(`{invalid json}`),
				},
			},
			lang:    "bad",
			wantErr: true,
			errMsg:  "failed to unmarshal language file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewI18nManager(&tt.fs)
			err := m.Load(tt.lang)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && tt.errMsg != "" {
				if err.Error() != tt.errMsg && !contains(err.Error(), tt.errMsg) {
					t.Errorf("Load() error = %v, want error containing %q", err, tt.errMsg)
				}
			}
		})
	}
}

func TestT(t *testing.T) {
	fs := fstest.MapFS{
		"assets/lang/en.json": &fstest.MapFile{
			Data: []byte(`{"greeting":"Hello","welcome":"Welcome, %s","simple":"Simple value"}`),
		},
	}
	m := NewI18nManager(&fs)
	m.Load("en")

	tests := []struct {
		name string
		key  string
		args []any
		want string
	}{
		{
			name: "AC4: T() known key",
			key:  "greeting",
			args: nil,
			want: "Hello",
		},
		{
			name: "AC5: T() missing key fallback",
			key:  "unknown_key",
			args: nil,
			want: "unknown_key",
		},
		{
			name: "AC6: T() with formatting",
			key:  "welcome",
			args: []any{"Alice"},
			want: "Welcome, Alice",
		},
		{
			name: "AC7: T() no formatting",
			key:  "simple",
			args: nil,
			want: "Simple value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := m.T(tt.key, tt.args...)
			if got != tt.want {
				t.Errorf("T(%q, %v) = %q, want %q", tt.key, tt.args, got, tt.want)
			}
		})
	}
}

func TestLoadReadError(t *testing.T) {
	// AC8: Coverage for read error path
	fs := &errorFS{}
	m := NewI18nManager(fs)
	err := m.Load("en")
	if err == nil {
		t.Errorf("Load() expected error, got nil")
	}
	if !contains(err.Error(), "failed to read language file") {
		t.Errorf("Load() error = %v, want error containing 'failed to read language file'", err)
	}
}

type errorFS struct{}

func (e *errorFS) Open(name string) (fs.File, error) {
	return &errorFile{}, nil
}

type errorFile struct{}

func (f *errorFile) Stat() (fs.FileInfo, error) {
	return nil, io.EOF
}

func (f *errorFile) Read(b []byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

func (f *errorFile) Close() error {
	return nil
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
