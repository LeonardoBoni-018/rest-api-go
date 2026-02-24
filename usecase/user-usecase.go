package usecase

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"go-api/model"
	"go-api/repository"
)

type AuthUseCase struct {
	userRepo *repository.UserRepository
}

func NewAuthUserCase(userRepo *repository.UserRepository) *AuthUseCase {
	return &AuthUseCase{userRepo: userRepo}
}

func (a *AuthUseCase) Authenticate(username, password string) (*model.User, error) {
	user, hash, err := a.userRepo.GetuserByName(username)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil // usuário não encontrado
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return nil, errors.New("Invalid credentials")
	}
	return user, nil
}
