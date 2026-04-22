package terminal

import (
	"omniport-api/internal/helper"
	"time"
)

type TerminalRequest struct {
	BranchCode    string     `json:"branch_code" binding:"required"`
	TerminalCode  string     `json:"terminal_code" binding:"required"`
	TerminalName  string     `json:"terminal_name" binding:"required"`
	GoLiveDate    *time.Time `json:"go_live_date"`
	IsGoLive      string     `json:"is_go_live"`
	ProfitCenter  string     `json:"profit_center"`
	Latitude      string     `json:"latitude"`
	Longitude     string     `json:"longitude"`
	Status        string     `json:"status"`
	VersionCode   int64      `json:"version_code"`
	VersionName   string     `json:"version_name"`
	DocumentCode  string     `json:"document_code"`
	VesselVersion int64      `json:"vessel_version"`
	LogoURL       string     `json:"logo_url"`
	LogoMiniURL   string     `json:"logo_mini_url"`
	Address       string     `json:"address"`
	CompanyType   int64      `json:"company_type"`
	PortCode      string     `json:"port_code"`
	CompanyCode   string     `json:"company_code"`
	CompanyName   string     `json:"company_name"`
}

type SearchTerminalRequest struct {
	Page     int                  `json:"page"`
	Limit    int                  `json:"limit"`
	Search   string               `json:"search"`
	Filters  map[string]string    `json:"filters"`
	Sort     helper.SortQuery     `json:"sort"`
	Download helper.DownloadQuery `json:"download"`
}

func (r SearchTerminalRequest) ToPaginationQuery() helper.PaginationQuery {
	return helper.PaginationQuery{
		Page:     r.Page,
		Limit:    r.Limit,
		Search:   r.Search,
		Filters:  r.Filters,
		Sort:     r.Sort,
		Download: r.Download,
	}
}

type TerminalStats struct {
	TotalTerminals int64 `json:"total_terminals"`
	GoLiveTerminals int64 `json:"go_live_terminals"`
}
