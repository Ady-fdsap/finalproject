package main

import (
	geo "github.com/kellydunn/golang-geo"
)

func (geofence *Geofence) setInclusionTiles() {
	xVertices := geofence.getXVertices()
	yVertices := geofence.getYVertices()

	geofence.minX = getMin(xVertices)
	geofence.minY = getMin(yVertices)
	geofence.maxX = getMax(xVertices)
	geofence.maxY = getMax(yVertices)

	xRange := geofence.maxX - geofence.minX
	yRange := geofence.maxY - geofence.minY
	geofence.tileWidth = xRange / float64(geofence.granularity)
	geofence.tileHeight = yRange / float64(geofence.granularity)

	geofence.minTileX = project(geofence.minX, geofence.tileWidth)
	geofence.minTileY = project(geofence.minY, geofence.tileHeight)
	geofence.maxTileX = project(geofence.maxX, geofence.tileWidth)
	geofence.maxTileY = project(geofence.maxY, geofence.tileHeight)

	geofence.setExclusionTiles(geofence.vertices, true)
	if len(geofence.holes) > 0 {
		for _, hole := range geofence.holes {
			geofence.setExclusionTiles(hole, false)
		}
	}
}

func (geofence *Geofence) setExclusionTiles(vertices []*geo.Point, inclusive bool) {
	var tileHash float64
	var bBoxPoly []*geo.Point
	for tileX := geofence.minTileX; tileX <= geofence.maxTileX; tileX++ {
		for tileY := geofence.minTileY; tileY <= geofence.maxTileY; tileY++ {
			tileHash = (tileY-geofence.minTileY)*float64(geofence.granularity) + (tileX - geofence.minTileX)
			bBoxPoly = []*geo.Point{geo.NewPoint(tileX*geofence.tileWidth, tileY*geofence.tileHeight), geo.NewPoint((tileX+1)*geofence.tileWidth, tileY*geofence.tileHeight), geo.NewPoint((tileX+1)*geofence.tileWidth, (tileY+1)*geofence.tileHeight), geo.NewPoint(tileX*geofence.tileWidth, (tileY+1)*geofence.tileHeight), geo.NewPoint(tileX*geofence.tileWidth, tileY*geofence.tileHeight)}

			if haveIntersectingEdges(bBoxPoly, vertices) || hasPointInPolygon(vertices, bBoxPoly) {
				geofence.tiles[tileHash] = "x"
			} else if hasPointInPolygon(bBoxPoly, vertices) {
				if inclusive {
					geofence.tiles[tileHash] = "i"
				} else {
					geofence.tiles[tileHash] = "o"
				}
			} // else all points are outside the poly
		}
	}
}
