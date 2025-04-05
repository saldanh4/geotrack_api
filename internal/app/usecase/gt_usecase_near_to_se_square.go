package usecase

import (
	e "geotrack_api/config/customerrors"
	m "geotrack_api/model"
)

func (ipUsecase *GeotrackUsecaseImpl) GetNearestIpToSeSquareUsecase() (*m.GeoLocationData, *e.CustomError) {
	result, err := ipUsecase.repository.GetNearestIpToSeSquareRepository()
	if err != nil {
		return nil, err
	}
	return result, nil
}
