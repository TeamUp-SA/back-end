package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Bulletin struct {
	BulletinID primitive.ObjectID `json:"bulletinID"`
	AuthorID   primitive.ObjectID `json:"authorID"`
	Title      string             `json:"title"`
	Content    string             `json:"content"`
	GroupID    primitive.ObjectID `json:"groupID,omitempty"`
	CreatedAt  time.Time          `json:"createdAt"`
}

type BulletinCreateRequest struct {
	AuthorID primitive.ObjectID `json:"authorID"`
	Title    string             `json:"title"`
	Content  string             `json:"content"`
	GroupID  primitive.ObjectID `json:"groupID,omitempty"`
}

type BulletinUpdateRequest struct {
	Title   string             `json:"title"`
	Content string             `json:"content"`
	GroupID primitive.ObjectID `json:"groupID,omitempty"`
}
