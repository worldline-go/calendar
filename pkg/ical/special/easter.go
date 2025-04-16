package special

import "time"

// CalculateEasterDate uses the Butcher's algorithm to determine Easter Sunday for a given year.
func CalculateEasterDate(year int) time.Time {
	a := year % 19
	b := year / 100
	c := year % 100
	d := b / 4
	e := b % 4
	f := (b + 8) / 25
	g := (b - f + 1) / 3
	h := (19*a + b - d - g + 15) % 30
	i := c / 4
	k := c % 4
	l := (32 + 2*e + 2*i - h - k) % 7
	m := (a + 11*h + 22*l) / 451
	month := (h + l - 7*m + 114) / 31 // 3=March, 4=April
	day := ((h + l - 7*m + 114) % 31) + 1

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func GoodFriday(year int) time.Time {
	return CalculateEasterDate(year).AddDate(0, 0, -2)
}

func EasterSunday(year int) time.Time {
	return CalculateEasterDate(year)
}

func EasterMonday(year int) time.Time {
	return CalculateEasterDate(year).AddDate(0, 0, 1)
}

func AscensionDay(year int) time.Time {
	return CalculateEasterDate(year).AddDate(0, 0, 39)
}

func WhitSunday(year int) time.Time {
	return CalculateEasterDate(year).AddDate(0, 0, 49)
}

func WhitMonday(year int) time.Time {
	return CalculateEasterDate(year).AddDate(0, 0, 50)
}
