package op

import "time"

type LoadingUnloadingPlan struct {
	ID                int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	BranchCode        int        `gorm:"column:branch_code;not null" json:"branch_code"`
	TerminalCode      int        `gorm:"column:terminal_code;not null" json:"terminal_code"`
	BranchName        string     `gorm:"column:branch_name" json:"branch_name"`
	TerminalName      string     `gorm:"column:terminal_name" json:"terminal_name"`
	VesselCode        string     `gorm:"column:vessel_code" json:"vessel_code"`
	VesselName        string     `gorm:"column:vessel_name" json:"vessel_name"`
	VesselType        string     `gorm:"column:vessel_type" json:"vessel_type"`
	GRT               *float64   `gorm:"column:grt" json:"grt"`
	LOA               *float64   `gorm:"column:loa" json:"loa"`
	ShippingType      string     `gorm:"column:shipping_type" json:"shipping_type"`
	AgentName         string     `gorm:"column:agent_name" json:"agent_name"`
	PPKNumber         string     `gorm:"column:ppk_number" json:"ppk_number"`
	PlanNumber        string     `gorm:"column:plan_code;not null" json:"plan_code"`
	PlanDate          time.Time  `gorm:"column:plan_date;not null" json:"plan_date"`
	ETA               *time.Time `gorm:"column:eta" json:"eta"`
	ETD               *time.Time `gorm:"column:etd" json:"etd"`
	BilledTo          *int       `gorm:"column:billed_to" json:"billed_to"`
	AssignedTo        *int       `gorm:"column:assigned_to" json:"assigned_to"`
	ActivityCode      string     `gorm:"column:activity_code" json:"activity_code"`
	ActivityName      string     `gorm:"column:activity_name" json:"activity_name"`
	Remarks           string     `gorm:"column:remarks" json:"remarks"`
	Status            *int       `gorm:"column:status" json:"status"`
	Cycle             string     `gorm:"column:cycle" json:"cycle"`
	TotalDays         *int       `gorm:"column:total_days" json:"total_days"`
	TotalShifts       *int       `gorm:"column:total_shifts" json:"total_shifts"`
	ActivityStartDate *time.Time `gorm:"column:activity_start_date" json:"activity_start_date"`
	ActivityEndDate   *time.Time `gorm:"column:activity_end_date" json:"activity_end_date"`
	VesselFacing      string     `gorm:"column:vessel_facing" json:"vessel_facing"`
	MooringLimit      *float64   `gorm:"column:mooring_limit" json:"mooring_limit"`
	BT                *float64   `gorm:"column:bt" json:"bt"`
	CreationDate      time.Time  `gorm:"column:creation_date;not null" json:"creation_date"`
	CreationBy        string     `gorm:"column:creation_by;not null" json:"creation_by"`
	LastUpdatedDate   *time.Time `gorm:"column:last_updated_date" json:"last_updated_date"`
	LastUpdatedBy     string     `gorm:"column:last_updated_by" json:"last_updated_by"`
	ProgramName       string     `gorm:"column:program_name;not null" json:"program_name"`
	VesselRpkID       uint64     `gorm:"-" json:"vessel_rpk_id"`
}

func (LoadingUnloadingPlan) TableName() string { return "plan.post_vessel_plan" }

