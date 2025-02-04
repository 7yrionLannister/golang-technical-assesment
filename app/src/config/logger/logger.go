package logger

type Logger interface {
	InitLogger(levelStr string)    // Initialize the logger with the given level
	Debug(msg string, args ...any) // Log a debug message
	Info(msg string, args ...any)  // Log an info message
	Warn(msg string, args ...any)  // Log a warning message
	Error(msg string, args ...any) // Log an error message
}

var L Logger = &slogLogger{} // Global logger
