package cli

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/ralugr/datacollector/pkg/log"
)

const (
	PlainEncoding = "plain"
	JSONEncoding  = "json"
)

// supports json and plain text
type Writer struct {
	encoding string
	mu       sync.Mutex
}

func NewWriter() *Writer {
	return &Writer{
		encoding: PlainEncoding,
	}
}

func (w *Writer) SetEncoding(encoding string) {
	if encoding != PlainEncoding && encoding != JSONEncoding {
		return // thow an error?
	}
	w.encoding = encoding
}

func (w *Writer) RecordLog(logInfo log.Entry) {
	w.mu.Lock()
	defer w.mu.Unlock()

	switch w.encoding {
	case PlainEncoding:
		fmt.Println(w.logEntryToString(logInfo))
	case JSONEncoding:
		fmt.Println(w.logEntryToJson(logInfo))
	default:
		fmt.Println("Unknown encoding")
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
		fmt.Println("Error:", err) // throw an error?
		return ""
	}

	return string(jsonData)
}
