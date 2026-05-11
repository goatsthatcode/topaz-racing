package topazracing

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// Task 5.3: Event annotations

func TestRaceVizBootstrapImplementsEventAnnotations(t *testing.T) {
	data, err := os.ReadFile("assets/js/race-viz.js")
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)

	expectedSnippets := []string{
		// Constants
		`DEFAULT_EVENTS_SOURCE_ID`,
		`DEFAULT_EVENTS_LAYER_ID`,

		// Config section
		`events: {`,
		`sourceID: root.dataset.raceVizEventsSource ?? DEFAULT_EVENTS_SOURCE_ID`,
		`layerID: root.dataset.raceVizEventsLayer ?? DEFAULT_EVENTS_LAYER_ID`,

		// State section
		`events: {`,
		`activePopup: null`,

		// State setter
		`setEventsState`,
		`raceVizEventsState`,

		// Feature builder
		`buildEventFeatures`,
		`featureType: "event-annotation"`,

		// Time-only events are resolved via isSelf boat interpolation
		`selfBoat`,
		`interpolateBoatPosition(selfBoat`,

		// Source and layer management
		`upsertEventsSource`,
		`renderEventLayers`,

		// Interaction handler
		`attachEventInteractions`,
		`maplibregl.Popup`,
		`race-viz-event-popup`,

		// Human-readable time formatting
		`formatEventTime`,

		// Loader wiring
		`loadEvents`,
		`boatsReadyPromise`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, source, snippet)
	}
}

func TestRaceVizBootstrapRendersEventLayersWithCorrectStructure(t *testing.T) {
	data, err := os.ReadFile("assets/js/race-viz.js")
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)

	expectedSnippets := []string{
		// Halo circle behind the marker
		"`${layerID}-halo`",

		// Amber/gold color for event markers
		`rgba(255, 200, 60`,

		// Text label layer for events
		"`${layerID}-label`",

		// Mouse cursor change on hover
		`mouseenter`,
		`mouseleave`,
		`style.cursor = "pointer"`,
		`style.cursor = ""`,

		// Popup HTML includes event fields
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
	data, err := os.ReadFile("assets/js/race-viz.js")
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)

	// loadBoats result stored as boatsReadyPromise and passed to loadEvents
	expectedSnippets := []string{
		`const boatsReadyPromise = loadBoats(root, stage, state, mapReadyPromise)`,
		`void loadEvents(root, stage, state, mapReadyPromise, boatsReadyPromise)`,
		// Events waits for boats before building features so timeline is available
		`await boatsReadyPromise?.catch`,
		// Then uses state.replay.timeline for interpolation
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
	cmd.Dir = "."

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("hugo build failed: %v\n%s", err, output)
	}

	racePage := readBuiltHTML(t, filepath.Join(outputDir, "races", "dan-byrne-2025", "bishop-rock-race", "index.html"))

	expectedSnippets := []string{
		// Events state attribute rendered by shortcode
		`data-race-viz-events-state="idle"`,
		// Events URL wired to the events.json resource
		`data-events-url=`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, racePage, snippet)
	}
}
