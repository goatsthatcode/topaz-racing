package topazracing

import (
	"os"
	"testing"
)

// Task 4.5: Render moving boats and persistent full track tails

func TestRaceVizBootstrapImplementsMovingBoatsAndTrackTails(t *testing.T) {
	data, err := os.ReadFile("assets/js/race-viz.js")
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)

	expectedSnippets := []string{
		// Replay tail source and layer IDs
		`DEFAULT_REPLAY_TAILS_SOURCE_ID`,
		`DEFAULT_REPLAY_TAILS_LAYER_ID`,

		// Boat marker source and layer IDs
		`DEFAULT_BOAT_MARKERS_SOURCE_ID`,
		`DEFAULT_BOAT_MARKERS_LAYER_ID`,

		// Track tail computation
		`buildTrackTailCoordinates`,
		`buildReplayTailFeatures`,
		`featureType: "replay-tail"`,

		// Boat marker computation
		`buildBoatMarkerFeatures`,
		`featureType: "boat-marker"`,

		// Source upsert helpers
		`upsertReplayTailsSource`,
		`upsertBoatMarkersSource`,

		// Layer setup helpers
		`renderReplayTailLayers`,
		`renderBoatMarkerLayers`,

		// Frame render called from setReplayTime during playback
		`renderReplayFrame`,
		`state.replay.started && state.map.instance`,

		// Boat markers use per-boat color and isSelf discriminant
		`featureType: "boat-marker"`,
		`circle-color`,
		`circle-radius`,
		`circle-stroke-width`,

		// Replay tails use per-boat color and isSelf discriminant
		`featureType: "replay-tail"`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, source, snippet)
	}
}

func TestRaceVizBootstrapUpdatesMapOnEachReplayTick(t *testing.T) {
	data, err := os.ReadFile("assets/js/race-viz.js")
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)

	expectedSnippets := []string{
		// renderReplayFrame is called from setReplayTime
		`renderReplayFrame(state.map.instance, state)`,

		// Frame updates both tails and markers
		`buildReplayTailFeatures(state.replay.timeline, state.replay.currentTimeMs)`,
		`upsertReplayTailsSource(map, state, tailFeatures)`,
		`buildBoatMarkerFeatures(state.replay.snapshot)`,
		`upsertBoatMarkersSource(map, state, markerFeatures)`,

		// Empty feature collections are used for initial source state
		`emptyFeatureCollection`,

		// Tail coordinates are trimmed to current time then extended with interpolated tip
		`interpolateBoatPosition(boat, timeMs)`,
		`buildTrackTailCoordinates`,

		// Replay layers are initialized with visibility none, shown on entering playing mode
		`visibility: "none"`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, source, snippet)
	}
}
