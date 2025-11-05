package model

import (
	"time"

	"github.com/Ntchah/TeamUp-application-service/internal/enum/grouptag"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Bulletin struct {
	BulletinID  primitive.ObjectID   `json:"bulletinID" bson:"_id,omitempty"`
	AuthorID    primitive.ObjectID   `json:"authorID" bson:"author_id"`
	Title       string               `json:"title" bson:"title"`
	Description string               `json:"description" bson:"description"`
	GroupID     []primitive.ObjectID `json:"groupID,omitempty" bson:"group_id,omitempty"`
	Date        time.Time            `json:"date" bson:"date"`
	Image       string               `json:"image" bson:"image"`
	Tags        []grouptag.GroupTag  `json:"tags" bson:"tags"`
	CreatedAt   time.Time            `json:"createdAt" bson:"createdAt"`
}
