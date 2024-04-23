package log

import (
	"github.com/sirupsen/logrus"
	microLogger "go-micro.dev/v4/logger"
)

type LogrusAdapter struct{}

func (l *LogrusAdapter) Init(opts ...microLogger.Option) error {
	return nil
}

func (l *LogrusAdapter) Options() microLogger.Options {
	return microLogger.Options{}
}

func (l *LogrusAdapter) Fields(fields map[string]interface{}) microLogger.Logger {
	return l
}

func (l *LogrusAdapter) Log(level microLogger.Level, args ...interface{}) {
	switch level {
	case microLogger.DebugLevel:
		logrus.Debug(args...)
	case microLogger.InfoLevel:
		logrus.Info(args...)
	case microLogger.WarnLevel:
		logrus.Warn(args...)
	case microLogger.ErrorLevel:
		logrus.Error(args...)
	case microLogger.FatalLevel:
		logrus.Fatal(args...)
	default:
		panic("unhandled default case")
	}
}

func (l *LogrusAdapter) Logf(level microLogger.Level, format string, args ...interface{}) {
	switch level {
	case microLogger.DebugLevel:
		logrus.Debugf(format, args...)
	case microLogger.InfoLevel:
		logrus.Infof(format, args...)
	case microLogger.WarnLevel:
		logrus.Warnf(format, args...)
	case microLogger.ErrorLevel:
		logrus.Errorf(format, args...)
	case microLogger.FatalLevel:
		logrus.Fatalf(format, args...)
	default:
		panic("unhandled default case")
	}
}

func (l *LogrusAdapter) String() string {
	return "logrus"
}
