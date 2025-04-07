package app

import (
	"bookstore_oauth_api/src/domain/access_token"
	"bookstore_oauth_api/src/http"
	"bookstore_oauth_api/src/repository/db"

	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication() {

	atService := access_token.NewService(db.NewRepository())
	atHandler := http.NewHandler(atService)

	router.Handle("GET", "/oauth/access_token/:access_token_id", atHandler.GetById)
	router.Handle("POST", "/oauth/access_token", atHandler.Create)

	router.Run(":8080")
}
