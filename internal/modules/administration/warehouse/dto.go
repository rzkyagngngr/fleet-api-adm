package warehouse

import "omniport-api/internal/helper"

type WarehouseDetailReq struct {
	WarehouseCodeD        *string `json:"warehouse_code_d" binding:"omitempty,max=50"`
	WarehouseNameD        *string `json:"warehouse_name_d" binding:"omitempty,max=150"`
	WerehouseDType        *string `json:"werehouse_d_type" binding:"omitempty,max=50"`
	WarehouseCapacityDM3  *int    `json:"warehouse_capacity_d_m3"`
	WarehouseCapacityDTon *int    `json:"warehouse_capacity_d_ton"`
	XPosition             *int    `json:"x_position"`
	YPosition             *int    `json:"y_position"`
	WSize                 *int    `json:"w_size"`
	HSize                 *int    `json:"h_size"`
	Status                *int    `json:"status"`
}

type WarehouseReq struct {
	WarehouseCode     *string              `json:"warehouse_code" binding:"omitempty,max=50"`
	WarehouseName     *string              `json:"warehouse_name" binding:"omitempty,max=150"`
	WarehouseType     *string              `json:"warehouse_type" binding:"omitempty,max=100"`
	WarehouseCapacity *int                 `json:"warehouse_capacity"`
	Status            *int                 `json:"status"`
	Details           []WarehouseDetailReq `json:"details"`
}

type SearchWarehouseRequest struct {
	Page     int                  `json:"page"`
	Limit    int                  `json:"limit"`
	Search   string               `json:"search"`
	Filters  map[string]string    `json:"filters"`
	Sort     helper.SortQuery     `json:"sort"`
	Download helper.DownloadQuery `json:"download"`
}

func (r SearchWarehouseRequest) ToPaginationQuery() helper.PaginationQuery {
	return helper.PaginationQuery{
		Page:     r.Page,
		Limit:    r.Limit,
		Search:   r.Search,
		Filters:  r.Filters,
		Sort:     r.Sort,
		Download: r.Download,
	}
}
