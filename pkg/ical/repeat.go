package ical

import (
	"fmt"
	"strings"
	"time"

	"github.com/worldline-go/calendar/pkg/ical/special"
)

type Repeat struct {
	RRule []*RRule
	Func  []func(int) time.Time
}

// ParseRepeat parses a repeat string and returns a Repeat struct.
// The repeat string can be in the format of "RRULE:FREQ=DAILY;INTERVAL=1" or "FUNC:GoodFriday" or both with space/new line.
func ParseRepeat(rruleStr string) (*Repeat, error) {
	var rrule Repeat
	// Split the string by space or new line
	parts := strings.Fields(rruleStr)
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty repeat string")
	}

	for _, part := range parts {
		if strings.HasPrefix(part, "RRULE:") {
			rruleStr := strings.TrimPrefix(part, "RRULE:")
			rule, err := ParseRRule(rruleStr)
			if err != nil {
				return nil, fmt.Errorf("failed to parse rrule: %w", err)
			}
			rrule.RRule = append(rrule.RRule, rule)
		} else if strings.HasPrefix(part, "FUNC:") {
			funcName := strings.TrimPrefix(part, "FUNC:")
			if fn, ok := special.GetFunc(funcName); ok {
				rrule.Func = append(rrule.Func, fn)
			} else {
				return nil, fmt.Errorf("unknown function: %s", funcName)
			}
		} else {
			return nil, fmt.Errorf("invalid repeat string format")
		}
	}

	return &rrule, nil
}
