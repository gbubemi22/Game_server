package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID             primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	ConversationID primitive.ObjectID   `bson:"conversationId" json:"conversationId"`
	SenderID       primitive.ObjectID   `bson:"senderId" json:"senderId"`
	Content        string               `bson:"content" json:"content"`
	AttachmentURL  *string              `bson:"attachmentUrl,omitempty" json:"attachmentUrl,omitempty"` // optional
	CreatedAt      time.Time            `bson:"createdAt" json:"createdAt"`
}
