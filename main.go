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
	log.Println("API up and running :) ")
	go menu()
	log.Fatal(http.ListenAndServe(":8080", nil))

}
