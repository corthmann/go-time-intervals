package timeinterval_test

import (
	"fmt"
	"time"

	"github.com/corthmann/go-time-intervals/timeinterval"
)

func ExampleInterval() {
	now, err := time.Parse(time.RFC3339, "2019-01-02T21:00:00Z")
	fmt.Println(err)
	startsAt := now.Add(-1 * time.Hour)
	endsAt := now.Add(5 * time.Hour)

	in, err := timeinterval.NewInterval(&startsAt, &endsAt, nil, nil)
	fmt.Println(err)
	fmt.Println(in.StartsAt.Format(time.RFC3339))
	fmt.Println(in.EndsAt.Format(time.RFC3339))

	fmt.Println(in.Started(startsAt))
	fmt.Println(in.Started(now))
	fmt.Println(in.In(now))
	fmt.Println(in.In(endsAt))
	fmt.Println(in.Ended(now))
	fmt.Println(in.Ended(endsAt))

	// Output:
	// <nil>
	// <nil>
	// 2019-01-02T20:00:00Z
	// 2019-01-03T02:00:00Z
	// true
	// true
	// true
	// true
	// false
	// false
}

func ExampleRepeating() {
	now, err := time.Parse(time.RFC3339, "2019-01-02T21:00:00Z")
	fmt.Println(err)
	startsAt := now.Add(-15 * time.Minute)
	duration := 15 * time.Minute
	repetitions := uint32(5)

	i, err := timeinterval.NewInterval(&startsAt, nil, &duration, nil)
	fmt.Println(err)
	in := timeinterval.Repeating{
		Interval:    *i,
		Repetitions: &repetitions}

	fmt.Println(in.StartsAt().Format(time.RFC3339))
	fmt.Println(in.EndsAt().Format(time.RFC3339))

	fmt.Println(in.Started(startsAt))
	fmt.Println(in.Started(now))
	fmt.Println(in.In(now))
	fmt.Println(in.Ended(now))
	fmt.Println(in.Next(now).Format(time.RFC3339))

	// Output:
	// <nil>
	// <nil>
	// 2019-01-02T20:45:00Z
	// 2019-01-02T22:00:00Z
	// true
	// true
	// true
	// false
	// 2019-01-02T21:15:00Z
}
