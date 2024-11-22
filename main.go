package main

import (
	"database/sql"
	"log"
	"net/http"
)


var db *sql.DB

func main() {
	log.Println("API up and running :) ")
	initDB()
	api := &API{database: db}
	go menu()
	http.HandleFunc("/geofence/check", corsMiddleware(Logger(handleGeofenceCheckRequest)))
	http.HandleFunc("/employee/login", corsMiddleware(api.handleEmployeeLogin))
	log.Fatal(http.ListenAndServe(":8080", nil))

}
