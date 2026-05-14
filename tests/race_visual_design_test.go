package tests

import (
	"os"
	"testing"
)

// ─── DESIGN-1: Side-by-side layout ───────────────────────────────────────────

func TestSideBySideLayoutMediaQueryPresent(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "css", "race-viz.css"))
	if err != nil {
		t.Fatalf("failed to read race-viz.css: %v", err)
	}
	assertContains(t, string(data), "min-width: 56rem")
}

func TestSideBySideLayoutShellUsesGrid(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "css", "race-viz.css"))
	if err != nil {
		t.Fatalf("failed to read race-viz.css: %v", err)
	}
	text := string(data)
	assertContains(t, text, "grid-template-columns: 1fr 19rem")
	assertContains(t, text, "grid-template-rows: 1fr auto")
}

func TestSideBySideLayoutGridAreasAssigned(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "css", "race-viz.css"))
	if err != nil {
		t.Fatalf("failed to read race-viz.css: %v", err)
	}
	text := string(data)
	// stage and sidebar appear in the grid-template-areas value
	assertContains(t, text, `"stage  sidebar"`)
	assertContains(t, text, `"caption caption"`)
	// each element is assigned to its area
	assertContains(t, text, "grid-area: stage")
	assertContains(t, text, "grid-area: sidebar")
	assertContains(t, text, "grid-area: caption")
}

func TestSideBySideLayoutSidebarHasBorderLeft(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "css", "race-viz.css"))
	if err != nil {
		t.Fatalf("failed to read race-viz.css: %v", err)
	}
	assertContains(t, string(data), "border-left: 1px solid var(--site-border-subtle)")
}

func TestSideBySideLayoutSidebarScrollable(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "css", "race-viz.css"))
	if err != nil {
		t.Fatalf("failed to read race-viz.css: %v", err)
	}
	// The sidebar within the media query gets overflow-y: auto so it can
	// scroll independently when the map stage is shorter than the fleet list.
	assertContains(t, string(data), "overflow-y: auto")
}

func TestSideBySideLayoutReducedStageHeightInRacePage(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "css", "race-viz.css"))
	if err != nil {
		t.Fatalf("failed to read race-viz.css: %v", err)
	}
	// A 65svh min-height is set in the side-by-side context so the stage does
	// not shrink below a usable height while avoiding excessive vertical space.
	assertContains(t, string(data), "clamp(24rem, 65svh, 56rem)")
}

// ─── DESIGN-4: Small visual polish items ─────────────────────────────────────

func TestDesign4BoatSwatchIsTickMark(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "css", "race-viz.css"))
	if err != nil {
		t.Fatalf("failed to read race-viz.css: %v", err)
	}
	text := string(data)
	// The swatch is a horizontal tick (1.2rem × 3px) that matches the track
	// line appearance on the map, rather than a generic color dot.
	assertContains(t, text, "width: 1.2rem")
	assertContains(t, text, "height: 3px")
}

func TestDesign4SpeedSelectBorderRadius(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "css", "race-viz.css"))
	if err != nil {
		t.Fatalf("failed to read race-viz.css: %v", err)
	}
	// The speed selector uses 0.35rem radius — crisper than the old 0.75rem
	// and consistent with the flat design direction.
	assertContains(t, string(data), "border-radius: 0.35rem")
}

func TestDesign4SidebarTitleLetterSpacing(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "css", "race-viz.css"))
	if err != nil {
		t.Fatalf("failed to read race-viz.css: %v", err)
	}
	// 0.1em is tighter than the old 0.16em — more legible in a narrow column.
	assertContains(t, string(data), "letter-spacing: 0.1em")
}

func TestDesign4ShellBorderRadius(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "css", "race-viz.css"))
	if err != nil {
		t.Fatalf("failed to read race-viz.css: %v", err)
	}
	// 0.5rem shell border-radius is sharper than 1rem for embedded instances.
	assertContains(t, string(data), "border-radius: 0.5rem")
}
