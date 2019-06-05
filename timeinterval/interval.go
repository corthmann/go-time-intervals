package timeinterval

import "time"

// Interval describes an interval which can be bounded by a startsAt and/or an endsAt time.
// If startsAt is unset it will be interpreted as "unbounded" (goes infinitely long back in time).
// If endsAt is unset it will be interpreted as "unbounded" (goes infinitely long into the future).
// If both startsAt and endsAt is unset, then it will span all of time and be fairly pointless ;-)
type Interval struct {
	startsAt *time.Time
	endsAt   *time.Time
	duration *time.Duration
}

// StartsAt returns the time the interval starts or nil if it does not have a lower bound.
func (in Interval) StartsAt() *time.Time {
	if in.startsAt == nil && in.endsAt != nil && in.duration != nil {
		startsAt := in.endsAt.Add(-*in.duration)
		return &startsAt
	}
	return in.startsAt
}

// EndsAt returns the time the interval ends or nil if it does not have an upper bound.
func (in Interval) EndsAt() *time.Time {
	if in.endsAt == nil && in.startsAt != nil && in.duration != nil {
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
