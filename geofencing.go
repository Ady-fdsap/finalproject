package main

import (
	geo "github.com/kellydunn/golang-geo"
)

// ito ung structure or parang blueprint ng mga factors/variables na ginagamit ng geofence
type Geofence struct {
	vertices    []*geo.Point
	holes       [][]*geo.Point
	tiles       map[float64]string
	granularity int64
	minX        float64
	maxX        float64
	minY        float64
	maxY        float64
	tileWidth   float64
	tileHeight  float64
	minTileX    float64
	maxTileX    float64
	minTileY    float64
	maxTileY    float64
}

// higher value = more tiles = more accuracy
const defaultGranularity = 50
