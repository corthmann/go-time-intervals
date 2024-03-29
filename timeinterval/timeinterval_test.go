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
	duration := 7 * 24 * time.Hour
	expectations := map[string]Interval{
		"2019-01-02T21:00:00Z/2022-01-03T21:00:00Z": {StartsAt: startsAt, EndsAt: endsAt, Format: ISOFormatTimeAndTime},                     // Time - Time
		"2019-01-02T21:00:00Z/P1W":                  {StartsAt: startsAt, EndsAt: startsAt.Add(duration), Format: ISOFormatTimeAndDuration}, // Time - Duration
		"P1W/2022-01-03T21:00:00Z":                  {StartsAt: endsAt.Add(-duration), EndsAt: endsAt, Format: ISOFormatDurationAndTime},    // Duration - Time
	}
	for given, expected := range expectations {
		result, err := ParseIntervalISO8601(given)
		assert.Nil(t, err)
		assert.Equal(t, &expected, result)
		iso, err := result.ISO8601()
		assert.Nil(t, err)
		assert.Equal(t, given, iso)
	}
}

func TestParseRepeatingIntervalISO8601(t *testing.T) {
	startsAt, err := time.Parse(time.RFC3339, "2019-01-02T21:00:00Z")
	assert.Nil(t, err)
	endsAt, err := time.Parse(time.RFC3339, "2022-01-03T21:00:00Z")
	assert.Nil(t, err)
	duration := 7 * 24 * time.Hour
	repetitions := uint32(10)
	expectations := map[string]Repeating{
		"R/2019-01-02T21:00:00Z/2022-01-03T21:00:00Z": {
			Repetitions: nil,
			Interval:    Interval{StartsAt: startsAt, EndsAt: endsAt, Format: ISOFormatTimeAndTime},
		}, // Time - Time
		"R/2019-01-02T21:00:00Z/P1W": {
			Repetitions: nil,
			Interval:    Interval{StartsAt: startsAt, EndsAt: startsAt.Add(duration), Format: ISOFormatTimeAndDuration},
		}, // Time - Duration
		"R/P1W/2022-01-03T21:00:00Z": {
			Repetitions: nil,
			Interval:    Interval{StartsAt: endsAt.Add(-duration), EndsAt: endsAt, Format: ISOFormatDurationAndTime},
		}, // Duration - Time
		"R10/P1W/2022-01-03T21:00:00Z": {
			Repetitions: &repetitions,
			Interval:    Interval{StartsAt: endsAt.Add(-duration), EndsAt: endsAt, Format: ISOFormatDurationAndTime},
		}, // Duration - Time
	}
	for given, expected := range expectations {
		result, err := ParseRepeatingIntervalISO8601(given)
		assert.Nil(t, err)
		assert.Equal(t, &expected, result)
	}
}
