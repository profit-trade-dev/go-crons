package crons

import "time"

// defaultLocation holds the package-wide default location used by new crons.
var defaultLocation *time.Location

func init() {
	var err error
	// Backward-compat default remains Asia/Kolkata (IST)
	defaultLocation, err = time.LoadLocation("Asia/Kolkata")
	if err != nil {
		panic(err)
	}
}

// SetDefaultLocation allows callers to change the default location for new crons.
func SetDefaultLocation(loc *time.Location) {
	if loc != nil {
		defaultLocation = loc
	}
}

// nowIn returns the current time in the provided location (or default location if nil).
func nowIn(loc *time.Location) time.Time {
	if loc == nil {
		loc = defaultLocation
	}
	return time.Now().In(loc)
}

// timeFromHHMMInLoc parses "HH:MM" and returns a time on the same day as base, in loc.
func timeFromHHMMInLoc(hhmm string, base time.Time, loc *time.Location) time.Time {
	if loc == nil {
		loc = defaultLocation
	}
	t, err := time.Parse("15:04", hhmm)
	if err != nil {
		return base.In(loc)
	}
	b := base.In(loc)
	return time.Date(b.Year(), b.Month(), b.Day(), t.Hour(), t.Minute(), 0, 0, loc)
}

// nextTimeFromHHMMInLoc returns the next occurrence (today or tomorrow) of HH:MM in loc.
func nextTimeFromHHMMInLoc(hhmm string, loc *time.Location) time.Time {
	ct := nowIn(loc)
	tt := timeFromHHMMInLoc(hhmm, ct, loc)
	if !tt.After(ct) {
		return tt.Add(24 * time.Hour)
	}
	return tt
}

// Backward-compat helpers (preserved but now implemented via location-aware helpers)
func getCurrentIndianTime() time.Time { return nowIn(defaultLocation) }
func getIndianTimeFromTiming(t string) time.Time {
	return timeFromHHMMInLoc(t, nowIn(defaultLocation), defaultLocation)
}
func getNextIndianTimeFromTiming(t string) time.Time {
	return nextTimeFromHHMMInLoc(t, defaultLocation)
}
