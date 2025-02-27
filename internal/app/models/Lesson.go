package models

import "awesomeProject/internal/app/ds"

// Запросы

type GetLessonsRequest struct {
	DateFrom string `form:"date_from"` // дата начала диапазона
	DateTo   string `form:"date_to"`   // дата конца диапазона
	Status   string `form:"status"`    // статус
	UserID   uint   `form:"user_id"`   // идентификатор пользователя
}

type GetMyLessonCardsRequest struct {
	UserId string `form:"user_id"` // идентификатор пользователя
}

type UpdateLessonRequest struct {
	ID         uint   `json:"id"`          // идентификатор
	Address    string `json:"address"`     // адрес
	LessonDate string `json:"lesson_date"` // дата доставки
	LessonType string `json:"lesson_type"` // тип доставки
}

type FinishLessonRequest struct {
	ID uint `json:"id"` // идентификатор
}

type CompleteOrRejectLessonRequest struct {
	ID          uint `json:"id"`           // идентификатор
	ModeratorID uint `json:"moderator_id"` // модератор
	IsComplete  bool `json:"is_complete"`  // завершен
}

// Ответы

type GetLessonsResponse struct {
	Lessons []*ds.LessonRequest `json:"lessons"` // список звонков
}

type GetMyLessonCardsResponse struct {
	LessonRequest *ds.LessonRequest `json:"lesson_request"` // заявка на доставку
	TaskItems     []*ds.TaskItem    `json:"task_items"`     // карточки доставки
	Count         int               `json:"count"`          // количество карточек
}

type GetLessonResponse struct {
	LessonRequest *ds.LessonRequest   `json:"lesson_request"` // заявка на доставку
	TaskItems     []TaskItemWithCount `json:"fines"`          // карточки доставки
	// общее число доставок
	Count int `json:"count"`
}

type TaskItemWithCount struct {
	ds.TaskItem
	Forced bool `json:"forced"`
}

type UpdateLessonResponse struct {
	CallRequest *ds.LessonRequest `json:"lesson_request"` // заявка на доставку
}

type FinishLessonResponse struct {
	CallRequest *ds.LessonRequest `json:"lesson_request"` // заявка на доставку
}

type CompleteOrRejectLessonResponse struct {
	CallRequest *ds.LessonRequest `json:"lessson_request"` // заявка на доставку
	TotalCount  int               `json:"total_count"`     // общее количество
}
