package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"

	"game_tcpserver/internal/database"
	"game_tcpserver/internal/service"
	"game_tcpserver/internal/tcp"
)

type Server struct {
	port int
	db   *mongo.Database
	//ws   *websocket.WebSocketServer
}

func NewServer() *http.Server {
	portStr := os.Getenv("PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 5001
	}

	db, err := database.New()
	if err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		os.Exit(1)
	}

	// Initialize your services
	conversationService := service.NewConversationService(db)

	// Start TCP server
	go tcp.StartTCPServer(tcp.Dependencies{
		ConversationService: conversationService,
	})

	newServer := &Server{
		port: port,
		db:   db,
		//ws:   ws,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", newServer.port),
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
