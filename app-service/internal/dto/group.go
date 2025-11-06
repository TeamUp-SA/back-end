package dto

import (
	"time"

	"github.com/Ntchah/TeamUp-application-service/internal/enum/grouptag"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Group struct {
	GroupID     primitive.ObjectID   `json:"groupID"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	OwnerID     primitive.ObjectID   `json:"ownerID"`
	Members     []primitive.ObjectID `json:"members"`
	Tags        []grouptag.GroupTag  `json:"tags"`
	Closed      bool                 `json:"closed"`
	Date        string            `json:"date"`
	CreatedAt   time.Time            `json:"createdAt"`
}

type GroupCreateRequest struct {
	Title       string               `json:"title"`
	Description string               `json:"description"`
	OwnerID     primitive.ObjectID   `json:"ownerID"`
	Members     []primitive.ObjectID `json:"members"`
	Tags        []grouptag.GroupTag  `json:"tags"`
	Closed      bool                 `json:"closed"`
	Date        string            `json:"date"`
	CreatedAt   time.Time            `json:"createdAt"`
}

type GroupUpdateRequest struct {
	Title       string              `json:"title"`
	Description string              `json:"description"`
	OwnerID     string              `json:"ownerID"`
	Members     []string            `json:"members"`
	Tags        []grouptag.GroupTag `json:"tags"`
	Closed      bool                `json:"closed"`
	Date        string            `json:"date"`
	CreatedAt   time.Time           `json:"createdAt"`
}