type LoadingUnloadingPlanDetail struct {
	ID                int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	BranchCode        int        `gorm:"column:branch_code;not null" json:"branch_code"`
	TerminalCode      int        `gorm:"column:terminal_code;not null" json:"terminal_code"`
	PlanNumber        string     `gorm:"column:plan_code;not null" json:"plan_code"`
	SequenceNo        int        `gorm:"column:sequence_no;not null" json:"sequence_no"`
	ActivityDate      *time.Time `gorm:"column:activity_date" json:"activity_date"`
	Stowage           string     `gorm:"column:stowage" json:"stowage"`
	CargoCode         string     `gorm:"column:cargo_code" json:"cargo_code"`
	CargoName         string     `gorm:"column:cargo_name" json:"cargo_name"`
	TotalQuantity     *float64   `gorm:"column:total_quantity" json:"total_quantity"`
	PlannedQuantity   *float64   `gorm:"column:planned_quantity" json:"planned_quantity"`
	CargoUnit         string     `gorm:"column:cargo_unit" json:"cargo_unit"`
	CargoPackaging    string     `gorm:"column:cargo_packaging" json:"cargo_packaging"`
	DayNo             *int       `gorm:"column:day_no" json:"day_no"`
	BerthCode         string     `gorm:"column:berth_code" json:"berth_code"`
	BerthName         string     `gorm:"column:berth_name" json:"berth_name"`
	DockCode          string     `gorm:"column:dock_code" json:"dock_code"`
	DockName          string     `gorm:"column:dock_name" json:"dock_name"`
	Shift1            string     `gorm:"column:shift_1" json:"shift_1"`
	Shift2            string     `gorm:"column:shift_2" json:"shift_2"`
	Shift3            string     `gorm:"column:shift_3" json:"shift_3"`
	PBMCode           string     `gorm:"column:pbm_code" json:"pbm_code"`
	PBMName           string     `gorm:"column:pbm_name" json:"pbm_name"`
	ConsigneeCode     string     `gorm:"column:consignee_code" json:"consignee_code"`
	ConsigneeName     string     `gorm:"column:consignee_name" json:"consignee_name"`
	TruckCount        *int       `gorm:"column:truck_count" json:"truck_count"`
	TruckCapacity     *int       `gorm:"column:truck_capacity" json:"truck_capacity"`
	GangCount         *int       `gorm:"column:gang_count" json:"gang_count"`
	EquipmentCode     string     `gorm:"column:equipment_code" json:"equipment_code"`
	EquipmentName     string     `gorm:"column:equipment_name" json:"equipment_name"`
	EquipmentGroup    string     `gorm:"column:equipment_group" json:"equipment_group"`
	Attrib1           string     `gorm:"column:attrib1" json:"attrib1"`
	Attrib2           string     `gorm:"column:attrib2" json:"attrib2"`
	Attrib3           string     `gorm:"column:attrib3" json:"attrib3"`
	Val1              *float64   `gorm:"column:val1" json:"val1"`
	Val2              *float64   `gorm:"column:val2" json:"val2"`
	Val3              *float64   `gorm:"column:val3" json:"val3"`
	Status            *int       `gorm:"column:status" json:"status"`
	CargoNature       string     `gorm:"column:cargo_nature" json:"cargo_nature"`
	CargoNatureDesc   string     `gorm:"column:cargo_nature_desc" json:"cargo_nature_desc"`
	PlanDetailCode    string     `gorm:"column:plan_detail_code" json:"plan_detail_code"`
	DeterminationCode string     `gorm:"column:confirmed_plan_code" json:"confirmed_plan_code"`
	WorkOrderCode     string     `gorm:"column:work_order_code" json:"work_order_code"`
	FromDockCode      string     `gorm:"column:from_dock_code" json:"from_dock_code"`
	FromBerthCode     string     `gorm:"column:from_berth_code" json:"from_berth_code"`
	ProgramName       string     `gorm:"column:program_name;not null" json:"program_name"`
	CreationDate      time.Time  `gorm:"column:creation_date;not null" json:"creation_date"`
	CreationBy        string     `gorm:"column:creation_by;not null" json:"creation_by"`
	LastUpdatedDate   *time.Time `gorm:"column:last_updated_date" json:"last_updated_date"`
	LastUpdatedBy     string     `gorm:"column:last_updated_by" json:"last_updated_by"`
}

func (LoadingUnloadingPlanDetail) TableName() string {
	return "plan.post_vessel_plan_d"
}

