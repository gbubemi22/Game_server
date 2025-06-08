package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"game_tcpserver/internal/utils"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // execute all handlers

		// Check for errors returned by handlers
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				// Check if it's our CustomError
				if customErr, ok := err.Err.(*utils.CustomError); ok {
					c.JSON(customErr.HTTPStatusCode, customErr)
					return
				}
			}

			// If not our CustomError, fallback
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal Server Error",
				"success": false,
			})
		}
	}
}
