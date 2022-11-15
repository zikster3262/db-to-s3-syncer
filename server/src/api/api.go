package api

import (
	"concurrency/src/routes"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func Register(db *sqlx.DB, ctx context.Context) *gin.Engine {

	router := gin.New()
	routes.Routes(router, db)

	return router
}
