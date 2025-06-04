package bootstrap

import (
	"log"

	"gorm.io/gorm"
	models_auth "pashmak.com/pashmak/models/auth"
	models "pashmak.com/pashmak/models/openai"
	models_place "pashmak.com/pashmak/models/place"
	models_report "pashmak.com/pashmak/models/report"
)

func getModels() []interface{} {
	// [INFO] add your model here to be migrated
	all_models := []interface{}{
		// authentication
		&models_auth.User{},
		&models_auth.JWTBlacklist{},
		&models_place.Comment{},
		&models_place.Place{},
		&models_place.Reaction{},
		&models_report.Report{},
		&models.SearchHistory{},
	}
	return all_models
}

func MakeMigrations(db *gorm.DB) {
	all_models := getModels()
	for _, model := range all_models {
		if err := db.AutoMigrate(model); err != nil {
			log.Println("Error migrating model: ", err.Error())
		}
	}
	// INFO: Database indexing for efficiency of comments querying
	// db.Exec("CREATE INDEX IF NOT EXISTS idx_comments_place_id ON comments (place_id)")
}
