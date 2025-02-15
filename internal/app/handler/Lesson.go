package handler

import (
	"awesomeProject/internal/app/ds"
	"awesomeProject/internal/app/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// GetLessons - возращает все заявки
func (h *Handler) GetLessons(ctx *gin.Context) {
	var request models.GetLessonsRequest
	dateFromQuery := ctx.Query("date_from")
	dateToQuery := ctx.Query("date_to")
	statusQuery := ctx.Query("status")

	request.DateFrom = dateFromQuery
	request.DateTo = dateToQuery
	request.Status = statusQuery
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	layout := "2006-01-02"
	dateFrom, err := time.Parse(layout, request.DateFrom)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dateTo, err := time.Parse(layout, request.DateTo)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(dateTo, dateFrom)
	lessons, err := h.Repository.GetLessons(dateFrom, dateTo, request.Status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.GetLessonsResponse{Lessons: lessons})
}

// DeleteLesson - устанавливает статус "удалено" для звонка
func (h *Handler) DeleteLesson(ctx *gin.Context) {
	id := ctx.Param("id")

	err := h.Repository.DeleteLesson(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

// GetMyLessonCards рисует страницу с заявкой
func (h *Handler) GetMyLessonCards(ctx *gin.Context) {
	if lessonRequestId, err := strconv.Atoi(ctx.Param("id")); err == nil {
		// Предполагаем, что пользователь идентификатор равен 1
		userId := 1

		// Получаем заявку по ID
		lessonRequest, err := h.Repository.GetLessonRequestById(uint(lessonRequestId))
		if err != nil || lessonRequest.Status == ds.DeletedStatus {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Call request not found or deleted"})
			return
		}

		// Получаем карточки доставки для этой заявки
		tasks, err := h.Repository.GetLessonItemsByUserAndStatus(ds.DraftStatus, uint(userId))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// получаем колическво карточек из м-м таблицы item_request
		count, err := h.Repository.GetLessonReqCount(ds.DraftStatus, uint(userId))

		ctx.JSON(http.StatusOK, models.GetMyLessonCardsResponse{
			LessonRequest: lessonRequest,
			TaskItems:     tasks,
			Count:         int(count),
		})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// GetCall возвращает заявку на звонок-заявку
func (h *Handler) GetCall(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	lesson, err := h.Repository.GetLessonRequestById(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Получаем все доставки для этой заявки на звонок
	itemRequests, err := h.Repository.GetItemRequestsByLessonRequestID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Формируем список доставок с количеством
	var taskItemsWithCount []models.TaskItemWithCount
	for _, itemRequest := range itemRequests {
		taskItemsWithCount = append(taskItemsWithCount, models.TaskItemWithCount{
			TaskItem: itemRequest.Item,
			Forced:   itemRequest.Forced,
		})
	}

	ctx.JSON(http.StatusOK, models.GetLessonResponse{
		LessonRequest:   lesson,
		TaskItems:       taskItemsWithCount,
		DeliveriesCount: len(taskItemsWithCount),
	})
}

// UpdateCall обновляет заявку  по теме
func (h *Handler) UpdateCall(ctx *gin.Context) {
	var request models.UpdateLessonRequest
	id, _ := strconv.Atoi(ctx.Param("id"))
	request.ID = uint(id)
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("даты:", request.LessonDate)
	layout := "2006-01-02" // Измените формат на этот, если вы передаете дату без времени
	lessonDate, err := time.Parse(layout, request.LessonDate)
	if err != nil {
		fmt.Println("Ошибка при парсинге даты:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}
	lesson := &ds.LessonRequest{
		ID: request.ID,
		//Address:    request.Address,
		LessonDate: lessonDate,
		LessonType: request.LessonType,
	}
	resp, err := h.Repository.UpdateLessons(lesson)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, models.UpdateLessonResponse{
		CallRequest: resp,
	})
}

// FormCall - формирует заявку
func (h *Handler) FormCall(ctx *gin.Context) {
	var request models.FinishLessonRequest
	id, _ := strconv.Atoi(ctx.Param("id"))
	request.ID = uint(id)

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(request.UserID, "")
	resp, err := h.Repository.FormLesson(request.ID, request.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, models.UpdateLessonResponse{
		CallRequest: resp,
	})
}

// CompleteOrRejectCall - завершает заявку
func (h *Handler) CompleteOrRejectCall(ctx *gin.Context) {
	var request models.CompleteOrRejectLessonRequest
	id, _ := strconv.Atoi(ctx.Param("id"))
	request.ID = uint(id)
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Получаем полный объект call из базы данных
	call, err := h.Repository.GetLessonRequestById(request.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Устанавливаем ModeratorID из запроса
	call.ModeratorID = request.ModeratorID

	resp, totalCount, err := h.Repository.CompleteOrRejectLesson(call, request.IsComplete)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, models.CompleteOrRejectLessonResponse{
		CallRequest: resp,
		TotalCount:  totalCount,
	})
}
