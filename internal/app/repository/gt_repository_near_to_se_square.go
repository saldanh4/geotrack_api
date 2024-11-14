package repository

import (
	e "geotrack_api/config/customerrors"
	l "geotrack_api/config/logger"
	m "geotrack_api/model"

	"go.uber.org/zap"
)

func (ipRepo *GeotrackRepository) GetNearestIpToSeSquareRepository() (*m.GeoLocationData, *e.CustomError) {
	var ipData m.GeoLocationData
	query, err := ipRepo.connection.Prepare(NEAREST_SE_SQUARE_QUERY)
	if err != nil {
		l.Logger.Error("erro ao preparar consulta no banco de dados", zap.Error(err))
		return nil, e.CustomErr(e.ErrDataBase, "erro interno")
	}
	defer query.Close()

	rows, err := query.Query()
	if err != nil {
		l.Logger.Error("erro ao efetuar consulta no banco de dados", zap.Error(err))
		return nil, e.CustomErr(e.ErrDataBase, "falha ao efetuar consulta no banco de dados")
	}

	for rows.Next() {
		if err = rows.Scan(&ipData.Query,
			&ipData.As,
			&ipData.City,
			&ipData.Country,
			&ipData.DistanceSeSquare); err != nil {
			l.Logger.Error("erro ao processar os dados", zap.Error(err))
			return nil, e.CustomErr(e.ErrDataBase, "Erro interno")
		}
	}

	return &ipData, nil
}
