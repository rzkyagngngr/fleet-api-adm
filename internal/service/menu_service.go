package service

import (
	"gin-boilerplate/internal/model/entity"
	"gin-boilerplate/internal/repository"
)

type MenuService interface {
	Create(menu *entity.Menu) error
	FindAll() ([]entity.Menu, error)
	FindByID(id uint64) (*entity.Menu, error)
	Update(id uint64, menu *entity.Menu) error
	Delete(id uint64) error
}

type menuService struct {
	menuRepo      repository.MenuRepository
	accessService AccessService
}

func NewMenuService(menuRepo repository.MenuRepository, accessService AccessService) MenuService {
	return &menuService{
		menuRepo:      menuRepo,
		accessService: accessService,
	}
}

func (s *menuService) Create(menu *entity.Menu) error {
	err := s.menuRepo.Create(menu)
	if err != nil {
		return err
	}
	return s.accessService.InitializeMenuAccessForAllRoles(menu)
}

func (s *menuService) FindAll() ([]entity.Menu, error) {
	return s.menuRepo.FindAll()
}

func (s *menuService) FindByID(id uint64) (*entity.Menu, error) {
	return s.menuRepo.FindByID(id)
}

func (s *menuService) Update(id uint64, menu *entity.Menu) error {
	return s.menuRepo.Update(id, menu)
}

func (s *menuService) Delete(id uint64) error {
	return s.menuRepo.Delete(id)
}
