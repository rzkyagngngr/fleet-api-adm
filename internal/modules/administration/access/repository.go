package access

import (
	"context"

	"gorm.io/gorm"
)

type AccessRepository interface {
	FindByRoleID(ctx context.Context, roleID uint64) ([]Access, error)
	FindAllMenuByRole(ctx context.Context, roleID uint64) ([]Access, error)
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
func (r *accessRepository) FindAllMenuByRole(ctx context.Context, roleID uint64) ([]Access, error) {
	var list []Access
	err := r.db.WithContext(ctx).Table("adm.posm_menus m").
		Select(`
			CAST(a.access_id AS bigint) as access_id,
			CAST(? AS bigint) as roles_id,
			m.id as menu_id,
			m.menu_text as menu_text,
			m.menu_url as menu_url,
			COALESCE(a.status, 0) as status,
			m.application_id as application_id,
			m.parent_menu_id as parent_menu_id,
			COALESCE(a.can_insert, 0) as can_insert,
			COALESCE(a.can_update, 0) as can_update,
			COALESCE(a.can_delete, 0) as can_delete,
			m.menu_order as menu_order,
			m.menu_icon as menu_icon
		`, roleID).
		Joins("LEFT JOIN adm.posm_access a ON m.id = a.menu_id AND a.roles_id = ?", roleID).
		Where("m.menu_status = ?", 1).
		Order("COALESCE(m.parent_menu_id, 0), m.menu_order").
		Find(&list).Error
	return list, err
}
func (r *accessRepository) DeleteByRoleID(ctx context.Context, roleID uint64) error {
	return r.db.WithContext(ctx).Where("roles_id = ?", roleID).Delete(&Access{}).Error
}
func (r *accessRepository) BulkCreate(ctx context.Context, accessList []Access) error {
	return r.db.WithContext(ctx).Create(&accessList).Error
}
