package topazracing

import (
	"path/filepath"
	"testing"
	"time"
)

func TestReferenceRaceBundleProvidesCompleteV1Dataset(t *testing.T) {
	bundleDir := filepath.Join("content", "races", "dan-byrne-2025", "bishop-rock-race")

	var course raceCourseFile
	readJSONFixture(t, filepath.Join(bundleDir, "course.json"), &course)

	var boats raceBoatsFile
	readJSONFixture(t, filepath.Join(bundleDir, "boats.json"), &boats)

	var events raceEventsFile
	readJSONFixture(t, filepath.Join(bundleDir, "events.json"), &events)

	if len(course.Elements) < 3 {
		t.Fatalf("expected sample course to include at least 3 elements, got %d", len(course.Elements))
	}
	if len(boats.Boats) < 2 {
		t.Fatalf("expected sample race to include at least 2 boats, got %d", len(boats.Boats))
	}
	if len(events.Events) < 1 {
		t.Fatalf("expected sample race to include at least 1 event, got %d", len(events.Events))
	}
	manualShapingPoints := 0
	for _, element := range course.Elements {
		manualShapingPoints += len(element.ControlPointsToNext)
	}
	if manualShapingPoints < 1 {
		t.Fatal("expected sample course to exercise the manual land-routing fallback with at least one shaping point")
	}

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

	for _, event := range events.Events {
		if event.Time == "" {
			continue
		}

		eventTime, err := time.Parse(time.RFC3339, event.Time)
		if err != nil {
			t.Fatalf("failed to parse event %q time: %v", event.ID, err)
		}
		if eventTime.Before(replayStart) || eventTime.After(replayEnd) {
			t.Fatalf(
				"expected event %q at %s to fall within replay bounds %s - %s",
				event.ID,
				eventTime.Format(time.RFC3339),
				replayStart.Format(time.RFC3339),
				replayEnd.Format(time.RFC3339),
			)
		}
	}
}
