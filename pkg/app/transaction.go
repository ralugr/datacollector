package app

import (
	"github.com/ralugr/datacollector/pkg/config"
	"github.com/ralugr/datacollector/pkg/log"
)

type Transaction struct {
	*txn
}

func (t *Transaction) End() {
	t.end()
}

func (t *Transaction) Debug(msg string, attributes ...log.Attrb) {

	t.log(log.DebugLevel, msg, attributes...)
}

func (t *Transaction) Info(msg string, attributes ...log.Attrb) {
	t.log(log.InfoLevel, msg, attributes...)
}

func (t *Transaction) Warning(msg string, attributes ...log.Attrb) {
	t.log(log.WarnLevel, msg, attributes...)
}

func (t *Transaction) Error(msg string, attributes ...log.Attrb) {
	t.log(log.ErrorLevel, msg, attributes...)
}

func newTransaction(driver Driver, cfg config.Config, attributes ...log.Attrb) *Transaction {
	return &Transaction{
		txn: newPrivateTxn(driver, cfg, attributes...),
	}
}
