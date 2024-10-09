package main

import (
	//"testing"
	"fmt"

	geo "github.com/kellydunn/golang-geo"
)

func NewGeofence(points [][]*geo.Point, args ...interface{}) *Geofence {

	points = [][]*geo.Point{
		{
			geo.NewPoint(14.068165238581306, 121.32714315177356),
			geo.NewPoint(14.068096291876111, 121.32722764135187),
			geo.NewPoint(14.067866035747677, 121.32688163641218),
			geo.NewPoint(14.067797739394011, 121.32696612599041),
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
			geo.NewPoint(14.068165238581306, 121.32714315177356),
			geo.NewPoint(14.068096291876111, 121.32722764135187),
			geo.NewPoint(14.067866035747677, 121.32688163641218),
			geo.NewPoint(14.067797739394011, 121.32696612599041),
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
