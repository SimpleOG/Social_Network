package pool

import (
	"context"
	"errors"
	db "github.com/SimpleOG/Social_Network/internal/repositories/database/postgresql/sqlc"
	"github.com/SimpleOG/Social_Network/internal/service"
	"github.com/SimpleOG/Social_Network/internal/service/chatService/client"
	room2 "github.com/SimpleOG/Social_Network/internal/service/chatService/room"
	"github.com/SimpleOG/Social_Network/pkg/util/httpResponse"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

type PoolHandlersInterface interface {
	ServeRoomsConnections(ctx *gin.Context)
	ServePools()
}

func (p *Pool) ServePools() {
	p.service.Pool.StartPools()
}

type Pool struct {
	service  service.Service
	upgrader *websocket.Upgrader
}

func NewPoolHandlers(service service.Service, upgrader *websocket.Upgrader) PoolHandlersInterface {
	return &Pool{
		service:  service,
		upgrader: upgrader,
	}
}
func (p *Pool) GetIdFromSet(user any) (db.User, error) {
	//получаем текущего пользователя( ужасный процесс)
	current_id, ok := user.(int32)
	if !ok {
		return db.User{}, errors.New("ошибка обработки")
	}
	currentUser, err := p.service.Querier.GetUsersById(context.Background(), current_id)
	if err != nil {
		return db.User{}, err
	}
	return currentUser, nil
}
func (p *Pool) ServeRoomsConnections(ctx *gin.Context) {
	//чекаем что в принципе указано с кем хотим попиздеть
	query := ctx.Query("id")
	if len(query) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "partners id must be specified"})
		return
	}
	//берем юзера из бд
	user, ok := ctx.Get("id")
	if !ok {
		ctx.JSON(http.StatusUnprocessableEntity, "Пользователь не найден")
	}
	CurrentUser, err := p.GetIdFromSet(user)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, httpResponse.ErrorResponse(err))
	}
	//чекаем что собеседник  есть в системе
	id, err := strconv.Atoi(query)
	conversator, err := p.service.Querier.GetUsersById(context.Background(), int32(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, httpResponse.ErrorResponse(err))
		return
	}
	p.upgrader.CheckOrigin = func(c *http.Request) bool { return true }
	//если есть в системе то открываем коннект
	ws, err := p.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpResponse.ErrorResponse(err))
	}
	//Создаем клиента

	CurrentClient := client.NewClient(&CurrentUser, ws)
	//Проверяем есть ли комната для текущих юзеров в бд
	//Я кстати рот ебал, из за того что создаются в ифах комнаты их вне ифов не видно
	var room *room2.Room
	//Если комната есть то добавляем юзера в неё и пусть пиздят
	room, ok = p.service.Pool.CheckIfRoomExists([]int32{CurrentUser.ID, conversator.ID})
	if ok {
		CurrentClient.RClient = room.RedisClient
		room.Clients[CurrentClient] = struct{}{}
	}
	//Если комнаты нет , то генерируем комнату
	if !ok {
		//Сначала добавляем руму в бд с инфой про uuid и про юзеров которые в ней должны быть

		room, err = p.service.Pool.CreateRoom([2]int32{CurrentUser.ID, conversator.ID})
		if err != nil {
			return
		}
		//прокидываем редис до клиента
		CurrentClient.RClient = room.RedisClient
		//Докидываем в руму пользователей
		room.Clients[CurrentClient] = struct{}{}
	}
	//Теперь пусть смски читаются
	go CurrentClient.Write()
	CurrentClient.Read(room.RoomUUID)

	defer func() {
		room.DeleteClientChan <- CurrentClient
		ws.Close()
	}()

}
