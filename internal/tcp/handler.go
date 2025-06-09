package tcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"game_tcpserver/internal/service"
)

type Dependencies struct {
	ConversationService *service.ConversationService
}

func HandleConnection(conn net.Conn, deps Dependencies) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client disconnected:", conn.RemoteAddr())
			break
		}

		msg = strings.TrimSpace(msg)
		fmt.Printf("[%v] %s\n", conn.RemoteAddr(), msg)

		// Route message to appropriate handler
		response := handleCommand(msg, deps)

		conn.Write([]byte(response + "\n"))
	}
}

type TCPCommand struct {
	Type           string   `json:"type"`
	SenderID       string   `json:"senderId,omitempty"`
	ReceiverID     string   `json:"receiverId,omitempty"`
	Title          string   `json:"title,omitempty"`
	CreatorID      string   `json:"creatorId,omitempty"`
	ParticipantIDs []string `json:"participantIds,omitempty"`
}

func handleCommand(msg string, deps Dependencies) string {
	var cmd TCPCommand
	if err := json.Unmarshal([]byte(msg), &cmd); err != nil {
		return "Invalid JSON format"
	}

	switch cmd.Type {
	case "create_conversation":
		if cmd.SenderID == "" || cmd.ReceiverID == "" {
			return "Missing senderId or receiverId"
		}

		sender, err1 := primitive.ObjectIDFromHex(cmd.SenderID)
		receiver, err2 := primitive.ObjectIDFromHex(cmd.ReceiverID)
		if err1 != nil || err2 != nil {
			return "Invalid sender or receiver ObjectID"
		}

		convo, err := deps.ConversationService.CreatePrivateConversation(sender, receiver)
		if err != nil {
			return fmt.Sprintf("Error creating conversation: %v", err)
		}

		return fmt.Sprintf("Conversation created with ID: %s", convo.ID.Hex())

	case "create_group_conversation":
		if cmd.Title == "" || cmd.CreatorID == "" || len(cmd.ParticipantIDs) == 0 {
			return "Missing title, creatorId, or participantIds"
		}

		creatorID, err := primitive.ObjectIDFromHex(cmd.CreatorID)
		if err != nil {
			return "Invalid creatorId"
		}

		var participantIDs []primitive.ObjectID
		for _, idStr := range cmd.ParticipantIDs {
			objID, err := primitive.ObjectIDFromHex(strings.TrimSpace(idStr))
			if err != nil {
				return "Invalid participant ID: " + idStr
			}
			participantIDs = append(participantIDs, objID)
		}

		input := service.PrivateConv{
			CreatorID:    creatorID,
			Title:        cmd.Title,
			Participants: participantIDs,
		}

		convo, err := deps.ConversationService.CreateGroupConversation(input)
		if err != nil {
			return "Error creating group conversation: " + err.Error()
		}

		return "Group conversation created with ID: " + convo.ID.Hex()

	default:
		return "Unknown command"
	}
}
