package handler

import (
	"awesomeProject/internal/app/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RegUser
// @Description create user
// @Tags user
// @Produce  json
// @Param user body models.CreateUserRequest true "User info"
// @Success 201 {object} models.CreateUserResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/register [post]
func (h *Handler) RegUser(ctx *gin.Context) {
	var request models.CreateUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Repository.CreateUser(request.Login, request.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, models.CreateUserResponse{
		User: user,
	})
}

// UpdateUser
// @Description update user
// @Tags user
// @Produce  json
// @Param user body models.UpdateUserRequest true "User info"
// @Success 200 {object} models.CreateUserResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/update [put]
func (h *Handler) UpdateUser(ctx *gin.Context) {
	var request models.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Repository.UpdateUser(request.ID, request.Login, request.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, models.CreateUserResponse{
		User: user,
	})
}

// AuthUser
// @Description auth user
// @Tags user
// @Produce  json
// @Param user body models.CreateUserRequest true "User info"
// @Success 200 {object} models.AuthUserRequest
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/login [post]
func (h *Handler) AuthUser(ctx *gin.Context) {
	fmt.Println("логинимс")

	var request models.AuthUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(request.Login, request.Password)
	token, isAdmin, err := h.Repository.AuthUser(request.Login, request.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token":   token,
		"login":   request.Login,
		"isAdmin": isAdmin,
	})
}

// LogoutUser
// @Description logout user
// @Tags user
// @Produce  json
// @Param user body models.LogoutUserRequest true "User info"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/logout [post]
func (h *Handler) LogoutUser(ctx *gin.Context) {
	raw, _ := ctx.GetRawData()
	fmt.Println("RAW REQUEST BODY:", string(raw)) // Логируем сырое тело запроса

	var request models.LogoutUserRequest

	// Пробуем сразу распарсить JSON
	if err := json.Unmarshal(raw, &request); err != nil {
		fmt.Println("Ошибка парсинга JSON:", err)

		// КОСТЫЛЬ: если пришла строка (например, "lexa"), пытаемся обернуть её в JSON
		if json.Valid([]byte(fmt.Sprintf(`{"login": %s}`, raw))) {
			fixedJSON := []byte(fmt.Sprintf(`{"login": %s}`, raw))
			fmt.Println("Исправленный JSON:", string(fixedJSON))
			_ = json.Unmarshal(fixedJSON, &request) // Повторная попытка парсинга
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный JSON"})
			return
		}
	}

	// Проверяем, не пустой ли логин после обработки
	if request.Login == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Поле login отсутствует"})
		return
	}

	// Извлекаем токен из заголовка
	tokenString := extractTokenFromHeader(ctx.Request)
	if tokenString == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Вызываем логику выхода пользователя
	err := h.Repository.LogoutUser(request.Login)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"logout": "success"})
}
