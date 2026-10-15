package config

import (
	"fmt"
	"os"

	"todo-tui/backend"

	"go.yaml.in/yaml/v3"
)

const (
	defaultPanelWidth  = 50
	defaultPanelHeight = 43
)

// Config is the on-disk application configuration. It holds data only; it does
// not build any pane or model — the todo root does that in NewTodo.
type Config struct {
	Layout LayoutConfig   `yaml:"layout"`
	Colors ColorsConfig   `yaml:"colors"`
	Labels backend.Labels `yaml:"labels"`
}

// ColorsConfig is the subset of the theme that config.yaml may override. Any
// field left empty falls back to backend.DefaultTheme.
type ColorsConfig struct {
	Text    backend.HexCode `yaml:"text"`
	Border  backend.HexCode `yaml:"border"`
	Weekday backend.HexCode `yaml:"weekday"`
	Weekend backend.HexCode `yaml:"weekend"`
}

// LayoutConfig centralizes dashboard panel dimensions in terminal cells.
type LayoutConfig struct {
	Tasks   PanelSize `yaml:"tasks"`
	Preview PanelSize `yaml:"preview"`
	Editor  PanelSize `yaml:"editor"`
}

// PanelSize describes the outer width and height of a rendered panel.
type PanelSize struct {
	Width  int `yaml:"width"`
	Height int `yaml:"height"`
}

// Theme overlays the configured colors onto the default theme, so a partial
// config still yields a complete, usable color scheme.
func (c Config) Theme() backend.Theme {
	t := backend.DefaultTheme()
	if c.Colors.Text != "" {
		t.Text = c.Colors.Text
	}
	if c.Colors.Border != "" {
		t.Border = c.Colors.Border
	}
	if c.Colors.Weekday != "" {
		t.Weekday = c.Colors.Weekday
	}
	if c.Colors.Weekend != "" {
		t.Weekend = c.Colors.Weekend
	}
	return t
}

func (c *Config) setLayoutDefaults() {
	if c.Layout.Tasks.Width <= 0 {
		c.Layout.Tasks.Width = defaultPanelWidth
	}
	if c.Layout.Tasks.Height <= 0 {
		c.Layout.Tasks.Height = defaultPanelHeight
	}
	if c.Layout.Preview.Width <= 0 {
		c.Layout.Preview.Width = defaultPanelWidth
	}
	if c.Layout.Preview.Height <= 0 {
		c.Layout.Preview.Height = defaultPanelHeight
	}
	if c.Layout.Editor.Width <= 0 {
		c.Layout.Editor.Width = c.Layout.Preview.Width
	}
	if c.Layout.Editor.Height <= 0 {
		c.Layout.Editor.Height = c.Layout.Preview.Height
	}
}

func LoadConfig(filePath string) (Config, error) {
	var config Config

	data, err := os.ReadFile(filePath)
	if err != nil {
		return config, fmt.Errorf("read config: %w", err)
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("parse config: %w", err)
	}

	config.setLayoutDefaults()
	return config, nil
}
