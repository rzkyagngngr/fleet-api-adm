package cargo

import (
	"omniport-api/internal/helper"
	"time"
)

type CargoResponse struct {
	ID                    uint64     `json:"id"`
	BranchCode            int        `json:"branch_code"`
	TerminalCode          int        `json:"terminal_code"`
	CargoCode             string     `json:"cargo_code"`
	CargoSitcCode         string     `json:"cargo_sitc_code"`
	CargoHsHarmonizedCode string     `json:"cargo_hs_harmonized_code"`
	CargoName             string     `json:"cargo_name"`
	CargoGroup            string     `json:"cargo_group"`
	CargoCommodity        string     `json:"cargo_commodity"`
	CargoCharacteristic   string     `json:"cargo_characteristic"`
	CargoImdgCode         *int16     `json:"cargo_imdg_code"`
	CargoImdgDescription  string     `json:"cargo_imdg_description"`
	CargoPackaging1       string     `json:"cargo_packaging_1"`
	CargoConversion1      float64    `json:"cargo_conversion_1"`
	CargoDimension1       float64    `json:"cargo_dimension_1"`
	CargoUnit1            string     `json:"cargo_unit_1"`
	CargoPackaging2       string     `json:"cargo_packaging_2"`
	CargoConversion2      float64    `json:"cargo_conversion_2"`
	CargoDimension2       float64    `json:"cargo_dimension_2"`
	CargoUnit2            string     `json:"cargo_unit_2"`
	CargoPackaging3       string     `json:"cargo_packaging_3"`
	CargoConversion3      float64    `json:"cargo_conversion_3"`
	CargoDimension3       float64    `json:"cargo_dimension_3"`
	CargoUnit3            string     `json:"cargo_unit_3"`
	CargoMooringType      string     `json:"cargo_mooring_type"`
	CargoNotes            string     `json:"cargo_notes"`
	CargoCommodityGroup   string     `json:"cargo_commodity_group"`
	CargoCommodityType    string     `json:"cargo_commodity_type"`
	IsActive              string     `json:"is_active"`
	CargoDocument         string     `json:"cargo_document"`
	HsCode                string     `json:"hs_code"`
	HsDescription         string     `json:"hs_description"`
	CargoProductName      string     `json:"cargo_product_name"`
	CreationDate          *time.Time `json:"creation_date"`

	// Mapping for Frontend Compatibility (Legacy/Atomic)
	ItemCode         string `json:"item_code"`
	ItemName         string `json:"item_name"`
	IsDangerousGoods int    `json:"is_dangerous_goods"`
	Status           int    `json:"status"`
}

type CargoRequest struct {
	BranchCode            int     `json:"branch_code"`
	TerminalCode          int     `json:"terminal_code"`
	CargoCode             string  `json:"cargo_code"`
	CargoSitcCode         string  `json:"cargo_sitc_code"`
	CargoHsHarmonizedCode string  `json:"cargo_hs_harmonized_code"`
	CargoName             string  `json:"cargo_name"`
	CargoGroup            string  `json:"cargo_group"`
	CargoCommodity        string  `json:"cargo_commodity"`
	CargoCharacteristic   string  `json:"cargo_characteristic"`
	CargoImdgCode         *int16  `json:"cargo_imdg_code"`
	CargoImdgDescription  string  `json:"cargo_imdg_description"`
	CargoPackaging1       string  `json:"cargo_packaging_1"`
	CargoConversion1      float64 `json:"cargo_conversion_1"`
	CargoDimension1       float64 `json:"cargo_dimension_1"`
	CargoUnit1            string  `json:"cargo_unit_1"`
	CargoPackaging2       string  `json:"cargo_packaging_2"`
	CargoConversion2      float64 `json:"cargo_conversion_2"`
	CargoDimension2       float64 `json:"cargo_dimension_2"`
	CargoUnit2            string  `json:"cargo_unit_2"`
	CargoPackaging3       string  `json:"cargo_packaging_3"`
	CargoConversion3      float64 `json:"cargo_conversion_3"`
	CargoDimension3       float64 `json:"cargo_dimension_3"`
	CargoUnit3            string  `json:"cargo_unit_3"`
	CargoMooringType      string  `json:"cargo_mooring_type"`
	CargoNotes            string  `json:"cargo_notes"`
	CargoCommodityGroup   string  `json:"cargo_commodity_group"`
	CargoCommodityType    string  `json:"cargo_commodity_type"`
	IsActive              string  `json:"is_active"`
	CargoDocument         string  `json:"cargo_document"`
	HsCode                string  `json:"hs_code"`
	HsDescription         string  `json:"hs_description"`
	CargoProductName      string  `json:"cargo_product_name"`

	// Legacy Mapping for Compatibility
	ItemCode         string `json:"item_code"`
	ItemName         string `json:"item_name"`
	IsDangerousGoods int    `json:"is_dangerous_goods"`
	Status           int    `json:"status"`
	Category         string `json:"category"`
	UOM              string `json:"uom"`
	StorageType      string `json:"storage_type"`
	Brand            string `json:"brand"`
	Remark           string `json:"remark"`
}

