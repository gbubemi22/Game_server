package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConversationUserSetting struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ConversationID primitive.ObjectID `bson:"conversationId" json:"conversationId"`
	UserID         primitive.ObjectID `bson:"userId" json:"userId"`
	IsMuted        bool               `bson:"isMuted" json:"isMuted"`
	MutedAt        *time.Time         `bson:"mutedAt,omitempty" json:"mutedAt,omitempty"`
}
