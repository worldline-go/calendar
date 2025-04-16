package ical

import "time"

func TimeParse(format, v string, defaultTZ *time.Location) time.Time {
	t, _ := time.ParseInLocation(format, v, defaultTZ)

	return t
}
