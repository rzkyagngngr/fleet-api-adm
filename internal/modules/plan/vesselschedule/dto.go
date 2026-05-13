package vesselschedule

import (
	"omniport-api/internal/helper"
	"time"
)

type VesselScheduleRequest struct {
	VesselName             *string    `json:"vessel_name" binding:"omitempty,max=100"`
	VesselCode             *string    `json:"vessel_code" binding:"omitempty,max=50"`
	VesselType             *string    `json:"vessel_type" binding:"omitempty,max=50"`
	VesselHatchNumber      *int       `json:"vessel_hatch_number"`
	VoyageNumber           string     `json:"voyage_number" binding:"required,max=50"`
	PKKNumber              *string    `json:"pkk_number" binding:"omitempty,max=50"`
	PPKNumber              *string    `json:"ppk_number" binding:"omitempty,max=100"`
	VoyageType             string     `json:"voyage_type" binding:"required,max=50"`
	GRT                    *int       `json:"grt"`
	LOA                    *float64   `json:"loa"`
	AgencyName             *string    `json:"agency_name" binding:"omitempty,max=100"`
	PortAgent              *string    `json:"port_agent" binding:"omitempty,max=100"`
	EmergencyContact       *string    `json:"emergency_contact" binding:"omitempty,max=50"`
	OriginPortCode         *string    `json:"origin_port_code" binding:"omitempty,max=100"`
	OriginPortName         *string    `json:"origin_port_name" binding:"omitempty,max=100"`
	DestinationPortCode    *string    `json:"destination_port_code" binding:"omitempty,max=100"`
	DestinationPortName    *string    `json:"destination_port_name" binding:"omitempty,max=100"`
	DischargePortCode      *string    `json:"discharge_port_code" binding:"omitempty,max=100"`
	DischargePortName      *string    `json:"discharge_port_name" binding:"omitempty,max=100"`
	AssignedBerthName      *string    `json:"assigned_berth_name" binding:"omitempty,max=100"`
	DockID                 *int       `json:"dock_id"`
	DockCode               *string    `json:"dock_code" binding:"omitempty,max=50"`
	DockName               *string    `json:"dock_name" binding:"omitempty,max=100"`
	BerthCode              *string    `json:"berth_code" binding:"omitempty,max=50"`
	BerthName              *string    `json:"berth_name" binding:"omitempty,max=100"`
	BerthLatitude          *string    `json:"berth_latitude" binding:"omitempty,max=50"`
	BerthLongitude         *string    `json:"berth_longitude" binding:"omitempty,max=50"`
	CodeInaportnet         *string    `json:"code_inaportnet" binding:"omitempty,max=100"`
	LocationNameInaportnet *string    `json:"location_name_inaportnet" binding:"omitempty,max=100"`
	StartBerthPosition     *string    `json:"start_berth_position" binding:"omitempty,max=100"`
	EndBerthPosition       *string    `json:"end_berth_position" binding:"omitempty,max=100"`
	ETA                    *time.Time `json:"eta"`
	ETB                    *time.Time `json:"etb"`
	ETC                    *time.Time `json:"etc"`
	ETD                    *time.Time `json:"etd"`
	Status                 *int       `json:"status"`
}

type SearchVesselScheduleRequest struct {
	Page     int                  `json:"page"`
	Limit    int                  `json:"limit"`
	Search   string               `json:"search"`
	Filters  map[string]string    `json:"filters"`
	Sort     helper.SortQuery     `json:"sort"`
	Download helper.DownloadQuery `json:"download"`
}

type VesselScheduleSearchResponse struct {
	VesselSchedule
	Plans []VesselSchedulePlanResponse `json:"plans"`
}

type VesselScheduleDetailResponse struct {
	VesselSchedule
	Vessel       interface{}   `json:"vessel"`
	HatchDetails []interface{} `json:"hatch_details"`
}

type VesselSchedulePlanResponse struct {
	PPKNumber         string     `json:"ppk_number" gorm:"column:ppk_number"`
	PlanCode          string     `json:"plan_code" gorm:"column:plan_code"`
	PlanDate          *time.Time `json:"plan_date" gorm:"column:plan_date"`
	ActivityCode      string     `json:"activity_code" gorm:"column:activity_code"`
	ActivityName      string     `json:"activity_name" gorm:"column:activity_name"`
	ActivityStartDate *time.Time `json:"activity_start_date" gorm:"column:activity_start_date"`
	ActivityEndDate   *time.Time `json:"activity_end_date" gorm:"column:activity_end_date"`
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

type UpdateVesselScheduleStatusRequest struct {
	ScheduleCode string `json:"schedule_code" binding:"required"`
	Status       *int   `json:"status" binding:"required"`
}

type InitChatGroupRequest struct {
	ScheduleCode string `json:"schedule_code" binding:"required"`
}
