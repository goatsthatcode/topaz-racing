package topazracing

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRaceVizAssetsOnlyLoadOnPagesThatRenderTheEmbed(t *testing.T) {
	outputDir := t.TempDir()

	cmd := exec.Command("hugo", "--destination", outputDir)
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), "HUGO_ENVIRONMENT=development")

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("hugo build failed: %v\n%s", err, output)
	}

	racePage := readBuiltHTML(t, filepath.Join(outputDir, "races", "dan-byrne-2025", "bishop-rock-race", "index.html"))
	if !strings.Contains(racePage, `data-race-viz-asset="css"`) {
		t.Fatal("expected race page to load race visualization CSS")
	}
	if !strings.Contains(racePage, `data-race-viz-asset="js"`) {
		t.Fatal("expected race page to load race visualization JS")
	}
	if !strings.Contains(racePage, `data-race-viz-asset="maplibre-css"`) {
		t.Fatal("expected race page to load MapLibre CSS")
	}
	if !strings.Contains(racePage, `data-race-viz-asset="maplibre-js"`) {
		t.Fatal("expected race page to load MapLibre JS")
	}

	homePage := readBuiltHTML(t, filepath.Join(outputDir, "index.html"))
	if strings.Contains(homePage, `data-race-viz-asset="css"`) ||
		strings.Contains(homePage, `data-race-viz-asset="js"`) ||
		strings.Contains(homePage, `data-race-viz-asset="maplibre-css"`) ||
		strings.Contains(homePage, `data-race-viz-asset="maplibre-js"`) {
		t.Fatal("expected non-embedded home page to omit race visualization assets")
	}
}

func readBuiltHTML(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read built HTML %s: %v", path, err)
	}

	return string(data)
}
