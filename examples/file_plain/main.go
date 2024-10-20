package main

import (
	"fmt"
	"os"

	"github.com/ralugr/datacollector/pkg/app"
	"github.com/ralugr/datacollector/pkg/config"
	"github.com/ralugr/datacollector/pkg/drivers/file"
	"github.com/ralugr/datacollector/pkg/log"
)

func main() {
	driver, err := file.NewWriter("log_plain.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer driver.Close()
	driver.SetEncoding("plain")

	app, err := app.NewDataCollector(
		driver,
		config.AppName("Example"),
		config.LogLevel(log.DebugLevel),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	app.Debug("Application started",
		log.Attr("userID", "12345"),
		log.Attr("attempt", 3),
		log.Attr("success", true),
	)

	transaction := app.StartTransaction()
	transaction.Debug("Transaction started",
		log.Attr("database_name", "products"),
		log.Attr("active_connections", 5),
		log.Attr("sql", false))

	transaction.End()
	transaction.Info("Attemping to write to a finished transaction")
}
