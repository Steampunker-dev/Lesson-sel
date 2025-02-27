package handler

import (
	"awesomeProject/internal/app/ds"
	"awesomeProject/internal/app/models"
	"awesomeProject/internal/app/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// GetLessons
// @Description get all lessons
// @Tags lesson
// @Produce json
// @Success 200 {object} models.GetLessonsResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /lesson [get]
func (h *Handler) GetLessons(ctx *gin.Context) {
	fmt.Println(ctx.GetRawData())
	var request models.GetLessonsRequest
	userID, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "userID not found"})
		return
	}

	request.UserID = userID.(uint)
	dateFromQuery := ctx.Query("date_from")
	dateToQuery := ctx.Query("date_to")
	statusQuery := ctx.Query("status")
	fmt.Println("dfsfdsf")
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
	lessons, err := h.Repository.GetLessons(dateFrom, dateTo, request.Status, request.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	// 3. Генерируем QR для каждого урока, собираем их в отдельный массив qrs
	qrs := make([]string, len(lessons))
	for i, lesson := range lessons {
		qrBase64, err := services.GenerateLessonQR(*lesson)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации QR"})
			return
		}
		qrs[i] = qrBase64
	}

	// 4. Возвращаем уроки и QR-коды (раздельно):
	ctx.JSON(http.StatusOK, gin.H{
		"lessons": lessons,
		"qrs":     qrs,
	})
}

// DeleteLesson
// @Description delete lesson
// @Tags call
// @Produce json
// @Param id path string true "Lesson ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /lesson/{id} [delete]
func (h *Handler) DeleteLesson(ctx *gin.Context) {
	id := ctx.Param("id")

	err := h.Repository.DeleteLesson(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

// GetMyLessonCards
// @Description get my lesson cards
// @Tags lesson
// @Produce json
// @Param id path string true "Lesson ID"
// @Success 200 {object} models.GetMyLessonCardsResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /lesson/{id} [get]
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

// GetLesson
// @Description get lesson by id
// @Tags lesson
// @Produce json
// @Param id path string true "Lesson ID"
// @Success 200 {object} models.GetLessonResponse
// @Failure 403 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /lesson/{id} [get]
func (h *Handler) GetLesson(ctx *gin.Context) {
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
		LessonRequest: lesson,
		TaskItems:     taskItemsWithCount,
		Count:         len(taskItemsWithCount),
	})
}

// UpdateLesson
// @Description update lesson
// @Tags lesson
// @Produce json
// @Param id path string true "Lesson ID"
// @Success 200 {object} models.UpdateLessonResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /Lesson/{id} [put]
func (h *Handler) UpdateLesson(ctx *gin.Context) {
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

// FormLesson
// @Description form lesson
// @Tags lesson
// @Produce json
// @Param id path string true "Lesson ID"
// @Success 200 {object} models.UpdateLessonResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /lesson/form/{id} [put]
func (h *Handler) FormLesson(ctx *gin.Context) {
	var request models.FinishLessonRequest
	id, _ := strconv.Atoi(ctx.Param("id"))
	request.ID = uint(id)
	fmt.Println(request.ID)

	fmt.Println("request")

	resp, err := h.Repository.FormLesson(request.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, models.UpdateLessonResponse{
		CallRequest: resp,
	})
}

// CompleteOrRejectLesson
// @Description complete or reject lesson
// @Tags lesson
// @Produce json
// @Param id path string true "Lesson ID"
// @Success 200 {object} models.CompleteOrRejectLessonResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /lesson/complete/{id} [put]
func (h *Handler) CompleteOrRejectLesson(ctx *gin.Context) {
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
	qrCode, err := services.GenerateLessonQR(*resp)
	println(qrCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка генерации QR-кода: " + err.Error(),
		})
		return
	}
	fmt.Println(qrCode)
	ctx.JSON(http.StatusOK, gin.H{
		"CallRequest": resp,
		"TotalCount":  totalCount,
		"qr":          qrCode,
	})
}
