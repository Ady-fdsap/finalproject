package main

import (
	"bytes"
	"log"
	"net/http"
	"strconv"
	"time"

	"database/sql"

	_ "github.com/lib/pq"
)

type RequestLog struct {
	Timestamp time.Time `json:"timestamp"`
	Method    string    `json:"method"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	IPAddress string    `json:"ip_address"`
}

func initDB() {
	var err error
	db, err := sql.Open("postgres", "user=ady dbname=Requests sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("[DB] Database connection established")
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS logs (
    id SERIAL PRIMARY KEY,
    timestamp TIMESTAMP NOT NULL,
    method VARCHAR(10) NOT NULL,
    latitude FLOAT NOT NULL,
    longitude FLOAT NOT NULL,
    ip_address VARCHAR(255) NOT NULL,
    response TEXT NOT NULL
);


    `)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("[DB] Logs table initialized")
	log.Println("[DB] Listening for requests...")

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS employees (
        id VARCHAR(15) PRIMARY KEY UNIQUE,
        first_name VARCHAR(50) NOT NULL,
		last_name VARCHAR(50) NOT NULL,
        date_added DATE NOT NULL,
        password VARCHAR(255) NOT NULL,
        role VARCHAR(50) NOT NULL DEFAULT 'employee'
    );
`)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("[DB] Employees table initialized")

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS login_attempts (
        id SERIAL PRIMARY KEY,
        timestamp TIMESTAMP NOT NULL,
        ip_address VARCHAR(255) NOT NULL,
        employee_id VARCHAR(15) NOT NULL,
        success BOOLEAN NOT NULL
    );
`)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("[DB] Login attempts table initialized")

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS attendance (
        id SERIAL PRIMARY KEY,
        employee_id VARCHAR(15) NOT NULL,
        check_in_time TIMESTAMP NOT NULL,
        log_out_time TIMESTAMP,
        FOREIGN KEY (employee_id) REFERENCES employees(id)
    );
`)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("[DB] Attendance table initialized")

}

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the IP address of the request
		ipAddress := r.Header.Get("X-Forwarded-For")
		if ipAddress == "" {
			ipAddress = r.RemoteAddr
		}

		// Create a custom response writer that captures the response body
		captureWriter := &responseCaptureWriter{ResponseWriter: w, body: &bytes.Buffer{}}

		if r.Method == "GET" || r.Method == "POST" || r.Method == "PATCH" || r.Method == "DELETE" {
			// Log the request
			lat := r.URL.Query().Get("lat")
			lng := r.URL.Query().Get("lng")

			// Convert lat and lng to float64
			latFloat, err := strconv.ParseFloat(lat, 64)
			if err != nil {
				http.Error(w, "Invalid latitude", http.StatusBadRequest)
				return
			}

			// watermark ni ady
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
				IPAddress: ipAddress, // Add the IP address to the log entry
			}

			// Call the next handler with the custom response writer
			next(captureWriter, r)

			// Log the response
			responseBody := captureWriter.body.String()

			// Insert the log entry into the database
			_, err = db.Exec(`
                INSERT INTO logs (timestamp, method, latitude, longitude, ip_address, response)
                VALUES ($1, $2, $3, $4, $5, $6);
            `, logEntry.Timestamp, logEntry.Method, logEntry.Latitude, logEntry.Longitude, logEntry.IPAddress, responseBody)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

type responseCaptureWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (w *responseCaptureWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
