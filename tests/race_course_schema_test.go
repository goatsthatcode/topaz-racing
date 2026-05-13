package tests

import (
	"os"
	"slices"
	"strings"
	"testing"
)

func TestRaceCourseSchemaDeclaresV1Contract(t *testing.T) {
	var schema raceCourseSchemaDocument
	readJSONFixture(t, repoFile("schemas", "race-course-v1.schema.json"), &schema)

	assertStringSet(t, schema.Required, []string{"id", "name", "elements"})
	assertStringSet(t, schema.Defs.CourseElement.Required, []string{"id", "type", "lat", "lon", "rounding"})
	assertStringSet(t, schema.Defs.CourseElement.Properties.Type.Enum, []string{"mark", "start_line", "finish_line", "waypoint"})
	assertStringSet(t, schema.Defs.CourseElement.Properties.Rounding.Enum, []string{"port", "starboard", "none"})

	if len(schema.Defs.CourseElement.Properties.ControlPointsToNext) == 0 {
		t.Fatal("expected controlPointsToNext to be defined for manual route shaping")
	}
}

func TestReferenceRaceCourseSatisfiesV1Contract(t *testing.T) {
	var course raceCourseFile
	readJSONFixture(
		t,
		repoFile("content", "races", "dan-byrne-2025", "bishop-rock-race", "course.json"),
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
		if !slices.Contains([]string{"mark", "start_line", "finish_line", "waypoint"}, element.Type) {
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

func TestRaceCourseSchemaDocReferencesSchemaArtifact(t *testing.T) {
	content, err := os.ReadFile(repoFile("docs", "race-course-schema.md"))
	if err != nil {
		t.Fatalf("failed to read course schema doc: %v", err)
	}

	if !strings.Contains(string(content), "schemas/race-course-v1.schema.json") {
		t.Fatal("expected course schema doc to reference the machine-readable schema")
	}
}
