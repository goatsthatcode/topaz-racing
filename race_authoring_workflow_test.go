package topazracing

import (
	"os"
	"path/filepath"
	"testing"
)

// ─── Task 8.1: Race authoring workflow documentation ─────────────────────────

// TestAuthoringWorkflowDocExists verifies that the race authoring guide is
// present so authors can discover the workflow without reading source code.
func TestAuthoringWorkflowDocExists(t *testing.T) {
	_, err := os.Stat(filepath.Join("docs", "race-authoring-workflow.md"))
	if err != nil {
		t.Fatalf("docs/race-authoring-workflow.md must exist: %v", err)
	}
}

// TestAuthoringWorkflowDocCoversDirectoryStructure verifies the guide documents
// the canonical leaf-bundle directory layout so an author knows exactly where
// to put new files.
func TestAuthoringWorkflowDocCoversDirectoryStructure(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("docs", "race-authoring-workflow.md"))
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

// TestAuthoringWorkflowDocCoversCourseJSON verifies the guide explains how to
// author course.json including element types and required fields.
func TestAuthoringWorkflowDocCoversCourseJSON(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("docs", "race-authoring-workflow.md"))
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

// TestAuthoringWorkflowDocCoversBoatsJSON verifies the guide explains the boat
// track format including the isSelf boat, competitor boats, and track points.
func TestAuthoringWorkflowDocCoversBoatsJSON(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("docs", "race-authoring-workflow.md"))
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

// TestAuthoringWorkflowDocCoversEventsJSON verifies the guide explains event
// annotation authoring including common event types and optional fields.
func TestAuthoringWorkflowDocCoversEventsJSON(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("docs", "race-authoring-workflow.md"))
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

// TestAuthoringWorkflowDocCoversEmbedShortcode verifies the guide shows how to
// embed the visualization in markdown using the race-viz shortcode, including
// both the self-page form and the cross-page reference form.
func TestAuthoringWorkflowDocCoversEmbedShortcode(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("docs", "race-authoring-workflow.md"))
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

// TestAuthoringWorkflowDocReferencesReferenceRace verifies the guide points to
// the canonical reference race bundle so authors have a working example to copy.
func TestAuthoringWorkflowDocReferencesReferenceRace(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("docs", "race-authoring-workflow.md"))
	if err != nil {
		t.Fatalf("failed to read docs/race-authoring-workflow.md: %v", err)
	}
	text := string(data)

	assertContains(t, text, "bishop-rock-race")
}

// TestAuthoringWorkflowDocReferencesSchemas verifies the guide points to the
// machine-readable JSON schemas so authors can validate their files.
func TestAuthoringWorkflowDocReferencesSchemas(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("docs", "race-authoring-workflow.md"))
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
