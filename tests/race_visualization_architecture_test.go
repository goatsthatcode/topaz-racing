package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestRaceVizUsesSharedEngineContractAcrossModes(t *testing.T) {
	outputDir := t.TempDir()

	tempPagePath := repoFile("content", "post", "race-viz-architecture-test.md")
	tempPage := `---
title: "Race Viz Architecture Test"
---

{{< race-viz race="dan-byrne-2025/bishop-rock-race" mode="course" >}}
`
	if err := os.WriteFile(tempPagePath, []byte(tempPage), 0o644); err != nil {
		t.Fatalf("failed to write temporary test page: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Remove(tempPagePath)
	})

	cmd := exec.Command("hugo", "--destination", outputDir)
	cmd.Env = os.Environ()
	cmd.Dir = repoRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("hugo build failed: %v\n%s", err, output)
	}

	replayPage := readBuiltHTML(t, filepath.Join(outputDir, "races", "dan-byrne-2025", "bishop-rock-race", "index.html"))
	coursePage := readBuiltHTML(t, filepath.Join(outputDir, "post", "race-viz-architecture-test", "index.html"))

	assertContains(t, replayPage, `data-race-viz-engine="shared-v1"`)
	assertContains(t, replayPage, `data-race-viz-shared-layers="map course tracks boats events controls"`)
	assertContains(t, replayPage, `data-race-viz-active-layers="map course tracks boats events controls"`)

	assertContains(t, coursePage, `data-race-viz-engine="shared-v1"`)
	assertContains(t, coursePage, `data-race-viz-shared-layers="map course tracks boats events controls"`)
	assertContains(t, coursePage, `data-race-viz-active-layers="map course events"`)

	for _, layer := range []string{"map", "course", "tracks", "boats", "events", "controls"} {
		snippet := `data-race-viz-layer="` + layer + `"`
		assertContains(t, replayPage, snippet)
		assertContains(t, coursePage, snippet)
	}
}

func TestRaceVisualizationArchitectureNoteDocumentsSharedModeBoundaries(t *testing.T) {
	data, err := os.ReadFile(repoFile("docs", "race-visualization-architecture.md"))
	if err != nil {
		t.Fatalf("failed to read architecture note: %v", err)
	}

	text := string(data)
	expectedSnippets := []string{
		"shared race visualization engine",
		"`course` and `replay` are configuration modes",
		"`map`",
		"`course`",
		"`tracks`",
		"`boats`",
		"`events`",
		"`controls`",
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, text, snippet)
	}
}
