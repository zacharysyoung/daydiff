package main

import (
	"errors"
	"fmt"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: daydiff date1 date2")
		os.Exit(2)
	}
	sa, sb := os.Args[1], os.Args[2]

	a, err := toTime(sa)
	if err != nil {
		fatalf("bad time: %s:", sa)
	}
	b, err := toTime(sb)
	if err != nil {
		fatalf("bad time: %s:", sb)
	}

	if a.Before(b) {
		a, b = b, a
	}
	i := 0
	for ; b.Before(a); i++ {
		b = b.Add(24 * time.Hour)
	}

	fmt.Printf("days=%d\n", i)
}

func toTime(s string) (time.Time, error) {
	var (
		t   time.Time
		err error

		year = fmt.Sprintf("%d", time.Now().Year())
	)

	t, err = time.Parse("2006/01/02", s)
	if err == nil {
		return t, nil
	}
	t, err = time.Parse("2006/01/02", year+"/"+s)
	if err == nil {
		return t, nil
	}

	t, err = time.Parse("1/2/2006", s)
	if err == nil {
		return t, nil
	}
	t, err = time.Parse("1/2/2006", s+"/"+year)
	if err == nil {
		return t, nil
	}

	return t, errors.New("could not parse")
}

func fatalf(format string, args ...any) {
	format = "error: " + format + "\n"
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}
