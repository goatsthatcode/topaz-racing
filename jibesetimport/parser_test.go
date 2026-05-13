package jibesetimport

import (
	"strings"
	"testing"
	"time"

	"goatsthatcode.github.io/topaz-racing/gpximport"
)

// minimalJibeset is a two-boat, four-point Jibeset export using \r line terminators.
var minimalJibeset = strings.Join([]string{
	"// Start of track 1234 ALPHA",
	"//DE enabled at: 2025-03-01T00:00:00Z_1234_ALPHA_1740787200_1741046400_share.garmin.com/alpha",
	"//Z_FINH_Start:_2025-03-01_08:00:00_1740808800_Finish:_2025-03-02_10:00:00_1740895200",
	"//Z_ Trackstart:_1740808800_Trackend:_1740895200",
	"2025-03-01_08:00:00_33.900000_-118.500000",
	"2025-03-01_09:00:00_33.800000_-118.600000",
	"2025-03-01_10:00:00_33.700000_-118.700000",
	"// End of track 1234 ALPHA",
	"// Start of track usa99 BETA",
	"//DE enabled at: 2025-03-01T00:00:00Z_usa99_BETA_1740787200_1741046400_share.garmin.com/beta",
	"//Z_FINH_Start:_2025-03-01_08:00:00_1740808800_Finish:_2025-03-02_12:00:00_1740902400",
	"//Z_ Trackstart:_1740808800_Trackend:_1740902400",
	"2025-03-01_08:00:00_33.901000_-118.501000",
	"2025-03-02_12:00:00_33.500000_-118.600000",
	"// End of track usa99 BETA",
}, "\r")

func TestParseFile_TwoBoats(t *testing.T) {
	tracks, err := ParseFile(strings.NewReader(minimalJibeset))
	if err != nil {
		t.Fatalf("ParseFile returned error: %v", err)
	}
	if len(tracks) != 2 {
		t.Fatalf("expected 2 tracks, got %d", len(tracks))
	}
}

func TestParseFile_SailNumberAndName(t *testing.T) {
	tracks, err := ParseFile(strings.NewReader(minimalJibeset))
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if tracks[0].SailNumber != "1234" {
		t.Errorf("track[0].SailNumber = %q, want %q", tracks[0].SailNumber, "1234")
	}
	if tracks[0].Name != "ALPHA" {
		t.Errorf("track[0].Name = %q, want %q", tracks[0].Name, "ALPHA")
	}
	if tracks[1].SailNumber != "usa99" {
		t.Errorf("track[1].SailNumber = %q, want %q", tracks[1].SailNumber, "usa99")
	}
}

func TestParseFile_PointCount(t *testing.T) {
	tracks, err := ParseFile(strings.NewReader(minimalJibeset))
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if len(tracks[0].Points) != 3 {
		t.Errorf("track[0] points = %d, want 3", len(tracks[0].Points))
	}
	if len(tracks[1].Points) != 2 {
		t.Errorf("track[1] points = %d, want 2", len(tracks[1].Points))
	}
}

func TestParseFile_PointValues(t *testing.T) {
	tracks, err := ParseFile(strings.NewReader(minimalJibeset))
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	pt := tracks[0].Points[0]
	wantTime := time.Date(2025, 3, 1, 8, 0, 0, 0, time.UTC)
	if !pt.Time.Equal(wantTime) {
		t.Errorf("point time = %v, want %v", pt.Time, wantTime)
	}
	if pt.Lat != 33.9 {
		t.Errorf("point lat = %v, want 33.9", pt.Lat)
	}
	if pt.Lon != -118.5 {
		t.Errorf("point lon = %v, want -118.5", pt.Lon)
	}
}

func TestParseFile_StartAndFinishTime(t *testing.T) {
	tracks, err := ParseFile(strings.NewReader(minimalJibeset))
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	wantStart := time.Date(2025, 3, 1, 8, 0, 0, 0, time.UTC)
	if !tracks[0].StartTime.Equal(wantStart) {
		t.Errorf("track[0].StartTime = %v, want %v", tracks[0].StartTime, wantStart)
	}
	wantFinish := time.Date(2025, 3, 2, 10, 0, 0, 0, time.UTC)
	if !tracks[0].FinishTime.Equal(wantFinish) {
		t.Errorf("track[0].FinishTime = %v, want %v", tracks[0].FinishTime, wantFinish)
	}
}

func TestParseFile_MultipleDatesDistinctStartTimes(t *testing.T) {
	// The two boats in minimalJibeset have different finish times but same start date.
	tracks, err := ParseFile(strings.NewReader(minimalJibeset))
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	// Both share the same race start date (2025-03-01).
	for _, tr := range tracks {
		if tr.StartTime.Year() != 2025 || tr.StartTime.Month() != 3 || tr.StartTime.Day() != 1 {
			t.Errorf("sail=%s unexpected start date %v", tr.SailNumber, tr.StartTime)
		}
	}
}

