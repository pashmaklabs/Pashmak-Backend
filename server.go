package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"pashmak.com/pashmak/bootstrap"
)

func init() {
	bootstrap.LoadEnvVars()
	db := bootstrap.SetUpPostgres()
	bootstrap.MakeMigrations(db)
}

func main() {
	router := gin.Default()
	router.Run(fmt.Sprintf("localhost:%s", bootstrap.SERVER_PORT))
}