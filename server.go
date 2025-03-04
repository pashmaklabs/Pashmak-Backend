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
	router := gin.Default()
	router.Run(fmt.Sprintf("localhost:%s", initializers.ServerPort))
}