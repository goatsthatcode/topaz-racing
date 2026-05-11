package topazracing

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

type courseFile struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Elements []struct {
		ID       string  `json:"id"`
		Type     string  `json:"type"`
		Lat      float64 `json:"lat"`
		Lon      float64 `json:"lon"`
		Rounding string  `json:"rounding"`
	} `json:"elements"`
}

type boatsFile struct {
	Boats []struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Color    string `json:"color"`
		BoatType string `json:"boatType"`
		Source   string `json:"source"`
		IsSelf   bool   `json:"isSelf"`
		Track    []struct {
			Time string  `json:"time"`
			Lat  float64 `json:"lat"`
			Lon  float64 `json:"lon"`
		} `json:"track"`
	} `json:"boats"`
}

type eventsFile struct {
	Events []struct {
		ID    string  `json:"id"`
		Type  string  `json:"type"`
		Time  string  `json:"time"`
		Lat   float64 `json:"lat"`
		Lon   float64 `json:"lon"`
		Label string  `json:"label"`
	} `json:"events"`
}

func TestReferenceRaceBundleHasCanonicalFiles(t *testing.T) {
	bundleDir := filepath.Join("content", "races", "dan-byrne-2025", "bishop-rock-race")

	requiredFiles := []string{
		"index.md",
		"course.json",
		"boats.json",
		"events.json",
	}

	for _, name := range requiredFiles {
		path := filepath.Join(bundleDir, name)
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected %s to exist: %v", path, err)
		}
	}
}

func TestReferenceRaceDataParsesAndIncludesCoreEntities(t *testing.T) {
	bundleDir := filepath.Join("content", "races", "dan-byrne-2025", "bishop-rock-race")

	var course courseFile
	readJSONFixture(t, filepath.Join(bundleDir, "course.json"), &course)
	if course.ID == "" || course.Name == "" {
		t.Fatal("expected course metadata to be populated")
	}
	if len(course.Elements) < 3 {
		t.Fatalf("expected at least 3 course elements, got %d", len(course.Elements))
	}

	var boats boatsFile
	readJSONFixture(t, filepath.Join(bundleDir, "boats.json"), &boats)
	if len(boats.Boats) < 2 {
		t.Fatalf("expected at least 2 boats, got %d", len(boats.Boats))
	}
	selfBoats := 0
	for _, boat := range boats.Boats {
		if boat.ID == "" || boat.Name == "" || boat.Color == "" || boat.BoatType == "" || boat.Source == "" {
			t.Fatal("expected required boat metadata fields to be populated")
		}
		if len(boat.Track) < 2 {
			t.Fatal("expected each boat track to contain multiple points")
		}
		if boat.IsSelf {
			selfBoats++
		}
	}
	if selfBoats != 1 {
		t.Fatalf("expected exactly one self boat, got %d", selfBoats)
	}

	var events eventsFile
	readJSONFixture(t, filepath.Join(bundleDir, "events.json"), &events)
	if len(events.Events) < 1 {
		t.Fatal("expected at least one event annotation")
	}
	if events.Events[0].ID == "" || events.Events[0].Type == "" || events.Events[0].Time == "" {
		t.Fatal("expected required event fields to be populated")
	}
}

func readJSONFixture(t *testing.T, path string, target any) {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read %s: %v", path, err)
	}
	if err := json.Unmarshal(data, target); err != nil {
		t.Fatalf("failed to parse %s: %v", path, err)
	}
}
