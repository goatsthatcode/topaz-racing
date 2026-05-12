package gpximport

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// gpx11 is a minimal GPX 1.1 document with two track points.
const gpx11 = `<?xml version="1.0" encoding="UTF-8"?>
<gpx version="1.1" creator="Garmin" xmlns="http://www.topografix.com/GPX/1/1">
  <trk>
    <name>Bishop Rock Race 2025</name>
    <trkseg>
      <trkpt lat="34.4086" lon="-119.6932">
        <time>2025-02-11T18:00:00Z</time>
      </trkpt>
      <trkpt lat="34.4310" lon="-119.8940">
        <time>2025-02-11T19:15:00Z</time>
      </trkpt>
      <trkpt lat="34.4562" lon="-120.1198">
        <time>2025-02-11T20:40:00Z</time>
      </trkpt>
    </trkseg>
  </trk>
</gpx>`

// gpx10 uses no XML namespace (GPX 1.0 style).
const gpx10 = `<?xml version="1.0"?>
<gpx version="1.0" creator="iNavX">
  <trk>
    <name>Leg One</name>
    <trkseg>
      <trkpt lat="33.7000" lon="-118.3000">
        <time>2025-04-05T10:00:00Z</time>
      </trkpt>
      <trkpt lat="33.7500" lon="-118.3500">
        <time>2025-04-05T11:30:00Z</time>
      </trkpt>
    </trkseg>
  </trk>
</gpx>`

// multiSegment has two segments and two tracks to verify concatenation.
const multiSegment = `<?xml version="1.0"?>
<gpx version="1.1" xmlns="http://www.topografix.com/GPX/1/1">
  <trk>
    <name>First Leg</name>
    <trkseg>
      <trkpt lat="34.0" lon="-119.0"><time>2025-01-01T08:00:00Z</time></trkpt>
      <trkpt lat="34.1" lon="-119.1"><time>2025-01-01T09:00:00Z</time></trkpt>
    </trkseg>
    <trkseg>
      <trkpt lat="34.2" lon="-119.2"><time>2025-01-01T10:00:00Z</time></trkpt>
    </trkseg>
  </trk>
  <trk>
    <name>Second Leg</name>
    <trkseg>
      <trkpt lat="34.3" lon="-119.3"><time>2025-01-01T11:00:00Z</time></trkpt>
    </trkseg>
  </trk>
</gpx>`

