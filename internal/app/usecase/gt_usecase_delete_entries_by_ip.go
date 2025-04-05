package usecase

import (
	e "geotrack_api/config/customerrors"
)

func (ipUsecase *GeotrackUsecaseImpl) DeleteEntriesByIpUsecase(givenIp string) (string, *e.CustomError) {
	result, err := ipUsecase.repository.DeleteEntriesByIpRepository(givenIp)
	if err != nil {
		return "", err
	}
	return result, nil
}
