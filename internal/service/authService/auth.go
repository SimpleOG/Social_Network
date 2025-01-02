package authService

import (
	"context"
	"errors"
	db "github.com/SimpleOG/Social_Network/internal/repositories/database/postgresql/sqlc"
	"github.com/SimpleOG/Social_Network/pkg/jwt"
)

type AuthorizationServiceInterface interface {
	SingIn(userParams db.CreateUserParams) (db.User, error)
	Login(params db.GetUserByDataParams) (db.User, error)
}
type AuthService struct {
	querier       db.Querier
	authInterface jwt.Authorization
}

// Регистрация пользователя
func (a *AuthService) SingIn(userParams db.CreateUserParams) (db.User, error) {
	//проверка что пользователя нет в системе
	arg := db.GetUserByDataParams{
		Username: userParams.Username,
		Password: userParams.Password,
	}
	user, err := a.querier.GetUserByData(context.Background(), arg)
	if err != nil {
		return db.User{}, nil
	}
	if user.ID != 0 {
		return db.User{}, errors.New("user already exists")
	}
	//создание пользователя
	hashedPassword, err := a.authInterface.CreateHashPass(userParams.Password)
	if err != nil {
		return db.User{}, err
	}
	userParams.Password = hashedPassword
	NewUser, err := a.querier.CreateUser(context.Background(), userParams)
	if err != nil {
		return db.User{}, err
	}
	return NewUser, nil
}
func (a *AuthService) Login(params db.GetUserByDataParams, SecretKey string) (string, error) {
	user, err := a.querier.GetUserByData(context.Background(), params)
	if err != nil {
		return "", err
	}
	token, err := a.authInterface.GenerateToken(user.Password, SecretKey)
	if err != nil {
		return "", err
	}
	return token, nil
}
