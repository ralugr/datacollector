package file

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ralugr/datacollector/pkg/log"
	"github.com/stretchr/testify/assert"
)

func TestNewWriter(t *testing.T) {
	tmpFile := filepath.Join(os.TempDir(), "testlog.txt")
	defer os.Remove(tmpFile)

	writer, err := NewWriter(tmpFile)
	assert.NoError(t, err)
	assert.NotNil(t, writer)
	assert.Equal(t, PlainEncoding, writer.encoding)
}

func TestSetEncoding(t *testing.T) {
	tmpFile := filepath.Join(os.TempDir(), "testlog.txt")
	defer os.Remove(tmpFile)

	writer, err := NewWriter(tmpFile)
	assert.NoError(t, err)

	writer.SetEncoding(JSONEncoding)
	assert.Equal(t, JSONEncoding, writer.encoding)

	writer.SetEncoding("unknown")
	assert.Equal(t, JSONEncoding, writer.encoding) // Encoding shouldn't change
}

func TestRecordLogPlainEncoding(t *testing.T) {
	tmpFile := filepath.Join(os.TempDir(), "testlog.txt")
	defer os.Remove(tmpFile)

	writer, err := NewWriter(tmpFile)
	assert.NoError(t, err)

	entry := log.Entry{
		Timestamp:  time.Now(),
		Level:      log.InfoLevel,
		AppName:    "TestApp",
		Message:    "Test log message",
		Attributes: []log.Attrb{log.Attr("key", "value")},
	}

	writer.RecordLog(entry)
	writer.Close()

	content, err := os.ReadFile(tmpFile)
	assert.NoError(t, err)

	expected := "time:" + entry.Timestamp.UTC().Format(time.RFC3339) +
		", level:INFO, app_name:TestApp, message:Test log message, attributes:[{key value}]\n"
	assert.Contains(t, string(content), expected)
}

func TestRecordLogJSONEncoding(t *testing.T) {
	tmpFile := filepath.Join(os.TempDir(), "testlog.txt")
	defer os.Remove(tmpFile)

	writer, err := NewWriter(tmpFile)
	assert.NoError(t, err)

	writer.SetEncoding(JSONEncoding)

	entry := log.Entry{
		Timestamp:  time.Now(),
		Level:      log.InfoLevel,
		AppName:    "TestApp",
		Message:    "Test log message",
		Attributes: []log.Attrb{log.Attr("key", "value")},
	}

	writer.RecordLog(entry)
	writer.Close()

	content, err := os.ReadFile(tmpFile)
	assert.NoError(t, err)

	assert.Contains(t, string(content), `"app_name": "TestApp"`)
	assert.Contains(t, string(content), `"message": "Test log message"`)
}

func TestRotateFile(t *testing.T) {
	tmpFile := filepath.Join(os.TempDir(), "testlog.txt")
	defer os.Remove(tmpFile)

	writer, err := NewWriter(tmpFile)
	assert.NoError(t, err)

	// Simulate exceeding max file size
	writer.currentSize = maxFileSize + 1

	entry := log.Entry{
		Timestamp:  time.Now(),
		Level:      log.InfoLevel,
		AppName:    "TestApp",
		Message:    "Test log message",
		Attributes: []log.Attrb{log.Attr("key", "value")},
	}

	writer.RecordLog(entry)
	writer.Close()

	rotatedFile := fmt.Sprintf("%s.%d", tmpFile, time.Now().Unix())

	assert.FileExists(t, rotatedFile)
}

func TestClose(t *testing.T) {
	tmpFile := filepath.Join(os.TempDir(), "testlog.txt")
	defer os.Remove(tmpFile)

	writer, err := NewWriter(tmpFile)
	assert.NoError(t, err)

	writer.Close()

	_, err = writer.file.Write([]byte("test"))
	assert.Error(t, err)
}
