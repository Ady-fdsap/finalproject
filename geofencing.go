package main

import (
	"net/http"
	"strconv"

	geo "github.com/kellydunn/golang-geo"
)

// ito ung structure or parang blueprint ng mga factors/variables na ginagamit ng geofence
type Geofence struct {
	ID          string
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

type GeofenceManager struct {
	geofences map[string]*Geofence
}

func (gm *GeofenceManager) GetGeofence(id string) *Geofence {
	return gm.geofences[id]
}

var Geofences = map[string]*Geofence{
	"geofence1": NewGeofence("geofence1", [][]*geo.Point{
		{
			geo.NewPoint(14.067694194798804, 121.32708640042505),
			geo.NewPoint(14.06800445535538, 121.32742234709286),
			geo.NewPoint(14.068129707532552, 121.32719002650704),
			geo.NewPoint(14.06788253996273, 121.32688224303081),
		},
	}),
	//"geofence2": NewGeofence("geofence2", [][]*geo.Point{
	//	{
	//		geo.NewPoint(14.067694194798804, 121.32708640042505),
	//		geo.NewPoint(14.06800445535538, 121.32742234709286),
	//		geo.NewPoint(14.068129707532552, 121.32719002650704),
	//		geo.NewPoint(14.06788253996273, 121.32688224303081),
	//	},
	//}),
}

func NewGeofence(id string, points [][]*geo.Point) *Geofence {
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

// request wrapper for the handleGeofenceCheck function
func handleGeofenceCheckRequestWrapper(w http.ResponseWriter, r *http.Request) {
	gm := &GeofenceManager{} // You need to get the GeofenceManager from somewhere
	var _ *GeofenceManager = gm
	handleGeofenceCheckRequest(w, r)
}

func handleGeofenceCheckRequest(w http.ResponseWriter, r *http.Request) {
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

	// Check if the point is inside any geofence
	for _, geofence := range Geofences {
		if geofence.Inside(point) {
			w.Write([]byte("true"))
			return
		}
	}

	// If the point is not inside any geofence, return false
	w.Write([]byte("false"))
}
