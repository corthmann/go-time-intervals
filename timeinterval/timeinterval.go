package timeinterval

import "time"

// TimeInterval is an interface that describes the API supported
// for all intervals implemented in the timeinterval package.
type TimeInterval interface {
	// StartsAt returns the time the interval begins or nil if it does not have a lower bound.
	StartsAt() *time.Time
	// EndsAt returns the time the interval ends or nil if it does not have an upper bound.
	EndsAt() *time.Time

	// Started returns a boolean indicating if the interval have begun at the given time.
	Started(t time.Time) bool
	// Ended returns a boolean indicating if the interval have ended at the given time.
	Ended(t time.Time) bool
	// In returns a boolean indicating if the given time is while the interval is active (Started and not Ended)
	In(t time.Time) bool
}

// RepeatingTimeInterval is an interface that describes the API supported by intervals with repeating events.
type RepeatingTimeInterval interface {
	// TimeInterval means that this interface supports the TimeInterval API.
	TimeInterval
	// Next returns the time of the next occurrence/event relative to the given time.
	// the StartsAt time is returned if the interval have not Started yet.
	// nil is returned if the interval has Ended.
	Next(t time.Time) *time.Time
}
