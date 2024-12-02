package main

import (
	"bytes"
	"io"
	"io/ioutil"
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

	log.Printf("Received %s request from %s to %s?%s\n", r.Method, r.RemoteAddr, r.URL.Path, r.URL.Query())

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
	geofence := Geofences["geofence1"]
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
			if r.ContentLength > 0 {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					log.Printf("[DEBUG] Error reading request body: %v", err)
				} else {
					log.Printf("[DEBUG] Request body: %s", string(body))
					// Restore the request body so that the next handler can read it
					r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
				}
			} else {
				log.Printf("[DEBUG] Request body: (empty)")
			}
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")

		next(w, r)
	}
}

// Define the StartAPI function
func StartAPI() {
	api := &API{}
	http.HandleFunc("/geofence/check", api.handleGeofenceCheck)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
