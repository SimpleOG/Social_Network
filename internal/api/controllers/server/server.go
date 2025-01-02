package server

import (
	"github.com/SimpleOG/Social_Network/internal/service/authService"
	"github.com/gin-gonic/gin"
)

type Server struct {
	authService.AuthorizationServiceInterface
}

func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
