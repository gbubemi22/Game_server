package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Conversation struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Title        string               `bson:"title,omitempty" json:"title,omitempty"`
	IsGroup      bool                 `bson:"isGroup" json:"isGroup"`
	GroupAvatar  string               `bson:"groupAvatar,omitempty" json:"groupAvatar,omitempty"` // Optional image URL
	Participants []primitive.ObjectID `bson:"participants" json:"participants"`
	CreatedBy    primitive.ObjectID   `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	LastMessage  *LastMessage         `bson:"lastMessage,omitempty" json:"lastMessage,omitempty"`
	CreatedAt    time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time            `bson:"updatedAt" json:"updatedAt"`
}

// A snapshot of the last message for faster loading
type LastMessage struct {
	Content   string             `bson:"content" json:"content"`
	SenderID  primitive.ObjectID `bson:"senderId" json:"senderId"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
}
