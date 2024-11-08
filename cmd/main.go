package main

import (
	"geotrack_api/config"
	"geotrack_api/config/db"
	l "geotrack_api/config/logger"
	"geotrack_api/config/route"
	"geotrack_api/internal/app/controller"
	"geotrack_api/internal/app/repository"
	"geotrack_api/internal/app/usecase"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	l.LoggerInit()
	defer l.Logger.Sync()
	l.Logger.Info("Aplicação iniciada")

	server := gin.Default()

	cfg := config.LoadConfig()
	dbConnection, err := db.Init(cfg)
	if err != nil {
		l.Logger.Panic("Falha ao estabelecer conexão com o Banco de Dados", zap.Error(err))
	}

	GtRepository := repository.NewGeotrackRepository(dbConnection)
	GtUsecase := usecase.NewGeotrackUsecase(GtRepository)
	GtController := controller.NewGeotrackController(GtUsecase)
	route.Endpoints(server, &GtController)
}
