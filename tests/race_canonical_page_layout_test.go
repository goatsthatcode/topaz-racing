package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRacePageUsesMapFirstLayout(t *testing.T) {
	outputDir := t.TempDir()

	cmd := exec.Command("hugo", "--destination", outputDir)
	cmd.Env = append(os.Environ(), "HUGO_ENVIRONMENT=development")
	cmd.Dir = repoRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("hugo build failed: %v\n%s", err, output)
	}

	html := readBuiltHTML(t, filepath.Join(outputDir, "races", "dan-byrne-2025", "bishop-rock-race", "index.html"))

	if !strings.Contains(html, `class="race-page"`) {
		t.Error("expected race page to have class=\"race-page\" from races/single.html layout")
	}
	if !strings.Contains(html, `data-race-page`) {
		t.Error("expected race page to have data-race-page attribute")
	}

	if !strings.Contains(html, `data-race-viz-asset="css"`) {
		t.Error("expected race page to load race-viz CSS")
	}
	if !strings.Contains(html, `data-race-viz-asset="js"`) {
		t.Error("expected race page to load race-viz JS")
	}

	if !strings.Contains(html, `class="race-page-map"`) {
		t.Error("expected race page to have a race-page-map section")
	}
	if !strings.Contains(html, `data-race-viz`) {
		t.Error("expected race-viz figure to be present on the race page")
	}

	headerIdx := strings.Index(html, "race-page-header")
	mapIdx := strings.Index(html, "race-page-map")
	if headerIdx < 0 {
		t.Error("expected race-page-header to be present")
	}
	if mapIdx < 0 {
		t.Error("expected race-page-map to be present")
	}
	if headerIdx >= 0 && mapIdx >= 0 && headerIdx > mapIdx {
		t.Error("expected race-page-header to appear before race-page-map in the DOM")
	}

	figureIdx := strings.Index(html, `class="race-viz-shell"`)
	if figureIdx < 0 {
		t.Error("expected race-viz-shell figure to be present")
	}
	if mapIdx >= 0 && figureIdx >= 0 && figureIdx < mapIdx {
		t.Error("expected race-viz-shell to appear inside race-page-map (after its opening tag)")
	}
}

func TestRacePageLayoutDoesNotLoadOnNonRacePages(t *testing.T) {
	outputDir := t.TempDir()

	cmd := exec.Command("hugo", "--destination", outputDir)
	cmd.Env = append(os.Environ(), "HUGO_ENVIRONMENT=development")
	cmd.Dir = repoRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("hugo build failed: %v\n%s", err, output)
	}

	home := readBuiltHTML(t, filepath.Join(outputDir, "index.html"))

	if strings.Contains(home, `data-race-page`) {
		t.Error("expected non-race home page to omit data-race-page attribute")
	}
	if strings.Contains(home, `class="race-page"`) {
		t.Error("expected non-race home page to omit race-page layout class")
	}
}

func TestRacePageMapSectionContainsVisualizationBeforeProse(t *testing.T) {
	outputDir := t.TempDir()

	cmd := exec.Command("hugo", "--destination", outputDir)
	cmd.Env = append(os.Environ(), "HUGO_ENVIRONMENT=development")
	cmd.Dir = repoRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("hugo build failed: %v\n%s", err, output)
	}

	html := readBuiltHTML(t, filepath.Join(outputDir, "races", "dan-byrne-2025", "bishop-rock-race", "index.html"))

	bodyStart := strings.Index(html, "<body")
	if bodyStart < 0 {
		t.Fatal("expected <body> element in rendered page")
	}
	body := html[bodyStart:]

	figureIdx := strings.Index(body, `class="race-viz-shell"`)
	proseIdx := strings.Index(body, "reference pattern")
	if figureIdx < 0 {
		t.Fatal("expected race-viz-shell figure in rendered page body")
	}
	if proseIdx < 0 {
		t.Fatal("expected prose text in rendered page body")
	}
	if figureIdx > proseIdx {
		t.Errorf("expected race-viz figure (idx=%d) to appear before prose text (idx=%d) for map-first ordering", figureIdx, proseIdx)
	}
}
