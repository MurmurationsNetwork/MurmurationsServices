package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log logger
)

type commonLogger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
}

func init() {
	var err error
	if log.log, err = newConfig().Build(); err != nil {
		panic(err)
	}
}

func newConfig() zap.Config {
	return zap.Config{
		Level:    zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:      "time",
			LevelKey:     "level",
			MessageKey:   "msg",
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

// logger implements commonLogger interface.
type logger struct {
	log *zap.Logger
}

func GetLogger() commonLogger {
	return log
}

// Printf logs formatted messages.
func (l logger) Printf(format string, v ...interface{}) {
	if len(v) == 0 {
		Info(format)
	} else {
		Info(fmt.Sprintf(format, v...))
	}
}

func (l logger) Print(v ...interface{}) {
	Info(fmt.Sprintf("%v", v))
}

// Info logs informational messages.
func Info(msg string, tags ...zap.Field) {
	log.log.Info(msg, tags...)
	_ = log.log.Sync()
}

// Warn logs warning messages.
func Warn(msg string, tags ...zap.Field) {
	if log.log == nil {
		return
	}
	log.log.Warn(msg, tags...)
	_ = log.log.Sync()
}

// Error logs error messages.
func Error(msg string, err error, tags ...zap.Field) {
	if err != nil {
		tags = append(tags, zap.NamedError("error", err))
	}
	log.log.Error(msg, tags...)
	_ = log.log.Sync()
}

// Panic logs panic messages.
func Panic(msg string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("error", err))
	log.log.Panic(msg, tags...)
	_ = log.log.Sync()
}

// Fatal logs fatal messages and exits the program.
func Fatal(msg string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("error", err))
	log.log.Fatal(msg, tags...)
	_ = log.log.Sync()
}
