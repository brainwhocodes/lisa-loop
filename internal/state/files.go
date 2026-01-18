package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ReadStateFile reads a JSON state file
func ReadStateFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file %s: %w", path, err)
	}
	return data, nil
}

// WriteStateFile writes data to a file atomically (write to temp, then rename)
func WriteStateFile(path string, data []byte) error {
	// Write to temporary file
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write temp file %s: %w", tmpPath, err)
	}

	// Atomic rename
	if err := os.Rename(tmpPath, path); err != nil {
		// Clean up temp file if rename fails
		os.Remove(tmpPath)
		return fmt.Errorf("failed to rename %s to %s: %w", tmpPath, path, err)
	}

	return nil
}

// AtomicWrite is a generic atomic write function
func AtomicWrite(path string, data []byte) error {
	return WriteStateFile(path, data)
}

// LoadCallCount loads the call count from .call_count
func LoadCallCount() (int, error) {
	data, err := ReadStateFile(".call_count")
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}

	var count int
	if err := json.Unmarshal(data, &count); err != nil {
		return 0, fmt.Errorf("failed to parse call count: %w", err)
	}

	return count, nil
}

// SaveCallCount saves the call count to .call_count atomically
func SaveCallCount(count int) error {
	data, err := json.Marshal(count)
	if err != nil {
		return fmt.Errorf("failed to marshal call count: %w", err)
	}

	return WriteStateFile(".call_count", data)
}

// LoadLastReset loads last reset time from .last_reset
func LoadLastReset() (time.Time, error) {
	data, err := ReadStateFile(".last_reset")
	if err != nil {
		if os.IsNotExist(err) {
			return time.Now(), nil
		}
		return time.Time{}, err
	}

	var t time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return time.Time{}, fmt.Errorf("failed to parse last reset: %w", err)
	}

	return t, nil
}

// SaveLastReset saves last reset time to .last_reset atomically
func SaveLastReset(t time.Time) error {
	data, err := json.Marshal(t)
	if err != nil {
		return fmt.Errorf("failed to marshal last reset: %w", err)
	}

	return WriteStateFile(".last_reset", data)
}

// LoadCodexSession loads Codex session ID from .codex_session_id
func LoadCodexSession() (string, error) {
	data, err := os.ReadFile(".codex_session_id")
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	return string(data), nil
}

// SaveCodexSession saves Codex session ID to .codex_session_id atomically
func SaveCodexSession(id string) error {
	return AtomicWrite(".codex_session_id", []byte(id))
}

// LoadRalphSession loads Ralph session metadata from .ralph_session
func LoadRalphSession() (map[string]interface{}, error) {
	data, err := ReadStateFile(".ralph_session")
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]interface{}{}, nil
		}
		return nil, err
	}

	var session map[string]interface{}
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to parse Ralph session: %w", err)
	}

	return session, nil
}

// SaveRalphSession saves Ralph session metadata atomically
func SaveRalphSession(session map[string]interface{}) error {
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal Ralph session: %w", err)
	}

	return WriteStateFile(".ralph_session", data)
}

// LoadExitSignals loads recent exit signals from .exit_signals
func LoadExitSignals() ([]string, error) {
	data, err := ReadStateFile(".exit_signals")
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var signals []string
	if err := json.Unmarshal(data, &signals); err != nil {
		return nil, fmt.Errorf("failed to parse exit signals: %w", err)
	}

	return signals, nil
}

// SaveExitSignals saves exit signals atomically
func SaveExitSignals(signals []string) error {
	data, err := json.Marshal(signals)
	if err != nil {
		return fmt.Errorf("failed to marshal exit signals: %w", err)
	}

	return WriteStateFile(".exit_signals", data)
}

// LoadCircuitBreakerState loads circuit breaker state from .circuit_breaker_state
func LoadCircuitBreakerState() (map[string]interface{}, error) {
	data, err := ReadStateFile(".circuit_breaker_state")
	if err != nil {
		if os.IsNotExist(err) {
			// Default to CLOSED state
			return map[string]interface{}{
				"state":           "CLOSED",
				"last_check_time": time.Now().Format(time.RFC3339),
			}, nil
		}
		return nil, err
	}

	var state map[string]interface{}
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse circuit breaker state: %w", err)
	}

	return state, nil
}

// SaveCircuitBreakerState saves circuit breaker state atomically
func SaveCircuitBreakerState(state map[string]interface{}) error {
	data, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal circuit breaker state: %w", err)
	}

	return WriteStateFile(".circuit_breaker_state", data)
}

// EnsureStateDir ensures the directory for state files exists
func EnsureStateDir() error {
	if err := os.MkdirAll(".", 0755); err != nil {
		return fmt.Errorf("failed to create state directory: %w", err)
	}
	return nil
}

// CleanupOldFiles removes old temporary state files
func CleanupOldFiles() error {
	entries, err := os.ReadDir(".")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".tmp" {
			os.Remove(filepath.Join(".", entry.Name()))
		}
	}

	return nil
}
