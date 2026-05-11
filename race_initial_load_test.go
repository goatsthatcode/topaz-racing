package topazracing

import (
	"os"
	"testing"
)

// Task 4.4: Implement initial load behavior

func TestRaceVizBootstrapImplementsInitialLoadBehavior(t *testing.T) {
	data, err := os.ReadFile("assets/js/race-viz.js")
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)

	expectedSnippets := []string{
		// Pre-play / playing mode transitions
		`replay.started`,
		`enterPrePlayMode`,
		`enterPlayingMode`,

		// Pre-play: show only isSelf track, hide competitors
		`selfOnlyFilter`,
		`setFilter`,
		`setLayerVisibility`,

		// Entering playing mode hides static tracks, shows animated layers
		`setLayerVisibility(map, tracksLayerID, false)`,
		`setLayerVisibility(map, replayTailsLayerID, true)`,
		`setLayerVisibility(map, boatMarkersLayerID, true)`,

		// Reset returns to pre-play state
		`state.replay.started = false`,
		`enterPrePlayMode(state.map.instance, state)`,

		// Play and scrub both transition to playing mode
		`enterPlayingMode(state.map.instance, state)`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, source, snippet)
	}
}

func TestRaceVizBootstrapInitializesReplayAtTimeZero(t *testing.T) {
	data, err := os.ReadFile("assets/js/race-viz.js")
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)

	expectedSnippets := []string{
		// Replay starts at the fleet's earliest time
		`state.replay.currentTimeMs = timeline.startTimeMs`,
		`buildReplaySnapshot(timeline, timeline.startTimeMs)`,

		// started flag begins false
		`started: false`,

		// Pre-play mode is entered after boats are loaded
		`enterPrePlayMode(map, state)`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, source, snippet)
	}
}
