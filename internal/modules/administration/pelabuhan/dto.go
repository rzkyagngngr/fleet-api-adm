package pelabuhan

import "omniport-api/internal/helper"

type PortReq struct {
	PortCode    string  `json:"port_code" binding:"required,max=10"`
	PortName    string  `json:"port_name" binding:"max=40"`
	PortCity    string  `json:"port_city" binding:"max=40"`
	CountryCode string  `json:"country_code" binding:"max=3"`
	Status      *string `json:"status" binding:"omitempty,oneof=A I"` // A = Active, I = Inactive
}

type SearchPortRequest struct {
	Page     int                  `json:"page"`
	Limit    int                  `json:"limit"`
	Search   string               `json:"search"`
	Filters  map[string]string    `json:"filters"`
	Sort     helper.SortQuery     `json:"sort"`
	Download helper.DownloadQuery `json:"download"`
}

func (r SearchPortRequest) ToPaginationQuery() helper.PaginationQuery {
	return helper.PaginationQuery{
		Page:     r.Page,
		Limit:    r.Limit,
		Search:   r.Search,
		Filters:  r.Filters,
		Sort:     r.Sort,
		Download: r.Download,
	}
}
