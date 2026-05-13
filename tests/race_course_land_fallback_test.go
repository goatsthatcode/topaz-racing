package tests

import (
	"testing"
)

func TestReferenceRaceCourseUsesStraightRouteSegment(t *testing.T) {
	var course raceCourseFile
	readJSONFixture(
		t,
		repoFile("content", "races", "dan-byrne-2025", "bishop-rock-race", "course.json"),
		&course,
	)

	expectedCoordinates := []raceCoordinate{
		{Lat: 33.95916, Lon: -118.46722},
		{Lat: 33.91147, Lon: -118.4596},
		{Lat: 32.43059, Lon: -119.12433},
		{Lat: 33.48193, Lon: -118.61635},
	}

	actualCoordinates := expandCourseRouteCoordinates(course)
	if len(actualCoordinates) != len(expectedCoordinates) {
		t.Fatalf("expected %d route coordinates, got %d", len(expectedCoordinates), len(actualCoordinates))
	}

	for i, expected := range expectedCoordinates {
		actual := actualCoordinates[i]
		if actual != expected {
			t.Fatalf("expected route coordinate %d to be %+v, got %+v", i, expected, actual)
		}
	}

	for i, element := range course.Elements {
		if len(element.ControlPointsToNext) != 0 {
			t.Fatalf("expected element %d (%s) to have no shaping points, got %d", i, element.ID, len(element.ControlPointsToNext))
		}
	}
}

func TestCourseRouteExpansionIncludesManualShapingPoints(t *testing.T) {
	course := raceCourseFile{
		ID:   "manual-shaping-test",
		Name: "Manual Shaping Test",
		Elements: []raceCourseElement{
			{
				ID:       "start",
				Type:     "start_line",
				Lat:      33.0,
				Lon:      -118.0,
				Rounding: "none",
				ControlPointsToNext: []raceCoordinate{
					{Lat: 33.1, Lon: -118.1},
					{Lat: 33.2, Lon: -118.2},
				},
			},
			{
				ID:       "mark",
				Type:     "mark",
				Lat:      33.3,
				Lon:      -118.3,
				Rounding: "port",
			},
			{
				ID:       "finish",
				Type:     "finish_line",
				Lat:      33.4,
				Lon:      -118.4,
				Rounding: "none",
			},
		},
	}

	expectedCoordinates := []raceCoordinate{
		{Lat: 33.0, Lon: -118.0},
		{Lat: 33.1, Lon: -118.1},
		{Lat: 33.2, Lon: -118.2},
		{Lat: 33.3, Lon: -118.3},
		{Lat: 33.4, Lon: -118.4},
	}

	actualCoordinates := expandCourseRouteCoordinates(course)
	if len(actualCoordinates) != len(expectedCoordinates) {
		t.Fatalf("expected %d route coordinates, got %d", len(expectedCoordinates), len(actualCoordinates))
	}

	for i, expected := range expectedCoordinates {
		actual := actualCoordinates[i]
		if actual != expected {
			t.Fatalf("expected route coordinate %d to be %+v, got %+v", i, expected, actual)
		}
	}
}
