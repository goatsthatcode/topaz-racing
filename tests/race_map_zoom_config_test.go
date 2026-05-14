package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestRaceVizShortcodeEmitsFitMaxZoomAttribute verifies that when the fitMaxZoom
// shortcode parameter is provided, the data-race-viz-fit-max-zoom attribute is
// emitted on the root element with the configured value.
func TestRaceVizShortcodeEmitsFitMaxZoomAttribute(t *testing.T) {
	outputDir := t.TempDir()

	courseData, err := os.ReadFile(repoFile("content", "races", "dan-byrne-2025", "bishop-rock-race", "course.json"))
	if err != nil {
		t.Fatalf("failed to read reference course fixture: %v", err)
	}
	boatsData, err := os.ReadFile(repoFile("content", "races", "dan-byrne-2025", "bishop-rock-race", "boats.json"))
	if err != nil {
		t.Fatalf("failed to read reference boats fixture: %v", err)
	}

	bundleDir := repoFile("content", "races", "test-season", "zoom-fit-test-race")
	if err := os.MkdirAll(bundleDir, 0o755); err != nil {
		t.Fatalf("failed to create temp race bundle: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(repoFile("content", "races", "test-season")) })

	if err := os.WriteFile(filepath.Join(bundleDir, "index.md"), []byte("+++\ntitle = 'Zoom Fit Test Race'\n+++\n"), 0o644); err != nil {
		t.Fatalf("failed to write index.md: %v", err)
	}
	if err := os.WriteFile(filepath.Join(bundleDir, "course.json"), courseData, 0o644); err != nil {
		t.Fatalf("failed to write course.json: %v", err)
	}
	if err := os.WriteFile(filepath.Join(bundleDir, "boats.json"), boatsData, 0o644); err != nil {
		t.Fatalf("failed to write boats.json: %v", err)
	}

	pagePath := repoFile("content", "post", "zoom-fit-config-test.md")
	pageContent := `---
title: "Zoom Fit Config Test"
---
{{< race-viz race="test-season/zoom-fit-test-race" fitMaxZoom="9" >}}
`
	if err := os.WriteFile(pagePath, []byte(pageContent), 0o644); err != nil {
		t.Fatalf("failed to write test page: %v", err)
	}
	t.Cleanup(func() { _ = os.Remove(pagePath) })

	cmd := exec.Command("hugo", "--destination", outputDir)
	cmd.Env = os.Environ()
	cmd.Dir = repoRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("hugo build failed: %v\n%s", err, output)
	}

	html := readBuiltHTML(t, filepath.Join(outputDir, "post", "zoom-fit-config-test", "index.html"))
	assertContains(t, html, `data-race-viz-fit-max-zoom="9"`)
}

// TestRaceVizShortcodeEmitsMapMinZoomAttribute verifies that when the mapMinZoom
// shortcode parameter is provided, the data-race-viz-map-min-zoom attribute is
// emitted on the root element with the configured value.
func TestRaceVizShortcodeEmitsMapMinZoomAttribute(t *testing.T) {
	outputDir := t.TempDir()

	courseData, err := os.ReadFile(repoFile("content", "races", "dan-byrne-2025", "bishop-rock-race", "course.json"))
	if err != nil {
		t.Fatalf("failed to read reference course fixture: %v", err)
	}
	boatsData, err := os.ReadFile(repoFile("content", "races", "dan-byrne-2025", "bishop-rock-race", "boats.json"))
	if err != nil {
		t.Fatalf("failed to read reference boats fixture: %v", err)
	}

	bundleDir := repoFile("content", "races", "test-season", "zoom-min-test-race")
	if err := os.MkdirAll(bundleDir, 0o755); err != nil {
		t.Fatalf("failed to create temp race bundle: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(repoFile("content", "races", "test-season")) })

	if err := os.WriteFile(filepath.Join(bundleDir, "index.md"), []byte("+++\ntitle = 'Zoom Min Test Race'\n+++\n"), 0o644); err != nil {
		t.Fatalf("failed to write index.md: %v", err)
	}
	if err := os.WriteFile(filepath.Join(bundleDir, "course.json"), courseData, 0o644); err != nil {
		t.Fatalf("failed to write course.json: %v", err)
	}
	if err := os.WriteFile(filepath.Join(bundleDir, "boats.json"), boatsData, 0o644); err != nil {
		t.Fatalf("failed to write boats.json: %v", err)
	}

	pagePath := repoFile("content", "post", "zoom-min-config-test.md")
	pageContent := `---
title: "Zoom Min Config Test"
---
{{< race-viz race="test-season/zoom-min-test-race" mapMinZoom="5" >}}
`
	if err := os.WriteFile(pagePath, []byte(pageContent), 0o644); err != nil {
		t.Fatalf("failed to write test page: %v", err)
	}
	t.Cleanup(func() { _ = os.Remove(pagePath) })

	cmd := exec.Command("hugo", "--destination", outputDir)
	cmd.Env = os.Environ()
	cmd.Dir = repoRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("hugo build failed: %v\n%s", err, output)
	}

	html := readBuiltHTML(t, filepath.Join(outputDir, "post", "zoom-min-config-test", "index.html"))
	assertContains(t, html, `data-race-viz-map-min-zoom="5"`)
}

// TestRaceVizShortcodeOmitsMapMinZoomByDefault verifies that the
// data-race-viz-map-min-zoom attribute is absent when mapMinZoom is not provided,
// so the map imposes no explicit minimum zoom.
func TestRaceVizShortcodeOmitsMapMinZoomByDefault(t *testing.T) {
	outputDir := t.TempDir()

	cmd := exec.Command("hugo", "--destination", outputDir)
	cmd.Env = os.Environ()
	cmd.Dir = repoRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("hugo build failed: %v\n%s", err, output)
	}

	html := readBuiltHTML(t, filepath.Join(outputDir, "races", "dan-byrne-2025", "bishop-rock-race", "index.html"))
	assertNotContains(t, html, `data-race-viz-map-min-zoom=`)
}

// TestRaceVizBootstrapRoutesZoomConfigThroughConfig verifies that createRaceVizConfig
// reads both fitMaxZoom and mapMinZoom from the root dataset and exposes them as
// first-class config fields rather than having consumers reach into root.dataset directly.
func TestRaceVizBootstrapRoutesZoomConfigThroughConfig(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "js", "race-viz.js"))
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)
	expectedSnippets := []string{
		// fitMaxZoom is read from dataset inside createRaceVizConfig
		`fitMaxZoom: parseFloat(root.dataset.raceVizFitMaxZoom`,
		`DEFAULT_MAP_FIT_MAX_ZOOM`,
		// mapMinZoom is read from dataset inside createRaceVizConfig map section
		`raceVizMapMinZoom`,
		`minZoom`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, source, snippet)
	}
}

