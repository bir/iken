package dates

import (
	"strings"
	"time"
)

var nowFunc = time.Now

func EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()

	return time.Date(year, month, day, 23, 59, 0, 0, t.Location())
}

func EndOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()

	first := time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
	last := first.AddDate(0, 1, -1)

	return time.Date(year, month, last.Day(), 23, 59, 0, 0, t.Location())
}

func EndOfYear(t time.Time) time.Time {
	year, _, _ := t.Date()

	return time.Date(year, 12, 31, 23, 59, 0, 0, t.Location())
}

// ToTime parses an offset duration string.
// The calculation is based on an optional anchor (default is NOW) plus an offset duration.
//
// shortcut for TimeToTime(time.Now()).
func ToTime(duration string, location *time.Location) (time.Time, error) {
	t := nowFunc()

	return TimeToTime(t, duration, location)
}

// TimeToTime parses an offset duration string.
// The calculation is based on an optional anchor (default is the passed in time) plus an offset duration.
// Valid Anchors:
//
//	EOD will calculate the end of day
//	EOM will calculate the end of month
//	EOY will calculate end of year
//
// A duration string is a number and unit suffix, such as "300m", "1.5h" or "2h45m".
// Valid time units are "m", "h".
// If location is provided, calculations are set to the corresponding time zone.
// Examples:
//
//	"EOD+72h" = end of day in 3 days
//	"3h" = 3 hours
//	"EOY" = end of year
//	"30m" = 30 minutes
//	"" = now
func TimeToTime(t time.Time, duration string, location *time.Location) (time.Time, error) {
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
	case "EOM":
		t = EndOfMonth(t)
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
		return time.Time{}, err //nolint
	}

	return t.Add(d), nil
}

// As returns the T as if it was read in the location originally.  It is similar
// to time.In except this function treats the input as zone less.  Input of 6pm
// will become 6pm in the new location.
func As(t time.Time, location *time.Location) time.Time {
	if t.Location() == location {
		return t
	}

	_, offset := t.Zone()
	out := t.In(location)
	_, offsetDest := out.Zone()

	return out.Add(time.Duration(offset-offsetDest) * time.Second)
}
