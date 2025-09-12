# go-crons

A lightweight and feature-rich cron job scheduler library for Go applications with Indian timezone support and error thresholds.

## Features

- **Time-based Scheduling**: Schedule jobs at specific times (e.g., "09:30" for 9:30 AM IST)
- **Interval Scheduling**: Run jobs every N duration (e.g., every 30 seconds)
- **Time Windows**: Define start and end times for job execution
- **Indian Timezone**: Built-in support for Asia/Kolkata timezone
- **Error Thresholds**: Configure maximum consecutive errors before stopping
- **Custom Logging**: Pluggable logger interface
- **Context Support**: Jobs receive context for cancellation and timeouts
- **Zero Dependencies**: Uses only standard Go packages

## Installation

```bash
go get github.com/profit-trade-dev/go-crons
```

## Usage

### Basic Example

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/profit-trade-dev/go-crons"
)

// SimpleLogger implements the Logger interface
type SimpleLogger struct{}

func (l *SimpleLogger) Debug(ctx context.Context, msg string, args ...any) {
    log.Printf("[DEBUG] %s", fmt.Sprintf(msg, args...))
}

func (l *SimpleLogger) Info(ctx context.Context, msg string, args ...any) {
    log.Printf("[INFO] %s", fmt.Sprintf(msg, args...))
}

func (l *SimpleLogger) Error(ctx context.Context, err error, msg string, args ...any) {
    if err != nil {
        log.Printf("[ERROR] %s: %v", fmt.Sprintf(msg, args...), err)
    } else {
        log.Printf("[ERROR] %s", fmt.Sprintf(msg, args...))
    }
}

func main() {
    ctx := context.Background()

    // Job that runs at 9:30 AM IST daily
    cron := crons.New().
        At("09:30").
        WithLogger(&SimpleLogger{})

    cron.Start(ctx, func(ctx context.Context) error {
        log.Println("Running daily report at 9:30 AM IST...")
        // Your job logic here
        return nil
    })

    // Your application logic here...

    // Stop the cron when done
    cron.Stop()
}
```

### Advanced Example

```go
// Job that runs every 30 seconds between 10:00 AM and 6:00 PM IST
monitoringCron := crons.New().
    Every(30 * time.Second).
    StartsAt("10:00").
    EndsAt("18:00").
    WithErrorThreshold(3).
    WithLogger(&SimpleLogger{})

monitoringCron.Start(ctx, func(ctx context.Context) error {
    log.Println("Checking system health...")
    // Your monitoring logic here
    return nil
})
```

## API Reference

### New() *Cron
Creates a new cron instance.

### At(time string) *Cron
Sets the cron to run at a specific time every day. Time format: "15:04" (e.g., "09:30" for 9:30 AM).

### Every(duration time.Duration) *Cron
Sets the cron to run continuously every specified duration.

### StartsAt(time string) *Cron
Sets the start time for the cron execution window. Time format: "15:04".

### EndsAt(time string) *Cron
Sets the end time for the cron execution window. Time format: "15:04".

### WithErrorThreshold(threshold int) *Cron
Sets the maximum number of consecutive errors before stopping the cron. Use -1 for no threshold.

### WithLogger(logger Logger) *Cron
Sets a custom logger for the cron.

### Start(ctx context.Context, function RunFunction)
Starts the cron with the given function.

### Stop()
Stops the cron.

### StopAll()
Stops all running crons.

## Logger Interface

```go
type Logger interface {
    Debug(ctx context.Context, msg string, args ...any)
    Info(ctx context.Context, msg string, args ...any)
    Error(ctx context.Context, err error, msg string, args ...any)
}
```

## RunFunction Type

```go
type RunFunction func(ctx context.Context) error
```

## Time Format

All time strings should be in "15:04" format (24-hour format):
- "09:30" = 9:30 AM IST
- "14:15" = 2:15 PM IST
- "23:45" = 11:45 PM IST

## Examples

See `examples/original_cron_usage.go` for comprehensive usage examples.
