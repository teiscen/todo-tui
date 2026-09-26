package calendar

//
// import (
// 	"fmt"
// 	"os"
// 	"time"
//
// 	"todo-tui/artifacts"
// 	"todo-tui/backend"
//
// 	"charm.land/bubbles/v2/textarea"
// 	tea "charm.land/bubbletea/v2"
// 	lipgloss "charm.land/lipgloss/v2"
// )
//
// type Focus bool
//
// const (
// 	FocusCal  Focus = false
// 	FocusNote       = true
// )
//
// func (f Focus) toggleFocus() Focus {
// 	return Focus(!f)
// }
//
// type Model struct {
// 	selDate  backend.CalendarDate
// 	selLabel int
// 	labels   []backend.Label
// 	style    CalendarStyle
// 	calendar backend.Calendar
// 	rInfo    RenderInfo
// 	config   artifacts.Config
// 	focus    Focus
// 	notes    textarea.Model
// }
//
// // Need to populate the calenderPages from a config file
// func initialModel() Model {
// 	config := artifacts.defaultConfig()
//
// 	ti := textarea.New()
// 	ti.Prompt = ""
// 	ti.ShowLineNumbers = false
// 	// ti.EndOfBufferCharacter = '~'
// 	ti.MaxWidth = 30
// 	ti.SetHeight(23) // <- this is what actually controls the fixed height
// 	ti.KeyMap = textarea.DefaultKeyMap()
//
// 	styles := textarea.DefaultDarkStyles()
// 	styles.Cursor.Shape = tea.CursorUnderline
// 	ti.SetStyles(styles)
//
// 	labels := make([]Label, 0, len(config.Labels))
// 	for label := range config.Labels {
// 		labels = append(labels, label)
// 	}
//
// 	m := Model{
// 		timeToCalendarDate(time.Now()),
// 		0,
// 		labels,
// 		doubleStyle,
// 		mockCalendar(config.Labels),
// 		RenderInfo{},
// 		config,
// 		FocusCal,
// 		ti,
// 	}
//
// 	m.rInfo = CreateGrid(m)
//
// 	return m
// }
//
// func (m Model) Init() tea.Cmd {
// 	return tea.Batch(textarea.Blink, tea.RequestBackgroundColor)
// }
//
// func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	var cmd tea.Cmd
// 	var cmds []tea.Cmd
//
// 	switch msg := msg.(type) {
// 	case tea.BackgroundColorMsg:
// 		// m.textarea.SetStyles(textarea.DefaultStyles(msg.IsDark()))
//
// 	case tea.KeyPressMsg:
// 		switch msg.String() {
// 		case "ctrl+c", "q", "shift-c":
// 			return m, tea.Quit
// 		case "tab":
// 			m.focus = m.focus.toggleFocus()
// 			if m.focus == FocusCal {
// 				m.notes.Blur()
// 			} else {
// 				cmd = m.notes.Focus()
// 				cmds = append(cmds, cmd)
// 			}
// 			return m, tea.Batch(cmds...)
// 		case "p", "left":
// 			if m.focus == FocusCal {
// 				m.movePage(Left)
// 				return m, nil
// 			}
// 		case "n", "right":
// 			if m.focus == FocusCal {
// 				m.movePage(Right)
// 				return m, nil
// 			}
// 		case "k", "h", "l", "j":
// 			if m.focus == FocusCal {
// 				switch msg.String() {
// 				case "k":
// 					m.moveCursor(Up)
// 				case "h":
// 					m.moveCursor(Left)
// 				case "l":
// 					m.moveCursor(Right)
// 				case "j":
// 					m.moveCursor(Down)
// 				}
// 				m.rInfo = CreateGrid(m)
// 				m.notes.SetValue(m.calendar.Months[m.selDate][m.labels[m.selLabel]].Msg)
// 				return m, tea.ClearScreen
// 			}
// 		}
// 	}
//
// 	// Anything not already handled above (typed characters, arrow keys
// 	// inside the textarea, etc.) goes to whichever child has focus.
// 	if m.focus == FocusNote {
// 		m.notes, cmd = m.notes.Update(msg)
// 		cmds = append(cmds, cmd)
// 	}
//
// 	return m, tea.Batch(cmds...)
// }
//
// func (m Model) View() tea.View {
// 	label := m.config.Labels[m.labels[m.selLabel]]
// 	labelColor := m.config.Colors[label.ColorName]
//
// 	calendarComplete := genBorderCalendar(
// 		m.RenderCalendar(),
// 		m.cursorColor(),
// 		calendarDateToTime(m.selDate).Format("January 02"),
// 		m.cursorColor(),
// 		label.Name,
// 		labelColor,
// 	)
//
// 	targetWidth := lipgloss.Width(calendarComplete)
// 	m.notes.SetWidth(targetWidth - 6) // 6 = genBorderNotes' border+padding overhead, confirmed above
//
// 	color := m.config.Colors["valid"]
// 	notesComplete := genBorderNotes(m.notes.View(), color)
//
// 	complete := genBorder(calendarComplete, notesComplete, color)
// 	return tea.NewView(complete)
// }
//
// func TestCalendar() {
// 	p := tea.NewProgram(initialModel())
// 	if _, err := p.Run(); err != nil {
// 		fmt.Printf("Alas, there's been an error: %v", err)
// 		os.Exit(1)
// 	}
// }
