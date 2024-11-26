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
	http.HandleFunc("/geofence/check", corsMiddleware(Logger(handleGeofenceCheckRequest)))
	http.HandleFunc("/employee/login", corsMiddleware(api.handleEmployeeLogin))
	log.Println("API up and running :) ")
	go menu()
	log.Fatal(http.ListenAndServe(":8080", nil))

}
