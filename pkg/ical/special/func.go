package special

import (
	"strings"
	"time"
)

var Funcs = map[string]func(int) time.Time{
	"GOODFRIDAY":   GoodFriday,
	"EASTERSUNDAY": EasterSunday,
	"EASTERMONDAY": EasterMonday,
	"ASCENSIONDAY": AscensionDay,
	"WHITSUNDAY":   WhitSunday,
	"WHITMONDAY":   WhitMonday,
}

func GetFunc(name string) (func(int) time.Time, bool) {
	fn, ok := Funcs[strings.ToUpper(name)]

	return fn, ok
}
