package timeinterval

import (
	"time"
)

// RepeatingInterval describes an interval with recurring events at every Duration.
// The interval can be bounded by either:
// a fixed startsAt and endsAt
// or by a fixed startsAt with a fixed number of Repetitions from which the endsAt will be derived.
// or by a fixed endsAt with a fixed number of Repetitions from which the startsAt will be derived.
type RepeatingInterval struct {
	Interval TimeInterval
	Duration time.Duration
	Repetitions *uint32
}

// StartsAt returns the time the interval begins.
// When possible StartsAt will be derived using the Duration and Repetitions fields if Interval.StartsAt is unset.
func (in RepeatingInterval) StartsAt() *time.Time {
	if in.isStartsAtBoundedByRepetitions() {
		startsAt := in.Interval.EndsAt().Add(-time.Duration(*in.Repetitions)  * in.Duration)
		return &startsAt
	}
	return in.Interval.StartsAt()
}

// EndsAt returns the time the interval ends.
// When possible EndsAt will be derived using the Duration and Repetitions fields if Interval.EndsAt is unset.
func (in RepeatingInterval) EndsAt() *time.Time {
	if in.isEndsAtBoundedByRepetitions() {
		endsAt := in.Interval.StartsAt().Add(time.Duration(*in.Repetitions)  * in.Duration)
		return &endsAt
	}
	return in.Interval.EndsAt()
}

// Started returns a boolean indicating if the interval has begun at the given time.
func (in RepeatingInterval) Started(t time.Time) bool {
	if in.isStartsAtBoundedByRepetitions() {
		startsAt := in.StartsAt()
		return t.Equal(*startsAt) || t.After(*startsAt)
	}
	return in.Interval.Started(t)
}

// Ended returns a boolean indicating if the interval has ended at the given time.
func (in RepeatingInterval) Ended(t time.Time) bool {
	if in.isEndsAtBoundedByRepetitions() {
		endsAt := in.EndsAt()
		return t.After(*endsAt)
	}
	return in.Interval.Ended(t)
}

// In returns a boolean indicating if the given time is when the interval is active (Started and not Ended)
func (in RepeatingInterval) In(t time.Time) bool {
	return in.Started(t) && !in.Ended(t)
}

// Next returns the time of the next interval-occurrence relative to the given time.
// It returns the startsAt time if the interval have not started yet and nil if the interval has ended.
func (in RepeatingInterval) Next(t time.Time) *time.Time {
	if !in.Started(t) {
		return in.StartsAt()
	}
	if in.Ended(t) {
		return nil
	}
	startsAt := in.StartsAt()
	diff := t.Sub(*startsAt)
	mod := diff % in.Duration
	nxt := t.Add(in.Duration - mod)
	if in.Ended(nxt) {
		return nil
	}
	return &nxt
}

// isStartsAtBoundedByRepetitions returns a boolean which indicate if startsAt is unset
// and can be derived by using the Duration and Repetitions fields.
func (in RepeatingInterval) isStartsAtBoundedByRepetitions() bool {
	return in.Repetitions != nil && in.Interval.StartsAt() == nil && in.Interval.EndsAt() != nil
}

// isEndsAtBoundedByRepetitions returns a boolean which indicate if endsAt is unset
// and can be derived by using the Duration and Repetitions fields.
func (in RepeatingInterval) isEndsAtBoundedByRepetitions() bool {
	return in.Repetitions != nil && in.Interval.EndsAt() == nil && in.Interval.StartsAt() != nil
}
