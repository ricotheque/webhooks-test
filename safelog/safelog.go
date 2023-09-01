// Package safelog provides a thread-safe way to write logs to a file.
package safelog

import (
	"fmt"
	"os"
	"sync"
)

type SafeLogger struct {
	mu   sync.Mutex
	file *os.File
}

var (
	defaultLogger *SafeLogger
	once          sync.Once
)

// New creates a new SafeLogger. It returns an error if the file cannot be opened.
func New(filename string) (*SafeLogger, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &SafeLogger{file: file}, nil
}

// InitDefaultLogger initializes the default logger instance.
func InitDefaultLogger(filename string) error {
	var err error
	once.Do(func() {
		defaultLogger, err = New(filename)
	})
	return err
}

// Log writes a message to the default log file. It is safe to call from multiple goroutines.
func Log(message string) error {
	if defaultLogger == nil {
		return fmt.Errorf("default logger not initialized")
	}
	return defaultLogger.Log(message)
}

// CloseDefaultLogger closes the underlying file of the default logger.
func CloseDefaultLogger() error {
	if defaultLogger == nil {
		return fmt.Errorf("default logger not initialized")
	}
	return defaultLogger.Close()
}

// Methods for SafeLogger instance

func (l *SafeLogger) Log(message string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	_, err := l.file.WriteString(fmt.Sprintf("%s\n", message))
	return err
}

func (l *SafeLogger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.file.Close()
}
