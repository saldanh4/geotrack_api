package controller

import (
	e "geotrack_api/config/customerrors"
	"geotrack_api/internal/app/usecase"
	m "geotrack_api/model"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

type GeotrackController struct {
	GeotrackUsecase usecase.GeotrackUsecase
}

func NewGeotrackController(usecase usecase.GeotrackUsecase) GeotrackController {
	return GeotrackController{
		GeotrackUsecase: usecase,
	}
}

func CheckEntryData(input string, c *gin.Context) (*m.GivenIP, *e.CustomError) {
	var givenIP m.GivenIP

	switch c.Request.Method {
	case http.MethodGet:
		switch input {
		case "ip":
			givenIP.Ip = c.Query("ip")
			if givenIP.Ip == "" {
				return nil, e.CustomErr(e.ErrInvalidInput, "campo 'ip' é obrigatório")
			}
		case "country":
			givenIP.Country = c.Query("country")
			if givenIP.Country == "" {
				return nil, e.CustomErr(e.ErrInvalidInput, "campo 'country' é obrigatório")
			}
		}
		if c.Request.ContentLength > 0 {
			return nil, e.CustomErr(e.ErrInvalidInput, "solicitações GET não devem ter dados enviados via body")
		}
	case http.MethodPost:
		if err := c.ShouldBindJSON(&givenIP); err != nil {
			return nil, e.CustomErr(e.ErrInvalidInput, "campo 'ip' é obrigatório e deve estar no formato correto")
		}
	default:
		return nil, e.CustomErr(e.ErrInvalidInput, "método não suportado")
	}

	switch {
	case givenIP.Ip != "":
		if err := ValidateIp(givenIP.Ip); err != nil {
			return nil, err
		}
	case givenIP.Country != "":
		if err := ValidateCountry(givenIP.Country); err != nil {
			return nil, err
		}
	}

	return &givenIP, nil
}

func ValidateIp(ip string) *e.CustomError {
	padraoIP := `^(\d{1,3}\.){3}\d{1,3}$`
	match, _ := regexp.MatchString(padraoIP, ip)
	if !match {
		return e.CustomErr(e.ErrInvalidInput, "formato de ip inválido")
	}
	return nil
}

func ValidateCountry(country string) *e.CustomError {
	if len(country) < 2 {
		return e.CustomErr(e.ErrInvalidInput, "nome ou código de país invalido")
	}
	return nil
}
