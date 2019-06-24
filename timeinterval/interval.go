package timeinterval

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Interval describes an interval bounded by a StartsAt and EndsAt time.
// the unexported "iso8601" is used to store the user's ISO8601 string. This makes it possible to marshal/unmarshal
// the interval to/from the same ISO8601 representation originally provided if desired.
type Interval struct {
	iso8601  string
	StartsAt time.Time
	EndsAt   time.Time
}

// NewInterval returns an Interval instance based upon the given options
// and an error if the given options are insufficient to construct the interval.
func NewInterval(startsAt, endsAt *time.Time, duration *time.Duration, iso8601 *string) (*Interval, error) {
	in := Interval{}
	if startsAt == nil && endsAt == nil {
		return nil, errors.New("invalid interval")
	}
	if startsAt != nil && endsAt != nil {
		in.StartsAt = *startsAt
		in.EndsAt = *endsAt
	} else {
		if duration == nil {
			return nil, errors.New("invalid interval")
		}
		if startsAt != nil {
			in.StartsAt = *startsAt
			in.EndsAt = in.StartsAt.Add(*duration)
		}
		if endsAt != nil {
			in.EndsAt = *endsAt
			in.StartsAt = in.EndsAt.Add(-*duration)
		}
	}
	if iso8601 == nil {
		in.iso8601 = in.ISO8601()
	} else {
		in.iso8601 = *iso8601
	}
	return &in, nil
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

// MarshalJSON marshals Interval into an ISO8601 "interval" string.
func (in Interval) MarshalJSON() ([]byte, error) {
	iso := in.iso8601
	if len(iso) < 1 {
		s := in.ISO8601()
		iso = s
	}
	return json.Marshal(in.iso8601)
}

// Duration returns the duration of the interval.
func (in Interval) Duration() time.Duration {
	return in.EndsAt.Sub(in.StartsAt)
}

// Started returns a boolean indicating if the interval has begun at the given time.
func (in Interval) Started(t time.Time) bool {
	return in.StartsAt.Before(t) || in.StartsAt.Equal(t)
}

// Ended returns a boolean indicating if the interval has ended at the given time.
func (in Interval) Ended(t time.Time) bool {
	return in.EndsAt.Before(t)
}

// In returns a boolean indicating if the given time is when the interval is active (Started and not Ended)
func (in Interval) In(t time.Time) bool {
	return in.Started(t) && !in.Ended(t)
}

// ISO8691 returns the interval formatted as an ISO8601 interval string.
func (in Interval) ISO8601() string {
	if len(in.iso8601) > 1 {
		return in.iso8601
	}
	return fmt.Sprintf("%s/%s", in.StartsAt.Format(time.RFC3339), in.EndsAt.Format(time.RFC3339))
}
