package role

import "time"

type Role struct {
	HakAksesID      uint64     `gorm:"primaryKey;autoIncrement;column:hak_akses_id" json:"hak_akses_id"`
	HakAksesNama    *string    `gorm:"column:hak_akses_nama;size:50" json:"hak_akses_nama"`
	Status          *int16     `gorm:"column:status;default:0" json:"status"`
	CreationBy      *int64     `gorm:"column:creation_by" json:"creation_by"`
	CreationDate    *time.Time `gorm:"column:creation_date;default:CURRENT_TIMESTAMP" json:"creation_date"`
	LastUpdateBy    *int64     `gorm:"column:last_update_by" json:"last_update_by"`
	LastUpdateDate  *time.Time `gorm:"column:last_update_date" json:"last_update_date"`
	ApplicationID   *int       `gorm:"column:application_id" json:"application_id"`
	ApplicationNama *string    `gorm:"column:application_nama;size:50" json:"application_nama"`
	HakAksesSuper   *bool      `gorm:"column:hak_akses_super;default:false" json:"hak_akses_super"`
	DefaultMenuURL  *string    `gorm:"column:default_menu_url;size:200" json:"default_menu_url"`
}

func (Role) TableName() string { return "posm_roles" }
