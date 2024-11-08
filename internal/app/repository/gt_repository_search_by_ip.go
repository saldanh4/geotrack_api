package repository

import (
	"database/sql"
	e "geotrack_api/config/customerrors"
	l "geotrack_api/config/logger"
	"geotrack_api/model"

	"go.uber.org/zap"
)

func (ipRepo *GeotrackRepository) GetEntriesByIpRepository(giVenIp string) (*model.GeoLocationData, *e.CustomError) {

	var ipData model.GeoLocationData

	query, err := ipRepo.connection.Prepare(SEARCH_BY_IP_QUERY)
	if err != nil {
		l.Logger.Error("erro ao efetuar consulta no banco de dados", zap.Error(err))
		return &model.GeoLocationData{}, &e.CustomError{BaseError: e.ErrDataBase, CustomMsg: "falha ao efetuar consulta no banco de dados"}
	}

	if err := query.QueryRow(giVenIp).Scan(
		&ipData.Query,
		&ipData.Isp,
		&ipData.Country,
		&ipData.Count); err != nil {
		if err == sql.ErrNoRows {
			return &model.GeoLocationData{}, &e.CustomError{BaseError: e.ErrNotFound, CustomMsg: "registro não localizado no banco de dados"}
		}
		l.Logger.Error("erro ao efetuar consulta no banco de dados", zap.Error(err))
		return &model.GeoLocationData{}, &e.CustomError{BaseError: e.ErrDataBase, CustomMsg: "erro interno"}
	}

	query.Close()
	return &ipData, nil

}
