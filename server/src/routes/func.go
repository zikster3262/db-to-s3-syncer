package routes

import (
	"concurrency/src/models"
	"concurrency/src/utils"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	ErrDBDoesNotExists = errors.New("database does not exists")
	ErrRecordExists    = errors.New("record exists in the Database")
)

func healthEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Status OK",
	})
}

func createRequest(c *gin.Context) {

	db, ok := getDBfromCtx(c)
	if !ok {
		log.Error(ErrDBDoesNotExists)
	}

	var r models.Request
	err := bindJson(c, &r)
	if err != nil {
		utils.FailOnError(err, err.Error())
	}

	sql := fmt.Sprintf("SELECT * FROM request WHERE uuid = \"%v\"", r.Uuid)
	var req models.SQLRequest
	err = db.Get(&req, sql)
	if err == nil {
		utils.FailOnError(err, "internal server error")
	}

	if req.Uuid == r.Uuid {
		log.Fatal(ErrRecordExists)
	} else {
		_, err = db.NamedExec(`INSERT INTO request (uuid, time) VALUES (:uuid, :time)`, r)
		if err != nil {
			utils.FailOnError(err, "record not created")
		}
	}

	c.JSON(http.StatusOK, gin.H{"id": r.Uuid, "time": r.Time})

}
