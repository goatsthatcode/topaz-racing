package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestRaceVizBootstrapImplementsBoatLegendToggles(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "js", "race-viz.js"))
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)

	expectedSnippets := []string{
		`visibility`,
		`hiddenBoatIds: new Set()`,
		`race-viz-boat-toggle`,
		`data-race-viz-boat-toggle`,
		`aria-pressed`,
		`raceVizBoatHidden`,
		`syncBoatLegendVisibility`,
		`attachBoatLegendToggles`,
		`applyBoatVisibilityToLayers`,
		`state.visibility.hiddenBoatIds.has`,
		`state.visibility.hiddenBoatIds.delete`,
		`state.visibility.hiddenBoatIds.add`,
		`applyBoatVisibilityToLayers(state.map.instance, state)`,
		`attachBoatLegendToggles(root, state)`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, source, snippet)
	}
}

func TestRaceVizBootstrapFiltersHiddenBoatsFromReplayLayers(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "js", "race-viz.js"))
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)

	expectedSnippets := []string{
		`state.visibility.hiddenBoatIds`,
		`buildReplayTailFeatures(state.replay.timeline, state.replay.currentTimeMs, hiddenBoatIds)`,
		`buildBoatMarkerFeatures(state.replay.snapshot, hiddenBoatIds)`,
		`hiddenBoatIds = null`,
		`hiddenBoatIds !== null && hiddenBoatIds.has`,
		`hiddenBoatIds === null || !hiddenBoatIds.has`,
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
	cmd.Dir = repoRoot

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
