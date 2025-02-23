package handler

import (
	"awesomeProject/internal/app/ds"
	"awesomeProject/internal/app/models"
	"awesomeProject/internal/app/storage"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// GetAllTasks
// @Description get all tasks
// @Tags task
// @Produce  json
// @Param price_from query string false "Minutes from"
// @Param price_to query string false "Minutes to"
// @Success 200 {object} models.GetAllTaskResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /task [get]
func (h *Handler) GetAllTasks(ctx *gin.Context) {
	var request models.GetAllTaskRequest
	minutesFrom := ctx.Query("minutesFrom")
	minutesTo := ctx.Query("minutesTo")
	request.MinutesFrom = minutesFrom
	request.MinutesTo = minutesTo

	userId := 1
	reqCount, _ := h.Repository.GetLessonReqCount(ds.DraftStatus, uint(userId))

	// Проверка на наличие заявки
	reqID, _ := h.Repository.HasRequestByUserID(uint(userId))
	// Если заявки нет, нужно вывести заявку с нулевыми полями, пустую

	var cards *[]ds.TaskItem
	var err error
	if request.MinutesFrom == "" && request.MinutesTo == "" {
		cards, err = h.Repository.TaskItemList()
	} else {
		fmt.Println(request.MinutesFrom)
		fmt.Println(request.MinutesTo)

		cards, err = h.Repository.SearchTaskItem(request.MinutesFrom, request.MinutesTo)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	response := models.GetAllTaskResponse{
		ReqID:          int(reqID),
		ReqLessonCount: int(reqCount),
		Card:           cards,
	}

	ctx.JSON(http.StatusOK, response)
}

// GetTask
// @Description get task by id
// @Tags task
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} models.GetTaskResponse
// @Failure 500 {object} map[string]string
// @Router /task/{id} [get]
func (h *Handler) GetTask(ctx *gin.Context) {
	var request models.GetTaskRequest
	request.ID = ctx.Param("id")
	card, err := h.Repository.GetTaskItemByID(request.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, models.GetTaskResponse{
		Card: card,
	})
}

// CreateTask
// @Description create task
// @Tags task
// @Produce json
// @Success 200 {object} models.CreateTaskResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /task [post]
func (h *Handler) CreateTask(ctx *gin.Context) {
	var request models.CreateTaskRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	card := &ds.TaskItem{
		Title:       request.Title,
		Minutes:     request.Minutes,
		Description: request.Description,
	}
	newCard, err_ := h.Repository.CreateTaskItem(card)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err_.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, models.CreateTaskResponse{
		Lesson: newCard,
	})
}

// UploadImage
// @Description load image to task
// @Tags task
// @Produce json
// @Success 200 {object} models.UploadImageResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /task/img/{id} [post]
func (h *Handler) UploadImage(ctx *gin.Context) {
	var request models.UploadImageRequest
	// считываем id из запроса
	id, _ := strconv.Atoi(ctx.Param("id"))
	request.ID = uint(id)

	// Привязать данные из запроса к структуре
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Проверка, что поле Image не является nil
	if request.Image == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No image in request"})
		return
	}

	// Инициализация Minio хранилища
	minioStorage, err := storage.NewMinioStorage(
		os.Getenv("MINIO_ENDPOINT_URL"),
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
		os.Getenv("MINIO_SECURE") == "true",
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Minio client"})
		return
	}
	fmt.Println(request.ID, "IDDDDDDDDDDDD")

	// Извлечение файла из запроса
	file, err := request.Image.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to open image"})
		return
	}
	defer file.Close()

	// Генерация имени файла
	fileExtension := filepath.Ext(request.Image.Filename)
	fileName := strconv.Itoa(int(request.ID)) + fileExtension

	// Загрузка файла в Minio
	err = minioStorage.LoadImg(os.Getenv("MINIO_BUCKET_NAME"), fileName, file, request.Image.Size)
	fmt.Println(err)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load image"})
		return
	}

	// Генерация URL изображения
	imageURL := "http://" + os.Getenv("MINIO_ENDPOINT_URL") + "/" + os.Getenv("MINIO_BUCKET_NAME") + "/" + fileName
	fmt.Println(imageURL)
	delivery, err := h.Repository.GetTaskItemByID(strconv.Itoa(int(request.ID)))
	if delivery == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Delivery not found"})
		return
	}
	strURL, _ := h.Repository.UploadImage(strconv.Itoa(int(request.ID)), imageURL)

	// Ответ с URL изображения
	ctx.JSON(http.StatusOK, gin.H{"image_url": strURL})
}

// UpdateTask
// @Description update task
// @Tags task
// @Produce json
// @Success 200 {object} models.CreateTaskResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /task/{id} [put]
func (h *Handler) UpdateTask(ctx *gin.Context) {
	var request models.CreateTaskRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	id, _ := strconv.Atoi(ctx.Param("id"))
	request.ID = uint(id)
	card := &ds.TaskItem{
		ID:          request.ID,
		Title:       request.Title,
		Minutes:     request.Minutes,
		Description: request.Description,
	}
	updatedCard, err_ := h.Repository.UpdateTaskItem(card)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err_.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, models.CreateTaskResponse{
		Lesson: updatedCard,
	})

}

// DeleteTask
// @Description Delete task
// @Tags task
// @Produce json
// @Success 200 {object} models.CreateTaskResponse
// @Failure 500 {object} map[string]string
// @Router /task/{id} [delete]
func (h *Handler) DeleteTask(ctx *gin.Context) {
	id := ctx.Param("id")
	err := h.Repository.DeleteTaskItem(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Deleted",
	})
}

// AddTaskToLesson
// @Description Add task to lesson
// @Tags task
// @Produce json
// @Success 200 {object} models.AddTasktoLessonResponse
// @Failure 500 {object} map[string]string
// @Router /task/add/{id} [post]
func (h *Handler) AddTaskToLesson(ctx *gin.Context) {
	itemID := ctx.Param("id")
	intItemID, _ := strconv.Atoi(itemID)
	userID := 4

	err := h.Repository.LinkItemToDraftRequest(uint(userID), uint(intItemID))
	if err != nil {
	}
	task, err_ := h.Repository.GetTaskItemByID(itemID)
	if err_ != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err_.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, models.AddTasktoLessonResponse{
		TaskItem: task,
	})
}
