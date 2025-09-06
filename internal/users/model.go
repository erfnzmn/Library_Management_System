package users

import (
	"time"

	"gorm.io/gorm"
)
const (
	RoleMember  = "member"
	RoleStudent = "student"
)

func IsValidRole(role string) bool {
	return role == RoleMember || role == RoleStudent
}


// مدل دیتابیسی کاربر
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"size:100;not null" json:"name"`
	Email        string         `gorm:"size:190;not null;uniqueIndex" json:"email"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`                
	Role         string         `gorm:"size:20;not null;default:member" json:"role"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string { return "users" }

// ورودیِ ثبت‌نام
type SignupRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role" binding:"required"` // member | student

}

// ورودیِ لاگین
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
