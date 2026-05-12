package tests

import (
	"os"
	"testing"
)

func TestAuthoringWorkflowDocExists(t *testing.T) {
	_, err := os.Stat(repoFile("docs", "race-authoring-workflow.md"))
	if err != nil {
		t.Fatalf("docs/race-authoring-workflow.md must exist: %v", err)
	}
}

func TestAuthoringWorkflowDocCoversDirectoryStructure(t *testing.T) {
	data, err := os.ReadFile(repoFile("docs", "race-authoring-workflow.md"))
	if err != nil {
		t.Fatalf("failed to read docs/race-authoring-workflow.md: %v", err)
	}
	text := string(data)

	required := []string{
		"content/races/",
		"index.md",
		"course.json",
		"boats.json",
		"events.json",
	}
	for _, tok := range required {
		assertContains(t, text, tok)
	}
}

func TestAuthoringWorkflowDocCoversCourseJSON(t *testing.T) {
	data, err := os.ReadFile(repoFile("docs", "race-authoring-workflow.md"))
	if err != nil {
		t.Fatalf("failed to read docs/race-authoring-workflow.md: %v", err)
	}
	text := string(data)

	required := []string{
		"course.json",
		"start_line",
		"finish_line",
		"mark",
		"rounding",
		"port",
		"starboard",
		"controlPointsToNext",
	}
	for _, tok := range required {
		assertContains(t, text, tok)
	}
}

func TestAuthoringWorkflowDocCoversBoatsJSON(t *testing.T) {
	data, err := os.ReadFile(repoFile("docs", "race-authoring-workflow.md"))
	if err != nil {
		t.Fatalf("failed to read docs/race-authoring-workflow.md: %v", err)
	}
	text := string(data)

	required := []string{
		"boats.json",
		"isSelf",
		"boatType",
		"color",
		"track",
		"hand-authored",
	}
	for _, tok := range required {
		assertContains(t, text, tok)
	}
}

func TestAuthoringWorkflowDocCoversEventsJSON(t *testing.T) {
	data, err := os.ReadFile(repoFile("docs", "race-authoring-workflow.md"))
	if err != nil {
		t.Fatalf("failed to read docs/race-authoring-workflow.md: %v", err)
	}
	text := string(data)

	required := []string{
		"events.json",
		"gybe",
		"label",
		"description",
	}
	for _, tok := range required {
		assertContains(t, text, tok)
	}
}

func TestAuthoringWorkflowDocCoversEmbedShortcode(t *testing.T) {
	data, err := os.ReadFile(repoFile("docs", "race-authoring-workflow.md"))
	if err != nil {
		t.Fatalf("failed to read docs/race-authoring-workflow.md: %v", err)
	}
	text := string(data)

	required := []string{
		"race-viz",
		`race="`,
	}
	for _, tok := range required {
		assertContains(t, text, tok)
	}
}

func TestAuthoringWorkflowDocReferencesReferenceRace(t *testing.T) {
	data, err := os.ReadFile(repoFile("docs", "race-authoring-workflow.md"))
	if err != nil {
		t.Fatalf("failed to read docs/race-authoring-workflow.md: %v", err)
	}
	text := string(data)

	assertContains(t, text, "bishop-rock-race")
}

func TestAuthoringWorkflowDocReferencesSchemas(t *testing.T) {
	data, err := os.ReadFile(repoFile("docs", "race-authoring-workflow.md"))
	if err != nil {
		t.Fatalf("failed to read docs/race-authoring-workflow.md: %v", err)
	}
	text := string(data)

	required := []string{
		"schemas/race-course-v1.schema.json",
		"schemas/race-boats-v1.schema.json",
		"schemas/race-events-v1.schema.json",
	}
	for _, tok := range required {
		assertContains(t, text, tok)
	}
}
