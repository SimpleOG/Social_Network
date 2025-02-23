package authService

import (
	"context"
	"errors"
	db "github.com/SimpleOG/Social_Network/internal/repositories/database/postgresql/sqlc"
	"github.com/SimpleOG/Social_Network/pkg/jwt"
)

type AuthorizationServiceInterface interface {
	SingIn(userParams db.CreateUserParams) (db.User, error)
	Login(userParams db.GetUserForLoginParams) (string, error)
	ValidateToken(token string) (int32, error)
}
type AuthService struct {
	querier       db.Querier
	authInterface jwt.Authorization
}

func NewAuthService(q db.Querier, auth jwt.Authorization) AuthorizationServiceInterface {
	return &AuthService{
		querier:       q,
		authInterface: auth,
	}
}

// Регистрация пользователя
func (a *AuthService) SingIn(userParams db.CreateUserParams) (db.User, error) {
	//Проверяем что пользователя нет в бд для регистрации

	user, err := a.querier.GetUserByUsername(context.Background(), userParams.Username)

	if err != nil {
		if err.Error() != "no rows in result set" {
			return db.User{}, err
		}
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
func (a *AuthService) Login(userParams db.GetUserForLoginParams) (string, error) {
	user, err := a.querier.GetUserByUsername(context.Background(), userParams.Username)
	if err != nil {
		return "", err
	}
	token, err := a.authInterface.GenerateToken(userParams.Password, user.Password, user.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}
func (a *AuthService) ValidateToken(token string) (int32, error) {
	claimsMap, err := a.authInterface.ValidateToken(token)
	if err != nil {
		return 0, err
	}
	//Проверяем наличие пользователя в системе
	user, err := a.querier.GetUserByUsername(context.Background(), claimsMap["sub"].(string))
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}
