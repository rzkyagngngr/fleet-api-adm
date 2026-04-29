package vesselschedule

import (
	"omniport-api/internal/helper"
	"time"
)

type VesselScheduleRequest struct {
	VesselName          *string    `json:"vessel_name" binding:"omitempty,max=100"`
	VesselCode          *string    `json:"vessel_code" binding:"omitempty,max=50"`
	VesselType          *string    `json:"vessel_type" binding:"omitempty,max=50"`
	VoyageNo            string     `json:"voyage_no" binding:"required,max=50"`
	GRT                 *int       `json:"grt"`
	LOA                 *float64   `json:"loa"`
	AgencyName          *string    `json:"agency_name" binding:"omitempty,max=100"`
	PortAgent           *string    `json:"port_agent" binding:"omitempty,max=100"`
	EmergencyContact    *string    `json:"emergency_contact" binding:"omitempty,max=50"`
	OriginPortCode      *string    `json:"origin_port_code" binding:"omitempty,max=100"`
	OriginPortName      *string    `json:"origin_port_name" binding:"omitempty,max=100"`
	DestinationPortCode *string    `json:"destination_port_code" binding:"omitempty,max=100"`
	DestinationPortName *string    `json:"destination_port_name" binding:"omitempty,max=100"`
	DischargePortCode   *string    `json:"discharge_port_code" binding:"omitempty,max=100"`
	DischargePortName   *string    `json:"discharge_port_name" binding:"omitempty,max=100"`
	AssignedBerthName   *string    `json:"assigned_berth_name" binding:"omitempty,max=100"`
	DockID              *int       `json:"dock_id"`
	DockCode            *string    `json:"dock_code" binding:"omitempty,max=50"`
	DockName            *string    `json:"dock_name" binding:"omitempty,max=100"`
	BerthCode           *string    `json:"berth_code" binding:"omitempty,max=50"`
	BerthName           *string    `json:"berth_name" binding:"omitempty,max=100"`
	BerthPosition       *string    `json:"berth_position" binding:"omitempty,max=100"`
	PositionRange       *string    `json:"position_range" binding:"omitempty,max=50"`
	ETA                 *time.Time `json:"eta"`
	ETB                 *time.Time `json:"etb"`
	ETC                 *time.Time `json:"etc"`
	ETD                 *time.Time `json:"etd"`
	Status              *int       `json:"status"`
}

type SearchVesselScheduleRequest struct {
	Page     int                  `json:"page"`
	Limit    int                  `json:"limit"`
	Search   string               `json:"search"`
	Filters  map[string]string    `json:"filters"`
	Sort     helper.SortQuery     `json:"sort"`
	Download helper.DownloadQuery `json:"download"`
}

func (r SearchVesselScheduleRequest) ToPaginationQuery() helper.PaginationQuery {
	return helper.PaginationQuery{
		Page:     r.Page,
		Limit:    r.Limit,
		Search:   r.Search,
		Filters:  r.Filters,
		Sort:     r.Sort,
		Download: r.Download,
	}
}
