package ds

type TaskItem struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Image       string `json:"image" gorm:"type:varchar(255);default:''"`
	Title       string `json:"title" gorm:"type:varchar(255)"`
	Minutes     int    `json:"minutes" gorm:"not null"`
	Description string `json:"description" gorm:"type:text"`
	Answ        string `json:"answer" gorm:"type:varchar(255)"`
	IsDelete    bool   `json:"-" gorm:"default:false"`
}
