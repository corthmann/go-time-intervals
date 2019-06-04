package timeinterval

import "time"

// Interval describes an interval which can be bounded by a startsAt and/or an endsAt time.
// If startsAt is unset it will be interpreted as "unbounded" (goes infinitely long back in time).
// If endsAt is unset it will be interpreted as "unbounded" (goes infinitely long into the future).
// If both startsAt and endsAt is unset, then it will span all of time and be fairly pointless ;-)
type Interval struct {
	startsAt *time.Time
	endsAt   *time.Time
}

// StartsAt returns the time the interval starts or nil if it does not have a lower bound.
func (in Interval) StartsAt() *time.Time {
	return in.startsAt
}

// EndsAt returns the time the interval ends or nil if it does not have an upper bound.
func (in Interval) EndsAt() *time.Time {
	return in.endsAt
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
