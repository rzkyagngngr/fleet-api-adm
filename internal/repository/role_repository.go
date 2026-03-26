package repository

import (
	"gin-boilerplate/internal/model/entity"
	"gorm.io/gorm"
)

type RoleRepository interface {
	Create(role *entity.Role) error
	FindAll() ([]entity.Role, error)
	FindByID(id uint64) (*entity.Role, error)
	Update(id uint64, role *entity.Role) error
	Delete(id uint64) error
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(role *entity.Role) error {
	return r.db.Create(role).Error
}

func (r *roleRepository) FindAll() ([]entity.Role, error) {
	var roles []entity.Role
	err := r.db.Find(&roles).Error
	return roles, err
}

func (r *roleRepository) FindByID(id uint64) (*entity.Role, error) {
	var role entity.Role
	err := r.db.Where("hak_akses_id = ?", id).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) Update(id uint64, role *entity.Role) error {
	return r.db.Where("hak_akses_id = ?", id).Updates(role).Error
}

func (r *roleRepository) Delete(id uint64) error {
	return r.db.Where("hak_akses_id = ?", id).Delete(&entity.Role{}).Error
}
