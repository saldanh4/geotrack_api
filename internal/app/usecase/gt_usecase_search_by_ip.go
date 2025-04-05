package usecase

import (
	e "geotrack_api/config/customerrors"
	"geotrack_api/model"
)

func (ipUsecase *GeotrackUsecaseImpl) GetEntriesByIpUsecase(givenIp string) (*model.GeoLocationData, *e.CustomError) {

	ipData, err := ipUsecase.repository.GetEntriesByIpRepository(givenIp)
	if err != nil {
		return &model.GeoLocationData{}, err
	}

	return ipData, nil
}
