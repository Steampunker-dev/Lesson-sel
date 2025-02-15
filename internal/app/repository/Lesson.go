package repository

import (
	"awesomeProject/internal/app/ds"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

// DeleteLesson  удаляет заявку
func (r *Repository) DeleteLesson(id string) error {
	query := "UPDATE lesson_requests SET status = 'удален' WHERE id = $1"
	result := r.db.Exec(query, id)
	fmt.Println("ID del req   ", id, " stetus ")
	r.logger.Info("Rows affected:", result.RowsAffected)

	return nil
}

// CreateOrUpdateLessonReq создает или обновляет заявку на доставку
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
	itemRequest := ds.TaskLesson{
		ItemID:    itemID,
		RequestID: order.ID,
		Forced:    false,
	}
	if err := r.db.Create(&itemRequest).Error; err != nil {
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

// GetLessonRequestById возвращает заявку по ID
func (r *Repository) GetLessonRequestById(id uint) (*ds.LessonRequest, error) {
	var callRequest ds.LessonRequest
	err := r.db.Where("id = ?", id).First(&callRequest).Error

	if err != nil {
		return nil, fmt.Errorf("error fetching call request: %w", err)
	}

	// Выводим информацию о найденной записи
	r.logger.Infof("Found call request ID: %d", id)

	return &callRequest, nil
}

// GetDeliveryItemsByLessonRequestID возвращает элементы заявки по ID заявки
func (r *Repository) GetDeliveryItemsByLessonRequestID(lessonRequestID uint) ([]*ds.TaskItem, error) {
	var items []*ds.TaskItem

	// Используем GORM для выполнения запроса
	err := r.db.Model(&ds.TaskItem{}).
		Select("task_items.*").
		Joins("INNER JOIN task_lessons ON task_items.id = task_lessons.item_id").
		Where("task_lessons.request_id = ?", lessonRequestID).
		Find(&items).Error

	if err != nil {
		return nil, err
	}

	return items, nil
}

// CreateDraftRequestAndGetID создает черновик заявки и возвращает ID
func (r *Repository) CreateDraftRequestAndGetID(userID uint) (uint, error) {
	draftRequest := ds.LessonRequest{
		Status: ds.DraftStatus,
		UserID: userID,
		//Address:     "",
		DateCreated: time.Now(),
		LessonType:  ds.Common_lesson,
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
	// нужно проверить, является ли удаленной
	var item ds.TaskItem
	err_ := r.db.Where("id = ?", itemId).First(&item).Error
	if err_ != nil {
		return fmt.Errorf("error fetching item: %w", err_)
	}
	if item.IsDelete == true {
		return fmt.Errorf("item with id %d is deleted", itemId)
	}

	// поик существующей заявки пользователя со статусом 'черновик'
	var draftRequest ds.LessonRequest
	err := r.db.Where("user_id = ? AND status = ?", userID, ds.DraftStatus).First(&draftRequest).Error
	if err == gorm.ErrRecordNotFound {
		// если заявки нет, создаем новую
		draftRequest.UserID = userID
		draftRequest.Status = ds.DraftStatus
		//draftRequest.Address = ""
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

// HasRequestByUserID проверяет наличие заявки пользователя
func (r *Repository) HasRequestByUserID(userID uint) (uint, error) {
	var req ds.LessonRequest
	err := r.db.Where("user_id = ? AND status = ?", userID, ds.DraftStatus).First(&req).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Если ошибка о том, то записи нет, то нет ошибки, тк нужно потом вывести заявку с null+0 полями
			return 0, nil
		}
		return 0, err
	}
	return req.ID, nil
}

// GetLessons - возвращает пары с учетом фильтров
func (r *Repository) GetLessons(dateFrom, dateTo time.Time, status string) ([]*ds.LessonRequest, error) {
	var lessons []*ds.LessonRequest
	fmt.Println(dateFrom, dateTo, status)
	query := "SELECT * FROM lesson_requests WHERE date_formed BETWEEN ? AND ? AND status = ?"
	result := r.db.Raw(query, dateFrom, dateTo, status).Scan(&lessons)
	if result.Error != nil {
		return nil, result.Error
	}
	return lessons, nil
}

// UpdateLessons обновляет Lesson
func (r *Repository) UpdateLessons(lesson *ds.LessonRequest) (*ds.LessonRequest, error) {
	result := r.db.Model(&ds.LessonRequest{}).Where("id = ?", lesson.ID).Updates(map[string]interface{}{
		//"Address":      call.Address,
		"LessonDate": lesson.LessonDate,
		"LessonType": lesson.LessonType,
	})
	if result.Error != nil {
		return nil, result.Error
	}

	// Загрузка обновленной записи из базы данных
	var updatedLesson ds.LessonRequest
	if err := r.db.First(&updatedLesson, lesson.ID).Error; err != nil {
		return nil, err
	}

	return &updatedLesson, nil
}

// FormLesson- формирование Lesson
func (r *Repository) FormLesson(lessonID uint, userID uint) (*ds.LessonRequest, error) {
	// Получение звонка из базы данных
	var existingLesson ds.LessonRequest
	if err := r.db.Where("id = ? AND user_id = ?", lessonID, userID).First(&existingLesson).Error; err != nil {
		return nil, err
	}
	fmt.Println("call", existingLesson.ID, existingLesson.UserID, existingLesson.Status)

	b := !(existingLesson.UserID == userID) || !(r.IsAdmin(userID))
	a := !(existingLesson.UserID == userID)
	c := r.IsAdmin(userID)
	fmt.Println(b, a, c, existingLesson.UserID, userID)
	// Проверка, что пользователь является владельцем заявки или модератором
	if !(existingLesson.UserID == userID) {
		if !(r.IsAdmin(userID)) {
			return nil, fmt.Errorf("user with id %d is not the owner of the call request or a moderator", userID)
		}
	}

	// Проверка, что статус звонка является "черновиком"
	if existingLesson.Status != ds.DraftStatus {
		return nil, fmt.Errorf("call status is not draft")
	}

	result := r.db.Model(&existingLesson).Updates(map[string]interface{}{
		"Status":     ds.FormedStatus,
		"DateFormed": time.Now(),
	})
	if result.Error != nil {
		return nil, result.Error
	}

	// Загрузка обновленной записи из базы данных
	var updatedLesson ds.LessonRequest
	if err := r.db.First(&updatedLesson, lessonID).Error; err != nil {
		return nil, err
	}

	return &updatedLesson, nil
}

// /////////////////////////////////////////////////////////////////////////////////////////
// CompleteOrRejectLesson - завершает или отклоняет заявку
func (r *Repository) CompleteOrRejectLesson(lesson *ds.LessonRequest, isComplete bool) (*ds.LessonRequest, int, error) {
	// Проверка, является ли пользователь администратором
	isAdmin := r.IsAdmin(lesson.ModeratorID)

	if !isAdmin {
		return nil, 0, fmt.Errorf("user with id %d is not an admin", lesson.ModeratorID)
	}

	// Получение lesson из базы данных
	var existingLesson ds.LessonRequest
	if err := r.db.First(&existingLesson, lesson.ID).Error; err != nil {
		return nil, 0, err
	}

	// Проверка, что статус звонка является "сформированным"
	if existingLesson.Status != ds.FormedStatus {
		return nil, 0, fmt.Errorf("call status is not formed")
	}

	// Обновление статуса звонка
	newStatus := ds.RejectedStatus
	if isComplete {
		newStatus = ds.CompletedStatus
	}

	result := r.db.Model(&existingLesson).Updates(map[string]interface{}{
		"ModeratorID":  lesson.ModeratorID,
		"Status":       newStatus,
		"DateAccepted": time.Now(),
	})
	if result.Error != nil {
		return nil, 0, result.Error
	}

	var forcedCount int64

	err := r.db.Model(&ds.TaskLesson{}).
		Where("request_id = ? AND forced = ?", lesson.ID, true).
		Count(&forcedCount).Error

	if err != nil {
		return nil, 0, err
	}

	// Загрузка обновлённой записи из базы данных
	var updatedCall ds.LessonRequest
	if err := r.db.First(&updatedCall, lesson.ID).Error; err != nil {
		return nil, 0, err
	}

	return &updatedCall, int(forcedCount), nil
}

// GetItemRequestsByLessonRequestID возвращает элементы заявки по ID заявки
func (r *Repository) GetItemRequestsByLessonRequestID(lessonRequestID uint) ([]ds.TaskLesson, error) {
	var itemRequests []ds.TaskLesson
	err := r.db.Where("request_id = ?", lessonRequestID).Preload("Item").Find(&itemRequests).Error
	if err != nil {
		return nil, err
	}
	return itemRequests, nil
}
