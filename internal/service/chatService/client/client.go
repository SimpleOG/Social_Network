package client

import (
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
}

func NewClient(user *db.User, ws *websocket.Conn) *Client {
	return &Client{
		MsgChan:  make(chan string, 1024),
		Socket:   ws,
		UserInfo: user,
	}
}

// написать В КАНАЛ юзера
func (c *Client) Write() {
	defer c.Socket.Close()

	for {
		fmt.Println("Я жду сообщение")
		select {
		case message, ok := <-c.MsgChan:
			if !ok {
				continue
			}
			log.Println("Вот и сообщения пишутся")
			// Отправляем сообщение через вебсокет
			if err := c.Socket.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				log.Println("Write error:", err)
				continue
			}
		}

	}

}

// ПРОЧИТАТЬ что юзер написал
func (c *Client) Read(redisChan string) {
	defer c.Socket.Close()
	log.Println("Я читаю")
	for {
		_, msg, err := c.Socket.ReadMessage()
		if err != nil {
			log.Println(err)
		}
		log.Println(string(msg))
		err = c.RClient.SendMsgToChan(redisChan, msg)
		if err != nil {
			log.Println(err)
		}
	}

}
