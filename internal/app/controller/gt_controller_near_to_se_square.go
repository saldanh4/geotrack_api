package controller

import (
	"fmt"
	u "geotrack_api/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Result struct {
	Ip       string
	As       string
	City     string
	Country  string
	Distance string
}

func (ipController *GeotrackControllerImpl) GetNearestIpToSeSquare(c *gin.Context) {
	var checkJson map[string]interface{}
	checkParam := c.Request.URL.RawQuery
	if err := c.ShouldBindJSON(&checkJson); err != nil {
		checkJson = nil
	}

	if checkParam != "" || checkJson != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "endpoint não espera nenhum parâmetro via boy ou url")
		return
	}

	if err := c.ShouldBindJSON(&givenIp); err != nil {
		givenIp.Ip = c.Query("ip")
	}

	result, err := ipController.GeotrackUsecase.GetNearestIpToSeSquareUsecase()
	if err != nil {
		status := u.ErrorHandler(err)
		c.AbortWithStatusJSON(status, gin.H{"message": err.CustomMsg})
	}

	formatedDistance := fmt.Sprintf("%.2f", result.DistanceSeSquare)

	resposta := &Result{
		Ip:       result.Query,
		As:       result.As,
		City:     result.City,
		Country:  result.Country,
		Distance: formatedDistance + "km",
	}

	c.IndentedJSON(200, gin.H{
		"Ip mais próximo da praça da Sé é": result.Query,
		"dados":                            resposta})
}
