package service

import (
	db "github.com/SimpleOG/Social_Network/internal/repositories/database/postgresql/sqlc"
	"github.com/SimpleOG/Social_Network/internal/repositories/redis"
	"github.com/SimpleOG/Social_Network/internal/service/authService"
	"github.com/SimpleOG/Social_Network/internal/service/chatService/pool"
	"github.com/SimpleOG/Social_Network/pkg/jwt"
)

type Service struct {
	Auth    authService.AuthorizationServiceInterface
	Pool    pool.PoolInterface
	Querier db.Querier
}

func NewService(q db.Querier, auth jwt.Auth, redis *redis.RedisStore) Service {
	return Service{
		authService.NewAuthService(q, auth),
		pool.NewPool(q, redis),
		q,
	}
}
