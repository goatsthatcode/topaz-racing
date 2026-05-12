// Package gpximport converts GPX track files to the Topaz Racing V1 boats.json format.
package gpximport

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"time"
)

// BoatOptions carries metadata for the boat entry produced from a GPX file.
type BoatOptions struct {
	ID       string
	Name     string
	Color    string
	BoatType string
	IsSelf   bool
}

// TrackPoint is a single time-stamped position in a boat's track.
type TrackPoint struct {
	Time time.Time `json:"time"`
	Lat  float64   `json:"lat"`
	Lon  float64   `json:"lon"`
}

// MarshalJSON formats time as RFC3339 UTC string.
func (p TrackPoint) MarshalJSON() ([]byte, error) {
	type wire struct {
		Time string  `json:"time"`
		Lat  float64 `json:"lat"`
		Lon  float64 `json:"lon"`
	}
	return json.Marshal(wire{
		Time: p.Time.UTC().Format(time.RFC3339),
		Lat:  p.Lat,
		Lon:  p.Lon,
	})
}

// Boat is one entry in boats.json.
type Boat struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Color    string       `json:"color"`
	BoatType string       `json:"boatType"`
	Source   string       `json:"source"`
	IsSelf   bool         `json:"isSelf"`
	Track    []TrackPoint `json:"track"`
}

// BoatsFile is the top-level boats.json container.
type BoatsFile struct {
	Boats []Boat `json:"boats"`
}

// GPX parsing types. Both GPX 1.0 and 1.1 share the same element names.
type gpxDoc struct {
	XMLName xml.Name   `xml:"gpx"`
	Tracks  []gpxTrack `xml:"trk"`
}

type gpxTrack struct {
	Name     string       `xml:"name"`
	Segments []gpxSegment `xml:"trkseg"`
}

type gpxSegment struct {
	Points []gpxPoint `xml:"trkpt"`
}

type gpxPoint struct {
	Lat  float64 `xml:"lat,attr"`
	Lon  float64 `xml:"lon,attr"`
	Time string  `xml:"time"`
}

// ConvertGPX parses a GPX stream and produces a Boat entry using the supplied options.
// All track segments across all tracks in the file are concatenated in order.
// If opts.Name is empty and the GPX file contains a track name, that name is used.
func ConvertGPX(r io.Reader, opts BoatOptions) (*Boat, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading GPX input: %w", err)
	}

	var doc gpxDoc
	if err := xml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("parsing GPX XML: %w", err)
	}

	// Collect all track points across all tracks and segments.
	var points []TrackPoint
	var firstTrackName string
	for _, trk := range doc.Tracks {
		if firstTrackName == "" {
			firstTrackName = trk.Name
		}
		for _, seg := range trk.Segments {
			for _, pt := range seg.Points {
				t, err := parseGPXTime(pt.Time)
				if err != nil {
					return nil, fmt.Errorf("parsing timestamp %q: %w", pt.Time, err)
				}
				points = append(points, TrackPoint{Time: t, Lat: pt.Lat, Lon: pt.Lon})
			}
		}
	}

	if len(points) < 2 {
		return nil, fmt.Errorf("GPX file contains %d track point(s); at least 2 are required", len(points))
	}

	name := opts.Name
	if name == "" {
		name = firstTrackName
	}
	if name == "" {
		name = opts.ID
	}

	color := opts.Color
	if color == "" {
		color = "#4fd1ff"
	}

	boatType := opts.BoatType
	if boatType == "" {
		boatType = "unknown"
	}

	boat := &Boat{
		ID:       opts.ID,
		Name:     name,
		Color:    color,
		BoatType: boatType,
		Source:   "gpx",
		IsSelf:   opts.IsSelf,
		Track:    points,
	}
	return boat, nil
}

// MergeBoat adds or replaces a boat entry in an existing BoatsFile.
// If a boat with the same ID already exists it is replaced; otherwise the new boat is appended.
func MergeBoat(file *BoatsFile, boat *Boat) {
	for i, b := range file.Boats {
		if b.ID == boat.ID {
			file.Boats[i] = *boat
			return
		}
	}
	file.Boats = append(file.Boats, *boat)
}

// gpxTimeFormats lists the timestamp formats commonly found in GPX files.
var gpxTimeFormats = []string{
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02T15:04:05.999999999Z07:00",
}

func parseGPXTime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, fmt.Errorf("empty timestamp")
	}
	for _, layout := range gpxTimeFormats {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognized timestamp format %q", s)
}
