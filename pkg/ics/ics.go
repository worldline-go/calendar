package ics

import (
	"fmt"
	"strings"
	"time"

	"github.com/worldline-go/calendar/pkg/models"
)

// GenerateICS generates an iCalendar (ICS) file content from a list of events.
func GenerateICS(events []models.Event) (string, error) {
	var b strings.Builder
	b.WriteString("BEGIN:VCALENDAR\r\n")
	b.WriteString("VERSION:2.0\r\n")
	b.WriteString("PRODID:-//worldline-go//calendar//EN\r\n")

	for _, e := range events {
		b.WriteString("BEGIN:VEVENT\r\n")
		b.WriteString(fmt.Sprintf("SUMMARY:%s\r\n", escapeICS(e.Name)))
		b.WriteString(fmt.Sprintf("DESCRIPTION:%s\r\n", escapeICS(e.Description)))

		from := e.DateFrom.Time
		to := e.DateTo.Time

		isAllDay := from.Hour() == 0 && from.Minute() == 0 && from.Second() == 0 &&
			to.Hour() == 0 && to.Minute() == 0 && to.Second() == 0 &&
			to.Sub(from) == 24*time.Hour

		if isAllDay {
			// All-day event: DTSTART/DTEND in DATE format (YYYYMMDD)
			b.WriteString(fmt.Sprintf("DTSTART;VALUE=DATE:%s\r\n", from.Format("20060102")))
			b.WriteString(fmt.Sprintf("DTEND;VALUE=DATE:%s\r\n", to.Format("20060102")))
		} else {
			// Timed event: include TZID if not UTC
			fromLoc, toLoc := from.Location(), to.Location()
			if fromLoc != time.UTC {
				b.WriteString(fmt.Sprintf("DTSTART;TZID=%s:%s\r\n", fromLoc.String(), from.Format("20060102T150405")))
			} else {
				b.WriteString(fmt.Sprintf("DTSTART:%s\r\n", from.UTC().Format("20060102T150405Z")))
			}
			if toLoc != time.UTC {
				b.WriteString(fmt.Sprintf("DTEND;TZID=%s:%s\r\n", toLoc.String(), to.Format("20060102T150405")))
			} else {
				b.WriteString(fmt.Sprintf("DTEND:%s\r\n", to.UTC().Format("20060102T150405Z")))
			}
		}

		if e.RRule != "" {
			b.WriteString(fmt.Sprintf("RRULE:%s\r\n", e.RRule))
		}
		b.WriteString("END:VEVENT\r\n")
	}

	b.WriteString("END:VCALENDAR\r\n")

	return b.String(), nil
}

// escapeICS escapes special characters for ICS fields
func escapeICS(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, ";", "\\;")
	s = strings.ReplaceAll(s, ",", "\\,")
	s = strings.ReplaceAll(s, "\n", "\\n")

	return s
}
