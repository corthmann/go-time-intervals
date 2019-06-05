package timeinterval

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRepeatingInterval_Next(t *testing.T) {
	duration := 15 * time.Minute
	startsAt := time.Now().Add(-1*time.Hour)
	endsAt := time.Now().Add(5*time.Hour)
	diff := endsAt.Sub(startsAt)
	in := RepeatingInterval{
		Interval: Interval{
			startsAt: &startsAt,
			endsAt:   &endsAt,
		},
		RepeatIn: duration,
	}
	expectations := map[time.Time]time.Time{
		startsAt.Add(-5*time.Hour): startsAt,
		startsAt: startsAt.Add(duration),
		startsAt.Add(7 * time.Minute): startsAt.Add(duration),
		startsAt.Add(7 * time.Minute + duration): startsAt.Add(2* duration),
		endsAt.Add(-duration): startsAt.Add(diff - (diff % duration)),
	}
	for given, expected := range expectations {
		result := in.Next(given)
		assert.True(t, expected.Equal(*result))
	}
	assert.Nil(t, in.Next(endsAt))
}

func TestRepeatingInterval_NextWithoutStartsAt(t *testing.T) {
	duration := 15 * time.Minute
	repetitions := uint32(5)
	endsAt := time.Now().Add(5*time.Hour)
	in := RepeatingInterval{
		Interval: Interval{
			startsAt: nil,
			endsAt:   &endsAt,
		},
		RepeatIn: duration,
		Repetitions: &repetitions,
	}

	assert.Nil(t, in.Next(endsAt))
	assert.Equal(t, &endsAt,in.Next(endsAt.Add(-duration)))
	assert.Equal(t, endsAt.Add(-time.Duration(repetitions-1) * duration),*in.Next(endsAt.Add(-time.Duration(repetitions) * duration)))
	assert.Equal(t, endsAt.Add(-time.Duration(repetitions) * duration), *in.Next(endsAt.Add(-time.Duration(repetitions+1) * duration)))
}

func TestRepeatingInterval_Started(t *testing.T) {
	endsAt := time.Now().Add(-1*time.Hour)

	duration := 15 * time.Minute
	repetitions := uint32(5)
	in := RepeatingInterval{
		Interval: Interval{
			startsAt: nil,
			endsAt:   &endsAt,
		},
		RepeatIn:duration,
		Repetitions: &repetitions}


	assert.False(t, in.Started(endsAt.Add(-time.Duration(repetitions+1) * duration)))
	assert.True(t, in.Started(endsAt.Add(-time.Duration(repetitions) * duration)))
	in.Repetitions = nil
	assert.True(t, in.Started(endsAt.Add(-time.Duration(repetitions+1) * duration)))
}


func TestRepeatingInterval_Ended(t *testing.T) {
	startsAt := time.Now().Add(-1*time.Hour)

	duration := 15 * time.Minute
	repetitions := uint32(5)
	in := RepeatingInterval{
		Interval: Interval{
			startsAt: &startsAt,
			endsAt:   nil,
		},
		RepeatIn:duration,
		Repetitions: &repetitions}


	assert.True(t, in.Ended(startsAt.Add(time.Duration(repetitions+1) * duration)))
	assert.False(t, in.Ended(startsAt.Add(time.Duration(repetitions) * duration)))
	in.Repetitions = nil
	assert.False(t, in.Ended(startsAt.Add(time.Duration(repetitions+1) * duration)))
}

func TestRepeatingInterval_ISO8601(t *testing.T) {
	expectations := []string{
		"R/2019-01-02T21:00:00Z/2022-01-03T21:00:00Z",
		"R/2019-01-02T21:00:00Z/P1W",
		"R/P1W/2022-01-03T21:00:00Z",
		"R10/P1W/2022-01-03T21:00:00Z",
	}
	for _, expectation := range expectations {
		in, err := ParseRepeatingIntervalISO8601(expectation)
		assert.Nil(t, err)
		result, err := in.ISO8601()
		assert.Nil(t, err)
		assert.Equal(t, expectation, result)
	}
}

func TestRepeatingInterval_MarshalJSON(t *testing.T) {
	expectations := []string{
		"R/2019-01-02T21:00:00Z/2022-01-03T21:00:00Z",
		"R/2019-01-02T21:00:00Z/P1W",
		"R/P1W/2022-01-03T21:00:00Z",
		"R10/P1W/2022-01-03T21:00:00Z",
	}
	for _, expected := range expectations {
		// Parse & Marshal interval
		in, err := ParseRepeatingIntervalISO8601(expected)
		assert.Nil(t, err)
		b, err := json.Marshal(in)
		assert.Nil(t, err)
		// Unqoute result and compare to input
		result, err := strconv.Unquote(string(b))
		assert.Nil(t, err)
		assert.Equal(t, expected, result)
	}
}

func TestRepeatingInterval_UnmarshalJSON(t *testing.T) {
	expectations := []string{
		"R/2019-01-02T21:00:00Z/2022-01-03T21:00:00Z",
		"R/2019-01-02T21:00:00Z/P1W",
		"R/P1W/2022-01-03T21:00:00Z",
		"R10/P1W/2022-01-03T21:00:00Z",
	}
	for _, input := range expectations {
		// Parse & Marshal interval
		expected, err := ParseRepeatingIntervalISO8601(input)
		assert.Nil(t, err)
		b, err := json.Marshal(expected)
		assert.Nil(t, err)
		// Unmarshal and evaluate the result
		result := RepeatingInterval{}
		err = json.Unmarshal(b, &result)
		assert.Nil(t, err)
		assert.Equal(t, expected, &result)
	}
}
