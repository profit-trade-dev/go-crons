package crons

import (
	"context"
	"testing"
	"time"
)

// MockLogger implements the Logger interface for testing
type MockLogger struct {
	debugLogs []string
	infoLogs  []string
	errorLogs []string
}

func (m *MockLogger) Debug(ctx context.Context, msg string, args ...any) {
	m.debugLogs = append(m.debugLogs, msg)
}

func (m *MockLogger) Info(ctx context.Context, msg string, args ...any) {
	m.infoLogs = append(m.infoLogs, msg)
}

func (m *MockLogger) Error(ctx context.Context, err error, msg string, args ...any) {
	m.errorLogs = append(m.errorLogs, msg)
}

func TestNewCron(t *testing.T) {
	cron := New()
	if cron == nil {
		t.Fatal("New() returned nil")
	}
	if cron.errThreshold != -1 {
		t.Errorf("Expected errThreshold to be -1, got %d", cron.errThreshold)
	}
	if cron.ch == nil {
		t.Fatal("Channel not initialized")
	}
}

func TestCronAt(t *testing.T) {
	cron := New()
	cron.At("09:30")
	if cron.at != "09:30" {
		t.Errorf("Expected at to be '09:30', got '%s'", cron.at)
	}
	if cron.every != 0 {
		t.Errorf("Expected every to be 0, got %v", cron.every)
	}
}

func TestCronEvery(t *testing.T) {
	cron := New()
	duration := 30 * time.Second
	cron.Every(duration)
	if cron.every != duration {
		t.Errorf("Expected every to be %v, got %v", duration, cron.every)
	}
	if cron.at != empty {
		t.Errorf("Expected at to be empty, got '%s'", cron.at)
	}
}

func TestCronStartsAt(t *testing.T) {
	cron := New()
	cron.StartsAt("10:00")
	if cron.startsAt != "10:00" {
		t.Errorf("Expected startsAt to be '10:00', got '%s'", cron.startsAt)
	}
}

func TestCronEndsAt(t *testing.T) {
	cron := New()
	cron.EndsAt("18:00")
	if cron.endsAt != "18:00" {
		t.Errorf("Expected endsAt to be '18:00', got '%s'", cron.endsAt)
	}
}

func TestCronWithErrorThreshold(t *testing.T) {
	cron := New()
	cron.WithErrorThreshold(3)
	if cron.errThreshold != 3 {
		t.Errorf("Expected errThreshold to be 3, got %d", cron.errThreshold)
	}
}

func TestCronWithLogger(t *testing.T) {
	cron := New()
	logger := &MockLogger{}
	cron.WithLogger(logger)
	if cron.logger != logger {
		t.Error("Logger not set correctly")
	}
}

func TestCronStartStop(t *testing.T) {
	cron := New()
	logger := &MockLogger{}
	cron.WithLogger(logger)

	ctx := context.Background()

	cron.Start(ctx, func(ctx context.Context) error {
		return nil
	})

	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)

	if !cron.isRunning {
		t.Error("Cron should be running")
	}

	cron.Stop()

	// Give it a moment to stop
	time.Sleep(100 * time.Millisecond)

	if cron.isRunning {
		t.Error("Cron should be stopped")
	}
}

func TestStopAll(t *testing.T) {
	// Create multiple crons
	cron1 := New()
	cron2 := New()
	cron3 := New()

	ctx := context.Background()

	// Start all crons
	cron1.Start(ctx, func(ctx context.Context) error { return nil })
	cron2.Start(ctx, func(ctx context.Context) error { return nil })
	cron3.Start(ctx, func(ctx context.Context) error { return nil })

	// Verify they're running
	if !cron1.isRunning || !cron2.isRunning || !cron3.isRunning {
		t.Error("All crons should be running")
	}

	// Stop all
	StopAll()

	// Give them time to stop
	time.Sleep(100 * time.Millisecond)

	// Verify they're stopped
	if cron1.isRunning || cron2.isRunning || cron3.isRunning {
		t.Error("All crons should be stopped")
	}
}
