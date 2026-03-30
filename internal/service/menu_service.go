package service

import (
	"context"
	"gin-boilerplate/internal/model/entity"
	"gin-boilerplate/internal/repository"
)

type MenuService interface {
	Create(ctx context.Context, menu *entity.Menu) error
	FindAll(ctx context.Context) ([]entity.Menu, error)
	FindByID(ctx context.Context, id uint64) (*entity.Menu, error)
	Update(ctx context.Context, id uint64, menu *entity.Menu) error
	Delete(ctx context.Context, id uint64) error
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

func (s *menuService) Create(ctx context.Context, menu *entity.Menu) error {
	err := s.menuRepo.Create(ctx, menu)
	if err != nil {
		return err
	}
	return s.accessService.InitializeMenuAccessForAllRoles(ctx, menu)
}

func (s *menuService) FindAll(ctx context.Context) ([]entity.Menu, error) {
	return s.menuRepo.FindAll(ctx)
}

func (s *menuService) FindByID(ctx context.Context, id uint64) (*entity.Menu, error) {
	return s.menuRepo.FindByID(ctx, id)
}

func (s *menuService) Update(ctx context.Context, id uint64, menu *entity.Menu) error {
	return s.menuRepo.Update(ctx, id, menu)
}

func (s *menuService) Delete(ctx context.Context, id uint64) error {
	return s.menuRepo.Delete(ctx, id)
}
