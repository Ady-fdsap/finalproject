package main

import (
	"database/sql"
	"log"
	"net/http"

	geo "github.com/kellydunn/golang-geo"
)

func NewGeofence(points [][]*geo.Point, args ...interface{}) *Geofence {

	points = [][]*geo.Point{
		{
			//true latitude longitude
			geo.NewPoint(14.067694194798804, 121.32708640042505),
			geo.NewPoint(14.06800445535538, 121.32742234709286),
			geo.NewPoint(14.068129707532552, 121.32719002650704),
			geo.NewPoint(14.06788253996273, 121.32688224303081),

			//test false lat lng
			/*
				geo.NewPoint(15, 20),
				geo.NewPoint(20, 15),
				geo.NewPoint(20, 20),
				geo.NewPoint(15, 20),
			*/
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

var db *sql.DB

func main() {
	log.Println("API up and running :) ")
	initDB()
	//go menu()
	http.HandleFunc("/geofence/check", corsMiddleware(Logger(handleGeofenceCheckRequest)))
	log.Fatal(http.ListenAndServe(":8080", nil))

}
