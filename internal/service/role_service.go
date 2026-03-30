package service

import (
	"context"
	"gin-boilerplate/internal/model/entity"
	"gin-boilerplate/internal/repository"
)

type RoleService interface {
	Create(ctx context.Context, role *entity.Role) error
	FindAll(ctx context.Context) ([]entity.Role, error)
	FindByID(ctx context.Context, id uint64) (*entity.Role, error)
	Update(ctx context.Context, id uint64, role *entity.Role) error
	Delete(ctx context.Context, id uint64) error
}

type roleService struct {
	roleRepo      repository.RoleRepository
	accessService AccessService
}

func NewRoleService(roleRepo repository.RoleRepository, accessService AccessService) RoleService {
	return &roleService{
		roleRepo:      roleRepo,
		accessService: accessService,
	}
}

func (s *roleService) Create(ctx context.Context, role *entity.Role) error {
	err := s.roleRepo.Create(ctx, role)
	if err != nil {
		return err
	}

	// Auto-initialize access rules for all menus
	return s.accessService.InitializeRoleAccess(ctx, role.HakAksesID)
}

func (s *roleService) FindAll(ctx context.Context) ([]entity.Role, error) {
	return s.roleRepo.FindAll(ctx)
}

func (s *roleService) FindByID(ctx context.Context, id uint64) (*entity.Role, error) {
	return s.roleRepo.FindByID(ctx, id)
}

func (s *roleService) Update(ctx context.Context, id uint64, role *entity.Role) error {
	return s.roleRepo.Update(ctx, id, role)
}

func (s *roleService) Delete(ctx context.Context, id uint64) error {
	return s.roleRepo.Delete(ctx, id)
}
