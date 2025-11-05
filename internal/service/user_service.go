package service

import (
	"context"
	"errors"
	"view-list/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	uRepo domain.UserRepo
}

func NewUserService(uRepo domain.UserRepo) domain.UserService {
	return &userService{uRepo: uRepo}
}

func (s *userService) Register(ctx context.Context, user *domain.User) error {
	_, err := s.uRepo.GetByEmail(ctx, user.Email)
	if err == nil {
		return errors.New("User already exists")
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashed)
	return s.uRepo.Create(ctx, user)
}

func (s *userService) Login(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := s.uRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("User not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("Invalid password")
	}

	return user, nil
}
