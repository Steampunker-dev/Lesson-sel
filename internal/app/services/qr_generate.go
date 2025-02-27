package services

import (
	"awesomeProject/internal/app/ds"
	"encoding/base64"
	"fmt"
	"github.com/skip2/go-qrcode"
	"time"
)

// getString возвращает значение строки, если оно не пустое, иначе — "не указан".
func getString(s string) string {
	if s == "" {
		return "не указан"
	}
	return s
}

// getInt возвращает строковое представление числа, если оно не равно 0, иначе — "не указан".
func getInt(n uint) string {
	if n == 0 {
		return "не указан"
	}
	return fmt.Sprintf("%d", n)
}

// getTime возвращает отформатированное значение времени, если оно не нулевое, иначе — "не указан".
func getTime(t time.Time) string {
	if t.IsZero() {
		return "не указан"
	}
	// Формат: ГГГГ-ММ-ДД ЧЧ:ММ:СС
	return t.Format("2006-01-02 15:04:05")
}

// GenerateLessonQR формирует строку из полей структуры ds.LessonRequest,
// подставляя "не указан" для пустых значений, генерирует QR-код (base64).
func GenerateLessonQR(l ds.LessonRequest) (string, error) {
	info := fmt.Sprintf(
		"Занятие (ID): %s\n"+
			"Статус: %s\n"+
			"Дата создания: %s\n"+
			"Дата формирования: %s\n"+
			"Дата принятия: %s\n"+
			"Тип занятия: %s\n"+
			"Дата занятия: %s\n"+
			// Ниже – поля пользователя, если нужны
			"Пользователь (ID): %s\n"+
			"Модератор (ID): %s",
		getInt(l.ID),
		getString(l.Status),
		getTime(l.DateCreated),
		getTime(l.DateFormed),
		getTime(l.DateAccepted),
		getString(l.LessonType),
		getTime(l.LessonDate),
		getInt(l.UserID),
		getInt(l.ModeratorID),
	)

	// Генерация QR-кода размером 256x256 с уровнем коррекции ошибок Medium.
	png, err := qrcode.Encode(info, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}

	// Преобразование PNG-изображения в строку base64.
	base64Image := base64.StdEncoding.EncodeToString(png)
	return base64Image, nil
}
