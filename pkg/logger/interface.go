package logger

type ILogger interface {
	WithField(key string, value interface{}) ILogger
	WithFields(fields map[string]interface{}) ILogger
	Info(msg string)
	Error(msg string)
	Panic(msg string)
	// Add other log levels and methods as needed}
}
