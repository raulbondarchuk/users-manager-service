package logger

import (
	"os"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	LOCAL = "local"
	DEV   = "dev"
	TEST  = "test"
	PROD  = "prod"
)

const (
	INFO_LOGGER_INIT  = "Logger initialized"
	WARN_INVALID_MODE = "Invalid logger mode, using default log level."
)

const (
	DEBUG          = "debug"
	INFO           = "info"
	DATABASE       = "database"
	SERVICE        = "service"
	HANDLER        = "handler"
	ROUTE          = "route"
	OTHER          = "other"
	GORM           = "gorm"
	WEBSOCKET      = "websocket"
	VERIFICACIONES = "verificaciones"
)

var (
	instance *Logger
	once     sync.Once
)

// Logger - structure for our logger
type Logger struct {
	*logrus.Logger
}

// Initialize creates and sets up the global logger instance
func MustLoad() {
	once.Do(func() {
		invalidMode := false
		mode := strings.ToLower(viper.GetString("logger.mode"))
		var logLevel string

		switch mode {
		case DEV, LOCAL, TEST:
			logLevel = DEBUG
		case PROD:
			logLevel = INFO
		default:

			invalidMode = true
			logLevel = DEBUG
		}
		instance = New(logLevel)

		if invalidMode {
			instance.Warn(WARN_INVALID_MODE, map[string]interface{}{"logger_level": logLevel})
		}
		instance.Info(INFO_LOGGER_INIT, map[string]interface{}{"logger_level": logLevel})

	})
}

// GetLogger returns the global logger instance
func GetLogger() *Logger {
	if instance == nil {
		MustLoad()
	}
	return instance
}

// New creates a new instance of the logger
func New(level string) *Logger {
	log := logrus.New()

	log.SetFormatter(&CustomFormatter{})
	log.SetOutput(os.Stdout)

	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		lvl = logrus.InfoLevel
	}
	log.SetLevel(lvl)

	return &Logger{log}
}

// logWithType - base method for logging with a specific type
func (l *Logger) logWithType(logType string, level logrus.Level, msg string, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["type"] = logType

	entry := l.WithFields(fields)
	switch level {
	case logrus.InfoLevel:
		entry.Info(msg)
	case logrus.ErrorLevel:
		entry.Error(msg)
	case logrus.WarnLevel:
		entry.Warn(msg)
	case logrus.DebugLevel:
		entry.Debug(msg)
	}
}

// Base logging methods
func (l *Logger) Info(msg string, fields map[string]interface{})  { l.WithFields(fields).Info(msg) }
func (l *Logger) Error(msg string, fields map[string]interface{}) { l.WithFields(fields).Error(msg) }
func (l *Logger) Debug(msg string, fields map[string]interface{}) { l.WithFields(fields).Debug(msg) }
func (l *Logger) Warn(msg string, fields map[string]interface{})  { l.WithFields(fields).Warn(msg) }

// Database logging methods
func (l *Logger) DatabaseInfo(msg string, fields map[string]interface{}) {
	l.logWithType(DATABASE, logrus.InfoLevel, msg, fields)
}

func (l *Logger) DatabaseError(msg string, fields map[string]interface{}) {
	l.logWithType(DATABASE, logrus.ErrorLevel, msg, fields)
}

func (l *Logger) DatabaseWarn(msg string, fields map[string]interface{}) {
	l.logWithType(DATABASE, logrus.WarnLevel, msg, fields)
}

func (l *Logger) DatabaseDebug(msg string, fields map[string]interface{}) {
	l.logWithType(DATABASE, logrus.DebugLevel, msg, fields)
}

// Service logging methods
func (l *Logger) ServiceInfo(msg string, fields map[string]interface{}) {
	l.logWithType(SERVICE, logrus.InfoLevel, msg, fields)
}

func (l *Logger) ServiceError(msg string, fields map[string]interface{}) {
	l.logWithType(SERVICE, logrus.ErrorLevel, msg, fields)
}

func (l *Logger) ServiceWarn(msg string, fields map[string]interface{}) {
	l.logWithType(SERVICE, logrus.WarnLevel, msg, fields)
}

func (l *Logger) ServiceDebug(msg string, fields map[string]interface{}) {
	l.logWithType(SERVICE, logrus.DebugLevel, msg, fields)
}

