package usecase

import (
	"errors"
	"geotrack_api/internal/app/repository"
	"net"
)

type GeotrackUsecase struct {
	repository repository.GeotrackRepository
}

func NewGeotrackUsecase(repo repository.GeotrackRepository) GeotrackUsecase {
	return GeotrackUsecase{
		repository: repo,
	}
}

func CheckIpPatern(givenIp string) (string, error) {

	checkIP, _, err := net.ParseCIDR(givenIp + "/32")
	if err != nil {
		return "", errors.New("o parâmetro informado não possui o padrão de um IP válido")
	}

	if checkIP.To4() == nil {
		return "", errors.New("o parâmetro informado não é um IP válido")
	}

	return givenIp, nil
}
