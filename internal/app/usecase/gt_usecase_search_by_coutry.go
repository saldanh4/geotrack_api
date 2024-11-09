package usecase

import (
	e "geotrack_api/config/customerrors"
	m "geotrack_api/model"
	"strings"
)

func (ipUsecase *GeotrackUsecase) GetEntriesByCountryUsecase(givenCountry string) (*[]m.GeoLocationData, *e.CustomError) {
	country := strings.ToLower(givenCountry)

	countryList, err := ipUsecase.repository.GetEntriesByCountryRepository(country)
	if err != nil {
		return nil, err
	}
	return countryList, nil
}
