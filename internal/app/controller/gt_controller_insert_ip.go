package controller

import (
	l "geotrack_api/config/logger"
	u "geotrack_api/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (ipController *GeotrackController) CreateIP(c *gin.Context) {

	givenIp, err := CheckIpEntryData(c)
	if err != nil {
		status := u.ErrorHandler(err)
		c.AbortWithStatusJSON(status, gin.H{"message": err.CustomMsg})
		return
	}

	result, err := ipController.GeotrackUsecase.CreateIP(givenIp)
	if err != nil {
		status := u.ErrorHandler(err)
		l.Logger.Warn("%v", zap.Error(err))
		c.AbortWithStatusJSON(status, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"ip":         result.Query,
		"country":    result.Country,
		"insertDate": result.TimeStamp},
	)
}
