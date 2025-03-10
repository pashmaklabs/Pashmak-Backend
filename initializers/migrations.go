package initializers

import (
	"gorm.io/gorm"
	"pashmak.com/pashmak/models"
)

func getModels() []interface{}{
	// [INFO] add your model here to be migrated
	all_models := []interface{}{
		// authentication
		&authentication.User{},
		&authentication.UserOTP{},
	}
	return all_models
}

func MakeMigrations(db *gorm.DB) {
	all_models := getModels()
	for _, model := range all_models{
		db.AutoMigrate(model)
	}
}