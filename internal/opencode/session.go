package opencode

import (
	"os"
	"time"
)

const sessionFile = ".opencode_session_id"

// LoadSessionID loads the OpenCode session ID from .opencode_session_id
func LoadSessionID() (string, error) {
	data, err := os.ReadFile(sessionFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return string(data), nil
}

// SaveSessionID saves the OpenCode session ID to .opencode_session_id atomically
func SaveSessionID(id string) error {
	tmpPath := sessionFile + ".tmp"
	if err := os.WriteFile(tmpPath, []byte(id), 0644); err != nil {
		return err
	}
	return os.Rename(tmpPath, sessionFile)
}

// ClearSession removes the session file
func ClearSession() error {
	err := os.Remove(sessionFile)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// SessionExists checks if a session ID file exists
func SessionExists() bool {
	id, err := LoadSessionID()
	if err != nil {
		return false
	}
	return id != ""
}

// SessionAgeHours calculates the age of the session in hours
func SessionAgeHours() (int, error) {
	if !SessionExists() {
		return 0, nil
	}

	info, err := os.Stat(sessionFile)
	if err != nil {
		return 0, err
	}

	age := time.Since(info.ModTime()).Hours()
	return int(age), nil
}

// IsSessionExpired checks if the session has expired
func IsSessionExpired(expiryHours int) bool {
	if expiryHours <= 0 {
		return false
	}

	age, err := SessionAgeHours()
	if err != nil {
		return false
	}

	return age >= expiryHours
}
