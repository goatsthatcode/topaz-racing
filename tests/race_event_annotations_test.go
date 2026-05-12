package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestRaceVizBootstrapImplementsEventAnnotations(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "js", "race-viz.js"))
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)

	expectedSnippets := []string{
		`DEFAULT_EVENTS_SOURCE_ID`,
		`DEFAULT_EVENTS_LAYER_ID`,
		`events: {`,
		`sourceID: root.dataset.raceVizEventsSource ?? DEFAULT_EVENTS_SOURCE_ID`,
		`layerID: root.dataset.raceVizEventsLayer ?? DEFAULT_EVENTS_LAYER_ID`,
		`events: {`,
		`activePopup: null`,
		`setEventsState`,
		`raceVizEventsState`,
		`buildEventFeatures`,
		`featureType: "event-annotation"`,
		`selfBoat`,
		`interpolateBoatPosition(selfBoat`,
		`upsertEventsSource`,
		`renderEventLayers`,
		`attachEventInteractions`,
		`maplibregl.Popup`,
		`race-viz-event-popup`,
		`formatEventTime`,
		`loadEvents`,
		`boatsReadyPromise`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, source, snippet)
	}
}

func TestRaceVizBootstrapRendersEventLayersWithCorrectStructure(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "js", "race-viz.js"))
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)

	expectedSnippets := []string{
		"`${layerID}-halo`",
		`rgba(255, 200, 60`,
		"`${layerID}-label`",
		`mouseenter`,
		`mouseleave`,
		`style.cursor = "pointer"`,
		`style.cursor = ""`,
		`race-viz-event-popup-label`,
		`race-viz-event-popup-description`,
		`race-viz-event-popup-time`,
		`race-viz-event-type`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, source, snippet)
	}
}

func TestRaceVizBootstrapWiresEventsAfterBoats(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "js", "race-viz.js"))
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)

	expectedSnippets := []string{
		`const boatsReadyPromise = loadBoats(root, stage, state, mapReadyPromise)`,
		`void loadEvents(root, stage, state, mapReadyPromise, boatsReadyPromise)`,
		`await boatsReadyPromise?.catch`,
		`buildEventFeatures(payload, state.replay.timeline)`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, source, snippet)
	}
}

func TestRaceVizShortcodeEmitsEventStateAttribute(t *testing.T) {
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
		`data-race-viz-events-state="idle"`,
		`data-events-url=`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, racePage, snippet)
	}
}
