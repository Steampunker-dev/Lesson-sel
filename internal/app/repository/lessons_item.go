package repository

import (
	"awesomeProject/internal/app/ds"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// услуги

// LessonItemList возвращает список услуг
func (r *Repository) TaskItemList() (*[]ds.TaskItem, error) {
	var taskItems []ds.TaskItem
	r.db.Where("is_delete = ?", false).Find(&taskItems)
	return &taskItems, nil
}

// TaskLessonList получает список карточек и возвращает список forced по последним вхождениям
func (r *Repository) TaskLessonList(cards []*ds.TaskItem) ([]bool, error) {
	var taskLessons []ds.TaskLesson
	var forcedList []bool

	// Получаем все TaskLessons из базы данных
	if err := r.db.Find(&taskLessons).Error; err != nil {
		return nil, err
	}

	// Создаем мапу для быстрого доступа к Forced значению по ItemID
	forcedMap := make(map[int]bool)
	for _, lesson := range taskLessons {
		forcedMap[int(lesson.ItemID)] = lesson.Forced
	}

	// Заполняем forcedList на основе карточек
	for _, item := range cards {
		if forced, exists := forcedMap[int(item.ID)]; exists {
			forcedList = append(forcedList, forced)
		} else {
			forcedList = append(forcedList, false)
		}
	}

	return forcedList, nil
}

// SearchTaskItem возвращает список услуг, отфильтрованный по минутам
func (r *Repository) SearchTaskItem(minutesFrom, minutesTo string) (*[]ds.TaskItem, error) {
	intMinutesFrom, _ := strconv.Atoi(minutesFrom)
	intMinutesTo, _ := strconv.Atoi(minutesTo)

	var taskItems []ds.TaskItem
	// сохраняем данные из бд в массив
	r.db.Find(&taskItems)

	var filteredItems []ds.TaskItem
	for _, item := range taskItems {
		if item.Minutes <= intMinutesTo && item.Minutes >= intMinutesFrom {
			filteredItems = append(filteredItems, item)
			fmt.Println(filteredItems)
		}
	}
	return &filteredItems, nil
}

// DeleteTaskItem  удаляет услугу
func (r *Repository) DeleteTaskItem(id string) error {
	query := "UPDATE task_items SET is_delete = true WHERE id = $1"
	result := r.db.Exec(query, id)
	r.logger.Info("Rows affected:", result.RowsAffected)
	return nil
}

// DeleteLessonReq  удаляет заявку
func (r *Repository) DeleteLessonReq(id string) error {
	query := "UPDATE lesson_requests SET status = 'удален' WHERE id = $1"
	result := r.db.Exec(query, id)
	fmt.Println("ID del req   ", id, " stetus ")
	r.logger.Info("Rows affected:", result.RowsAffected)

	return nil
}

// GetTaskItemByID возвращает услугу по ID
func (r *Repository) GetTaskItemByID(id string) (*ds.TaskItem, error) {
	var TItem ds.TaskItem
	intID, _ := strconv.Atoi(id)
	r.db.Find(&TItem, intID)
	print(TItem.ID, "ID")
	return &TItem, nil
}

// HasRequestByUserID проверяет наличие заявки пользователя
func (r *Repository) HasRequestByUserID(userID uint) (uint, error) {
	var req ds.LessonRequest
	err := r.db.Where("user_id = ? AND status = ?", userID, ds.DraftStatus).First(&req).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	return req.ID, nil
}

// CreateOrUpdateLessonReq создает или обновляет заявку на занятие
func (r *Repository) CreateOrUpdateLessonReq(itemID, userID uint) (*ds.LessonRequest, error) {
	var order ds.LessonRequest
	err := r.db.Where("user_id = ? AND status = ?", userID, ds.DraftStatus).First(&order).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// создать
		order = ds.LessonRequest{
			UserID:      userID,
			Status:      ds.DraftStatus,
			DateCreated: time.Now(),
		}
		if err := r.db.Create(&order).Error; err != nil {
			return nil, err
		}
	}

	// добавить в заявку
	task_lessons := ds.TaskLesson{
		ItemID:    itemID,
		RequestID: order.ID,
		Forced:    false,
	}
	if err := r.db.Create(&task_lessons).Error; err != nil {
		return nil, err
	}

	return &order, nil
}

