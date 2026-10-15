package todo

import (
	"database/sql"

	"todo-tui/backend"
	"todo-tui/database"
	"todo-tui/leftbox"
	"todo-tui/preview"
	"todo-tui/taskform"
)

// detailFromRow translates a resolved leftbox row into the data the preview
// pane draws, so preview never has to know about leftbox.
func detailFromRow(row leftbox.Row) preview.Detail {
	return preview.Detail{
		ID:         row.ID,
		Title:      row.Title,
		Done:       row.Done,
		LabelID:    row.LabelID,
		Color:      row.Color,
		Schedule:   row.Schedule,
		Due:        row.Due,
		DueTone:    row.DueTone,
		StepsDone:  row.StepsDone,
		StepsTotal: row.StepsTotal,
	}
}

// formTaskFromDB converts a persisted task into the form's own representation,
// so the editor never has to know about database.Task.
func formTaskFromDB(task database.Task) taskform.Task {
	return taskform.Task{
		ID:           task.ID,
		Title:        task.Title,
		LabelID:      backend.LabelID(task.LabelID.String),
		Horizon:      task.Horizon,
		DueDate:      task.DueDate.String,
		Notes:        task.Notes,
		Recurring:    task.Recurring,
		Kind:         task.Kind,
		Weekdays:     task.Weekdays,
		IntervalDays: int(task.IntervalDays.Int64),
		StartDate:    task.StartDate,
	}
}

// dbTaskFromForm converts the form's task back into the database
// representation so it can be persisted.
func dbTaskFromForm(task taskform.Task) database.Task {
	return database.Task{
		ID:        task.ID,
		Title:     task.Title,
		LabelID:   sql.NullString{String: string(task.LabelID), Valid: task.LabelID != ""},
		Horizon:   task.Horizon,
		DueDate:   sql.NullString{String: task.DueDate, Valid: task.DueDate != ""},
		Notes:     task.Notes,
		Recurring: task.Recurring,
		Kind:      task.Kind,
		Weekdays:  task.Weekdays,
		IntervalDays: sql.NullInt64{
			Int64: int64(task.IntervalDays),
			Valid: task.IntervalDays > 0,
		},
		StartDate: task.StartDate,
	}
}
