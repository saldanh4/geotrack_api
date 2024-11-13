package model

import (
	"time"

	goip "github.com/jpiontek/go-ip-api"
)

type GeoLocationData struct {
	Id               int     `json:"id"`
	As               string  `json:"as"`
	City             string  `json:"city"`
	Country          string  `json:"country"`
	CountryCode      string  `json:"countryCode"`
	Isp              string  `json:"isp"`
	Lat              float32 `json:"lat"`
	Lon              float32 `json:"lon"`
	Org              string  `json:"org"`
	Query            string  `json:"query"`
	Region           string  `json:"region"`
	RegionName       string  `json:"regionName"`
	Status           string  `json:"status"`
	Timezone         string  `json:"timezone"`
	Zip              string  `json:"zip"`
	TimeStamp        string  `json:"timeStamp"`
	Count            int8    `json:"count"`
	DistanceSeSquare float64 `json:"distanceSeSquare"`
}

type GivenIP struct {
	Ip string `json:"ip" binding:"required"`
}
type GivenCountry struct {
	Country string `json:"country" binding:"required"`
}

type GivenData struct {
	Ip      string
	Country string
}

func SetIpData(result *goip.Location, h time.Time, distance float64) GeoLocationData {
	ipData := GeoLocationData{
		As:               result.As,
		City:             result.City,
		Country:          result.Country,
		CountryCode:      result.CountryCode,
		Isp:              result.Isp,
		Lat:              result.Lat,
		Lon:              result.Lon,
		Org:              result.Org,
		Query:            result.Query,
		Region:           result.Region,
		RegionName:       result.RegionName,
		Status:           result.Status,
		Timezone:         result.Timezone,
		Zip:              result.Zip,
		TimeStamp:        h.Format("02/Jan/2006 15:04:05"),
		DistanceSeSquare: distance,
	}
	return ipData
}
