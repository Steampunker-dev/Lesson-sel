package models

import (
	"awesomeProject/internal/app/ds"
	"mime/multipart"
)

// Запрос

type GetAllTaskRequest struct {
	MinutesFrom string `json:"minutes_from"`
	MinutesTo   string `json:"minutes_to"`
}

type GetTaskRequest struct {
	ID string `json:"id"`
}

type CreateTaskRequest struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Minutes     int    `json:"minutes"`
	Description string `json:"description"`
}

type UploadImageRequest struct {
	ID    uint                  `json:"id"`
	Image *multipart.FileHeader `form:"image"` // Поле для файла
}

// Ответ

type GetAllTaskResponse struct {
	ReqID          int            `json:"req_id"`
	ReqLessonCount int            `json:"req_less_count"`
	Card           *[]ds.TaskItem `json:"cards"`
}

type GetTaskResponse struct {
	Card *ds.TaskItem `json:"cards"`
}

type CreateTaskResponse struct {
	Lesson *ds.TaskItem `json:"lesson"`
}

type UploadImageResponse struct {
	ImageURL string `json:"image_url"`
}

type AddTasktoLessonResponse struct {
	TaskItem *ds.TaskItem `json:"task"`
}
