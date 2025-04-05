package controller

import (
	"geotrack_api/config/customerrors"
	e "geotrack_api/config/customerrors"
	"geotrack_api/internal/app/usecase"
	"geotrack_api/model"
	m "geotrack_api/model"
	"log"
	"net"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

type GeotrackController interface {
	CreateIP(c *gin.Context)
	DeleteEntriesByIp(c *gin.Context)
	GetNearestIpToSeSquare(c *gin.Context)
	GetEntriesByCountry(c *gin.Context)
	GetEntriesByIp(c *gin.Context)
}

type CheckService interface {
	CheckEntryData(input string, c *gin.Context) (*model.GivenData, *e.CustomError)
}

type CheckData interface {
	CheckInputData(input string, c *gin.Context) (*m.GivenData, *e.CustomError)
	ValidateIp(ip string) *customerrors.CustomError
	ValidateCountry(country string) *customerrors.CustomError
}

type GeotrackControllerImpl struct {
	GeotrackUsecase usecase.GeotrackUsecase
	Service         CheckService
	Data            CheckData
}
type CheckServiceImpl struct {
	//GeotrackController GeotrackController
	CheckData CheckData
}

type CheckDataImpl struct {
}

func NewDefaultCheckService() CheckService {
	return &CheckServiceImpl{
		CheckData: &CheckDataImpl{},
	}
}

func NewDefaultCheckData() CheckData {
	return &CheckDataImpl{}
}

func NewGeotrackController(usecase usecase.GeotrackUsecase, service CheckService, data CheckData) GeotrackController {
	if service == nil {
		service = NewDefaultCheckService()
		data = NewDefaultCheckData()
	}
	return &GeotrackControllerImpl{
		GeotrackUsecase: usecase,
		Service:         service,
		Data:            data,
	}
}

var givenIp m.GivenIP
var givenCountry m.GivenCountry
var givenData m.GivenData

func (control *CheckServiceImpl) CheckEntryData(input string, c *gin.Context) (*m.GivenData, *e.CustomError) {
	log.Printf("Método: %s, Input: %s, Query: %s, Body: %s", c.Request.Method, input, c.Request.URL.Query(), c.Request.Body)
	givenIp.Ip = ""
	givenCountry.Country = ""
	givenData.Ip = ""
	givenData.Country = ""
	switch c.Request.Method {
	case http.MethodGet:
		if c.Request.ContentLength > 0 {
			return nil, e.CustomErr(e.ErrInvalidInput, "solicitações "+http.MethodGet+" não devem ter dados enviados via body")
		}
		result, err := control.CheckData.CheckInputData(input, c)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		givenData = *result
	case http.MethodPost:
		result, err := control.CheckData.CheckInputData(input, c)
		if err != nil {
			return nil, err
		}
		givenData = *result
		if len(c.Request.URL.Query()) > 0 {
			return nil, e.CustomErr(e.ErrInvalidInput, "solicitações "+http.MethodPost+" não devem ter dados enviados via url")
		}
	case http.MethodDelete:
		result, err := control.CheckData.CheckInputData(input, c)
		if err != nil {
			return nil, err
		}
		givenData = *result
		if c.Request.ContentLength > 0 {
			return nil, e.CustomErr(e.ErrInvalidInput, "solicitações "+http.MethodDelete+" não devem ter dados enviados via body")
		}
	default:
		return nil, e.CustomErr(e.ErrInvalidInput, "") //método não suportado")
	}

	switch {
	case givenData.Ip != "":
		if err := control.CheckData.ValidateIp(givenData.Ip); err != nil {
			return nil, err
		}
	case givenData.Country != "":
		if err := control.CheckData.ValidateCountry(givenData.Country); err != nil {
			return nil, err
		}
	}

	return &givenData, nil
}

func (dtControl *CheckDataImpl) CheckInputData(input string, c *gin.Context) (*m.GivenData, *e.CustomError) {

	switch input {
	case "ip":
		if err := c.ShouldBindJSON(&givenIp); err != nil {
			givenIp.Ip = c.Query("ip")
		}
		givenData.Ip = givenIp.Ip
		if givenData.Ip == "" {
			return nil, e.CustomErr(e.ErrInvalidInput, "campo 'ip' é obrigatório")
		}
		//givenData.Country = ""
	case "country":
		if err := c.ShouldBindJSON(&givenCountry); err != nil {
			givenCountry.Country = c.Query("country")
		}
		givenData.Country = givenCountry.Country
		if givenData.Country == "" {
			return nil, e.CustomErr(e.ErrInvalidInput, "campo 'country' é obrigatório")
		}
		//givenData.Ip = ""
	default:
		return &givenData, e.CustomErr(e.ErrInternalServer, "falha ao checar se a entrada é 'ip' ou 'país'")
	}
	return &givenData, nil
}

func (dtControl *CheckDataImpl) ValidateIp(ip string) *e.CustomError {
	padraoIP := `^(\d{1,3}\.){3}\d{1,3}$`
	match, _ := regexp.MatchString(padraoIP, ip)
	if !match {
		return e.CustomErr(e.ErrInvalidInput, "formato de ip inválido")
	}
	_, _, err := net.ParseCIDR(ip + "/32")
	if err != nil {
		return e.CustomErr(e.ErrInvalidInput, "formato de ip inválido")
	}
	return nil
}

func (dtControl *CheckDataImpl) ValidateCountry(country string) *e.CustomError {
	regex := `^[a-zA-Z\s]+$`
	validCountry, _ := regexp.MatchString(regex, country)
	if len(country) < 2 || !validCountry {
		return e.CustomErr(e.ErrInvalidInput, "nome ou código de país invalido")
	}
	return nil
}
