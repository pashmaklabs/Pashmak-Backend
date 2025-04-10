package bootstrap

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	PashmakApiUrl string
	ServerHost string
	ServerPort string
	
	EmailPassword string
	EmailHost string
	EmailPort string
	EmailAddr string

	PostgresHost     string
	PostgresUser     string
	PostgresPassword string
	PostgresDBName string
	PostgresPort string

	TokenAge int64
	PrivateKeyPath string

	RedisPort        string
	RedisHost        string
	RedisPassword    string

	AllowdOrigins	[]string
	CookieDomain string

	Environment string
}

var (
	PASHMAK_API_URL string
	SERVER_HOST string
	SERVER_PORT string
	
	EMAIL_PASSWORD string
	EMAIL_HOST string
	EMAIL_PORT string
	EMAIL_ADDR string

	POSTGRES_HOST     string
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_DBNAME string
	POSTGRES_PORT string

	TOKEN_AGE int64
	PRIVATE_KEY_PATH string

	REDIS_PORT        string
	REDIS_HOST        string
	REDIS_PASSWORD    string

	CookieDomain string
)

func LoadEnvVars() *AppConfig {
	// [INFO] overwrite existing envs
	err := godotenv.Overload()
	if err != nil {
		log.Println("Error loading .env file: ", err.Error())
	}
	PASHMAK_API_URL = os.Getenv("PASHMAK_API_URL")
	SERVER_HOST = os.Getenv("SERVER_HOST")
	SERVER_PORT = os.Getenv("SERVER_PORT")

	POSTGRES_HOST = os.Getenv("POSTGRES_HOST")
	POSTGRES_USER = os.Getenv("POSTGRES_USER")
	POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_DBNAME = os.Getenv("POSTGRES_DBNAME")
	POSTGRES_PORT = os.Getenv("POSTGRES_PORT")

	TOKEN_AGE, err := strconv.ParseInt(os.Getenv("TOKEN_AGE"), 10, 64)
	TOKEN_AGE++ // Should be removed
	if err != nil {
        log.Println("Error converting TOKEN_AGE to int64: ", err.Error())
    }
	PRIVATE_KEY_PATH = os.Getenv("PRIVATE_KEY_PATH")


	EMAIL_PASSWORD = os.Getenv("EMAIL_PASSWORD")
	EMAIL_HOST = os.Getenv("EMAIL_HOST")
	EMAIL_PORT = os.Getenv("EMAIL_PORT")
	EMAIL_ADDR = os.Getenv("EMAIL_ADDR")

	REDIS_HOST = os.Getenv("REDIS_HOST")
	REDIS_PORT = os.Getenv("REDIS_PORT")
	REDIS_PASSWORD = os.Getenv("REDIS_PASSWORD")

	AllowdOrigins := []string {
		"http://localhost:5174", 
		"http://localhost:5173",
		"http://localhost:8080", 
		"https://pashmak.darkube.app",
	}

	Environment := os.Getenv("GO_ENV")
	if Environment == "development" {
		CookieDomain = "darkube.app"
	}else{
		CookieDomain = ""
	}

	return &AppConfig{
		PashmakApiUrl: PASHMAK_API_URL,
		ServerHost: SERVER_HOST,
		ServerPort: SERVER_PORT,
		EmailPassword: EMAIL_PASSWORD,
		EmailHost: EMAIL_HOST,
		EmailPort: EMAIL_PORT,
		EmailAddr: EMAIL_ADDR,
		PostgresHost: POSTGRES_HOST,
		PostgresUser: POSTGRES_USER,
		PostgresPassword: POSTGRES_PASSWORD,
		PostgresDBName: POSTGRES_DBNAME,
		PostgresPort: POSTGRES_PORT,
		TokenAge: TOKEN_AGE,
		PrivateKeyPath: PRIVATE_KEY_PATH,
		RedisPort: REDIS_PORT,
		RedisHost: REDIS_HOST,
		RedisPassword: REDIS_PASSWORD,
		AllowdOrigins: AllowdOrigins,
		CookieDomain: CookieDomain,
	}
}
