package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
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

		// Marshal log entry to JSON
		jsonLog, err := json.Marshal(logEntry)
		if err != nil {
			log.Println(err)
			return
		}

		// Write log entry to file
		logFile, err := os.OpenFile("./requests.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
			return
		}
		defer logFile.Close()

		_, err = logFile.Write(append(jsonLog, '\n'))
		if err != nil {
			log.Println(err)
			return

		}
		phTimezone, err := time.LoadLocation("Asia/Manila")
		if err != nil {
			log.Println(err)
			return
		}

		timestamp := time.Now().In(phTimezone)
		timestampStr := timestamp.Format("2006-01-02,15:04:05-07:00")
		log.Println(timestampStr)

		logEntry = RequestLog{
			Timestamp: timestamp,
			Method:    r.Method,
			Latitude:  latFloat,
			Longitude: lngFloat,
		}
		// Call next handler
		next(w, r)
	}
}
