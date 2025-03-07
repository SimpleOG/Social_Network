package room

import (
	"encoding/json"
	db "github.com/SimpleOG/Social_Network/internal/repositories/database/postgresql/sqlc"
	"github.com/SimpleOG/Social_Network/internal/repositories/redis"
	"github.com/SimpleOG/Social_Network/internal/service/chatService/client"
	"golang.org/x/net/context"
	"log"
	"sync"
)

var MaxBuffSize = 1 << 10

type Room struct {
	Clients             map[*client.Client]struct{}
	NewClientChan       chan *client.Client
	DeleteClientChan    chan int32
	msgChan             chan client.ClientMessage
	stopRoom            chan struct{}
	RoomUUID            string
	RedisClient         *redis.RedisStore
	Querier             db.Querier
	Stopped             bool
	UndeliveredMessages chan client.ClientMessage
	mu                  sync.RWMutex
}

func NewRoom(UUID string, redisClient *redis.RedisStore, querier db.Querier) *Room {
	return &Room{
		Clients:             make(map[*client.Client]struct{}),
		NewClientChan:       make(chan *client.Client, MaxBuffSize),
		DeleteClientChan:    make(chan int32, MaxBuffSize),
		msgChan:             make(chan client.ClientMessage, MaxBuffSize),
		stopRoom:            make(chan struct{}, MaxBuffSize),
		RoomUUID:            UUID,
		RedisClient:         redisClient,
		Stopped:             false,
		Querier:             querier,
		UndeliveredMessages: make(chan client.ClientMessage, 2*MaxBuffSize),
		mu:                  sync.RWMutex{},
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

			var NewMsg client.ClientMessage
			if err := json.Unmarshal([]byte(msg.Payload), &NewMsg); err != nil {
				log.Println("Сообщение не удалось размаршалить: ", err.Error())
			}
			msgArg := db.CreateMessageParams{
				RoomID:         r.RoomUUID,
				MessageContent: NewMsg.Content,
				MessageOwner:   NewMsg.MessageOwner,
				WasDelivered:   func() bool { return len(r.Clients) == 2 }(),
			}

			if err := r.Querier.CreateMessage(context.Background(), msgArg); err != nil {
				log.Println("Возникла ошибка создания сообщения в базе данных :", err.Error())
			}
			r.msgChan <- NewMsg
		}

	}

}

// SendUndeliveredMessages Функция проверяет есть ли сообщения, которые
// были недоставлены
func (r *Room) SendUndeliveredMessages(id int32) {
	arg := db.GetAllUndeliveredMessagesParams{
		RoomID:       r.RoomUUID,
		MessageOwner: id,
	}
	undelivered_messages, err := r.Querier.GetAllUndeliveredMessages(context.Background(), arg)
	if err != nil {
		log.Println(err)
		return
	}
	if len(undelivered_messages) > 0 {
		log.Println("Начат процесс отправки неотправленных сообщений ")
		for _, v := range undelivered_messages {
			msg := client.ClientMessage{
				MessageOwner: v.ID,
				Content:      v.MessageContent,
			}
			for clt := range r.Clients {
				if clt.UserInfo.ID != msg.MessageOwner {
					clt.MsgChan <- msg.Content
					err := r.Querier.ChangeDeliveryTipe(context.Background(), v.ID)
					if err != nil {
						log.Println("Ошибка изменения вида письма на отправленное")
					}
				}
			}
		}
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
			r.mu.Lock()
			r.Clients[NewClient] = struct{}{}
			r.mu.Unlock()
			r.SendUndeliveredMessages(NewClient.UserInfo.ID)
		case _ = <-r.stopRoom:
			break
		case clientId := <-r.DeleteClientChan:
			for i := range r.Clients {
				if i.UserInfo.ID == clientId {
					log.Println("Удаление клиента с id", clientId)
					delete(r.Clients, i)
				}
			}
		case msg := <-r.msgChan:
			//log.Println("В канал пришло новое сообщение:", msg)
			for client := range r.Clients {
				go func() {
					//log.Printf("Сообщение отправлено в канал клиента %v", client.UserInfo.ID)
					if client.UserInfo.ID != msg.MessageOwner {
						client.MsgChan <- msg.Content
					}

				}()
			}

		}
	}
}
