package repository

import (
	"gin-boilerplate/internal/model/entity"
	"gorm.io/gorm"
)

type MenuRepository interface {
	Create(menu *entity.Menu) error
	FindAll() ([]entity.Menu, error)
	FindByID(id uint64) (*entity.Menu, error)
	Update(id uint64, menu *entity.Menu) error
	Delete(id uint64) error
}

type menuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{db: db}
}

func (r *menuRepository) Create(menu *entity.Menu) error {
	return r.db.Create(menu).Error
}

func (r *menuRepository) FindAll() ([]entity.Menu, error) {
	var menus []entity.Menu
	err := r.db.Order("menu_level asc, id asc").Find(&menus).Error
	return menus, err
}

func (r *menuRepository) FindByID(id uint64) (*entity.Menu, error) {
	var menu entity.Menu
	err := r.db.Where("id = ?", id).First(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *menuRepository) Update(id uint64, menu *entity.Menu) error {
	return r.db.Where("id = ?", id).Updates(menu).Error
}

func (r *menuRepository) Delete(id uint64) error {
	return r.db.Where("id = ?", id).Delete(&entity.Menu{}).Error
}
