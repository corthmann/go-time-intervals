package timeinterval

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInterval_Started(t *testing.T) {
	startsAt := time.Now().Add(-1*time.Hour)
	endsAt := time.Now().Add(5*time.Hour)
	in := Interval{startsAt: &startsAt, endsAt: &endsAt}
	expectations := map[time.Time]bool{
		startsAt: true,
		startsAt.Add(-1* time.Hour): false,
		endsAt: true,
	}
	for given, expected := range expectations {
		assert.Equal(t, expected, in.Started(given))
	}
	assert.Equal(t, true, Interval{startsAt: nil}.Started(time.Now()))
}

func TestInterval_Ended(t *testing.T) {
	startsAt := time.Now().Add(-1*time.Hour)
	endsAt := time.Now().Add(5*time.Hour)
	in := Interval{startsAt: &startsAt, endsAt: &endsAt}
	expectations := map[time.Time]bool{
		startsAt: false,
		startsAt.Add(-1* time.Hour): false,
		endsAt: false,
		endsAt.Add(1*time.Hour): true,
	}
	for given, expected := range expectations {
		assert.Equal(t, expected, in.Ended(given))
	}
	assert.Equal(t, false, Interval{endsAt: nil}.Ended(time.Now()))
}

func TestInterval_In(t *testing.T) {
	startsAt := time.Now().Add(-1*time.Hour)
	endsAt := time.Now().Add(5*time.Hour)
	in := Interval{startsAt: &startsAt, endsAt: &endsAt}
	expectations := map[time.Time]bool{
		startsAt.Add(-1* time.Hour): false,
		startsAt: true,
		startsAt.Add(1*time.Hour): true,
		endsAt: true,
		endsAt.Add(1*time.Hour): false,
	}
	for given, expected := range expectations {
		assert.Equal(t, expected, in.In(given))
	}
}

func TestInterval_ISO8601(t *testing.T) {
	expectations := []string{
		"2019-01-02T21:00:00Z/2022-01-03T21:00:00Z",
		"2019-01-02T21:00:00Z/P1W",
		"P1W/2022-01-03T21:00:00Z",
	}
	for _, expectation := range expectations {
		in, err := ParseIntervalISO8601(expectation)
		assert.Nil(t, err)
		result, err := in.ISO8601()
		assert.Nil(t, err)
		assert.Equal(t, expectation, result)
	}
}
