package tests

import (
	"os"
	"testing"
)

func TestRaceVizUsageDocCoversShortcodeAndBundleExamples(t *testing.T) {
	data, err := os.ReadFile(repoFile("docs", "race-viz-usage.md"))
	if err != nil {
		t.Fatalf("failed to read race viz usage doc: %v", err)
	}

	text := string(data)
	expectedSnippets := []string{
		"`race-viz`",
		`content/races/my-season/my-race/`,
		`{{< race-viz >}}`,
		`{{< race-viz race="dan-byrne-2025/bishop-rock-race" mode="replay"`,
		"`course.json`",
		"`boats.json`",
		"`events.json`",
		"`course`",
		"`replay`",
		"`martin`",
	}

	for _, snippet := range expectedSnippets {
		assertContains(t, text, snippet)
	}
}
