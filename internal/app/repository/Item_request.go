package repository

import (
	"awesomeProject/internal/app/ds"
	"fmt"
)

// DeleteTL - удаляет услугу из заявки
func (r *Repository) DeleteTL(taskID uint, callID uint) error {
	// Получение услуги из базы данных
	var existingTL ds.TaskLesson
	if err := r.db.Where("item_id = ? AND request_id = ?", taskID, callID).First(&existingTL).Error; err != nil {
		return err
	}

	// Получение заявки из базы данных
	var existingRequest ds.LessonRequest
	if err := r.db.First(&existingRequest, callID).Error; err != nil {
		return err
	}

	// Проверка, что статус заявки является "сформирован" или "черновик"
	if existingRequest.Status != ds.FormedStatus && existingRequest.Status != ds.DraftStatus {
		return fmt.Errorf("request status is not formed or draft")
	}

	// Удаление услуги из базы данных
	result := r.db.Delete(&existingTL)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateTLCount - обновляет важность услуг в заявке
func (r *Repository) UpdateTLCount(taskID uint, lessonID uint, forced bool) error {
	// Получение услуги из базы данных
	var existingTL ds.TaskLesson
	if err := r.db.Where("item_id = ? AND request_id = ?", taskID, lessonID).First(&existingTL).Error; err != nil {
		return err
	}

	// Получение заявки из базы данных
	var existingRequest ds.LessonRequest
	if err := r.db.First(&existingRequest, lessonID).Error; err != nil {
		return err
	}

	// Проверка, что статус заявки является "сформирован" или "черновик"
	if existingRequest.Status != ds.FormedStatus && existingRequest.Status != ds.DraftStatus {
		return fmt.Errorf("request status is not formed or draft")
	}

	// Обновление количества услуг в заявке
	// Обновление параметра forced в заявке
	result := r.db.Model(&existingTL).Update("Forced", forced)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