func TestConvertGPX_BasicParsing(t *testing.T) {
	boat, err := ConvertGPX(strings.NewReader(gpx11), BoatOptions{
		ID:       "topaz",
		Name:     "Topaz",
		Color:    "#4fd1ff",
		BoatType: "Express 27",
		IsSelf:   true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if boat.ID != "topaz" {
		t.Errorf("ID: got %q, want %q", boat.ID, "topaz")
	}
	if boat.Name != "Topaz" {
		t.Errorf("Name: got %q, want %q", boat.Name, "Topaz")
	}
	if boat.Color != "#4fd1ff" {
		t.Errorf("Color: got %q, want %q", boat.Color, "#4fd1ff")
	}
	if boat.BoatType != "Express 27" {
		t.Errorf("BoatType: got %q, want %q", boat.BoatType, "Express 27")
	}
	if !boat.IsSelf {
		t.Error("IsSelf: got false, want true")
	}
	if boat.Source != "gpx" {
		t.Errorf("Source: got %q, want %q", boat.Source, "gpx")
	}

	if len(boat.Track) != 3 {
		t.Fatalf("Track length: got %d, want 3", len(boat.Track))
	}

	first := boat.Track[0]
	if first.Lat != 34.4086 || first.Lon != -119.6932 {
		t.Errorf("first point: got (%v,%v), want (34.4086,-119.6932)", first.Lat, first.Lon)
	}
	wantTime := time.Date(2025, 2, 11, 18, 0, 0, 0, time.UTC)
	if !first.Time.Equal(wantTime) {
		t.Errorf("first point time: got %v, want %v", first.Time, wantTime)
	}
}

func TestConvertGPX_UsesTrackNameWhenNameOptionOmitted(t *testing.T) {
	boat, err := ConvertGPX(strings.NewReader(gpx11), BoatOptions{ID: "topaz"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if boat.Name != "Bishop Rock Race 2025" {
		t.Errorf("Name: got %q, want GPX track name", boat.Name)
	}
}

func TestConvertGPX_UsesIDWhenNoNameAvailable(t *testing.T) {
	noName := `<?xml version="1.0"?>
<gpx version="1.1" xmlns="http://www.topografix.com/GPX/1/1">
  <trk>
    <trkseg>
      <trkpt lat="34.0" lon="-119.0"><time>2025-01-01T08:00:00Z</time></trkpt>
      <trkpt lat="34.1" lon="-119.1"><time>2025-01-01T09:00:00Z</time></trkpt>
    </trkseg>
  </trk>
</gpx>`
	boat, err := ConvertGPX(strings.NewReader(noName), BoatOptions{ID: "my-boat"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if boat.Name != "my-boat" {
		t.Errorf("Name fallback: got %q, want %q", boat.Name, "my-boat")
	}
}

func TestConvertGPX_DefaultColor(t *testing.T) {
	boat, err := ConvertGPX(strings.NewReader(gpx10), BoatOptions{ID: "x", Color: ""})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if boat.Color != "#4fd1ff" {
		t.Errorf("default color: got %q, want #4fd1ff", boat.Color)
	}
}

func TestConvertGPX_DefaultBoatType(t *testing.T) {
	boat, err := ConvertGPX(strings.NewReader(gpx10), BoatOptions{ID: "x"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if boat.BoatType != "unknown" {
		t.Errorf("default boatType: got %q, want unknown", boat.BoatType)
	}
}

func TestConvertGPX_GPX10Format(t *testing.T) {
	boat, err := ConvertGPX(strings.NewReader(gpx10), BoatOptions{ID: "inavx-boat"})
	if err != nil {
		t.Fatalf("GPX 1.0 parse error: %v", err)
	}
	if len(boat.Track) != 2 {
		t.Fatalf("Track length: got %d, want 2", len(boat.Track))
	}
}

func TestConvertGPX_MultiSegmentAndMultiTrackConcatenation(t *testing.T) {
	boat, err := ConvertGPX(strings.NewReader(multiSegment), BoatOptions{ID: "fleet"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(boat.Track) != 4 {
		t.Fatalf("Track length: got %d, want 4 (2 segs + second trk)", len(boat.Track))
	}
}

func TestConvertGPX_TooFewPoints(t *testing.T) {
	onePoint := `<?xml version="1.0"?>
<gpx version="1.1" xmlns="http://www.topografix.com/GPX/1/1">
  <trk><trkseg>
    <trkpt lat="34.0" lon="-119.0"><time>2025-01-01T08:00:00Z</time></trkpt>
  </trkseg></trk>
</gpx>`
	_, err := ConvertGPX(strings.NewReader(onePoint), BoatOptions{ID: "x"})
	if err == nil {
		t.Error("expected error for 1 track point, got nil")
	}
}

func TestConvertGPX_InvalidXML(t *testing.T) {
	_, err := ConvertGPX(strings.NewReader("<not valid xml"), BoatOptions{ID: "x"})
	if err == nil {
		t.Error("expected error for invalid XML, got nil")
	}
}

func TestConvertGPX_BadTimestamp(t *testing.T) {
	badTime := `<?xml version="1.0"?>
<gpx version="1.1" xmlns="http://www.topografix.com/GPX/1/1">
  <trk><trkseg>
    <trkpt lat="34.0" lon="-119.0"><time>not-a-date</time></trkpt>
    <trkpt lat="34.1" lon="-119.1"><time>also-not-a-date</time></trkpt>
  </trkseg></trk>
</gpx>`
	_, err := ConvertGPX(strings.NewReader(badTime), BoatOptions{ID: "x"})
	if err == nil {
		t.Error("expected error for unparseable timestamp, got nil")
	}
}

func TestConvertGPX_JSONOutputMatchesV1Schema(t *testing.T) {
	boat, err := ConvertGPX(strings.NewReader(gpx11), BoatOptions{
		ID:       "topaz",
		Name:     "Topaz",
		Color:    "#4fd1ff",
		BoatType: "Express 27",
		IsSelf:   true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	file := BoatsFile{Boats: []Boat{*boat}}
	data, err := json.Marshal(file)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	// Verify required V1 schema fields are present and properly typed.
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	boats, ok := raw["boats"].([]interface{})
	if !ok || len(boats) == 0 {
		t.Fatal("missing or empty boats array in JSON output")
	}

	b := boats[0].(map[string]interface{})
	for _, field := range []string{"id", "name", "color", "boatType", "source", "isSelf", "track"} {
		if _, exists := b[field]; !exists {
			t.Errorf("V1 required field %q missing from JSON output", field)
		}
	}

	track, ok := b["track"].([]interface{})
	if !ok {
		t.Fatal("track is not an array")
	}
	pt := track[0].(map[string]interface{})
	for _, field := range []string{"time", "lat", "lon"} {
		if _, exists := pt[field]; !exists {
			t.Errorf("V1 track point field %q missing from JSON output", field)
		}
	}
	// time must be a string (RFC3339).
	if _, ok := pt["time"].(string); !ok {
		t.Error("track point 'time' should be a string")
	}
}

func TestMergeBoat_AppendNewBoat(t *testing.T) {
	file := &BoatsFile{
		Boats: []Boat{{ID: "existing", Name: "Existing", Color: "#fff", BoatType: "x", Source: "hand-authored", Track: []TrackPoint{}}},
	}
	newBoat := &Boat{ID: "new", Name: "New"}
	MergeBoat(file, newBoat)

	if len(file.Boats) != 2 {
		t.Fatalf("expected 2 boats after merge, got %d", len(file.Boats))
	}
	if file.Boats[1].ID != "new" {
		t.Errorf("appended boat ID: got %q, want new", file.Boats[1].ID)
	}
}

func TestMergeBoat_ReplacesExistingBoat(t *testing.T) {
	file := &BoatsFile{
		Boats: []Boat{
			{ID: "topaz", Name: "Old Name"},
			{ID: "wildcard", Name: "Wildcard"},
		},
	}
	updated := &Boat{ID: "topaz", Name: "New Name"}
	MergeBoat(file, updated)

	if len(file.Boats) != 2 {
		t.Fatalf("expected 2 boats after replace, got %d", len(file.Boats))
	}
	if file.Boats[0].Name != "New Name" {
		t.Errorf("replaced boat name: got %q, want New Name", file.Boats[0].Name)
	}
}

func TestMergeBoat_EmptyFile(t *testing.T) {
	file := &BoatsFile{}
	boat := &Boat{ID: "solo"}
	MergeBoat(file, boat)

	if len(file.Boats) != 1 {
		t.Fatalf("expected 1 boat, got %d", len(file.Boats))
	}
}
