package repository

import (
	"context"
	"gin-boilerplate/internal/model/entity"
	"gorm.io/gorm"
)

type AccessRepository interface {
	Create(ctx context.Context, access *entity.Access) error
	FindByRoleID(ctx context.Context, roleID uint64) ([]entity.Access, error)
	Update(ctx context.Context, id uint64, access *entity.Access) error
	DeleteByRoleID(ctx context.Context, roleID uint64) error
	BulkCreate(ctx context.Context, accessList []entity.Access) error
	FindAllMenus(ctx context.Context) ([]entity.Menu, error)
	FindAllRoles(ctx context.Context) ([]entity.Role, error)
}

type accessRepository struct {
	db *gorm.DB
}

func NewAccessRepository(db *gorm.DB) AccessRepository {
	return &accessRepository{db: db}
}

func (r *accessRepository) Create(ctx context.Context, access *entity.Access) error {
	return r.db.WithContext(ctx).Create(access).Error
}

func (r *accessRepository) FindByRoleID(ctx context.Context, roleID uint64) ([]entity.Access, error) {
	var accessList []entity.Access
	err := r.db.WithContext(ctx).Where("roles_id = ?", roleID).Find(&accessList).Error
	return accessList, err
}

func (r *accessRepository) Update(ctx context.Context, id uint64, access *entity.Access) error {
	return r.db.WithContext(ctx).Where("access_id = ?", id).Updates(access).Error
}

func (r *accessRepository) DeleteByRoleID(ctx context.Context, roleID uint64) error {
	return r.db.WithContext(ctx).Where("roles_id = ?", roleID).Delete(&entity.Access{}).Error
}

func (r *accessRepository) BulkCreate(ctx context.Context, accessList []entity.Access) error {
	return r.db.WithContext(ctx).Create(&accessList).Error
}

func (r *accessRepository) FindAllMenus(ctx context.Context) ([]entity.Menu, error) {
	var menus []entity.Menu
	err := r.db.WithContext(ctx).Find(&menus).Error
	return menus, err
}

func (r *accessRepository) FindAllRoles(ctx context.Context) ([]entity.Role, error) {
	var roles []entity.Role
	err := r.db.WithContext(ctx).Find(&roles).Error
	return roles, err
}
