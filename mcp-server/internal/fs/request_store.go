package fs

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// TODO: move to config
const logFileName = "requests.log"

type RequestStore struct {
	dir  string
	mu   sync.Mutex
	file *os.File
	buf  *bufio.Writer
}

type RequestEntry struct {
	Timestamp time.Time `json:"timestamp"`
	ToolName  string    `json:"tool_name"`
	SessionID string    `json:"session_id,omitempty"`
	RequestID string    `json:"request_id,omitempty"`
	Result    any       `json:"result,omitempty"`
	Error     string    `json:"error,omitempty"`
}

func NewRequestStore(dir string) (*RequestStore, error) {
	filePath := filepath.Join(dir, logFileName)
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	return &RequestStore{
		dir:  dir,
		file: f,
		buf:  bufio.NewWriter(f),
	}, nil
}

func (rs *RequestStore) Close() error {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	if err := rs.buf.Flush(); err != nil {
		return err
	}
	return rs.file.Close()
}

func (rs *RequestStore) save(entry RequestEntry) error {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("marshal entry: %w", err)
	}

	rs.mu.Lock()
	defer rs.mu.Unlock()

	if _, err := rs.buf.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write log entry: %w", err)
	}

	if err := rs.buf.Flush(); err != nil {
		return fmt.Errorf("failed to flush buffer: %w", err)
	}

	return nil
}

func (rs *RequestStore) SaveResult(toolName, sessionID string, result any) error {
	return rs.save(RequestEntry{
		Timestamp: time.Now(),
		ToolName:  toolName,
		SessionID: sessionID,
		Result:    result,
	})
}

func (rs *RequestStore) SaveError(toolName, sessionID string, err error) error {
	return rs.save(RequestEntry{
		Timestamp: time.Now(),
		ToolName:  toolName,
		SessionID: sessionID,
		Error:     err.Error(),
	})
}
