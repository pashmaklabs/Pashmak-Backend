package bootstrap

import (
	"log"

	"gorm.io/gorm"
	"pashmak.com/pashmak/models/auth"
	"pashmak.com/pashmak/models/comment"
)

func getModels() []interface{}{
	// [INFO] add your model here to be migrated
	all_models := []interface{}{
		// authentication
		&models_auth.User{},
		&models_auth.JWTBlacklist{},
		&models_comment.Comment{},
	}
	return all_models
}

func MakeMigrations(db *gorm.DB) {
	all_models := getModels()
	for _, model := range all_models{
		if err := db.AutoMigrate(model); err != nil {
			log.Println("Error migrating model: ", err.Error())
		}
	}
}