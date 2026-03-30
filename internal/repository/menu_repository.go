package repository

import (
	"context"
	"gin-boilerplate/internal/model/entity"
	"gorm.io/gorm"
)

type MenuRepository interface {
	Create(ctx context.Context, menu *entity.Menu) error
	FindAll(ctx context.Context) ([]entity.Menu, error)
	FindByID(ctx context.Context, id uint64) (*entity.Menu, error)
	Update(ctx context.Context, id uint64, menu *entity.Menu) error
	Delete(ctx context.Context, id uint64) error
}

type menuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{db: db}
}

func (r *menuRepository) Create(ctx context.Context, menu *entity.Menu) error {
	return r.db.WithContext(ctx).Create(menu).Error
}

func (r *menuRepository) FindAll(ctx context.Context) ([]entity.Menu, error) {
	var menus []entity.Menu
	err := r.db.WithContext(ctx).Order("menu_level asc, id asc").Find(&menus).Error
	return menus, err
}

func (r *menuRepository) FindByID(ctx context.Context, id uint64) (*entity.Menu, error) {
	var menu entity.Menu
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *menuRepository) Update(ctx context.Context, id uint64, menu *entity.Menu) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Updates(menu).Error
}

func (r *menuRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.Menu{}).Error
}
