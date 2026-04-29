package tariff

import (
	"omniport-api/internal/helper"
	"time"
)

type TariffServiceReq struct {
	SequenceNo     *int     `json:"sequence_no"`
	ServiceType    *string  `json:"service_type" binding:"omitempty,max=50"`
	ServiceName    *string  `json:"service_name" binding:"omitempty,max=50"`
	CustomerName   *string  `json:"customer_name" binding:"omitempty,max=100"`
	CustomerCode   *string  `json:"customer_code" binding:"omitempty,max=25"`
	CargoCode      *string  `json:"cargo_code" binding:"omitempty,max=25"`
	CargoName      *string  `json:"cargo_name" binding:"omitempty,max=50"`
	CargoPackaging *string  `json:"cargo_packaging" binding:"omitempty,max=12"`
	CargoUnit      *string  `json:"cargo_unit" binding:"omitempty,max=20"`
	EquipmentCode  *string  `json:"equipment_code" binding:"omitempty,max=25"`
	EquipmentName  *string  `json:"equipment_name" binding:"omitempty,max=50"`
	EquipmentGroup *string  `json:"equipment_group" binding:"omitempty,max=255"`
	EquipmentUnit  *string  `json:"equipment_unit" binding:"omitempty,max=255"`
	BaseTariff     *float64 `json:"base_tariff"`
	CurrencyCode   *string  `json:"currency_code" binding:"omitempty,max=25"`
	Discount       *float64 `json:"discount"`
	Attrib1        *string  `json:"attrib1" binding:"omitempty,max=50"`
	Attrib2        *string  `json:"attrib2" binding:"omitempty,max=50"`
	Attrib3        *string  `json:"attrib3" binding:"omitempty,max=50"`
}

type TariffReq struct {
	NameTariff      string             `json:"name_tariff" binding:"required,max=100"`
	Description     *string            `json:"description"`
	Status          *int               `json:"status"`
	AgreementNumber *string            `json:"agreement_number" binding:"omitempty,max=100"`
	StartDate       *time.Time         `json:"start_date"`
	EndDate         *time.Time         `json:"end_date"`
	Details         []TariffServiceReq `json:"details"`
}

type UpdateTariffStatusRequest struct {
	Status *int `json:"status" binding:"required"`
}

type SearchTariffRequest struct {
	Page     int                  `json:"page"`
	Limit    int                  `json:"limit"`
	Search   string               `json:"search"`
	Filters  map[string]string    `json:"filters"`
	Sort     helper.SortQuery     `json:"sort"`
	Download helper.DownloadQuery `json:"download"`
}

func (r SearchTariffRequest) ToPaginationQuery() helper.PaginationQuery {
	return helper.PaginationQuery{
		Page:     r.Page,
		Limit:    r.Limit,
		Search:   r.Search,
		Filters:  r.Filters,
		Sort:     r.Sort,
		Download: r.Download,
	}
}
