package server

import (
	"github.com/SimpleOG/Social_Network/internal/api/controllers/ChatControllers/pool"
	"github.com/SimpleOG/Social_Network/internal/api/controllers/userControllers"
	"github.com/SimpleOG/Social_Network/internal/service"
	"github.com/gorilla/websocket"
	"log"

	"github.com/gin-gonic/gin"
)

type Controllers struct {
	AuthHandlers userControllers.AuthHandlersInterface
	PoolHandlers pool.PoolHandlersInterface
}

func NewControllers(service service.Service, upgrader *websocket.Upgrader) Controllers {
	return Controllers{
		AuthHandlers: userControllers.NewAuthHandlers(service),
		PoolHandlers: pool.NewPoolHandlers(service, upgrader),
	}
}

type Server struct {
	router      *gin.Engine
	controllers Controllers
}

func NewServer(router *gin.Engine, service service.Service, upgrader *websocket.Upgrader) (*Server, error) {
	server := &Server{
		router:      router,
		controllers: NewControllers(service, upgrader),
	}

	server.InitRoutes()
	return server, nil
}
func (s *Server) Start(address string) error {
	go s.controllers.PoolHandlers.ServePools()
	log.Println("Сервер запустился")
	return s.router.Run(address)
}

func (s *Server) InitRoutes() {

	auth := s.router.Group("/auth")
	{
		auth.POST("/sign_in", s.controllers.AuthHandlers.SingIn)
		auth.POST("/login", s.controllers.AuthHandlers.Login)
		auth.GET("/validate", s.controllers.AuthHandlers.RequireAuth, s.controllers.AuthHandlers.Validate)
	}
	chat := s.router.Group("/chat", s.controllers.AuthHandlers.RequireAuth)
	{
		chat.GET("/createChat", s.controllers.PoolHandlers.ServeRoomsConnections)
	}
}
