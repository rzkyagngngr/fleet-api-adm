package access

import (
	"context"

	"gorm.io/gorm"
)

type AccessRepository interface {
	FindByRoleID(ctx context.Context, roleID uint64) ([]Access, error)
	DeleteByRoleID(ctx context.Context, roleID uint64) error
	BulkCreate(ctx context.Context, accessList []Access) error
}

type accessRepository struct{ db *gorm.DB }

func NewAccessRepository(db *gorm.DB) AccessRepository { return &accessRepository{db: db} }
func (r *accessRepository) FindByRoleID(ctx context.Context, roleID uint64) ([]Access, error) {
	var list []Access
	err := r.db.WithContext(ctx).Where("roles_id = ?", roleID).Find(&list).Error
	return list, err
}
func (r *accessRepository) DeleteByRoleID(ctx context.Context, roleID uint64) error {
	return r.db.WithContext(ctx).Where("roles_id = ?", roleID).Delete(&Access{}).Error
}
func (r *accessRepository) BulkCreate(ctx context.Context, accessList []Access) error {
	return r.db.WithContext(ctx).Create(&accessList).Error
}
