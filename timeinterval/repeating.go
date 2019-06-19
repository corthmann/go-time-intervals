package timeinterval

import (
	"encoding/json"
	"fmt"
	"time"
)

// Repeating describes an interval with recurring events distributed evenly by a fixed duration.
// The interval can be bounded by either:
// a fixed startsAt and endsAt
// or by a fixed startsAt with a fixed number of Repetitions from which the endsAt will be derived.
// or by a fixed endsAt with a fixed number of Repetitions from which the startsAt will be derived.
type Repeating struct {
	Interval    Interval
	RepeatEvery time.Duration
	Repetitions *uint32
}

// UnmarshalJSON unmarshal Repeating from an ISO8601 "repeating interval" string.
func (in *Repeating) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	ri, err := ParseRepeatingIntervalISO8601(s)
	if err != nil {
		return err
	}
	*in = *ri
	return nil
}

// MarshalJSON marshal Repeating into an ISO8601 "repeating interval" string.
func (in Repeating) MarshalJSON() ([]byte, error) {
	s, err := in.ISO8601()
	if err != nil {
		return nil, err
	}
	return json.Marshal(s)
}

// StartsAt returns the time the interval begins.
// When possible StartsAt will be derived using the Duration and Repetitions fields if Interval.StartsAt is unset.
func (in Repeating) StartsAt() *time.Time {
	if in.isStartsAtBoundedByRepetitions() {
		startsAt := in.Interval.EndsAt().Add(-time.Duration(*in.Repetitions) * in.RepeatEvery)
		return &startsAt
	}
	return in.Interval.StartsAt()
}

// EndsAt returns the time the interval ends.
// When possible EndsAt will be derived using the Duration and Repetitions fields if Interval.EndsAt is unset.
func (in Repeating) EndsAt() *time.Time {
	if in.isEndsAtBoundedByRepetitions() {
		endsAt := in.Interval.StartsAt().Add(time.Duration(*in.Repetitions) * in.RepeatEvery)
		return &endsAt
	}
	return in.Interval.EndsAt()
}

// Duration returns the duration the repeating interval will be active for or nil if it is unbounded.
func (in Repeating) Duration() *time.Duration {
	endsAt := in.EndsAt()
	startsAt := in.StartsAt()
	if startsAt == nil || endsAt == nil {
		return nil
	}
	d := endsAt.Sub(*startsAt)
	return &d
}

// Started returns a boolean indicating if the interval has begun at the given time.
func (in Repeating) Started(t time.Time) bool {
	if in.isStartsAtBoundedByRepetitions() {
		startsAt := in.StartsAt()
		return t.Equal(*startsAt) || t.After(*startsAt)
	}
	return in.Interval.Started(t)
}

// Ended returns a boolean indicating if the interval has ended at the given time.
func (in Repeating) Ended(t time.Time) bool {
	if in.isEndsAtBoundedByRepetitions() {
		endsAt := in.EndsAt()
		return t.After(*endsAt)
	}
	return in.Interval.Ended(t)
}

// In returns a boolean indicating if the given time is when the interval is active (Started and not Ended)
func (in Repeating) In(t time.Time) bool {
	return in.Started(t) && !in.Ended(t)
}

// Next returns the time of the next interval-occurrence relative to the given time.
// It returns the startsAt time if the interval have not started yet and nil if the interval has ended.
func (in Repeating) Next(t time.Time) *time.Time {
	if !in.Started(t) {
		return in.StartsAt()
	}
	if in.Ended(t) || in.RepeatEvery == 0 {
		return nil
	}
	startsAt := in.StartsAt()
	diff := t.Sub(*startsAt)
	mod := diff % in.RepeatEvery
	nxt := t.Add(in.RepeatEvery - mod)
	if in.Ended(nxt) {
		return nil
	}
	return &nxt
}

// ISO8691 returns the repeating interval formatted as an ISO8601 repeating interval string.
// An error is returned if formatting fails.
func (in Repeating) ISO8601() (string, error) {
	startsAt := in.Interval.StartsAt()
	endsAt := in.Interval.EndsAt()
	var startString string
	var endString string
	if in.Interval.StartsAtDerivedFromDuration() {
		d := in.RepeatEvery
		s, err := durationToISO8601(d)
		startString = s
		if err != nil {
			return "", err
		}
		s = endsAt.Format(time.RFC3339)
		endString = s

	} else if in.Interval.EndsAtDerivedFromDuration() {
		d := in.RepeatEvery
		s, err := durationToISO8601(d)
		endString = s
		if err != nil {
			return "", err
		}
		s = startsAt.Format(time.RFC3339)
		startString = s
	} else {
		startString = startsAt.Format(time.RFC3339)
		endString = endsAt.Format(time.RFC3339)
	}
	if in.Repetitions != nil {
		return fmt.Sprintf("R%d/%s/%s", *in.Repetitions, startString, endString), nil
	}
	return fmt.Sprintf("R/%s/%s", startString, endString), nil
}

// isStartsAtBoundedByRepetitions returns a boolean which indicate if startsAt is unset
// and can be derived by using the Duration and Repetitions fields.
func (in Repeating) isStartsAtBoundedByRepetitions() bool {
	return in.Repetitions != nil && in.Interval.StartsAt() == nil && in.Interval.EndsAt() != nil
}

// isEndsAtBoundedByRepetitions returns a boolean which indicate if endsAt is unset
// and can be derived by using the Duration and Repetitions fields.
func (in Repeating) isEndsAtBoundedByRepetitions() bool {
	return in.Repetitions != nil && in.Interval.EndsAt() == nil && in.Interval.StartsAt() != nil
}
