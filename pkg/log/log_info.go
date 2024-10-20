package log

import (
	"slices"
	"time"
)

type Entry struct {
	Timestamp     time.Time `json:"timestamp"`
	Level         Level     `json:"level"`
	AppName       string    `json:"app_name"`
	Message       string    `json:"message"`
	Attributes    []Attrb   `json:"attributes,omitempty"`
	TransactionID string    `json:"transaction_id,omitempty"`
}

type Attrb struct {
	Key   string
	Value any
}

// Define the exported LogLevel type as a string
type Level string

// Exported constants for valid log levels
const (
	DebugLevel Level = "DEBUG"
	InfoLevel  Level = "INFO"
	WarnLevel  Level = "WARNING"
	ErrorLevel Level = "ERROR"
)

var logLevels = []Level{DebugLevel, InfoLevel, WarnLevel, ErrorLevel}

// Single function to create a log attribute
func Attr(key string, value any) Attrb {
	return Attrb{Key: key, Value: value}
}

// make private
// check how zap logger is implemented because it uses zapcore for fields and zap for logging
func IsValid(appLogLevel Level, currentLogLevel Level) bool {
	return slices.Index(logLevels, currentLogLevel) >= slices.Index(logLevels, appLogLevel)
}
