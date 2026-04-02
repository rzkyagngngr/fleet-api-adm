package role

import (
	"context"

	"gorm.io/gorm"
)

type RoleRepository interface {
	Create(ctx context.Context, role *Role) error
	FindAll(ctx context.Context) ([]Role, error)
	FindByID(ctx context.Context, id uint64) (*Role, error)
	Update(ctx context.Context, id uint64, role *Role) error
	Delete(ctx context.Context, id uint64) error
}

type roleRepository struct{ db *gorm.DB }

func NewRoleRepository(db *gorm.DB) RoleRepository { return &roleRepository{db: db} }
func (r *roleRepository) Create(ctx context.Context, role *Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}
func (r *roleRepository) FindAll(ctx context.Context) ([]Role, error) {
	var roles []Role
	err := r.db.WithContext(ctx).Find(&roles).Error
	return roles, err
}
func (r *roleRepository) FindByID(ctx context.Context, id uint64) (*Role, error) {
	var role Role
	if err := r.db.WithContext(ctx).Where("hak_akses_id = ?", id).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}
func (r *roleRepository) Update(ctx context.Context, id uint64, role *Role) error {
	return r.db.WithContext(ctx).Where("hak_akses_id = ?", id).Updates(role).Error
}
func (r *roleRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Where("hak_akses_id = ?", id).Delete(&Role{}).Error
}
