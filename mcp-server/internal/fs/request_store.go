package fs

import (
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
	dir string
	mu  sync.Mutex
}

type RequestEntry struct {
	Timestamp time.Time `json:"timestamp"`
	ToolName  string    `json:"tool_name"`
	SessionID string    `json:"session_id,omitempty"`
	RequestID string    `json:"request_id,omitempty"`
	Result    any       `json:"result,omitempty"`
	Error     string    `json:"error,omitempty"`
}

func NewRequestStore(dir string) *RequestStore {
	return &RequestStore{dir: dir}
}

func (rs *RequestStore) save(entry RequestEntry) error {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("marshal entry: %w", err)
	}

	filePath := filepath.Join(rs.dir, logFileName)

	rs.mu.Lock()
	defer rs.mu.Unlock()

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer f.Close()

	if _, err := f.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write log entry: %w", err)
	}

	return nil
}

// Helpers
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
