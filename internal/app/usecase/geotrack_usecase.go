package usecase

import (
	//"errors"
	e "geotrack_api/config/customerrors"
	"geotrack_api/internal/app/repository"
	"geotrack_api/model"
	m "geotrack_api/model"
	//"net"
)

type GeotrackUsecase interface {
	CreateIP(givenIp *model.GivenData) (*model.GeoLocationData, *e.CustomError)
	DeleteEntriesByIpUsecase(givenIp string) (string, *e.CustomError)
	GetNearestIpToSeSquareUsecase() (*m.GeoLocationData, *e.CustomError)
	GetEntriesByCountryUsecase(givenCountry string) (*[]m.GeoLocationData, *e.CustomError)
	GetEntriesByIpUsecase(givenIp string) (*model.GeoLocationData, *e.CustomError)
}

type GeotrackUsecaseImpl struct {
	repository repository.GeotrackRepository
}

func NewGeotrackUsecase(repo repository.GeotrackRepository) GeotrackUsecase {
	return &GeotrackUsecaseImpl{
		repository: repo,
	}
}

// func CheckIpPatern(givenIp string) (string, error) {

// 	checkIP, _, err := net.ParseCIDR(givenIp + "/32")
// 	if err != nil {
// 		return "", errors.New("o parâmetro informado não possui o padrão de um IP válido")
// 	}

// 	if checkIP.To4() == nil {
// 		return "", errors.New("o parâmetro informado não é um IP válido")
// 	}

// 	return givenIp, nil
// }
