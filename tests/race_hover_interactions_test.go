package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestRaceVizBootstrapImplementsHoverState(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "js", "race-viz.js"))
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)

	expectedSnippets := []string{
		`hover: {`,
		`activeTooltip: null`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, source, snippet)
	}
}

func TestRaceVizBootstrapImplementsBoatMarkerHoverInteractions(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "js", "race-viz.js"))
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)

	expectedSnippets := []string{
		`attachBoatMarkerHoverInteractions`,
		`mouseenter`,
		`mouseleave`,
		`style.cursor = "pointer"`,
		`race-viz-hover-name`,
		`race-viz-hover-time`,
		`formatReplayClockLabel(state.replay.currentTimeMs)`,
		`race-viz-hover-tooltip`,
		`window.maplibregl.Popup`,
		`state.hover.activeTooltip`,
		`state.hover.activeTooltip.remove()`,
		`state.hover.activeTooltip = null`,
		`attachBoatMarkerHoverInteractions(map, state)`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, source, snippet)
	}
}

func TestRaceVizBootstrapImplementsTrackHoverInteractions(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "js", "race-viz.js"))
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)

	expectedSnippets := []string{
		`attachTrackHoverInteractions`,
		`attachTrackHoverInteractions(map, state)`,
		`style.cursor = "crosshair"`,
		`state.config.tracks.layerID`,
		`state.config.replayTails.layerID`,
		`race-viz-hover-tooltip-content`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, source, snippet)
	}
}

func TestRaceVizHoverTooltipCSSIsDefined(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "css", "race-viz.css"))
	if err != nil {
		t.Fatalf("failed to read race viz CSS: %v", err)
	}

	source := string(data)

	expectedSnippets := []string{
		`.race-viz-hover-tooltip .maplibregl-popup-content`,
		`.race-viz-hover-tooltip-content`,
		`.race-viz-hover-name`,
		`.race-viz-hover-time`,
		`pointer-events: none`,
		`font-variant-numeric: tabular-nums`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, source, snippet)
	}
}

func TestRaceVizHoverInteractionsWiredInLoadBoats(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "js", "race-viz.js"))
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)

	expectedSnippets := []string{
		`renderBoatMarkerLayers(map, state)`,
		`attachBoatMarkerHoverInteractions(map, state)`,
		`attachTrackHoverInteractions(map, state)`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, source, snippet)
	}
}

func TestRaceVizBuiltRacePageRendersWithHoverReadyLayers(t *testing.T) {
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
		`data-race-viz-boats-state="idle"`,
		`data-race-viz-replay-state="idle"`,
		`data-boats-url=`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, racePage, snippet)
	}
}
