package repository

import (
	"database/sql"
	e "geotrack_api/config/customerrors"
	l "geotrack_api/config/logger"
	m "geotrack_api/model"

	"go.uber.org/zap"
)

func (ipRrepo *GeotrackRepository) GetEntriesByCountryRepository(givenCountry string) (*[]m.GeoLocationData, *e.CustomError) {
	var results []m.GeoLocationData
	query := &sql.Stmt{}
	if len(givenCountry) == 2 {
		stmt, err := ipRrepo.connection.Prepare(SEARCH_BY_COUNTRY_CODE_QUERY)
		if err != nil {
			l.Logger.Error("erro ao preparar consulta no banco de dados", zap.Error(err))
			return nil, e.CustomErr(e.ErrDataBase, "falha ao preparar consulta no banco de dados")
		}
		query = stmt
	} else {
		stmt, err := ipRrepo.connection.Prepare(SEARCH_BY_COUNTRY_QUERY)
		if err != nil {
			l.Logger.Error("erro ao efetuar consulta no banco de dados", zap.Error(err))
			return nil, e.CustomErr(e.ErrDataBase, "falha ao preparar consulta no banco de dados")
		}
		query = stmt
	}
	defer query.Close()

	rows, err := query.Query(givenCountry)
	if err != nil {
		l.Logger.Error("erro ao efetuar consulta no banco de dados", zap.Error(err))
		return nil, e.CustomErr(e.ErrDataBase, "falha ao efetuar consulta no banco de dados")
	}
	defer rows.Close()

	for rows.Next() {
		var ipData m.GeoLocationData
		if err := rows.Scan(&ipData.Id,
			&ipData.As,
			&ipData.City,
			&ipData.Country,
			&ipData.CountryCode,
			&ipData.Isp,
			&ipData.Lat,
			&ipData.Lon,
			&ipData.Org,
			&ipData.Query,
			&ipData.Region,
			&ipData.RegionName,
			&ipData.Status,
			&ipData.Timezone,
			&ipData.Zip,
			&ipData.TimeStamp,
			&ipData.DistanceSeSquare); err != nil {
			l.Logger.Error("erro ao processar os dados", zap.Error(err))
			return nil, e.CustomErr(e.ErrDataBase, "Erro interno")
		}
		results = append(results, ipData)
	}

	if len(results) == 0 {
		return nil, e.CustomErr(e.ErrNotFound, "Não foram localizados dados para "+givenCountry+" em nosso banco de dados.")
	}

	return &results, nil
}
