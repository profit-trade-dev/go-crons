package crons

import "context"

// Logger is the set of logger methods for cron.
type Logger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, err error, msg string, args ...any)
}

func debugLog(ctx context.Context, c *Cron, msg string, args ...any) {
	if c.logger != nil {
		c.logger.Debug(ctx, msg, args...)
	}
}

func infoLog(ctx context.Context, c *Cron, msg string, args ...any) {
	if c.logger != nil {
		c.logger.Info(ctx, msg, args...)
	}
}

func errorLog(ctx context.Context, c *Cron, err error, msg string, args ...any) {
	if c.logger != nil {
		c.logger.Error(ctx, err, msg, args...)
	}
}
