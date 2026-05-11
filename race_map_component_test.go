package topazracing

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestRaceVizBuildsEmbeddableMapCanvasForRacePages(t *testing.T) {
	outputDir := t.TempDir()

	cmd := exec.Command("hugo", "--destination", outputDir)
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), "HUGO_ENVIRONMENT=development")

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("hugo build failed: %v\n%s", err, output)
	}

	racePage := readBuiltHTML(t, filepath.Join(outputDir, "races", "dan-byrne-2025", "bishop-rock-race", "index.html"))
	assertContains(t, racePage, `class="race-viz-stage" role="img" aria-label="Bishop Rock Race race visualization"`)
	assertContains(t, racePage, `data-race-viz-layer="map"`)
	assertContains(t, racePage, `class="race-viz-map-canvas" data-race-viz-map-canvas aria-hidden="true"`)
	assertContains(t, racePage, `src="https://unpkg.com/maplibre-gl@5.0.1/dist/maplibre-gl.js"`)
	assertContains(t, racePage, `href="https://unpkg.com/maplibre-gl@5.0.1/dist/maplibre-gl.css"`)
}
