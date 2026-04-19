package access

import "time"

type Access struct {
	AccessID       uint64     `gorm:"primaryKey;autoIncrement;column:access_id" json:"access_id"`
	RolesID        *int64     `gorm:"column:roles_id" json:"roles_id"`
	MenuID         *int64     `gorm:"column:menu_id" json:"menu_id"`
	MenuText       string     `gorm:"column:menu_text;size:250;not null" json:"menu_text"`
	MenuURL        *string    `gorm:"column:menu_url;size:250" json:"menu_url"`
	Status         *int16     `gorm:"column:status" json:"status"`
	ApplicationID  *int64     `gorm:"column:application_id" json:"application_id"`
	ParentMenuID   *int64     `gorm:"column:parent_menu_id" json:"parent_menu_id"`
	CanInsert      *int16     `gorm:"column:can_insert" json:"can_insert"`
	CanUpdate      *int16     `gorm:"column:can_update" json:"can_update"`
	CanDelete      *int16     `gorm:"column:can_delete" json:"can_delete"`
	MenuOrder      *int       `gorm:"column:menu_order" json:"menu_order"`
	MenuIcon       *string    `gorm:"column:menu_icon;size:100" json:"menu_icon"`
	CreationBy     *int64     `gorm:"column:creation_by" json:"creation_by"`
	CreationDate   *time.Time `gorm:"column:creation_date" json:"creation_date"`
	LastUpdateBy   *int64     `gorm:"column:last_update_by" json:"last_update_by"`
	LastUpdateDate *time.Time `gorm:"column:last_update_date" json:"last_update_date"`
}

func (Access) TableName() string { return "adm.posm_access" }
