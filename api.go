package main

import (
	//"encoding/json"
	"fmt"
	"log"
	"net/http"

	geo "github.com/kellydunn/golang-geo"
)

type API struct{}

func (api *API) handleGeofence (w http.ResponseWriter, r *http.Request) {
	points := [][]*geo.Point{
		{
			geo.NewPoint(14.067694194798804, 121.32708640042505),
			geo.NewPoint(14.06800445535538, 121.32742234709286),
			geo.NewPoint(14.068129707532552, 121.32719002650704),
			geo.NewPoint(14.06788253996273, 121.32688224303081),
		},
	}

	geofence := NewGeofence(points)
	point := geo.NewPoint(14.067694194798804, 121.32708640042505)
	if geofence.Inside(point) {
		fmt.Println("Point is inside the geofence")
	} else {
		fmt.Println("Point is outside the geofence")
	}
}
func StartAPI() {
	api := &API{}

	http.HandleFunc("/geofence", api.handleGeofence)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
