package app

import (
	"github.com/ralugr/datacollector/pkg/config"
	"github.com/ralugr/datacollector/pkg/log"
)

type App struct {
	*application
}

type Driver interface {
	RecordLog(logInfo log.Entry)
	SetEncoding(encoding string)
}

func NewDataCollector(driver Driver, opts ...config.ConfigOption) (*App, error) {
	cfg := config.DefaultConfig()
	for _, fn := range opts {
		if fn != nil {
			fn(&cfg)
			if cfg.Error != nil {
				return nil, cfg.Error
			}
		}
	}

	return &App{
		application: newApplication(driver, cfg),
	}, nil
}

func (a *App) StartTransaction(attributes ...log.Attrb) *Transaction {

	return a.startTransaction(attributes...)
}

// check log level
func (a *App) Debug(msg string, attributes ...log.Attrb) {
	a.log(log.DebugLevel, msg, attributes...)
}

func (a *App) Info(msg string, attributes ...log.Attrb) {
	a.log(log.InfoLevel, msg, attributes...)
}

func (a *App) Warning(msg string, attributes ...log.Attrb) {
	a.log(log.WarnLevel, msg, attributes...)
}

func (a *App) Error(msg string, attributes ...log.Attrb) {
	a.log(log.ErrorLevel, msg, attributes...)
}
