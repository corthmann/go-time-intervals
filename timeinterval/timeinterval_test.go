package timeinterval

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseISO8601(t *testing.T) {
	startsAt, err := time.Parse(time.RFC3339, "2019-01-02T21:00:00Z")
	assert.Nil(t, err)
	endsAt, err := time.Parse(time.RFC3339, "2022-01-03T21:00:00Z")
	assert.Nil(t, err)
	duration := 7*24*time.Hour
	expectations := map[string]Interval{
		"2019-01-02T21:00:00Z/2022-01-03T21:00:00Z": {startsAt:&startsAt,endsAt:&endsAt}, // Time - Time
		"2019-01-02T21:00:00Z/P1W": {startsAt:&startsAt, endsAt:nil, duration:&duration}, // Time - Duration
		"P1W/2022-01-03T21:00:00Z": {startsAt:nil, duration:&duration, endsAt:&endsAt}, // Duration - Time
	}
	for given, expected := range expectations {
		result, err := ParseIntervalISO8601(given)
		assert.Nil(t, err)
		assert.Equal(t, &expected, result)
		assert.NotNil(t, result.StartsAt())
		assert.NotNil(t, result.EndsAt())
	}
}

func TestParseRepeatingIntervalISO8601(t *testing.T) {
	startsAt, err := time.Parse(time.RFC3339, "2019-01-02T21:00:00Z")
	assert.Nil(t, err)
	endsAt, err := time.Parse(time.RFC3339, "2022-01-03T21:00:00Z")
	assert.Nil(t, err)
	duration := 7*24*time.Hour
	repetitions := uint32(10)
	diff := endsAt.Sub(startsAt)
	expectations := map[string]RepeatingInterval{
		"R/2019-01-02T21:00:00Z/2022-01-03T21:00:00Z": {
			Repetitions: nil,
			RepeatEvery: diff,
			Interval:    Interval{ startsAt: &startsAt, endsAt: &endsAt, duration: nil},
		}, // Time - Time
		"R/2019-01-02T21:00:00Z/P1W": {
			Repetitions: nil,
			RepeatEvery: duration,
			Interval:    Interval{ startsAt: &startsAt, endsAt: nil, duration: &duration},
		}, // Time - Duration
		"R/P1W/2022-01-03T21:00:00Z": {
			Repetitions: nil,
			RepeatEvery: duration,
			Interval:    Interval{ startsAt: nil, endsAt: &endsAt, duration: &duration},
		}, // Duration - Time
		"R10/P1W/2022-01-03T21:00:00Z": {
			Repetitions: &repetitions,
			RepeatEvery: duration,
			Interval:    Interval{ startsAt: nil, endsAt: &endsAt, duration: &duration},
		}, // Duration - Time
	}
	for given, expected := range expectations {
		result, err := ParseRepeatingIntervalISO8601(given)
		assert.Nil(t, err)
		assert.Equal(t, &expected, result)
	}
}
