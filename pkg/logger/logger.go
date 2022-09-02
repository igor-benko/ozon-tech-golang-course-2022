package logger

import (
	"log"

	"go.uber.org/zap"
)

var defaultLogger *zap.SugaredLogger

func init() {
	nonSugaredLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	defaultLogger = nonSugaredLogger.Sugar()
}

func WithAppName(name string) {
	defaultLogger.With("app_name", name)
}

// Key-Value support
func InfoKV(message string, args ...interface{}) {
	defaultLogger.Infow(message, args...)
}

func WarnKV(message string, args ...interface{}) {
	defaultLogger.Warnw(message, args...)
}

func DebugKV(message string, args ...interface{}) {
	defaultLogger.Debugw(message, args...)
}

func ErrorKV(message string, args ...interface{}) {
	defaultLogger.Errorw(message, args...)
}

func FatalKV(message string, args ...interface{}) {
	defaultLogger.Fatalw(message, args...)
}

// Formatted
func Infof(message string, args ...interface{}) {
	defaultLogger.Infof(message, args...)
}

func Warnf(message string, args ...interface{}) {
	defaultLogger.Warnf(message, args...)
}

func Debugf(message string, args ...interface{}) {
	defaultLogger.Debugf(message, args...)
}

func Errorf(message string, args ...interface{}) {
	defaultLogger.Errorf(message, args...)
}

func Fatalf(message string, args ...interface{}) {
	defaultLogger.Fatalf(message, args...)
}