// GetLessonReqCount возвращает количество элементов в заявке
func (r *Repository) GetLessonReqCount(status string, userId uint) (int64, error) {
	var count int64
	var req ds.LessonRequest

	if err := r.db.Where("user_id = ? AND status = ?", userId, status).First(&req).Error; err != nil {
		return 0, err
	}

	reqID := req.ID

	err := r.db.Model(&ds.TaskLesson{}).Where("request_id = ?", reqID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	fmt.Println("Count ", count,
		"reqID ", reqID,
		"status ", status)

	return count, nil
}

// GetLessonItemsByUserAndStatus возвращает элементы заявки пользователя по статусу
func (r *Repository) GetLessonItemsByUserAndStatus(status string, userID uint) ([]*ds.TaskItem, error) {
	var items []*ds.TaskItem

	// Используем GORM для выполнения запроса
	err := r.db.Model(&ds.TaskItem{}).
		Select("task_items.*").
		Joins("INNER JOIN task_lessons ON task_items.id = task_lessons.item_id").
		Joins("INNER JOIN lesson_requests ON task_lessons.request_id = lesson_requests.id").
		Where("lesson_requests.user_id = ?", userID).
		Where("lesson_requests.status = ?", status).
		Find(&items).Error

	if err != nil {
		return nil, err
	}
	fmt.Println("Len items  ", len(items))

	return items, nil
}

// GettaskRequestById возвращает заявку по ID
func (r *Repository) GettaskRequestById(id uint) (*ds.LessonRequest, error) {
	var taskRequest ds.LessonRequest
	err := r.db.Where("id = ?", id).First(&taskRequest).Error

	if err != nil {
		return nil, fmt.Errorf("error fetching task request: %w", err)
	}

	// Выводим информацию о найденной записи
	r.logger.Infof("Found task request ID: %d", id)

	return &taskRequest, nil
}

// CreateDraftRequestAndGetID создает черновик заявки и возвращает ID
func (r *Repository) CreateDraftRequestAndGetID(userID uint) (uint, error) {
	draftRequest := ds.LessonRequest{
		Status:     ds.DraftStatus,
		UserID:     userID,
		LessonDate: time.Now(),
		LessonType: ds.Common_lesson,
	}

	err := r.db.Create(&draftRequest).Error
	if err != nil {
		return 0, fmt.Errorf("error creating draft request: %w", err)
	}

	r.logger.Infof("Created new draft request ID: %d", draftRequest.ID)

	return draftRequest.ID, nil
}

// LinkItemToDraftRequest связывает элемент с черновиком заявки
func (r *Repository) LinkItemToDraftRequest(userID uint, itemId uint) error {

	// поик существующей заявки пользователя со статусом 'черновик'
	var draftRequest ds.LessonRequest
	err := r.db.Where("user_id = ? AND status = ?", userID, ds.DraftStatus).First(&draftRequest).Error
	if err == gorm.ErrRecordNotFound {

		// если заявки нет, создаем новую
		draftRequest.UserID = userID
		draftRequest.ModeratorID = userID
		draftRequest.Status = ds.DraftStatus
		draftRequest.LessonDate = time.Now()
		draftRequest.LessonType = ds.Common_lesson
		err = r.db.Create(&draftRequest).Error

		if err != nil {
			return fmt.Errorf("error creating new draft request: %w", err)
		}
		r.logger.Infof("Created new draft request ID: %d for user ID: %d", draftRequest.ID, userID)
	} else {
		r.logger.Infof("Found existing draft request ID: %d for user ID: %d", draftRequest.ID, userID)
	}
	// Добавляем элемент в существующую заявку
	itemRequest := ds.TaskLesson{
		ItemID:    itemId,
		RequestID: draftRequest.ID,
		Forced:    false,
	}
	err = r.db.Create(&itemRequest).Error
	if err != nil {
		return fmt.Errorf("error linking item to draft request: %w", err)
	}

	return nil
}
