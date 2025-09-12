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
	fmt.Println("=== Cron Package Example ===")

	ctx := context.Background()

	// Example 1: Job that runs at 9:30 AM IST daily
	fmt.Println("\n1. Creating daily report job at 9:30 AM IST...")
	dailyReportCron := crons.New().
		At("09:30").
		WithLogger(&ExampleLogger{})

	dailyReportCron.Start(ctx, func(ctx context.Context) error {
		fmt.Println("Generating daily report at 9:30 AM IST...")
		time.Sleep(1 * time.Second) // Simulate work
		fmt.Println("Daily report generated successfully!")
		return nil
	})

	// Example 2: Job that runs every 30 seconds between 10:00 AM and 6:00 PM IST
	fmt.Println("\n2. Creating monitoring job every 30 seconds (10:00 AM - 6:00 PM IST)...")
	monitoringCron := crons.New().
		Every(30 * time.Second).
		StartsAt("10:00").
		EndsAt("18:00").
		WithErrorThreshold(3).
		WithLogger(&ExampleLogger{})

	monitoringCron.Start(ctx, func(ctx context.Context) error {
		fmt.Println("Checking system health...")
		time.Sleep(500 * time.Millisecond) // Simulate work
		fmt.Println("System health check completed!")
		return nil
	})

	// Example 3: Job that runs every 5 minutes with error threshold
	fmt.Println("\n3. Creating data sync job every 5 minutes...")
	dataSyncCron := crons.New().
		Every(5 * time.Minute).
		WithErrorThreshold(2).
		WithLogger(&ExampleLogger{})

	dataSyncCron.Start(ctx, func(ctx context.Context) error {
		fmt.Println("Syncing data...")
		time.Sleep(1 * time.Second) // Simulate work
		fmt.Println("Data sync completed!")
		return nil
	})

	fmt.Println("\nAll cron jobs started!")
	fmt.Println("Jobs running:")
	fmt.Println("- Daily report at 9:30 AM IST")
	fmt.Println("- System monitoring every 30 seconds (10:00 AM - 6:00 PM IST)")
	fmt.Println("- Data sync every 5 minutes")

	// Run for a short time to demonstrate
	fmt.Println("\nRunning for 2 minutes to demonstrate...")
	time.Sleep(2 * time.Minute)

	// Stop all crons
	fmt.Println("\nStopping all cron jobs...")
	crons.StopAll()
	fmt.Println("All cron jobs stopped!")
}
