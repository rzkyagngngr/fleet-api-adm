package service

import (
	"gin-boilerplate/internal/model/entity"
	"gin-boilerplate/internal/repository"
)

type RoleService interface {
	Create(role *entity.Role) error
	FindAll() ([]entity.Role, error)
	FindByID(id uint64) (*entity.Role, error)
	Update(id uint64, role *entity.Role) error
	Delete(id uint64) error
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

func (s *roleService) Create(role *entity.Role) error {
	err := s.roleRepo.Create(role)
	if err != nil {
		return err
	}

	// Auto-initialize access rules for all menus
	return s.accessService.InitializeRoleAccess(role.HakAksesID)
}

func (s *roleService) FindAll() ([]entity.Role, error) {
	return s.roleRepo.FindAll()
}

func (s *roleService) FindByID(id uint64) (*entity.Role, error) {
	return s.roleRepo.FindByID(id)
}

func (s *roleService) Update(id uint64, role *entity.Role) error {
	return s.roleRepo.Update(id, role)
}

func (s *roleService) Delete(id uint64) error {
	return s.roleRepo.Delete(id)
}