func TestParseFile_UnixNewlines(t *testing.T) {
	// Same data but using \n — ParseFile must normalise line terminators.
	data := strings.ReplaceAll(minimalJibeset, "\r", "\n")
	tracks, err := ParseFile(strings.NewReader(data))
	if err != nil {
		t.Fatalf("ParseFile with \\n terminators: %v", err)
	}
	if len(tracks) != 2 {
		t.Errorf("expected 2 tracks with \\n terminators, got %d", len(tracks))
	}
}

func TestConvertTrack_BasicConversion(t *testing.T) {
	tracks, _ := ParseFile(strings.NewReader(minimalJibeset))
	boat, err := ConvertTrack(tracks[0], gpximport.BoatOptions{
		ID:    "alpha",
		Name:  "Alpha",
		Color: "#ff8a5b",
	}, 0)
	if err != nil {
		t.Fatalf("ConvertTrack: %v", err)
	}
	if boat.ID != "alpha" {
		t.Errorf("ID = %q, want alpha", boat.ID)
	}
	if boat.Source != "jibeset" {
		t.Errorf("Source = %q, want jibeset", boat.Source)
	}
	if len(boat.Track) != 3 {
		t.Errorf("track points = %d, want 3", len(boat.Track))
	}
}

func TestConvertTrack_DefaultsFromTrack(t *testing.T) {
	tracks, _ := ParseFile(strings.NewReader(minimalJibeset))
	// Empty opts — ID and Name should fall back to sail number / track name.
	boat, err := ConvertTrack(tracks[0], gpximport.BoatOptions{}, 0)
	if err != nil {
		t.Fatalf("ConvertTrack: %v", err)
	}
	if boat.ID != "1234" {
		t.Errorf("ID = %q, want 1234", boat.ID)
	}
	if boat.Name != "ALPHA" {
		t.Errorf("Name = %q, want ALPHA", boat.Name)
	}
	if boat.Color != "#4fd1ff" {
		t.Errorf("Color = %q, want default #4fd1ff", boat.Color)
	}
}

func TestConvertTrack_IsSelf(t *testing.T) {
	tracks, _ := ParseFile(strings.NewReader(minimalJibeset))
	boat, _ := ConvertTrack(tracks[0], gpximport.BoatOptions{IsSelf: true}, 0)
	if !boat.IsSelf {
		t.Error("IsSelf should be true")
	}
}

func TestConvertTrack_Downsampling(t *testing.T) {
	// ALPHA has 3 points at 1-hour intervals; a 90-minute interval keeps first + last.
	tracks, _ := ParseFile(strings.NewReader(minimalJibeset))
	boat, err := ConvertTrack(tracks[0], gpximport.BoatOptions{}, 5400) // 90 min
	if err != nil {
		t.Fatalf("ConvertTrack: %v", err)
	}
	// Should keep point[0] (08:00), skip point[1] (09:00, only 60 min later),
	// but include point[2] (10:00, 120 min after point[0]) — and always the last point.
	if len(boat.Track) != 2 {
		t.Errorf("after 90min downsampling: %d points, want 2", len(boat.Track))
	}
	// Last point must be the actual last track point.
	want := time.Date(2025, 3, 1, 10, 0, 0, 0, time.UTC)
	if !boat.Track[len(boat.Track)-1].Time.Equal(want) {
		t.Errorf("last point time = %v, want %v", boat.Track[len(boat.Track)-1].Time, want)
	}
}

func TestConvertTrack_TooFewPoints(t *testing.T) {
	raw := strings.Join([]string{
		"// Start of track 9 SOLO",
		"2025-01-01_00:00:00_33.0_-118.0",
		"// End of track 9 SOLO",
	}, "\r")
	tracks, _ := ParseFile(strings.NewReader(raw))
	_, err := ConvertTrack(tracks[0], gpximport.BoatOptions{}, 0)
	if err == nil {
		t.Error("expected error for single-point track, got nil")
	}
}

func TestParseFile_EmptyFile(t *testing.T) {
	_, err := ParseFile(strings.NewReader(""))
	if err == nil {
		t.Error("expected error for empty file, got nil")
	}
}

func TestParseFile_MalformedPointsSkipped(t *testing.T) {
	raw := strings.Join([]string{
		"// Start of track 5 VALID",
		"not_a_valid_line",
		"2025-01-01_12:00:00_33.5_-118.5",
		"2025-01-01_13:00:00_33.6_-118.6",
		"// End of track 5 VALID",
	}, "\r")
	tracks, err := ParseFile(strings.NewReader(raw))
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if len(tracks[0].Points) != 2 {
		t.Errorf("expected 2 valid points (malformed skipped), got %d", len(tracks[0].Points))
	}
}
