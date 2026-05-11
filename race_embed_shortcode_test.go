package topazracing

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRaceVizShortcodeBuildsReferenceRacePage(t *testing.T) {
	outputDir := t.TempDir()

	cmd := exec.Command("hugo", "--destination", outputDir)
	cmd.Env = os.Environ()
	cmd.Dir = "."

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("hugo build failed: %v\n%s", err, output)
	}

	pagePath := filepath.Join(outputDir, "races", "dan-byrne-2025", "bishop-rock-race", "index.html")
	htmlBytes, err := os.ReadFile(pagePath)
	if err != nil {
		t.Fatalf("failed to read rendered race page: %v", err)
	}

	html := string(htmlBytes)
	expectedSnippets := []string{
		`data-race-viz`,
		`data-race-id="dan-byrne-2025/bishop-rock-race"`,
		`data-race-mode="replay"`,
		`data-course-url="races/dan-byrne-2025/bishop-rock-race/course.json"`,
		`data-boats-url="races/dan-byrne-2025/bishop-rock-race/boats.json"`,
		`data-events-url="races/dan-byrne-2025/bishop-rock-race/events.json"`,
	}

	for _, snippet := range expectedSnippets {
		if !strings.Contains(html, snippet) {
			t.Fatalf("rendered race page missing %q", snippet)
		}
	}
}

func TestRaceVizShortcodeBuildsCourseModeWithoutReplayPayloads(t *testing.T) {
	outputDir := t.TempDir()

	courseData, err := os.ReadFile(filepath.Join("content", "races", "dan-byrne-2025", "bishop-rock-race", "course.json"))
	if err != nil {
		t.Fatalf("failed to read reference course fixture: %v", err)
	}

	tempBundleDir := filepath.Join("content", "races", "test-season", "course-only-race")
	if err := os.MkdirAll(tempBundleDir, 0o755); err != nil {
		t.Fatalf("failed to create temporary race bundle: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempBundleDir, "index.md"), []byte("+++\ntitle = 'Course Only Race'\n+++\n"), 0o644); err != nil {
		t.Fatalf("failed to write temporary race bundle index: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempBundleDir, "course.json"), courseData, 0o644); err != nil {
		t.Fatalf("failed to write temporary race bundle course: %v", err)
	}
	t.Cleanup(func() {
		_ = os.RemoveAll(filepath.Join("content", "races", "test-season"))
	})

	tempPagePath := filepath.Join("content", "post", "race-viz-course-only-test.md")
	tempPage := `---
title: "Race Viz Course Only Test"
---

{{< race-viz race="test-season/course-only-race" mode="course" >}}
`
	if err := os.WriteFile(tempPagePath, []byte(tempPage), 0o644); err != nil {
		t.Fatalf("failed to write temporary host page: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Remove(tempPagePath)
	})

	cmd := exec.Command("hugo", "--destination", outputDir)
	cmd.Env = os.Environ()
	cmd.Dir = "."

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("hugo build failed: %v\n%s", err, output)
	}

	pagePath := filepath.Join(outputDir, "post", "race-viz-course-only-test", "index.html")
	htmlBytes, err := os.ReadFile(pagePath)
	if err != nil {
		t.Fatalf("failed to read rendered course-only page: %v", err)
	}

	html := string(htmlBytes)
	expectedSnippets := []string{
		`data-race-id="test-season/course-only-race"`,
		`data-race-mode="course"`,
		`data-course-url="races/test-season/course-only-race/course.json"`,
	}

	for _, snippet := range expectedSnippets {
		if !strings.Contains(html, snippet) {
			t.Fatalf("rendered course-only page missing %q", snippet)
		}
	}

	if strings.Contains(html, `data-boats-url=`) {
		t.Fatal("expected course-mode embed without boats.json to omit data-boats-url")
	}
	if strings.Contains(html, `data-events-url=`) {
		t.Fatal("expected course-mode embed without events.json to omit data-events-url")
	}
}
