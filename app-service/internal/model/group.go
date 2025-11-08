package model

import (
	"time"

	"app-service/internal/enum/grouptag"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Group struct {
	GroupID     primitive.ObjectID   `json:"groupID" bson:"_id,omitempty"`
	Title       string               `json:"title" bson:"title"`
	Description string               `json:"description" bson:"description"`
	OwnerID     primitive.ObjectID   `json:"ownerID" bson:"owner_id"`
	Members     []primitive.ObjectID `json:"members" bson:"members"`
	Tags        []grouptag.GroupTag  `json:"tags" bson:"tags"`
	Closed      bool                 `json:"closed" bson:"closed"`
	Date        string               `json:"date" bson:"date"`
	CreatedAt   time.Time            `json:"createdAt" bson:"createdAt"`
}
