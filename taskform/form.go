package taskform

import (
	"strconv"
	"strings"
	"time"

	"todo-tui/backend"
	"todo-tui/ui"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	fieldTitle = iota
	fieldLabel
	fieldHorizon
	fieldDueDate
	fieldNotes
	fieldRecurring
	fieldSchedule
	fieldInterval
	fieldCount
)

var recurrenceKinds = []string{backend.RecurDaily, backend.RecurWeekdays, backend.RecurInterval}

// Task is the form's own representation of the task being edited. It carries
// the fields the editor needs without depending on the database package; the
// root translates to and from database.Task.
type Task struct {
	ID           int
	Title        string
	LabelID      backend.LabelID
	Horizon      string
	DueDate      string
	Notes        string
	Recurring    bool
	Kind         string
	Weekdays     int
	IntervalDays int
	StartDate    string
}

// SavedMsg is emitted when the form has valid data ready to persist.
type SavedMsg struct {
	Task Task
}

// CancelMsg tells the root model to close the form without saving.
type CancelMsg struct{}

// Form edits the fields shared by new and existing tasks.
type Form struct {
	Theme  backend.Theme
	Width  int
	Height int
	Labels backend.Labels

	Task    Task
	Editing bool
	Focus   int

	Title     string
	Due       string
	Notes     string
	Interval  string
	Recurring bool
	Kind      string
	Error     string
}

// New returns a form configured with the supplied labels and default values.
func New(theme backend.Theme, labels backend.Labels, width, height int) Form {
	return Form{
		Theme:    theme,
		Width:    width,
		Height:   height,
		Labels:   labels,
		Task:     Task{Horizon: "small"},
		Kind:     backend.RecurDaily,
		Interval: strconv.Itoa(backend.DefaultIntervalDays),
	}
}

// Begin loads a task into the form. A zero-value task starts a new entry.
func (f *Form) Begin(task Task, editing bool) {
	f.Task = task
	f.Editing = editing
	f.Title = task.Title
	f.Due = task.DueDate
	f.Notes = task.Notes
	f.Recurring = task.Recurring
	f.Kind = task.Kind
	if f.Kind == "" {
		f.Kind = backend.RecurDaily
	}
	f.Interval = strconv.Itoa(task.IntervalDays)
	if f.Interval == "0" {
		f.Interval = strconv.Itoa(backend.DefaultIntervalDays)
	}
	f.Error = ""
	f.Focus = fieldTitle
}

// Init satisfies tea.Model-style initialization; the form needs no command.
func (f Form) Init() tea.Cmd { return nil }

// Update handles form navigation, editing, validation, and save/cancel requests.
func (f Form) Update(msg tea.Msg) (Form, tea.Cmd) {
	key, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return f, nil
	}

	switch key.String() {
	case "esc":
		return f, func() tea.Msg { return CancelMsg{} }
	case "ctrl+s":
		return f, f.saveCommand()
	case "tab", "down":
		f.Focus = (f.Focus + 1) % fieldCount
		return f, nil
	case "shift+tab", "up":
		f.Focus = (f.Focus + fieldCount - 1) % fieldCount
		return f, nil
	case "left", "right":
		f.changeChoice(key.String() == "right")
		return f, nil
	case " ", "space":
		if f.Focus == fieldRecurring {
			f.Recurring = !f.Recurring
			return f, nil
		}
		if f.isTextField() {
			f.appendText(" ")
			return f, nil
		}
	case "backspace", "delete":
		f.removeLastCharacter()
		return f, nil
	}

	// Only append printable ASCII characters to editable text fields.
	if value := key.String(); len(value) == 1 && value >= " " && value <= "~" {
		f.appendText(value)
	}
	return f, nil
}

// saveCommand validates the required fields before returning the task to the root model.
func (f *Form) saveCommand() tea.Cmd {
	if strings.TrimSpace(f.Title) == "" {
		f.Error = "Title is required"
		return nil
	}

	task := f.Task
	task.Title = strings.TrimSpace(f.Title)
	task.Notes = f.Notes
	if task.Horizon != "long" {
		task.Horizon = "small"
	}
	task.DueDate = ""
	if strings.TrimSpace(f.Due) != "" {
		if _, err := time.Parse("2006-01-02", f.Due); err != nil {
			f.Error = "Due date must use YYYY-MM-DD"
			return nil
		}
		task.DueDate = f.Due
	}

	task.Recurring = f.Recurring
	task.Kind = f.Kind
	task.StartDate = time.Now().Format("2006-01-02")
	task.IntervalDays = 0
	task.Weekdays = 0
	if task.Recurring {
		switch task.Kind {
		case backend.RecurWeekdays:
			task.Weekdays = backend.WeekdaysMonFri
		case backend.RecurInterval:
			days, err := strconv.ParseInt(f.Interval, 10, 64)
			if err != nil || days < 1 {
				f.Error = "Interval must be a positive number"
				return nil
			}
			task.IntervalDays = int(days)
		}
	}
	f.Error = ""
	return func() tea.Msg { return SavedMsg{Task: task} }
}

