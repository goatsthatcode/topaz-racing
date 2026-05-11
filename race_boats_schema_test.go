package topazracing

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

type raceBoatsSchemaDocument struct {
	Required []string `json:"required"`
	Defs     struct {
		Boat struct {
			Required []string `json:"required"`
		} `json:"boat"`
		TrackPoint struct {
			Required   []string `json:"required"`
			Properties struct {
				Time struct {
					Format string `json:"format"`
				} `json:"time"`
			} `json:"properties"`
		} `json:"trackPoint"`
	} `json:"$defs"`
}

type raceBoatsFile struct {
	Boats []raceBoat `json:"boats"`
}

type raceBoat struct {
	ID       string           `json:"id"`
	Name     string           `json:"name"`
	Color    string           `json:"color"`
	BoatType string           `json:"boatType"`
	Source   string           `json:"source"`
	IsSelf   bool             `json:"isSelf"`
	Track    []raceTrackPoint `json:"track"`
}

type raceTrackPoint struct {
	Time string  `json:"time"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

func TestRaceBoatsSchemaDeclaresV1Contract(t *testing.T) {
	var schema raceBoatsSchemaDocument
	readJSONFixture(t, filepath.Join("schemas", "race-boats-v1.schema.json"), &schema)

	assertStringSet(t, schema.Required, []string{"boats"})
	assertStringSet(t, schema.Defs.Boat.Required, []string{"id", "name", "color", "boatType", "source", "isSelf", "track"})
	assertStringSet(t, schema.Defs.TrackPoint.Required, []string{"time", "lat", "lon"})

	if schema.Defs.TrackPoint.Properties.Time.Format != "date-time" {
		t.Fatalf("expected track point time format to be date-time, got %q", schema.Defs.TrackPoint.Properties.Time.Format)
	}
}

func TestReferenceRaceBoatsSatisfyV1Contract(t *testing.T) {
	var boats raceBoatsFile
	readJSONFixture(
		t,
		filepath.Join("content", "races", "dan-byrne-2025", "bishop-rock-race", "boats.json"),
		&boats,
	)

	if len(boats.Boats) < 2 {
		t.Fatalf("expected at least 2 boats, got %d", len(boats.Boats))
	}

	selfBoats := 0
	for i, boat := range boats.Boats {
		if boat.ID == "" || boat.Name == "" || boat.Color == "" || boat.BoatType == "" || boat.Source == "" {
			t.Fatalf("expected boat %d to have all required metadata populated", i)
		}
		if len(boat.Track) < 2 {
			t.Fatalf("expected boat %q to have at least 2 track points", boat.ID)
		}
		if boat.IsSelf {
			selfBoats++
		}

		var previous time.Time
		for j, point := range boat.Track {
			assertCoordinateInRange(t, point.Lat, point.Lon)

			parsed, err := time.Parse(time.RFC3339, point.Time)
			if err != nil {
				t.Fatalf("expected boat %q track point %d to use RFC3339 time: %v", boat.ID, j, err)
			}
			if j > 0 && !parsed.After(previous) {
				t.Fatalf("expected boat %q track times to be strictly increasing", boat.ID)
			}
			previous = parsed
		}
	}

	if selfBoats != 1 {
		t.Fatalf("expected exactly one self boat, got %d", selfBoats)
	}
}

func TestRaceBoatsSchemaDocReferencesSchemaArtifact(t *testing.T) {
	content, err := os.ReadFile(filepath.Join("docs", "race-boats-schema.md"))
	if err != nil {
		t.Fatalf("failed to read boats schema doc: %v", err)
	}

	if !strings.Contains(string(content), "schemas/race-boats-v1.schema.json") {
		t.Fatal("expected boats schema doc to reference the machine-readable schema")
	}
}

func TestReferenceRaceBoatsIncludeSelfAndCompetitor(t *testing.T) {
	var boats raceBoatsFile
	readJSONFixture(
		t,
		filepath.Join("content", "races", "dan-byrne-2025", "bishop-rock-race", "boats.json"),
		&boats,
	)

	ids := make([]string, 0, len(boats.Boats))
	hasCompetitor := false
	for _, boat := range boats.Boats {
		ids = append(ids, boat.ID)
		if !boat.IsSelf {
			hasCompetitor = true
		}
	}

	if !slices.Contains(ids, "topaz") {
		t.Fatal("expected reference boats.json to include the self boat")
	}
	if !hasCompetitor {
		t.Fatal("expected reference boats.json to include at least one competitor")
	}
}
