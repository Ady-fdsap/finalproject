package main

import (
	"fmt"

	"os"

	"github.com/tidwall/geojson"
)

func main() {
	// Parse GeoJSON
	geojson.NewPoint([]float64{1, 2})
	fmt.Fprintln(os.Stdout, []any{geojson}...)
}