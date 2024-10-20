# Data Collector

The Data Collector is an open source observability framework for collecting telemetry data such as logs and transactions.
It provides a set of API to directly add structured and level logging to your application, as grouping several logs into transactions.

Logs are used for debugging, tracking system behavior, and providing insights into how the application is functioning or why failures occurred.


# Requirements

Use go 1.22.5 or newer.


# Debug logger issues
stderr
clients output


# Drivers
The available drivers default to Plaint Text encoding. For changing the encoding use the SetEncoding(encoding string) method from the driver inte


Over time, the log file could grow very large, which could lead to issues with storage or performance. Many logging libraries handle this by rotating log files based on size, time, or both.

	•	Opening and closing a file are relatively expensive operations. In each Log call, you would open the file, write a single log entry, and then close the file. This introduces significant overhead, especially for high-frequency logging (e.g., in microservices or applications that log a lot of data).

    the file needs to be closed by the client.

    Based on the code snippets provided, your package already implements:
Efficient Encoding: By supporting both plain text and JSON encodings.
Leveled Logging: By checking log levels before processing entries.
Buffered I/O: By using a buffered writer for file operations.

Further improvements:
- add more documentation.
- run in a container with all the dependencies installed.
- improve error handling (return errors and handle them at a higher level)
- 