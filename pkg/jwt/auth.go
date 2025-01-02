package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Authorization interface {
	CreateHashPass(password string) (string, error)
	GenerateToken(password, secretKey string) (string, error)
}

type Auth struct {
}

func (a Auth) CreateHashPass(password string) (string, error) {

	hashPasswordString, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}

	return string(hashPasswordString), nil
}
func (a Auth) GenerateToken(password, secretKey string) (string, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(password)); err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "id",
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	//получаем полный токен пользователя
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, err
}
