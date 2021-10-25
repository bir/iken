package dates

import (
	"strings"
	"time"
)

func EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()

	return time.Date(year, month, day, 23, 59, 0, 0, t.Location())
}

func EndOfYear(t time.Time) time.Time {
	year, _, _ := t.Date()

	return time.Date(year, 12, 31, 23, 59, 0, 0, t.Location())
}

// ToTime parses an offset duration string.
// The calculation is based on an optional anchor (default is NOW) plus a duration.
// Valid Anchors:
//   EOD will calculate the end of day
//   EOY will calculate end of year
// A duration string is a number and unit suffix, such as "300m", "1.5h" or "2h45m".
// Valid time units are "m", "h".
// If location is provided, calculations are set to the corresponding time zone.
// Examples:
//   "EOD+72h" = end of day in 3 days
//   "3h" = 3 hours
//   "EOY" = end of year
//   "30m" = 30 minutes
//   "" = now
func ToTime(duration string, location *time.Location) (time.Time, error) {
	t := time.Now()
	if location != nil {
		t = t.In(location)
	}

	anchor := strings.ToUpper(duration)
	if len(duration) > 5 && duration[3] == '+' {
		anchor = duration[:3]
	}

	switch anchor {
	case "EOD":
		t = EndOfDay(t)
		duration = duration[3:]
	case "EOY":
		t = EndOfYear(t)
		duration = duration[3:]
	}

	if duration == "" {
		return t, nil
	}

	d, err := time.ParseDuration(duration)
	if err != nil {
		return time.Time{}, err // nolint
	}

	return t.Add(d), nil
}
