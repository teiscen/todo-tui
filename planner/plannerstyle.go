package planner

import (
	backend "todo-tui/backend"
)

type PlannerStyle struct {
	Theme backend.Theme
}

// NewStyle builds a planner style from the shared theme.
func NewStyle(theme backend.Theme) PlannerStyle {
	return PlannerStyle{
		Theme: theme,
	}
}
