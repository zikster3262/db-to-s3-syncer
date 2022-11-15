package routes

import (
	"concurrency/src/models"
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
		fmt.Println(err)
	}

	sql := fmt.Sprintf("SELECT * FROM request WHERE uuid = \"%v\"", r.Uuid)
	var req models.SQLRequest
	err = db.Get(&req, sql)
	if err == nil {
		log.Info("Internal server errror.")
	}

	if req.Uuid == r.Uuid {
		log.Fatal(ErrRecordExists)
	} else {
		_, err = db.NamedExec(`INSERT INTO request (uuid, time) VALUES (:uuid, :time)`, r)
		if err != nil {
			log.Error("Record was not created")
		}
	}

	c.JSON(http.StatusOK, gin.H{"id": r.Uuid, "time": r.Time})

}
