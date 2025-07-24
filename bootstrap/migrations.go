package bootstrap

import (
	"log"

	"gorm.io/gorm"
	models_auth "pashmak.com/pashmak/models/auth"
	models "pashmak.com/pashmak/models/openai"
	models_pgvector "pashmak.com/pashmak/models/pgvector"
	models_place "pashmak.com/pashmak/models/place"
	models_report "pashmak.com/pashmak/models/report"
)

func getModels() []interface{} {
	// [INFO] add your model here to be migrated
	all_models := []interface{}{
		&models_auth.User{},
		&models_auth.JWTBlacklist{},
		&models_place.Comment{},
		&models_place.Place{},
		&models_place.Reaction{},
		&models_place.PlaceLabel{},
		&models_place.SavedLocation{},
		&models_auth.Role{},
		&models_auth.Permission{},
		&models_auth.RolePermission{},
		&models_report.Report{},
		&models.SearchHistory{},
		&models_place.Image{},
	}
	return all_models
}


func getPGVectorModels() []interface{} {
	// [INFO] add your model here to be migrated
	all_models := []interface{}{
		&models_pgvector.Gplace{},
		&models_pgvector.Greview{},
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

	SetupRoleAndPermissions(db)

	
	// TODO: Database indexing for efficiency of comments querying
}


func MakePGVectorMigrations(db *gorm.DB) {
	all_models := getPGVectorModels()
	for _, model := range all_models {
		if err := db.AutoMigrate(model); err != nil {
			log.Println("Error migrating model: ", err.Error())
		}
	}

	SetupRoleAndPermissions(db)

	
	// TODO: Database indexing for efficiency of comments querying
}
