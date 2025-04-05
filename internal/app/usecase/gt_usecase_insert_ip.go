package usecase

import (
	"geotrack_api/internal/lib"
	"geotrack_api/model"
	"time"

	e "geotrack_api/config/customerrors"
)

func (ipUsecase *GeotrackUsecaseImpl) CreateIP(givenIp *model.GivenData) (*model.GeoLocationData, *e.CustomError) {
	result, err := lib.GetGeoData(givenIp.Ip)
	if err != nil {
		return &model.GeoLocationData{}, err
	}

	distance := lib.CalculateDistanceToPracaDaSe(result)

	h := time.Now()
	ipData := model.SetIpData(result, h, distance)

	err = ipUsecase.repository.CreateIP(&ipData)
	if err != nil {
		return &model.GeoLocationData{}, err
	}

	return &ipData, nil
}
