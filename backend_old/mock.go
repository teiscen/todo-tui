package backend

import (
	"math/rand"
	"time"
)

func MockCalendar() Calendar {
	labels := Labels{
		Selected: "uni",
		Info: map[LabelID]LabelInfo{
			"uni":  {Name: "University", Color: HexCode("#81c8be")},
			"work": {Name: "Work", Color: HexCode("#ca9ee6")},
			"gym":  {Name: "Gym", Color: HexCode("#ef9f76")},
			"app":  {Name: "Application", Color: HexCode("#babbf1")},
			"vol":  {Name: "Volunteer", Color: HexCode("#f2d5cf")},
		},
		Order: []LabelID{
			"uni", "work", "gym", "app", "vol",
		},
	}

	// Fixed seed so the mock data is identical on every run.
	rng := rand.New(rand.NewSource(2026))
	months := make(map[Date]map[LabelID]Entry)

	set := func(d Date, id LabelID, s Status, msg string) {
		if months[d] == nil {
			months[d] = make(map[LabelID]Entry)
		}
		months[d][id] = Entry{Status: s, Msg: msg}
	}
	pick := func(s []string) string { return s[rng.Intn(len(s))] }
	// Full most of the time, Partial with probability p.
	roll := func(p float64) Status {
		if rng.Float64() < p {
			return Partial
		}
		return Full
	}

	uniMsgs := []string{
		"Lectures: systems programming",
		"Operating systems lab",
		"Databases lecture + tutorial",
		"Algorithms lecture",
		"Worked on assignment",
		"Study group in the library",
	}
	uniPartial := []string{
		"Missed the morning lecture",
		"Left early, caught up at night",
		"Only made it to the lab",
	}
	examMsgs := []string{"Exam review", "Practice exams", "Final exam"}

	gymByDay := map[time.Weekday]string{
		time.Monday:    "Push day",
		time.Wednesday: "Pull day",
		time.Friday:    "Leg day",
		time.Sunday:    "Cardio + core",
	}
	gymPartial := []string{"Cut short, felt tired", "Warm-up and one lift only", "20 min session"}

	companies := []string{
		"Submitted application: backend developer co-op",
		"Submitted application: infrastructure intern",
		"Submitted application: data engineering role",
		"Submitted application: SWE new grad",
		"Submitted application: platform team intern",
	}
	appPartial := []string{
		"Drafted cover letter",
		"Tailored resume for posting",
		"Started application, finish tomorrow",
		"Prepped for coding screen",
	}

	volMsgs := []string{
		"Food bank shift",
		"Community garden cleanup",
		"Youth coding workshop",
		"Neighbourhood event setup",
	}
	volPartial := []string{"Helped for an hour", "Left early"}

	start := Date{Year: 2026, Month: 9, Date: 1}
	// Sep (30) + Oct (31) + Nov (30) = 91 days.
	for i := 0; i < 91; i++ {
		d := start.AddDate(0, 0, i)
		wd := d.Weekday()
		isWeekday := wd >= time.Monday && wd <= time.Friday
		examSeason := d.Month == 11 && d.Date >= 17

		// University: weekdays, with the odd light Sunday study session.
		switch {
		case isWeekday && examSeason:
			set(d, "uni", Full, pick(examMsgs))
		case isWeekday:
			if s := roll(0.15); s == Partial {
				set(d, "uni", s, pick(uniPartial))
			} else {
				set(d, "uni", s, pick(uniMsgs))
			}
		case wd == time.Sunday && rng.Float64() < 0.4:
			set(d, "uni", Partial, "Light study session")
		}

		// Work: Tue/Thu evenings, Saturday day shift.
		switch wd {
		case time.Tuesday, time.Thursday:
			if s := roll(0.1); s == Partial {
				set(d, "work", s, "Left early")
			} else {
				set(d, "work", s, "Evening shift 4pm-10pm")
			}
		case time.Saturday:
			if s := roll(0.1); s == Partial {
				set(d, "work", s, "Left early")
			} else {
				set(d, "work", s, "Day shift 9am-5pm")
			}
		}

		// Gym: Mon/Wed/Fri/Sun, ~20% of sessions skipped entirely.
		if msg, ok := gymByDay[wd]; ok && rng.Float64() > 0.2 {
			if s := roll(0.15); s == Partial {
				set(d, "gym", s, pick(gymPartial))
			} else {
				set(d, "gym", s, msg)
			}
		}

		// Applications: ~30% of weekdays.
		if isWeekday && rng.Float64() < 0.3 {
			if s := roll(0.4); s == Partial {
				set(d, "app", s, pick(appPartial))
			} else {
				set(d, "app", s, pick(companies))
			}
		}

		// Volunteering: most Sundays.
		if wd == time.Sunday && rng.Float64() < 0.6 {
			if s := roll(0.15); s == Partial {
				set(d, "vol", s, pick(volPartial))
			} else {
				set(d, "vol", s, pick(volMsgs))
			}
		}
	}

	return Calendar{
		Months: months,
		Labels: labels,
	}
}
