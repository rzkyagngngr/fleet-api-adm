package equipment

import "time"

type Equipment struct {
	ID                  uint64     `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	BranchCode          *int       `gorm:"column:branch_code" json:"branch_code"`
	BranchName          *string    `gorm:"column:branch_name;size:100" json:"branch_name"`
	TerminalCode        *int       `gorm:"column:terminal_code" json:"terminal_code"`
	TerminalName        *string    `gorm:"column:terminal_name;size:100" json:"terminal_name"`
	EquipmentCode       *string    `gorm:"column:equipment_code;size:50" json:"equipment_code"`
	EquipmentName       *string    `gorm:"column:equipment_name;size:150" json:"equipment_name"`
	EquipmentGroup      *string    `gorm:"column:equipment_group;size:100" json:"equipment_group"`
	EquipmentType       *string    `gorm:"column:equipment_type;size:100" json:"equipment_type"`
	Capacity            *int       `gorm:"column:capacity" json:"capacity"`
	MinimalLoadCapacity *int       `gorm:"column:minimal_load_capacity" json:"minimal_load_capacity"`
	MaxLoadCapacity     *int       `gorm:"column:max_load_capacity" json:"max_load_capacity"`
	OwnershipStatus     *string    `gorm:"column:ownership_status;size:50" json:"ownership_status"`
	OwnerName           *string    `gorm:"column:owner_name;size:150" json:"owner_name"`
	OwnerCode           *string    `gorm:"column:owner_code;size:50" json:"owner_code"`
	StartDate           *time.Time `gorm:"column:start_date" json:"start_date"`
	EndDate             *time.Time `gorm:"column:end_date" json:"end_date"`
	EquipmentCondition  *string    `gorm:"column:equipment_condition;size:100" json:"equipment_condition"`
	Status              *int       `gorm:"column:status" json:"status"`
	CreationDate        *time.Time `gorm:"column:creation_date" json:"creation_date"`
	CreationBy          *string    `gorm:"column:creation_by;size:100" json:"creation_by"`
	LastUpdatedDate     *time.Time `gorm:"column:last_updated_date" json:"last_updated_date"`
	LastUpdatedBy       *string    `gorm:"column:last_updated_by;size:100" json:"last_updated_by"`
}

func (Equipment) TableName() string { return "adm.posm_equipments" }
