package userControllers

import (
	"errors"
	"github.com/SimpleOG/Social_Network/internal/api/controllers/server"
	db "github.com/SimpleOG/Social_Network/internal/repositories/database/postgresql/sqlc"
	"github.com/SimpleOG/Social_Network/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthHandlersInterface interface {
	SingIn(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type AuthHandlers struct {
	service service.Service
}

// Реализация эндпоинта регистрации
func (a AuthHandlers) SingIn(ctx *gin.Context) {
	var userParams db.CreateUserParams
	if err := ctx.ShouldBindJSON(&userParams); err != nil {
		ctx.JSON(http.StatusBadRequest, server.ErrorResponse(err))
		return
	}
	user, err := a.service.SingIn(userParams)
	//проверка что пользователя добавило в систему
	if err != nil {
		if errors.Is(err, errors.New("user already exists")) {
			ctx.JSON(http.StatusBadRequest, server.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusBadRequest, server.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"User data": user})
}
func (a AuthHandlers) Login(ctx *gin.Context) {
}
