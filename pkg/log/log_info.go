package log

import (
	"slices"
	"time"
)

// Entry represents a log entry containing metadata about a specific log event.
// Includes json formatting
type Entry struct {
	Timestamp     time.Time `json:"timestamp"`
	Level         Level     `json:"level"`
	AppName       string    `json:"app_name"`
	Message       string    `json:"message"`
	Attributes    []Attrb   `json:"attributes,omitempty"`
	TransactionID string    `json:"transaction_id,omitempty"`
}

// Attrb represents a single key-value pair for log attributes.
type Attrb struct {
	Key   string
	Value any
}

// Level defines a custom type.
// See constants below for available levels.
type Level string

// Predefined log levels for standard log severity.
// Should be passed when creating a new data collector application.
const (
	DebugLevel Level = "DEBUG"
	InfoLevel  Level = "INFO"
	WarnLevel  Level = "WARNING"
	ErrorLevel Level = "ERROR"
)

// A list of log levels in order of increasing severity.
// This is used internally to compare log levels.
var logLevels = []Level{DebugLevel, InfoLevel, WarnLevel, ErrorLevel}

// Attr creates a new key-value pair attribute for a log entry.
// This function should be used to add attributes to logs and transactions.
func Attr(key string, value any) Attrb {
	return Attrb{Key: key, Value: value}
}

// IsValid determines if a log entry's level is valid according to the application's logging configuration.
// It compares the current log level with the application's configured log level.
func IsValid(appLogLevel Level, currentLogLevel Level) bool {
	appLevelIndex := slices.Index(logLevels, appLogLevel)
	currentLevelIndex := slices.Index(logLevels, currentLogLevel)

	return appLevelIndex >= 0 && currentLevelIndex >= appLevelIndex
}