type SearchCargoRequest struct {
	Page     int                  `json:"page"`
	Limit    int                  `json:"limit"`
	Search   string               `json:"search"`
	Filters  map[string]string    `json:"filters"`
	Sort     helper.SortQuery     `json:"sort"`
	Download helper.DownloadQuery `json:"download"`
}

func (r SearchCargoRequest) ToPaginationQuery() helper.PaginationQuery {
	return helper.PaginationQuery{
		Page:     r.Page,
		Limit:    r.Limit,
		Search:   r.Search,
		Filters:  r.Filters,
		Sort:     r.Sort,
		Download: r.Download,
	}
}

type CargoStatsResponse struct {
	TotalCargoMasters  int64 `json:"total_cargo_masters"`
	ActiveCommodities int64 `json:"active_commodities"`
	HazmatRegistry     int64 `json:"hazmat_registry"`
}

func (c *Cargo) ToResponse() CargoResponse {
	isDG := 0
	if c.CargoImdgCode != nil && *c.CargoImdgCode > 0 {
		isDG = 1
	}

	status := 0
	if c.IsActive == "1" || c.IsActive == "Y" {
		status = 1
	}

	return CargoResponse{
		ID:                    c.ID,
		BranchCode:            c.BranchCode,
		TerminalCode:          c.TerminalCode,
		CargoCode:             c.CargoCode,
		CargoSitcCode:         c.CargoSitcCode,
		CargoHsHarmonizedCode: c.CargoHsHarmonizedCode,
		CargoName:             c.CargoName,
		CargoGroup:            c.CargoGroup,
		CargoCommodity:        c.CargoCommodity,
		CargoCharacteristic:   c.CargoCharacteristic,
		CargoImdgCode:         c.CargoImdgCode,
		CargoImdgDescription:  c.CargoImdgDescription,
		CargoPackaging1:       c.CargoPackaging1,
		CargoConversion1:      c.CargoConversion1,
		CargoDimension1:       c.CargoDimension1,
		CargoUnit1:            c.CargoUnit1,
		CargoPackaging2:       c.CargoPackaging2,
		CargoConversion2:      c.CargoConversion2,
		CargoDimension2:       c.CargoDimension2,
		CargoUnit2:            c.CargoUnit2,
		CargoPackaging3:       c.CargoPackaging3,
		CargoConversion3:      c.CargoConversion3,
		CargoDimension3:       c.CargoDimension3,
		CargoUnit3:            c.CargoUnit3,
		CargoMooringType:      c.CargoMooringType,
		CargoNotes:            c.CargoNotes,
		CargoCommodityGroup:   c.CargoCommodityGroup,
		CargoCommodityType:    c.CargoCommodityType,
		IsActive:              c.IsActive,
		CargoDocument:         c.CargoDocument,
		HsCode:                c.HsCode,
		HsDescription:         c.HsDescription,
		CargoProductName:      c.CargoProductName,
		CreationDate:          c.CreatedDate,

		// Legacy Mappings
		ItemCode:         c.CargoCode,
		ItemName:         c.CargoName,
		IsDangerousGoods: isDG,
		Status:           status,
	}
}
