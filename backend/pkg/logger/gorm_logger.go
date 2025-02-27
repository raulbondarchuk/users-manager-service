package logger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormLogger struct {
	logger *Logger
}

func NewGormLogger(logger *Logger) *GormLogger {
	return &GormLogger{
		logger: logger,
	}
}

// LogMode implementation of gorm.Logger interface
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

// Info implementation of gorm.Logger interface
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.logger.GormInfo(fmt.Sprintf(msg, data...), nil)
}

// Warn implementation of gorm.Logger interface
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.logger.GormWarn(fmt.Sprintf(msg, data...), nil)
}

// Error implementation of gorm.Logger interface
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.logger.GormError(fmt.Sprintf(msg, data...), nil)
}

// Trace implementation of gorm.Logger interface
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := map[string]interface{}{
		"elapsed": elapsed,
		"rows":    rows,
		"sql":     sql,
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		l.logger.GormError("GORM Error", fields)
		return
	}

	// Log slow SQL queries if they take more than 1 second
	if elapsed > time.Second {
		l.logger.GormWarn("SLOW SQL", fields)
	} else {
		l.logger.GormDebug("SQL", fields)
	}
}
