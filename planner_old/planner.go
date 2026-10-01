package calendar

// import (
// 	"todo-tui/backend"
// )
//
// type Planner struct {
// 	Cal  Calendar
// 	Note Note
// 	pRI  backend.PlannerRenderInfo
// }
//
// func NewPlanner(s backend.State) Planner {
// 	return Planner{
// 		NewCalendar(s),
// 		NewNote(s),
// 		s.GetPlannerRenderInfo(),
// 	}
// }
//
// // func (p *Planner) Update(pRI backend.PlannerRenderInfo) {
// func (p *Planner) Update(s backend.State) {
// 	p.Cal.Update(s)
// 	p.Note.Update(s)
// 	p.pRI = s.GetPlannerRenderInfo()
// }
//
// func (p Planner) Render() string {
// 	calRender := p.Cal.Render()
// 	noteRender := p.Note.Render()
// 	planRender := p.RenderBorder(calRender, noteRender)
//
// 	return planRender
// } // calendar/planner.go
