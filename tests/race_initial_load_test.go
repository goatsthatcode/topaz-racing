package tests

import (
	"os"
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
