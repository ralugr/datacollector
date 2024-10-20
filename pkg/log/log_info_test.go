package log

import (
	"testing"
	"time")

func TestAttr(t *testing.T) {
	tests := []struct {
		key      string
		value    any
		expected Attrb
	}{
		{"key1", "value1", Attrb{Key: "key1", Value: "value1"}},
		{"key2", 123, Attrb{Key: "key2", Value: 123}},
		{"empty", nil, Attrb{Key: "empty", Value: nil}},
	}

	for _, tt := range tests {
		attr := Attr(tt.key, tt.value)
		if attr.Key != tt.expected.Key || attr.Value != tt.expected.Value {
			t.Errorf("Attr() failed. Expected %v, got %v", tt.expected, attr)
		}
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		appLogLevel     Level
		currentLogLevel Level
		expected        bool
	}{
		{DebugLevel, DebugLevel, true},
		{InfoLevel, WarnLevel, true},
		{WarnLevel, DebugLevel, false},
		{ErrorLevel, DebugLevel, false},
		{ErrorLevel, ErrorLevel, true},
		{"UNKNOWN", DebugLevel, false}, // Invalid app log level
		{InfoLevel, "UNKNOWN", false},  // Invalid current log level
	}

	for _, tt := range tests {
		result := IsValid(tt.appLogLevel, tt.currentLogLevel)
		if result != tt.expected {
			t.Errorf("IsValid() failed. For appLogLevel: %v and currentLogLevel: %v, expected %v, got %v", tt.appLogLevel, tt.currentLogLevel, tt.expected, result)
		}
	}
}

func TestEntryInitialization(t *testing.T) {
	timestamp := time.Now()

	entry := Entry{
		Timestamp:     timestamp,
		Level:         InfoLevel,
		AppName:       "TestApp",
		Message:       "Test Message",
		Attributes:    []Attrb{{Key: "key", Value: "value"}},
		TransactionID: "12345",
	}

	if entry.Timestamp != timestamp {
		t.Errorf("Entry Timestamp mismatch. Expected %v, got %v", timestamp, entry.Timestamp)
	}

	if entry.Level != InfoLevel {
		t.Errorf("Entry Level mismatch. Expected %v, got %v", InfoLevel, entry.Level)
	}

	if entry.AppName != "TestApp" {
		t.Errorf("Entry AppName mismatch. Expected %v, got %v", "TestApp", entry.AppName)
	}

	if entry.Message != "Test Message" {
		t.Errorf("Entry Message mismatch. Expected %v, got %v", "Test Message", entry.Message)
	}

	if entry.TransactionID != "12345" {
		t.Errorf("Entry TransactionID mismatch. Expected %v, got %v", "12345", entry.TransactionID)
	}

	if len(entry.Attributes) != 1 || entry.Attributes[0].Key != "key" || entry.Attributes[0].Value != "value" {
		t.Errorf("Entry Attributes mismatch. Expected key-value pair {key, value}, got %v", entry.Attributes)
	}
}
