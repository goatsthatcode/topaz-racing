package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestRaceVizBootstrapImplementsInitialLoadBehavior(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "js", "race-viz.js"))
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)

	expectedSnippets := []string{
		`replay.started`,
		`enterPrePlayMode`,
		`enterPlayingMode`,
		`selfOnlyFilter`,
		`setFilter`,
		`setLayerVisibility`,
		`setLayerVisibility(map, tracksLayerID, false)`,
		`setLayerVisibility(map, replayTailsLayerID, true)`,
		`setLayerVisibility(map, boatMarkersLayerID, true)`,
		`state.replay.started = false`,
		`enterPrePlayMode(state.map.instance, state)`,
		`enterPlayingMode(state.map.instance, state)`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, source, snippet)
	}
}

func TestRaceVizBootstrapInitializesReplayAtTimeZero(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "js", "race-viz.js"))
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)

	expectedSnippets := []string{
		`state.replay.currentTimeMs = timeline.startTimeMs`,
		`buildReplaySnapshot(timeline, timeline.startTimeMs)`,
		`started: false`,
		`enterPrePlayMode(map, state)`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, source, snippet)
	}
}

// TestRaceVizReplaySpeedReadsFromConfig verifies that the initial replay speed
// is read from the config (which reads the shortcode attribute) rather than
// hardcoded to 1. A hardcoded 1 would make a 26-hour race replay in real time.
func TestRaceVizReplaySpeedReadsFromConfig(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "js", "race-viz.js"))
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)

	expectedSnippets := []string{
		// Config reads speed from the DOM attribute
		`raceVizReplaySpeed`,
		`replaySpeed`,
		// State initializes speed from config, not a literal 1
		`speed: config.replaySpeed`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, source, snippet)
	}
}

// TestRaceVizShortcodeEmitsReplaySpeed60 verifies that the shortcode sets the
// default replay speed to 60 on the root element so JS picks it up correctly.
func TestRaceVizShortcodeEmitsReplaySpeed60(t *testing.T) {
	outputDir := t.TempDir()

	cmd := exec.Command("hugo", "--destination", outputDir)
	cmd.Env = os.Environ()
	cmd.Dir = repoRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("hugo build failed: %v\n%s", err, output)
	}

	racePage := readBuiltHTML(t, filepath.Join(outputDir, "races", "dan-byrne-2025", "bishop-rock-race", "index.html"))
	assertContains(t, racePage, `data-race-viz-replay-speed="60"`)
}

// TestRaceVizBoatsFallbackIsImplemented verifies that a boats-load failure
// renders a visible error message rather than silently disabling controls.
func TestRaceVizBoatsFallbackIsImplemented(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "js", "race-viz.js"))
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)

	expectedSnippets := []string{
		`renderBoatsFallback`,
		`race-viz-boats-fallback`,
		`raceVizBoatsFallback`,
		// setBoatsState forwards message to renderBoatsFallback
		`setBoatsState(root, stage, state, "error"`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, source, snippet)
	}
}
