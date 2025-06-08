package server

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"game_tcpserver/internal/controller"
	"game_tcpserver/internal/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(gin.Logger())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	r.GET("/", s.HelloWorldHandler)

	r.Use(middleware.ErrorHandlerMiddleware)
	r.NoRoute(middleware.HandleNotFound)
	r.Use(middleware.ErrorHandler())

	userController := controller.NewUserController()

	r.POST("/v1/auth/users/create", userController.CreateUserHttp)
	r.POST("/v1/auth/users/verify-email", userController.VerifyEmailHandler)

	r.POST("/v1/auth/users/send-email", userController.SendEmailHandler)

	r.POST("/v1/auth/users/login", userController.LoginHandler)

	//authorized := r.Group("/v1/auth")
	// authorized.Use(middleware.VerifyToken())
	// {
	// 	authorized.POST("/users/upload-image", userController.UploadImageHandler)
	// }

	// r.GET("/ws", func(c *gin.Context) {
	// 	s.ws.HandleConnections(c.Writer, c.Request)
	// })

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}
