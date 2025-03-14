package bootstrap

import (
	"fmt"
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
  
var DB *gorm.DB

func SetUpPostgres() *gorm.DB {
	var err error
	postgres_dsn := DatabaseDSN{
		host: POSTGRES_HOST, 
		user: POSTGRES_USER,
		password: POSTGRES_PASSWORD,
		dbname: POSTGRES_DBNAME,
		port: POSTGRES_PORT,
	}	
	var dsn string = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tehran",
									postgres_dsn.host, postgres_dsn.user, postgres_dsn.password, postgres_dsn.dbname, postgres_dsn.port)
	DB, err = gorm.Open(postgres.Open(dsn),  &gorm.Config{})
	if err != nil {
		// TODO: set logger instead of Println
		fmt.Println("failed to initialize database")
	}
	return DB
}
