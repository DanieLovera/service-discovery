package results

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type Writer struct {
	file   *os.File
	writer *csv.Writer
	mu     sync.Mutex
}

func Create(dir, name string, header ...string) (*Writer, error) {
	file, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		return nil, fmt.Errorf("create %s: %w", name, err)
	}

	writer := csv.NewWriter(file)
	if err := writer.Write(header); err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("write %s header: %w", name, err)
	}

	return &Writer{file: file, writer: writer}, nil
}

func (w *Writer) Write(record ...string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	_ = w.writer.Write(record)
}

func (w *Writer) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.writer.Flush()

	return errors.Join(w.writer.Error(), w.file.Close())
}

func FormatTimestamp(t time.Time) string {
	micros := t.UnixMicro()
	return fmt.Sprintf("%d.%03d", micros/1000, micros%1000)
}

func FormatMilliseconds(d time.Duration) string {
	return strconv.FormatFloat(float64(d)/float64(time.Millisecond), 'f', 3, 64)
}
