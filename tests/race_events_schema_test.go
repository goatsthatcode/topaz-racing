package tests

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestRaceEventsSchemaDeclaresV1Contract(t *testing.T) {
	var schema raceEventsSchemaDocument
	readJSONFixture(t, repoFile("schemas", "race-events-v1.schema.json"), &schema)

	assertStringSet(t, schema.Required, []string{"events"})
	assertStringSet(t, schema.Defs.Event.Required, []string{"id", "type"})
	assertStringSet(t, schema.Defs.Event.Properties.Type.Enum, []string{"gybe", "wipeout", "sail_change", "tack", "mark_rounding", "note"})

	if schema.Defs.Event.Properties.Time.Format != "date-time" {
		t.Fatalf("expected event time format to be date-time, got %q", schema.Defs.Event.Properties.Time.Format)
	}
	if len(schema.Defs.Event.AnyOf) != 2 {
		t.Fatalf("expected event anchor anyOf with 2 options, got %d", len(schema.Defs.Event.AnyOf))
	}
	if len(schema.Defs.Event.Properties.Label) == 0 || len(schema.Defs.Event.Properties.Description) == 0 {
		t.Fatal("expected optional label and description fields to be declared")
	}
}

func TestReferenceRaceEventsSatisfyV1Contract(t *testing.T) {
	bundleDir := repoFile("content", "races", "dan-byrne-2025", "bishop-rock-race")

	var events raceEventsFile
	readJSONFixture(t, filepath.Join(bundleDir, "events.json"), &events)

	validTypes := []string{"gybe", "wipeout", "sail_change", "tack", "mark_rounding", "note"}

	if len(events.Events) < len(validTypes) {
		t.Fatalf("expected reference events.json to have at least %d events (one per type), got %d", len(validTypes), len(events.Events))
	}

	seenTypes := make(map[string]bool)
	for i, event := range events.Events {
		if event.ID == "" || event.Type == "" {
			t.Fatalf("expected event %d to have required metadata", i)
		}
		if !slices.Contains(validTypes, event.Type) {
			t.Fatalf("unexpected event type %q", event.Type)
		}
		seenTypes[event.Type] = true

		hasTime := event.Time != ""
		hasPosition := event.Lat != nil && event.Lon != nil
		if !hasTime && !hasPosition {
			t.Fatalf("expected event %q to be anchored by time or position", event.ID)
		}
		if hasTime {
			if _, err := time.Parse(time.RFC3339, event.Time); err != nil {
				t.Fatalf("expected event %q to use RFC3339 time: %v", event.ID, err)
			}
		}
		if hasPosition {
			assertCoordinateInRange(t, *event.Lat, *event.Lon)
		}
		if event.Label == "" && event.Description == "" {
			t.Fatalf("expected event %q to include editorial label or description", event.ID)
		}
	}

	for _, typ := range validTypes {
		if !seenTypes[typ] {
			t.Errorf("reference events.json is missing an event of type %q — dataset must exercise all schema-defined types", typ)
		}
	}

	// Verify time-anchored events fall within the boat track replay bounds.
	var boats raceBoatsFile
	readJSONFixture(t, filepath.Join(bundleDir, "boats.json"), &boats)

	var replayStart, replayEnd time.Time
	for i, boat := range boats.Boats {
		start, err := time.Parse(time.RFC3339, boat.Track[0].Time)
		if err != nil {
			t.Fatalf("failed to parse boat %q start time: %v", boat.ID, err)
		}
		end, err := time.Parse(time.RFC3339, boat.Track[len(boat.Track)-1].Time)
		if err != nil {
			t.Fatalf("failed to parse boat %q end time: %v", boat.ID, err)
		}
		if i == 0 || start.Before(replayStart) {
			replayStart = start
		}
		if i == 0 || end.After(replayEnd) {
			replayEnd = end
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

func TestCatalinaBacksideRaceEventsSatisfyV1Contract(t *testing.T) {
	bundleDir := repoFile("content", "races", "dan-byrne-2025", "catalina-backside-race")

	var events raceEventsFile
	readJSONFixture(t, filepath.Join(bundleDir, "events.json"), &events)

	if len(events.Events) < 3 {
		t.Fatalf("expected at least 3 events to exercise multiple annotation types, got %d", len(events.Events))
	}

	seenTypes := make(map[string]bool)
	validTypes := []string{"gybe", "wipeout", "sail_change", "tack", "mark_rounding", "note"}

	for i, event := range events.Events {
		if event.ID == "" || event.Type == "" {
			t.Fatalf("expected event %d to have id and type", i)
		}
		if !slices.Contains(validTypes, event.Type) {
			t.Fatalf("event %q has unknown type %q", event.ID, event.Type)
		}
		seenTypes[event.Type] = true

		hasTime := event.Time != ""
		hasPosition := event.Lat != nil && event.Lon != nil
		if !hasTime && !hasPosition {
			t.Fatalf("expected event %q to be anchored by time or position", event.ID)
		}
		if hasTime {
			if _, err := time.Parse(time.RFC3339, event.Time); err != nil {
				t.Fatalf("expected event %q to use RFC3339 time: %v", event.ID, err)
			}
		}
		if hasPosition {
			assertCoordinateInRange(t, *event.Lat, *event.Lon)
		}
		if event.Label == "" && event.Description == "" {
			t.Fatalf("expected event %q to include editorial label or description", event.ID)
		}
	}

	if len(seenTypes) < 2 {
		t.Fatalf("expected events to cover at least 2 distinct event types, got %d: %v", len(seenTypes), seenTypes)
	}

	// Verify event times fall within the boat track replay bounds.
	var boats raceBoatsFile
	readJSONFixture(t, filepath.Join(bundleDir, "boats.json"), &boats)

	var replayStart, replayEnd time.Time
	for i, boat := range boats.Boats {
		start, err := time.Parse(time.RFC3339, boat.Track[0].Time)
		if err != nil {
			t.Fatalf("failed to parse boat %q start time: %v", boat.ID, err)
		}
		end, err := time.Parse(time.RFC3339, boat.Track[len(boat.Track)-1].Time)
		if err != nil {
			t.Fatalf("failed to parse boat %q end time: %v", boat.ID, err)
		}
		if i == 0 || start.Before(replayStart) {
			replayStart = start
		}
		if i == 0 || end.After(replayEnd) {
			replayEnd = end
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

func TestRaceEventsSchemaDocReferencesSchemaArtifact(t *testing.T) {
	content, err := os.ReadFile(repoFile("docs", "race-events-schema.md"))
	if err != nil {
		t.Fatalf("failed to read events schema doc: %v", err)
	}

	if !strings.Contains(string(content), "schemas/race-events-v1.schema.json") {
		t.Fatal("expected events schema doc to reference the machine-readable schema")
	}
}
