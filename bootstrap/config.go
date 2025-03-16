package bootstrap

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	SERVER_PORT string
	
	EMAIL_PASSWORD string
	EMAIL_HOST string
	EMAIL_PORT string
	EMAIL_ADDR string
)

var (
	POSTGRES_HOST     string
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_DBNAME   string
	POSTGRES_PORT     string

	REDIS_PORT        string
	REDIS_HOST        string
	REDIS_PASSWORD    string
)

func LoadEnvVars() {
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

	EMAIL_PASSWORD = os.Getenv("EMAIL_PASSWORD")
	EMAIL_HOST = os.Getenv("EMAIL_HOST")
	EMAIL_PORT = os.Getenv("EMAIL_PORT")
	EMAIL_ADDR = os.Getenv("EMAIL_ADDR")

	REDIS_HOST = os.Getenv("REDIS_HOST")
	REDIS_PORT = os.Getenv("REDIS_PORT")
	REDIS_HOST = os.Getenv("REDIS_HOST")
	REDIS_PASSWORD = os.Getenv("REDIS_PASSWORD")
}
