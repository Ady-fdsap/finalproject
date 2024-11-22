package main

import (
	"fmt"
	"net/http"
	"strconv"

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

func NewGeofence() *Geofence {
	points := [][]*geo.Point{
		{
			geo.NewPoint(14.067694194798804, 121.32708640042505),
			geo.NewPoint(14.06800445535538, 121.32742234709286),
			geo.NewPoint(14.068129707532552, 121.32719002650704),
			geo.NewPoint(14.06788253996273, 121.32688224303081),
		},
	}
	geofence := &Geofence{}
	geofence.tiles = make(map[float64]string)
	geofence.granularity = 50
	geofence.vertices = points[0]
	if len(points) > 1 {
		geofence.holes = points[1:]
	}
	geofence.setInclusionTiles()
	return geofence
}

func handleGeofenceCheckRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Received %s request from %s to %s?%s\n", r.Method, r.RemoteAddr, r.URL.Path, r.URL.Query())
	lat := r.URL.Query().Get("lat")
	lng := r.URL.Query().Get("lng")

	// Convert lat and lng to float64
	latFloat, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		http.Error(w, "Invalid latitude", http.StatusBadRequest)
		return
	}
	lngFloat, err := strconv.ParseFloat(lng, 64)
	if err != nil {
		http.Error(w, "Invalid longitude", http.StatusBadRequest)
		return
	}

	// Create a new geofence point
	point := geo.NewPoint(latFloat, lngFloat)

	// Check if the point is inside the geofence
	geofence := NewGeofence()
	if geofence.Inside(point) {
		w.Write([]byte("true"))
	} else {
		w.Write([]byte("false"))
	}

}
