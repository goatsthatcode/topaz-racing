package topazracing

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

type raceCourseSchemaDocument struct {
	Required []string `json:"required"`
	Defs     struct {
		CourseElement struct {
			Required   []string `json:"required"`
			Properties struct {
				Type struct {
					Enum []string `json:"enum"`
				} `json:"type"`
				Rounding struct {
					Enum []string `json:"enum"`
				} `json:"rounding"`
				ControlPointsToNext json.RawMessage `json:"controlPointsToNext"`
			} `json:"properties"`
		} `json:"courseElement"`
	} `json:"$defs"`
}

type raceCourseFile struct {
	ID       string              `json:"id"`
	Name     string              `json:"name"`
	Elements []raceCourseElement `json:"elements"`
}

type raceCourseElement struct {
	ID                  string           `json:"id"`
	Type                string           `json:"type"`
	Lat                 float64          `json:"lat"`
	Lon                 float64          `json:"lon"`
	Name                string           `json:"name"`
	Rounding            string           `json:"rounding"`
	ControlPointsToNext []raceCoordinate `json:"controlPointsToNext"`
}

type raceCoordinate struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func TestRaceCourseSchemaDeclaresV1Contract(t *testing.T) {
	var schema raceCourseSchemaDocument
	readJSONFixture(t, filepath.Join("schemas", "race-course-v1.schema.json"), &schema)

	assertStringSet(t, schema.Required, []string{"id", "name", "elements"})
	assertStringSet(t, schema.Defs.CourseElement.Required, []string{"id", "type", "lat", "lon", "rounding"})
	assertStringSet(t, schema.Defs.CourseElement.Properties.Type.Enum, []string{"mark", "start_line", "finish_line"})
	assertStringSet(t, schema.Defs.CourseElement.Properties.Rounding.Enum, []string{"port", "starboard", "none"})

	if len(schema.Defs.CourseElement.Properties.ControlPointsToNext) == 0 {
		t.Fatal("expected controlPointsToNext to be defined for manual route shaping")
	}
}

func TestReferenceRaceCourseSatisfiesV1Contract(t *testing.T) {
	var course raceCourseFile
	readJSONFixture(
		t,
		filepath.Join("content", "races", "dan-byrne-2025", "bishop-rock-race", "course.json"),
		&course,
	)

	if course.ID == "" || course.Name == "" {
		t.Fatal("expected course metadata to be populated")
	}
	if len(course.Elements) < 2 {
		t.Fatalf("expected at least 2 ordered course elements, got %d", len(course.Elements))
	}

	for i, element := range course.Elements {
		if element.ID == "" {
			t.Fatalf("expected element %d to have an id", i)
		}
		if !slices.Contains([]string{"mark", "start_line", "finish_line"}, element.Type) {
			t.Fatalf("unexpected element type %q", element.Type)
		}
		if !slices.Contains([]string{"port", "starboard", "none"}, element.Rounding) {
			t.Fatalf("unexpected rounding value %q", element.Rounding)
		}
		assertCoordinateInRange(t, element.Lat, element.Lon)

		for _, point := range element.ControlPointsToNext {
			assertCoordinateInRange(t, point.Lat, point.Lon)
		}
	}
}

func assertCoordinateInRange(t *testing.T, lat, lon float64) {
	t.Helper()

	if lat < -90 || lat > 90 {
		t.Fatalf("latitude %f out of range", lat)
	}
	if lon < -180 || lon > 180 {
		t.Fatalf("longitude %f out of range", lon)
	}
}

func assertStringSet(t *testing.T, actual, expected []string) {
	t.Helper()

	if len(actual) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}

	for _, value := range expected {
		if !slices.Contains(actual, value) {
			t.Fatalf("expected %q in %v", value, actual)
		}
	}
}

func TestRaceCourseSchemaDocReferencesSchemaArtifact(t *testing.T) {
	content, err := os.ReadFile(filepath.Join("docs", "race-course-schema.md"))
	if err != nil {
		t.Fatalf("failed to read course schema doc: %v", err)
	}

	if !strings.Contains(string(content), "schemas/race-course-v1.schema.json") {
		t.Fatal("expected course schema doc to reference the machine-readable schema")
	}
}
