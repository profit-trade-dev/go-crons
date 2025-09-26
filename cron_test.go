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

// Location support tests
func TestWithLocation(t *testing.T) {
	cron := New()

	// Test setting IST location
	istLocation, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		t.Fatalf("Failed to load IST location: %v", err)
	}

	cron.WithLocation(istLocation)
	if cron.location != istLocation {
		t.Error("IST location not set correctly")
	}

	// Test setting UTC location
	cron.WithLocation(time.UTC)
	if cron.location != time.UTC {
		t.Error("UTC location not set correctly")
	}

	// Test setting nil location (should not change)
	cron.WithLocation(nil)
	if cron.location != time.UTC {
		t.Error("Location should not change when setting nil")
	}
}

func TestUTCConvenienceMethods(t *testing.T) {
	cron := New()

	// Test AtUTC
	cron.AtUTC("14:30")
	if cron.location != time.UTC {
		t.Error("AtUTC should set location to UTC")
	}
	if cron.at != "14:30" {
		t.Error("AtUTC should set the time correctly")
	}

	// Test StartsAtUTC
	cron = New()
	cron.StartsAtUTC("09:00")
	if cron.location != time.UTC {
		t.Error("StartsAtUTC should set location to UTC")
	}
	if cron.startsAt != "09:00" {
		t.Error("StartsAtUTC should set the start time correctly")
	}

	// Test EndsAtUTC
	cron = New()
	cron.EndsAtUTC("17:00")
	if cron.location != time.UTC {
		t.Error("EndsAtUTC should set location to UTC")
	}
	if cron.endsAt != "17:00" {
		t.Error("EndsAtUTC should set the end time correctly")
	}
}

func TestDefaultLocation(t *testing.T) {
	// Test that new crons use IST as default
	cron := New()
	istLocation, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		t.Fatalf("Failed to load IST location: %v", err)
	}

	if cron.location.String() != istLocation.String() {
		t.Errorf("New cron should use IST as default location, got %s, expected %s", cron.location.String(), istLocation.String())
	}
}

func TestSetDefaultLocation(t *testing.T) {
	// Save original default location
	originalLocation := defaultLocation

	// Test setting new default location
	utcLocation := time.UTC
	SetDefaultLocation(utcLocation)

	// Create new cron and verify it uses UTC
	cron := New()
	if cron.location != utcLocation {
		t.Error("New cron should use UTC as default location after SetDefaultLocation")
	}

	// Test setting nil (should not change)
	SetDefaultLocation(nil)
	cron2 := New()
	if cron2.location != utcLocation {
		t.Error("Location should not change when setting nil default")
	}

	// Restore original location
	SetDefaultLocation(originalLocation)
}

func TestTimezoneAwareTimeCalculations(t *testing.T) {
	istLocation, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		t.Fatalf("Failed to load IST location: %v", err)
	}

	// Test timeFromHHMMInLoc with IST
	baseTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	istTime := timeFromHHMMInLoc("14:30", baseTime, istLocation)

	expectedIST := time.Date(2024, 1, 15, 14, 30, 0, 0, istLocation)
	if !istTime.Equal(expectedIST) {
		t.Errorf("Expected IST time %v, got %v", expectedIST, istTime)
	}

	// Test timeFromHHMMInLoc with UTC
	utcTime := timeFromHHMMInLoc("14:30", baseTime, time.UTC)
	expectedUTC := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
	if !utcTime.Equal(expectedUTC) {
		t.Errorf("Expected UTC time %v, got %v", expectedUTC, utcTime)
	}
}

func TestNextTimeFromHHMMInLoc(t *testing.T) {
	istLocation, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		t.Fatalf("Failed to load IST location: %v", err)
	}

	// Test with IST location
	nextIST := nextTimeFromHHMMInLoc("14:30", istLocation)
	nowIST := nowIn(istLocation)

	if !nextIST.After(nowIST) {
		t.Error("Next IST time should be after current IST time")
	}

	// Test with UTC location
	nextUTC := nextTimeFromHHMMInLoc("14:30", time.UTC)
	nowUTC := nowIn(time.UTC)

	if !nextUTC.After(nowUTC) {
		t.Error("Next UTC time should be after current UTC time")
	}
}

func TestCronWithLocationIntegration(t *testing.T) {
	istLocation, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		t.Fatalf("Failed to load IST location: %v", err)
	}

	// Test cron with IST location
	cronIST := New().WithLocation(istLocation).At("14:30")
	if cronIST.location != istLocation {
		t.Error("Cron should maintain IST location")
	}
	if cronIST.at != "14:30" {
		t.Error("Cron should maintain the time setting")
	}

	// Test cron with UTC location
	cronUTC := New().WithLocation(time.UTC).At("14:30")
	if cronUTC.location != time.UTC {
		t.Error("Cron should maintain UTC location")
	}
	if cronUTC.at != "14:30" {
		t.Error("Cron should maintain the time setting")
	}
}

func TestCronLocationWithEvery(t *testing.T) {
	istLocation, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		t.Fatalf("Failed to load IST location: %v", err)
	}

	// Test cron with IST location and Every
	cronIST := New().WithLocation(istLocation).Every(1 * time.Minute)
	if cronIST.location != istLocation {
		t.Error("Cron should maintain IST location with Every")
	}
	if cronIST.every != 1*time.Minute {
		t.Error("Cron should maintain the duration setting")
	}

	// Test cron with UTC location and Every
	cronUTC := New().WithLocation(time.UTC).Every(1 * time.Minute)
	if cronUTC.location != time.UTC {
		t.Error("Cron should maintain UTC location with Every")
	}
	if cronUTC.every != 1*time.Minute {
		t.Error("Cron should maintain the duration setting")
	}
}

func TestCronLocationWithStartEndTimes(t *testing.T) {
	istLocation, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		t.Fatalf("Failed to load IST location: %v", err)
	}

	// Test cron with IST location and start/end times
	cronIST := New().
		WithLocation(istLocation).
		Every(1 * time.Minute).
		StartsAt("09:00").
		EndsAt("17:00")

	if cronIST.location != istLocation {
		t.Error("Cron should maintain IST location with start/end times")
	}
	if cronIST.startsAt != "09:00" {
		t.Error("Cron should maintain start time")
	}
	if cronIST.endsAt != "17:00" {
		t.Error("Cron should maintain end time")
	}

	// Test cron with UTC location and start/end times
	cronUTC := New().
		WithLocation(time.UTC).
		Every(1 * time.Minute).
		StartsAt("09:00").
		EndsAt("17:00")

	if cronUTC.location != time.UTC {
		t.Error("Cron should maintain UTC location with start/end times")
	}
	if cronUTC.startsAt != "09:00" {
		t.Error("Cron should maintain start time")
	}
	if cronUTC.endsAt != "17:00" {
		t.Error("Cron should maintain end time")
	}
}
