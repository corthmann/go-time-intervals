package timeinterval

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRepeating_StartsAt(t *testing.T) {
	duration := 15 * time.Minute
	repetitions := uint32(8)
	endsAt := time.Now().Add(1 * time.Hour)
	i, err := NewInterval(nil, &endsAt, &duration, nil)
	assert.Nil(t, err)
	in := Repeating{
		Interval:    *i,
		RepeatEvery: duration,
		Repetitions: &repetitions,
	}
	result := in.StartsAt()
	assert.NotNil(t, result)
	assert.Equal(t, endsAt.Add(-duration).Format(time.RFC3339), result.Format(time.RFC3339))
}

func TestRepeating_EndsAt(t *testing.T) {
	duration := 15 * time.Minute
	repetitions := uint32(8)
	startsAt := time.Now().Add(-1 * time.Hour)
	i, err := NewInterval(&startsAt, nil, &duration, nil)
	assert.Nil(t, err)
	in := Repeating{
		Interval:    *i,
		RepeatEvery: duration,
		Repetitions: &repetitions,
	}
	result := in.EndsAt()
	assert.NotNil(t, result)
	assert.Equal(t, startsAt.Add(2*time.Hour).Format(time.RFC3339), result.Format(time.RFC3339))
}

func TestRepeating_Next(t *testing.T) {
	startsAt := time.Now().Add(-1 * time.Hour)
	endsAt := time.Now().Add(5 * time.Hour)
	duration := endsAt.Sub(startsAt)
	repetitions := uint32(3)
	diff := endsAt.Sub(startsAt)
	i, err := NewInterval(&startsAt, &endsAt, nil, nil)
	assert.Nil(t, err)
	in := Repeating{
		Interval:    *i,
		RepeatEvery: duration,
		Repetitions: &repetitions,
	}
	expectations := map[time.Time]time.Time{
		startsAt.Add(-5 * time.Hour):           startsAt,
		startsAt:                               startsAt.Add(duration),
		startsAt.Add(7 * time.Minute):          startsAt.Add(duration),
		startsAt.Add(7*time.Minute + duration): startsAt.Add(2 * duration),
		endsAt.Add(-duration):                  startsAt.Add(diff - (diff % duration)),
	}
	for given, expected := range expectations {
		result := in.Next(given)
		assert.True(t, expected.Equal(*result))
	}
	assert.Nil(t, in.Next(startsAt.Add(time.Duration(repetitions)*duration)))
}

func TestRepeating_NextUnbounded(t *testing.T) {
	startsAt := time.Now().Add(-1 * time.Hour)
	endsAt := time.Now().Add(5 * time.Hour)
	duration := endsAt.Sub(startsAt)
	i, err := NewInterval(&startsAt, &endsAt, nil, nil)
	assert.Nil(t, err)
	in := Repeating{
		Interval:    *i,
		RepeatEvery: duration,
	}
	assert.Equal(t, startsAt.Add(duration), *in.Next(startsAt))
	assert.Equal(t, startsAt, *in.Next(startsAt.Add(-duration)))
	assert.Equal(t, startsAt.Add(-duration), *in.Next(startsAt.Add(-2 * duration)))
	assert.Equal(t, endsAt.Add(duration), *in.Next(endsAt))
}

func TestRepeating_Started(t *testing.T) {
	endsAt := time.Now().Add(-1 * time.Hour)

	duration := 15 * time.Minute
	repetitions := uint32(5)
	i, err := NewInterval(nil, &endsAt, &duration, nil)
	assert.Nil(t, err)
	in := Repeating{
		Interval:    *i,
		RepeatEvery: duration,
		Repetitions: &repetitions}

	assert.False(t, in.Started(i.EndsAt.Add(-time.Duration(repetitions+1)*duration)))
	assert.True(t, in.Started(i.StartsAt))
	in.Repetitions = nil
	assert.True(t, in.Started(endsAt.Add(-time.Duration(repetitions+1)*duration)))
}

func TestRepeating_Ended(t *testing.T) {
	startsAt := time.Now().Add(-1 * time.Hour)

	duration := 15 * time.Minute
	repetitions := uint32(5)
	i, err := NewInterval(&startsAt, nil, &duration, nil)
	assert.Nil(t, err)
	in := Repeating{
		Interval:    *i,
		RepeatEvery: duration,
		Repetitions: &repetitions}

	assert.True(t, in.Ended(startsAt.Add(time.Duration(repetitions+1)*duration)))
	assert.False(t, in.Ended(startsAt.Add(time.Duration(repetitions)*duration)))
	in.Repetitions = nil
	assert.False(t, in.Ended(startsAt.Add(time.Duration(repetitions+1)*duration)))
}

func TestRepeating_ISO8601(t *testing.T) {
	expectations := []string{
		"R/2019-01-02T21:00:00Z/2022-01-03T21:00:00Z",
		"R/2019-01-02T21:00:00Z/P1W",
		"R/P1W/2022-01-03T21:00:00Z",
		"R10/P1W/2022-01-03T21:00:00Z",
	}
	for _, expectation := range expectations {
		in, err := ParseRepeatingIntervalISO8601(expectation)
		assert.Nil(t, err)
		assert.Equal(t, expectation, in.ISO8601())
	}
}

func TestRepeating_MarshalJSON(t *testing.T) {
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

func TestRepeating_UnmarshalJSON(t *testing.T) {
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
		result := Repeating{}
		err = json.Unmarshal(b, &result)
		assert.Nil(t, err)
		assert.Equal(t, expected, &result)
	}
}
