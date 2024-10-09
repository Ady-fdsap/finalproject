package main

import (
	//"fmt"

	geo "github.com/kellydunn/golang-geo"
)

type Geofence struct {
	vertices     []*geo.Point
	holes        [][]*geo.Point
	tiles        map[float64]string
	granularity  int64
	geofenceType string
	minX         float64
	maxX         float64
	minY         float64
	maxY         float64
	tileWidth    float64
	tileHeight   float64
	minTileX     float64
	maxTileX     float64
	minTileY     float64
	maxTileY     float64
}

const defaultGranularity = 20

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
	geofence.setInclusionTiles() // This function is not defined
	return geofence
}
