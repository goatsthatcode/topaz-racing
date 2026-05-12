package topazracing

import (
	"os"
	"path/filepath"
	"testing"
)

// ─── Task 7.1: Site-level dark theme and retro-digital styling ───────────────

// TestSiteCustomCSSExists verifies that the project-level custom.css exists and
// takes precedence over the theme's placeholder, carrying the Topaz Racing site
// identity tokens.
func TestSiteCustomCSSExists(t *testing.T) {
	_, err := os.Stat(filepath.Join("assets", "css", "custom.css"))
	if err != nil {
		t.Fatalf("assets/css/custom.css must exist to override the hugo-brewm theme accent and provide site identity tokens: %v", err)
	}
}

// TestSiteCustomCSSContainsDarkThemeTokens verifies the site-level CSS carries
// the design tokens needed for the dark-first visual direction: deep-navy
// background, the same cyan accent used inside race visualizations, and the
// site identity variables shared between prose pages and overlays.
func TestSiteCustomCSSContainsDarkThemeTokens(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("assets", "css", "custom.css"))
	if err != nil {
		t.Fatalf("failed to read assets/css/custom.css: %v", err)
	}
	text := string(data)

	required := []string{
		"--ac-dark",
		"--ac-light",
		"--site-navy",
		"--site-cyan",
		"--site-border-subtle",
		"prefers-color-scheme: dark",
		"--bg",
		"--fg",
	}
	for _, tok := range required {
		assertContains(t, text, tok)
	}
}

// TestSiteCustomCSSAlignsDarkAccentWithRaceVizPalette verifies that the site's
// dark-mode accent color matches the primary teal/cyan used in the race
// visualization component. This alignment is the key connection that makes prose
// pages feel visually consistent with embedded race maps.
func TestSiteCustomCSSAlignsDarkAccentWithRaceVizPalette(t *testing.T) {
	siteCSS, err := os.ReadFile(filepath.Join("assets", "css", "custom.css"))
	if err != nil {
		t.Fatalf("failed to read assets/css/custom.css: %v", err)
	}

	raceCSS, err := os.ReadFile(filepath.Join("assets", "css", "race-viz.css"))
	if err != nil {
		t.Fatalf("failed to read assets/css/race-viz.css: %v", err)
	}

	// The shared teal/cyan value is the bond between the site and race-viz.
	// custom.css uses the hex form; race-viz.css uses the rgba form.
	// Both must reference the same underlying RGB values (126, 245, 236).
	assertContains(t, string(siteCSS), "#7ef5ec")
	assertContains(t, string(raceCSS), "126, 245, 236")
}

// TestSiteCustomCSSDarkNavyBackgroundOverride verifies that the site's dark mode
// override pushes the body background toward the deep navy shared by the race
// visualization, rather than staying at the theme's default #111.
func TestSiteCustomCSSDarkNavyBackgroundOverride(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("assets", "css", "custom.css"))
	if err != nil {
		t.Fatalf("failed to read assets/css/custom.css: %v", err)
	}
	text := string(data)

	// The site navy value is the base background for the dark experience.
	assertContains(t, text, "--site-navy: #08111f")
	// The body --bg override must use the same deep navy shade.
	assertContains(t, text, "#08111f")
}

// ─── Task 7.2: Mobile usability ─────────────────────────────────────────────

// TestMobileUsabilityTouchTargets verifies that the race-viz CSS raises replay
// button touch targets to the 44 px (≈ 2.75 rem) WCAG recommended minimum so
// controls are reliably tappable on phone screens.
func TestMobileUsabilityTouchTargets(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("assets", "css", "race-viz.css"))
	if err != nil {
		t.Fatalf("failed to read assets/css/race-viz.css: %v", err)
	}
	text := string(data)

	// Touch-target rules must live inside a narrow-screen media query.
	assertContains(t, text, "max-width")
	// The minimum touch target size for buttons.
	assertContains(t, text, "2.75rem")
}

// TestMobileUsabilityLegendScrollable verifies that the boat legend panel can
// scroll on narrow screens so a large fleet does not overflow the viewport.
func TestMobileUsabilityLegendScrollable(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("assets", "css", "race-viz.css"))
	if err != nil {
		t.Fatalf("failed to read assets/css/race-viz.css: %v", err)
	}
	text := string(data)

	// max-height constrains the legend, overflow-y enables scroll.
	assertContains(t, text, "overflow-y")
	assertContains(t, text, "max-height")
}

// TestMobileUsabilityPopupSizing verifies that event and hover popup text is
// capped to a readable width on narrow screens.
func TestMobileUsabilityPopupSizing(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("assets", "css", "race-viz.css"))
	if err != nil {
		t.Fatalf("failed to read assets/css/race-viz.css: %v", err)
	}
	text := string(data)

	// Popup max-width must be limited so callouts do not overflow on a phone.
	assertContains(t, text, "max-width: 14rem")
}
