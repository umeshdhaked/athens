package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	once   sync.Once
	logger *Logger
)

type Logger struct {
	l *zap.Logger
}

func GetLogger() ILogger {
	return logger
}

func Build() {
	once.Do(func() {
		// Define encoder configuration
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		// Configure output paths
		stdoutSyncer := zapcore.AddSync(os.Stdout)

		// Create a console encoder
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

		// Create a core
		core := zapcore.NewCore(consoleEncoder, stdoutSyncer, zapcore.DebugLevel)

		// Create a logger
		l := zap.New(core)

		// Flush the logger
		defer l.Sync()

		logger = &Logger{
			l: l,
		}
	})
}

// WithField adds a single field to the logger.
func (l *Logger) WithField(key string, value interface{}) ILogger {
	return &Logger{l: l.l.With(zap.Any(key, value))}
}

// WithFields adds multiple fields to the logger.
func (l *Logger) WithFields(fields map[string]interface{}) ILogger {
	zapFields := make([]zap.Field, 0, len(fields))
	for key, value := range fields {
		zapFields = append(zapFields, zap.Any(key, value))
	}
	return &Logger{l: l.l.With(zapFields...)}
}

// Info logs a message at the Info level.
func (l *Logger) Info(msg string) {
	l.l.Info(msg)
}

// Error logs a message at the Error level.
func (l *Logger) Error(msg string) {
	l.l.Error(msg)
}

// Error logs a message at the Error level.
func (l *Logger) Panic(msg string) {
	l.l.Panic(msg)
}
