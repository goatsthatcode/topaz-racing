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
		`race-viz-course-start-finish`,
		`race-viz-course-labels`,
		`text-field": ["upcase", ["get", "name"]]`,
		`fitBounds`,
		`setMinZoom`,
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
