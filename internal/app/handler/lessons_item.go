package handler

import (
	"awesomeProject/internal/app/ds"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// вызываются функции из репы, которые идут в бд
// то есть как бы прослойка между эндпоинтами и данными, которые идут их бд

// LessonItemList рисует главную страницу
func (h *Handler) LessonItemList(ctx *gin.Context) {
	minutesFrom := ctx.Query("minutes_from")
	minutesTo := ctx.Query("minutes_to")
	userId := 1
	reqCount, _ := h.Repository.GetLessonReqCount(ds.DraftStatus, uint(userId))
	reqID, _ := h.Repository.HasRequestByUserID(uint(userId))
	if minutesFrom == "" && minutesTo == "" {
		cards, err := h.Repository.TaskItemList()

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"NoCards":          "",
			"tasks":            cards,
			"SearchFrom":       minutesFrom,
			"SearchUp":         minutesTo,
			"SpecialTaskCount": reqCount,
			"ReqID":            reqID,
		})
		return
	}
	cards, err := h.Repository.SearchTaskItem(minutesFrom, minutesTo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"NoCards":          "",
		"tasks":            cards,
		"SearchFrom":       minutesFrom,
		"SearchUp":         minutesTo,
		"SpecialTaskCount": reqCount,
		"ReqID":            reqID,
	})

}

// LessonItemByID рисует страницу с карточкой
func (h *Handler) LessonItemByID(ctx *gin.Context) {
	id := ctx.Param("id")
	card, err := h.Repository.GetTaskItemByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	print("card", card.ID, card.Title)
	ctx.HTML(http.StatusOK, "cardDetails.html", gin.H{
		"task": card,
	})
}

// DeleteLessonReq удаляет заявку и редиректит на главную
func (h *Handler) DeleteLessonReq(ctx *gin.Context) {
	id := ctx.Param("id")
	err := h.Repository.DeleteLessonReq(id)
	if err != nil {
	}
	fmt.Println("1ID del req   ", id, " stetus ")
	ctx.Redirect(http.StatusFound, "/")

}

// AddLessonItem добавляет карточку в заявку
func (h *Handler) AddLessonItem(ctx *gin.Context) {
	itemID := ctx.Param("id")
	intItemID, _ := strconv.Atoi(itemID)
	userID := 1

	err := h.Repository.LinkItemToDraftRequest(uint(userID), uint(intItemID))
	if err != nil {
	}
	fmt.Println("1ID del req   ", itemID, " status ")
	ctx.Redirect(http.StatusFound, "/")
}

// GetMytaskCards рисует страницу с заявкой
func (h *Handler) GetMytaskCards(ctx *gin.Context) {
	if taskRequestId, err := strconv.Atoi(ctx.Param("taskrequest_id")); err == nil {
		print("id req = ", taskRequestId)

		// Предполагаем, что пользователь идентификатор равен 1
		user_id := 1

		// Получаем заявку по ID
		taskRequest, err := h.Repository.GettaskRequestById(uint(taskRequestId))
		if err != nil || taskRequest.Status == ds.DeletedStatus {
			// Если заявка не найдена или удалена, перенаправляем на главную страницу
			ctx.Redirect(http.StatusFound, "/")
			return
		}

		// Получаем карточки задачек для этой заявки
		cards, err := h.Repository.GetLessonItemsByUserAndStatus(ds.DraftStatus, uint(user_id))
		if err != nil {
			// Если произошла ошибка, перенаправляем на главную страницу
			ctx.Redirect(http.StatusFound, "/")
			return
		}
		forcedList, err := h.Repository.TaskLessonList(cards)
		if err != nil {
			// Если произошла ошибка, перенаправляем на главную страницу
			ctx.Redirect(http.StatusFound, "/")
			return
		}
		timestamp := taskRequest.LessonDate
		formattedTime := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", timestamp.Year(), timestamp.Month(), timestamp.Day(), timestamp.Hour(), timestamp.Minute(), timestamp.Second())

		// получаем колическво карточек из м-м таблицы item_request
		//count, err := h.Repository.GetLessonReqCount(ds.DraftStatus, uint(user_id))

		ctx.HTML(http.StatusOK, "mycards.html", gin.H{
			"tasks":      cards,
			"Data":       formattedTime,
			"LessonType": taskRequest.LessonType,
			"ReqID":      taskRequestId,
			"ForcedList": forcedList,
		})
	} else {
		h.errorHandler(ctx, http.StatusBadRequest, err)
	}
}
