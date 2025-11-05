package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Bulletin struct {
	BulletinID primitive.ObjectID `json:"bulletinID" bson:"_id,omitempty"`
	AuthorID   primitive.ObjectID `json:"authorID" bson:"author_id"`
	Title      string             `json:"title" bson:"title"`
	Content    string             `json:"content" bson:"content"`
	GroupID    primitive.ObjectID `json:"groupID,omitempty" bson:"group_id,omitempty"`
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
}
