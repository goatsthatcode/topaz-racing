package tests

import (
	"os"
	"testing"
)

func TestSiteCustomCSSExists(t *testing.T) {
	_, err := os.Stat(repoFile("assets", "css", "custom.css"))
	if err != nil {
		t.Fatalf("assets/css/custom.css must exist to override the hugo-brewm theme accent and provide site identity tokens: %v", err)
	}
}

func TestSiteCustomCSSContainsDarkThemeTokens(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "css", "custom.css"))
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

func TestSiteCustomCSSAlignsDarkAccentWithRaceVizPalette(t *testing.T) {
	siteCSS, err := os.ReadFile(repoFile("assets", "css", "custom.css"))
	if err != nil {
		t.Fatalf("failed to read assets/css/custom.css: %v", err)
	}

	raceCSS, err := os.ReadFile(repoFile("assets", "css", "race-viz.css"))
	if err != nil {
		t.Fatalf("failed to read assets/css/race-viz.css: %v", err)
	}

	assertContains(t, string(siteCSS), "#7ef5ec")
	assertContains(t, string(raceCSS), "126, 245, 236")
}

func TestSiteCustomCSSDarkNavyBackgroundOverride(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "css", "custom.css"))
	if err != nil {
		t.Fatalf("failed to read assets/css/custom.css: %v", err)
	}
	text := string(data)

	assertContains(t, text, "--site-navy: #08111f")
	assertContains(t, text, "#08111f")
}

func TestMobileUsabilityTouchTargets(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "css", "race-viz.css"))
	if err != nil {
		t.Fatalf("failed to read assets/css/race-viz.css: %v", err)
	}
	text := string(data)

	assertContains(t, text, "max-width")
	assertContains(t, text, "2.75rem")
}

func TestMobileUsabilityLegendScrollable(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "css", "race-viz.css"))
	if err != nil {
		t.Fatalf("failed to read assets/css/race-viz.css: %v", err)
	}
	text := string(data)

	assertContains(t, text, "overflow-y")
	assertContains(t, text, "max-height")
}

func TestMobileUsabilityPopupSizing(t *testing.T) {
	data, err := os.ReadFile(repoFile("assets", "css", "race-viz.css"))
	if err != nil {
		t.Fatalf("failed to read assets/css/race-viz.css: %v", err)
	}
	text := string(data)

	assertContains(t, text, "max-width: 14rem")
}
