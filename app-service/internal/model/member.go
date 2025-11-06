package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Member struct {
	MemberID    primitive.ObjectID `json:"memberID,omitempty" bson:"_id,omitempty"`
	Username    string             `json:"username" bson:"username"`
	Password    string             `json:"password,omitempty" bson:"password,omitempty" copier:"-"`
	FirstName   string             `json:"firstName" bson:"firstName"`
	LastName    string             `json:"lastName" bson:"lastName"`
	Email       string             `json:"email" bson:"email"`
	PhoneNumber string             `json:"phoneNumber,omitempty" bson:"phoneNumber,omitempty"`

	// Portfolio Section
	Bio          string   `json:"bio,omitempty" bson:"bio,omitempty"`               // Short self intro
	Skills       []string `json:"skills,omitempty" bson:"skills,omitempty"`         // e.g. ["React", "Go", "UI/UX"]
	LinkedIn     string   `json:"linkedIn,omitempty" bson:"linkedIn,omitempty"`
	GitHub       string   `json:"github,omitempty" bson:"github,omitempty"`
	Website      string   `json:"website,omitempty" bson:"website,omitempty"`
	ProfileImage string   `json:"profileImage,omitempty" bson:"profileImage,omitempty"`

	// Experience Section
	Experience []Experience `json:"experience,omitempty" bson:"experience,omitempty"`
	Education  []Education  `json:"education,omitempty" bson:"education,omitempty"`
}

// Work / project experience
type Experience struct {
	Title       string `json:"title" bson:"title"`             // e.g. "Frontend Developer Intern"
	Company     string `json:"company" bson:"company"`         // or "ABC Inc."
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	StartYear   int `json:"startDate,omitempty" bson:"startDate,omitempty"`
	EndYear     int `json:"endDate,omitempty" bson:"endDate,omitempty"`
}

// Education info
type Education struct {
	School      string `json:"school" bson:"school"`
	Degree      string `json:"degree,omitempty" bson:"degree,omitempty"`
	Field       string `json:"field,omitempty" bson:"field,omitempty"`
	StartYear   int    `json:"startYear,omitempty" bson:"startYear,omitempty"`
	EndYear     int    `json:"endYear,omitempty" bson:"endYear,omitempty"`
}