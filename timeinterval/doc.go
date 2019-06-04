/*
Package timeinterval provides an API for evaluating time intervals and for using repeating intervals.

Interval Example:

	startsAt := time.Now().Add(-1*time.Hour)
	endsAt := time.Now().Add(5*time.Hour)

	in := Interval{startsAt: &startsAt, endsAt: &endsAt}

	in.StartsAt() // => *time.Time
	in.EndsAt() // => *time.Time

	t := time.Now()
	in.Started(t) // => bool
	in.In(t) // => bool
	in.Ended(t) // => bool

Repeating Interval Example:

	duration := 15 * time.Minute
	repetitions := uint32(5)

	in := RepeatingInterval{
		Interval: Interval{
			startsAt: &startsAt,
			endsAt:   nil,
		},
		Duration:duration,
		Repetitions: &repetitions}

	in.StartsAt() // => *time.Time
	in.EndsAt() // => *time.Time

	t := time.Now()
	in.Started(t) // => bool
	in.In(t) // => bool
	in.Ended(t) // => bool
	in.Next(t) // => *time.Time
*/
package timeinterval
