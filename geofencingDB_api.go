package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	_ "github.com/lib/pq"
)

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Coordinates struct {
	Point1 Coordinate `json:"point1"`
	Point2 Coordinate `json:"point2"`
	Point3 Coordinate `json:"point3"`
	Point4 Coordinate `json:"point4"`
}

func createGeofenceDatabase(db *sql.DB) error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS geofence_coordinates (
            id SERIAL PRIMARY KEY,
			latitude FLOAT NOT NULL,
            longitude FLOAT NOT NULL,
        );
    `)
	return err
}

func handleGetCoordinates(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the query string
	queryString := r.URL.Query().Get("data")

	// Decode the URL query string
	decodedQueryString, err := url.QueryUnescape(queryString)
	if err != nil {
		http.Error(w, "Failed to decode query string", http.StatusInternalServerError)
		return
	}

	// Unmarshal the JSON data
	var coords Coordinates
	err = json.Unmarshal([]byte(decodedQueryString), &coords)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	// Connect to the database
	db, err = sql.Open("postgres", "postgresql://requests_0lsz_user:HD2YXsKbv57ceqtC1vCV920SLuH1D7E4@dpg-ct4h34lumphs73e62f1g-a.singapore-postgres.render.com/requests_0lsz")
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("TRUNCATE TABLE geofence_coordinates")
	if err != nil {
		http.Error(w, "Failed to delete existing coordinates", http.StatusInternalServerError)
		return
	}

	// Insert the coordinates into the database
	for _, point := range []Coordinate{coords.Point1, coords.Point2, coords.Point3, coords.Point4} {
		_, err = db.Exec(`
            INSERT INTO geofence_coordinates (longitude, latitude)
            VALUES ($1, $2);
        `, point.Longitude, point.Latitude)
		if err != nil {
			http.Error(w, "Failed to insert coordinate into database", http.StatusInternalServerError)
			return
		}
	}

	// Update the geofences
	updateGeofences(db)
	w.Write([]byte("Coordinates inserted successfully, updated the geofence"))
	log.Println("Coordinates inserted successfully, updated the geofence")
}
