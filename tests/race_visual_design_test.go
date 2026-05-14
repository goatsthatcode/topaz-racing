package tests

import (
	"os"
	"testing"
)

// ─── DESIGN-2: Flat/crisp container style ────────────────────────────────────

func TestDesign2ShellHasFlatBackground(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "css", "race-viz.css"))
	if err != nil {
		t.Fatalf("failed to read race-viz.css: %v", err)
	}
	text := string(data)
	// Shell uses a single flat fill instead of layered radial+linear gradients.
	assertContains(t, text, "background: rgba(6, 14, 25, 0.97)")
	// No multi-layer gradient stack on the shell.
	assertNotContains(t, text, "radial-gradient(circle at top, rgba(64, 224, 208")
}

func TestDesign2ShellHasNoBoxShadow(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "css", "race-viz.css"))
	if err != nil {
		t.Fatalf("failed to read race-viz.css: %v", err)
	}
	// The large drop-shadow glow beneath the shell is removed.
	assertNotContains(t, string(data), "box-shadow: 0 1.25rem 3rem")
}

func TestDesign2PanelHasTransparentBackground(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "css", "race-viz.css"))
	if err != nil {
		t.Fatalf("failed to read race-viz.css: %v", err)
	}
	// Panels are flat sections — the sidebar column background provides the fill.
	assertContains(t, string(data), "background: transparent")
}

func TestDesign2PanelHasTopOnlyHairlineBorder(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "css", "race-viz.css"))
	if err != nil {
		t.Fatalf("failed to read race-viz.css: %v", err)
	}
	text := string(data)
	// Panels are separated by a subtle top hairline, not a full card border.
	assertContains(t, text, "border-top: 1px solid rgba(126, 245, 236, 0.09)")
	// The full card border and border-radius are removed from panels.
	assertNotContains(t, text, "border-radius: 0.85rem")
}

func TestDesign2ButtonHasFlatBackground(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "css", "race-viz.css"))
	if err != nil {
		t.Fatalf("failed to read race-viz.css: %v", err)
	}
	text := string(data)
	// Buttons use a flat dark fill instead of the gradient pair.
	assertContains(t, text, "background: rgba(8, 20, 34, 0.95)")
	// The gradient background on buttons is removed.
	assertNotContains(t, text, "linear-gradient(180deg, rgba(13, 32, 46")
}

func TestDesign2ButtonHasNoInnerGlow(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "css", "race-viz.css"))
	if err != nil {
		t.Fatalf("failed to read race-viz.css: %v", err)
	}
	// The inset glow ring on buttons is removed.
	assertNotContains(t, string(data), "inset 0 0 0 1px rgba(255, 255, 255, 0.04)")
}

func TestDesign2FigcaptionDotHasNoGlowRing(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "css", "race-viz.css"))
	if err != nil {
		t.Fatalf("failed to read race-viz.css: %v", err)
	}
	// The glow bloom around the caption dot is removed — radial fill alone suffices.
	assertNotContains(t, string(data), "box-shadow: 0 0 0.75rem rgba(126, 245, 236, 0.4)")
}

func TestDesign2SiteHeaderHasNoBoxShadow(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "css", "custom.css"))
	if err != nil {
		t.Fatalf("failed to read custom.css: %v", err)
	}
	// Site header border-bottom is enough; glow bloom removed.
	assertNotContains(t, string(data), "box-shadow: 0 0.5rem 2rem var(--site-cyan-dim)")
}

func TestDesign2RacePageHeaderHasNoGradient(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "css", "custom.css"))
	if err != nil {
		t.Fatalf("failed to read custom.css: %v", err)
	}
	// Race page header gradient is removed — flat background is more editorial.
	assertNotContains(t, string(data), "rgba(126, 245, 236, 0.04)")
}

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
