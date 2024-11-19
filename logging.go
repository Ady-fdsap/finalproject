package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type RequestLog struct {
	Timestamp time.Time `json:"timestamp"`
	Method    string    `json:"method"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
}

func Logger(next http.HandlerFunc) http.HandlerFunc {
	println("Logging initialized")
	return func(w http.ResponseWriter, r *http.Request) {
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

		logEntry := RequestLog{
			Timestamp: time.Now(),
			Method:    r.Method,
			Latitude:  latFloat,
			Longitude: lngFloat,
		}

		// Connect to the PostgreSQL database
		db, err := sql.Open("postgres", "user=ady dbname=Requests sslmode=disable")
		if err != nil {
			log.Println(err)
			return
		}
		defer db.Close()
		log.Println("Database connection established")

		// Create the logs table if it doesn't exist
		_, err = db.Exec(`
		CREATE TABLE logs (
			id SERIAL PRIMARY KEY,
			timestamp TIMESTAMP NOT NULL,
			method VARCHAR(10) NOT NULL,
			latitude FLOAT NOT NULL,
			longitude FLOAT NOT NULL
		);
	`)
		if err != nil {
			if err.Error() == "pq: relation \"logs\" already exists" {
				log.Println("Logs table already exists")
			} else {
				log.Println(err)
				return
			}
		}

		// Insert the log entry into the database
		_, err = db.Exec(`
			INSERT INTO logs (timestamp, method, latitude, longitude)
			VALUES ($1, $2, $3, $4);
		`, logEntry.Timestamp, logEntry.Method, logEntry.Latitude, logEntry.Longitude)
		if err != nil {
			log.Println(err)
			return
		}

		// Call the next handler in the chain
		next(w, r)
	}
}
