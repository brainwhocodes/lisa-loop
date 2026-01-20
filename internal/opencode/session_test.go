package opencode

import (
	"os"
	"testing"
)

func TestSessionPersistence(t *testing.T) {
	// Clean up before and after test
	os.Remove(sessionFile)
	defer os.Remove(sessionFile)

	// Initially no session should exist
	if SessionExists() {
		t.Error("expected no session to exist initially")
	}

	// Load should return empty string
	id, err := LoadSessionID()
	if err != nil {
		t.Errorf("LoadSessionID failed: %v", err)
	}
	if id != "" {
		t.Errorf("expected empty session ID, got %s", id)
	}

	// Save a session ID
	testID := "test-session-abc123"
	if err := SaveSessionID(testID); err != nil {
		t.Fatalf("SaveSessionID failed: %v", err)
	}

	// Session should now exist
	if !SessionExists() {
		t.Error("expected session to exist after save")
	}

	// Load should return the saved ID
	loadedID, err := LoadSessionID()
	if err != nil {
		t.Errorf("LoadSessionID failed: %v", err)
	}
	if loadedID != testID {
		t.Errorf("expected session ID %s, got %s", testID, loadedID)
	}

	// Clear the session
	if err := ClearSession(); err != nil {
		t.Errorf("ClearSession failed: %v", err)
	}

	// Session should no longer exist
	if SessionExists() {
		t.Error("expected no session after clear")
	}
}

func TestClearSession_NoFile(t *testing.T) {
	os.Remove(sessionFile)

	// Should not error if file doesn't exist
	if err := ClearSession(); err != nil {
		t.Errorf("ClearSession should not error for non-existent file: %v", err)
	}
}

func TestSessionAgeHours_NoSession(t *testing.T) {
	os.Remove(sessionFile)

	age, err := SessionAgeHours()
	if err != nil {
		t.Errorf("SessionAgeHours failed: %v", err)
	}
	if age != 0 {
		t.Errorf("expected age 0 for no session, got %d", age)
	}
}

func TestIsSessionExpired_NoExpiry(t *testing.T) {
	os.Remove(sessionFile)
	defer os.Remove(sessionFile)

	if err := SaveSessionID("test-session"); err != nil {
		t.Fatalf("SaveSessionID failed: %v", err)
	}

	// Zero expiry means never expire
	if IsSessionExpired(0) {
		t.Error("expected session not to be expired with 0 expiry")
	}

	// Negative expiry means never expire
	if IsSessionExpired(-1) {
		t.Error("expected session not to be expired with negative expiry")
	}
}

func TestIsSessionExpired_NotExpired(t *testing.T) {
	os.Remove(sessionFile)
	defer os.Remove(sessionFile)

	if err := SaveSessionID("test-session"); err != nil {
		t.Fatalf("SaveSessionID failed: %v", err)
	}

	// Session was just created, should not be expired in 24 hours
	if IsSessionExpired(24) {
		t.Error("expected new session not to be expired")
	}
}
