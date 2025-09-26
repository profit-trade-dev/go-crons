package crons

import (
	"context"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// RunFunction is used to run the given function with cron.
type RunFunction func(ctx context.Context) error

// Cron is used to denote an instance of the cron.
type Cron struct {
	at string

	startsAt string
	endsAt   string
	every    time.Duration

	ch chan struct{}

	errThreshold int

	logger Logger

	isRunning bool

	mu sync.Mutex

	// location controls timezone for scheduling and comparisons
	location *time.Location
}

var crons []*Cron
var cronsMu sync.Mutex

// runIDKey is the unique context key used to store the cron run ID.
const runIDKey = "cronTraceId"

// New is used to create a new cron instance.
func New() *Cron {
	cronsMu.Lock()
	defer cronsMu.Unlock()
	c := &Cron{
		ch:           make(chan struct{}, 2),
		errThreshold: -1,
		location:     defaultLocation,
	}
	crons = append(crons, c)
	return c
}

// StopAll is used to stop all the running crons.
func StopAll() {
	cronsMu.Lock()
	defer cronsMu.Unlock()
	for _, c := range crons {
		c.Stop()
	}
}

// At is used to specify cron runs at an exact time every day.
// The time t should be of the following format, 15:04
func (c *Cron) At(t string) *Cron {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.isRunning {
		return c
	}
	c.at = strings.TrimSpace(t)
	c.every = 0
	return c
}

// Every is used to run the cron continuously every d duration.
func (c *Cron) Every(d time.Duration) *Cron {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.isRunning {
		return c
	}
	c.at = empty
	c.every = d
	return c
}

// StartsAt is used to run the cron from the given starting time.
// If current time has passed start for the day, then it will
// start to run at the next applicable time. Otherwise, it will run
// starting at the given start time.
func (c *Cron) StartsAt(s string) *Cron {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.isRunning {
		return c
	}
	c.startsAt = s
	return c
}

// EndsAt is used to run the cron till the given ending time.
func (c *Cron) EndsAt(s string) *Cron {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.isRunning {
		return c
	}
	c.endsAt = s
	return c
}

// WithErrorThreshold is used to add the max threshold to sustain till.
func (c *Cron) WithErrorThreshold(threshold int) *Cron {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.isRunning {
		return c
	}
	c.errThreshold = threshold
	return c
}

// WithLogger is used to put logger for cron.
func (c *Cron) WithLogger(l Logger) *Cron {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.isRunning {
		return c
	}
	c.logger = l
	return c
}

// WithLocation sets the timezone location for this cron.
func (c *Cron) WithLocation(loc *time.Location) *Cron {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.isRunning {
		return c
	}
	if loc != nil {
		c.location = loc
	}
	return c
}

// Start is used to start the cron run.
func (c *Cron) Start(ctx context.Context, f RunFunction) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.isRunning {
		return
	}
	c.isRunning = true
	if c.at != empty {
		go c.runAt(ctx, f)
		return
	}
	if c.every > 0 {
		go c.runEvery(ctx, f)
	}
}

// Stop is used to stop the cron run.
func (c *Cron) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.isRunning {
		c.ch <- struct{}{}
		c.isRunning = false
	}
}

func (c *Cron) runAt(ctx context.Context, f RunFunction) {
	defer c.recovery(ctx)
	errCnt := 0
	for {
		select {
		case <-c.ch:
			infoLog(ctx, c, "exiting cron run")
			return
		default:
			time.Sleep(nextTimeFromHHMMInLoc(c.at, c.location).
				Sub(nowIn(c.location)))
			runCtx := withRunID(ctx)
			debugLog(runCtx, c, "executing cron run")
			err := f(runCtx)
			if err != nil {
				errCnt++
				errorLog(runCtx, c, err, "error executing cron run")
				if c.errThreshold >= 0 && errCnt > c.errThreshold {
					os.Exit(1)
					return
				}
				continue
			}
			errCnt = 0
		}
	}
}

func (c *Cron) runEvery(ctx context.Context, f RunFunction) {
	defer c.recovery(ctx)
	errCnt := 0
	var t *time.Ticker
	defer stopTicker(t)
	for {
		time.Sleep(getSleepDurationForRun(c.startsAt, c.endsAt, c.every, c.location))
		t = time.NewTicker(c.every)
		for {
			runCtx := withRunID(ctx)
			debugLog(runCtx, c, "executing cron run")
			now := nowIn(c.location)
			if (c.endsAt != empty && now.After(timeFromHHMMInLoc(c.endsAt, now, c.location))) ||
				(c.startsAt != empty && now.Before(timeFromHHMMInLoc(c.startsAt, now, c.location))) {
				break
			}
			err := f(runCtx)
			if err != nil {
				errCnt++
				errorLog(runCtx, c, err, "error executing cron run")
				if c.errThreshold >= 0 && errCnt > c.errThreshold {
					os.Exit(1)
					return
				}
			} else {
				errCnt = 0
			}
			select {
			case <-t.C:
				continue
			case <-c.ch:
				infoLog(ctx, c, "exiting cron run")
				return
			}
		}
	}
}

func (c *Cron) recovery(ctx context.Context) {
	if r := recover(); r != nil {
		errorLog(ctx, c, nil, "cron panic occurred: %v", r)
	}
}

func stopTicker(t *time.Ticker) {
	if t != nil {
		t.Stop()
	}
}

// UTC convenience methods
func (c *Cron) AtUTC(t string) *Cron       { return c.WithLocation(time.UTC).At(t) }
func (c *Cron) StartsAtUTC(s string) *Cron { return c.WithLocation(time.UTC).StartsAt(s) }
func (c *Cron) EndsAtUTC(s string) *Cron   { return c.WithLocation(time.UTC).EndsAt(s) }

func getSleepDurationForRun(startsAt, endsAt string, every time.Duration, loc *time.Location) time.Duration {
	ct := nowIn(loc)
	st := timeFromHHMMInLoc(startsAt, ct, loc)
	et := timeFromHHMMInLoc(endsAt, ct, loc)
	if ct.Equal(et) {
		et = ct.Add(24 * time.Hour * 365)
	}
	if ct.Before(st) || ct.Equal(st) {
		return st.Sub(ct)
	}
	if ct.Before(et) || ct.Equal(et) {
		d := every - (ct.Sub(st))%every
		if ct.Add(d).Before(et) || ct.Add(d).Equal(et) {
			return d
		}
	}
	return st.Add(24 * time.Hour).Sub(ct)
}

func newRunID() string { return uuid.NewString() }

func withRunID(ctx context.Context) context.Context {
	return context.WithValue(ctx, runIDKey, newRunID())
}
