package ical

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
)

// https://datatracker.ietf.org/doc/html/rfc5545

// RRule represents a parsed RRULE according to RFC5545 section 3.3.10
// Only common fields are included for brevity, but can be extended.
type RRule struct {
	Freq       string
	Until      *time.Time
	Count      *int
	Interval   int
	BySecond   []int
	ByMinute   []int
	ByHour     []int
	ByDay      []string
	ByMonthDay []int
	ByYearDay  []int
	ByWeekNo   []int
	ByMonth    []int
	BySetPos   []int
	Wkst       string

	org string
}

func (r *RRule) Org() string {
	return r.org
}

// ParseRRule parses an RRULE string into an RRule struct.
func ParseRRule(s string) (*RRule, error) {
	rule := &RRule{Interval: 1, org: s}
	parts := strings.Split(s, ";")
	for _, part := range parts {
		if part == "" {
			continue
		}
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid RRULE part: %q", part)
		}
		key := strings.ToUpper(kv[0])
		val := kv[1]
		switch key {
		case "FREQ":
			rule.Freq = strings.ToUpper(val)
		case "UNTIL":
			t, err := parseTime(val)
			if err != nil {
				return nil, fmt.Errorf("invalid UNTIL: %w", err)
			}
			rule.Until = &t
		case "COUNT":
			count, err := strconv.Atoi(val)
			if err != nil {
				return nil, fmt.Errorf("invalid COUNT: %w", err)
			}
			rule.Count = &count
		case "INTERVAL":
			interval, err := strconv.Atoi(val)
			if err != nil {
				return nil, fmt.Errorf("invalid INTERVAL: %w", err)
			}
			rule.Interval = interval
		case "BYSECOND":
			rule.BySecond = parseIntList(val)
		case "BYMINUTE":
			rule.ByMinute = parseIntList(val)
		case "BYHOUR":
			rule.ByHour = parseIntList(val)
		case "BYDAY":
			rule.ByDay = strings.Split(val, ",")
		case "BYMONTHDAY":
			rule.ByMonthDay = parseIntList(val)
		case "BYYEARDAY":
			rule.ByYearDay = parseIntList(val)
		case "BYWEEKNO":
			rule.ByWeekNo = parseIntList(val)
		case "BYMONTH":
			rule.ByMonth = parseIntList(val)
		case "BYSETPOS":
			rule.BySetPos = parseIntList(val)
		case "WKST":
			rule.Wkst = strings.ToUpper(val)
		default:
			// ignore unknown keys for now
		}
	}
	return rule, nil
}

func parseIntList(s string) []int {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	res := make([]int, 0, len(parts))
	for _, p := range parts {
		i, err := strconv.Atoi(p)
		if err == nil {
			res = append(res, i)
		}
	}

	return res
}

// parseTime parses RFC5545 date-time (UTC or local)
func parseTime(s string) (time.Time, error) {
	// RFC5545 allows both DATE-TIME (with T and Z) and DATE (YYYYMMDD)
	if strings.HasSuffix(s, "Z") {
		return time.Parse("20060102T150405Z", s)
	}

	if len(s) == 8 {
		return time.Parse("20060102", s)
	}

	return time.Parse("20060102T150405", s)
}

// nextFreq returns the next occurrence for the given freq and interval
func nextFreq(t time.Time, freq string, interval int) time.Time {
	if interval < 1 {
		interval = 1
	}
	switch freq {
	case "DAILY":
		return t.AddDate(0, 0, interval)
	case "WEEKLY":
		return t.AddDate(0, 0, 7*interval)
	case "MONTHLY":
		return t.AddDate(0, interval, 0)
	case "YEARLY":
		return t.AddDate(interval, 0, 0)
	default:
		return t.AddDate(0, 0, interval)
	}
}

// MatchRRuleAt checks if the search time matches any occurrence of the RRule event.
// Returns the start and stop time of the matching occurrence, and true if found.
func MatchRRuleAt(rrule *RRule, dtstart, dtend, search time.Time) (time.Time, time.Time, bool) {
	if rrule == nil || rrule.Freq == "" {
		return time.Time{}, time.Time{}, false
	}
	start := dtstart
	count := 0
	maxCount := -1
	if rrule.Count != nil {
		maxCount = *rrule.Count
	}
	// Use a reasonable search window
	until := search.AddDate(10, 0, 0)
	if rrule.Until != nil && rrule.Until.Before(until) {
		until = *rrule.Until
	}
	occ := start
	duration := time.Duration(0)
	if !dtend.IsZero() {
		duration = dtend.Sub(dtstart)
	}

	for occ.Before(until) || occ.Equal(until) {
		// Generate all candidates for the current period (for BYSETPOS)
		candidates := generateCandidatesForPeriod(rrule, occ)
		// Apply BYSETPOS if present
		if len(rrule.BySetPos) > 0 {
			candidates = filterBySetPos(candidates, rrule.BySetPos)
		}
		for _, candidate := range candidates {
			occEnd := candidate
			if duration > 0 {
				occEnd = candidate.Add(duration)
			}
			// Check if search is within this occurrence
			if !search.Before(candidate) && search.Before(occEnd) {
				return candidate, occEnd, true
			}
			count++
			if maxCount > 0 && count >= maxCount {
				return time.Time{}, time.Time{}, false
			}
		}
		occ = nextFreq(occ, rrule.Freq, rrule.Interval)
	}

	return time.Time{}, time.Time{}, false
}

