package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	geo "github.com/kellydunn/golang-geo"
)

func NewGeofence(points [][]*geo.Point, args ...interface{}) *Geofence {

	points = [][]*geo.Point{
		{
			geo.NewPoint(14.067694194798804, 121.32708640042505),
			geo.NewPoint(14.06800445535538, 121.32742234709286),
			geo.NewPoint(14.068129707532552, 121.32719002650704),
			geo.NewPoint(14.06788253996273, 121.32688224303081),
		},
	}
	geofence := &Geofence{}
	geofence.tiles = make(map[float64]string)
	if len(args) > 0 {
		geofence.granularity = args[0].(int64)
	} else {
		geofence.granularity = defaultGranularity
	}
	geofence.vertices = points[(0)]
	if len(points) > 1 {
		geofence.holes = points[1:]
		geofence.tiles = make(map[float64]string)
	}
	geofence.setInclusionTiles()
	return geofence
}

/*
	func main() {
		go StartAPI()
		fmt.Println("API up and running :) ")

		//uncomment if testing inside the code

			fmt.Print("Enter latitude: ")
			var lat float64
			fmt.Scanln(&lat)

			fmt.Print("Enter longitude: ")
			var lng float64
			fmt.Scanln(&lng)

		points := [][]*geo.Point{
			{
				geo.NewPoint(14.067694194798804, 121.32708640042505),
				geo.NewPoint(14.06800445535538, 121.32742234709286),
				geo.NewPoint(14.068129707532552, 121.32719002650704),
				geo.NewPoint(14.06788253996273, 121.32688224303081),
			},
		}

		geofence := NewGeofence(points)
		point := geo.NewPoint(lat, lng)
		if geofence.Inside(point) {
			fmt.Println("Point is inside the geofence")
		} else {
			fmt.Println("Point is outside the geofence")
		}
	}
*/
func handleGeofenceCheckRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Received %s request from %s to %s?%s\n", r.Method, r.RemoteAddr, r.URL.Path, r.URL.Query())
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

	// Create a new geofence point
	point := geo.NewPoint(latFloat, lngFloat)

	// Check if the point is inside the geofence
	geofence := NewGeofence(points)
	if geofence.Inside(point) {
		w.Write([]byte("true"))
	} else {
		w.Write([]byte("false"))
	}
}

func main() {
	fmt.Println("API up and running :) ")
	http.HandleFunc("/geofence/check", handleGeofenceCheckRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
