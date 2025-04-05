package controller

import (
	u "geotrack_api/internal/utils"

	"github.com/gin-gonic/gin"
)

func (ipController *GeotrackControllerImpl) GetEntriesByCountry(c *gin.Context) {
	inputControl := "country"
	givenCountry, err := ipController.Service.CheckEntryData(inputControl, c)
	if err != nil {
		status := u.ErrorHandler(err)
		c.AbortWithStatusJSON(status, gin.H{"message": err.CustomMsg})
		return
	}

	result, err := ipController.GeotrackUsecase.GetEntriesByCountryUsecase(givenCountry.Country)
	if err != nil {
		status := u.ErrorHandler(err)
		c.AbortWithStatusJSON(status, gin.H{"message": err.CustomMsg})
		return
	}

	c.IndentedJSON(200, result)
}
