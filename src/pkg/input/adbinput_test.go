package input

import (
	"testing"
)

func TestEscapeForInput(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello", "hello"},
		{"hello world", "hello%sworld"},
		{"it's", "it\\'s"},
		{"a&b", "a\\&b"},
		{"test(1)", "test\\(1\\)"},
	}

	for _, tt := range tests {
		got := escapeForInput(tt.input)
		if got != tt.want {
			t.Errorf("escapeForInput(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestResolveKeycode(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"HOME", "KEYCODE_HOME"},
		{"home", "KEYCODE_HOME"},
		{"BACK", "KEYCODE_BACK"},
		{"ENTER", "KEYCODE_ENTER"},
		{"KEYCODE_HOME", "KEYCODE_HOME"},
		{"VOLUME_UP", "KEYCODE_VOLUME_UP"},
		{"SPACE", "KEYCODE_SPACE"},
		{"UNKNOWN_KEY", "KEYCODE_UNKNOWN_KEY"},
	}

	for _, tt := range tests {
		got := resolveKeycode(tt.input)
		if got != tt.want {
			t.Errorf("resolveKeycode(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