// TestRaceVizBootstrapFitCourseBoundsUsesConfig verifies that fitCourseBounds
// reads fitMaxZoom from the config object rather than directly from root.dataset,
// keeping the config object the single source of truth for component configuration.
func TestRaceVizBootstrapFitCourseBoundsUsesConfig(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "js", "race-viz.js"))
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)

	// fitCourseBounds should read fitMaxZoom from config, not from root.dataset
	assertContains(t, source, `config?.fitMaxZoom`)
	// The old direct dataset read must not be present in fitCourseBounds
	assertNotContains(t, source, `root?.dataset?.raceVizFitMaxZoom`)

	// fitCourseBounds call site must pass state.config, not root
	assertContains(t, source, `fitCourseBounds(map, courseFeatures, state.config)`)
}

// TestRaceVizBootstrapPassesMinZoomToMapConstructor verifies that mapMinZoom from
// the config is conditionally applied to the MapLibre Map constructor options,
// allowing per-embed control of the minimum zoom level independent of maxBounds.
func TestRaceVizBootstrapPassesMinZoomToMapConstructor(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "js", "race-viz.js"))
	if err != nil {
		t.Fatalf("failed to read race viz bootstrap: %v", err)
	}

	source := string(data)
	expectedSnippets := []string{
		`state.config.map.minZoom`,
		`mapOptions.minZoom`,
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, source, snippet)
	}
}
