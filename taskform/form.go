package taskform

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"todo-tui/backend"
	"todo-tui/database"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type SavedMsg struct { Task database.Task }
type CancelMsg struct{}

type Form struct {
	Style backend.HexCode
	Labels backend.Labels
	Task database.Task
	Editing bool
	Focus int
	Title, Due, Notes, Interval string
	Recurring bool
	Kind string
	Error string
}

func New(labels backend.Labels) Form {
	return Form{Style: backend.HexCode("#C6D0F5"), Labels: labels, Task: database.Task{Horizon:"small"}, Kind:"daily", Interval:"3"}
}

func (f *Form) Begin(task database.Task, editing bool) {
	f.Editing, f.Task = editing, task
	f.Title, f.Due, f.Notes = task.Title, task.DueDate.String, task.Notes
	f.Recurring, f.Kind = task.Recurring, task.Kind
	if f.Kind == "" { f.Kind = "daily" }
	f.Interval = strconv.FormatInt(task.IntervalDays.Int64, 10)
	if f.Interval == "0" { f.Interval = "3" }
	f.Error, f.Focus = "", 0
}

func (f Form) Init() tea.Cmd { return nil }

func (f Form) Update(msg tea.Msg) (Form, tea.Cmd) {
	key, ok := msg.(tea.KeyPressMsg)
	if !ok { return f, nil }
	k := key.String()
	switch k {
	case "esc": return f, func() tea.Msg { return CancelMsg{} }
	case "ctrl+s", "ctrl+S":
		if strings.TrimSpace(f.Title) == "" { f.Error = "Title is required"; return f, nil }
		t := f.Task
		t.Title, t.Notes = strings.TrimSpace(f.Title), f.Notes
		if t.Horizon == "" { t.Horizon = "small" }
		t.DueDate = sql.NullString{}
		if f.Due != "" { t.DueDate = sql.NullString{String:f.Due,Valid:true} }
		if !t.LabelID.Valid || t.LabelID.String == "" { t.LabelID = sql.NullString{} }
		t.Recurring, t.Kind = f.Recurring, f.Kind
		t.StartDate = time.Now().Format("2006-01-02")
		if f.Kind == "interval" {
			n, err := strconv.ParseInt(f.Interval, 10, 64)
			if err != nil || n < 1 { f.Error = "Interval must be a positive number"; return f, nil }
			t.IntervalDays = sql.NullInt64{Int64:n,Valid:true}
		}
		if f.Kind == "weekdays" { t.Weekdays = 62 }
		return f, func() tea.Msg { return SavedMsg{Task:t} }
	case "tab", "down":
		f.Focus = (f.Focus + 1) % 8
		return f, nil
	case "shift+tab", "up":
		f.Focus = (f.Focus + 7) % 8
		return f,nil
	case "left", "right":
		if f.Focus == 1 && len(f.Labels.Order) > 0 {
			i := 0
			for n,id := range f.Labels.Order { if string(id)==f.Task.LabelID.String { i=n; break } }
			if k=="right" { i=(i+1)%len(f.Labels.Order) } else { i=(i+len(f.Labels.Order)-1)%len(f.Labels.Order) }
			f.Task.LabelID = sql.NullString{String:string(f.Labels.Order[i]),Valid:true}
		} else if f.Focus == 5 {
			kinds:=[]string{"daily","weekdays","interval"}
			i:=0; for n,v:=range kinds { if v==f.Kind { i=n } }
			if k=="right" { i=(i+1)%len(kinds) } else { i=(i+len(kinds)-1)%len(kinds) }
			f.Kind=kinds[i]
		} else if f.Focus == 7 {
			if f.Task.Horizon == "long" { f.Task.Horizon = "small" } else { f.Task.Horizon = "long" }
		}
		return f,nil
	case " ":
		if f.Focus==4 { f.Recurring=!f.Recurring; return f,nil }
		if f.Focus==0 { f.Title+=" "; return f,nil }
		if f.Focus==2 { f.Due+=" "; return f,nil }
		if f.Focus==3 { f.Notes+=" "; return f,nil }
		if f.Focus==6 { f.Interval+=" "; return f,nil }
	}
	if k=="backspace" || k=="delete" {
		switch f.Focus { case 0:f.Title=dropLast(f.Title); case 2:f.Due=dropLast(f.Due); case 3:f.Notes=dropLast(f.Notes); case 6:f.Interval=dropLast(f.Interval) }
		return f,nil
	}
	if len([]rune(k))==1 && k>=" " && k<="~" {
		switch f.Focus { case 0:f.Title+=k; case 2:f.Due+=k; case 3:f.Notes+=k; case 6:f.Interval+=k }
	}
	return f,nil
}

func (f Form) View() tea.View {
	border:=f.Style.ToStyle()
	lines:=[]string{
		border.Bold(true).Render("TASK EDITOR"),"",
		f.field(0,"TITLE",f.Title),
		f.field(1,"LABEL",f.labelText()),
		f.field(7,"HORIZON",f.Task.Horizon),
		f.field(2,"DUE DATE (YYYY-MM-DD)",f.Due),
		f.field(3,"NOTES",f.Notes),
		f.field(4,"RECURRING",map[bool]string{true:"Yes",false:"No"}[f.Recurring]),
	}
	if f.Recurring {
		lines=append(lines,f.field(5,"SCHEDULE",f.Kind))
		if f.Kind=="interval" { lines=append(lines,f.field(6,"EVERY N DAYS",f.Interval)) }
	}
	lines=append(lines,"",border.Render("Tab/↑↓ move · ←/→ change choice · Space toggle"),border.Render("Ctrl+S save · Esc cancel"))
	if f.Error!="" { lines=append(lines,"",backend.HexCode("#E78284").ToStyle().Render(f.Error)) }
	box:=lipgloss.NewStyle().Width(50).Height(42).Padding(1,2).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(string(f.Style)))
	return tea.NewView(box.Render(strings.Join(lines,"\n")))
}

func (f Form) field(index int,label,value string) string {
	prefix:="  "; if f.Focus==index { prefix="> " }
	return prefix+label+": "+value
}
func (f Form) labelText() string {
	if !f.Task.LabelID.Valid || f.Task.LabelID.String=="" { return "(none)" }
	if info,ok:=f.Labels.Info[backend.LabelID(f.Task.LabelID.String)]; ok { return info.Name }
	return f.Task.LabelID.String
}
func dropLast(s string) string { r:=[]rune(s); if len(r)>0 { return string(r[:len(r)-1]) }; return s }
