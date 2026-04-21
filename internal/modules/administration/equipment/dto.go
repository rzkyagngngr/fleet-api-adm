package equipment

import (
	"omniport-api/internal/helper"
	"time"
)

type EquipmentReq struct {
	BranchCode          *int       `json:"branch_code"`
	BranchName          *string    `json:"branch_name" binding:"omitempty,max=100"`
	TerminalCode        *int       `json:"terminal_code"`
	TerminalName        *string    `json:"terminal_name" binding:"omitempty,max=100"`
	EquipmentCode       *string    `json:"equipment_code" binding:"omitempty,max=50"`
	EquipmentName       *string    `json:"equipment_name" binding:"omitempty,max=150"`
	EquipmentGroup      *string    `json:"equipment_group" binding:"omitempty,max=100"`
	EquipmentType       *string    `json:"equipment_type" binding:"omitempty,max=100"`
	Capacity            *int       `json:"capacity"`
	MinimalLoadCapacity *int       `json:"minimal_load_capacity"`
	MaxLoadCapacity     *int       `json:"max_load_capacity"`
	OwnershipStatus     *string    `json:"ownership_status" binding:"omitempty,max=50"`
	OwnerName           *string    `json:"owner_name" binding:"omitempty,max=150"`
	OwnerCode           *string    `json:"owner_code" binding:"omitempty,max=50"`
	StartDate           *time.Time `json:"start_date"`
	EndDate             *time.Time `json:"end_date"`
	EquipmentCondition  *string    `json:"equipment_condition" binding:"omitempty,max=100"`
	Status              *int       `json:"status"`
}

type SearchEquipmentRequest struct {
	Page     int                  `json:"page"`
	Limit    int                  `json:"limit"`
	Search   string               `json:"search"`
	Filters  map[string]string    `json:"filters"`
	Sort     helper.SortQuery     `json:"sort"`
	Download helper.DownloadQuery `json:"download"`
}

type CustomerOptionRequest struct {
	Q     string `json:"q"`
	Limit int    `json:"limit"`
}

type CustomerOption struct {
	CustomerID   uint64 `json:"customer_id"`
	CustomerCode string `json:"customer_code"`
	CustomerName string `json:"customer_name"`
	OwnerCode    string `json:"owner_code"`
	OwnerName    string `json:"owner_name"`
	Value        string `json:"value"`
	Label        string `json:"label"`
}

type EquipmentGroupOptionRequest struct {
	Q     string `json:"q"`
	Limit int    `json:"limit"`
}

type EquipmentGroupOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

func (r SearchEquipmentRequest) ToPaginationQuery() helper.PaginationQuery {
	return helper.PaginationQuery{
		Page:     r.Page,
		Limit:    r.Limit,
		Search:   r.Search,
		Filters:  r.Filters,
		Sort:     r.Sort,
		Download: r.Download,
	}
}
