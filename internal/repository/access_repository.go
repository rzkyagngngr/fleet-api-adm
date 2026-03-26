package repository

import (
	"gin-boilerplate/internal/model/entity"
	"gorm.io/gorm"
)

type AccessRepository interface {
	Create(access *entity.Access) error
	FindByRoleID(roleID uint64) ([]entity.Access, error)
	Update(id uint64, access *entity.Access) error
	DeleteByRoleID(roleID uint64) error
	BulkCreate(accessList []entity.Access) error
	FindAllMenus() ([]entity.Menu, error)
	FindAllRoles() ([]entity.Role, error)
}

type accessRepository struct {
	db *gorm.DB
}

func NewAccessRepository(db *gorm.DB) AccessRepository {
	return &accessRepository{db: db}
}

func (r *accessRepository) Create(access *entity.Access) error {
	return r.db.Create(access).Error
}

func (r *accessRepository) FindByRoleID(roleID uint64) ([]entity.Access, error) {
	var accessList []entity.Access
	err := r.db.Where("roles_id = ?", roleID).Find(&accessList).Error
	return accessList, err
}

func (r *accessRepository) Update(id uint64, access *entity.Access) error {
	return r.db.Where("access_id = ?", id).Updates(access).Error
}

func (r *accessRepository) DeleteByRoleID(roleID uint64) error {
	return r.db.Where("roles_id = ?", roleID).Delete(&entity.Access{}).Error
}

func (r *accessRepository) BulkCreate(accessList []entity.Access) error {
	return r.db.Create(&accessList).Error
}

func (r *accessRepository) FindAllMenus() ([]entity.Menu, error) {
	var menus []entity.Menu
	err := r.db.Find(&menus).Error
	return menus, err
}

func (r *accessRepository) FindAllRoles() ([]entity.Role, error) {
	var roles []entity.Role
	err := r.db.Find(&roles).Error
	return roles, err
}
