package topazracing

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestRacePageUsesMapFirstLayout verifies that pages in the races/ section use
// the dedicated race layout and render the map visualization as the primary
// visual element, with the prose section below it.
func TestRacePageUsesMapFirstLayout(t *testing.T) {
	outputDir := t.TempDir()

	cmd := exec.Command("hugo", "--destination", outputDir)
	cmd.Env = append(os.Environ(), "HUGO_ENVIRONMENT=development")
	cmd.Dir = "."

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("hugo build failed: %v\n%s", err, output)
	}

	html := readBuiltHTML(t, filepath.Join(outputDir, "races", "dan-byrne-2025", "bishop-rock-race", "index.html"))

	// Layout marker: the race section layout wraps content in a race-page article.
	if !strings.Contains(html, `class="race-page"`) {
		t.Error("expected race page to have class=\"race-page\" from races/single.html layout")
	}
	if !strings.Contains(html, `data-race-page`) {
		t.Error("expected race page to have data-race-page attribute")
	}

	// Race assets must always load on race pages (unconditional in the layout).
	if !strings.Contains(html, `data-race-viz-asset="css"`) {
		t.Error("expected race page to load race-viz CSS")
	}
	if !strings.Contains(html, `data-race-viz-asset="js"`) {
		t.Error("expected race page to load race-viz JS")
	}

	// Map-first composition: the race-viz figure must be present inside the
	// race-page-map section.
	if !strings.Contains(html, `class="race-page-map"`) {
		t.Error("expected race page to have a race-page-map section")
	}
	if !strings.Contains(html, `data-race-viz`) {
		t.Error("expected race-viz figure to be present on the race page")
	}

	// The page header with the title must come before the map section in the DOM
	// so that the title is visible at the top while the map dominates below.
	// The header carries multiple classes ("race-page-header pagewidth") so search
	// for the distinguishing substring rather than an exact class value.
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

	// The race viz figure must appear inside the race-page-map section, i.e. the
	// figure index must be after the section opening tag.
	figureIdx := strings.Index(html, `class="race-viz-shell"`)
	if figureIdx < 0 {
		t.Error("expected race-viz-shell figure to be present")
	}
	if mapIdx >= 0 && figureIdx >= 0 && figureIdx < mapIdx {
		t.Error("expected race-viz-shell to appear inside race-page-map (after its opening tag)")
	}
}

// TestRacePageLayoutDoesNotLoadOnNonRacePages verifies that the race-page
// layout markers are absent from ordinary blog posts.
func TestRacePageLayoutDoesNotLoadOnNonRacePages(t *testing.T) {
	outputDir := t.TempDir()

	cmd := exec.Command("hugo", "--destination", outputDir)
	cmd.Env = append(os.Environ(), "HUGO_ENVIRONMENT=development")
	cmd.Dir = "."

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

// TestRacePageMapSectionContainsVisualizationBeforeProse verifies that the
// race-viz shortcode output appears before any prose paragraphs inside the
// race-page-map section, achieving the map-first visual ordering.
func TestRacePageMapSectionContainsVisualizationBeforeProse(t *testing.T) {
	outputDir := t.TempDir()

	cmd := exec.Command("hugo", "--destination", outputDir)
	cmd.Env = append(os.Environ(), "HUGO_ENVIRONMENT=development")
	cmd.Dir = "."

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("hugo build failed: %v\n%s", err, output)
	}

	html := readBuiltHTML(t, filepath.Join(outputDir, "races", "dan-byrne-2025", "bishop-rock-race", "index.html"))

	// Restrict search to the <body> so we ignore meta-description content in
	// <head> which may also contain the prose snippet.
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
