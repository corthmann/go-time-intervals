package timeinterval

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type isoFormat uint8

// ISOFormatUnknown indicates that the isoFormat is unset
const ISOFormatUnknown isoFormat = 0

// ISOFormatTimeAndTime means the interval.ISO8601() output will have the format Time/Time.
const ISOFormatTimeAndTime isoFormat = 1

// ISOFormatTimeAndDuration means the interval.ISO8601() output will have the format Time/Duration.
const ISOFormatTimeAndDuration isoFormat = 2

// ISOFormatTimeAndDuration means the interval.ISO8601() output will have the format Duration/Time.
const ISOFormatDurationAndTime isoFormat = 3

// Interval describes an interval bounded by a StartsAt and EndsAt time.
// the unexported "iso8601" is used to store the user's ISO8601 string. This makes it possible to marshal/unmarshal
// the interval to/from the same ISO8601 representation originally provided if desired.
type Interval struct {
	Format   isoFormat
	StartsAt time.Time
	EndsAt   time.Time
}

// NewInterval returns an Interval instance with set StartsAt, EndsAt and Format fields
// and an error if:
//
// 1) both the given startsAt and endsAt are nil
// 2) the given duration is nil while either startsAt or endsAt are nil
// 3) the resulting interval is invalid. See: Interval#Validate()
func NewInterval(startsAt, endsAt *time.Time, duration *time.Duration) (*Interval, error) {
	in := Interval{}
	if startsAt == nil && endsAt == nil {
		return nil, errors.New("invalid interval")
	}
	if startsAt != nil && endsAt != nil {
		in.StartsAt = *startsAt
		in.EndsAt = *endsAt
		in.Format = ISOFormatTimeAndTime
	} else {
		if duration == nil {
			return nil, errors.New("invalid interval")
		}
		if startsAt != nil {
			in.StartsAt = *startsAt
			in.EndsAt = in.StartsAt.Add(*duration)
			in.Format = ISOFormatTimeAndDuration
		}
		if endsAt != nil {
			in.EndsAt = *endsAt
			in.StartsAt = in.EndsAt.Add(-*duration)
			in.Format = ISOFormatDurationAndTime
		}
	}
	return &in, in.Validate()
}

// Validate verifies that validity of the interval and returns an error if the:
//
// 1) EndsAt time is before the StartsAt time
// 2) ISO8601 output format is unset
func (in Interval) Validate() error {
	if in.EndsAt.Before(in.StartsAt) {
		return errors.New("interval must start before it ends")
	}
	if in.Format == ISOFormatUnknown {
		return errors.New("unknown ISO8601 output format")
	}
	return nil
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

// String returns a string that describes of the interval.
func (in Interval) String() string {
	return fmt.Sprintf("%v -> %v", in.StartsAt, in.EndsAt)
}

// MarshalJSON marshals Interval into an ISO8601 "interval" string.
func (in Interval) MarshalJSON() ([]byte, error) {
	s, err := in.ISO8601()
	if err != nil {
		return nil, err
	}
	return json.Marshal(s)
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
func (in Interval) ISO8601() (string, error) {
	switch in.Format {
	case ISOFormatDurationAndTime:
		d, err := durationToISO8601(in.Duration())
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s/%s", d, in.EndsAt.Format(time.RFC3339)), nil
	case ISOFormatTimeAndDuration:
		d, err := durationToISO8601(in.Duration())
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s/%s", in.StartsAt.Format(time.RFC3339), d), nil
	default:
		return fmt.Sprintf("%s/%s", in.StartsAt.Format(time.RFC3339), in.EndsAt.Format(time.RFC3339)), nil
	}
}