// Handler logging methods
func (l *Logger) HandlerInfo(msg string, fields map[string]interface{}) {
	l.logWithType(HANDLER, logrus.InfoLevel, msg, fields)
}

func (l *Logger) HandlerError(msg string, fields map[string]interface{}) {
	l.logWithType(HANDLER, logrus.ErrorLevel, msg, fields)
}

func (l *Logger) HandlerWarn(msg string, fields map[string]interface{}) {
	l.logWithType(HANDLER, logrus.WarnLevel, msg, fields)
}

func (l *Logger) HandlerDebug(msg string, fields map[string]interface{}) {
	l.logWithType(HANDLER, logrus.DebugLevel, msg, fields)
}

// Route logging methods
func (l *Logger) RouteInfo(msg string, fields map[string]interface{}) {
	l.logWithType(ROUTE, logrus.InfoLevel, msg, fields)
}

func (l *Logger) RouteError(msg string, fields map[string]interface{}) {
	l.logWithType(ROUTE, logrus.ErrorLevel, msg, fields)
}

func (l *Logger) RouteWarn(msg string, fields map[string]interface{}) {
	l.logWithType(ROUTE, logrus.WarnLevel, msg, fields)
}

func (l *Logger) RouteDebug(msg string, fields map[string]interface{}) {
	l.logWithType(ROUTE, logrus.DebugLevel, msg, fields)
}

// Other logging methods
func (l *Logger) OtherInfo(msg string, fields map[string]interface{}) {
	l.logWithType(OTHER, logrus.InfoLevel, msg, fields)
}

func (l *Logger) OtherError(msg string, fields map[string]interface{}) {
	l.logWithType(OTHER, logrus.ErrorLevel, msg, fields)
}

func (l *Logger) OtherWarn(msg string, fields map[string]interface{}) {
	l.logWithType(OTHER, logrus.WarnLevel, msg, fields)
}

func (l *Logger) OtherDebug(msg string, fields map[string]interface{}) {
	l.logWithType(OTHER, logrus.DebugLevel, msg, fields)
}

// GORM logging methods
func (l *Logger) GormInfo(msg string, fields map[string]interface{}) {
	l.logWithType(GORM, logrus.InfoLevel, msg, fields)
}

func (l *Logger) GormError(msg string, fields map[string]interface{}) {
	l.logWithType(GORM, logrus.ErrorLevel, msg, fields)
}

func (l *Logger) GormWarn(msg string, fields map[string]interface{}) {
	l.logWithType(GORM, logrus.WarnLevel, msg, fields)
}

func (l *Logger) GormDebug(msg string, fields map[string]interface{}) {
	l.logWithType(GORM, logrus.DebugLevel, msg, fields)
}

// WebSocket logging methods
func (l *Logger) WebSocketInfo(msg string, fields map[string]interface{}) {
	l.logWithType(WEBSOCKET, logrus.InfoLevel, msg, fields)
}

func (l *Logger) WebSocketError(msg string, fields map[string]interface{}) {
	l.logWithType(WEBSOCKET, logrus.ErrorLevel, msg, fields)
}

func (l *Logger) WebSocketWarn(msg string, fields map[string]interface{}) {
	l.logWithType(WEBSOCKET, logrus.WarnLevel, msg, fields)
}

func (l *Logger) WebSocketDebug(msg string, fields map[string]interface{}) {
	l.logWithType(WEBSOCKET, logrus.DebugLevel, msg, fields)
}

// Verificaciones logging methods
func (l *Logger) VerificacionesInfo(msg string, fields map[string]interface{}) {
	l.logWithType(VERIFICACIONES, logrus.InfoLevel, msg, fields)
}

func (l *Logger) VerificacionesError(msg string, fields map[string]interface{}) {
	l.logWithType(VERIFICACIONES, logrus.ErrorLevel, msg, fields)
}

func (l *Logger) VerificacionesWarn(msg string, fields map[string]interface{}) {
	l.logWithType(VERIFICACIONES, logrus.WarnLevel, msg, fields)
}

func (l *Logger) VerificacionesDebug(msg string, fields map[string]interface{}) {
	l.logWithType(VERIFICACIONES, logrus.DebugLevel, msg, fields)
}
