package repository

import (
	"awesomeProject/internal/app/ds"
	"github.com/go-playground/validator/v10"
	"strconv"
)

// услуги

// TaskItemList возвращает список услуг
func (r *Repository) TaskItemList() (*[]ds.TaskItem, error) {
	var taskItems []ds.TaskItem
	r.db.Where("is_delete = ?", false).Order("title ASC").Find(&taskItems)
	return &taskItems, nil
}

// SearchTaskItem возвращает список услуг, отфильтрованный по цене
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

// GetTaskItemByID возвращает услугу по ID
func (r *Repository) GetTaskItemByID(id string) (*ds.TaskItem, error) {
	var TasItem ds.TaskItem
	intID, _ := strconv.Atoi(id)
	r.db.Find(&TasItem, intID)
	print(TasItem.ID, "ID")
	return &TasItem, nil
}

// ////////////////////////////////////////////////////////////////////////////////////////////
// CreateTaskItem создает услугу
func (r *Repository) CreateTaskItem(task *ds.TaskItem) (*ds.TaskItem, error) {
	validate := validator.New()
	err := validate.Struct(task)
	if err != nil {
		return nil, err
	}

	result := r.db.Create(task)
	if result.Error != nil {
		return nil, result.Error
	}

	return task, nil
}

// ////////////////////////////////////////////////////////////////////////////
// UploadImage загружает изображение в Minio
func (r *Repository) UploadImage(id string, img string) (string, error) {
	query := "UPDATE task_items SET image = $1 WHERE id = $2"
	result := r.db.Exec(query, img, id)
	r.logger.Info("Rows affected:", result.RowsAffected, img)

	// получить строку, которая в итоге у нас получилась
	var imageURL string
	r.db.Model(&ds.TaskItem{}).Where("id = ?", id).Select("image").Row().Scan(&imageURL)
	return imageURL, nil
}

// UpdateTaskItem обновляет услугу
func (r *Repository) UpdateTaskItem(task *ds.TaskItem) (*ds.TaskItem, error) {
	// Получаем текущий элемент доставки
	currentTaskItem := &ds.TaskItem{}
	result := r.db.First(currentTaskItem, task.ID)
	if result.Error != nil {
		return nil, result.Error
	}

	// Жизнь без костыля - не жизнь
	task.Image = currentTaskItem.Image

	validate := validator.New()
	err := validate.Struct(task)
	if err != nil {
		return nil, err
	}

	result = r.db.Save(task)
	if result.Error != nil {
		return nil, result.Error
	}

	return task, nil
}
