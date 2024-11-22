package main

import (
	"database/sql"
	"log"
	"net/http"

	geo "github.com/kellydunn/golang-geo"
)

func NewGeofence() *Geofence {
	points := [][]*geo.Point{
		{
			geo.NewPoint(14.067694194798804, 121.32708640042505),
			geo.NewPoint(14.06800445535538, 121.32742234709286),
			geo.NewPoint(14.068129707532552, 121.32719002650704),
			geo.NewPoint(14.06788253996273, 121.32688224303081),
		},
	}
	geofence := &Geofence{}
	geofence.tiles = make(map[float64]string)
	geofence.granularity = 50
	geofence.vertices = points[0]
	if len(points) > 1 {
		geofence.holes = points[1:]
	}
	geofence.setInclusionTiles()
	return geofence
}

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
