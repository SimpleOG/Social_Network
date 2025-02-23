package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Authorization interface {
	CreateHashPass(password string) (string, error)
	GenerateToken(password, hashedPassword, username string) (string, error)
	ValidateToken(tokenString string) (jwt.MapClaims, error)
}

type Auth struct {
	SecretKey string `json:"secret_key"`
}

func NewJwtAuth(key string) Auth {
	return Auth{
		SecretKey: key,
	}
}
func (a Auth) CreateHashPass(password string) (string, error) {

	hashPasswordString, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}

	return string(hashPasswordString), nil
}
func (a Auth) GenerateToken(password, hashedPassword, username string) (string, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	//получаем полный токен пользователя
	tokenString, err := token.SignedString([]byte(a.SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, err
}
func (a Auth) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method : %v", token.Header["alg"])
		}
		return []byte(a.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return nil, errors.New("time limit expired")
		}
		return claims, nil
	} else {
		return nil, errors.New("cannot parse claims")
	}
}
