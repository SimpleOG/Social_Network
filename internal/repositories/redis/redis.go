package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type RedisStore struct {
	RedisClient *redis.Client
}
type RedisInterface interface {
	SendMsgToChan(chanName string, msg any) error
}

// client *redis.client
func NewRedisClient() (*RedisStore, error) {

	client := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
	store := RedisStore{RedisClient: client}
	_, err := store.RedisClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &store, err
}

func (r *RedisStore) SendMsgToChan(chanName string, msg any) error {
	err := r.RedisClient.Publish(context.Background(), chanName, msg).Err()
	if err != nil {
		return errors.New(fmt.Sprintf("message %v couldnt be delivered caz %v", msg, err.Error()))
	}
	return nil
}
