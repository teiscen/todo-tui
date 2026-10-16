package config

import (
	"os"
	"path/filepath"
	"testing"

	"todo-tui/backend"
)

func TestLoadConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := `
layout:
  tasks:
    width: 60
    height: 40
  preview:
    width: 50
colors:
  text: "#111111"
  border: "#222222"
  weekday: "#333333"
  weekend: "#444444"
labels:
  selected: work
  order: [uni, work]
  info:
    uni:
      name: University
      color: "#81c8be"
    work:
      name: Work
      color: "#ca9ee6"
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}

	if cfg.Layout.Tasks.Width != 60 || cfg.Layout.Tasks.Height != 40 {
		t.Errorf("Layout.Tasks = %+v, want {60 40}", cfg.Layout.Tasks)
	}
	if cfg.Colors.Text != backend.HexCode("#111111") {
		t.Errorf("Colors.Text = %q, want #111111", cfg.Colors.Text)
	}
	if cfg.Labels.Selected != "work" {
		t.Errorf("Labels.Selected = %q, want work", cfg.Labels.Selected)
	}
}

func TestLayoutDefaults(t *testing.T) {
	var cfg Config
	cfg.setLayoutDefaults()
	if cfg.Layout.Tasks.Width != 50 || cfg.Layout.Tasks.Height != 43 {
		t.Errorf("Tasks defaults = %+v, want {50 43}", cfg.Layout.Tasks)
	}
	if cfg.Layout.Editor.Width != cfg.Layout.Preview.Width {
		t.Errorf("Editor.Width = %d, want Preview.Width %d", cfg.Layout.Editor.Width, cfg.Layout.Preview.Width)
	}
}

func TestThemeOverridesPartialConfig(t *testing.T) {
	cfg := Config{
		Colors: ColorsConfig{
			Text:    "#111111",
			Weekend: "#444444",
		},
	}
	th := cfg.Theme()

	if th.Text != "#111111" {
		t.Errorf("Theme.Text = %q, want #111111", th.Text)
	}
	if th.Weekend != "#444444" {
		t.Errorf("Theme.Weekend = %q, want #444444", th.Weekend)
	}
	// Unset colors keep the default theme values.
	if th.Border != backend.DefaultTheme().Border {
		t.Errorf("Theme.Border = %q, want default %q", th.Border, backend.DefaultTheme().Border)
	}
	if th.Weekday != backend.DefaultTheme().Weekday {
		t.Errorf("Theme.Weekday = %q, want default %q", th.Weekday, backend.DefaultTheme().Weekday)
	}
	if th.Dim != backend.DefaultTheme().Dim {
		t.Errorf("Theme.Dim = %q, want default %q", th.Dim, backend.DefaultTheme().Dim)
	}
}

func TestThemeEmptyConfigUsesDefaults(t *testing.T) {
	var cfg Config
	th := cfg.Theme()
	if th != backend.DefaultTheme() {
		t.Errorf("Theme() from empty config = %+v, want %+v", th, backend.DefaultTheme())
	}
}
