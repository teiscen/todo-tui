package todo

import (
	"database/sql"
	"fmt"
	"time"

	"todo-tui/backend"
	"todo-tui/config"
	"todo-tui/database"
	"todo-tui/leftbox"
	"todo-tui/planner"
	"todo-tui/preview"
	"todo-tui/taskform"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Todo struct {
	db *sql.DB
	cfg config.Config
	err error
	plannerModel planner.Planner
	leftboxModel leftbox.Leftbox
	previewModel preview.Preview
	form taskform.Form
	formActive bool
}

func NewTodo(cfg config.Config, db *sql.DB) *Todo {
	t := &Todo{
		db: db, cfg: cfg, plannerModel: cfg.LoadPlannerModel(),
		leftboxModel: cfg.LoadLeftboxModel(), previewModel: cfg.LoadPreviewModel(),
		form: taskform.New(cfg.Labels),
	}
	for _, id := range cfg.Labels.Order {
		info := cfg.Labels.Info[id]
		_, _ = db.Exec("INSERT INTO labels(name,color) VALUES(?,?) ON CONFLICT(name) DO UPDATE SET color=excluded.color", string(id), string(info.Color))
	}
	t.reloadTasks()
	t.syncPreview()
	return t
}

func (t *Todo) Init() tea.Cmd {
	return func() tea.Msg {
		return planner.ChangeCalendarMsg{LabelID:*t.plannerModel.SelectedLabel(),NewDate:*t.plannerModel.SelectedDate()}
	}
}

func (t *Todo) View() tea.View {
	combined := lipgloss.JoinHorizontal(lipgloss.Top,t.leftboxModel.View().Content,t.previewModel.View().Content,t.plannerModel.View().Content)
	view:=tea.NewView(combined)
	view.MouseMode=tea.MouseModeAllMotion
	return view
}

func (t *Todo) Update(msg tea.Msg) (tea.Model,tea.Cmd) {
	switch m:=msg.(type) {
	case tea.MouseMotionMsg:
		if !t.formActive && m.X<t.leftboxModel.Style.Width && t.leftboxModel.HoverAt(m.Y) { t.syncPreview() }
		return t,nil
	case taskform.CancelMsg:
		t.formActive=false
		t.syncPreview()
		return t,nil
	case taskform.SavedMsg:
		_,err:=database.SaveTask(t.db,m.Task)
		if err!=nil { t.form.Error=err.Error(); t.formActive=true; t.previewModel.SetForm(&t.form); return t,nil }
		t.formActive=false
		t.reloadTasks()
		t.syncPreview()
		return t,nil
	case leftbox.ToggleTodoMsg:
		if err:=database.ToggleTask(t.db,m.ID,m.Recurring,m.Done); err!=nil { t.err=err } else { t.err=nil }
		t.reloadTasks(); t.syncPreview()
		return t,nil
	case leftbox.SelectTodoMsg:
		if !t.formActive { t.syncPreview() }
		return t,nil
	case planner.ChangeCalendarMsg:
		t.ReloadPlannerState(m.LabelID,m.NewDate); return t,nil
	case planner.SaveEntryMsg:
		t.AddEntry(m.Label,m.Date,m.Msg); t.ReloadPlannerState(m.Label,m.Date); return t,nil
	case planner.RemoveEntryMsg:
		t.RemoveEntry(m.Label,m.Date); t.ReloadPlannerState(m.Label,m.Date); return t,nil
	}
	if t.formActive {
		var cmd tea.Cmd
		t.form,cmd=t.form.Update(msg)
		t.previewModel.SetForm(&t.form)
		return t,cmd
	}
	if key,ok:=msg.(tea.KeyPressMsg); ok {
		switch key.String() {
		case "a":
			t.beginAdd(); return t,nil
		case "e":
			t.beginEdit(); return t,nil
		case "j","down","k","up","[","]","left_bracket","right_bracket","open_bracket","close_bracket"," ","space","f","enter":
			var cmd tea.Cmd
			t.leftboxModel,cmd=t.leftboxModel.Update(msg); t.syncPreview(); return t,cmd
		}
	}
	var cmd tea.Cmd
	t.plannerModel,cmd=t.plannerModel.Update(msg)
	return t,cmd
}

func (t *Todo) beginAdd() {
	task:=database.Task{Horizon:"small",LabelID:sql.NullString{String:string(t.cfg.Labels.Selected),Valid:t.cfg.Labels.Selected!=""}}
	t.form.Begin(task,false); t.formActive=true; t.previewModel.SetForm(&t.form)
}

func (t *Todo) beginEdit() {
	row,ok:=t.leftboxModel.Selected(); if !ok { return }
	tasks,err:=database.FetchTasks(t.db); if err!=nil { t.err=err; return }
	for _,task:=range tasks {
		if task.ID==row.ID { t.form.Begin(task,true); t.formActive=true; t.previewModel.SetForm(&t.form); return }
	}
}

func (t *Todo) reloadTasks() {
	tasks,err:=database.FetchTasks(t.db); if err!=nil { t.err=err; return }
	today:=backend.NewDate(time.Now())
	rows:=make([]leftbox.Row,0,len(tasks))
	for _,task:=range tasks {
		recurring:=task.Recurring
		done,err:=database.IsTaskDoneToday(t.db,task.ID,recurring,task.DoneAt)
		if err!=nil { t.err=err; continue }
		tab:=leftbox.TabSmall; if task.Horizon=="long" { tab=leftbox.TabLong }; if recurring { tab=leftbox.TabRecurring }
		row:=leftbox.Row{ID:task.ID,Tab:tab,Title:task.Title,Done:done,LabelID:backend.LabelID(task.LabelID.String)}
		if row.LabelID!="" { row.Color=t.cfg.Labels.Info[row.LabelID].Color }
		if recurring {
			switch task.Kind { case "daily": row.Schedule="daily"; case "weekdays": row.Schedule="M-F"; case "interval": row.Schedule=fmt.Sprintf("every %dd",task.IntervalDays.Int64) }
			row.Group=leftbox.GroupToday
			if done { row.Group=leftbox.GroupDone }
		} else if task.DueDate.Valid && task.DueDate.String!="" {
			due,parseErr:=backend.ParseSQL(task.DueDate.String)
			if parseErr==nil { row.Group=leftbox.GroupFor(&due,done,today); if !done { row.Due,row.DueTone=leftbox.DueInfo(due,today) } }
		} else { row.Group=leftbox.GroupFor(nil,done,today) }
		rows=append(rows,row)
	}
	t.leftboxModel.ReloadState(rows)
}

func (t *Todo) syncPreview() {
	if t.formActive { t.previewModel.SetForm(&t.form); return }
	if row,ok:=t.leftboxModel.Selected(); ok { t.previewModel.SetTask(row,true); return }
	t.previewModel.SetTask(leftbox.Row{},false)
}
