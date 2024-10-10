package main

import (
	//"testing"
	"fmt"

	geo "github.com/kellydunn/golang-geo"
)

func NewGeofence(points [][]*geo.Point, args ...interface{}) *Geofence {

	points = [][]*geo.Point{
		{
			geo.NewPoint(14.067723824810402, 121.32709045412432),
			geo.NewPoint(14.068017824302677, 121.32736605108208),
			geo.NewPoint(14.068111487778342, 121.32721718849176),
			geo.NewPoint(14.067799926485147, 121.32696640196573),
		},
	}

	geofence := &Geofence{}
	geofence.tiles = make(map[float64]string)
	if len(args) > 0 {
		geofence.granularity = args[0].(int64)
	} else {
		geofence.granularity = defaultGranularity
	}
	geofence.vertices = points[(0)]
	if len(points) > 1 {
		geofence.holes = points[1:]
		geofence.tiles = make(map[float64]string)
	}
	geofence.setInclusionTiles()
	return geofence
}

func main() {
	fmt.Print("Enter latitude: ")
	var lat float64
	fmt.Scanln(&lat)

	fmt.Print("Enter longitude: ")
	var lng float64
	fmt.Scanln(&lng)

	points := [][]*geo.Point{
		{
			geo.NewPoint(14.067723824810402, 121.32709045412432),
			geo.NewPoint(14.068017824302677, 121.32736605108208),
			geo.NewPoint(14.068111487778342, 121.32721718849176),
			geo.NewPoint(14.067799926485147, 121.32696640196573),
		},
	}

	geofence := NewGeofence(points)
	point := geo.NewPoint(lat, lng)
	if geofence.Inside(point) {
		fmt.Println("Point is inside the geofence")
	} else {
		fmt.Println("Point is outside the geofence")
	}
}
