package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/profit-trade-dev/go-crons"
)

// ExampleLogger implements the Logger interface
type ExampleLogger struct{}

func (l *ExampleLogger) Debug(ctx context.Context, msg string, args ...any) {
	log.Printf("[DEBUG] %s", fmt.Sprintf(msg, args...))
}

func (l *ExampleLogger) Info(ctx context.Context, msg string, args ...any) {
	log.Printf("[INFO] %s", fmt.Sprintf(msg, args...))
}

func (l *ExampleLogger) Error(ctx context.Context, err error, msg string, args ...any) {
	if err != nil {
		log.Printf("[ERROR] %s: %v", fmt.Sprintf(msg, args...), err)
	} else {
		log.Printf("[ERROR] %s", fmt.Sprintf(msg, args...))
	}
}

func main() {
	fmt.Println("=== Cron Package Example with Location Support ===")

	ctx := context.Background()

	// Example 1: Job that runs at 9:30 AM IST daily (default location)
	fmt.Println("\n1. Creating daily report job at 9:30 AM IST (default location)...")
	dailyReportCron := crons.New().
		At("09:30").
		WithLogger(&ExampleLogger{})

	dailyReportCron.Start(ctx, func(ctx context.Context) error {
		fmt.Printf("Generating daily report at 9:30 AM IST (current time: %s)...\n", time.Now().Format("15:04:05 MST"))
		time.Sleep(1 * time.Second) // Simulate work
		fmt.Println("Daily report generated successfully!")
		return nil
	})

	// Example 2: Job that runs at 2:00 PM UTC using UTC convenience method
	fmt.Println("\n2. Creating backup job at 2:00 PM UTC...")
	backupCron := crons.New().
		AtUTC("14:00").
		WithLogger(&ExampleLogger{})

	backupCron.Start(ctx, func(ctx context.Context) error {
		fmt.Printf("Running backup at 2:00 PM UTC (current time: %s)...\n", time.Now().UTC().Format("15:04:05 MST"))
		time.Sleep(1 * time.Second) // Simulate work
		fmt.Println("Backup completed successfully!")
		return nil
	})

	// Example 3: Job that runs every 30 seconds between 10:00 AM and 6:00 PM IST
	fmt.Println("\n3. Creating monitoring job every 30 seconds (10:00 AM - 6:00 PM IST)...")
	monitoringCron := crons.New().
		Every(30 * time.Second).
		StartsAt("10:00").
		EndsAt("18:00").
		WithErrorThreshold(3).
		WithLogger(&ExampleLogger{})

	monitoringCron.Start(ctx, func(ctx context.Context) error {
		fmt.Printf("Checking system health (current time: %s)...\n", time.Now().Format("15:04:05 MST"))
		time.Sleep(500 * time.Millisecond) // Simulate work
		fmt.Println("System health check completed!")
		return nil
	})

	// Example 4: Job that runs every 2 minutes in UTC timezone
	fmt.Println("\n4. Creating data sync job every 2 minutes in UTC...")
	dataSyncCron := crons.New().
		WithLocation(time.UTC).
		Every(2 * time.Minute).
		WithErrorThreshold(2).
		WithLogger(&ExampleLogger{})

	dataSyncCron.Start(ctx, func(ctx context.Context) error {
		fmt.Printf("Syncing data in UTC (current time: %s)...\n", time.Now().UTC().Format("15:04:05 MST"))
		time.Sleep(1 * time.Second) // Simulate work
		fmt.Println("Data sync completed!")
		return nil
	})

	// Example 5: Job that runs every 5 minutes with custom IST location
	fmt.Println("\n5. Creating cleanup job every 5 minutes with explicit IST location...")
	istLocation, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		fmt.Printf("Failed to load IST location: %v\n", err)
		return
	}

	cleanupCron := crons.New().
		WithLocation(istLocation).
		Every(5 * time.Minute).
		WithErrorThreshold(2).
		WithLogger(&ExampleLogger{})

	cleanupCron.Start(ctx, func(ctx context.Context) error {
		fmt.Printf("Running cleanup in IST (current time: %s)...\n", time.Now().In(istLocation).Format("15:04:05 MST"))
		time.Sleep(1 * time.Second) // Simulate work
		fmt.Println("Cleanup completed!")
		return nil
	})

	fmt.Println("\nAll cron jobs started!")
	fmt.Println("Jobs running:")
	fmt.Println("- Daily report at 9:30 AM IST (default location)")
	fmt.Println("- Backup at 2:00 PM UTC")
	fmt.Println("- System monitoring every 30 seconds (10:00 AM - 6:00 PM IST)")
	fmt.Println("- Data sync every 2 minutes in UTC")
	fmt.Println("- Cleanup every 5 minutes in IST")

	// Run for a short time to demonstrate
	fmt.Println("\nRunning for 1 minute to demonstrate...")
	time.Sleep(1 * time.Minute)

	// Stop all crons
	fmt.Println("\nStopping all cron jobs...")
	crons.StopAll()
	fmt.Println("All cron jobs stopped!")
}
