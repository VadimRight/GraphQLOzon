package service

import "golang.org/x/crypto/bcrypt"

// PasswordService интерфейс для работы с паролями
type PasswordService interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashed string, normal string) bool
}

type passwordService struct{}

func NewPasswordService() PasswordService {
	return &passwordService{}
}

func (s *passwordService) HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (s *passwordService) ComparePassword(hashed string, normal string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(normal))
	return err == nil
}
