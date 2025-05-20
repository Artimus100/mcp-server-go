package utils

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// LogLevel represents the severity level of a log message
type LogLevel int

// Log levels
const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
	FATAL
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger provides a simple logging interface
type Logger struct {
	prefix   string
	minLevel LogLevel
	logger   *log.Logger
	mu       sync.Mutex
}

var (
	// Default logger
	defaultLogger *Logger
	once          sync.Once
)

// initDefaultLogger initializes the default logger
func initDefaultLogger() {
	defaultLogger = &Logger{
		prefix:   "",
		minLevel: INFO,
		logger:   log.New(os.Stdout, "", 0),
	}
}

// NewLogger creates a new logger with the given prefix
func NewLogger(prefix string) *Logger {
	once.Do(initDefaultLogger)

	return &Logger{
		prefix:   prefix,
		minLevel: defaultLogger.minLevel,
		logger:   log.New(os.Stdout, "", 0),
	}
}

// WithPrefix returns a new logger with an additional prefix
func (l *Logger) WithPrefix(prefix string) *Logger {
	return &Logger{
		prefix:   fmt.Sprintf("%s.%s", l.prefix, prefix),
		minLevel: l.minLevel,
		logger:   l.logger,
	}
}

// SetLevel sets the minimum log level
func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.minLevel = level
}

// log logs a message with the given level and arguments
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if level < l.minLevel {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	prefix := l.prefix
	if prefix != "" {
		prefix = "[" + prefix + "] "
	}

	levelStr := level.String()
	message := fmt.Sprintf(format, args...)

	l.logger.Printf("%s %s %s%s", timestamp, levelStr, prefix, message)

	// If this is a fatal message, exit the program
	if level == FATAL {
		os.Exit(1)
	}
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warning logs a warning message
func (l *Logger) Warning(format string, args ...interface{}) {
	l.log(WARNING, format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// Fatal logs a fatal message and exits the program
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
}

// TODO: Consider adding additional features:
// - Log rotation
// - Output to multiple destinations (file, stdout, etc.)
// - Structured logging (JSON)
// - Log filtering by module/component
