package vessel

import (
	"omniport-api/internal/helper"
	"time"
)

type VesselResponse struct {
	ID                    uint64     `json:"id"`
	VesselCode            string     `json:"vessel_code"`
	VesselName            string     `json:"vessel_name"`
	VesselType            string     `json:"vessel_type"`
	VesselCallSign        string     `json:"vessel_call_sign"`
	VesselImo             string     `json:"vessel_imo"`
	VesselGrt             string     `json:"vessel_grt"`
	VesselLoa             string     `json:"vessel_loa"`
	VesselOwnerName       string     `json:"vessel_owner_name"`
	VesselShippingRoute   string     `json:"vessel_shipping_route"`
	VesselFlag            string     `json:"vessel_flag"`
	VesselCountry         string     `json:"vessel_country"`
	VesselYearMade        string     `json:"vessel_year_made"`
	VesselHatchNumber     int        `json:"vessel_hatch_number"`
	VesselHatchType       string     `json:"vessel_hatch_type"`
	VesselOwnershipStatus string     `json:"vessel_ownership_status"`
	VesselOperationStatus string     `json:"vessel_operation_status"`
	Status                string     `json:"status"`
	Remark                string     `json:"remark"`
	PortCode              int64      `json:"port_code"`
	BranchCode            int64      `json:"branch_code"`
	TerminalCode          int64      `json:"terminal_code"`
	CreationDate          time.Time  `json:"creation_date"`
}

type VesselStatsResponse struct {
	TotalFleet     int64 `json:"total_fleet"`
	ActiveVessels  int64 `json:"active_vessels"`
	Maintenance    int64 `json:"maintenance"`
	Deactivated    int64 `json:"deactivated"`
	CargoCount     int64 `json:"cargo_count"`
	TankerCount    int64 `json:"tanker_count"`
	ContainerCount int64 `json:"container_count"`
	OtherCount     int64 `json:"other_count"`
}

type VesselRequest struct {
	VesselCode            string `json:"vessel_code" binding:"required"`
	VesselName            string `json:"vessel_name" binding:"required"`
	VesselType            string `json:"vessel_type" binding:"required"`
	VesselCallSign        string `json:"vessel_call_sign"`
	VesselImo             string `json:"vessel_imo"`
	VesselGrt             string `json:"vessel_grt"`
	VesselLoa             string `json:"vessel_loa"`
	VesselOwnerName       string `json:"vessel_owner_name"`
	VesselShippingRoute   string `json:"vessel_shipping_route"`
	VesselFlag            string `json:"vessel_flag"`
	VesselCountry         string `json:"vessel_country"`
	VesselYearMade        string `json:"vessel_year_made"`
	VesselHatchNumber     int    `json:"vessel_hatch_number"`
	VesselHatchType       string `json:"vessel_hatch_type"`
	VesselOwnershipStatus string `json:"vessel_ownership_status"`
	VesselOperationStatus string `json:"vessel_operation_status"`
	Status                string `json:"status"`
	Remark                string `json:"remark"`
	PortCode              int64  `json:"port_code"`
	BranchCode            int64  `json:"branch_code"`
	TerminalCode          int64  `json:"terminal_code"`
}

type SearchVesselsRequest struct {
	Page     int                  `json:"page"`
	Limit    int                  `json:"limit"`
	Search   string               `json:"search"`
	Filters  map[string]string    `json:"filters"`
	Sort     helper.SortQuery     `json:"sort"`
	Download helper.DownloadQuery `json:"download"`
}

func (r SearchVesselsRequest) ToPaginationQuery() helper.PaginationQuery {
	return helper.PaginationQuery{
		Page:     r.Page,
		Limit:    r.Limit,
		Search:   r.Search,
		Filters:  r.Filters,
		Sort:     r.Sort,
		Download: r.Download,
	}
}

func ToResponse(v *Vessel) VesselResponse {
	return VesselResponse{
		ID:                    v.ID,
		VesselCode:            v.VesselCode,
		VesselName:            v.VesselName,
		VesselType:            v.VesselType,
		VesselCallSign:        v.VesselCallSign,
		VesselImo:             v.VesselImo,
		VesselGrt:             v.VesselGrt,
		VesselLoa:             v.VesselLoa,
		VesselOwnerName:       v.VesselOwnerName,
		VesselShippingRoute:   v.VesselShippingRoute,
		VesselFlag:            v.VesselFlag,
		VesselCountry:         v.VesselCountry,
		VesselYearMade:        v.VesselYearMade,
		VesselHatchNumber:     v.VesselHatchNumber,
		VesselHatchType:       v.VesselHatchType,
		VesselOwnershipStatus: v.VesselOwnershipStatus,
		VesselOperationStatus: v.VesselOperationStatus,
		Status:                v.Status,
		Remark:                v.Remark,
		PortCode:              v.PortCode,
		BranchCode:            v.BranchCode,
		TerminalCode:          v.TerminalCode,
		CreationDate:          v.CreationDate,
	}
}
