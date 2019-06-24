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

// typeUnknown indicates that the given string has a format that is unknown and unsupported.
const typeUnknown formatType = 0

// typeTime indicates that the given string is an ISO8601 time string
const typeTime formatType = 1

// typeDuration indicates that the given string is am ISO8601 duration string
const typeDuration formatType = 2

const durationWeek = 7 * 24 * time.Hour
const durationDay = 24 * time.Hour

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
	if partTypes[0] == typeDuration && partTypes[1] == typeDuration {
		return nil, errors.New("interval cannot consist of two durations")
	}
	var startsAt, endsAt *time.Time
	var duration *time.Duration
	for i := 0; i < len(partTypes); i++ {
		switch partTypes[i] {
		case typeDuration:
			d, err := parseDurationString(parts[i])
			if err != nil {
				return nil, err
			}
			duration = &d
		case typeTime:
			t, err := parseTimeString(parts[i])
			if err != nil {
				return nil, err
			}
			if i == 0 {
				startsAt = &t
			} else {
				endsAt = &t
			}
		}
	}
	return NewInterval(startsAt, endsAt, duration, &s)
}

// ParseRepeatingIntervalISO8601 accepts a string with the ISO8601 "repeating interval" format
// and returns a Repeating and an error if parsing of the string failed.
// See: ref: https://en.wikipedia.org/wiki/ISO_8601#Repeating_intervals
func ParseRepeatingIntervalISO8601(s string) (*Repeating, error) {
	if !strings.HasPrefix(s, "R") {
		return nil, errors.New("invalid repeating interval format")
	}
	ri := Repeating{}
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
	ri.RepeatEvery = d
	return &ri, nil
}

func identifyIntervalTypes(parts []string) ([]formatType, error) {
	types := make([]formatType, len(parts))
	for i := 0; i < len(parts); i++ {
		ft, err := identifyType(parts[i])
		if err != nil {
			return nil, err
		}
		if ft == typeUnknown {
			return types, errors.New("invalid interval format")
		}
		types[i] = ft
	}
	return types, nil
}

func identifyType(s string) (formatType, error) {
	if regexTimeStringISO.MatchString(s) {
		return typeTime, nil
	}
	if strings.HasPrefix(s, "P") {
		return typeDuration, nil
	}
	return typeUnknown, errors.New("invalid/unknown format")
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
				d += time.Duration(currentCount) * durationWeek
			}
		case "D":
			{
				countStr = ""
				d += time.Duration(currentCount) * durationDay
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
