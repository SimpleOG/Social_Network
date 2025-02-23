package pool

import (
	"context"
	db "github.com/SimpleOG/Social_Network/internal/repositories/database/postgresql/sqlc"
	"github.com/SimpleOG/Social_Network/internal/repositories/redis"
	"github.com/SimpleOG/Social_Network/internal/service/chatService/room"
	uuid2 "github.com/google/uuid"
	"log"
)

type PoolInterface interface {
	//Сразу добавляем в бд и всю хуйню
	CreateRoom([2]int32) (*room.Room, error)
	StartPools()
	CheckIfRoomExists(arr []int32) (*room.Room, bool)
}

// ключ это юзеры, значение рума для них
type Pool struct {
	Rooms       map[string]*room.Room
	querier     db.Querier
	redisClient *redis.RedisStore
	addRooms    chan *room.Room
}

func NewPool(q db.Querier, store *redis.RedisStore) *Pool {
	return &Pool{
		Rooms:       make(map[string]*room.Room, 0),
		querier:     q,
		redisClient: store,
		addRooms:    make(chan *room.Room, 1024),
	}
}
func (p *Pool) CheckIfRoomExists(arr []int32) (*room.Room, bool) {
	arg := db.GetRoomByUsersParams{
		User1: arr[0],
		User2: arr[1],
	}
	roomID, err := p.querier.GetRoomByUsers(context.Background(), arg)
	if err != nil {
		return nil, false
	}
	currentRoom := p.Rooms[roomID]
	log.Printf("Комната для %v, %v найдена", arg.User1, arg.User2)
	return currentRoom, true
}
func (p *Pool) CreateRoom(users [2]int32) (*room.Room, error) {
	uuid, err := uuid2.NewUUID()
	if err != nil {
		return nil, err
	}
	arg := db.CreateRoomParams{
		RoomUnique: uuid.String(),
		User1:      users[0],
		User2:      users[1],
	}
	db_room, err := p.querier.CreateRoom(context.Background(), arg)
	if err != nil {
		return nil, err
	}
	Room := room.NewRoom(db_room.RoomUnique, p.redisClient)
	p.Rooms[Room.RoomUUID] = Room
	log.Println("Комната добавляется в пул")
	p.addRooms <- Room
	return Room, nil
}
func (p *Pool) StartPools() {
	log.Println("Пул запущен")
	//запуск для уже существующих комнат
	for _, room := range p.Rooms {
		go room.Run()
	}
	//горутина которая запускает новые каналы
	go func() {
		for room := range p.addRooms {
			log.Println("Новая комната создана")
			go room.Run()
		}
	}()
}
