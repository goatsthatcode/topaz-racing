// Package jibesetimport parses Jibeset multi-boat track export files into the
// Topaz Racing V1 boats.json format.
//
// The Jibeset export format is a plain-text file using carriage returns (\r)
// as line terminators. It contains one or more track blocks delimited by:
//
//	// Start of track SAILNUM NAME
//	...header comment lines...
//	YYYY-MM-DD_HH:MM:SS_LAT_LON
//	...
//	// End of track SAILNUM NAME
//
// Track timestamps are UTC. Header lines beginning with //Z_FINH carry the
// official race start and finish times for each boat.
package jibesetimport

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"goatsthatcode.github.io/topaz-racing/gpximport"
)

// RawTrack holds the parsed contents of a single track block.
type RawTrack struct {
	SailNumber string
	Name       string
	// StartTime is the official race start time from the FINH header.
	StartTime time.Time
	// FinishTime is the declared finish time from the FINH header, if present.
	FinishTime time.Time
	// Points contains every track point found in the block, in file order.
	Points []gpximport.TrackPoint
}

// ParseFile reads a Jibeset export and returns one RawTrack per block.
// The reader may use either \r or \n (or both) as line terminators.
func ParseFile(r io.Reader) ([]RawTrack, error) {
	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading jibeset input: %w", err)
	}

	// Normalise to \n so we can split uniformly regardless of source encoding.
	normalised := bytes.ReplaceAll(raw, []byte("\r\n"), []byte("\n"))
	normalised = bytes.ReplaceAll(normalised, []byte("\r"), []byte("\n"))
	lines := strings.Split(string(normalised), "\n")

	var tracks []RawTrack
	var current *RawTrack

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, "// Start of track "):
			rest := strings.TrimPrefix(line, "// Start of track ")
			parts := strings.SplitN(rest, " ", 2)
			sailNum := parts[0]
			name := ""
			if len(parts) == 2 {
				name = strings.TrimSpace(parts[1])
			}
			current = &RawTrack{SailNumber: sailNum, Name: name}

		case strings.HasPrefix(line, "// End of track "):
			if current != nil {
				tracks = append(tracks, *current)
				current = nil
			}

		case strings.HasPrefix(line, "//Z_FINH_"):
			// //Z_FINH_Start:_YYYY-MM-DD_HH:MM:SS_unix_Finish:_YYYY-MM-DD_HH:MM:SS_unix
			if current == nil {
				continue
			}
			st, ft, parseErr := parseFINHHeader(line)
			if parseErr == nil {
				current.StartTime = st
				current.FinishTime = ft
			}

		case strings.HasPrefix(line, "//"):
			// Other comment lines — skip.

		default:
			if current == nil {
				continue
			}
			pt, parseErr := parseTrackPoint(line)
			if parseErr != nil {
				continue // tolerate malformed lines
			}
			current.Points = append(current.Points, pt)
		}
	}

	// If the file ended without an explicit End marker, keep the last block.
	if current != nil && len(current.Points) > 0 {
		tracks = append(tracks, *current)
	}

	if len(tracks) == 0 {
		return nil, fmt.Errorf("no track blocks found in jibeset file")
	}

	return tracks, nil
}

// ConvertTrack converts a RawTrack into a Boat using the provided metadata.
// minIntervalSecs specifies the minimum seconds between consecutive output
// points; pass 0 to include every parsed point.
// The boat ID is derived from the sail number if opts.ID is empty.
func ConvertTrack(track RawTrack, opts gpximport.BoatOptions, minIntervalSecs int) (*gpximport.Boat, error) {
	if len(track.Points) < 2 {
		return nil, fmt.Errorf("track %q has %d point(s); at least 2 are required",
			track.SailNumber, len(track.Points))
	}

	points := downsample(track.Points, minIntervalSecs)
	if len(points) < 2 {
		return nil, fmt.Errorf("track %q has fewer than 2 points after downsampling to %ds interval",
			track.SailNumber, minIntervalSecs)
	}

	id := opts.ID
	if id == "" {
		id = sailToID(track.SailNumber)
	}

	name := opts.Name
	if name == "" {
		name = track.Name
	}
	if name == "" {
		name = id
	}

	color := opts.Color
	if color == "" {
		color = "#4fd1ff"
	}

	boatType := opts.BoatType
	if boatType == "" {
		boatType = "unknown"
	}

	return &gpximport.Boat{
		ID:       id,
		Name:     name,
		Color:    color,
		BoatType: boatType,
		Source:   "jibeset",
		IsSelf:   opts.IsSelf,
		Track:    points,
	}, nil
}

