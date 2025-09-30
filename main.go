package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func usage() {
	fmt.Fprintf(os.Stderr, `usage: timediff time1 time2

Acceptable time formats:
  - just time:   15:04(:05)
  - spelled out: Jan(uary) 2 ([2006|06]) (15:04(:05))
  - int'l:       ((20)06)-1-2 (15:04) | ((20)06)/1/2 (15:04(:05))
  - USA:         1/2((20)06) (15:04(:05))

which loosely translates to a time (hours, minutes, and optionally seconds), and/or a month/date, and optionally a year.
`)
	os.Exit(2)
}

func main() {
	if len(os.Args) != 3 {
		usage()
	}
	sa, sb := os.Args[1], os.Args[2]

	layouts := makeLayouts()

	a, err := toTime(sa, layouts)
	if err != nil {
		fatalf("bad time %s: %v", sa, err)
	}
	b, err := toTime(sb, layouts)
	if err != nil {
		fatalf("bad time %s: %v", sb, err)
	}

	if b.Before(a) {
		a, b = b, a
	}

	d := b.Sub(a)

	// layout := "Jan 2 2006 15:04:00"
	// fmt.Println(a.Format(layout), "to", b.Format(layout))

	i := 0
	for ; a.Before(b.Truncate(24 * time.Hour)); i++ {
		a = a.Add(24 * time.Hour)
	}

	var days string
	switch i {
	default:
		days = fmt.Sprintf("%d days", i)
	case 1:
		days = fmt.Sprintf("%d day", i)
	}

	switch i {
	case 0:
		fmt.Printf("%s\n", d)
	default:
		fmt.Printf("%s (%s)\n", d, days)
	}

}

func makeLayouts() []string {
	normSpaces := func(s string) string { return strings.Join(strings.Fields(s), " ") }

	const (
		yearLong  = "2006"
		yearShort = "06"

		monthLong  = "January"
		monthShort = "Jan"
		monthNum   = "1"

		day = "2"

		hourMin    = "15:04"
		hourMinSec = "15:04:05"
	)

	var (
		spelledOut []string // Jan(uary) 2 ([2006|'06]) (15:04)
		euro       []string // ((20)06)-1-2 (15:04) | ((20)06)/1/2 (15:04)
		usa        []string // 1/2((20)06) (15:04)

		layouts = []string{hourMin, hourMinSec}
	)

	for _, _m := range []string{monthShort, monthLong} {
		for _, _y := range []string{"", yearLong, yearShort} {
			for _, _t := range []string{"", hourMin, hourMinSec} {
				s := strings.Join([]string{_m, day, _y, _t}, " ")
				s = normSpaces(s)
				spelledOut = append(spelledOut, s)
			}
		}
	}

	layouts = append(layouts, spelledOut...)

	for _, sep := range []string{"-", "/"} {
		for _, _y := range []string{"", yearLong, yearShort} {
			for _, _t := range []string{"", hourMin, hourMinSec} {
				var s string
				switch _y {
				default:
					s = strings.Join([]string{_y, monthNum, day}, sep)
				case "":
					s = strings.Join([]string{monthNum, day}, sep)
				}
				s += " " + _t
				s = normSpaces(s)
				euro = append(euro, s)
			}
		}
	}

	layouts = append(layouts, euro...)

	for _, _y := range []string{"", yearLong, yearShort} {
		for _, _t := range []string{"", hourMin, hourMinSec} {
			var s string
			switch _y {
			default:
				s = strings.Join([]string{monthNum, day, _y}, "/")
			case "":
				s = strings.Join([]string{monthNum, day}, "/")
			}
			s += " " + _t
			s = normSpaces(s)
			usa = append(usa, s)
		}
	}

	layouts = append(layouts, usa...)

	return layouts
}

func toTime(s string, layouts []string) (time.Time, error) {
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("could not find layout that matched %s", s)
}

func fatalf(format string, args ...any) {
	format = "error: " + format + "\n"
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}
