package timeinterval

import (
	"encoding/json"
	"fmt"
	"time"
)

// Interval describes an interval which can be bounded by a startsAt and/or an endsAt time.
// If startsAt is unset it will be interpreted as "unbounded" (goes infinitely long back in time).
// If endsAt is unset it will be interpreted as "unbounded" (goes infinitely long into the future).
// If both startsAt and endsAt is unset, then it will span all of time and be fairly pointless ;-)
type Interval struct {
	startsAt *time.Time
	endsAt   *time.Time
	duration *time.Duration
}

// UnmarshalJSON unmarshal Interval from an ISO8601 "interval" string.
func (in *Interval) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	i, err := ParseIntervalISO8601(s)
	if err != nil {
		return err
	}
	*in = *i
	return nil
}

// MarshalJSON marshal Interval into an ISO8601 "interval" string.
func (in Interval) MarshalJSON() ([]byte, error) {
	s, err := in.ISO8601()
	if err != nil {
		return nil, err
	}
	return json.Marshal(s)
}

// StartsAt returns the time the interval starts or nil if it does not have a lower bound.
func (in Interval) StartsAt() *time.Time {
	if in.StartsAtDerivedFromDuration() {
		startsAt := in.endsAt.Add(-*in.duration)
		return &startsAt
	}
	return in.startsAt
}

// EndsAt returns the time the interval ends or nil if it does not have an upper bound.
func (in Interval) EndsAt() *time.Time {
	if in.EndsAtDerivedFromDuration() {
		endsAt := in.startsAt.Add(*in.duration)
		return &endsAt
	}
	return in.endsAt
}

// Duration returns the duration of the interval or nil if it is unbounded.
func (in Interval) Duration() *time.Duration {
	if in.duration != nil {
		return in.duration
	}
	endsAt := in.EndsAt()
	startsAt := in.StartsAt()
	if startsAt == nil || endsAt == nil {
		return nil
	}
	d := endsAt.Sub(*startsAt)
	return &d
}

// Started returns a boolean indicating if the interval has begun at the given time.
func (in Interval) Started(t time.Time) bool {
	if in.startsAt == nil {
		return true
	}
	return in.startsAt.Before(t) || in.startsAt.Equal(t)
}

// Ended returns a boolean indicating if the interval has ended at the given time.
func (in Interval) Ended(t time.Time) bool {
	if in.endsAt == nil {
		return false
	}
	return in.endsAt.Before(t)
}

// In returns a boolean indicating if the given time is when the interval is active (Started and not Ended)
func (in Interval) In(t time.Time) bool {
	return in.Started(t) && !in.Ended(t)
}

// StartsAtDerivedFromDuration returns a boolean that indicates if Interval#StartsAt() is derived by
// the combination of EndsAt() and duration.
func (in Interval) StartsAtDerivedFromDuration() bool {
	return in.startsAt == nil && in.endsAt != nil && in.duration != nil
}

// EndsAtDerivedFromDuration returns a boolean that indicates if Interval#EndsAt() is derived by
// the combination of StartsAt() and duration.
func (in Interval) EndsAtDerivedFromDuration() bool {
	return in.endsAt == nil && in.startsAt != nil && in.duration != nil
}

// ISO8691 returns the interval formatted as an ISO8601 interval string.
// An error is returned if formatting fails.
func (in Interval) ISO8601() (string, error) {
	startsAt := in.StartsAt()
	endsAt := in.EndsAt()
	var startString string
	var endString string
	if in.StartsAtDerivedFromDuration() {
		s, err := durationToISO8601(*in.duration)
		if err != nil {
			return "", err
		}
		startString = s
		endString = endsAt.Format(time.RFC3339)
	} else if in.EndsAtDerivedFromDuration() {
		s, err := durationToISO8601(*in.duration)
		if err != nil {
			return "", err
		}
		endString = s
		startString = startsAt.Format(time.RFC3339)
	} else {
		startString = startsAt.Format(time.RFC3339)
		endString = endsAt.Format(time.RFC3339)
	}
	return fmt.Sprintf("%s/%s", startString, endString), nil
}
