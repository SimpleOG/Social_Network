package client

import (
	"encoding/json"
	"fmt"
	db "github.com/SimpleOG/Social_Network/internal/repositories/database/postgresql/sqlc"
	"github.com/SimpleOG/Social_Network/internal/repositories/redis"
	"github.com/gorilla/websocket"
	"log"
)

type ClientInterface interface {
	Write()
	Read()
}

type Client struct {
	MsgChan  chan string
	Socket   *websocket.Conn
	UserInfo *db.User
	RClient  *redis.RedisStore
	Done     chan int32
}
type ClientMessage struct {
	MessageOwner int32  `json:"message_owner"`
	Content      string `json:"content"`
}

func NewClient(user *db.User, ws *websocket.Conn, done chan int32, redis *redis.RedisStore) *Client {
	return &Client{
		MsgChan:  make(chan string, 1024),
		Socket:   ws,
		UserInfo: user,
		Done:     done,
		RClient:  redis,
	}
}

// написать В КАНАЛ юзера
func (c *Client) Write() {
	defer c.Socket.Close()

	for {

		select {
		case message, ok := <-c.MsgChan:
			if !ok {
				continue
			}
			// Отправляем сообщение через вебсокет
			if err := c.Socket.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("Ошибка неожиданного закрытия", err)
					return
				} else if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					log.Println("Нормальное завершение  ", err)
					return
				}
				log.Println(err)
			}
		}

	}

}

// ПРОЧИТАТЬ что юзер написал
func (c *Client) Read(redisChan string) {
	defer c.Socket.Close()
	for {
		if c.UserInfo.ID == 1 {
			fmt.Println("Я читаю  че юзер пишет")
		}
		_, msg, err := c.Socket.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Ошибка неожиданного закрытия", err)
				return
			} else if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Println("Нормальное завершение  ", err)
				return
			}
			log.Println(err)

		}

		clientMsg, err := json.Marshal(ClientMessage{
			MessageOwner: c.UserInfo.ID,
			Content:      string(msg),
		})

		if err != nil {
			log.Println("Ошибка маршалинга :", err.Error())
			continue
		}
		err = c.RClient.SendMsgToChan(redisChan, clientMsg)

		if err != nil {
			log.Println(err)
		}

	}

}
