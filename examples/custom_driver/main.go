package main

import (
	"fmt"
	"os"

	"github.com/ralugr/datacollector/pkg/app"
	"github.com/ralugr/datacollector/pkg/config"
	"github.com/ralugr/datacollector/pkg/log"
)

type CustomDriver struct{}

func (d *CustomDriver) RecordLog(logInfo log.Entry) {
	fmt.Printf("Custom driver log - %v\n", logInfo)
}

func (d *CustomDriver) SetEncoding(encoding string) {
	// No encoding needed for this example
}

func main() {
	driver := &CustomDriver{}
	app, err := app.NewDataCollector(
		driver,
		config.ConfigAppName("Custom Driver Example"),
		config.ConfigLogLevel(log.DebugLevel),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	app.Info("Application started")

	app.Error("Sample error message")
}
