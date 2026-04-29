package postrequest

import "time"

// PostRequest maps to plan.post_requests
type PostRequest struct {
	ID               int64      `gorm:"primaryKey;autoIncrement"                     json:"id"`
	BranchCode       *int       `gorm:"column:branch_code"                           json:"branch_code"`
	TerminalCode     *int       `gorm:"column:terminal_code"                         json:"terminal_code"`
	BranchName       string     `gorm:"column:branch_name"                           json:"branch_name"`
	TerminalName     string     `gorm:"column:terminal_name"                         json:"terminal_name"`
	PPKNumber        string     `gorm:"column:ppk_number"                            json:"ppk_number"`
	VesselCode       string     `gorm:"column:vessel_code"                           json:"vessel_code"`
	VesselName       string     `gorm:"column:vessel_name"                           json:"vessel_name"`
	VesselType       string     `gorm:"column:vessel_type"                           json:"vessel_type"`
	VoyageType       string     `gorm:"column:voyage_type"                           json:"voyage_type"`
	AgentName        string     `gorm:"column:agent_name"                            json:"agent_name"`
	RequestCode      string     `gorm:"column:request_code"                          json:"request_code"`
	RequestDate      time.Time  `gorm:"column:request_date;not null"                 json:"request_date"`
	PBMCode          string     `gorm:"column:pbm_code"                              json:"pbm_code"`
	PBMName          string     `gorm:"column:pbm_name"                              json:"pbm_name"`
	NoBC11           string     `gorm:"column:no_bc11"                               json:"no_bc11"`
	DateBC11         *time.Time `gorm:"column:date_bc11"                             json:"date_bc11"`
	Description      string     `gorm:"column:description"                           json:"description"`
	Status           *int       `gorm:"column:status"                                json:"status"`
	PlanStatus       *int       `gorm:"column:plan_status"                           json:"plan_status"`
	ProgramName      string     `gorm:"column:program_name;not null"                 json:"program_name"`
	RefNumber        string     `gorm:"column:ref_number"                            json:"ref_number"`
	RefDate          *time.Time `gorm:"column:ref_date"                              json:"ref_date"`
	Ref1             string     `gorm:"column:ref1"                                  json:"ref1"`
	Ref2             string     `gorm:"column:ref2"                                  json:"ref2"`
	Val1             *float64   `gorm:"column:val1"                                  json:"val1"`
	Val2             *float64   `gorm:"column:val2"                                  json:"val2"`
	TotalManifest    *float64   `gorm:"column:total_manifest"                        json:"total_manifest"`
	BillableCode     string     `gorm:"column:billable_code"                         json:"billable_code"`
	BillableName     string     `gorm:"column:billable_name"                         json:"billable_name"`
	VesselCodeDst    string     `gorm:"column:vessel_code_dst"                       json:"vessel_code_dst"`
	VesselNameDst    string     `gorm:"column:vessel_name_dst"                       json:"vessel_name_dst"`
	ActivityCode     string     `gorm:"column:activity_code"                         json:"activity_code"`
	ActivityName     string     `gorm:"column:activity_name"                         json:"activity_name"`
	ToPPKNumber      string     `gorm:"column:to_ppk_number"                         json:"to_ppk_number"`
	ApprovalDate     *time.Time `gorm:"column:approval_date"                         json:"approval_date"`
	CreationDate     *time.Time `gorm:"column:creation_date"                         json:"creation_date"`
	CreationBy       string     `gorm:"column:creation_by"                           json:"creation_by"`
	LastUpdatedDate  *time.Time `gorm:"column:last_updated_date"                     json:"last_updated_date"`
	LastUpdatedBy    string     `gorm:"column:last_updated_by"                       json:"last_updated_by"`
}

func (PostRequest) TableName() string { return "plan.post_requests" }

