package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestRaceVizPublishesCourseRenderingContract(t *testing.T) {
	outputDir := t.TempDir()

	cmd := exec.Command("hugo", "--destination", outputDir)
	cmd.Env = append(os.Environ(), "HUGO_ENVIRONMENT=development")
	cmd.Dir = repoRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("hugo build failed: %v\n%s", err, output)
	}

	racePage := readBuiltHTML(t, filepath.Join(outputDir, "races", "dan-byrne-2025", "bishop-rock-race", "index.html"))

	expectedSnippets := []string{
		`data-race-viz-course-state="idle"`,
		`data-race-viz-course-palette="signal-v1"`,
		`data-race-viz-course-source="race-viz-course"`,
		`data-race-viz-course-route-layer="race-viz-course-route"`,
		`data-race-viz-course-marks-layer="race-viz-course-marks"`,
		`data-race-viz-course-rounding-layer="race-viz-course-marks-rounding"`,
		`data-race-viz-course-start-finish-layer="race-viz-course-start-finish"`,
		`data-race-viz-course-labels-layer="race-viz-course-labels"`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, racePage, snippet)
	}
}

func TestRaceVizBootstrapImplementsCourseMapLayers(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "js", "race-viz.js"))
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)
	expectedSnippets := []string{
		`DEFAULT_COURSE_LABELS_LAYER_ID`,
		`DEFAULT_COURSE_PALETTE`,
		`DEFAULT_MAP_FIT_PADDING`,
		`DEFAULT_MAP_FIT_MAX_ZOOM`,
		`signal-v1`,
		`routeGlowColor`,
		`controlPointsToNext`,
		`type: "geojson"`,
		`race-viz-course-route`,
		`race-viz-course-marks`,
		`race-viz-course-marks-rounding`,
		`race-viz-course-start-finish`,
		`race-viz-course-labels`,
		`text-field": ["upcase", ["get", "name"]]`,
		`fitBounds`,
		`maxBounds`,
		`line-dasharray`,
		`line-blur`,
		`circle-radius`,
		`fallbackTileEndpoint`,
		`replaceTileEndpointInStyle`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, source, snippet)
	}
}

func TestRaceVizBootstrapImplementsMarkRoundingDirectionLayer(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "js", "race-viz.js"))
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)
	expectedSnippets := []string{
		`DEFAULT_COURSE_ROUNDING_LAYER_ID`,
		`roundingLayerID`,
		`roundingPortColor`,
		`roundingStarboardColor`,
		`"port", courseStyle.roundingPortColor`,
		`"starboard", courseStyle.roundingStarboardColor`,
		`["match", ["get", "rounding"], ["port", "starboard"], true, false]`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, source, snippet)
	}
}

// TestWaypointTypeIsSchemaValid verifies that course files using "waypoint"
// type elements pass schema validation now that "waypoint" is in the enum.
func TestWaypointTypeIsSchemaValid(t *testing.T) {
	var schema raceCourseSchemaDocument
	readJSONFixture(t, repoFile("schemas", "race-course-v1.schema.json"), &schema)

	typeEnum := schema.Defs.CourseElement.Properties.Type.Enum
	found := false
	for _, v := range typeEnum {
		if v == "waypoint" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected schema type enum to include \"waypoint\"")
	}
}

// TestCatalinaBacksideCourseHasWaypointElement verifies the catalina-backside
// course file (which uses "waypoint") parses correctly under the updated schema.
func TestCatalinaBacksideCourseHasWaypointElement(t *testing.T) {
	var course raceCourseFile
	readJSONFixture(
		t,
		repoFile("content", "races", "dan-byrne-2025", "catalina-backside-race", "course.json"),
		&course,
	)

	validTypes := []string{"mark", "start_line", "finish_line", "waypoint"}
	var waypointCount int
	for i, el := range course.Elements {
		validType := false
		for _, vt := range validTypes {
			if el.Type == vt {
				validType = true
				break
			}
		}
		if !validType {
			t.Fatalf("element %d has unexpected type %q", i, el.Type)
		}
		if el.Type == "waypoint" {
			waypointCount++
		}
	}
	if waypointCount == 0 {
		t.Fatal("expected at least one waypoint element in catalina-backside-race course")
	}
}

// TestLabelsLayerFilterExcludesWaypoints verifies that the labels layer filter
// in race-viz.js includes a type guard so "waypoint" elements never emit labels.
func TestLabelsLayerFilterExcludesWaypoints(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "js", "race-viz.js"))
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)

	// The labels layer must restrict rendering to named mark/start/finish elements only.
	assertContains(t, source, `["match", ["get", "type"], ["mark", "start_line", "finish_line"], true, false]`)
}

func TestRaceVizStylesDefineEditorialCourseFrame(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "css", "race-viz.css"))
	if err != nil {
		t.Fatalf("failed to read race viz stylesheet: %v", err)
	}

	source := string(data)
	expectedSnippets := []string{
		`--race-viz-frame`,
		`--race-viz-stage-outline`,
		`--race-viz-caption-accent`,
		`linear-gradient(125deg`,
		`box-shadow:`,
		`pointer-events: none;`,
		`.race-viz-shell figcaption::before`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, source, snippet)
	}
}
