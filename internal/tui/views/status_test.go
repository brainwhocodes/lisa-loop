package views

import (
	"strings"
	"testing"
)

func TestFormatCircuitState(t *testing.T) {
	tests := []struct {
		name  string
		state string
		want  string
	}{
		{
			name:  "CLOSED state",
			state: "CLOSED",
			want:  "CLOSED",
		},
		{
			name:  "HALF_OPEN state",
			state: "HALF_OPEN",
			want:  "HALF_OPEN",
		},
		{
			name:  "OPEN state",
			state: "OPEN",
			want:  "OPEN",
		},
		{
			name:  "unknown state",
			state: "UNKNOWN",
			want:  "UNKNOWN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatCircuitState(tt.state)
			// Check that output contains state name (ignoring ANSI codes)
			if !strings.Contains(got, tt.want) {
				t.Errorf("FormatCircuitState() output doesn't contain %v, got %v", tt.want, got)
			}
		})
	}
}

func TestFormatWorkType(t *testing.T) {
	tests := []struct {
		name     string
		workType string
		want     string
	}{
		{
			name:     "IMPLEMENTATION work type",
			workType: "IMPLEMENTATION",
			want:     "IMPLEMENTATION",
		},
		{
			name:     "TESTING work type",
			workType: "TESTING",
			want:     "TESTING",
		},
		{
			name:     "DOCUMENTATION work type",
			workType: "DOCUMENTATION",
			want:     "DOCUMENTATION",
		},
		{
			name:     "REFACTORING work type",
			workType: "REFACTORING",
			want:     "REFACTORING",
		},
		{
			name:     "unknown work type",
			workType: "UNKNOWN",
			want:     "UNKNOWN",
		},
		{
			name:     "empty work type",
			workType: "",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatWorkType(tt.workType)
			if tt.want == "" {
				if got != tt.want {
					t.Errorf("FormatWorkType() = %v, want %v", got, tt.want)
				}
			} else if !strings.Contains(got, tt.want) {
				t.Errorf("FormatWorkType() output doesn't contain %v, got %v", tt.want, got)
			}
		})
	}
}

func TestUpdateProgressBar(t *testing.T) {
	tests := []struct {
		name     string
		progress float64
	}{
		{
			name:     "0% progress",
			progress: 0.0,
		},
		{
			name:     "50% progress",
			progress: 0.5,
		},
		{
			name:     "100% progress",
			progress: 1.0,
		},
		{
			name:     "25% progress",
			progress: 0.25,
		},
		{
			name:     "75% progress",
			progress: 0.75,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UpdateProgressBar(tt.progress)
			if got == "" {
				t.Errorf("UpdateProgressBar() returned empty string")
			}
			if len(got) < 2 {
				t.Errorf("UpdateProgressBar() returned string too short")
			}
		})
	}
}

func TestRender(t *testing.T) {
	tests := []struct {
		name           string
		loopNum        int
		callsMade      int
		callsRemaining int
		circuitState   string
		codexStatus    string
	}{
		{
			name:           "basic render",
			loopNum:        5,
			callsMade:      20,
			callsRemaining: 80,
			circuitState:   "CLOSED",
			codexStatus:    "running",
		},
		{
			name:           "zero loop",
			loopNum:        0,
			callsMade:      0,
			callsRemaining: 100,
			circuitState:   "OPEN",
			codexStatus:    "halted",
		},
		{
			name:           "high loop count",
			loopNum:        100,
			callsMade:      95,
			callsRemaining: 5,
			circuitState:   "HALF_OPEN",
			codexStatus:    "monitoring",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Render(tt.loopNum, tt.callsMade, tt.callsRemaining, tt.circuitState, tt.codexStatus, 0)
			if got == "" {
				t.Errorf("Render() returned empty string")
			}
		})
	}
}
