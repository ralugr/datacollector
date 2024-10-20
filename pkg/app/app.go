package app

import (
	"github.com/ralugr/datacollector/pkg/config"
	"github.com/ralugr/datacollector/pkg/log"
)

// App is a wrapper around the application struct that provides logging and transaction functionality.
// It enables logging at various levels (Debug, Info, Warning, Error) and supports starting new transactions
type App struct {
	*application
}

// Driver is an interface that defines the format and the output of the logs.
// - RecordLog: Used to record a log entry.
// - SetEncoding: Configures the encoding format for the log (e.g., JSON, plain text).
//
// Available drivers are - file.Writer for logging plain text or json logs into a file.
//   - cli.Writer for logging plain text or json into console output.
type Driver interface {
	RecordLog(logInfo log.Entry)
	SetEncoding(encoding string)
}

// NewDataCollector initializes a new App instance with a driver and configuration options.
// There are a couple of predefined drivers: file.Writer and cli.Writer or the user can define a custom deriver.
// If any configuration option sets an error, it returns that error and halts initialization.
// See example files for more information.
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

// StartTransaction begins a new transaction within the App.
// It accepts optional attributes that can be attached to the transaction for additional context.
func (a *App) StartTransaction(attributes ...log.Attrb) *Transaction {

	return a.startTransaction(attributes...)
}

// Debug logs a message at the Debug level.
// It allows optional attributes to be passed in for additional logging context.
func (a *App) Debug(msg string, attributes ...log.Attrb) {
	a.log(log.DebugLevel, msg, attributes...)
}

// Info logs a message at the Info level.
// It allows optional attributes to be passed in for additional logging context.
func (a *App) Info(msg string, attributes ...log.Attrb) {
	a.log(log.InfoLevel, msg, attributes...)
}

// Warning logs a message at the Warning level.
// It allows optional attributes to be passed in for additional logging context.
func (a *App) Warning(msg string, attributes ...log.Attrb) {
	a.log(log.WarnLevel, msg, attributes...)
}

// Error logs a message at the Error level.
// It allows optional attributes to be passed in for additional logging context.
func (a *App) Error(msg string, attributes ...log.Attrb) {
	a.log(log.ErrorLevel, msg, attributes...)
}
