package initializers

import (
    "log"
	"os"
    "github.com/joho/godotenv"
)

var (
	ServerPort = os.Getenv("SERVER_PORT")
)

func LoadEnvVars(){
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
}