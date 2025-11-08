package user

import (
	"time"

	"github.com/google/uuid"
)

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
    ImageUrl        string    `gorm:"type:varchar(255)"`

    CreatedAt time.Time
    UpdatedAt time.Time
}

func (User) TableName() string {
    return "users"
}