// changeChoice moves between label, horizon, and recurrence options.
func (f *Form) changeChoice(forward bool) {
	switch f.Focus {
	case fieldLabel:
		if len(f.Labels.Order) == 0 {
			return
		}
		index := 0
		for i, id := range f.Labels.Order {
			if id == f.Task.LabelID {
				index = i
				break
			}
		}
		if forward {
			index = (index + 1) % len(f.Labels.Order)
		} else {
			index = (index + len(f.Labels.Order) - 1) % len(f.Labels.Order)
		}
		f.Task.LabelID = f.Labels.Order[index]
	case fieldHorizon:
		if f.Task.Horizon == "long" {
			f.Task.Horizon = "small"
		} else {
			f.Task.Horizon = "long"
		}
	case fieldSchedule:
		index := 0
		for i, kind := range recurrenceKinds {
			if kind == f.Kind {
				index = i
				break
			}
		}
		if forward {
			index = (index + 1) % len(recurrenceKinds)
		} else {
			index = (index + len(recurrenceKinds) - 1) % len(recurrenceKinds)
		}
		f.Kind = recurrenceKinds[index]
	}
}

func (f Form) isTextField() bool {
	switch f.Focus {
	case fieldTitle, fieldDueDate, fieldNotes, fieldInterval:
		return true
	default:
		return false
	}
}

func (f *Form) appendText(value string) {
	switch f.Focus {
	case fieldTitle:
		f.Title += value
	case fieldDueDate:
		f.Due += value
	case fieldNotes:
		f.Notes += value
	case fieldInterval:
		f.Interval += value
	}
}

func (f *Form) removeLastCharacter() {
	switch f.Focus {
	case fieldTitle:
		f.Title = dropLast(f.Title)
	case fieldDueDate:
		f.Due = dropLast(f.Due)
	case fieldNotes:
		f.Notes = dropLast(f.Notes)
	case fieldInterval:
		f.Interval = dropLast(f.Interval)
	}
}

// View renders the same editor for both adding and editing a task.
func (f Form) View() string {
	style := ui.Fg(f.Theme.Text)
	title := "NEW TASK"
	if f.Editing {
		title = "EDIT TASK"
	}
	lines := []string{
		style.Bold(true).Render(title),
		"",
		f.field(fieldTitle, "TITLE", f.Title),
		f.field(fieldLabel, "LABEL", f.labelText()),
		f.field(fieldHorizon, "HORIZON", f.Task.Horizon),
		f.field(fieldDueDate, "DUE DATE (YYYY-MM-DD)", f.Due),
		f.field(fieldNotes, "NOTES", f.Notes),
		f.field(fieldRecurring, "RECURRING", map[bool]string{true: "Yes", false: "No"}[f.Recurring]),
	}
	if f.Recurring {
		lines = append(lines, f.field(fieldSchedule, "SCHEDULE", f.Kind))
		if f.Kind == backend.RecurInterval {
			lines = append(lines, f.field(fieldInterval, "EVERY N DAYS", f.Interval))
		}
	}
	lines = append(lines, "", style.Render("Tab/↑↓ move · ←/→ change choices · Space toggle"), style.Render("Ctrl+S save · Esc cancel"))
	if f.Error != "" {
		lines = append(lines, "", ui.Fg(f.Theme.Late).Render(f.Error))
	}

	box := lipgloss.NewStyle().
		Width(f.Width).
		Height(f.Height).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(string(f.Theme.Border)))
	return box.Render(strings.Join(lines, "\n"))
}

func (f Form) field(index int, label, value string) string {
	prefix := "  "
	if f.Focus == index {
		prefix = "> "
	}
	return prefix + label + ": " + value
}

func (f Form) labelText() string {
	if f.Task.LabelID == "" {
		return "(none)"
	}
	if info, ok := f.Labels.Info[f.Task.LabelID]; ok {
		return info.Name
	}
	return string(f.Task.LabelID)
}

func dropLast(value string) string {
	runes := []rune(value)
	if len(runes) == 0 {
		return value
	}
	return string(runes[:len(runes)-1])
}
