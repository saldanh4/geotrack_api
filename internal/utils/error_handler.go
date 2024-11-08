package utils

import (
	e "geotrack_api/config/customerrors"
	l "geotrack_api/config/logger"
	"net/http"

	"go.uber.org/zap"
)

func ErrorHandler(err error) int {
	var statusCode int

	if customErr, ok := err.(*e.CustomError); ok {
		switch customErr.BaseError {
		case e.ErrDataBase:
			statusCode = http.StatusInternalServerError
			//Quando é erro interno mapeado o log é gerado na raiz do erro
		case e.ErrInvalidInput:
			statusCode = http.StatusBadRequest
			l.Logger.Warn(customErr.CustomMsg, zap.Int("code", statusCode))
		case e.ErrNotFound:
			statusCode = http.StatusNotFound
			l.Logger.Warn(customErr.CustomMsg, zap.Int("code", statusCode))
		default:
			statusCode = http.StatusInternalServerError
			l.Logger.Error(customErr.CustomMsg, zap.Int("code", statusCode))
		}
	}
	return statusCode
}
