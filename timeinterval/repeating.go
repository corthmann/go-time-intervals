package timeinterval

import (
	"encoding/json"
	"fmt"
	"time"
)

// Repeating describes an interval with recurring events distributed evenly by the duration of the interval.
// The number of Repetitions determine the bounds of the repeating interval (from StartsAt).
// When Repetitions is unset, then the repeating interval will be unbounded and recur infinitely long into the future.
type Repeating struct {
	Interval    Interval
	Repetitions *uint32
}

// String returns a string that describes the repeating interval.
func (r Repeating) String() string {
	if r.Repetitions != nil {
		return fmt.Sprintf("%v, reps: %v, times: %v", r.Interval, r.RepeatEvery(), *r.Repetitions)
	}
	return fmt.Sprintf("%v, reps: %v", r.Interval, r.RepeatEvery())
}

// RepeatEvery returns duration of each repetition. This is identical to the duration of the interval.
func (r Repeating) RepeatEvery() time.Duration {
	return r.Interval.Duration()
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
	iso, err := in.ISO8601()
	if err != nil {
		return nil, err
	}
	return json.Marshal(iso)
}

// StartsAt returns the time the interval begins.
// If "Repetitions" is nil, then this indicates the repeating interval is unbounded
// and as a result StartsAt() will return nil.
func (in Repeating) StartsAt() *time.Time {
	if in.Repetitions == nil {
		return nil
	}
	return &in.Interval.StartsAt
}

// EndsAt returns the time the interval ends.
// If "Repetitions" is nil, then this indicates the repeating interval is unbounded
// and as a result EndsAt() will return nil.
func (in Repeating) EndsAt() *time.Time {
	if in.Repetitions == nil {
		return nil
	}
	endsAt := in.Interval.StartsAt.Add(time.Duration(*in.Repetitions) * in.RepeatEvery())
	return &endsAt
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
// When the repeating interval is unbounded, then this function will always return true.
func (in Repeating) Started(t time.Time) bool {
	startsAt := in.StartsAt()
	if startsAt == nil {
		return true
	}
	return t.Equal(*startsAt) || t.After(*startsAt)
}

// Ended returns a boolean indicating if the interval has ended at the given time.
// When the repeating interval is unbounded, then this function will always return false.
func (in Repeating) Ended(t time.Time) bool {
	endsAt := in.EndsAt()
	if endsAt == nil {
		return false
	}
	return t.After(*endsAt)
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
	if in.Ended(t) || in.RepeatEvery() == 0 {
		return nil
	}
	diff := t.Sub(in.Interval.StartsAt)
	mod := diff % in.RepeatEvery()
	nxt := t.Add(in.RepeatEvery() - mod)
	if in.Ended(nxt) {
		return nil
	}
	return &nxt
}

// ISO8691 returns the repeating interval formatted as an ISO8601 repeating interval string.
func (in Repeating) ISO8601() (string, error) {
	iso, err := in.Interval.ISO8601()
	if err != nil {
		return "", err
	}
	if in.Repetitions != nil {
		return fmt.Sprintf("R%d/%s", *in.Repetitions, iso), nil
	}
	return fmt.Sprintf("R/%s", iso), nil
}
