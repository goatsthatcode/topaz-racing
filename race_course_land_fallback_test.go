package topazracing

import (
	"path/filepath"
	"testing"
)

func TestReferenceRaceCourseUsesStraightRouteSegment(t *testing.T) {
	var course raceCourseFile
	readJSONFixture(
		t,
		filepath.Join("content", "races", "dan-byrne-2025", "bishop-rock-race", "course.json"),
		&course,
	)

	expectedCoordinates := []raceCoordinate{
		{Lat: 33.9769, Lon: -118.4451},
		{Lat: 32.475, Lon: -119.293},
		{Lat: 33.4445, Lon: -118.6074},
		{Lat: 33.4445, Lon: -118.6074},
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

	if len(course.Elements[1].ControlPointsToNext) != 0 {
		t.Fatalf("expected Cortez Bank leg to include no shaping points, got %d", len(course.Elements[1].ControlPointsToNext))
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
