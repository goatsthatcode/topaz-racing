package topazracing

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// Task 5.2: Hover interactions

func TestRaceVizBootstrapImplementsHoverState(t *testing.T) {
	data, err := os.ReadFile("assets/js/race-viz.js")
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)

	expectedSnippets := []string{
		// hover slot in createRaceVizState
		`hover: {`,
		`activeTooltip: null`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, source, snippet)
	}
}

func TestRaceVizBootstrapImplementsBoatMarkerHoverInteractions(t *testing.T) {
	data, err := os.ReadFile("assets/js/race-viz.js")
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)

	expectedSnippets := []string{
		// Function exists
		`attachBoatMarkerHoverInteractions`,

		// Cursor change on hover over boat markers
		`mouseenter`,
		`mouseleave`,
		`style.cursor = "pointer"`,

		// Popup created with boat name and current time
		`race-viz-hover-name`,
		`race-viz-hover-time`,
		`formatReplayClockLabel(state.replay.currentTimeMs)`,
		`race-viz-hover-tooltip`,

		// maplibregl.Popup used (same as event annotations)
		`window.maplibregl.Popup`,

		// Tooltip removed on mouse leave and tracked in state
		`state.hover.activeTooltip`,
		`state.hover.activeTooltip.remove()`,
		`state.hover.activeTooltip = null`,

		// Wired after renderBoatMarkerLayers in loadBoats
		`attachBoatMarkerHoverInteractions(map, state)`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, source, snippet)
	}
}

func TestRaceVizBootstrapImplementsTrackHoverInteractions(t *testing.T) {
	data, err := os.ReadFile("assets/js/race-viz.js")
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)

	expectedSnippets := []string{
		// Function exists and is wired
		`attachTrackHoverInteractions`,
		`attachTrackHoverInteractions(map, state)`,

		// Cursor change on hover over track/tail lines
		`style.cursor = "crosshair"`,

		// Covers both static tracks and replay tails layers
		`state.config.tracks.layerID`,
		`state.config.replayTails.layerID`,

		// Shows boat name in tooltip
		`race-viz-hover-tooltip-content`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, source, snippet)
	}
}

func TestRaceVizHoverTooltipCSSIsDefined(t *testing.T) {
	data, err := os.ReadFile("assets/css/race-viz.css")
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
	data, err := os.ReadFile("assets/js/race-viz.js")
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)

	// Both hover attachment calls appear together after renderBoatMarkerLayers
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
	cmd.Dir = "."

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("hugo build failed: %v\n%s", err, output)
	}

	racePage := readBuiltHTML(t, filepath.Join(outputDir, "races", "dan-byrne-2025", "bishop-rock-race", "index.html"))

	// Boat markers and track layers whose IDs are used by hover interactions must be in data attributes
	expectedSnippets := []string{
		`data-race-viz-boats-state="idle"`,
		`data-race-viz-replay-state="idle"`,
		`data-boats-url=`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, racePage, snippet)
	}
}
