package tests

import (
	"os"
	"testing"
)

func TestPreviewDeploymentDocExists(t *testing.T) {
	_, err := os.Stat(repoFile("docs", "race-preview-deployment.md"))
	if err != nil {
		t.Fatalf("docs/race-preview-deployment.md must exist: %v", err)
	}
}

func TestPreviewDeploymentDocCoversLocalTileServer(t *testing.T) {
	data, err := os.ReadFile(repoFile("docs", "race-preview-deployment.md"))
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

func TestPreviewDeploymentDocCoversProductionHosting(t *testing.T) {
	data, err := os.ReadFile(repoFile("docs", "race-preview-deployment.md"))
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

func TestPreviewDeploymentDocCoversTileURLContract(t *testing.T) {
	data, err := os.ReadFile(repoFile("docs", "race-preview-deployment.md"))
	if err != nil {
		t.Fatalf("failed to read docs/race-preview-deployment.md: %v", err)
	}
	text := string(data)

	assertContains(t, text, "{z}/{x}/{y}")
}

func TestPreviewDeploymentDocCoversPublishedMapArtifacts(t *testing.T) {
	data, err := os.ReadFile(repoFile("docs", "race-preview-deployment.md"))
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

func TestPreviewDeploymentDocCoversFutureTileGenerationPipeline(t *testing.T) {
	data, err := os.ReadFile(repoFile("docs", "race-preview-deployment.md"))
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

func TestPreviewDeploymentDocCoversFutureRaceDataImport(t *testing.T) {
	data, err := os.ReadFile(repoFile("docs", "race-preview-deployment.md"))
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

func TestMBTilesCatalogExists(t *testing.T) {
	_, err := os.Stat(repoFile("tiles", "mbtiles", "combined_socal.mbtiles"))
	if err != nil {
		t.Fatalf("tiles/mbtiles/combined_socal.mbtiles must exist for local preview: %v", err)
	}
}

func TestMartinConfigExists(t *testing.T) {
	_, err := os.Stat(repoFile("tiles", "martin-config"))
	if err != nil {
		t.Fatalf("tiles/martin-config must exist for local tile server preview: %v", err)
	}
}
