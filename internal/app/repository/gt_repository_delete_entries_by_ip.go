package repository

import (
	e "geotrack_api/config/customerrors"
	l "geotrack_api/config/logger"

	"go.uber.org/zap"
)

func (ipRepo *GeotrackRepository) DeleteEntriesByIpRepository(givenIp string) (string, *e.CustomError) {
	query, err := ipRepo.connection.Prepare(DELETE_ENTRIES_BY_IP)
	if err != nil {
		l.Logger.Error("erro ao preparar consulta no banco de dados", zap.Error(err))
		return "", e.CustomErr(e.ErrDataBase, "erro interno")
	}
	defer query.Close()

	result, err := query.Exec(givenIp)
	if err != nil {
		l.Logger.Error("erro ao deletar registro no banco de dados", zap.Error(err))
		return "", e.CustomErr(e.ErrDataBase, "erro interno")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		l.Logger.Error("falha ao verificar linhas afetadas", zap.Error(err))
		return "", e.CustomErr(e.ErrDataBase, "erro interno")
	}

	if rowsAffected == 0 {
		return "não foram localizados registros para o IP " + givenIp, nil
	}

	return "registros deletados com sucesso", nil
}
