package calender

import (
	"math"
	"time"
)

// Style Related

// Time related
func getFirstDay() time.Weekday {
	firstDay := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location())
	return firstDay.Weekday()
}

func getLastDay() time.Weekday {
	lastDay := time.Date(time.Now().Year(), time.Now().Month()+1, 0, 0, 0, 0, 0, time.Now().Location())
	return lastDay.Weekday()
}

func getNumOfDays() int {
	return time.Date(time.Now().Year(), time.Now().Month()+1, 0, 0, 0, 0, 0, time.Now().Location()).Day()
}

func getNumOfRows() int {
	totalDays := float64(getNumOfDays() + int(getFirstDay()))
	return int(math.Ceil(totalDays / 7.0))
}

func getTotalNumOfDays() int {
	return int(getFirstDay()) + getNumOfDays() + int(getLastDay()) - 1
}
