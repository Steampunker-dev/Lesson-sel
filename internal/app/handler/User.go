package handler

import (
	"awesomeProject/internal/app/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateUser - создает нового пользователя
func (h *Handler) CreateUser(ctx *gin.Context) {
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

	ctx.JSON(http.StatusOK, models.CreateUserResponse{
		Registration: "success",
		User:         user,
	})
}

// UpdateUser - обновляет данные пользователя
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
		Registration: "success",
		User:         user,
	})
}

// AuthUser - аутентификация пользователя по логину и паролю
func (h *Handler) AuthUser(ctx *gin.Context) {
	var request models.CreateUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Repository.AuthUser(request.Login, request.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, models.AuthUserResponse{
		Auth: "success",
		User: user,
	})
}

// LogoutUser - аутентификация пользователя по логину и паролю
func (h *Handler) LogoutUser(ctx *gin.Context) {
	var request models.LogoutUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_ = h.Repository.LogoutUser(request.ID)

	ctx.JSON(http.StatusOK, gin.H{
		"logout": "success",
	})
}
