package bootstrap

import (
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseDSN struct {
	host		string
	user		string
	password 	string
	dbname		string
	port		string
}

func SetUpPostgres() *gorm.DB {
	postgres_dsn := DatabaseDSN{
		host: POSTGRES_HOST, 
		user: POSTGRES_USER,
		password: POSTGRES_PASSWORD,
		dbname: POSTGRES_DBNAME,
		port: POSTGRES_PORT,
	}	
	var dsn string = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tehran",
									postgres_dsn.host, postgres_dsn.user, postgres_dsn.password, postgres_dsn.dbname, postgres_dsn.port)
	db, err := gorm.Open(postgres.Open(dsn),  &gorm.Config{})
	if err != nil {
		// TODO: set logger instead of Println
		log.Println("failed to initialize database", err.Error())
	}
	return db
}


func SetUpPGVector() *gorm.DB {
	postgres_dsn := DatabaseDSN{
		host: PGVECTOR_HOST, 
		user: POSTGRES_USER,
		password: POSTGRES_PASSWORD,
		dbname: PGVECTOR_DBNAME,
		port: PGVECTOR_PORT,
	}	
	var dsn string = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tehran",
									postgres_dsn.host, postgres_dsn.user, postgres_dsn.password, postgres_dsn.dbname, postgres_dsn.port)
	db, err := gorm.Open(postgres.Open(dsn),  &gorm.Config{})
	if err != nil {
		// TODO: set logger instead of Println
		log.Println("failed to initialize database", err.Error())
	}
	return db
}

func SetupRedis() *redis.Client{
	var RedisClient = redis.NewClient(&redis.Options{
        Addr:	  fmt.Sprintf("%s:%s", REDIS_HOST, REDIS_PORT),
        Password: REDIS_PASSWORD,
        DB:		  0,
        Protocol: 2,
    })
	return RedisClient
}
