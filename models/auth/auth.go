package models_auth

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	FirstName  string    `gorm:"default:'کاربر پشمک'"`
	LastName   string    `gorm:"default:''"`
	Email      string    `gorm:"unique;not null"`
	Password   string    `gorm:"default:''"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
	Avatar_url string    `gorm:"default:''"`
	Score      uint
	RoleID     uint
	Role       Role `gorm:"foreignKey:RoleID"` // Foreign key to roles table
}

type JWTBlacklist struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	JTI       string `gorm:"unique;not null"`
	ExpiresAt int64  `gorm:"not null"`
}

// Role represents the roles table
type Role struct {
	gorm.Model               // Includes ID, CreatedAt, UpdatedAt, DeletedAt
	Name        string       `gorm:"type:varchar(50);not null;unique"`
	Permissions []Permission `gorm:"many2many:role_permissions;"` // Many-to-many with permissions
}

// Permission represents the permissions table
type Permission struct {
	gorm.Model        // Includes ID, CreatedAt, UpdatedAt, DeletedAt
	Name       string `gorm:"type:varchar(50);not null;unique"`
}

// RolePermission represents the role_permissions join table (optional, for custom handling)
type RolePermission struct {
	RoleID       uint `gorm:"primaryKey"`
	PermissionID uint `gorm:"primaryKey"`
}
