package util

import "time"

// IsWeekend checks if time.Time object is in the weekend.
func IsWeekend(date time.Time) bool {
	return date.Weekday() == time.Saturday || date.Weekday() == time.Sunday
}
