package main

import (
	"fmt"
	"os"

	"todo-tui/backend"
	"todo-tui/calendar"
	"todo-tui/notes"
	"todo-tui/planner"

	"go.yaml.in/yaml/v3"
)

// TODO:
// Calendar is being loaded with MockCalendar
// until sql representation is finished.
// Once done: Change the way labels is being initialized
type Config struct {
	Colors struct {
		Text    backend.HexCode `yaml:"text"`
		Border  backend.HexCode `yaml:"border"`
		Weekday backend.HexCode `yaml:"weekday"`
		Weekend backend.HexCode `yaml:"weekend"`
	} `yaml:"colors"`
	Labels   backend.Labels   `yaml:"labels"`
	Calendar backend.Calendar `yaml:"-"`
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

	config.Calendar = MockCalendar()
	config.Calendar.Labels.Selected = config.Labels.Selected
	return config, nil
}

func (c Config) LoadCalendarModel() calendar.Calendar {
	style := calendar.GetDefaultStyle()
	style.ColorWeekday = c.Colors.Weekday
	style.ColorWeekend = c.Colors.Weekend
	currLabel := c.Labels.Selected
	style.ColorAccent = c.Labels.Info[currLabel].Color

	cal := calendar.NewCalendar()
	cal.Style = style
	cal.UpdateGrid(c.Calendar)
	return cal
}

func (c Config) LoadNotesModel() notes.Notes {
	style := notes.GetDefaultStyle()
	style.ColorText = c.Colors.Text
	currLabel := c.Labels.Selected
	style.ColorAccent = c.Labels.Info[currLabel].Color

	notes := notes.NewNotes("")
	notes.Style = style

	return notes
}

func (c Config) LoadPlannerModel() planner.Planner {
	p := planner.NewPlanner(c.Calendar)
	p.Style.BorderColor = c.Colors.Border
	p.CalendarModel = c.LoadCalendarModel()
	p.NotesModel = c.LoadNotesModel()
	return p
}
