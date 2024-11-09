package controller

import (
	u "geotrack_api/internal/utils"

	"github.com/gin-gonic/gin"
)

func (ipController *GeotrackController) GetEntriesByCountry(c *gin.Context) {
	inputControl := "country"
	givenCountry, err := CheckEntryData(inputControl, c)
	if err != nil {
		status := u.ErrorHandler(err)
		c.AbortWithStatusJSON(status, gin.H{"message": err.CustomMsg})
		return
	}

	result, err := ipController.GeotrackUsecase.GetEntriesByCountryUsecase(givenCountry.Country)
	if err != nil {
		return
	}

	c.IndentedJSON(200, result)
}
