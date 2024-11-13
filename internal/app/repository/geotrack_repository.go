package repository

import "database/sql"

type GeotrackRepository struct {
	connection *sql.DB
}

const (
	INSERT_IP_DATA_QUERY         = "INSERT INTO ip_data_endpoints (as_number, city, country, countrycode, isp, lat, lon, org, query, region, regionname, status, timezone, zip, time_stamp, distance_se_square) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)"
	SEARCH_BY_IP_QUERY           = "SELECT query, isp, country, COUNT(*) as qtd FROM ip_data_endpoints  WHERE query = $1 GROUP BY query, isp, country"
	SEARCH_BY_COUNTRY_QUERY      = "SELECT * FROM ip_data_endpoints  WHERE LOWER(country) = $1 GROUP BY id ORDER BY isp"
	SEARCH_BY_COUNTRY_CODE_QUERY = "SELECT * FROM ip_data_endpoints  WHERE LOWER(countrycode) = $1 GROUP BY id ORDER BY isp"
	DELETE_ENTRIES_BY_IP         = "DELETE FROM ip_data_endpoints WHERE query = $1"
)

func NewGeotrackRepository(connection *sql.DB) GeotrackRepository {
	return GeotrackRepository{
		connection: connection,
	}
}
