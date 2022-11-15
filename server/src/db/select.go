package db

import (
	"concurrency/src/models"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var sqlAllRequests string = "SELECT * FROM request;"

func GetAllRequest(db *sqlx.DB) (requests []models.DbRequest, err error) {

	err = db.Select(&requests, sqlAllRequests)
	if err != nil {
		log.Error().Msg("cant select from db")
	}
	return requests, err

}

func GetRequest(db *sqlx.DB, arg string) (request models.DbRequest, err error) {

	sql := fmt.Sprintf("SELECT * FROM request id = \"%v\";", arg)

	err = db.Select(&request, sql)
	if err != nil {
		log.Error().Msg("cant select record from db")
	}
	return request, err

}
