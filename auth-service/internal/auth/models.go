package auth

import (
	"time"

	"github.com/google/uuid"
)

// Users table in TeamUp DB
type User struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Username        string    `gorm:"type:varchar(50);unique;not null"`
	Name            string    `gorm:"type:varchar(100);not null"`
	Lastname        string    `gorm:"type:varchar(100);not null"`
	PhoneNumber     string    `gorm:"type:varchar(20)"`
	Email           string    `gorm:"type:varchar(255);unique;not null"`
	Password        string    `gorm:"type:text"`
	OAuthProvider   string    `gorm:"type:varchar(50)"`
	OAuthProviderID string    `gorm:"type:varchar(255)"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// Override table name so GORM uses lowercase `users`
func (User) TableName() string {
	return "users"
}