// downsample returns a subset of pts keeping only points where at least
// minIntervalSecs seconds have elapsed since the previous kept point.
// If minIntervalSecs <= 0 all points are returned unchanged.
func downsample(pts []gpximport.TrackPoint, minIntervalSecs int) []gpximport.TrackPoint {
	if minIntervalSecs <= 0 || len(pts) == 0 {
		return pts
	}

	interval := time.Duration(minIntervalSecs) * time.Second
	kept := []gpximport.TrackPoint{pts[0]}
	lastKept := pts[0].Time

	for _, pt := range pts[1:] {
		if pt.Time.Sub(lastKept) >= interval {
			kept = append(kept, pt)
			lastKept = pt.Time
		}
	}

	// Always include the final point so the track ends at the true finish.
	last := pts[len(pts)-1]
	if !kept[len(kept)-1].Time.Equal(last.Time) {
		kept = append(kept, last)
	}

	return kept
}

// sailToID converts a sail number string to a lowercase identifier suitable
// for use as a JSON boats.json boat ID.
func sailToID(sailNum string) string {
	return strings.ToLower(strings.ReplaceAll(sailNum, " ", "-"))
}

// parseTrackPoint parses a line of the form YYYY-MM-DD_HH:MM:SS_LAT_LON.
func parseTrackPoint(line string) (gpximport.TrackPoint, error) {
	parts := strings.Split(line, "_")
	if len(parts) != 4 {
		return gpximport.TrackPoint{}, fmt.Errorf("expected 4 underscore-separated fields, got %d: %q", len(parts), line)
	}

	t, err := time.Parse("2006-01-02T15:04:05Z", parts[0]+"T"+parts[1]+"Z")
	if err != nil {
		return gpximport.TrackPoint{}, fmt.Errorf("parsing timestamp %q %q: %w", parts[0], parts[1], err)
	}

	lat, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return gpximport.TrackPoint{}, fmt.Errorf("parsing lat %q: %w", parts[2], err)
	}

	lon, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return gpximport.TrackPoint{}, fmt.Errorf("parsing lon %q: %w", parts[3], err)
	}

	return gpximport.TrackPoint{Time: t, Lat: lat, Lon: lon}, nil
}

// parseFINHHeader extracts both the Start and Finish timestamps from a
// //Z_FINH_ header line.
// Example: //Z_FINH_Start:_2025-03-14_10:50:00_1741963800_Finish:_2025-03-17_03:05:23_1742195123
func parseFINHHeader(line string) (startTime, finishTime time.Time, err error) {
	const startMarker = "Start:_"
	si := strings.Index(line, startMarker)
	if si < 0 {
		return time.Time{}, time.Time{}, fmt.Errorf("no Start: field in %q", line)
	}
	sParts := strings.Split(line[si+len(startMarker):], "_")
	if len(sParts) < 2 {
		return time.Time{}, time.Time{}, fmt.Errorf("unexpected start format in %q", line)
	}
	startTime, err = time.Parse("2006-01-02T15:04:05Z", sParts[0]+"T"+sParts[1]+"Z")
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("parsing start time: %w", err)
	}

	const finishMarker = "Finish:_"
	fi := strings.Index(line, finishMarker)
	if fi < 0 {
		return startTime, time.Time{}, nil // finish may be absent for DNF
	}
	fParts := strings.Split(line[fi+len(finishMarker):], "_")
	if len(fParts) < 2 {
		return startTime, time.Time{}, nil
	}
	finishTime, err = time.Parse("2006-01-02T15:04:05Z", fParts[0]+"T"+fParts[1]+"Z")
	if err != nil {
		return startTime, time.Time{}, nil // non-fatal
	}
	return startTime, finishTime, nil
}
