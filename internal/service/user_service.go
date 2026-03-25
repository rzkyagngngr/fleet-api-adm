package service

import (
	"gin-boilerplate/internal/model/dto"
	"gin-boilerplate/internal/repository"
)

type UserService interface {
	GetProfile(userID uint64) (*dto.UserResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetProfile(userID uint64) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	resp := user.ToResponse()
	return &resp, nil
}
