package warehouse

import "time"

type Warehouse struct {
	ID                uint64            `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	BranchCode        *int              `gorm:"column:branch_code" json:"branch_code"`
	BranchName        *string           `gorm:"column:branch_name;size:100" json:"branch_name"`
	TerminalCode      *int              `gorm:"column:terminal_code" json:"terminal_code"`
	TerminalName      *string           `gorm:"column:terminal_name;size:100" json:"terminal_name"`
	WarehouseCode     *string           `gorm:"column:warehouse_code;size:50" json:"warehouse_code"`
	WarehouseName     *string           `gorm:"column:warehouse_name;size:150" json:"warehouse_name"`
	WarehouseType     *string           `gorm:"column:warehouse_type;size:100" json:"warehouse_type"`
	WarehouseCapacity *string           `gorm:"column:warehouse_capacity;size:100" json:"warehouse_capacity"`
	Status            *int              `gorm:"column:status" json:"status"`
	CreationDate      *time.Time        `gorm:"column:creation_date" json:"creation_date"`
	CreationBy        *string           `gorm:"column:creation_by;size:100" json:"creation_by"`
	LastUpdatedDate   *time.Time        `gorm:"column:last_updated_date" json:"last_updated_date"`
	LastUpdatedBy     *string           `gorm:"column:last_updated_by;size:100" json:"last_updated_by"`
	Details           []WarehouseDetail `gorm:"foreignKey:WarehouseID;references:ID" json:"details,omitempty"`
}

func (Warehouse) TableName() string { return "adm.posm_warehouses" }

type WarehouseDetail struct {
	ID                    uint64     `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	WarehouseID           uint64     `gorm:"column:warehouse_id" json:"warehouse_id"`
	WarehouseCodeD        *string    `gorm:"column:warehouse_code_d;size:50" json:"warehouse_code_d"`
	WarehouseNameD        *string    `gorm:"column:warehouse_name_d;size:150" json:"warehouse_name_d"`
	WerehouseDType        *string    `gorm:"column:werehouse_d_type;size:50" json:"werehouse_d_type"`
	WarehouseCapacityDM3  *int       `gorm:"column:warehouse_capacity_d_m3" json:"warehouse_capacity_d_m3"`
	WarehouseCapacityDTon *int       `gorm:"column:warehouse_capacity_d_ton" json:"warehouse_capacity_d_ton"`
	XPosition             *int       `gorm:"column:x_position" json:"x_position"`
	YPosition             *int       `gorm:"column:y_position" json:"y_position"`
	WSize                 *int       `gorm:"column:w_size" json:"w_size"`
	HSize                 *int       `gorm:"column:h_size" json:"h_size"`
	Status                *int       `gorm:"column:status" json:"status"`
	CreationDate          *time.Time `gorm:"column:creation_date" json:"creation_date"`
	CreationBy            *string    `gorm:"column:creation_by;size:100" json:"creation_by"`
	LastUpdatedDate       *time.Time `gorm:"column:last_updated_date" json:"last_updated_date"`
	LastUpdatedBy         *string    `gorm:"column:last_updated_by;size:100" json:"last_updated_by"`
}

func (WarehouseDetail) TableName() string { return "adm.posm_warehouses_d" }
