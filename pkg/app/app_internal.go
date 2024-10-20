package app

import (
	"fmt"
	"sync"
	"time"

	"github.com/ralugr/datacollector/pkg/config"
	"github.com/ralugr/datacollector/pkg/log"
)

const appName = "Data Collector"

type application struct {
	drv    Driver
	config config.Config
	mu     sync.Mutex
}

func newApplication(driver Driver, cfg config.Config) *application {
	return &application{
		drv:    driver,
		config: cfg,
	}
}

func (a *application) startTransaction(attributes ...log.Attrb) *Transaction {
	return newTransaction(a.drv, a.config, attributes...)
}

// check log level
func (a *application) log(level log.Level, msg string, attributes ...log.Attrb) {
	if !log.IsValid(a.config.LogLevel, level) {
		fmt.Printf("Log is not valid %v, %v", a.config.LogLevel, level)
		return
	}

	data := log.Entry{
		Timestamp:  time.Now(),
		Level:      level,
		AppName:    a.config.AppName,
		Message:    msg,
		Attributes: attributes,
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	a.drv.RecordLog(data)
}
