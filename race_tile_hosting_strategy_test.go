package topazracing

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestRaceVizPublishesTileHostingContractForProductionBuilds(t *testing.T) {
	outputDir := t.TempDir()

	cmd := exec.Command("hugo", "--destination", outputDir)
	cmd.Dir = "."
	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("hugo build failed: %v\n%s", err, output)
	}

	racePage := readBuiltHTML(t, filepath.Join(outputDir, "races", "dan-byrne-2025", "bishop-rock-race", "index.html"))
	assertContains(t, racePage, `data-race-viz-map-tile-manifest-url="https://topaz-racing.com/race-viz/tile-manifest.json"`)
	assertContains(t, racePage, `data-race-viz-map-tile-set="combined_socal"`)
	assertContains(t, racePage, `data-race-viz-map-serving-mode="external-static-vector-host"`)

	manifest := readBuiltFile(t, filepath.Join(outputDir, "race-viz", "tile-manifest.json"))
	assertContains(t, manifest, `"servingMode": "external-static-vector-host"`)
	assertContains(t, manifest, `"tileEndpoint": "https://topaz-racing.com/tiles"`)
	assertContains(t, manifest, `"defaultSet": "combined_socal"`)
	assertContains(t, manifest, `"previewCommand": "martin --config tiles/martin-config"`)
	assertContains(t, manifest, `"path": "combined_socal"`)
	assertContains(t, manifest, `"source": "tiles/mbtiles/combined_socal.mbtiles"`)
	assertContains(t, manifest, `"bounds": [-121, 31.5, -116, 35]`)
}

func TestRaceVizPublishesDevelopmentTileHostingContract(t *testing.T) {
	outputDir := t.TempDir()

	cmd := exec.Command("hugo", "--destination", outputDir)
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), "HUGO_ENVIRONMENT=development")

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("hugo build failed: %v\n%s", err, output)
	}

	manifest := readBuiltFile(t, filepath.Join(outputDir, "race-viz", "tile-manifest.json"))
	assertContains(t, manifest, `"tileEndpoint": "http://127.0.0.1:3000"`)
	assertContains(t, manifest, `"prototypePage": "/tiles/index.html"`)
	assertContains(t, manifest, `"prototypeStyle": "/tiles/style.json"`)
}

func TestRaceTileHostingStrategyDocCapturesPreviewProductionAndExpansionPlan(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("docs", "race-tile-hosting-strategy.md"))
	if err != nil {
		t.Fatalf("failed to read tile hosting strategy doc: %v", err)
	}

	text := string(data)
	expectedSnippets := []string{
		"`race-viz/tile-manifest.json`",
		"`external-static-vector-host`",
		"`martin --config tiles/martin-config`",
		"`http://127.0.0.1:3000`",
		"`https://topaz-racing.com/tiles`",
		"`/{tileset}/{z}/{x}/{y}`",
		"`params.raceViz.tiles.sets`",
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, text, snippet)
	}
}
