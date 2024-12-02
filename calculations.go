package main

import (
	//"fmt"
	"math"

	geo "github.com/kellydunn/golang-geo"
)

//calculations lang para sa tiles and stuff for the geofence

func getMin(slice []float64) float64 {
	var min float64
	if len(slice) > 0 {
		min = slice[0]
	}
	for i := 1; i < len(slice); i++ {
		if slice[i] < min {
			min = slice[i]
		}
	}
	return min
}

// watermark ni ady
func getMax(slice []float64) float64 {
	var max float64
	if len(slice) > 0 {
		max = slice[0]
	}
	for i := 1; i < len(slice); i++ {
		if slice[i] > max {
			max = slice[i]
		}
	}
	return max
}

func project(value float64, tileSize float64) float64 {
	return math.Floor(value / tileSize)
}

func (geofence *Geofence) getXVertices() []float64 {
	xVertices := make([]float64, len(geofence.vertices))
	for i := 0; i < len(geofence.vertices); i++ {
		xVertices[i] = geofence.vertices[i].Lat()
	}
	return xVertices
}

func (geofence *Geofence) getYVertices() []float64 {
	yVertices := make([]float64, len(geofence.vertices))
	for i := 0; i < len(geofence.vertices); i++ {
		yVertices[i] = geofence.vertices[i].Lng()
	}
	return yVertices
}

func hasPointInPolygon(sourcePoly []*geo.Point, targetPoly []*geo.Point) bool {
	tPolygon := geo.NewPolygon(targetPoly)
	for idx := 0; idx < len(sourcePoly)-1; idx++ {
		if tPolygon.Contains(sourcePoly[idx]) {
			return true
		}
	}
	return false
}

func haveIntersectingEdges(poly1 []*geo.Point, poly2 []*geo.Point) bool {
	for idx1 := 0; idx1 < len(poly1)-1; idx1++ {
		for idx2 := 0; idx2 < len(poly2)-1; idx2++ {
			if segmentsIntersect(poly1[idx1], poly1[idx1+1], poly2[idx2], poly2[idx2+1]) {
				return true
			}
		}
	}
	return false
}

func segmentsIntersect(s1p1 *geo.Point, s1p2 *geo.Point, s2p1 *geo.Point, s2p2 *geo.Point) bool {
	p := s1p1
	r := vectorDifference(s1p2, s1p1)
	q := s2p1
	s := vectorDifference(s2p2, s2p1)

	rCrossS := vectorCrossProduct(r, s)
	qMinusP := vectorDifference(q, p)

	if rCrossS == 0 {
		if vectorCrossProduct(qMinusP, r) == 0 {
			return true
		} else {
			return false
		}
	}

	t := vectorCrossProduct(qMinusP, s) / rCrossS
	u := vectorCrossProduct(qMinusP, r) / rCrossS
	return t >= 0 && t <= 1 && u >= 0 && u <= 1
}

func vectorDifference(p1 *geo.Point, p2 *geo.Point) *geo.Point {
	return geo.NewPoint(p1.Lat()-p2.Lat(), p1.Lng()-p2.Lng())
}

// sarap matulog gagi
func vectorCrossProduct(p1 *geo.Point, p2 *geo.Point) float64 {
	return p1.Lat()*p2.Lng() - p1.Lng()*p2.Lat()
}

func (geofence *Geofence) Inside(point *geo.Point) bool {
	// Bbox check first
	if point.Lat() < geofence.minX || point.Lat() > geofence.maxX || point.Lng() < geofence.minY || point.Lng() > geofence.maxY {
		//fmt.Println("Point is outside bounding box")
		return false
	}

	//watermark ni ady
	tileHash := (project(point.Lng(), geofence.tileHeight)-geofence.minTileY)*float64(geofence.granularity) + (project(point.Lat(), geofence.tileWidth) - geofence.minTileX)
	//fmt.Println("Tile hash:", tileHash)
	intersects := geofence.tiles[tileHash]
	//fmt.Println("Intersects:", intersects)

	if intersects == "i" {
		//fmt.Println("Point is inside tile")
		return true
	} else if intersects == "x" {
		//fmt.Println("Point is in tile with exclusion")
		polygon := geo.NewPolygon(geofence.vertices)
		inside := polygon.Contains(point)
		if !inside || len(geofence.holes) == 0 {
			//fmt.Println("Point is not inside polygon")
			return inside
		}

		// Check if point is inside any of the holes
		for _, hole := range geofence.holes {
			holePoly := geo.NewPolygon(hole)
			if holePoly.Contains(point) {
				//fmt.Println("Point is inside hole")
				return false
			}
		}

		//fmt.Println("Point is inside geofence")
		return true
	} else {
		//fmt.Println("Point is outside geofence")
		return false
	}
}

// password registration validation for capital letter
func hasCapitalLetter(s string) bool {
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			return true
		}
	}
	return false
}

// password registration validation for number
func hasNumber(s string) bool {
	for _, r := range s {
		if r >= '0' && r <= '9' {
			return true
		}
	}
	return false
}