type PostEquipmentPlan struct {
	ID              int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	BranchCode      int        `gorm:"column:branch_code;not null" json:"branch_code"`
	TerminalCode    int        `gorm:"column:terminal_code;not null" json:"terminal_code"`
	PlanNumber      string     `gorm:"column:plan_code;not null" json:"plan_code"`
	SequenceNo      int        `gorm:"column:sequence_no;not null" json:"sequence_no"`
	EquipmentCode   string     `gorm:"column:equipment_code;not null" json:"equipment_code"`
	EquipmentName   string     `gorm:"column:equipment_name" json:"equipment_name"`
	UnitCode        string     `gorm:"column:unit_code" json:"unit_code"`
	PBMCode         string     `gorm:"column:pbm_code" json:"pbm_code"`
	PBMName         string     `gorm:"column:pbm_name" json:"pbm_name"`
	ConsigneeCode   string     `gorm:"column:consignee_code" json:"consignee_code"`
	ConsigneeName   string     `gorm:"column:consignee_name" json:"consignee_name"`
	Description     string     `gorm:"column:description" json:"description"`
	EquipmentGroup  string     `gorm:"column:equipment_group" json:"equipment_group"`
	UnitTon         *float64   `gorm:"column:unit_ton" json:"unit_ton"`
	Attr1           string     `gorm:"column:attr1" json:"attr1"`
	Attr2           string     `gorm:"column:attr2" json:"attr2"`
	Attr3           string     `gorm:"column:attr3" json:"attr3"`
	Value1          *float64   `gorm:"column:value1" json:"value1"`
	Value2          *float64   `gorm:"column:value2" json:"value2"`
	Value3          *float64   `gorm:"column:value3" json:"value3"`
	HeaderID        int64      `gorm:"column:header_id" json:"header_id"`
	DayNo           *int       `gorm:"column:day_no" json:"day_no"`
	ActivityDate    *time.Time `gorm:"column:activity_date" json:"activity_date"`
	Stowage         string     `gorm:"column:stowage" json:"stowage"`
	Quantity        *float64   `gorm:"column:quantity" json:"quantity"`
	CreationDate    time.Time  `gorm:"column:creation_date;not null" json:"creation_date"`
	CreationBy      string     `gorm:"column:creation_by;not null" json:"creation_by"`
	LastUpdatedDate *time.Time `gorm:"column:last_updated_date" json:"last_updated_date"`
	LastUpdatedBy   string     `gorm:"column:last_updated_by" json:"last_updated_by"`
	ProgramName     string     `gorm:"column:program_name;not null" json:"program_name"`
}

func (PostEquipmentPlan) TableName() string {
	return "plan.post_vessel_equipment_plan"
}

type LoadingUnloadingDetermination struct {
	ID                int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	BranchCode        int        `gorm:"column:branch_code;not null" json:"branch_code"`
	TerminalCode      int        `gorm:"column:terminal_code;not null" json:"terminal_code"`
	BranchName        string     `gorm:"column:branch_name" json:"branch_name"`
	TerminalName      string     `gorm:"column:terminal_name" json:"terminal_name"`
	VesselCode        string     `gorm:"column:vessel_code" json:"vessel_code"`
	VesselName        string     `gorm:"column:vessel_name" json:"vessel_name"`
	VesselType        string     `gorm:"column:vessel_type" json:"vessel_type"`
	GRT               string     `gorm:"column:grt" json:"grt"`
	LOA               string     `gorm:"column:loa" json:"loa"`
	VoyageType        string     `gorm:"column:voyage_type" json:"voyage_type"`
	AgentName         string     `gorm:"column:agent_name" json:"agent_name"`
	PPKNumber         string     `gorm:"column:ppk_number" json:"ppk_number"`
	RequestCode       string     `gorm:"column:request_code" json:"request_code"`
	PlanCode          string     `gorm:"column:plan_code;not null" json:"plan_code"`
	PlanDate          time.Time  `gorm:"column:plan_date;not null" json:"plan_date"`
	DeterminationCode string     `gorm:"column:confirmed_plan_code;not null" json:"confirmed_plan_code"`
	DeterminationDate time.Time  `gorm:"column:confirmed_plan_date;not null" json:"confirmed_plan_date"`
	ETA               *time.Time `gorm:"column:eta" json:"eta"`
	ETD               *time.Time `gorm:"column:etd" json:"etd"`
	TGH               *int       `gorm:"column:tgh" json:"tgh"`
	TSD               *int       `gorm:"column:tsd" json:"tsd"`
	PBMCode           string     `gorm:"column:pbm_code" json:"pbm_code"`
	PBMName           string     `gorm:"column:pbm_name" json:"pbm_name"`
	ActivityCode      string     `gorm:"column:activity_code" json:"activity_code"`
	ActivityName      string     `gorm:"column:activity_name" json:"activity_name"`
	Remarks           string     `gorm:"column:remarks" json:"remarks"`
	Status            string     `gorm:"column:status" json:"status"`
	ProgramName       string     `gorm:"column:program_name;not null" json:"program_name"`
	Cycle             string     `gorm:"column:cycle" json:"cycle"`
	ActivityStatus    *int       `gorm:"column:activity_status" json:"activity_status"`
	TruckSequence     int        `gorm:"column:truck_sequence" json:"truck_sequence"`
	CreationDate      time.Time  `gorm:"column:creation_date;not null" json:"creation_date"`
	CreationBy        string     `gorm:"column:creation_by;not null" json:"creation_by"`
	LastUpdatedDate   *time.Time `gorm:"column:last_updated_date" json:"last_updated_date"`
	LastUpdatedBy     string     `gorm:"column:last_updated_by" json:"last_updated_by"`
}

