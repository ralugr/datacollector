package app

import (
	"testing"
	"time"

	"github.com/ralugr/datacollector/pkg/config"
	"github.com/ralugr/datacollector/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTxnEnd(t *testing.T) {
	cfg := config.DefaultConfig()
	driver := new(MockDriver)

	txn := newPrivateTxn(driver, cfg)

	assert.True(t, txn.active, "Transaction should be active when created")

	txn.end()

	assert.False(t, txn.active, "Transaction should be inactive after end")
}

func TestTxnLogWhenActive(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.LogLevel = log.DebugLevel
	driver := new(MockDriver)

	txn := newPrivateTxn(driver, cfg)

	expectedEntry := log.Entry{
		Timestamp:     time.Now(),
		Level:         log.DebugLevel,
		AppName:       cfg.AppName,
		Message:       "Test message",
		Attributes:    []log.Attrb{log.Attr("key", "value")},
		TransactionID: txn.id,
	}

	driver.On("RecordLog", mock.AnythingOfType("log.Entry")).Return()

	txn.log(log.DebugLevel, "Test message", log.Attr("key", "value"))

	driver.AssertCalled(t, "RecordLog", mock.MatchedBy(func(entry log.Entry) bool {
		return entry.Level == expectedEntry.Level &&
			entry.AppName == expectedEntry.AppName &&
			entry.Message == expectedEntry.Message &&
			entry.TransactionID == txn.id
	}))
}

func TestTxnLogWhenEnded(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.LogLevel = log.DebugLevel
	driver := new(MockDriver)

	txn := newPrivateTxn(driver, cfg)

	txn.end()

	expectedEntry := log.Entry{
		Timestamp: time.Now(),
		Level:     log.ErrorLevel,
		AppName:   cfg.AppName,
		Message:   "Transaction already ended!",
	}

	driver.On("RecordLog", mock.AnythingOfType("log.Entry")).Return()

	txn.log(log.DebugLevel, "This should not log")

	driver.AssertCalled(t, "RecordLog", mock.MatchedBy(func(entry log.Entry) bool {
		return entry.Level == expectedEntry.Level && entry.Message == expectedEntry.Message
	}))
}

func TestTxnLogInvalidLogLevel(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.LogLevel = log.InfoLevel // Only log Info level and above
	driver := new(MockDriver)

	txn := newPrivateTxn(driver, cfg)

	driver.On("RecordLog", mock.Anything).Return()

	txn.log(log.DebugLevel, "This should not log")

	driver.AssertNotCalled(t, "RecordLog", mock.Anything)
}

func TestGenerateID(t *testing.T) {
	id1, err1 := generateID()
	assert.NoError(t, err1)

	id2, err2 := generateID()
	assert.NoError(t, err2)

	assert.NotEqual(t, id1, id2, "Generated IDs should be unique")
}
