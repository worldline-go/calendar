package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/worldline-go/types"
)

// CheckYear checks if the given year is within the range defined by dateFrom and dateTo.
func CheckYear(year int, dateFrom, dateTo types.Null[types.Time], years string) (bool, error) {
	if years == "" {
		if dateFrom.Valid && dateFrom.V.Year() > year {
			return false, nil
		}

		if dateTo.Valid && dateTo.V.Year() < year {
			return false, nil
		}

		return true, nil
	}

	return includeYear(year, years)
}

func CheckDate(date types.Time, dateFrom, dateTo types.Null[types.Time], years string) (bool, error) {
	if years == "" {
		return checkDate(date, dateFrom, dateTo)
	}

	year := date.Year()
	ok, err := includeYear(year, years)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}

	var newDateFrom types.Null[types.Time]
	if dateFrom.Valid {
		newDateFrom = types.NewNull(types.Time{Time: ChangeYear(dateFrom.V.Time, year)})
	}

	var newDateTo types.Null[types.Time]
	if dateTo.Valid {
		newDateTo = types.NewNull(types.Time{Time: ChangeYear(dateTo.V.Time, year)})
	}

	return checkDate(date, newDateFrom, newDateTo)
}

func checkDate(date types.Time, dateFrom, dateTo types.Null[types.Time]) (bool, error) {
	if dateFrom.Valid && dateFrom.V.Time.After(date.Time) {
		return false, nil // Date is before the start date (included)
	}

	if dateTo.Valid && !date.Time.Before(dateTo.V.Time) {
		return false, nil // Date is on or after the end date (excluded)
	}

	return true, nil
}

func includeYear(year int, years string) (bool, error) {
	yearList := strings.Split(years, ",")
	for _, y := range yearList {
		// Case: wildcard only (any year is valid)
		if y == "*" {
			return true, nil
		}

		// Case: exact year match
		yearInt, err := strconv.Atoi(y)
		if err == nil && year == yearInt {
			return true, nil
		}

		// Case: year range with hyphen
		if strings.Contains(y, "-") {
			parts := strings.Split(y, "-")
			if len(parts) != 2 {
				return false, fmt.Errorf("invalid year range format: %s", y)
			}

			fromPart := parts[0]
			toPart := parts[1]

			// Process the start year
			var fromYear *int
			if fromPart != "*" {
				fromYearInt, err := strconv.Atoi(fromPart)
				if err != nil {
					return false, fmt.Errorf("invalid year format in range %s: %v", y, err)
				}
				if year < fromYearInt {
					continue // Year is below lower bound, check next pattern
				}

				fromYear = &fromYearInt
			}

			// Process the end year
			if toPart == "*" {
				// No upper bound, but year must be >= fromYear
				if fromYear == nil || year >= *fromYear {
					return true, nil
				}
			} else if len(toPart) > 1 && toPart[0] == '*' {
				// Handle pattern like "*4" - every 4 years starting from fromYear
				frequency, err := strconv.Atoi(toPart[1:])
				if err != nil {
					return false, fmt.Errorf("invalid frequency in range %s: %v", y, err)
				}

				// If fromYear is specified, check if the year is at the correct interval
				fromYearInt := 0
				if fromYear != nil {
					fromYearInt = *fromYear
				}

				if (year-fromYearInt)%frequency == 0 && (fromYear == nil || year >= *fromYear) {
					return true, nil
				}
			} else {
				// Regular upper bound
				toYear, err := strconv.Atoi(toPart)
				if err != nil {
					return false, fmt.Errorf("invalid year format in range %s: %v", y, err)
				}
				if year <= toYear && (fromYear == nil || year >= *fromYear) {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

// Function to change the year of a time.Time while preserving all other components
func ChangeYear(t time.Time, year int) time.Time {
	return time.Date(
		year,
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		t.Nanosecond(),
		t.Location(),
	)
}
