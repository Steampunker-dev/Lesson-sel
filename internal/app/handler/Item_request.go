package handler

import (
	"awesomeProject/internal/app/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// DeleteDC - удаляет услугу из заявки
func (h *Handler) DeleteDC(ctx *gin.Context) {
	lessonid, _ := strconv.Atoi(ctx.Param("id"))
	var request models.DeleteTLRequest
	request.LessonID = uint(lessonid)
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Repository.DeleteTL(request.TaskID, request.LessonID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

// UpdateDCCount - обновляет количество услуг в заявке
func (h *Handler) UpdateDCCount(ctx *gin.Context) {
	lessonid, _ := strconv.Atoi(ctx.Param("id"))
	var request models.UpdateTLCountRequest
	request.LessonID = uint(lessonid)
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Repository.UpdateTLCount(request.TaskID, request.LessonID, request.Forced)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
