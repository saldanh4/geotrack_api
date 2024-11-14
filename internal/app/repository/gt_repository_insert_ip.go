package repository

import (
	e "geotrack_api/config/customerrors"
	l "geotrack_api/config/logger"
	"geotrack_api/model"

	"go.uber.org/zap"
)

func (ipRepo *GeotrackRepository) CreateIP(ipData *model.GeoLocationData) *e.CustomError {
	var id int
	query, err := ipRepo.connection.Prepare(INSERT_IP_DATA_QUERY + " RETURNING id")
	if err != nil {
		l.Logger.Error("erro ao preparar requisição SQL", zap.Error(err))
		return e.CustomErr(e.ErrInternalServer, "erro ao preparar requisição SQL")
	}
	defer query.Close()

	err = query.QueryRow(
		ipData.As,
		ipData.City,
		ipData.Country,
		ipData.CountryCode,
		ipData.Isp,
		ipData.Lat,
		ipData.Lon,
		ipData.Org,
		ipData.Query,
		ipData.Region,
		ipData.RegionName,
		ipData.Status,
		ipData.Timezone,
		ipData.Zip,
		ipData.TimeStamp,
		ipData.DistanceSeSquare).Scan(&id)
	if err != nil {
		l.Logger.Error("erro ao executar consulta no bando de dados", zap.Error(err))
		return e.CustomErr(e.ErrInternalServer, "erro ao executar consulta no bando de dados")
	}

	query.Close()
	l.Logger.Info("registro inserido com sucesso no banco de dados")
	return nil
}