func (LoadingUnloadingDetermination) TableName() string {
	return "plan.post_vessel_confirmed_plan"
}

type LoadingUnloadingDeterminationDetail struct {
	ID                int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	BranchCode        int        `gorm:"column:branch_code;not null" json:"branch_code"`
	TerminalCode      int        `gorm:"column:terminal_code;not null" json:"terminal_code"`
	DeterminationCode string     `gorm:"column:confirmed_plan_code;not null" json:"confirmed_plan_code"`
	RequestCode       string     `gorm:"column:request_code;not null" json:"request_code"`
	WorkOrderCode     string     `gorm:"column:work_order_code;not null" json:"work_order_code"`
	SequenceNo        int        `gorm:"column:sequence_no;not null" json:"sequence_no"`
	ActivityDate      *time.Time `gorm:"column:activity_date" json:"activity_date"`
	Stowage           string     `gorm:"column:stowage" json:"stowage"`
	CargoCode         string     `gorm:"column:cargo_code" json:"cargo_code"`
	CargoName         string     `gorm:"column:cargo_name" json:"cargo_name"`
	TotalQuantity     *float64   `gorm:"column:total_quantity" json:"total_quantity"`
	CargoUnit         string     `gorm:"column:cargo_unit" json:"cargo_unit"`
	CargoPackaging    string     `gorm:"column:cargo_packaging" json:"cargo_packaging"`
	DayNo             *int       `gorm:"column:day_no" json:"day_no"`
	DockCode          string     `gorm:"column:dock_code" json:"dock_code"`
	DockName          string     `gorm:"column:dock_name" json:"dock_name"`
	BerthCode         string     `gorm:"column:berth_code" json:"berth_code"`
	BerthName         string     `gorm:"column:berth_name" json:"berth_name"`
	Shift1            string     `gorm:"column:shift_1" json:"shift_1"`
	Shift2            string     `gorm:"column:shift_2" json:"shift_2"`
	Shift3            string     `gorm:"column:shift_3" json:"shift_3"`
	PBMCode           string     `gorm:"column:pbm_code" json:"pbm_code"`
	PBMName           string     `gorm:"column:pbm_name" json:"pbm_name"`
	ConsigneeCode     string     `gorm:"column:consignee_code" json:"consignee_code"`
	ConsigneeName     string     `gorm:"column:consignee_name" json:"consignee_name"`
	TruckCount        *int       `gorm:"column:truck_count" json:"truck_count"`
	TruckCapacity     *int       `gorm:"column:truck_capacity" json:"truck_capacity"`
	GangCount         *int       `gorm:"column:gang_count" json:"gang_count"`
	Attribute1        string     `gorm:"column:attribute_1" json:"attribute_1"`
	Attribute2        string     `gorm:"column:attribute_2" json:"attribute_2"`
	Attribute3        string     `gorm:"column:attribute_3" json:"attribute_3"`
	Value1            *float64   `gorm:"column:value_1" json:"value_1"`
	Value2            *float64   `gorm:"column:value_2" json:"value_2"`
	Value3            *float64   `gorm:"column:value_3" json:"value_3"`
	Status            string     `gorm:"column:status" json:"status"`
	CreationDate      time.Time  `gorm:"column:creation_date;not null" json:"creation_date"`
	CreationBy        string     `gorm:"column:creation_by;not null" json:"creation_by"`
	ProgramName       string     `gorm:"column:program_name;not null" json:"program_name"`
	CargoNature       string     `gorm:"column:cargo_nature" json:"cargo_nature"`
	CargoNatureDesc   string     `gorm:"column:cargo_nature_desc" json:"cargo_nature_desc"`
	RequestDetailID   *string    `gorm:"column:request_detail_id" json:"request_detail_id"`
}

