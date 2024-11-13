package lib

import (
	e "geotrack_api/config/customerrors"
	l "geotrack_api/config/logger"

	goip "github.com/jpiontek/go-ip-api"
)

func GetGeoData(ip string) (*goip.Location, *e.CustomError) {

	client := goip.NewClient()

	result, err := client.GetLocationForIp(ip)
	if err != nil {
		l.Logger.Error(err.Error())
		return &goip.Location{}, e.CustomErr(e.ErrInternalServer, "não foram encontrados dados de geolocalização para o IP informado")
	}

	return result, nil
}
