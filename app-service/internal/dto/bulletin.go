package dto

import (
	"time"

	"github.com/Ntchah/TeamUp-application-service/internal/enum/grouptag"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Bulletin struct {
	BulletinID  primitive.ObjectID   `json:"bulletinID"`
	AuthorID    primitive.ObjectID   `json:"authorID"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	GroupID     []primitive.ObjectID `json:"groupID,omitempty"`
	Date        string            `json:"date"`
	Image       string               `json:"image"`
	Tags        []grouptag.GroupTag  `json:"tags"`
	CreatedAt   time.Time            `json:"createdAt"`
}

type BulletinCreateRequest struct {
	AuthorID    primitive.ObjectID   `json:"authorID"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	GroupID     []primitive.ObjectID `json:"groupID,omitempty"`
	Date        string            `json:"date"`
	Image       string               `json:"image"`
	Tags        []grouptag.GroupTag  `json:"tags"`
}

type BulletinUpdateRequest struct {
	Title       *string               `json:"title,omitempty"`
	Description *string               `json:"description,omitempty"`
	GroupID     *[]primitive.ObjectID `json:"groupID,omitempty"`
	Date        *string               `json:"date,omitempty"`
	Image       *string               `json:"image,omitempty"`
	Tags        *[]grouptag.GroupTag  `json:"tags,omitempty"`
}
