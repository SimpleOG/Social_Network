package pool

import (
	"github.com/SimpleOG/Social_Network/internal/service"
	"github.com/gin-gonic/gin"
)

type RoomHandlersInterface interface {
	ServeRoom(ctx *gin.Context)
}

type RoomHandlers struct {
	Service service.Service
}

func NewRoomHandlers() RoomHandlersInterface {
	return &RoomHandlers{}
}
func (r *RoomHandlers) ServeRoom(ctx *gin.Context) {

}