// MatchRRuleBetween returns the first occurrence (start and end) between dateFrom and dateTo, and true if found.
func MatchRRuleBetween(rrule *RRule, dtstart, dtend, dateFrom, dateTo time.Time) (time.Time, time.Time, bool) {
	if rrule == nil || rrule.Freq == "" {
		return time.Time{}, time.Time{}, false
	}

	duration := time.Duration(0)
	if !dtend.IsZero() {
		duration = dtend.Sub(dtstart)
	}

	start := dtstart
	count := 0
	maxCount := -1
	if rrule.Count != nil {
		maxCount = *rrule.Count
	}

	until := dateTo
	if rrule.Until != nil && rrule.Until.Before(until) {
		until = *rrule.Until
	}

	occ := start
	// Fast forward to the approximate first occurrence near dateFrom
	for occ.Before(dateFrom) {
		occ = nextFreq(occ, rrule.Freq, rrule.Interval)
		count++
		if maxCount > 0 && count >= maxCount {
			return time.Time{}, time.Time{}, false
		}
	}
	count = 0 // Reset count for actual iteration
	for occ.Before(until) || occ.Equal(until) {
		candidates := generateCandidatesForPeriod(rrule, occ)
		if len(rrule.BySetPos) > 0 {
			candidates = filterBySetPos(candidates, rrule.BySetPos)
		}
		for _, candidate := range candidates {
			occEnd := candidate
			if duration > 0 {
				occEnd = candidate.Add(duration)
			}
			// Check if this occurrence overlaps with our date range
			if (candidate.Before(dateTo) || candidate.Equal(dateTo)) &&
				(occEnd.After(dateFrom) || occEnd.Equal(dateFrom)) {
				return candidate, occEnd, true
			}
			count++
			if maxCount > 0 && count >= maxCount {
				return time.Time{}, time.Time{}, false
			}
		}
		occ = nextFreq(occ, rrule.Freq, rrule.Interval)
	}
	return time.Time{}, time.Time{}, false
}

// generateCandidatesForPeriod generates all possible candidates for the current period (e.g., week or month), applying BYxxx rules except BYSETPOS.
func generateCandidatesForPeriod(rrule *RRule, base time.Time) []time.Time {
	var candidates []time.Time
	freq := rrule.Freq
	wkst := parseWkst(rrule.Wkst)
	switch freq {
	case "SECONDLY":
		for i := range 60 {
			candidate := time.Date(base.Year(), base.Month(), base.Day(), base.Hour(), base.Minute(), i, 0, base.Location())
			if matchAllByRules(rrule, candidate) {
				candidates = append(candidates, candidate)
			}
		}
	case "MINUTELY":
		for i := range 60 {
			candidate := time.Date(base.Year(), base.Month(), base.Day(), base.Hour(), i, base.Second(), 0, base.Location())
			if matchAllByRules(rrule, candidate) {
				candidates = append(candidates, candidate)
			}
		}
	case "HOURLY":
		for i := range 24 {
			candidate := time.Date(base.Year(), base.Month(), base.Day(), i, base.Minute(), base.Second(), 0, base.Location())
			if matchAllByRules(rrule, candidate) {
				candidates = append(candidates, candidate)
			}
		}
	case "DAILY":
		if matchAllByRules(rrule, base) {
			candidates = append(candidates, base)
		}
	case "WEEKLY":
		startOfWeek := startOfWeek(base, wkst)
		for i := range 7 {
			candidate := startOfWeek.AddDate(0, 0, i)
			if candidate.Month() != base.Month() && base.Month() != 0 {
				continue
			}
			if matchAllByRules(rrule, candidate) {
				candidates = append(candidates, candidate)
			}
		}
	case "MONTHLY":
		if len(rrule.ByMonth) == 0 &&
			len(rrule.ByMonthDay) == 0 &&
			len(rrule.ByYearDay) == 0 &&
			len(rrule.ByWeekNo) == 0 &&
			len(rrule.ByDay) == 0 {
			if matchAllByRules(rrule, base) {
				candidates = append(candidates, base)
			}
			return candidates
		}

		first := time.Date(base.Year(), base.Month(), 1, base.Hour(), base.Minute(), base.Second(), base.Nanosecond(), base.Location())
		daysInMonth := daysInMonth(base.Year(), base.Month())
		for i := range daysInMonth {
			candidate := first.AddDate(0, 0, i)
			if matchAllByRules(rrule, candidate) {
				candidates = append(candidates, candidate)
			}
		}
	case "YEARLY":
		if len(rrule.ByMonth) == 0 &&
			len(rrule.ByMonthDay) == 0 &&
			len(rrule.ByYearDay) == 0 &&
			len(rrule.ByWeekNo) == 0 &&
			len(rrule.ByDay) == 0 {
			if matchAllByRules(rrule, base) {
				candidates = append(candidates, base)
			}
			return candidates
		}

		first := time.Date(base.Year(), 1, 1, base.Hour(), base.Minute(), base.Second(), base.Nanosecond(), base.Location())
		daysInYear := 365
		if isLeapYear(base.Year()) {
			daysInYear = 366
		}
		for i := range daysInYear {
			candidate := first.AddDate(0, 0, i)
			if matchAllByRules(rrule, candidate) {
				candidates = append(candidates, candidate)
			}
		}
	default:
		if matchAllByRules(rrule, base) {
			candidates = append(candidates, base)
		}
	}

	return candidates
}

func isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}

// filterBySetPos filters candidates by BYSETPOS (1-based, negative for from end)
func filterBySetPos(candidates []time.Time, setpos []int) []time.Time {
	var filtered []time.Time
	n := len(candidates)
	for _, pos := range setpos {
		idx := pos
		if pos > 0 {
			idx = pos - 1
		} else if pos < 0 {
			idx = n + pos
		}
		if idx >= 0 && idx < n {
			filtered = append(filtered, candidates[idx])
		}
	}

	return filtered
}

// parseWkst parses WKST (week start) string to time.Weekday, defaults to Monday
func parseWkst(wkst string) time.Weekday {
	switch strings.ToUpper(wkst) {
	case "SU":
		return time.Sunday
	case "MO":
		return time.Monday
	case "TU":
		return time.Tuesday
	case "WE":
		return time.Wednesday
	case "TH":
		return time.Thursday
	case "FR":
		return time.Friday
	case "SA":
		return time.Saturday
	default:
		return time.Monday
	}
}

// startOfWeek returns the start of the week for a given time and week start
func startOfWeek(t time.Time, wkst time.Weekday) time.Time {
	delta := (int(t.Weekday()) - int(wkst) + 7) % 7

	return t.AddDate(0, 0, -delta)
}

// daysInMonth returns the number of days in a month
func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// matchAllByRules checks all BYxxx rules except BYSETPOS for a candidate
func matchAllByRules(rrule *RRule, occ time.Time) bool {
	if len(rrule.BySecond) > 0 {
		if !slices.Contains(rrule.BySecond, occ.Second()) {
			return false
		}
	}
	if len(rrule.ByMinute) > 0 {
		if !slices.Contains(rrule.ByMinute, occ.Minute()) {
			return false
		}
	}
	if len(rrule.ByHour) > 0 {
		if !slices.Contains(rrule.ByHour, occ.Hour()) {
			return false
		}
	}
	if len(rrule.ByMonth) > 0 {
		if !slices.Contains(rrule.ByMonth, int(occ.Month())) {
			return false
		}
	}
	if len(rrule.ByMonthDay) > 0 {
		if !slices.Contains(rrule.ByMonthDay, occ.Day()) {
			return false
		}
	}
	if len(rrule.ByYearDay) > 0 {
		if !slices.Contains(rrule.ByYearDay, occ.YearDay()) {
			return false
		}
	}
	if len(rrule.ByWeekNo) > 0 {
		_, week := occ.ISOWeek()
		if !slices.Contains(rrule.ByWeekNo, week) {
			return false
		}
	}

	if len(rrule.ByDay) > 0 {
		found := false
		wday := occ.Weekday().String()[:2]
		for _, d := range rrule.ByDay {
			if len(d) > 2 {
				ord, day := d[:len(d)-2], d[len(d)-2:]
				if strings.EqualFold(day, wday) {
					ordInt, err := strconv.Atoi(ord)
					if err == nil && nthWeekdayOfMonth(occ, ordInt) {
						found = true
						break
					}
				}
			} else if strings.EqualFold(d, wday) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// nthWeekdayOfMonth checks if t is the nth weekday of its month (e.g., 2nd Monday, -1 Sunday)
func nthWeekdayOfMonth(t time.Time, n int) bool {
	weekday := t.Weekday()
	first := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	last := time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location())

	if n > 0 {
		count := 0
		for d := first; d.Month() == t.Month(); d = d.AddDate(0, 0, 1) {
			if d.Weekday() == weekday {
				count++
				if count == n && d.Day() == t.Day() {
					return true
				}
			}
		}
	} else if n < 0 {
		count := 0
		for d := last; d.Month() == t.Month(); d = d.AddDate(0, 0, -1) {
			if d.Weekday() == weekday {
				count--
				if count == n && d.Day() == t.Day() {
					return true
				}
			}
		}
	}

	return false
}
