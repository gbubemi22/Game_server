package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"game_tcpserver/internal/model"
)

type MessageService struct {
	collection *mongo.Collection
}

func NewMessageService(db *mongo.Database) *MessageService {
	return &MessageService{
		collection: db.Collection("messages"),
	}
}

func (s *MessageService) CreateMessage(msg model.Message) (*model.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	msg.ID = primitive.NewObjectID()
	msg.CreatedAt = time.Now()

	_, err := s.collection.InsertOne(ctx, msg)
	if err != nil {
		return nil, err
	}

	return &msg, nil
}
