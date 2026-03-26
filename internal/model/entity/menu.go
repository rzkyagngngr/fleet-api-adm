package entity

import (
	"time"
)

type Menu struct {
	ID              uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	MenuCode        string     `gorm:"column:menu_code;size:12;not null" json:"menu_code"`
	MenuText        string     `gorm:"column:menu_text;size:250;not null" json:"menu_text"`
	MenuDesc        *string    `gorm:"column:menu_desc;size:400" json:"menu_desc"`
	MenuUrl         *string    `gorm:"column:menu_url;size:250" json:"menu_url"`
	MenuLevel       int        `gorm:"column:menu_level;not null;default:1" json:"menu_level"`
	MenuOrder       *int       `gorm:"column:menu_order;default:0" json:"menu_order"`
	ParentMenuID    *int       `gorm:"column:parent_menu_id" json:"parent_menu_id"`
	MenuIcon        *string    `gorm:"column:menu_icon;size:30" json:"menu_icon"`
	ApplicationID   *int       `gorm:"column:application_id" json:"application_id"`
	MenuHeaderID    *int       `gorm:"column:menu_header_id" json:"menu_header_id"`
	MenuStatus      int16      `gorm:"column:menu_status;not null;default:1" json:"menu_status"`
	CreationBy      *string    `gorm:"column:creation_by" json:"creation_by"`
	CreationDate    *time.Time `gorm:"column:creation_date;default:CURRENT_TIMESTAMP" json:"creation_date"`
	LastUpdatedBy   *string    `gorm:"column:last_updated_by" json:"last_updated_by"`
	LastUpdatedDate *time.Time `gorm:"column:last_updated_date" json:"last_updated_date"`
}

func (Menu) TableName() string {
	return "posm_menus"
}
