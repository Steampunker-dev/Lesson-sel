package models

type DeleteTLRequest struct {
	TaskID uint `json:"id"`
}

type UpdateTLCountRequest struct {
	TaskID   uint `json:"task_id"`
	LessonID uint `json:"lesson_id"`
	Forced   bool `json:"forced"`
}
