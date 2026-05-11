package topazracing

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// Task 5.1: Boat legend and visibility toggles

func TestRaceVizBootstrapImplementsBoatLegendToggles(t *testing.T) {
	data, err := os.ReadFile("assets/js/race-viz.js")
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)

	expectedSnippets := []string{
		// Per-boat visibility state in createRaceVizState
		`visibility`,
		`hiddenBoatIds: new Set()`,

		// Toggle button element on each legend item
		`race-viz-boat-toggle`,
		`data-race-viz-boat-toggle`,
		`aria-pressed`,

		// Hidden state tracked via dataset property on legend items
		`raceVizBoatHidden`,

		// Core toggle functions
		`syncBoatLegendVisibility`,
		`attachBoatLegendToggles`,
		`applyBoatVisibilityToLayers`,

		// Toggle event wiring: clicking a toggle updates visibility then applies to map
		`state.visibility.hiddenBoatIds.has`,
		`state.visibility.hiddenBoatIds.delete`,
		`state.visibility.hiddenBoatIds.add`,
		`applyBoatVisibilityToLayers(state.map.instance, state)`,

		// attachBoatLegendToggles is called after renderBoatLegend in loadBoats
		`attachBoatLegendToggles(root, state)`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, source, snippet)
	}
}

func TestRaceVizBootstrapFiltersHiddenBoatsFromReplayLayers(t *testing.T) {
	data, err := os.ReadFile("assets/js/race-viz.js")
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)

	expectedSnippets := []string{
		// renderReplayFrame reads visibility state and passes it down
		`state.visibility.hiddenBoatIds`,
		`buildReplayTailFeatures(state.replay.timeline, state.replay.currentTimeMs, hiddenBoatIds)`,
		`buildBoatMarkerFeatures(state.replay.snapshot, hiddenBoatIds)`,

		// buildReplayTailFeatures accepts hiddenBoatIds and skips hidden boats
		`hiddenBoatIds = null`,
		`hiddenBoatIds !== null && hiddenBoatIds.has`,

		// buildBoatMarkerFeatures accepts hiddenBoatIds and filters before mapping
		`hiddenBoatIds === null || !hiddenBoatIds.has`,

		// Pre-play mode: visibility applied as a map filter on the tracks layer
		`applyBoatVisibilityToLayers`,
		`!state.replay.started`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, source, snippet)
	}
}

func TestRaceVizFleetPanelAppearsInBuiltRacePage(t *testing.T) {
	outputDir := t.TempDir()

	cmd := exec.Command("hugo", "--destination", outputDir)
	cmd.Env = os.Environ()
	cmd.Dir = "."

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("hugo build failed: %v\n%s", err, output)
	}

	racePage := readBuiltHTML(t, filepath.Join(outputDir, "races", "dan-byrne-2025", "bishop-rock-race", "index.html"))

	expectedSnippets := []string{
		`class="race-viz-panel race-viz-panel-fleet"`,
		`class="race-viz-sidebar-title"`,
		`class="race-viz-boat-legend" data-race-viz-boat-legend`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, racePage, snippet)
	}
}
