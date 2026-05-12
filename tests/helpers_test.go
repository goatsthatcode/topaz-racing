package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// repoRoot is the path from the tests/ package directory to the project root.
const repoRoot = ".."

// ─── Shared data types ────────────────────────────────────────────────────────

type raceCourseSchemaDocument struct {
	Required []string `json:"required"`
	Defs     struct {
		CourseElement struct {
			Required   []string `json:"required"`
			Properties struct {
				Type struct {
					Enum []string `json:"enum"`
				} `json:"type"`
				Rounding struct {
					Enum []string `json:"enum"`
				} `json:"rounding"`
				ControlPointsToNext json.RawMessage `json:"controlPointsToNext"`
			} `json:"properties"`
		} `json:"courseElement"`
	} `json:"$defs"`
}

type raceBoatsSchemaDocument struct {
	Required []string `json:"required"`
	Defs     struct {
		Boat struct {
			Required []string `json:"required"`
		} `json:"boat"`
		TrackPoint struct {
			Required   []string `json:"required"`
			Properties struct {
				Time struct {
					Format string `json:"format"`
				} `json:"time"`
			} `json:"properties"`
		} `json:"trackPoint"`
	} `json:"$defs"`
}

type raceEventsSchemaDocument struct {
	Required []string `json:"required"`
	Defs     struct {
		Event struct {
			Required   []string          `json:"required"`
			AnyOf      []json.RawMessage `json:"anyOf"`
			Properties struct {
				Type struct {
					Enum []string `json:"enum"`
				} `json:"type"`
				Time struct {
					Format string `json:"format"`
				} `json:"time"`
				Label       json.RawMessage `json:"label"`
				Description json.RawMessage `json:"description"`
			} `json:"properties"`
		} `json:"event"`
	} `json:"$defs"`
}

type raceCourseFile struct {
	ID       string              `json:"id"`
	Name     string              `json:"name"`
	Elements []raceCourseElement `json:"elements"`
}

type raceCourseElement struct {
	ID                  string           `json:"id"`
	Type                string           `json:"type"`
	Lat                 float64          `json:"lat"`
	Lon                 float64          `json:"lon"`
	Name                string           `json:"name"`
	Rounding            string           `json:"rounding"`
	ControlPointsToNext []raceCoordinate `json:"controlPointsToNext"`
}

type raceCoordinate struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type raceBoatsFile struct {
	Boats []raceBoat `json:"boats"`
}

type raceBoat struct {
	ID       string           `json:"id"`
	Name     string           `json:"name"`
	Color    string           `json:"color"`
	BoatType string           `json:"boatType"`
	Source   string           `json:"source"`
	IsSelf   bool             `json:"isSelf"`
	Track    []raceTrackPoint `json:"track"`
}

type raceTrackPoint struct {
	Time string  `json:"time"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

type raceEventsFile struct {
	Events []raceEvent `json:"events"`
}

type raceEvent struct {
	ID          string   `json:"id"`
	Type        string   `json:"type"`
	Time        string   `json:"time"`
	Lat         *float64 `json:"lat"`
	Lon         *float64 `json:"lon"`
	Label       string   `json:"label"`
	Description string   `json:"description"`
}

// courseFile, boatsFile, eventsFile are simpler anonymous-struct types used by
// the content structure tests that only need a subset of the full field set.
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

// ─── Shared helpers ───────────────────────────────────────────────────────────

func readBuiltHTML(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read built HTML %s: %v", path, err)
	}

	return string(data)
}

func readBuiltFile(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read built file %s: %v", path, err)
	}

	return strings.TrimSpace(string(data))
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

func assertContains(t *testing.T, text, snippet string) {
	t.Helper()

	if !strings.Contains(text, snippet) {
		t.Fatalf("expected content to include %q", snippet)
	}
}

func assertStringSet(t *testing.T, actual, expected []string) {
	t.Helper()

	if len(actual) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}

	for _, value := range expected {
		if !slices.Contains(actual, value) {
			t.Fatalf("expected %q in %v", value, actual)
		}
	}
}

func assertCoordinateInRange(t *testing.T, lat, lon float64) {
	t.Helper()

	if lat < -90 || lat > 90 {
		t.Fatalf("latitude %f out of range", lat)
	}
	if lon < -180 || lon > 180 {
		t.Fatalf("longitude %f out of range", lon)
	}
}

func expandCourseRouteCoordinates(course raceCourseFile) []raceCoordinate {
	coordinates := make([]raceCoordinate, 0, len(course.Elements))

	for i, element := range course.Elements {
		coordinates = append(coordinates, raceCoordinate{
			Lat: element.Lat,
			Lon: element.Lon,
		})

		if i == len(course.Elements)-1 {
			continue
		}

		coordinates = append(coordinates, element.ControlPointsToNext...)
	}

	return coordinates
}

// repoFile returns a path rooted at the project root so tests can reference
// source files without hard-coding absolute paths.
func repoFile(parts ...string) string {
	return filepath.Join(append([]string{repoRoot}, parts...)...)
}
