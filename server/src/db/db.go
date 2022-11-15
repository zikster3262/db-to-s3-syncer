package db

import (
	"concurrency/src/models"
	"log"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	SQLXConnection *sqlx.DB
	err            error
)

type SQLxProvider struct{}

var onceSQLx sync.Once

func OpenSQLx() *sqlx.DB {

	onceSQLx.Do(func() {
		dsn := os.Getenv("DB_URL")

		SQLXConnection, err = sqlx.Connect("mysql", dsn)
		if err != nil {
			log.Fatal(err)
		}
	})

	return SQLXConnection
}

func (s *SQLxProvider) Insert(r models.Request) {

}
