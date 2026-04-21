package branch

import (
	"omniport-api/internal/helper"
)

type BranchRequest struct {
	BranchCode  int64  `json:"branch_code" binding:"required"`
	BranchName  string `json:"branch_name" binding:"required"`
	KdPort      string `json:"kd_port"`
	Address     string `json:"address"`
	Status      string `json:"status"`
}

type SearchBranchRequest struct {
	Page     int                  `json:"page"`
	Limit    int                  `json:"limit"`
	Search   string               `json:"search"`
	Filters  map[string]string    `json:"filters"`
	Sort     helper.SortQuery     `json:"sort"`
	Download helper.DownloadQuery `json:"download"`
}

func (r SearchBranchRequest) ToPaginationQuery() helper.PaginationQuery {
	return helper.PaginationQuery{
		Page:     r.Page,
		Limit:    r.Limit,
		Search:   r.Search,
		Filters:  r.Filters,
		Sort:     r.Sort,
		Download: r.Download,
	}
}

type BranchStats struct {
	TotalBranches  int64 `json:"total_branches"`
	ActiveBranches int64 `json:"active_branches"`
}
