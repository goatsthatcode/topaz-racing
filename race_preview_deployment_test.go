package topazracing

import (
	"os"
	"path/filepath"
	"testing"
)

// ─── Task 8.2: Local preview and deployment workflow documentation ────────────

// TestPreviewDeploymentDocExists verifies that the preview/deployment guide is
// present so the map stack can be reproduced without reverse-engineering config.
func TestPreviewDeploymentDocExists(t *testing.T) {
	_, err := os.Stat(filepath.Join("docs", "race-preview-deployment.md"))
	if err != nil {
		t.Fatalf("docs/race-preview-deployment.md must exist: %v", err)
	}
}

// TestPreviewDeploymentDocCoversLocalTileServer verifies the guide documents
// how to start the Martin tile server so race maps render in local development.
func TestPreviewDeploymentDocCoversLocalTileServer(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("docs", "race-preview-deployment.md"))
	if err != nil {
		t.Fatalf("failed to read docs/race-preview-deployment.md: %v", err)
	}
	text := string(data)

	required := []string{
		"martin",
		"tiles/martin-config",
		"127.0.0.1:3000",
		"hugo server",
	}
	for _, tok := range required {
		assertContains(t, text, tok)
	}
}

// TestPreviewDeploymentDocCoversProductionHosting verifies the guide documents
// the production tile URL contract so the deployment assumption is explicit.
func TestPreviewDeploymentDocCoversProductionHosting(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("docs", "race-preview-deployment.md"))
	if err != nil {
		t.Fatalf("failed to read docs/race-preview-deployment.md: %v", err)
	}
	text := string(data)

	required := []string{
		"topaz-racing.com/tiles",
		"combined_socal",
	}
	for _, tok := range required {
		assertContains(t, text, tok)
	}
}

// TestPreviewDeploymentDocCoversTileURLContract verifies the guide documents
// the tile URL path shape so any compatible host can serve tiles correctly.
func TestPreviewDeploymentDocCoversTileURLContract(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("docs", "race-preview-deployment.md"))
	if err != nil {
		t.Fatalf("failed to read docs/race-preview-deployment.md: %v", err)
	}
	text := string(data)

	// The URL pattern /{tileset}/{z}/{x}/{y} is the stable contract.
	assertContains(t, text, "{z}/{x}/{y}")
}

// TestPreviewDeploymentDocCoversPublishedMapArtifacts verifies the guide
// documents the two Hugo-published map artifacts that decouple the frontend
// from hard-coded tile configuration.
func TestPreviewDeploymentDocCoversPublishedMapArtifacts(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("docs", "race-preview-deployment.md"))
	if err != nil {
		t.Fatalf("failed to read docs/race-preview-deployment.md: %v", err)
	}
	text := string(data)

	required := []string{
		"race-viz/map-style.json",
		"race-viz/tile-manifest.json",
	}
	for _, tok := range required {
		assertContains(t, text, tok)
	}
}

// TestPreviewDeploymentDocCoversFutureTileGenerationPipeline verifies the guide
// documents the interface for future chart-tile generation tooling so it can
// target a stable contract instead of requiring visualization rewrites.
func TestPreviewDeploymentDocCoversFutureTileGenerationPipeline(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("docs", "race-preview-deployment.md"))
	if err != nil {
		t.Fatalf("failed to read docs/race-preview-deployment.md: %v", err)
	}
	text := string(data)

	required := []string{
		"tiles/mbtiles/",
		"params.raceViz.tiles.sets",
	}
	for _, tok := range required {
		assertContains(t, text, tok)
	}
}

// TestPreviewDeploymentDocCoversFutureRaceDataImport verifies the guide
// documents the import interface for future Jibeset/Garmin/GPX tooling so
// it targets the stable internal schema without touching visualization code.
func TestPreviewDeploymentDocCoversFutureRaceDataImport(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("docs", "race-preview-deployment.md"))
	if err != nil {
		t.Fatalf("failed to read docs/race-preview-deployment.md: %v", err)
	}
	text := string(data)

	required := []string{
		"schemas/race-boats-v1.schema.json",
		"source",
		"Jibeset",
		"Garmin",
		"GPX",
	}
	for _, tok := range required {
		assertContains(t, text, tok)
	}
}

// TestMBTilesCatalogExists verifies that the local tile catalog directory and
// the default combined_socal MBTiles file are present so local preview works
// out of the box without extra setup.
func TestMBTilesCatalogExists(t *testing.T) {
	_, err := os.Stat(filepath.Join("tiles", "mbtiles", "combined_socal.mbtiles"))
	if err != nil {
		t.Fatalf("tiles/mbtiles/combined_socal.mbtiles must exist for local preview: %v", err)
	}
}

// TestMartinConfigExists verifies the Martin tile server configuration file is
// present so the preview command documented in race-preview-deployment.md works.
func TestMartinConfigExists(t *testing.T) {
	_, err := os.Stat(filepath.Join("tiles", "martin-config"))
	if err != nil {
		t.Fatalf("tiles/martin-config must exist for local tile server preview: %v", err)
	}
}
