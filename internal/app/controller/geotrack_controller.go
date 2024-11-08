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

func CheckIpEntryData(c *gin.Context) (*m.GivenIP, *e.CustomError) {
	var givenIP m.GivenIP

	switch c.Request.Method {
	case http.MethodGet:
		givenIP.Ip = c.Query("ip")
		if givenIP.Ip == "" {
			return nil, e.CustomErr(e.ErrInvalidInput, "campo 'ip' é obrigatório")
		} else if c.Request.ContentLength > 0 {
			return nil, e.CustomErr(e.ErrInvalidInput, "solicitações GET não devem ter dados enviados via body")
		}
	case http.MethodPost:
		if err := c.ShouldBindJSON(&givenIP); err != nil {
			return nil, e.CustomErr(e.ErrInvalidInput, "campo 'ip' é obrigatório e deve estar no formato correto")
		}
	default:
		return nil, e.CustomErr(e.ErrInvalidInput, "método não suportado")
	}

	padraoIP := `^(\d{1,3}\.){3}\d{1,3}$`
	match, _ := regexp.MatchString(padraoIP, givenIP.Ip)
	if !match {
		err := e.CustomErr(e.ErrInvalidInput, "formato de ip inválido")
		return nil, err
	}

	return &givenIP, nil
}
