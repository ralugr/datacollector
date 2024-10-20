package main

import (
	"sync"
	"time"

	"github.com/ralugr/datacollector/pkg/app"
	"github.com/ralugr/datacollector/pkg/config"
	"github.com/ralugr/datacollector/pkg/drivers/file"
	"github.com/ralugr/datacollector/pkg/log"
)

func main() {
	driver := file.NewWriter("log_file.txt") // Use pointer
	defer driver.Close()

	// Initialize DataCollector (also passed by pointer)
	app, err := app.NewDataCollector(
		driver, // Pass by pointer
		config.ConfigAppName("Example"),
		config.ConfigLogLevel(log.DebugLevel),
	)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	// Spawn multiple goroutines for concurrent logging
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// Create a log entry (can be passed by value)
			app.Debug("Concurrent log message",
				log.Attr("goroutine", i),
				log.Attr("timestamp", time.Now().String()))
		}(i)
	}

	wg.Wait() // Wait for all goroutines to finish
}
