package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Member struct {
	MemberID    primitive.ObjectID `json:"memberID,omitempty"`
	Username    string             `json:"username"`
	Password    string             `json:"password,omitempty"`
	FirstName   string             `json:"firstName"`
	LastName    string             `json:"lastName"`
	Email       string             `json:"email"`
	PhoneNumber string             `json:"phoneNumber,omitempty"`

	// Portfolio Section
	Bio          string   `json:"bio,omitempty"`
	Skills       []string `json:"skills,omitempty"`
	LinkedIn     string   `json:"linkedIn,omitempty"`
	GitHub       string   `json:"github,omitempty"`
	Website      string   `json:"website,omitempty"`
	ProfileImage string   `json:"profileImage,omitempty"`

	// Experience Section
	Experience []Experience `json:"experience,omitempty"`
	Education  []Education  `json:"education,omitempty"`
}

type MemberRegisterRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password,omitempty"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
}

type MemberUpdateRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password,omitempty"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber,omitempty"`

	// Portfolio Section
	Bio          string   `json:"bio,omitempty"`
	Skills       []string `json:"skills,omitempty"`
	LinkedIn     string   `json:"linkedIn,omitempty"`
	GitHub       string   `json:"github,omitempty"`
	Website      string   `json:"website,omitempty"`
	ProfileImage string   `json:"profileImage,omitempty"`

	// Experience Section
	Experience []Experience `json:"experience,omitempty"`
	Education  []Education  `json:"education,omitempty"`

}

// Work / project experience
type Experience struct {
	Title       string `json:"title" bson:"title"`     // e.g. "Frontend Developer Intern"
	Company     string `json:"company" bson:"company"` // or "ABC Inc."
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	StartYear   int    `json:"startYear,omitempty" bson:"startYear,omitempty"`
	EndYear     int    `json:"endYear,omitempty" bson:"endYear,omitempty"`
}

// Education info
type Education struct {
	School    string `json:"school" bson:"school"`
	Degree    string `json:"degree,omitempty" bson:"degree,omitempty"`
	Field     string `json:"field,omitempty" bson:"field,omitempty"`
	StartYear int    `json:"startYear,omitempty" bson:"startYear,omitempty"`
	EndYear   int    `json:"endYear,omitempty" bson:"endYear,omitempty"`
}
