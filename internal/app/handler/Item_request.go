package handler

import (
	"awesomeProject/internal/app/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"strconv"
)

// DeleteDC
// @Description delete dc
// @Tags DC
// @Produce  json
// @Param id path int true "lesson id"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tl/delete/{id} [delete]
func (h *Handler) DeleteDC(ctx *gin.Context) {
	id := ctx.Param("id")
	num, err1 := strconv.ParseUint(id, 10, 64) // 10 - система счисления, 64 - битность
	if err1 != nil {
		fmt.Println("Ошибка:", err1)
		return
	}
	fmt.Println(reflect.TypeOf(num), num)
	task_id := uint(num)
	err := h.Repository.DeleteTL(task_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "Deleted"})
}

// UpdateDCCount
// @Description update dc count
// @Tags DC
// @Produce  json
// @Param id path int true "lesson id"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tl/count/{id} [put]
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
