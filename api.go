package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	geo "github.com/kellydunn/golang-geo"
)

// Define the API struct
type API struct{}

// Define the geofence points
var points = [][]*geo.Point{
	{
		geo.NewPoint(14.067694194798804, 121.32708640042505),
		geo.NewPoint(14.06800445535538, 121.32742234709286),
		geo.NewPoint(14.068129707532552, 121.32719002650704),
		geo.NewPoint(14.06788253996273, 121.32688224303081),
	},
}

// Define the NewGeofence function

// Define the handleGeofenceCheck function
func (api *API) handleGeofenceCheck(w http.ResponseWriter, r *http.Request) {
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
	geofence := NewGeofence(points)
	if geofence.Inside(point) {
		w.Write([]byte("true"))
	} else {
		w.Write([]byte("false"))
	}
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		next(w, r)
	}
}

// Define the StartAPI function
func StartAPI() {
	api := &API{}

	http.HandleFunc("/geofence/check", api.handleGeofenceCheck)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
