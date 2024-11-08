package db

import (
	"database/sql"
	"fmt"
	"geotrack_api/config"

	l "geotrack_api/config/logger"

	_ "github.com/lib/pq"
)

func Init(cfg config.Config) (*sql.DB, error) {

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	l.Logger.Info("Successfully connected to " + cfg.DBName)
	return db, nil
}
