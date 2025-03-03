package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"pashmak.com/pashmak/initializers"
)

func init() {
	initializers.LoadEnvVars()
}

func main() {
	serverPort := os.Getenv("SERVER_PORT")
	router := gin.Default()
	router.Run(fmt.Sprintf("localhost:%s", serverPort))
}