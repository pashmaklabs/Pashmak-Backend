package bootstrap

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	models_auth "pashmak.com/pashmak/models/auth"
)

func SetupRoleAndPermissions(db *gorm.DB) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // 1. Create permissions
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

        // 2. Create roles with specific IDs
        roles := []struct {
            id          uint
            name        string
            permissions []models_auth.Permission
        }{
            {id: 1, name: "Admin", permissions: permissions},                    // Admin: ID 1, both permissions
            {id: 2, name: "Editor", permissions: []models_auth.Permission{permissions[0]}}, // Editor: ID 2, only "create_post"
            {id: 10, name: "User", permissions: []models_auth.Permission{}},     // User: ID 10, no permissions
        }

        for _, r := range roles {
            role := models_auth.Role{
                Model: gorm.Model{ID: r.id},
                Name:  r.name,
            }
            // Create or find the role
            if err := tx.Where("name = ?", role.Name).FirstOrCreate(&role).Error; err != nil {
                log.Printf("Failed to create role %s: %v", role.Name, err)
                return err
            }

            // Ensure the role has the correct ID
            if role.ID != r.id {
                // If the role exists with a different ID, update it
                if err := tx.Model(&role).Where("name = ?", role.Name).Update("id", r.id).Error; err != nil {
                    log.Printf("Failed to update ID for role %s: %v", role.Name, err)
                    return err
                }
                role.ID = r.id // Update the in-memory struct
            }

            // Explicitly clear existing associations and set new ones
            if err := tx.Model(&role).Association("Permissions").Clear(); err != nil {
                log.Printf("Failed to clear permissions for role %s: %v", role.Name, err)
                return err
            }
            if len(r.permissions) > 0 {
                if err := tx.Model(&role).Association("Permissions").Append(r.permissions); err != nil {
                    log.Printf("Failed to assign permissions to role %s: %v", role.Name, err)
                    return err
                }
            }
        }

        // 3. Create admin user
        if ADMIN_PASSWORD == "" || ADMIN_EMAIL == "" {
            log.Println("Warning: ADMIN_EMAIL or ADMIN_PASSWORD is empty")
            return nil // Optionally return an error if these are required
        }

        adminHashedPass, err := bcrypt.GenerateFromPassword([]byte(ADMIN_PASSWORD), 10)
        if err != nil {
            log.Printf("Failed to hash admin password: %v", err)
            return err
        }

        user := models_auth.User{
            FirstName:  "Admin",
            LastName:   "",
            Email:      ADMIN_EMAIL,
            Password:   string(adminHashedPass),
            Avatar_url: "",
            Score:      0,
            RoleID:     1, // Admin role ID
        }
        if err := tx.Where("email = ?", user.Email).FirstOrCreate(&user).Error; err != nil {
            log.Printf("Failed to create user %s: %v", user.Email, err)
            return err
        }

        return nil
    })
}