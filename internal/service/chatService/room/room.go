package room

import (
	"github.com/SimpleOG/Social_Network/internal/repositories/redis"
	"github.com/SimpleOG/Social_Network/internal/service/chatService/client"
	"golang.org/x/net/context"
	"log"
)

var MaxBuffSize = 1 << 10

type Room struct {
	Clients          map[*client.Client]struct{}
	NewClientChan    chan *client.Client
	DeleteClientChan chan *client.Client
	msgChan          chan string
	stopRoom         chan struct{}
	RoomUUID         string
	RedisClient      *redis.RedisStore
}

func NewRoom(UUID string, redisClient *redis.RedisStore) *Room {
	return &Room{
		Clients:          make(map[*client.Client]struct{}),
		NewClientChan:    make(chan *client.Client, MaxBuffSize),
		DeleteClientChan: make(chan *client.Client, MaxBuffSize),
		msgChan:          make(chan string, MaxBuffSize),
		stopRoom:         make(chan struct{}, MaxBuffSize),
		RoomUUID:         UUID,
		RedisClient:      redisClient,
	}
}

func (r *Room) ReadChan() {

	pubsub := r.RedisClient.RedisClient.Subscribe(context.Background(), r.RoomUUID)
	ch := pubsub.Channel()
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				continue
			}
			log.Println("Кладу сообщение в канал комнаты!")
			r.msgChan <- msg.Payload
		}

		//msg, err := pubsub.ReceiveMessage(context.Background())

	}

}
func (r *Room) Run() {
	log.Printf("Комната %v запущена", r.RoomUUID)
	go r.ReadChan()
	for {
		if len(r.Clients) == 0 {
			r.stopRoom <- struct{}{}
		}
		select {
		case NewClient := <-r.NewClientChan:
			r.Clients[NewClient] = struct{}{}
		case _ = <-r.stopRoom:
			break
		case client := <-r.DeleteClientChan:
			delete(r.Clients, client)
		case msg := <-r.msgChan:
			log.Println("В канал пришло новое сообщение:", msg)
			for client := range r.Clients {
				go func() {
					log.Printf("Сообщение отправленно в канал клиента %v", client.UserInfo.ID)
					client.MsgChan <- msg
				}()
			}

		}
	}
}
