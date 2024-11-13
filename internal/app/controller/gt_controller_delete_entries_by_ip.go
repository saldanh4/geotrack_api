package controller

import (
	u "geotrack_api/internal/utils"

	"github.com/gin-gonic/gin"
)

func (ipController *GeotrackController) DeleteEntriesByIp(c *gin.Context) {
	inputControl := "ip"
	givenIp, err := CheckEntryData(inputControl, c)
	if err != nil {
		status := u.ErrorHandler(err)
		c.AbortWithStatusJSON(status, gin.H{"message": err.CustomMsg})
		return
	}

	result, err := ipController.GeotrackUsecase.DeleteEntriesByIpUsecase(givenIp.Ip)
	if err != nil {
		status := u.ErrorHandler(err)
		c.AbortWithStatusJSON(status, err.CustomMsg)
	}

	c.IndentedJSON(200, result)
}