// PostRequestDetail maps to plan.post_requests_d
type PostRequestDetail struct {
	ID                   int64      `gorm:"primaryKey;autoIncrement"                json:"id"`
	BranchCode           int        `gorm:"column:branch_code;not null"             json:"branch_code"`
	TerminalCode         int        `gorm:"column:terminal_code;not null"           json:"terminal_code"`
	BranchName           string     `gorm:"column:branch_name"                      json:"branch_name"`
	TerminalName         string     `gorm:"column:terminal_name"                    json:"terminal_name"`
	RequestCode          string     `gorm:"column:request_code"                     json:"request_code"`
	SequenceNumber       *int       `gorm:"column:sequence_number"                  json:"sequence_number"`
	StackingType         string     `gorm:"column:stacking_type"                    json:"stacking_type"`
	CargoCode            string     `gorm:"column:cargo_code;not null"              json:"cargo_code"`
	CargoName            string     `gorm:"column:cargo_name;not null"              json:"cargo_name"`
	CargoUnit            string     `gorm:"column:cargo_unit"                       json:"cargo_unit"`
	Total                *float64   `gorm:"column:total"                            json:"total"`
	QuantityMT           *float64   `gorm:"column:quantity_mt"                      json:"quantity_mt"`
	Quantity             *float64   `gorm:"column:quantity"                         json:"quantity"`
	CargoNature          string     `gorm:"column:cargo_nature;type:char(1)"        json:"cargo_nature"`
	CargoNatureDesc      string     `gorm:"column:cargo_nature_desc"                json:"cargo_nature_desc"`
	CargoPackaging       string     `gorm:"column:cargo_packaging"                  json:"cargo_packaging"`
	Stowage              string     `gorm:"column:stowage"                          json:"stowage"`
	PlannedDate          *time.Time `gorm:"column:planned_date"                     json:"planned_date"`
	WarehouseID          string     `gorm:"column:warehouse_id"                     json:"warehouse_id"`
	BLAWBNumber          string     `gorm:"column:bl_awb_number"                    json:"bl_awb_number"`
	BLAWBDate            *time.Time `gorm:"column:bl_awb_date"                      json:"bl_awb_date"`
	Description          string     `gorm:"column:description"                      json:"description"`
	PackageCount         *float64   `gorm:"column:package_count"                    json:"package_count"`
	OriginPortCode       string     `gorm:"column:origin_port_code"                 json:"origin_port_code"`
	DestinationPortCode  string     `gorm:"column:destination_port_code"            json:"destination_port_code"`
	OriginPortName       string     `gorm:"column:origin_port_name"                 json:"origin_port_name"`
	DestinationPortName  string     `gorm:"column:destination_port_name"            json:"destination_port_name"`
	StorageReference     string     `gorm:"column:storage_reference"                json:"storage_reference"`
	StorageStackDate     *time.Time `gorm:"column:storage_stack_date"               json:"storage_stack_date"`
	WarehouseDetailID    string     `gorm:"column:warehouse_detail_id"              json:"warehouse_detail_id"`
	WarehouseDetailName  string     `gorm:"column:warehouse_detail_name"            json:"warehouse_detail_name"`
	WarehouseName        string     `gorm:"column:warehouse_name"                   json:"warehouse_name"`
	ConsigneeCode        string     `gorm:"column:consignee_code"                   json:"consignee_code"`
	ConsigneeName        string     `gorm:"column:consignee_name"                   json:"consignee_name"`
	CreationDate         *time.Time `gorm:"column:creation_date"                    json:"creation_date"`
	CreationBy           string     `gorm:"column:creation_by"                      json:"creation_by"`
	LastUpdatedDate      *time.Time `gorm:"column:last_updated_date"                json:"last_updated_date"`
	LastUpdatedBy        string     `gorm:"column:last_updated_by"                  json:"last_updated_by"`
	ProgramName          string     `gorm:"column:program_name;not null"            json:"program_name"`
}

func (PostRequestDetail) TableName() string { return "plan.post_requests_d" }

// PostVesselSchedule maps to plan.post_vessel_schedules
type PostVesselSchedule struct {
	ID                  int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	BranchCode          *int       `gorm:"column:branch_code" json:"branch_code"`
	TerminalCode        *int       `gorm:"column:terminal_code" json:"terminal_code"`
	BranchName          string     `gorm:"column:branch_name" json:"branch_name"`
	TerminalName        string     `gorm:"column:terminal_name" json:"terminal_name"`
	ScheduleCode        string     `gorm:"column:schedule_code" json:"schedule_code"`
	PKKNumber           string     `gorm:"column:pkk_number" json:"pkk_number"`
	VesselName          string     `gorm:"column:vessel_name" json:"vessel_name"`
	VesselCode          string     `gorm:"column:vessel_code" json:"vessel_code"`
	VesselType          string     `gorm:"column:vessel_type" json:"vessel_type"`
	VoyageNumber        string     `gorm:"column:voyage_number;not null" json:"voyage_number"`
	VoyageType          string     `gorm:"column:voyage_type;not null" json:"voyage_type"`

	GRT                 *int       `gorm:"column:grt" json:"grt"`
	LOA                 *float64   `gorm:"column:loa" json:"loa"`
	AgencyName          string     `gorm:"column:agency_name" json:"agency_name"`
	PortAgent           string     `gorm:"column:port_agent" json:"port_agent"`
	EmergencyContact    string     `gorm:"column:emergency_contact" json:"emergency_contact"`
	OriginPortCode      string     `gorm:"column:origin_port_code" json:"origin_port_code"`
	OriginPortName      string     `gorm:"column:origin_port_name" json:"origin_port_name"`
	DestinationPortCode string     `gorm:"column:destination_port_code" json:"destination_port_code"`
	DestinationPortName string     `gorm:"column:destination_port_name" json:"destination_port_name"`
	DischargePortCode   string     `gorm:"column:discharge_port_code" json:"discharge_port_code"`
	DischargePortName   string     `gorm:"column:discharge_port_name" json:"discharge_port_name"`
	AssignedBerthName   string     `gorm:"column:assigned_berth_name" json:"assigned_berth_name"`
	DockID              *int       `gorm:"column:dock_id" json:"dock_id"`
	DockCode            string     `gorm:"column:dock_code" json:"dock_code"`
	DockName            string     `gorm:"column:dock_name" json:"dock_name"`
	BerthCode           string     `gorm:"column:berth_code" json:"berth_code"`
	BerthName           string     `gorm:"column:berth_name" json:"berth_name"`
	BerthPosition       string     `gorm:"column:berth_position" json:"berth_position"`
	PositionRange       string     `gorm:"column:position_range" json:"position_range"`
	ETA                 *time.Time `gorm:"column:eta" json:"eta"`
	ETB                 *time.Time `gorm:"column:etb" json:"etb"`
	ETC                 *time.Time `gorm:"column:etc" json:"etc"`
	ETD                 *time.Time `gorm:"column:etd" json:"etd"`
	Status              *int       `gorm:"column:status" json:"status"`
	CreationDate        *time.Time `gorm:"column:creation_date" json:"creation_date"`
	CreationBy          string     `gorm:"column:creation_by" json:"creation_by"`
	LastUpdatedDate     *time.Time `gorm:"column:last_updated_date" json:"last_updated_date"`
	LastUpdatedBy       string     `gorm:"column:last_updated_by" json:"last_updated_by"`
}

func (PostVesselSchedule) TableName() string { return "plan.post_vessel_schedules" }

