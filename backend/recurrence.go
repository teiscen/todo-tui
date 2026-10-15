package backend

// Recurrence kinds shared by the task editor and the database layer.
const (
	RecurDaily    = "daily"
	RecurWeekdays = "weekdays"
	RecurInterval = "interval"
)

// WeekdaysMonFri is the weekday bitmask for Monday through Friday (bit 0 is
// Sunday). It is the default schedule for a "weekdays" recurrence.
const WeekdaysMonFri = 62

// DefaultIntervalDays is the interval the editor proposes for a new
// "interval" recurrence.
const DefaultIntervalDays = 3
