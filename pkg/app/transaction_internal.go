package app

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ralugr/datacollector/pkg/config"
	"github.com/ralugr/datacollector/pkg/log"
)

type txn struct {
	id     string
	drv    Driver
	config config.Config
	attr   []log.Attrb
	active bool
	mu     sync.Mutex
}

func newPrivateTxn(driver Driver, cfg config.Config, attributes ...log.Attrb) *txn {
	id, err := generateID()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to generate transaction ID due to error: %v. Please try again.\n", err)
		return &txn{
			active: false,
		}
	}
	return &txn{
		id:     id,
		drv:    driver,
		config: cfg,
		attr:   attributes,
		active: true,
	}
}

func (t *txn) end() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.active = false
}

func (t *txn) log(level log.Level, msg string, attributes ...log.Attrb) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.active {
		t.drv.RecordLog(log.Entry{
			Timestamp: time.Now(),
			Level:     log.ErrorLevel,
			AppName:   appName,
			Message:   "Transaction already ended!"})
		return
	}

	if !log.IsValid(t.config.LogLevel, level) {
		return
	}

	data := log.Entry{
		Timestamp:     time.Now(),
		Level:         level,
		AppName:       t.config.AppName,
		Message:       msg,
		Attributes:    attributes,
		TransactionID: t.id,
	}

	t.drv.RecordLog(data)
}

func generateID() (string, error) {
	timestamp := time.Now().UnixNano()

	randomBytes := make([]byte, 4)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("unable to generate transaction id: %v", err)
	}

	return fmt.Sprintf("%x%s", timestamp, hex.EncodeToString(randomBytes)), nil
}
