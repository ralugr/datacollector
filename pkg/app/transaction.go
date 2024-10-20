package app

import (
	"github.com/ralugr/datacollector/pkg/config"
	"github.com/ralugr/datacollector/pkg/log"
)

// Transaction represents a loggable transaction within the App.
// It enables logging at various levels (Debug, Info, Warning, Error).
// Make sure to call the End() function when the transaction is not needed.
type Transaction struct {
	*txn
}

// End marks the transaction as inactive and prevents further logging within this transaction.
// Once called, any subsequent attempts to log in this transaction will result in an error log entry.
func (t *Transaction) End() {
	t.end()
}

// Debug logs a message at the Debug level within the context of this transaction.
// Optional attributes can be provided for additional context.
func (t *Transaction) Debug(msg string, attributes ...log.Attrb) {

	t.log(log.DebugLevel, msg, attributes...)
}

// Info logs a message at the Info level within the context of this transaction.
// Optional attributes can be provided for additional context.
func (t *Transaction) Info(msg string, attributes ...log.Attrb) {
	t.log(log.InfoLevel, msg, attributes...)
}

// Warning logs a message at the Warning level within the context of this transaction.
// Optional attributes can be provided for additional context.
func (t *Transaction) Warning(msg string, attributes ...log.Attrb) {
	t.log(log.WarnLevel, msg, attributes...)
}

// Error logs a message at the Error level within the context of this transaction.
// Optional attributes can be provided for additional context.
func (t *Transaction) Error(msg string, attributes ...log.Attrb) {
	t.log(log.ErrorLevel, msg, attributes...)
}

// newTransaction creates and initializes a new Transaction instance.
// It takes a logging driver, configuration, and optional attributes for logging.
// The new transaction is returned to app.StartTransaction(), ready to log messages and eventually be ended.
func newTransaction(driver Driver, cfg config.Config, attributes ...log.Attrb) *Transaction {
	return &Transaction{
		txn: newPrivateTxn(driver, cfg, attributes...),
	}
}
