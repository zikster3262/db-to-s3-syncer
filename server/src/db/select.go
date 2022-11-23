package db

import (
	"concurrency/src/models"
	"concurrency/src/utils"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var sqlAllRequests string = "SELECT * FROM request;"

func GetAllRequest(db *sqlx.DB) (requests []models.DbRequest, err error) {

	err = db.Select(&requests, sqlAllRequests)
	if err != nil {
		utils.FailOnError(err, "cant select from db")
	}
	return requests, err

}

func GetRequest(db *sqlx.DB, arg string) (request models.DbRequest, err error) {

	sql := fmt.Sprintf("SELECT * FROM request id = \"%v\";", arg)

	err = db.Select(&request, sql)
	if err != nil {
		utils.FailOnError(err, "cant select record from db")
	}
	return request, err

}
