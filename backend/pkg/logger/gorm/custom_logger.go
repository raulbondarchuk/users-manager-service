package logger_gorm

import (
	"context"
	"log"
	"os"
	"time"

	"gorm.io/gorm/logger"
)

// CustomLogger is a custom GORM logger that ignores SQL logs
type CustomLogger struct {
	logger.Interface
}

// LogMode sets the log level
func (c *CustomLogger) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *c
	newlogger.Interface = c.Interface.LogMode(level)
	return &newlogger
}

// Info logs general information
func (c *CustomLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	// Do nothing to ignore info logs
}

// Warn logs warning messages
func (c *CustomLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	// Do nothing to ignore warning logs
}

// Error logs error messages
func (c *CustomLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	// Do nothing to ignore error logs
}

// NewCustomLogger creates a new custom logger
func NewCustomLogger() *CustomLogger {
	return &CustomLogger{
		Interface: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             200 * time.Millisecond, // Slow SQL threshold
				LogLevel:                  logger.Silent,          // Log level
				IgnoreRecordNotFoundError: true,                   // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,                   // Disable color
			},
		),
	}
}
