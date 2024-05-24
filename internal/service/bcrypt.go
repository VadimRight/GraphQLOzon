package service

import "golang.org/x/crypto/bcrypt"

// Функция хеширования пароля
func (s *userService) HashPassword(password string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed)
}

// Функция сравнения вводимого пароля и пароля полученного из базы данных при аутентификации
func (s *userService) ComparePassword(hashed string, normal string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(normal))
}
