# Data Collector

The Data Collector Application is an open source observability framework for collecting telemetry data such as logs and transactions.
It provides a set of APIs to directly add structured and leveled logging to your application, as well as grouping several logs into transactions.


## Features

* Structured logging with plain text or JSON encoding
* Leveled logging (Debug, Info, Warning, Error)
* Transaction-based logging
* Supports multiple drivers (CLI, File)
* Extensible with new drivers
* Customizable through config options
* Log file rotation - to avoid large log files
* Buffered I/O - to improve performance


## Installation

### Compatibility and Requirements

For the latest version of the library, Go 1.22+ is required.

### Getting started
Follow these steps to instrument your application. Use the [Examples](examples) as a starting point.

**Step 1: Installation**

Data Collector is a Go library. See dependencies in [go.mod](go.mod). Install the Data Collectorthe the same way you would install any other Go library. 
The simplest way is to run:

```
go get github.com/ralugr/datacollector
```

Then import the package in your application:
```go
import "github.com/ralugr/datacollector"
```

**Step 2: Create an Application**

In your `main` function, or an `init` block, create a [NewDataCollector](pkg/app/app.go) using [Config](pkg/config/config.go).

```go
func main() {
    // Initialize a driver(cli.Writer, file.Writer or create a custom driver)
    driver := &cli.Writer{}

    // Set the encoding
	driver.SetEncoding("plain")

    // Create an Application:
    app, err := app.NewDataCollector(
        // Pass the driver
        driver,
        // Name your application
        config.AppName("CLI Plain"),
        // Set the application log level
        config.LogLevel(log.DebugLevel),
    )
    // If an application could not be created then err will reveal why.
    if err != nil {
        fmt.Println("unable to create Data Collector", err)
    }
    // Now use the app to instrument everything!
}
```

**Step 3: Add Logging**
Use the Debug, Info, Warning, Error methods from the Application instance to log messages.

```go
func main() {
    // Initialize a driver(cli.Writer, file.Writer or create a custom driver)
    driver := &cli.Writer{}

    // Set the encoding
	driver.SetEncoding("plain")

    // Create an Application:
    app, err := app.NewDataCollector(
        // Pass the driver
        driver,
        // Name your application
        config.AppName("CLI Plain"),
        // Set the application log level
        config.LogLevel(log.DebugLevel),
    )
    // If an application could not be created then err will reveal why.
    if err != nil {
        fmt.Println("unable to create Data Collector", err)
    }
    // Use app.Debug, app.Info, app.Warning, app.Error to log messages.
    // Attributes are optional but can be used to add context to the log message.
    // Use log.Attr to create a new attribute with a string key and any value.
    app.Debug("Application started",
		log.Attr("userID", "12345"),
		log.Attr("attempt", 3),
		log.Attr("success", true),
	)
}
```

**Step 4: Add Transactions**
[Transactions](pkg/app/transaction.go) are a way to group several logs into a single unit of work.

```go
func main() {
    // Initialize a driver(cli.Writer, file.Writer or create a custom driver)
    driver := &cli.Writer{}

    // Set the encoding
	driver.SetEncoding("plain")

    // Create an Application:
    app, err := app.NewDataCollector(
        // Pass the driver
        driver,
        // Name your application
        config.AppName("CLI Plain"),
        // Set the application log level
        config.LogLevel(log.DebugLevel),
    )
    // If an application could not be created then err will reveal why.
    if err != nil {
        fmt.Println("unable to create Data Collector", err)
    }
    // Use app.Debug, app.Info, app.Warning, app.Error to log messages.
    // Attributes are optional but can be used to add context to the log message.
    // Use log.Attr to create a new attribute with a string key and any value.
    app.Debug("Application started",
		log.Attr("userID", "12345"),
		log.Attr("attempt", 3),
		log.Attr("success", true),
	)

    // Start a transaction
    // The transaction instance can be used to log messages and end the transaction.
    transaction := app.StartTransaction()

    // Log a message inside the transaction
	transaction.Debug("Transaction started",
		log.Attr("database_name", "products"),
		log.Attr("active_connections", 5),
		log.Attr("sql", false))

    // End the transaction
	transaction.End()

    // Logging to an inactive transaction will be ignored.
	transaction.Info("Attemping to write to a finished transaction")
}
```


## Drivers
Data Collector has some predefined drivers that can be plugged in to the application, but custom drivers can be created by implementing the [Driver](pkg/app/app.go/#Driver) interface.

The driver processes the logs by validating and converting them to the required format.
Then it outputs the logs to the desired location.

### Default Drivers

The data collector app comes with two default drivers: 
  * [cli.Writer](pkg/drivers/cli/writer.go) - for logging to the console
  * [file.Writer](pkg/drivers/file/writer.go) - for logging to a file
  
These drivers support both plain text and JSON encoding through the `SetEncoding` function. 
Their default encoding is plain text.

**CLI Driver Example**

```go
func main() {
    // Initialize a driver console writer
    driver := &cli.Writer{}

    // Set the encoding to plain text
	driver.SetEncoding("plain")

    // Alternatively, the client can set the encoding to json
	//driver.SetEncoding("json")

    // Create an Application:
    app, err := app.NewDataCollector(
        // Pass the driver
        driver,
        // Name your application
        config.AppName("CLI Plain"),
        // Set the application log level
        config.LogLevel(log.DebugLevel),
    )
    // If an application could not be created then err will reveal why.
    if err != nil {
        fmt.Println("unable to create Data Collector", err)
    }
    // Now use the app to instrument everything!
}
```

**File Driver Example**

```go
func main() {
    // Initialize a file writer as a driver
    driver, err := file.NewWriter("log_json.txt")

    // If the driver could not be created then err will reveal why.
    if err != nil {
        fmt.Println("unable to create Data Collector", err)
    }
    // Very important to close the driver when done
    defer driver.Close()

    // Set the encoding to plain text
	driver.SetEncoding("plain")

    // Alternatively, the client can set the encoding to json
	//driver.SetEncoding("json")

    // Create an Application:
    app, err := app.NewDataCollector(
        // Pass the driver
        driver,
        // Name your application
        config.AppName("File Driver"),
        // Set the application log level
        config.LogLevel(log.DebugLevel),
    )
    // If an application could not be created then err will reveal why.
    if err != nil {
        fmt.Println("unable to create Data Collector", err)
    }
    // Now use the app to instrument everything!
}
```

**Custom Driver Example**

```go
// Define a custom driver typw
type CustomDriver struct{}

// Implement the Driver interface: RecordLog
func (d *CustomDriver) RecordLog(logInfo log.Entry) {
	fmt.Printf("Custom driver log - %v\n", logInfo)
}

// Implement the Driver interface: SetEncoding
func (d *CustomDriver) SetEncoding(encoding string) {
	// No encoding needed for this example
}

func main() {
    // Initialize a custom driver
	driver := &CustomDriver{}
    
    // Create an Application
	app, err := app.NewDataCollector(
        // Pass the custom driver
		driver,
        // Name your application
		config.AppName("Custom Driver Example"),
        // Set the application log level
		config.LogLevel(log.DebugLevel),
	)
    // If an application could not be created then err will reveal why.
    if err != nil {
        fmt.Println("unable to create Data Collector", err)
    }
    // Now use the app to instrument everything!
}
```


## Roadmap
* Include transaction attributes in the structured logs.
* Add integration tests.
* Increase test coverage.
* Document the code.
* Create a containerized environment for testing.
* Improve error handling.
* Add github actions automated deploy and testing.
* Add versioning and update changelog.