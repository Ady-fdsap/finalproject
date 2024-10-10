package main

import (
	//"testing"
	"fmt"

	geo "github.com/kellydunn/golang-geo"
)

func NewGeofence(points [][]*geo.Point, args ...interface{}) *Geofence {

	points = [][]*geo.Point{
		{
			geo.NewPoint(14.067694194798804, 121.32708640042505),
			geo.NewPoint(14.06800445535538, 121.32742234709286),
			geo.NewPoint(14.068129707532552, 121.32719002650704),
			geo.NewPoint(14.06788253996273, 121.32688224303081),
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
			geo.NewPoint(14.067694194798804, 121.32708640042505),
			geo.NewPoint(14.06800445535538, 121.32742234709286),
			geo.NewPoint(14.068129707532552, 121.32719002650704),
			geo.NewPoint(14.06788253996273, 121.32688224303081),
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
