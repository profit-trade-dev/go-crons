package crons

import "time"

var il *time.Location

func init() {
	var err error
	il, err = time.LoadLocation("Asia/Kolkata")
	if err != nil {
		panic(err)
	}
}

func getCurrentIndianTime() time.Time {
	return time.Now().In(il)
}

func getIndianTimeFromTiming(t string) time.Time {
	ct := getCurrentIndianTime()
	tt := getTimeFromTiming(t)
	return time.Date(ct.Year(), ct.Month(), ct.Day(), tt.Hour(), tt.Minute(),
		tt.Second(), tt.Nanosecond(), ct.Location())
}

func getNextIndianTimeFromTiming(t string) time.Time {
	ct := getCurrentIndianTime()
	tt := getTimeFromTiming(t)
	r := time.Date(ct.Year(), ct.Month(), ct.Day(), tt.Hour(), tt.Minute(),
		tt.Second(), tt.Nanosecond(), ct.Location())
	if !r.After(ct) {
		return r.Add(time.Hour * 24)
	}
	return r
}

func getTimeFromTiming(t string) time.Time {
	ti, err := time.Parse("15:04", t)
	if err != nil {
		return getCurrentIndianTime()
	}
	return ti
}
