package file

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ralugr/datacollector/pkg/log"
)

const (
	PlainEncoding = "plain"
	JSONEncoding  = "json"
)

const maxFileSize = 10 * 1024 * 1024 // 10MB

// supports json and plain text
type Writer struct {
	encoding    string
	file        *os.File
	fileName    string
	currentSize int64 // why
	buffer      *bufio.Writer
	mu          sync.Mutex
}

func NewWriter(fileName string) *Writer {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file %v: %v\n", fileName, err)
		return nil
	}

	return &Writer{
		encoding:    PlainEncoding,
		file:        file,
		fileName:    fileName,
		currentSize: 0,
		buffer:      bufio.NewWriter(file),
	}
}

// Close method closes the file
func (w *Writer) Close() {
	w.buffer.Flush()
	w.file.Close()
}

func (w *Writer) SetEncoding(encoding string) {
	if encoding != PlainEncoding && encoding != JSONEncoding {
		fmt.Fprintf(os.Stderr, "Unknown encoding %v\n", encoding)
		return
	}
	w.encoding = encoding
}

func (w *Writer) RecordLog(logInfo log.Entry) {
	var err error
	var bytes int

	w.mu.Lock()
	defer w.mu.Unlock()

	switch w.encoding {
	case PlainEncoding:
		bytes, err = w.buffer.WriteString(w.logEntryToString(logInfo) + "\n")
	case JSONEncoding:
		bytes, err = w.buffer.WriteString(w.logEntryToJson(logInfo) + "\n")
	default:
		fmt.Fprintf(os.Stderr, "Unknown encoding: %v. Defaulting to %v\n", w.encoding, PlainEncoding)
		bytes, err = w.buffer.WriteString(w.logEntryToString(logInfo) + "\n")
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
		return
	}

	if logInfo.Level == log.ErrorLevel {
		w.buffer.Flush()
	}

	w.currentSize += int64(bytes) // includes buffer size as well
	if w.currentSize > maxFileSize {
		err := w.rotateFile()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error rotating log file: %v\n", err)
		}
	}

}

func (w *Writer) logEntryToString(log log.Entry) string {
	s := fmt.Sprintf("time:%v, level:%v, app_name:%v, message:%v, attributes:%v", log.Timestamp.UTC().Format(time.RFC3339), log.Level,
		log.AppName, log.Message, log.Attributes)

	if log.TransactionID != "" {
		s += fmt.Sprintf(" transaction_id:%v", log.TransactionID)
	}

	return s
}

func (w *Writer) logEntryToJson(log log.Entry) string {
	jsonData, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to marshal %v: %v\n", log, err)
		return ""
	}

	return string(jsonData)
}

// rotateFile closes the current file, rename and opens a new file
func (w *Writer) rotateFile() error {
	w.buffer.Flush()
	w.file.Close()

	newFilename := fmt.Sprintf("%s.%d", w.fileName, time.Now().Unix())
	err := os.Rename(w.fileName, newFilename)
	if err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}

	file, err := os.OpenFile(w.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	w.file = file
	w.buffer = bufio.NewWriter(file)
	w.currentSize = 0
	return nil
}
