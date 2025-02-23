package repository

import (
	"github.com/golang-jwt/jwt"
	"os"
	"strconv"
	"time"
)

const ExpTime = time.Hour * 1

func GenerateJWTToken(userID uint, isAdmin bool) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["is_admin"] = isAdmin
	claims["exp"] = ExpTime

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// CheckJWTInBlacklist - проверка токена в черном списке
func (r *Repository) CheckJWTInBlacklist(jwtStr string) (string, error) {
	res, err := r.rd.Get(jwtStr).Result()
	// если токен есть в черном списке, то возвращаем ошибку
	if err != nil {
		return "blacklist", err
	}
	return res, nil
}

// IsBlackListed - проверка на наличие токена в черном списке
func (r *Repository) IsBlackListed(userID float64, token string) bool {
	userIDStr := strconv.FormatUint(uint64(userID), 10)
	res, err := r.CheckJWTInBlacklist(userIDStr)
	if err != nil {
		return false
	}
	// если токен не равен токену в черном списке, то возвращаем false
	if res != token || res == "blacklist" {
		return false
	}
	return true
}

// SaveJWTToken - сохранение токена в БД
func (r *Repository) SaveJWTToken(userID uint, token string) error {
	exp := ExpTime
	idStr := strconv.FormatUint(uint64(userID), 10)

	err := r.rd.Set(idStr, token, exp).Err()
	if err != nil {
		return err
	}
	return nil
}
