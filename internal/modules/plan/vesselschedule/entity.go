package vesselschedule

import "time"

type VesselSchedule struct {
	ID                  uint64     `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	BranchCode          *int       `gorm:"column:branch_code" json:"branch_code"`
	TerminalCode        *int       `gorm:"column:terminal_code" json:"terminal_code"`
	BranchName          *string    `gorm:"column:branch_name;size:50" json:"branch_name"`
	TerminalName        *string    `gorm:"column:terminal_name;size:50" json:"terminal_name"`
	ScheduleCode        *string    `gorm:"column:schedule_code;size:20" json:"schedule_code"`
	VesselName          *string    `gorm:"column:vessel_name;size:100" json:"vessel_name"`
	VesselCode          *string    `gorm:"column:vessel_code;size:50" json:"vessel_code"`
	VesselType          *string    `gorm:"column:vessel_type;size:50" json:"vessel_type"`
	VoyageNumber        string     `gorm:"column:voyage_number;size:50;not null" json:"voyage_number"`
	PKKNumber           *string    `gorm:"column:pkk_number;size:50" json:"pkk_number"`
	VoyageType          string     `gorm:"column:voyage_type;size:50;not null" json:"voyage_type"`
	GRT                 *int       `gorm:"column:grt" json:"grt"`
	LOA                 *float64   `gorm:"column:loa" json:"loa"`
	AgencyName          *string    `gorm:"column:agency_name;size:100" json:"agency_name"`
	PortAgent           *string    `gorm:"column:port_agent;size:100" json:"port_agent"`
	EmergencyContact    *string    `gorm:"column:emergency_contact;size:50" json:"emergency_contact"`
	OriginPortCode      *string    `gorm:"column:origin_port_code;size:100" json:"origin_port_code"`
	OriginPortName      *string    `gorm:"column:origin_port_name;size:100" json:"origin_port_name"`
	DestinationPortCode *string    `gorm:"column:destination_port_code;size:100" json:"destination_port_code"`
	DestinationPortName *string    `gorm:"column:destination_port_name;size:100" json:"destination_port_name"`
	DischargePortCode   *string    `gorm:"column:discharge_port_code;size:100" json:"discharge_port_code"`
	DischargePortName   *string    `gorm:"column:discharge_port_name;size:100" json:"discharge_port_name"`
	AssignedBerthName   *string    `gorm:"column:assigned_berth_name;size:100" json:"assigned_berth_name"`
	DockID              *int       `gorm:"column:dock_id" json:"dock_id"`
	DockCode            *string    `gorm:"column:dock_code;size:50" json:"dock_code"`
	DockName            *string    `gorm:"column:dock_name;size:100" json:"dock_name"`
	BerthCode           *string    `gorm:"column:berth_code;size:50" json:"berth_code"`
	BerthName           *string    `gorm:"column:berth_name;size:100" json:"berth_name"`
	BerthPosition       *string    `gorm:"column:berth_position;size:100" json:"berth_position"`
	PositionRange       *string    `gorm:"column:position_range;size:50" json:"position_range"`
	ETA                 *time.Time `gorm:"column:eta" json:"eta"`
	ETB                 *time.Time `gorm:"column:etb" json:"etb"`
	ETC                 *time.Time `gorm:"column:etc" json:"etc"`
	ETD                 *time.Time `gorm:"column:etd" json:"etd"`
	Status              *int       `gorm:"column:status" json:"status"`
	CreationDate        *time.Time `gorm:"column:creation_date" json:"creation_date"`
	CreationBy          *string    `gorm:"column:creation_by;size:100" json:"creation_by"`
	LastUpdatedDate     *time.Time `gorm:"column:last_updated_date" json:"last_updated_date"`
	LastUpdatedBy       *string    `gorm:"column:last_updated_by;size:100" json:"last_updated_by"`
}

func (VesselSchedule) TableName() string { return "plan.post_vessel_schedules" }
