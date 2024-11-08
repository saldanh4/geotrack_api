package lib

import (
	goip "github.com/jpiontek/go-ip-api"
	"github.com/umahmood/haversine"
)

const (
	LATSE float64 = -23.5505
	LONSE float64 = -46.6333
)

func CalculateDistanceToPracaDaSe(geoData *goip.Location) float64 {

	givenCoord := haversine.Coord{float64(geoData.Lat), float64(geoData.Lon)}
	seSquareCoord := haversine.Coord{LATSE, LONSE}

	_, distance := haversine.Distance(seSquareCoord, givenCoord)

	return distance
}
