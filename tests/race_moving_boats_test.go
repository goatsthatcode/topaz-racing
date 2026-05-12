package tests

import (
	"os"
	"testing"
)

func TestRaceVizBootstrapImplementsMovingBoatsAndTrackTails(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "js", "race-viz.js"))
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)

	expectedSnippets := []string{
		`DEFAULT_REPLAY_TAILS_SOURCE_ID`,
		`DEFAULT_REPLAY_TAILS_LAYER_ID`,
		`DEFAULT_BOAT_MARKERS_SOURCE_ID`,
		`DEFAULT_BOAT_MARKERS_LAYER_ID`,
		`buildTrackTailCoordinates`,
		`buildReplayTailFeatures`,
		`featureType: "replay-tail"`,
		`buildBoatMarkerFeatures`,
		`featureType: "boat-marker"`,
		`upsertReplayTailsSource`,
		`upsertBoatMarkersSource`,
		`renderReplayTailLayers`,
		`renderBoatMarkerLayers`,
		`renderReplayFrame`,
		`state.replay.started && state.map.instance`,
		`featureType: "boat-marker"`,
		`circle-color`,
		`circle-radius`,
		`circle-stroke-width`,
		`featureType: "replay-tail"`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, source, snippet)
	}
}

func TestRaceVizBootstrapUpdatesMapOnEachReplayTick(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "js", "race-viz.js"))
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)

	expectedSnippets := []string{
		`renderReplayFrame(state.map.instance, state)`,
		`buildReplayTailFeatures(state.replay.timeline, state.replay.currentTimeMs, hiddenBoatIds)`,
		`upsertReplayTailsSource(map, state, tailFeatures)`,
		`buildBoatMarkerFeatures(state.replay.snapshot, hiddenBoatIds)`,
		`upsertBoatMarkersSource(map, state, markerFeatures)`,
		`emptyFeatureCollection`,
		`interpolateBoatPosition(boat, timeMs)`,
		`buildTrackTailCoordinates`,
		`visibility: "none"`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, source, snippet)
	}
}
