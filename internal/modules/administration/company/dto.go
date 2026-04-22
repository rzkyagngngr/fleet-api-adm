package company

import (
	"omniport-api/internal/helper"
	"time"
)

type CompanyResponse struct {
	ID              uint64     `json:"id"`
	CompanyCode     string     `json:"company_code"`
	CompanyName     string     `json:"company_name"`
	Npwp            string     `json:"npwp"`
	Address         string     `json:"address"`
	Email           string     `json:"email"`
	PhoneNumber     string     `json:"phone_number"`
	BusinessType    string     `json:"business_type"`
	Status          string     `json:"status"`
	CreatedBy       string     `json:"created_by"`
	CreatedDate     *time.Time `json:"created_date"`
	LastUpdatedBy   string     `json:"last_updated_by"`
	LastUpdatedDate *time.Time `json:"last_updated_date"`
	ProgramName     string     `json:"program_name"`
}

type CompanyRequest struct {
	CompanyCode  string `json:"company_code" binding:"required"`
	CompanyName  string `json:"company_name" binding:"required"`
	Npwp         string `json:"npwp"`
	Address      string `json:"address"`
	Email        string `json:"email"`
	PhoneNumber  string `json:"phone_number"`
	BusinessType string `json:"business_type"`
	Status       string `json:"status"`
}

type SearchCompaniesRequest struct {
	Page     int                  `json:"page"`
	Limit    int                  `json:"limit"`
	Search   string               `json:"search"`
	Filters  map[string]string    `json:"filters"`
	Sort     helper.SortQuery     `json:"sort"`
	Download helper.DownloadQuery `json:"download"`
}

func (r SearchCompaniesRequest) ToPaginationQuery() helper.PaginationQuery {
	return helper.PaginationQuery{
		Page:     r.Page,
		Limit:    r.Limit,
		Search:   r.Search,
		Filters:  r.Filters,
		Sort:     r.Sort,
		Download: r.Download,
	}
}

func ToResponse(c *Company) CompanyResponse {
	return CompanyResponse{
		ID:              c.ID,
		CompanyCode:     c.CompanyCode,
		CompanyName:     c.CompanyName,
		Npwp:            c.Npwp,
		Address:         c.Address,
		Email:           c.Email,
		PhoneNumber:     c.PhoneNumber,
		BusinessType:    c.BusinessType,
		Status:          c.Status,
		CreatedBy:       c.CreatedBy,
		CreatedDate:     c.CreatedDate,
		LastUpdatedBy:   c.LastUpdatedBy,
		LastUpdatedDate: c.LastUpdatedDate,
		ProgramName:     c.ProgramName,
	}
}
