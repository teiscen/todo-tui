package planner

import (
	backend "todo-tui/backend_old"
)

type PlannerStyle struct {
	BorderColor backend.HexCode
}

func GetDefaultStyle() PlannerStyle {
	return PlannerStyle{
		BorderColor: backend.HexCode("#C6D0F5"),
	}
}
