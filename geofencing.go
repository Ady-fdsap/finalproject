package main

import (
	//"fmt"

	geo "github.com/kellydunn/golang-geo"
)

//ito ung structure or parang blueprint ng mga factors/variables na ginagamit ng geofence
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

