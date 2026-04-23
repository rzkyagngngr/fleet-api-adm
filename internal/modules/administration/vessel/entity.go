package vessel

import "time"

type Vessel struct {
	ID                    uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	VesselCode            string     `gorm:"column:vessel_code;not null" json:"vessel_code"`
	VesselName            string     `gorm:"column:vessel_name;not null" json:"vessel_name"`
	VesselType            string     `gorm:"column:vessel_type;not null" json:"vessel_type"`
	VesselCallSign        string     `gorm:"column:vessel_call_sign" json:"vessel_call_sign"`
	VesselImo             string     `gorm:"column:vessel_imo" json:"vessel_imo"`
	VesselGrt             string     `gorm:"column:vessel_grt" json:"vessel_grt"`
	VesselLoa             string     `gorm:"column:vessel_loa" json:"vessel_loa"`
	VesselOwnerName       string     `gorm:"column:vessel_owner_name" json:"vessel_owner_name"`
	VesselShippingRoute   string     `gorm:"column:vessel_shipping_route" json:"vessel_shipping_route"`
	VesselFlag            string     `gorm:"column:vessel_flag" json:"vessel_flag"`
	VesselCountry         string     `gorm:"column:vessel_country" json:"vessel_country"`
	VesselYearMade        string     `gorm:"column:vessel_year_made" json:"vessel_year_made"`
	VesselHatchNumber     int        `gorm:"column:vessel_hatch_number" json:"vessel_hatch_number"`
	VesselHatchType       string     `gorm:"column:vessel_hatch_type" json:"vessel_hatch_type"`
	VesselOwnershipStatus string     `gorm:"column:vessel_ownership_status" json:"vessel_ownership_status"`
	VesselOperationStatus string     `gorm:"column:vessel_operation_status" json:"vessel_operation_status"`
	Status                string     `gorm:"column:status;default:ACTIVE" json:"status"`
	Remark                string     `gorm:"column:remark" json:"remark"`
	PortCode              int64      `gorm:"column:port_code" json:"port_code"`
	BranchCode            int64      `gorm:"column:branch_code" json:"branch_code"`
	TerminalCode          int64      `gorm:"column:terminal_code" json:"terminal_code"`
	CreationDate          time.Time  `gorm:"column:creation_date;default:CURRENT_TIMESTAMP" json:"creation_date"`
	CreationBy            string     `gorm:"column:creation_by" json:"creation_by"`
	LastUpdatedDate       *time.Time `gorm:"column:last_updated_date" json:"last_updated_date"`
	LastUpdatedBy         string     `gorm:"column:last_updated_by" json:"last_updated_by"`
}

func (Vessel) TableName() string { return "posm_vessel" }
