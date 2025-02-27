package repository

import (
	"awesomeProject/internal/app/ds"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
)

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
	// логин не может повторяться
	if err := r.db.Where("login = ?", name).First(&user).Error; err == nil {
		return nil, fmt.Errorf("user with login %s already exists", name)
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
func (r *Repository) AuthUser(name, password string) (string, bool, error) {
	var user ds.User
	// сперва проверим по логину
	if err := r.db.Where("login = ?", name).First(&user).Error; err != nil {
		return "", false, errors.New("user does not exist")

	}
	if user.Password != password {
		return "", false, errors.New("incorrect password")
	}

	token, err := GenerateJWTToken(user.ID, user.IsAdmin)
	if err != nil {
		return "", false, err
	}

	err = r.SaveJWTToken(user.ID, token)
	if err != nil {
		return "", false, err
	}

	return token, user.IsAdmin, nil
}

// LogoutUser - выход пользователя
func (r *Repository) LogoutUser(login string) error {
	var user ds.User
	if err := r.db.Where("login = ?", login).First(&user).Error; err != nil {
		return errors.New("user does not exist")
	}
	idStr := strconv.FormatUint(uint64(user.ID), 10)
	exp := ExpTime
	err := r.rd.Set(idStr, "blacklist", exp).Err()
	if err != nil {
		return err
	}
	return nil
}
