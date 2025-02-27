package logger

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format("2006-01-02 15:04:05")

	var levelColor, typeColor *color.Color
	levelColor = color.New(color.FgWhite).Add(color.Bold)

	switch {
	case entry.Level == logrus.DebugLevel:
		levelColor = color.New(color.FgMagenta).Add(color.Bold)
	case entry.Level == logrus.InfoLevel:
		levelColor = color.New(color.FgBlue).Add(color.Bold)
	case entry.Level == logrus.WarnLevel:
		levelColor = color.New(color.FgYellow).Add(color.Bold)
	case entry.Level == logrus.ErrorLevel, entry.Level == logrus.FatalLevel, entry.Level == logrus.PanicLevel:
		levelColor = color.New(color.FgRed).Add(color.Bold)
	}

	level := strings.ToUpper(entry.Level.String())
	prefix := ""
	isGorm := false

	if logType, exists := entry.Data["type"]; exists {
		switch logType {
		case "gorm":
			isGorm = true
			typeColor = color.New(color.FgBlue).Add(color.Bold)
			prefix = typeColor.Sprintf("[GORM]         ")
			levelColor = typeColor
		case "database":
			typeColor = color.New(color.FgGreen).Add(color.Bold)
			prefix = typeColor.Sprintf("[DATABASE]     ")
		case "service":
			typeColor = color.New(color.FgYellow).Add(color.Bold)
			prefix = typeColor.Sprintf("[SERVICE]      ")
		case "handler":
			typeColor = color.New(color.FgCyan).Add(color.Bold)
			prefix = typeColor.Sprintf("[HANDLER]      ")
		case "route":
			typeColor = color.New(color.FgWhite).Add(color.Bold)
			prefix = typeColor.Sprintf("[ROUTE] --> ")
		case "websocket":
			typeColor = color.New(color.FgBlue).Add(color.Bold)
			prefix = typeColor.Sprintf("[WEBSOCKET]    ")
		case "verificaciones":
			typeColor = color.New(color.FgHiGreen).Add(color.Bold)
			prefix = typeColor.Sprintf("[VERIFICACIONES]    ")
		case "other":
			typeColor = color.New(color.FgMagenta).Add(color.Bold)
			prefix = typeColor.Sprintf("[OTHER]        ")
		}
		delete(entry.Data, "type")
	}

	msg := fmt.Sprintf("%s %s%s %s",
		color.New(color.FgWhite).Add(color.Bold).Sprint(timestamp),
		prefix,
		levelColor.Sprintf("[%s]", level),
		func() string {
			if isGorm {
				return typeColor.Sprint(entry.Message)
			}
			return color.New(color.FgWhite).Add(color.Bold).Sprint(entry.Message)
		}())

	for k, v := range entry.Data {
		if isGorm {
			msg += fmt.Sprintf(" [%s:%v]",
				typeColor.Sprint(k),
				typeColor.Sprint(v))
		} else {
			msg += fmt.Sprintf(" [%s:%v]",
				color.New(color.FgCyan).Add(color.Bold).Sprint(k),
				color.New(color.FgWhite).Add(color.Bold).Sprint(v))
		}
	}

	return []byte(msg + "\n"), nil
}
