package main

import (
	"database/sql"
	"log"
	"net/http"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("postgres", "postgresql://requests_0lsz_user:HD2YXsKbv57ceqtC1vCV920SLuH1D7E4@dpg-ct4h34lumphs73e62f1g-a.singapore-postgres.render.com/requests_0lsz")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	initDB()
	api := &API{database: db}
	http.HandleFunc("/geofence/check", corsMiddleware(Logger(handleGeofenceCheckRequestWrapper)))
	http.HandleFunc("/employee/login", corsMiddleware(api.handleEmployeeLogin))
	http.HandleFunc("/employee/info", corsMiddleware(api.handleGetEmployeeInfo))
	http.HandleFunc("/register", corsMiddleware(api.handleRegisterEmployee))
	http.HandleFunc("/coordinates", corsMiddleware(handleGetCoordinates))
	createGeofenceDatabase(db) 
	

	log.Println("API up and running :) ")
	go func() {
		menu()
	}()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
