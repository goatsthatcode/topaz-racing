package tests

import (
	"os"
	"testing"
)

// TestRaceVizElapsedLabelFunctionExists verifies that formatElapsedLabel exists and
// has replaced the old UTC-based formatReplayClockLabel (UX-1).
func TestRaceVizElapsedLabelFunctionExists(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "js", "race-viz.js"))
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)

	assertContains(t, source, `function formatElapsedLabel(elapsedMs)`)
	assertNotContains(t, source, `function formatReplayClockLabel`)
}

// TestRaceVizElapsedLabelFormatsHoursMinutesSeconds verifies the elapsed formatter
// produces a "+HH:MM:SS" string using padded components (UX-1).
func TestRaceVizElapsedLabelFormatsHoursMinutesSeconds(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "js", "race-viz.js"))
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)

	// Elapsed formatter must compute hours/minutes/seconds from total ms.
	expectedSnippets := []string{
		`Math.floor(ms / 1000)`,
		`Math.floor(totalSeconds / 3600)`,
		`Math.floor((totalSeconds % 3600) / 60)`,
		`totalSeconds % 60`,
		// Result is a "+HH:MM:SS" string.
		"String(hours).padStart(2",
		"String(minutes).padStart(2",
		"String(seconds).padStart(2",
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, source, snippet)
	}
}

// TestRaceVizSyncReplayControlsUsesElapsedLabel verifies that syncReplayControls
// feeds elapsed milliseconds (not absolute timestamps) to formatElapsedLabel (UX-1).
func TestRaceVizSyncReplayControlsUsesElapsedLabel(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "js", "race-viz.js"))
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)

	// Current label is driven by elapsed time from race start.
	assertContains(t, source, `(state.replay.currentTimeMs ?? 0) - (state.replay.startTimeMs ?? 0)`)
	// End label uses durationMs, not endTimeMs.
	assertContains(t, source, `formatElapsedLabel(state.replay.durationMs ?? 0)`)
	// Start label is always zero elapsed.
	assertContains(t, source, `formatElapsedLabel(0)`)
}

// TestRaceVizTrackHoverUsesElapsedLabel verifies that the track hover tooltip
// expresses position time as elapsed from race start, not UTC wall-clock (UX-1).
func TestRaceVizTrackHoverUsesElapsedLabel(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "js", "race-viz.js"))
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)

	assertContains(t, source, `formatElapsedLabel(timeMs - (state.replay.startTimeMs ?? 0))`)
}