func (LoadingUnloadingDeterminationDetail) TableName() string {
	return "plan.post_vessel_confirmed_plan_d"
}

type PostEquipmentDetermination struct {
	ID                int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	BranchCode        int        `gorm:"column:branch_code;not null" json:"branch_code"`
	TerminalCode      int        `gorm:"column:terminal_code;not null" json:"terminal_code"`
	RequestCode       string     `gorm:"column:request_code;not null" json:"request_code"`
	DeterminationCode string     `gorm:"column:confirmed_plan_code;not null" json:"confirmed_plan_code"`
	SequenceNo        int        `gorm:"column:sequence_no;not null" json:"sequence_no"`
	EquipmentCode     string     `gorm:"column:equipment_code;not null" json:"equipment_code"`
	EquipmentName     string     `gorm:"column:equipment_name" json:"equipment_name"`
	UnitCode          string     `gorm:"column:unit_code" json:"unit_code"`
	PBMCode           string     `gorm:"column:pbm_code" json:"pbm_code"`
	PBMName           string     `gorm:"column:pbm_name" json:"pbm_name"`
	ConsigneeCode     string     `gorm:"column:consignee_code" json:"consignee_code"`
	ConsigneeName     string     `gorm:"column:consignee_name" json:"consignee_name"`
	Remarks           string     `gorm:"column:remarks" json:"remarks"`
	EquipmentGroup    string     `gorm:"column:equipment_group" json:"equipment_group"`
	UnitTon           *float64   `gorm:"column:unit_ton" json:"unit_ton"`
	Attribute1        string     `gorm:"column:attribute_1" json:"attribute_1"`
	Attribute2        string     `gorm:"column:attribute_2" json:"attribute_2"`
	Attribute3        string     `gorm:"column:attribute_3" json:"attribute_3"`
	Value1            *float64   `gorm:"column:value_1" json:"value_1"`
	Value2            *float64   `gorm:"column:value_2" json:"value_2"`
	Value3            *float64   `gorm:"column:value_3" json:"value_3"`
	CreationDate      time.Time  `gorm:"column:creation_date;not null" json:"creation_date"`
	CreationBy        string     `gorm:"column:creation_by;not null" json:"creation_by"`
	ProgramName       string     `gorm:"column:program_name;not null" json:"program_name"`
	RequestDetailID   string     `gorm:"column:request_detail_id" json:"request_detail_id"`
	DayNo             *int       `gorm:"column:day_no" json:"day_no"`
	ActivityDate      *time.Time `gorm:"column:activity_date" json:"activity_date"`
	Stowage           string     `gorm:"column:stowage" json:"stowage"`
}

func (PostEquipmentDetermination) TableName() string {
	return "plan.post_vessel_equipment_confirmed_plan"
}
