package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, `usage: daydiff date1 date2

Acceptable date formats:
  - January 2 2006
  - Jan 2 2006
  - 1/2/2006
  - 2006/01/02
  - 2006-01-02

and any of the above without the year.
`)
		os.Exit(2)
	}
	sa, sb := os.Args[1], os.Args[2]

	a, err := toTime(sa)
	if err != nil {
		fatalf("bad time %s: %v", sa, err)
	}
	b, err := toTime(sb)
	if err != nil {
		fatalf("bad time %s: %v", sb, err)
	}

	if a.Before(b) {
		a, b = b, a
	}

	layout := "Jan 2 2006"
	fmt.Println(b.Format(layout), "to", a.Format(layout))

	i := 0
	for ; b.Before(a); i++ {
		b = b.Add(24 * time.Hour)
	}

	fmt.Printf("%d days\n", i)
}

func toTime(s string) (time.Time, error) {
	var (
		t   time.Time
		err error

		parse = func(s string) error {
			for _, layout := range []string{
				"2006-01-02",
				"2006/01/02",
				"1/2/2006",
				"Jan 2 2006",
				"January 2 2006",
			} {
				t, err = time.Parse(layout, s)
				if err == nil {
					// log.Printf("parsed %s as %s", s, layout)
					return nil
				}
				// log.Printf("could not parse %s: %v", s, err)
			}
			return fmt.Errorf("could not parse %s", s)
		}
	)

	// try date s as-is
	if err = parse(s); err == nil {
		return t, nil
	}

	// that didn't work, try adding this year
	year := fmt.Sprintf("%d", time.Now().Year())

	for _, ss := range []string{
		year + "/" + s,
		year + "-" + s,
		s + " " + year,
	} {
		if err = parse(ss); err == nil {
			return t, nil
		}
	}

	return t, err
}

func fatalf(format string, args ...any) {
	format = "error: " + format + "\n"
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}
