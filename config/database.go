package config

import (
	"database/sql"
	"final-project-enigma/model/dto"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

func ConnectDb(in dto.ConfigData, logger zerolog.Logger) (*sql.DB, error) {

	logger.Info().Msg("Trying Connect to DB")

	var PsqlInfo = fmt.Sprintf("host=%s user= %s password=%s dbname=%s port=%s sslmode=disable", in.DbConfig.Host, in.DbConfig.User, in.DbConfig.Pass, in.DbConfig.Database, in.DbConfig.DbPort)

	// var PsqlInfo = fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable", in.DbConfig.Host)

	db, err := sql.Open("postgres", PsqlInfo)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to open database")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		logger.Fatal().Err(err).Msg("Failed to ping database")
		return nil, err
	}
	logger.Info().Msg("Succesfully Connected to DB")
	return db, nil
}
