package bootstrap

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	SERVER_PORT string
)

var (
	POSTGRES_HOST string
	POSTGRES_USER string
	POSTGRES_PASSWORD string
	POSTGRES_DBNAME string
	POSTGRES_PORT string
	TOKEN_AGE int64
	PRIVATE_KEY_PATH string
)

func LoadEnvVars(){
	// [INFO] overwrite existing envs
	err := godotenv.Overload()
	if err != nil {
		log.Println("Error loading .env file")
	}
	SERVER_PORT = os.Getenv("SERVER_PORT")
	POSTGRES_HOST = os.Getenv("POSTGRES_HOST")
	POSTGRES_USER = os.Getenv("POSTGRES_USER")
	POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_DBNAME = os.Getenv("POSTGRES_DBNAME")
	POSTGRES_PORT = os.Getenv("POSTGRES_PORT")
	TOKEN_AGE, err = strconv.ParseInt(os.Getenv("TOKEN_AGE"), 10, 64)
	if err != nil {
        log.Println("Error converting TOKEN_AGE to int64")
    }
	PRIVATE_KEY_PATH = os.Getenv("PRIVATE_KEY_PATH")
}