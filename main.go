package main

import (
	"database/sql"
	"log"
	"net/http"
)

var db *sql.DB

func main() {
	initDB()
	api := &API{database: db}
	http.HandleFunc("/geofence/check", corsMiddleware(Logger(handleGeofenceCheckRequestWrapper)))
	http.HandleFunc("/employee/login", corsMiddleware(api.handleEmployeeLogin))
	http.HandleFunc("/employee/info", corsMiddleware(api.handleGetEmployeeInfo))
	http.HandleFunc("/coordinates", corsMiddleware(handleGetCoordinates))
	createGeofenceDatabase(db)

	log.Println("API up and running :) ")
	//log.Println(" Created by Group 1, Batch 7 Interns 2024 :DD")
	go menu()
	log.Fatal(http.ListenAndServe(":8080", nil))

}
