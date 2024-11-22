package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"database/sql"

	geo "github.com/kellydunn/golang-geo"
)

// Define the API struct
type API struct {
	database *sql.DB
}

// Define the handleGeofenceCheck function
func (api *API) handleGeofenceCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

// Define the corsMiddleware function
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodOptions {
			log.Printf("[DEBUG] Request method: %s, URL: %s", r.Method, r.URL)
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

		next(w, r)
	}
}

// Define the StartAPI function
func StartAPI() {
	api := &API{}

	http.HandleFunc("/geofence/check", api.handleGeofenceCheck)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
