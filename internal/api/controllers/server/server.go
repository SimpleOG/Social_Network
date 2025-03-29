package server

import (
	"github.com/SimpleOG/Social_Network/internal/api/controllers/ChatControllers/pool"
	"github.com/SimpleOG/Social_Network/internal/api/controllers/MediaControllers"
	"github.com/SimpleOG/Social_Network/internal/api/controllers/userControllers"
	"github.com/SimpleOG/Social_Network/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
)

type Controllers struct {
	AuthHandlers  userControllers.AuthHandlersInterface
	PoolHandlers  pool.PoolHandlersInterface
	MediaHandlers MediaControllers.MediaControllersInteface
}

func NewControllers(service service.Service, upgrader *websocket.Upgrader) Controllers {
	return Controllers{
		AuthHandlers:  userControllers.NewAuthHandlers(service),
		PoolHandlers:  pool.NewPoolHandlers(service, upgrader),
		MediaHandlers: MediaControllers.NewMediaControllers(service),
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
	s.InitAuthRoutes()
	s.InitChatRoutes()
	s.InitMediaRoutes()

}
func (s *Server) InitAuthRoutes() {
	auth := s.router.Group("/auth")
	{
		auth.POST("/sign_in", s.controllers.AuthHandlers.SingIn)
		auth.POST("/login", s.controllers.AuthHandlers.Login)
		auth.GET("/validate", s.controllers.AuthHandlers.RequireAuth, s.controllers.AuthHandlers.Validate)
		auth.POST("logout")
	}
}

func (s *Server) InitChatRoutes() {
	chat := s.router.Group("/chat", s.controllers.AuthHandlers.RequireAuth)
	{
		chat.GET("/createChat", s.controllers.PoolHandlers.ServeRoomsConnections)
	}
}
func (s *Server) InitMediaRoutes() {
	media := s.router.Group("/media") //, s.controllers.AuthHandlers.RequireAuth
	{
		media.POST("/upload", s.controllers.MediaHandlers.UploadImage)
		media.POST("/download", s.controllers.MediaHandlers.DownloadImage)
	}
}
