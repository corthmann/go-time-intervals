package timeinterval

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var regexTimeStringISO = regexp.MustCompile("^(-?(?:[1-9][0-9]*)?[0-9]{4})-(1[0-2]|0[1-9])-(3[01]|0[1-9]|[12][0-9])T(2[0-3]|[01][0-9]):([0-5][0-9]):([0-5][0-9])(\\.[0-9]+)?(Z)?$")

type formatType uint8

// TypeUnknown indicates that the given string has a format that is unknown and unsupported.
const TypeUnknown formatType = 0
// TypeTime indicates that the given string is an ISO8601 time string
const TypeTime formatType = 1
// TypeDuration indicates that the given string is am ISO8601 duration string
const TypeDuration formatType = 2

// TimeInterval is an interface that describes the API supported
// for all intervals implemented in the timeinterval package.
type TimeInterval interface {
	// StartsAt returns the time the interval begins or nil if it does not have a lower bound.
	StartsAt() *time.Time
	// EndsAt returns the time the interval ends or nil if it does not have an upper bound.
	EndsAt() *time.Time
	// Duration returns the duration of the interval or nil if it is unbounded.
	Duration() *time.Duration

	// Started returns a boolean indicating if the interval have begun at the given time.
	Started(t time.Time) bool
	// Ended returns a boolean indicating if the interval have ended at the given time.
	Ended(t time.Time) bool
	// In returns a boolean indicating if the given time is while the interval is active (Started and not Ended)
	In(t time.Time) bool
}

// RepeatingTimeInterval is an interface that describes the API supported by intervals with repeating events.
type RepeatingTimeInterval interface {
	// TimeInterval means that this interface supports the TimeInterval API.
	TimeInterval
	// Next returns the time of the next occurrence/event relative to the given time.
	// the StartsAt time is returned if the interval have not Started yet.
	// nil is returned if the interval has Ended.
	Next(t time.Time) *time.Time
}

// ParseIntervalISO8601 accepts a string with the ISO8601 "interval" format
// and returns an Interval and an error if parsing of the string failed.
// See: ref: https://en.wikipedia.org/wiki/ISO_8601#Time_intervals
func ParseIntervalISO8601(s string) (*Interval, error) {
	// Interval
	parts := strings.Split(s, "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid interval format")
	}
	partTypes, err := identifyIntervalTypes(parts)
	if err != nil {
		return nil, err
	}
	if partTypes[0] == TypeDuration &&  partTypes[1] == TypeDuration {
		return nil, errors.New("interval cannot consist of two durations")
	}

	in := Interval{}
	for i := 0; i < len(partTypes); i++ {
		switch partTypes[i] {
		case TypeDuration:
			d, err := parseDurationString(parts[i])
			if err != nil {
				return nil, err
			}
			in.duration = &d
		case TypeTime:
			t, err := parseTimeString(parts[i])
			if err != nil {
				return nil, err
			}
			if i == 0 {
				in.startsAt = &t
			} else {
				in.endsAt = &t
			}
		}
	}
	return &in, nil
}

// ParseRepeatingIntervalISO8601 accepts a string with the ISO8601 "repeating interval" format
// and returns a RepeatingInterval and an error if parsing of the string failed.
// See: ref: https://en.wikipedia.org/wiki/ISO_8601#Repeating_intervals
func ParseRepeatingIntervalISO8601(s string) (*RepeatingInterval, error) {
	if !strings.HasPrefix(s, "R") {
		return nil, errors.New("invalid repeating interval format")
	}
	ri := RepeatingInterval{}
	// Split the "Repetition" and "Interval" parts of the string.
	parts := strings.SplitN(s, "/", 2)
	repetitionString := parts[0]
	intervalString := parts[1]
	// Set "Repetitions"
	if len(repetitionString) > 1 {
		n, err := strconv.Atoi(repetitionString[1:])
		if err != nil {
			return nil, err
		}
		repetitions := uint32(n)
		ri.Repetitions = &repetitions
	}
	// Set "Interval"
	in, err := ParseIntervalISO8601(intervalString)
	if err != nil {
		return nil, err
	}
	ri.Interval = *in
	// Set "Duration"
	d := ri.Interval.Duration()
	if d != nil {
		ri.RepeatIn = *d
	}
	return &ri, nil
}

func identifyIntervalTypes(parts []string) ([]formatType, error) {
	types := make([]formatType, len(parts))
	for i := 0; i < len(parts); i++ {
		ft, err := identifyType(parts[i])
		if err != nil {
			return nil, err
		}
		if ft == TypeUnknown {
			return types, errors.New("unvalid interval format")
		}
		types[i] = ft
	}
	return types, nil
}

func identifyType(s string) (formatType, error) {
	if regexTimeStringISO.MatchString(s) {
		return TypeTime, nil
	}
	if strings.HasPrefix(s, "P") {
		return TypeDuration, nil
	}
	return TypeUnknown, errors.New("invalid/unknown format")
}

func parseTimeString(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

func parseDurationString(s string) (time.Duration, error) {
	d := time.Duration(0)
	if !strings.HasPrefix(s, "P") {
		return d, errors.New("invalid duration format")
	}
	// Exclude Duration indicator-char
	countStr := ""
	runes := []rune(s[1:])
	// Iterate runes and calculate Duration
	var currentCount = 1
	for i := 0; i < len(runes); i++ {
		c := string(runes[i])
		switch c {
		case "W":
			{
				countStr = ""
				d += time.Duration(currentCount) * 7 * 24 * time.Hour
			}
		case "D":
			{
				countStr = ""
				d += time.Duration(currentCount) * 24 * time.Hour
			}
		default:
			countStr += c
			// Calculate "Count"
			count, err := strconv.Atoi(countStr)
			if err != nil {
				return d, err
			}
			currentCount = count
		}
	}
	return d, nil
}
