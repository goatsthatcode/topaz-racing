package tests

import (
	"os"
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
	var events raceEventsFile
	readJSONFixture(
		t,
		repoFile("content", "races", "dan-byrne-2025", "bishop-rock-race", "events.json"),
		&events,
	)

	if len(events.Events) < 1 {
		t.Fatal("expected at least one event annotation")
	}

	for i, event := range events.Events {
		if event.ID == "" || event.Type == "" {
			t.Fatalf("expected event %d to have required metadata", i)
		}
		if !slices.Contains([]string{"gybe", "wipeout", "sail_change", "tack", "mark_rounding", "note"}, event.Type) {
			t.Fatalf("unexpected event type %q", event.Type)
		}

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
