package service

import (
	db "github.com/SimpleOG/Social_Network/internal/repositories/database/postgresql/sqlc"
	"github.com/SimpleOG/Social_Network/internal/repositories/redis"
	"github.com/SimpleOG/Social_Network/internal/service/authService"
	"github.com/SimpleOG/Social_Network/internal/service/chatService/pool"
	"github.com/SimpleOG/Social_Network/pkg/jwt"
)

type Service struct {
	Auth authService.AuthorizationServiceInterface
	Pool pool.PoolInterface
	//MediaService MediaService.MediaServiceInterface
	Querier db.Querier
}

func NewService(q db.Querier, auth jwt.Auth, redis *redis.RedisStore) Service {
	return Service{
		Auth: authService.NewAuthService(q, auth),
		Pool: pool.NewPool(q, redis),
		//MediaService: MediaService.NewMediaService(q),
		Querier: q,
	}
}
