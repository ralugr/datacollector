package config

import (
	"testing"

	"github.com/ralugr/datacollector/pkg/log"
	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.Equal(t, "Test App", cfg.AppName, "Default AppName should be 'Test App'")
	assert.Equal(t, log.DebugLevel, cfg.LogLevel, "Default LogLevel should be DebugLevel")
	assert.Nil(t, cfg.Error, "Default Error should be nil")
}

func TestAppName(t *testing.T) {
	cfg := DefaultConfig()
	AppName("MyApp")(&cfg)

	assert.Equal(t, "MyApp", cfg.AppName, "AppName should be updated to 'MyApp'")
}

func TestLogLevelValidInputs(t *testing.T) {
	validLogLevels := []log.Level{log.DebugLevel, log.InfoLevel, log.WarnLevel, log.ErrorLevel}

	for _, level := range validLogLevels {
		cfg := DefaultConfig()
		LogLevel(level)(&cfg)

		assert.Equal(t, level, cfg.LogLevel, "LogLevel should be set to a valid level")
		assert.Nil(t, cfg.Error, "Error should be nil for valid LogLevel")
	}
}

func TestLogLevelInvalidInput(t *testing.T) {
	invalidLogLevel := log.Level("INVALID")

	cfg := DefaultConfig()
	LogLevel(invalidLogLevel)(&cfg)

	assert.NotNil(t, cfg.Error, "Error should not be nil for invalid LogLevel")
	assert.EqualError(t, cfg.Error, "Invalid value: INVALID", "Error message should indicate invalid log level")
	assert.Equal(t, log.DebugLevel, cfg.LogLevel, "LogLevel should not change for invalid input")
}

func TestLogLevelLowerBound(t *testing.T) {
	lowestValidLevel := log.DebugLevel

	cfg := DefaultConfig()
	LogLevel(lowestValidLevel)(&cfg)

	assert.Equal(t, lowestValidLevel, cfg.LogLevel, "LogLevel should be set to the lowest valid level")
	assert.Nil(t, cfg.Error, "Error should be nil for valid LogLevel")
}

func TestLogLevelUpperBound(t *testing.T) {
	highestValidLevel := log.ErrorLevel

	cfg := DefaultConfig()
	LogLevel(highestValidLevel)(&cfg)

	assert.Equal(t, highestValidLevel, cfg.LogLevel, "LogLevel should be set to the highest valid level")
	assert.Nil(t, cfg.Error, "Error should be nil for valid LogLevel")
}

func TestLogLevelCaseSensitivity(t *testing.T) {
	// Assuming log.Level is case-sensitive, test a lowercase equivalent
	lowercaseLogLevel := log.Level("debug")

	cfg := DefaultConfig()
	LogLevel(lowercaseLogLevel)(&cfg)

	assert.NotNil(t, cfg.Error, "Error should not be nil for invalid LogLevel case sensitivity")
	assert.EqualError(t, cfg.Error, "Invalid value: debug", "Error message should indicate invalid log level due to case sensitivity")
	assert.Equal(t, log.DebugLevel, cfg.LogLevel, "LogLevel should not change for invalid case sensitivity")
}

func TestCustomAppNameWithDefault(t *testing.T) {
	cfg := DefaultConfig()
	AppName("CustomApp")(&cfg)

	assert.Equal(t, "CustomApp", cfg.AppName, "AppName should be 'CustomApp'")
	assert.Nil(t, cfg.Error, "Error should remain nil when only changing AppName")
}
