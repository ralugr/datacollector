# Run unit tests
test:
	go test ./... -v

# Example: CLI output with json format
cli-json:
	go run ./examples/cli_json/main.go

# Example: CLI output with plain text format 
cli-plain:
	go run ./examples/cli_plain/main.go

# Example: Custom driver
custom-driver:
	go run ./examples/custom_driver/main.go

# Example: File output json format
file-json:
	go run ./examples/file_json/main.go

# Example: File output plain text format
file-plain:
	go run ./examples/file_plain/main.go

# Example: Multi threaded example
multi-thread:
	go run ./examples/multi_thread/main.go

# Cleanup step to remove all .txt files in the current directory
cleanup:
	@find . -maxdepth 1 -name "*.txt" -type f -exec rm -f {} \;
	@echo "All .txt files in the current directory have been removed."