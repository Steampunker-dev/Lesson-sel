package ds

import (
	"encoding/json"
	"time"
)

type LessonRequest struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	DateCreated  time.Time `json:"date_created"`
	DateFormed   time.Time `json:"date_formed"`
	DateAccepted time.Time `json:"date_accepted"`
	Status       string    `json:"status" gorm:"type:varchar(255)"`

	LessonDate  time.Time `json:"lesson_date"`
	LessonType  string    `json:"delivery_type" gorm:"type:varchar(255)"`
	UserID      uint      `json:"-"`
	ModeratorID uint      `json:"-"`
	User        User      `json:"-" gorm:"foreignKey:UserID"`
	Moderator   User      `json:"-" gorm:"foreignKey:ModeratorID"`
	LessonTime  int64     `json:"Lesson_time"`
}

const (
	DraftStatus     = "черновик"
	DeletedStatus   = "удален"
	FormedStatus    = "сформирован"
	CompletedStatus = "завершен"
	RejectedStatus  = "отклонен"
)

const (
	Test_lesson   = "Контрольная работа"
	Common_lesson = "Обычное занятие"
	Exam_lesson   = "Экзамен"
)

func (d LessonRequest) MarshalJSON() ([]byte, error) {
	type Alias LessonRequest
	return json.Marshal(&struct {
		DeliveryDate string `json:"delivery_date"`
		*Alias
	}{
		DeliveryDate: d.LessonDate.Format("2006-01-02"),
		Alias:        (*Alias)(&d),
	})
}
