package user

import "context"

type UserService interface {
	GetProfile(ctx context.Context, userID uint64) (*UserResponse, error)
}

type userService struct{ userRepo UserRepository }

func NewUserService(userRepo UserRepository) UserService { return &userService{userRepo: userRepo} }

func (s *userService) GetProfile(ctx context.Context, userID uint64) (*UserResponse, error) {
	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	res := ToResponse(u)
	return &res, nil
}
