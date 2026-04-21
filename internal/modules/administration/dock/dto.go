package dock

import "omniport-api/internal/helper"

type DockDetailReq struct {
	BerthCode  *string `json:"berth_code" binding:"omitempty,max=50"`
	BerthName  *string `json:"berth_name" binding:"omitempty,max=150"`
	MaxLoa     *int    `json:"max_loa"`
	XPosition  *int    `json:"x_position"`
	YPosition  *int    `json:"y_position"`
	WidthSize  *int    `json:"width_size"`
	HeightSize *int    `json:"height_size"`
	Status     *int    `json:"status"`
}

type DockReq struct {
	BranchCode      *int            `json:"branch_code"`
	BranchName      *string         `json:"branch_name" binding:"omitempty,max=100"`
	TerminalCode    *int            `json:"terminal_code"`
	TerminalName    *string         `json:"terminal_name" binding:"omitempty,max=100"`
	DockCode        *string         `json:"dock_code" binding:"omitempty,max=50"`
	DockName        *string         `json:"dock_name" binding:"omitempty,max=150"`
	DockType        *string         `json:"dock_type" binding:"omitempty,max=100"`
	DockLengthM     *float64        `json:"dock_length_m"`
	DockWidthM      *float64        `json:"dock_width_m"`
	DockCapacityTon *float64        `json:"dock_capacity_ton"`
	CodeInaportnet  *string         `json:"code_inaportnet" binding:"omitempty,max=50"`
	LocationNameIna *string         `json:"location_name_inaportnet" binding:"omitempty,max=150"`
	Status          *int            `json:"status"`
	Details         []DockDetailReq `json:"details"`
}

type SearchDockRequest struct {
	Page     int                  `json:"page"`
	Limit    int                  `json:"limit"`
	Search   string               `json:"search"`
	Filters  map[string]string    `json:"filters"`
	Sort     helper.SortQuery     `json:"sort"`
	Download helper.DownloadQuery `json:"download"`
}

func (r SearchDockRequest) ToPaginationQuery() helper.PaginationQuery {
	return helper.PaginationQuery{
		Page:     r.Page,
		Limit:    r.Limit,
		Search:   r.Search,
		Filters:  r.Filters,
		Sort:     r.Sort,
		Download: r.Download,
	}
}
