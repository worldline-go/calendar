package ical

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/worldline-go/calendar/pkg/models"
	"github.com/worldline-go/types"
)

// GenerateICS generates an iCalendar (ICS) file content from a list of events.
func GenerateICS(events []models.Event, category string) (string, error) {
	var b strings.Builder
	b.WriteString("BEGIN:VCALENDAR\r\n")
	b.WriteString("VERSION:2.0\r\n")
	b.WriteString("PRODID:-//worldline-go//calendar//EN\r\n")

	if category == "" {
		category = "Holidays"
	}

	for _, e := range events {
		b.WriteString("BEGIN:VEVENT\r\n")
		b.WriteString(fmt.Sprintf("UID:%s\r\n", e.ID))
		b.WriteString(fmt.Sprintf("CATEGORIES:%s\r\n", category))
		b.WriteString("CLASS:PUBLIC\r\n")

		name := escapeICS(e.Name)
		if strings.HasPrefix(name, "LANGUAGE=") {
			b.WriteString(fmt.Sprintf("SUMMARY;%s\r\n", name))
		} else {
			b.WriteString(fmt.Sprintf("SUMMARY:%s\r\n", name))
		}

		description := escapeICS(e.Description)
		if description != "" {
			b.WriteString(fmt.Sprintf("DESCRIPTION:%s\r\n", description))
		}

		from := e.DateFrom.Time
		to := e.DateTo.Time

		isAllDay := false
		if e.AllDay {
			isAllDay = from.Hour() == 0 && from.Minute() == 0 && from.Second() == 0 &&
				to.Hour() == 0 && to.Minute() == 0 && to.Second() == 0 &&
				to.Sub(from) == 24*time.Hour
		}

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
		b.WriteString("TRANSP:TRANSPARENT\r\n")
		b.WriteString("END:VEVENT\r\n")
	}

	b.WriteString("END:VCALENDAR\r\n")

	return b.String(), nil
}

// ParseICS parses ICS file data and returns a slice of models.Event.
func ParseICS(data io.Reader, tz *time.Location) ([]models.Event, error) {
	defaultTZ := time.UTC
	if tz != nil {
		defaultTZ = tz
	}

	reader := bufio.NewReader(data)
	var events []models.Event
	var e models.Event
	inEvent := false

	var current string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return nil, fmt.Errorf("failed to read line: %w", err)
			}

			break
		}

		if strings.HasPrefix(line, " ") {
			switch current {
			case "DESCRIPTION":
				e.Description += unescapeICS(strings.TrimSpace(line))
			case "SUMMARY":
				e.Name += unescapeICS(strings.TrimSpace(line))
			}

			continue
		}

		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		if line == "BEGIN:VEVENT" {
			inEvent = true
			e = models.Event{}
			current = ""

			continue
		}
		if line == "END:VEVENT" && inEvent {
			inEvent = false
			e.Tz = defaultTZ.String()
			if e.DateTo.Time.IsZero() {
				e.DateTo = types.Time{Time: e.DateFrom.AddDate(0, 0, 1)}
			}

			events = append(events, e)
			current = ""

			continue
		}
		if !inEvent {
			current = ""

			continue
		}

		if strings.HasPrefix(line, "UID:") {
			e.ID = strings.TrimPrefix(line, "UID:")
			current = "UID"
		} else if strings.HasPrefix(line, "SUMMARY:") {
			e.Name = unescapeICS(strings.TrimPrefix(line, "SUMMARY:"))
			current = "SUMMARY"
		} else if strings.HasPrefix(line, "SUMMARY;") {
			e.Name = unescapeICS(strings.TrimPrefix(line, "SUMMARY;"))
			current = "SUMMARY"
		} else if strings.HasPrefix(line, "DESCRIPTION:") {
			e.Description = unescapeICS(strings.TrimPrefix(line, "DESCRIPTION:"))
			current = "DESCRIPTION"
		} else if strings.HasPrefix(line, "DTSTART") {
			current = ""
			v := line[strings.Index(line, ":")+1:]
			if strings.Contains(line, ";VALUE=DATE") {
				e.DateFrom.Time = TimeParse("20060102", v, defaultTZ)
				e.AllDay = true
			} else if strings.Contains(line, "TZID=") {
				tzidStart := strings.Index(line, "TZID=") + len("TZID=")
				tzidEnd := strings.Index(line, ":")
				if tzidEnd > tzidStart {
					tzid := line[tzidStart:tzidEnd]
					loc, err := time.LoadLocation(tzid)
					if err == nil {
						e.DateFrom.Time = TimeParse("20060102T150405", v, loc)
					} else {
						// Fallback to parsing as local if TZID is invalid
						e.DateFrom.Time = TimeParse("20060102T150405", v, defaultTZ)
					}
				} else {
					// Fallback to parsing as local if TZID is malformed
					e.DateFrom.Time = TimeParse("20060102T150405", v, defaultTZ)
				}
			} else {
				e.DateFrom.Time = TimeParse("20060102T150405Z", v, defaultTZ)
			}
		} else if strings.HasPrefix(line, "DTEND") {
			current = ""
			v := line[strings.Index(line, ":")+1:]
			if strings.Contains(line, ";VALUE=DATE") {
				e.DateTo.Time = TimeParse("20060102", v, defaultTZ)
				e.AllDay = true
			} else if strings.Contains(line, "TZID=") {
				tzidStart := strings.Index(line, "TZID=") + len("TZID=")
				tzidEnd := strings.Index(line, ":")
				if tzidEnd > tzidStart {
					tzid := line[tzidStart:tzidEnd]
					loc, err := time.LoadLocation(tzid)
					if err == nil {
						e.DateTo.Time = TimeParse("20060102T150405", v, loc)
					} else {
						// Fallback to parsing as local if TZID is invalid
						e.DateTo.Time = TimeParse("20060102T150405", v, defaultTZ)
					}
				} else {
					// Fallback to parsing as local if TZID is malformed
					e.DateTo.Time = TimeParse("20060102T150405", v, defaultTZ)
				}
			} else {
				e.DateTo.Time = TimeParse("20060102T150405Z", v, defaultTZ)
			}
		} else if strings.HasPrefix(line, "RRULE:") {
			e.RRule = line
			current = ""
		}
	}

	return events, nil
}

// escapeICS escapes special characters for ICS fields
func escapeICS(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, ";", "\\;")
	s = strings.ReplaceAll(s, ",", "\\,")
	s = strings.ReplaceAll(s, "\n", "\\n")

	return s
}

// unescapeICS reverses escapeICS for ICS fields
func unescapeICS(s string) string {
	s = strings.ReplaceAll(s, "\\n", "\n")
	s = strings.ReplaceAll(s, "\\,", ",")
	s = strings.ReplaceAll(s, "\\;", ";")
	s = strings.ReplaceAll(s, "\\\\", "\\")
	return s
}
