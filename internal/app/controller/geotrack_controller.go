package controller

import (
	e "geotrack_api/config/customerrors"
	"geotrack_api/internal/app/usecase"
	m "geotrack_api/model"
	"net"
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

var givenIp m.GivenIP
var givenCountry m.GivenCountry
var givenData m.GivenData

func CheckEntryData(input string, c *gin.Context) (*m.GivenData, *e.CustomError) {
	switch c.Request.Method {
	case http.MethodGet:
		result, err := CheckInputData(input, c)
		if err != nil {
			return nil, err
		}
		givenData = *result
		if c.Request.ContentLength > 0 {
			return nil, e.CustomErr(e.ErrInvalidInput, "solicitações "+http.MethodGet+" não devem ter dados enviados via body")
		}
	case http.MethodPost:
		result, err := CheckInputData(input, c)
		if err != nil {
			return nil, err
		}
		givenData = *result
		if len(c.Request.URL.Query()) > 0 {
			return nil, e.CustomErr(e.ErrInvalidInput, "solicitações "+http.MethodPost+" não devem ter dados enviados via url")
		}

	case http.MethodDelete:
		result, err := CheckInputData(input, c)
		if err != nil {
			return nil, err
		}
		givenData = *result
		if c.Request.ContentLength > 0 {
			return nil, e.CustomErr(e.ErrInvalidInput, "solicitações "+http.MethodDelete+" não devem ter dados enviados via body")
		}
	default:
		return nil, e.CustomErr(e.ErrInvalidInput, "método não suportado")
	}

	switch {
	case givenData.Ip != "":
		if err := ValidateIp(givenData.Ip); err != nil {
			return nil, err
		}
	case givenData.Country != "":
		if err := ValidateCountry(givenCountry.Country); err != nil {
			return nil, err
		}
	}

	return &givenData, nil
}

func ValidateIp(ip string) *e.CustomError {
	padraoIP := `^(\d{1,3}\.){3}\d{1,3}$`
	match, _ := regexp.MatchString(padraoIP, ip)
	if !match {
		return e.CustomErr(e.ErrInvalidInput, "formato de ip inválido")
	}
	_, _, err := net.ParseCIDR(ip + "/32")
	if err != nil {
		return e.CustomErr(e.ErrInvalidInput, "formato de ip inválido")
	} // else if checkIp.To4() == nil {
	// 	return e.CustomErr(e.ErrInvalidInput, "formato de ip inválido")
	// }
	return nil
}

func ValidateCountry(country string) *e.CustomError {
	regex := `^[a-zA-Z\s]+$`
	validCountry, _ := regexp.MatchString(regex, country)
	if len(country) < 2 || !validCountry {
		return e.CustomErr(e.ErrInvalidInput, "nome ou código de país invalido")
	}
	return nil
}

func CheckInputData(input string, c *gin.Context) (*m.GivenData, *e.CustomError) {
	switch input {
	case "ip":
		if err := c.ShouldBindJSON(&givenIp); err != nil {
			givenIp.Ip = c.Query("ip")
		}
		givenData.Ip = givenIp.Ip
		if givenData.Ip == "" {
			return nil, e.CustomErr(e.ErrInvalidInput, "campo 'ip' é obrigatório")
		}
		//givenData.Ip = givenIp.Ip
	case "country":
		if err := c.ShouldBindJSON(&givenCountry); err != nil {
			givenCountry.Country = c.Query("country")
		}
		givenData.Country = givenCountry.Country
		if givenData.Country == "" {
			return nil, e.CustomErr(e.ErrInvalidInput, "campo 'country' é obrigatório")
		}
	}
	return &givenData, nil
}
