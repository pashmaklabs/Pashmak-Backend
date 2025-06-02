package bootstrap

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	models_auth "pashmak.com/pashmak/models/auth"
)

func SetupRoleAndPermissions(db *gorm.DB){
	db.Transaction(func(tx *gorm.DB) error {
        // 1. Create permissions first
        permissions := []models_auth.Permission{
            {Name: "create_post"},
            {Name: "view_reports"},
        }
        for i, perm := range permissions {
            if err := tx.Where("name = ?", perm.Name).FirstOrCreate(&permissions[i]).Error; err != nil {
                log.Printf("Failed to create permission %s: %v", perm.Name, err)
                return err
            }
        }

        // 2. Create roles and associate permissions
        adminRole := models_auth.Role{
            Name: "Admin",
            Permissions: permissions, // Assign both permissions
        }
        if err := tx.Where("name = ?", adminRole.Name).FirstOrCreate(&adminRole).Error; err != nil {
            log.Printf("Failed to create role %s: %v", adminRole.Name, err)
            return err
        }

        editorRole := models_auth.Role{
            Name: "Editor",
            Permissions: []models_auth.Permission{permissions[0]}, // Only "create_post"
        }
        if err := tx.Where("name = ?", editorRole.Name).FirstOrCreate(&editorRole).Error; err != nil {
            log.Printf("Failed to create role %s: %v", editorRole.Name, err)
            return err
        }

		AdminHashedPass, err := bcrypt.GenerateFromPassword([]byte(ADMIN_PASSWORD), 10)
		if err != nil {
			return err
		}
		
		user := models_auth.User{
            FirstName: "Admin",
            LastName:  "",
            Email:     ADMIN_EMAIL,
            Password:  string(AdminHashedPass),
            Avatar_url: "",
            Score:     0,
            RoleID:    adminRole.ID,
        }
        if err := tx.Where("email = ?", user.Email).FirstOrCreate(&user).Error; err != nil {
            log.Printf("Failed to create user %s: %v", user.Email, err)
            return err
        }

		return nil
    })
}