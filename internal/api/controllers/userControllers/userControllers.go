package userControllers

import (
	db "github.com/SimpleOG/Social_Network/internal/repositories/database/postgresql/sqlc"
	"github.com/SimpleOG/Social_Network/internal/service"
	"github.com/SimpleOG/Social_Network/pkg/util/httpResponse"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type AuthHandlersInterface interface {
	SingIn(ctx *gin.Context)
	Login(ctx *gin.Context)
	RequireAuth(ctx *gin.Context)
	Validate(ctx *gin.Context)
	Logout(ctx *gin.Context)
}

type AuthHandlers struct {
	service service.Service
}

func NewAuthHandlers(service service.Service) AuthHandlersInterface {
	return AuthHandlers{
		service: service,
	}
}

// Реализация эндпоинта регистрации
func (a AuthHandlers) SingIn(ctx *gin.Context) {
	var userParams db.CreateUserParams
	if err := ctx.ShouldBindJSON(&userParams); err != nil {
		ctx.JSON(http.StatusBadRequest, httpResponse.ErrorResponse(err))
		return
	}
	user, err := a.service.Auth.SingIn(userParams)
	//проверка что пользователя добавило в систему
	if err != nil {
		if err.Error() == "user already exists" {
			ctx.JSON(http.StatusBadRequest, httpResponse.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusBadRequest, httpResponse.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"User data": user})
}
func (a AuthHandlers) Logout(ctx *gin.Context) {
	// Удаляем куки
	ctx.SetCookie("Authorization", "", -1, "/", "localhost", false, true)

	ctx.JSON(200, gin.H{
		"message": "Сессия завершена",
	})
}
func (a AuthHandlers) Login(ctx *gin.Context) {
	var searchParams db.GetUserForLoginParams
	if err := ctx.ShouldBindJSON(&searchParams); err != nil {
		ctx.JSON(http.StatusBadRequest, httpResponse.ErrorResponse(err))
		return
	}
	token, err := a.service.Auth.Login(searchParams)
	if err != nil {
		ctx.JSON(http.StatusNotFound, httpResponse.ErrorResponse(err))
		return
	}
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", token, 3600*24*30, "", "", false, false)
	ctx.JSON(http.StatusOK, "Успешно")
}

func (a AuthHandlers) Validate(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "i'm logged in"})
}

func (a AuthHandlers) RequireAuth(ctx *gin.Context) {
	log.Println("Я внутри прослойки!!!!")
	//Получение токена из хедера авторизации
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
	//Парсинг токена
	user_id, err := a.service.Auth.ValidateToken(token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err})
	}
	ctx.Set("id", user_id)
	ctx.Next()

}
