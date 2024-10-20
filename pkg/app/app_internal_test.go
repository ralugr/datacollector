package app

import (
	"testing"
	"time"

	"github.com/ralugr/datacollector/pkg/config"
	"github.com/ralugr/datacollector/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDriver struct {
	mock.Mock
}

func (m *MockDriver) SetEncoding(encoding string) {
	m.Called(encoding)
}

func (m *MockDriver) RecordLog(entry log.Entry) {
	m.Called(entry)
}

func TestNewApplication(t *testing.T) {
	driver := new(MockDriver)
	cfg := config.DefaultConfig()

	app := newApplication(driver, cfg)

	assert.Equal(t, driver, app.drv, "Driver should be initialized correctly")
	assert.Equal(t, cfg, app.config, "Config should be initialized correctly")
}

func TestApplicationStartTransaction(t *testing.T) {
	driver := new(MockDriver)
	cfg := config.DefaultConfig()

	app := newApplication(driver, cfg)

	attributes := []log.Attrb{log.Attr("key", "value")}
	txn := app.startTransaction(attributes...)

	assert.NotNil(t, txn, "Transaction should be initialized")
	assert.Equal(t, attributes, txn.attr, "Transaction attributes should be set correctly")
	assert.Equal(t, app.drv, txn.drv, "Driver should be passed to the transaction")
	assert.Equal(t, app.config, txn.config, "Config should be passed to the transaction")
}

func TestApplicationLogValid(t *testing.T) {
	driver := new(MockDriver)
	cfg := config.DefaultConfig()
	cfg.LogLevel = log.DebugLevel

	app := newApplication(driver, cfg)

	expectedEntry := log.Entry{
		Timestamp:  time.Now(),
		Level:      log.DebugLevel,
		AppName:    cfg.AppName,
		Message:    "Test debug message",
		Attributes: []log.Attrb{log.Attr("key", "value")},
	}

	driver.On("RecordLog", mock.MatchedBy(func(entry log.Entry) bool {
		return entry.Level == expectedEntry.Level &&
			entry.AppName == expectedEntry.AppName &&
			entry.Message == expectedEntry.Message
	})).Return()

	app.log(log.DebugLevel, "Test debug message", log.Attr("key", "value"))

	driver.AssertCalled(t, "RecordLog", mock.MatchedBy(func(entry log.Entry) bool {
		return entry.Level == expectedEntry.Level && entry.Message == expectedEntry.Message
	}))
}

func TestApplicationLogInvalid(t *testing.T) {
	driver := new(MockDriver)
	cfg := config.DefaultConfig()
	cfg.LogLevel = log.InfoLevel

	app := newApplication(driver, cfg)

	driver.On("RecordLog", mock.Anything).Return()

	app.log(log.DebugLevel, "This should not log", log.Attr("key", "value"))

	driver.AssertNotCalled(t, "RecordLog", mock.Anything)
}

func TestApplicationLogConcurrency(t *testing.T) {
	driver := new(MockDriver)
	cfg := config.DefaultConfig()
	cfg.LogLevel = log.DebugLevel

	app := newApplication(driver, cfg)

	expectedEntry := log.Entry{
		Timestamp:  time.Now(),
		Level:      log.DebugLevel,
		AppName:    cfg.AppName,
		Message:    "Concurrent log message",
		Attributes: []log.Attrb{log.Attr("key", "value")},
	}

	driver.On("RecordLog", mock.AnythingOfType("log.Entry")).Return()

	done := make(chan bool)
	go func() {
		app.log(log.DebugLevel, "Concurrent log message", log.Attr("key", "value"))
		done <- true
	}()
	go func() {
		app.log(log.DebugLevel, "Concurrent log message", log.Attr("key", "value"))
		done <- true
	}()

	<-done
	<-done

	driver.AssertCalled(t, "RecordLog", mock.MatchedBy(func(entry log.Entry) bool {
		return entry.Level == expectedEntry.Level && entry.Message == expectedEntry.Message
	}))
}
