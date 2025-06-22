package tcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"game_tcpserver/internal/game"
	"game_tcpserver/internal/model"
	"game_tcpserver/internal/service"
)

type Dependencies struct {
	ConversationService *service.ConversationService
	MessageService      *service.MessageService
	Conn                net.Conn // Add net.Conn to Dependencies
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
		// Pass the connection to the dependencies for use in game-related functions
		deps.Conn = conn
		response := handleCommand(msg, deps) // Assuming handleCommand doesn't need conn directly
		response = handleSendMessage(msg, deps)
		response = handleRoom(msg, deps)


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

type TCPMessage struct {
	Type           string `json:"type"`
	SenderID       string `json:"senderId,omitempty"`
	ConversationID string `json:"conversationId,omitempty"`
	Content        string `json:"content,omitempty"`
}

func handleSendMessage(msg string, deps Dependencies) string {
	var cmd TCPMessage
	if err := json.Unmarshal([]byte(msg), &cmd); err != nil {
		return "Invalid JSON format"
	}

	switch cmd.Type {
	case "send_message":
		if cmd.SenderID == "" || cmd.ConversationID == "" || cmd.Content == "" {
			return "Missing senderId, conversationId, or content"
		}

		senderID, err1 := primitive.ObjectIDFromHex(cmd.SenderID)
		convID, err2 := primitive.ObjectIDFromHex(cmd.ConversationID)
		if err1 != nil || err2 != nil {
			return "Invalid ObjectID(s)"
		}

		message := model.Message{
			SenderID:       senderID,
			ConversationID: convID,
			Content:        cmd.Content,
		}

		saved, err := deps.MessageService.CreateMessage(message)
		if err != nil {
			return "Error saving message: " + err.Error()
		}

		return "Message saved with ID: " + saved.ID.Hex()

	default:
		return "Unknown command"
	}
}

type TCPRoom struct {
	Type       string `json:"type"`
	RoomID     string `json:"roomId,omitempty"`
	PlayerID   string `json:"playerId,omitempty"`
	PlayerName string `json:"playerName,omitempty"`
}

func handleRoom(msg string, deps Dependencies) string {
	var cmd TCPRoom
	if err := json.Unmarshal([]byte(msg), &cmd); err != nil {
		return "Invalid JSON format"
	}

	switch cmd.Type {
	case "create_room":

		if cmd.RoomID == "" || cmd.PlayerID == "" || cmd.PlayerName == "" {
			return "Missing roomId, playerId, or playerName"
		}

		player := &game.Player{
			ID:     cmd.PlayerID,
			Name:   cmd.PlayerName,
			Health: 100,
			Conn:   deps.Conn, // <-- you must pass conn in your Dependencies
		}

		game.AddPlayerToRoom(cmd.RoomID, player)

		return "Room created and player joined: " + cmd.RoomID

	default:
		return "Unknown command"
	}

}
