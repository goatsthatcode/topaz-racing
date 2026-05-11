package topazracing

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRaceVizPublishesMapFoundationContractForProductionBuilds(t *testing.T) {
	outputDir := t.TempDir()

	cmd := exec.Command("hugo", "--destination", outputDir)
	cmd.Dir = "."
	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("hugo build failed: %v\n%s", err, output)
	}

	racePage := readBuiltHTML(t, filepath.Join(outputDir, "races", "dan-byrne-2025", "bishop-rock-race", "index.html"))
	assertContains(t, racePage, `data-race-viz-map-style-url="https://topaz-racing.com/race-viz/map-style.json"`)
	assertContains(t, racePage, `data-race-viz-map-tile-endpoint="https://topaz-racing.com/tiles"`)
	assertContains(t, racePage, `data-race-viz-map-prototype-page="/tiles/index.html"`)
	assertContains(t, racePage, `data-race-viz-map-prototype-style="/tiles/style.json"`)

	style := readBuiltFile(t, filepath.Join(outputDir, "race-viz", "map-style.json"))
	assertContains(t, style, `Topaz Racing ENC Prototype`)
	assertContains(t, style, `https://topaz-racing.com/tiles/combined_socal/{z}/{x}/{y}`)
}

func TestRaceVizPublishesDevelopmentTileEndpointInMapStyle(t *testing.T) {
	outputDir := t.TempDir()

	cmd := exec.Command("hugo", "--destination", outputDir)
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), "HUGO_ENVIRONMENT=development")

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("hugo build failed: %v\n%s", err, output)
	}

	racePage := readBuiltHTML(t, filepath.Join(outputDir, "races", "dan-byrne-2025", "bishop-rock-race", "index.html"))
	assertContains(t, racePage, `data-race-viz-map-style-url="https://topaz-racing.com/race-viz/map-style.json"`)
	assertContains(t, racePage, `data-race-viz-map-tile-endpoint="http://127.0.0.1:3000"`)

	style := readBuiltFile(t, filepath.Join(outputDir, "race-viz", "map-style.json"))
	assertContains(t, style, `http://127.0.0.1:3000/combined_socal/{z}/{x}/{y}`)
}

func TestRaceMapFoundationDocCapturesPrototypeRelationshipAndEnvironmentChoice(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("docs", "race-map-foundation.md"))
	if err != nil {
		t.Fatalf("failed to read map foundation doc: %v", err)
	}

	text := string(data)
	expectedSnippets := []string{
		"`tiles/index.html`",
		"`tiles/style.json`",
		"`assets/race-viz/map/style.json.tmpl`",
		"`layouts/partials/race-viz/map-foundation.html`",
		"`data-race-viz-map-style-url`",
		"`http://127.0.0.1:3000`",
		"`https://topaz-racing.com/tiles`",
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, text, snippet)
	}
}

func readBuiltFile(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read built file %s: %v", path, err)
	}

	return strings.TrimSpace(string(data))
}
