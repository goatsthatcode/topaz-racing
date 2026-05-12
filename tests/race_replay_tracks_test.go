package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestRaceVizPublishesReplayTrackContract(t *testing.T) {
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
		`data-race-viz-replay-time="0"`,
		`data-race-viz-replay-speed="1"`,
		`data-race-viz-tracks-source="race-viz-tracks"`,
		`data-race-viz-tracks-layer="race-viz-tracks"`,
		`class="race-viz-sidebar" data-race-viz-sidebar`,
		`class="race-viz-boat-legend" data-race-viz-boat-legend`,
		`class="race-viz-panel race-viz-panel-controls" data-race-viz-controls`,
		`class="race-viz-button" data-race-viz-play-toggle disabled>Play</button>`,
		`class="race-viz-button race-viz-button-secondary" data-race-viz-replay-reset disabled>Reset</button>`,
		`class="race-viz-speed-select" data-race-viz-replay-speed-select disabled`,
		`class="race-viz-timeline-input"`,
		`data-race-viz-replay-timeline`,
		`data-race-viz-replay-current-label`,
		`data-race-viz-replay-start-label`,
		`data-race-viz-replay-end-label`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, racePage, snippet)
	}
}

func TestRaceVizBootstrapImplementsStaticBoatTrackRendering(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "js", "race-viz.js"))
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)
	expectedSnippets := []string{
		`DEFAULT_TRACKS_SOURCE_ID`,
		`DEFAULT_TRACKS_LAYER_ID`,
		`buildBoatTrackFeatures`,
		`buildReplayTimeline`,
		`interpolateBoatPosition`,
		`buildReplaySnapshot`,
		`attachReplayControls`,
		`startReplayPlayback`,
		`stopReplayPlayback`,
		`resetReplay`,
		`setReplayTime`,
		`raceVizReplayState`,
		`raceVizReplayStart`,
		`raceVizReplayEnd`,
		`raceVizReplayDuration`,
		`raceVizReplayTime`,
		`raceVizReplaySpeed`,
		`raceVizReplayPlaying`,
		`featureType: "boat-track"`,
		`pointCount: boat.track.length`,
		`data-race-viz-play-toggle`,
		`data-race-viz-replay-reset`,
		`data-race-viz-replay-speed-select`,
		`data-race-viz-replay-timeline`,
		`requestAnimationFrame`,
		`cancelAnimationFrame`,
		`data-race-viz-boat-legend`,
		`raceVizBoatColor`,
		`raceVizBoatCount`,
		`raceVizSelfBoatId`,
		`["coalesce", ["get", "color"], "#ffffff"]`,
		`["boolean", ["get", "isSelf"], false]`,
		`loadBoats`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, source, snippet)
	}
}

func TestReferenceRaceBoatsProvideRenderableTrackPolylines(t *testing.T) {
	var boats raceBoatsFile
	readJSONFixture(
		t,
		repoFile("content", "races", "dan-byrne-2025", "bishop-rock-race", "boats.json"),
		&boats,
	)

	for _, boat := range boats.Boats {
		if len(boat.Track) < 2 {
			t.Fatalf("expected boat %q to have at least two points for a static track polyline", boat.ID)
		}
		if boat.Track[0] == boat.Track[len(boat.Track)-1] {
			t.Fatalf("expected boat %q track to cover more than a single repeated position", boat.ID)
		}
	}
}

func TestReferenceRaceBoatsSupportInterpolatedReplayMidpoint(t *testing.T) {
	var boats raceBoatsFile
	readJSONFixture(
		t,
		repoFile("content", "races", "dan-byrne-2025", "bishop-rock-race", "boats.json"),
		&boats,
	)

	var replayStart time.Time
	var replayEnd time.Time
	for i, boat := range boats.Boats {
		boatStart, err := time.Parse(time.RFC3339, boat.Track[0].Time)
		if err != nil {
			t.Fatalf("failed to parse boat %q start time: %v", boat.ID, err)
		}
		boatEnd, err := time.Parse(time.RFC3339, boat.Track[len(boat.Track)-1].Time)
		if err != nil {
			t.Fatalf("failed to parse boat %q end time: %v", boat.ID, err)
		}

		if i == 0 || boatStart.Before(replayStart) {
			replayStart = boatStart
		}
		if i == 0 || boatEnd.After(replayEnd) {
			replayEnd = boatEnd
		}
	}

	midpoint := replayStart.Add(replayEnd.Sub(replayStart) / 2)
	for _, boat := range boats.Boats {
		straddlesMidpoint := false
		for i := 1; i < len(boat.Track); i++ {
			segmentStart, err := time.Parse(time.RFC3339, boat.Track[i-1].Time)
			if err != nil {
				t.Fatalf("failed to parse boat %q track point %d: %v", boat.ID, i-1, err)
			}
			segmentEnd, err := time.Parse(time.RFC3339, boat.Track[i].Time)
			if err != nil {
				t.Fatalf("failed to parse boat %q track point %d: %v", boat.ID, i, err)
			}

			if !midpoint.Before(segmentStart) && !midpoint.After(segmentEnd) {
				straddlesMidpoint = true
				break
			}
		}

		if !straddlesMidpoint {
			t.Fatalf("expected boat %q to have a track segment covering replay midpoint %s", boat.ID, midpoint.Format(time.RFC3339))
		}
	}
}
