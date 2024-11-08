package controller

import (
	"fmt"
	l "geotrack_api/config/logger"
	u "geotrack_api/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type IpData struct {
	IP      string
	Isp     string
	Country string
	Count   int8
}

func (ipController *GeotrackController) GetEntriesByIp(c *gin.Context) {
	givenIp, err := CheckIpEntryData(c)
	if err != nil {
		status := u.ErrorHandler(err)
		c.AbortWithStatusJSON(status, gin.H{"message": err.CustomMsg})
		return
	}

	result, err := ipController.GeotrackUsecase.GetEntriesByIpUsecase(givenIp.Ip)
	if err != nil {
		status := u.ErrorHandler(err)
		c.AbortWithStatusJSON(status, gin.H{"message": err.CustomMsg})
		return
	}

	ipData := &IpData{
		IP:      result.Query,
		Isp:     result.Isp,
		Country: result.Country,
		Count:   result.Count,
	}
	l.Logger.Info("Consulta realizada com sucesso", zap.Int("Status", http.StatusOK))
	c.IndentedJSON(http.StatusOK, ipData)
	fmt.Print(result)
}
