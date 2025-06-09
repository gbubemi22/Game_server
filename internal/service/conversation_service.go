package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"game_tcpserver/internal/model"
	"game_tcpserver/internal/utils"
)

type ConversationService struct {
	conversationCollection *mongo.Collection
	userCollection         *mongo.Collection
}

func NewConversationService(db *mongo.Database) *ConversationService {
	return &ConversationService{
		conversationCollection: db.Collection("conversation"),
		userCollection:         db.Collection("user"),
	}
}

func (s *ConversationService) CreatePrivateConversation(senderID, receiverID primitive.ObjectID) (*model.Conversation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if senderID == receiverID {
		return nil, utils.NewConflictError("You cannot create a conversation with yourself")
	}

	// Check if private conversation already exists
	filter := bson.M{
		"isGroup": false,
		"participants": bson.M{
			"$all": []primitive.ObjectID{senderID, receiverID},
		},
	}
	var existing model.Conversation
	err := s.conversationCollection.FindOne(ctx, filter).Decode(&existing)
	if err == nil {
		return &existing, nil
	}

	// Create new private conversation
	convo := model.Conversation{
		ID:           primitive.NewObjectID(),
		IsGroup:      false,
		Participants: []primitive.ObjectID{senderID, receiverID},
		CreatedBy:    senderID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err = s.conversationCollection.InsertOne(ctx, convo)
	if err != nil {
		return nil, err
	}
	return &convo, nil
}

type PrivateConv struct {
	CreatorID    primitive.ObjectID   `json:"creatorId" bson:"creatorId"`
	Title        string               `json:"title" bson:"title"`
	Participants []primitive.ObjectID `json:"participants" bson:"participants"`
	//Avatar       *string              `json:"avatar,omitempty" bson:"avatar,omitempty"` // Optional via pointer
}

func (s *ConversationService) CreateGroupConversation(input PrivateConv) (*model.Conversation, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if len(input.Participants) < 2 {
		return nil, utils.NewConflictError("group must have at least 2 participants")
	}

	filter := bson.M{
		"isGroup": true,
		"participants": bson.M{
			"$all": input.Title,
		},
	}

	var existing model.Conversation
	err := s.conversationCollection.FindOne(ctx, filter).Decode(&existing)
	if err == nil {
		return &existing, nil
	}

	convo := model.Conversation{
		ID:           primitive.NewObjectID(),
		Title:        input.Title,
		IsGroup:      true,
		Participants: input.Participants,
		CreatedBy:    input.CreatorID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Insert into DB
	_, err = s.conversationCollection.InsertOne(ctx, convo)
	if err != nil {
		return nil, err
	}

	return &convo, nil
}

type AddParticipantsInput struct {
	ConversationID primitive.ObjectID   `json:"conversationId" bson:"conversationId"`
	UserIDs        []primitive.ObjectID `json:"userIds" bson:"userIds"`
}

func (s *ConversationService) AddUsersToGroupConversation(input AddParticipantsInput) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Find the conversation
	var convo model.Conversation
	err := s.conversationCollection.FindOne(ctx, bson.M{"_id": input.ConversationID}).Decode(&convo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.NewNotFoundError("conversation not found")
		}
		return err
	}

	if !convo.IsGroup {
		return utils.NewConflictError("cannot add users to a private one-on-one conversation")
	}

	// 2. Filter out users who are already participants
	existingParticipantsMap := make(map[string]bool)
	for _, id := range convo.Participants {
		existingParticipantsMap[id.Hex()] = true
	}

	var newParticipants []primitive.ObjectID
	for _, userID := range input.UserIDs {
		if !existingParticipantsMap[userID.Hex()] {
			newParticipants = append(newParticipants, userID)
		}
	}

	if len(newParticipants) == 0 {
		return utils.NewConflictError("all users are already in the group")
	}

	// 3. Append new users
	convo.Participants = append(convo.Participants, newParticipants...)

	// 4. Update the conversation in DB
	update := bson.M{
		"$set": bson.M{
			"participants": convo.Participants,
			"updatedAt":    time.Now(),
		},
	}

	result, err := s.conversationCollection.UpdateOne(ctx, bson.M{"_id": convo.ID}, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return utils.NewConflictError("no changes were made")
	}

	return nil
}
