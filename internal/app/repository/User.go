package repository

import (
	"awesomeProject/internal/app/ds"
)

// IsAdmin проверяет, является ли пользователь администратором
func (r *Repository) IsAdmin(userID uint) bool {
	var user ds.User
	if err := r.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return false
	}
	return user.IsAdmin == true
}

// CreateUser - создает нового пользователя
func (r *Repository) CreateUser(name, password string) (*ds.User, error) {
	user := &ds.User{
		Login:    name,
		Password: password,
	}
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser - обновляет данные пользователя
func (r *Repository) UpdateUser(id uint, name, password string) (*ds.User, error) {
	var user ds.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	// Изменение необходимых полей
	user.Login = name
	user.Password = password
	// Сохранение изменений
	if err := r.db.Model(&user).Updates(user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// AuthUser - аутентификация пользователя по логину и паролю
func (r *Repository) AuthUser(name, password string) (*ds.User, error) {
	var user ds.User
	if err := r.db.Where("login = ? AND password = ?", name, password).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// LogoutUser - выход пользователя
func (r *Repository) LogoutUser(userID uint) error {
	return nil
}
