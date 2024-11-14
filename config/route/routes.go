package route

import (
	l "geotrack_api/config/logger"
	"geotrack_api/internal/app/controller"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func setTrustedProxies(server *gin.Engine) {
	err := server.SetTrustedProxies([]string{"192.168.1.1", "192.168.1.2"})
	if err != nil {
		l.Logger.Fatal("Falha ao definir proxies confiáveis: %v", zap.Error(err))
	}
}

func Endpoints(server *gin.Engine, controller *controller.GeotrackController) {

	setTrustedProxies(server)
	v1 := server.Group("api/v1/ips")
	{
		v1.POST("/add_ip", controller.CreateIP)
		v1.GET("/search_by_ip", controller.GetEntriesByIp)
		v1.GET("/search_by_country", controller.GetEntriesByCountry)
		v1.GET("/nearest_se_square", controller.GetNearestIpToSeSquare)
		v1.DELETE("delete_entries_by_ip", controller.DeleteEntriesByIp)
	}

	l.Logger.Info("Server is running on port: 8080")
	if err := server.Run(":8080"); err != nil {
		l.Logger.Fatal("Não foi possível iniciar o servidor: %v", zap.Error(err))
	}
}
