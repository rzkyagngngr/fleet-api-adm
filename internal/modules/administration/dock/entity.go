package dock

import "time"

type Dock struct {
	ID              uint64       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	BranchCode      *int         `gorm:"column:branch_code" json:"branch_code"`
	BranchName      *string      `gorm:"column:branch_name;size:100" json:"branch_name"`
	TerminalCode    *int         `gorm:"column:terminal_code" json:"terminal_code"`
	TerminalName    *string      `gorm:"column:terminal_name;size:100" json:"terminal_name"`
	DockCode        *string      `gorm:"column:dock_code;size:50" json:"dock_code"`
	DockName        *string      `gorm:"column:dock_name;size:150" json:"dock_name"`
	DockType        *string      `gorm:"column:dock_type;size:100" json:"dock_type"`
	DockLengthM     *float64     `gorm:"column:dock_length_m" json:"dock_length_m"`
	DockWidthM      *float64     `gorm:"column:dock_width_m" json:"dock_width_m"`
	DockCapacityTon *float64     `gorm:"column:dock_capacity_ton" json:"dock_capacity_ton"`
	CodeInaportnet  *string      `gorm:"column:code_inaportnet;size:50" json:"code_inaportnet"`
	LocationNameIna *string      `gorm:"column:location_name_inaportnet;size:150" json:"location_name_inaportnet"`
	Status          *int         `gorm:"column:status" json:"status"`
	CreationDate    *time.Time   `gorm:"column:creation_date" json:"creation_date"`
	CreationBy      *string      `gorm:"column:creation_by;size:100" json:"creation_by"`
	LastUpdatedDate *time.Time   `gorm:"column:last_updated_date" json:"last_updated_date"`
	LastUpdatedBy   *string      `gorm:"column:last_updated_by;size:100" json:"last_updated_by"`
	Details         []DockDetail `gorm:"foreignKey:DockID;references:ID" json:"details,omitempty"`
}

func (Dock) TableName() string { return "adm.posm_docks" }

type DockDetail struct {
	ID              uint64     `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	DockID          uint64     `gorm:"column:dock_id" json:"dock_id"`
	BerthCode       *string    `gorm:"column:berth_code;size:50" json:"berth_code"`
	BerthName       *string    `gorm:"column:berth_name;size:150" json:"berth_name"`
	BerthLatitude   *string    `gorm:"column:berth_latitude;size:50" json:"berth_latitude"`
	BerthLongitude  *string    `gorm:"column:berth_longitude;size:50" json:"berth_longitude"`
	MaxLoa          *int       `gorm:"column:max_loa" json:"max_loa"`
	XPosition       *int       `gorm:"column:x_position" json:"x_position"`
	YPosition       *int       `gorm:"column:y_position" json:"y_position"`
	WidthSize       *int       `gorm:"column:width_size" json:"width_size"`
	HeightSize      *int       `gorm:"column:height_size" json:"height_size"`
	Status          *int       `gorm:"column:status" json:"status"`
	CreationDate    *time.Time `gorm:"column:creation_date" json:"creation_date"`
	CreationBy      *string    `gorm:"column:creation_by;size:100" json:"creation_by"`
	LastUpdatedDate *time.Time `gorm:"column:last_updated_date" json:"last_updated_date"`
	LastUpdatedBy   *string    `gorm:"column:last_updated_by;size:100" json:"last_updated_by"`
}

func (DockDetail) TableName() string { return "adm.posm_docks_d" }
