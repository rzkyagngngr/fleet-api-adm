package op

import (
	"omniport-api/internal/helper"
	"time"
)

type SearchReadyOpsPlanInput struct {
	Page     int                  `json:"page"`
	Limit    int                  `json:"limit"`
	Search   string               `json:"search"`
	Filters  map[string]string    `json:"filters"`
	Sort     helper.SortQuery     `json:"sort"`
	Download helper.DownloadQuery `json:"download"`
}

func (r SearchReadyOpsPlanInput) ToPaginationQuery() helper.PaginationQuery {
	return helper.PaginationQuery{
		Page:     r.Page,
		Limit:    r.Limit,
		Search:   r.Search,
		Filters:  r.Filters,
		Sort:     r.Sort,
		Download: r.Download,
	}
}

type GetDataRequestInput struct {
	PKKNumber    string `json:"pkk_number"`
	PPKNumber    string `json:"ppk_number"`
	ActivityCode string `json:"activity_code" binding:"required"`
}

func (r GetDataRequestInput) RequestNumber() string {
	if r.PKKNumber != "" {
		return r.PKKNumber
	}
	return r.PPKNumber
}

type GetDataVesselLookupInput struct {
	PKKNumber  string `json:"pkk_number"`
	PPKNumber  string `json:"ppk_number"`
	VesselCode string `json:"vessel_code"`
}

func (r GetDataVesselLookupInput) RequestNumber() string {
	if r.PKKNumber != "" {
		return r.PKKNumber
	}
	return r.PPKNumber
}

type GetDataVeselInput struct {
	VesselCode string `json:"vessel_code" binding:"required"`
}

type GetDataOpInput struct {
	PPKNumber    string `json:"ppk_number"`
	PlanCode     string `json:"plan_code"`
	PlanNumber   string `json:"plan_number"`
	ActivityCode string `json:"activity_code"`
}

func (r GetDataOpInput) PlanIdentifier() string {
	if r.PlanCode != "" {
		return r.PlanCode
	}
	return r.PlanNumber
}

type GetDetailOpInput struct {
	PlanCode   string `json:"plan_code"`
	PlanNumber string `json:"plan_number"`
}

func (r GetDetailOpInput) PlanIdentifier() string {
	if r.PlanCode != "" {
		return r.PlanCode
	}
	return r.PlanNumber
}

type GetDetailDeterminationInput struct {
	DeterminationCode string `json:"determination_code"`
	WorkOrderCode     string `json:"work_order_code"`
	PlanCode          string `json:"plan_code"`
}

type ReadyOpsPlanResponse struct {
	BranchCode         *int       `json:"branch_code"`
	TerminalCode       *int       `json:"terminal_code"`
	TerminalName       string     `json:"terminal_name"`
	RequestDate        *time.Time `json:"request_date"`
	RequestCode        string     `json:"request_code"`
	PPKNumber          string     `json:"ppk_number"`
	VesselCode         string     `json:"vessel_code"`
	VesselName         string     `json:"vessel_name"`
	AgentName          string     `json:"agent_name"`
	VoyageType         string     `json:"voyage_type"`
	GRT                *float64   `json:"grt"`
	LOA                *float64   `json:"loa"`
	VesselType         string     `json:"vessel_type"`
	PBMCode            string     `json:"pbm_code"`
	PBMName            string     `json:"pbm_name"`
	ActivityName       string     `json:"activity_name"`
	ActivityCode       string     `json:"activity_code"`
	CargoNameList      string     `json:"cargo_name_list"`
	TotalList          string     `json:"total_list"`
	CargoUnitList      string     `json:"cargo_unit_list"`
	Total              *float64   `json:"total"`
}

type ReadyOpDetailResponse struct {
	BranchCode     *int     `json:"branch_code"`
	TerminalCode   *int     `json:"terminal_code"`
	PPKNumber      string   `json:"ppk_number"`
	ActivityCode   string   `json:"activity_code"`
	PBMCode        string   `json:"pbm_code"`
	PBMName        string   `json:"pbm_name"`
	CargoCode      string   `json:"cargo_code"`
	CargoName      string   `json:"cargo_name"`
	CargoUnit      string   `json:"cargo_unit"`
	CargoNature    string   `json:"cargo_nature"`
	CargoPackaging string   `json:"cargo_packaging"`
	StowageCode    string   `json:"stowage_code" gorm:"column:stowage_code"`
	Stowage        string   `json:"stowage"`
	ConsigneeCode  string   `json:"consignee_code"`
	ConsigneeName  string   `json:"consignee_name"`
	Total          *float64 `json:"total"`
}

type GetDataOpResponse struct {
	BranchCode        *int       `json:"branch_code"`
	TerminalCode      *int       `json:"terminal_code"`
	BranchName        string     `json:"branch_name"`
	TerminalName      string     `json:"terminal_name"`
	PlanNumber        string     `json:"plan_code" gorm:"column:plan_code"`
	PlanDate          *time.Time `json:"plan_date"`
	ETA               *time.Time `json:"eta"`
	PPKNumber         string     `json:"ppk_number"`
	ActivityCode      string     `json:"activity_code"`
	ActivityName      string     `json:"activity_name"`
	VesselType        string     `json:"vessel_type"`
	VesselCode        string     `json:"vessel_code"`
	VesselName        string     `json:"vessel_name"`
	GRT               *float64   `json:"grt"`
	LOA               *float64   `json:"loa"`
	ShippingType      string     `json:"shipping_type"`
	BerthName         string     `json:"berth_name"`
	PBMCode           string     `json:"pbm_code"`
	PBMName           string     `json:"pbm_name"`
	DeterminationCode string     `json:"determination_code"`
	WorkOrderCode     string     `json:"work_order_code"`
	Status            *int       `json:"status"`
}

type RawJSONResponse struct {
	Data []byte `json:"data" gorm:"column:data"`
}

func (r RawJSONResponse) MarshalJSON() ([]byte, error) {
	if len(r.Data) == 0 {
		return []byte("null"), nil
	}
	return r.Data, nil
}

type CreateLoadingUnloadingPlanInput struct {
	BranchCode        *int                               `json:"branch_code"`
	TerminalCode      *int                               `json:"terminal_code"`
	BranchName        string                             `json:"branch_name"`
	TerminalName      string                             `json:"terminal_name"`
	VesselCode        string                             `json:"vessel_code"`
	VesselName        string                             `json:"vessel_name"`
	VesselType        string                             `json:"vessel_type"`
	GRT               *float64                           `json:"grt"`
	LOA               *float64                           `json:"loa"`
	ShippingType      string                             `json:"shipping_type"`
	AgentName         string                             `json:"agent_name"`
	PPKNumber         string                             `json:"ppk_number" binding:"required"`
	PlanCode          string                             `json:"plan_code"`
	PlanNumber        string                             `json:"plan_number"`
	PlanDate          time.Time                          `json:"plan_date" binding:"required"`
	ETA               *time.Time                         `json:"eta"`
	ETD               *time.Time                         `json:"etd"`
	BilledTo          *int                               `json:"billed_to"`
	AssignedTo        *int                               `json:"assigned_to"`
	ActivityCode      string                             `json:"activity_code" binding:"required"`
	ActivityName      string                             `json:"activity_name"`
	Remarks           string                             `json:"remarks"`
	Status            *int                               `json:"status"`
	Cycle             string                             `json:"cycle"`
	TotalDays         *int                               `json:"total_days"`
	TotalShifts       *int                               `json:"total_shifts"`
	ActivityStartDate *time.Time                         `json:"activity_start_date"`
	ActivityEndDate   *time.Time                         `json:"activity_end_date"`
	VesselFacing      string                             `json:"vessel_facing"`
	MooringLimit      *float64                           `json:"mooring_limit"`
	BT                *float64                           `json:"bt"`
	Details           []CreateLoadingUnloadingPlanDInput `json:"details"`
	DetailsEquipement []CreatePostEquipmentPlanInput     `json:"detailsEquipement"`
}

type CreateLoadingUnloadingPlanDInput struct {
	SequenceNo      *int       `json:"sequence_no"`
	ActivityDate    *time.Time `json:"activity_date"`
	Stowage         string     `json:"stowage"`
	CargoCode       string     `json:"cargo_code"`
	CargoName       string     `json:"cargo_name"`
	TotalQuantity   *float64   `json:"total_quantity"`
	PlannedQuantity *float64   `json:"planned_quantity"`
	CargoUnit       string     `json:"cargo_unit"`
	CargoPackaging  string     `json:"cargo_packaging"`
	DayNo           *int       `json:"day_no"`
	BerthCode       string     `json:"berth_code"`
	BerthName       string     `json:"berth_name"`
	DockCode        string     `json:"dock_code"`
	DockName        string     `json:"dock_name"`
	Shift1          string     `json:"shift_1"`
	Shift2          string     `json:"shift_2"`
	Shift3          string     `json:"shift_3"`
	PBMCode         string     `json:"pbm_code"`
	PBMName         string     `json:"pbm_name"`
	ConsigneeCode   string     `json:"consignee_code"`
	ConsigneeName   string     `json:"consignee_name"`
	TruckCount      *int       `json:"truck_count"`
	TruckCapacity   *int       `json:"truck_capacity"`
	GangCount       *int       `json:"gang_count"`
	EquipmentCode   string     `json:"equipment_code"`
	EquipmentName   string     `json:"equipment_name"`
	EquipmentGroup  string     `json:"equipment_group"`
	Attrib1         string     `json:"attrib1"`
	Attrib2         string     `json:"attrib2"`
	Attrib3         string     `json:"attrib3"`
	Val1            *float64   `json:"val1"`
	Val2            *float64   `json:"val2"`
	Val3            *float64   `json:"val3"`
	Status          *int       `json:"status"`
	CargoNature     string     `json:"cargo_nature"`
	CargoNatureDesc string     `json:"cargo_nature_desc"`
	PlanDetailCode  string     `json:"plan_detail_code"`
	FromDockCode    string     `json:"from_dock_code"`
	FromBerthCode   string     `json:"from_berth_code"`
}

type CreatePostEquipmentPlanInput struct {
	SequenceNo     *int       `json:"sequence_no"`
	EquipmentCode  string     `json:"equipment_code"`
	EquipmentName  string     `json:"equipment_name"`
	UnitCode       string     `json:"unit_code"`
	PBMCode        string     `json:"pbm_code"`
	PBMName        string     `json:"pbm_name"`
	ConsigneeCode  string     `json:"consignee_code"`
	ConsigneeName  string     `json:"consignee_name"`
	Description    string     `json:"description"`
	EquipmentGroup string     `json:"equipment_group"`
	UnitTon        *float64   `json:"unit_ton"`
	Attr1          string     `json:"attr1"`
	Attr2          string     `json:"attr2"`
	Attr3          string     `json:"attr3"`
	Value1         *float64   `json:"value1"`
	Value2         *float64   `json:"value2"`
	Value3         *float64   `json:"value3"`
	DayNo          *int       `json:"day_no"`
	ActivityDate   *time.Time `json:"activity_date"`
	Stowage        string     `json:"stowage"`
	Quantity       *float64   `json:"quantity"`
}

type CreateLoadingUnloadingDeterminationInput struct {
	BranchCode        *int                                        `json:"branch_code"`
	TerminalCode      *int                                        `json:"terminal_code"`
	BranchName        string                                      `json:"branch_name"`
	TerminalName      string                                      `json:"terminal_name"`
	VesselCode        string                                      `json:"vessel_code"`
	VesselName        string                                      `json:"vessel_name"`
	VesselType        string                                      `json:"vessel_type"`
	GRT               string                                      `json:"grt"`
	LOA               string                                      `json:"loa"`
	VoyageType        string                                      `json:"voyage_type"`
	ShippingType      string                                      `json:"shipping_type"`
	AgentName         string                                      `json:"agent_name"`
	PPKNumber         string                                      `json:"ppk_number"`
	RequestCode       string                                      `json:"request_code"`
	PlanCode          string                                      `json:"plan_code"`
	PlanNumber        string                                      `json:"plan_number"`
	PlanDate          time.Time                                   `json:"plan_date"`
	DeterminationCode string                                      `json:"determination_code"`
	DeterminationDate time.Time                                   `json:"determination_date"`
	ETA               *time.Time                                  `json:"eta"`
	ETD               *time.Time                                  `json:"etd"`
	TGH               *int                                        `json:"tgh"`
	TSD               *int                                        `json:"tsd"`
	PBMCode           string                                      `json:"pbm_code"`
	PBMName           string                                      `json:"pbm_name"`
	ActivityCode      string                                      `json:"activity_code"`
	ActivityName      string                                      `json:"activity_name"`
	Remarks           string                                      `json:"remarks"`
	Status            string                                      `json:"status"`
	Cycle             string                                      `json:"cycle"`
	ActivityStatus    *int                                        `json:"activity_status"`
	TruckSequence     *int                                        `json:"truck_sequence"`
	Details           []CreateLoadingUnloadingDeterminationDInput `json:"details"`
	DetailsEquipement []CreatePostEquipmentDeterminationInput     `json:"detailsEquipement"`
	DetailsEquipment  []CreatePostEquipmentDeterminationInput     `json:"detailsEquipment"`
}

func (r CreateLoadingUnloadingDeterminationInput) PlanIdentifier() string {
	if r.PlanCode != "" {
		return r.PlanCode
	}
	return r.PlanNumber
}

func (r CreateLoadingUnloadingDeterminationInput) EquipmentInputs() []CreatePostEquipmentDeterminationInput {
	if r.DetailsEquipment != nil {
		return r.DetailsEquipment
	}
	return r.DetailsEquipement
}

type CreateLoadingUnloadingDeterminationDInput struct {
	WorkOrderCode   string     `json:"work_order_code"`
	SequenceNo      *int       `json:"sequence_no"`
	ActivityDate    *time.Time `json:"activity_date"`
	Stowage         string     `json:"stowage"`
	CargoCode       string     `json:"cargo_code"`
	CargoName       string     `json:"cargo_name"`
	TotalQuantity   *float64   `json:"total_quantity"`
	CargoUnit       string     `json:"cargo_unit"`
	CargoPackaging  string     `json:"cargo_packaging"`
	DayNo           *int       `json:"day_no"`
	DockCode        string     `json:"dock_code"`
	DockName        string     `json:"dock_name"`
	BerthCode       string     `json:"berth_code"`
	BerthName       string     `json:"berth_name"`
	Shift1          string     `json:"shift_1"`
	Shift2          string     `json:"shift_2"`
	Shift3          string     `json:"shift_3"`
	PBMCode         string     `json:"pbm_code"`
	PBMName         string     `json:"pbm_name"`
	ConsigneeCode   string     `json:"consignee_code"`
	ConsigneeName   string     `json:"consignee_name"`
	TruckCount      *int       `json:"truck_count"`
	TruckCapacity   *int       `json:"truck_capacity"`
	GangCount       *int       `json:"gang_count"`
	Attribute1      string     `json:"attribute_1"`
	Attribute2      string     `json:"attribute_2"`
	Attribute3      string     `json:"attribute_3"`
	Attrib1         string     `json:"attrib1"`
	Attrib2         string     `json:"attrib2"`
	Attrib3         string     `json:"attrib3"`
	Value1          *float64   `json:"value_1"`
	Value2          *float64   `json:"value_2"`
	Value3          *float64   `json:"value_3"`
	Val1            *float64   `json:"val1"`
	Val2            *float64   `json:"val2"`
	Val3            *float64   `json:"val3"`
	Status          string     `json:"status"`
	CargoNature     string     `json:"cargo_nature"`
	CargoNatureDesc string     `json:"cargo_nature_desc"`
	RequestDetailID string     `json:"request_detail_id"`
}

type CreatePostEquipmentDeterminationInput struct {
	SequenceNo      *int       `json:"sequence_no"`
	EquipmentCode   string     `json:"equipment_code"`
	EquipmentName   string     `json:"equipment_name"`
	UnitCode        string     `json:"unit_code"`
	PBMCode         string     `json:"pbm_code"`
	PBMName         string     `json:"pbm_name"`
	ConsigneeCode   string     `json:"consignee_code"`
	ConsigneeName   string     `json:"consignee_name"`
	Remarks         string     `json:"remarks"`
	Description     string     `json:"description"`
	EquipmentGroup  string     `json:"equipment_group"`
	UnitTon         *float64   `json:"unit_ton"`
	Attribute1      string     `json:"attribute_1"`
	Attribute2      string     `json:"attribute_2"`
	Attribute3      string     `json:"attribute_3"`
	Attr1           string     `json:"attr1"`
	Attr2           string     `json:"attr2"`
	Attr3           string     `json:"attr3"`
	Value1          *float64   `json:"value_1"`
	Value2          *float64   `json:"value_2"`
	Value3          *float64   `json:"value_3"`
	Val1            *float64   `json:"val1"`
	Val2            *float64   `json:"val2"`
	Val3            *float64   `json:"val3"`
	RequestDetailID string     `json:"request_detail_id"`
	DayNo           *int       `json:"day_no"`
	ActivityDate    *time.Time `json:"activity_date"`
	Stowage         string     `json:"stowage"`
}

type LoadingUnloadingDeterminationResponse struct {
	Header            *LoadingUnloadingDetermination        `json:"header"`
	Headers           []LoadingUnloadingDetermination       `json:"headers,omitempty"`
	Details           []LoadingUnloadingDeterminationDetail `json:"details"`
	DetailsEquipment  []PostEquipmentDetermination          `json:"detailsEquipment"`
	DetailsEquipement []PostEquipmentDetermination          `json:"detailsEquipement,omitempty"`
}

type UpdateLoadingUnloadingPlanInput struct {
	PlanCode          string                             `json:"plan_code"`
	PlanNumber        string                             `json:"plan_number"`
	TotalDays         *int                               `json:"total_days"`
	TotalShifts       *int                               `json:"total_shifts"`
	ActivityStartDate *time.Time                         `json:"activity_start_date"`
	ActivityEndDate   *time.Time                         `json:"activity_end_date"`
	Details           []CreateLoadingUnloadingPlanDInput `json:"details"`
	DetailsEquipement []CreatePostEquipmentPlanInput     `json:"detailsEquipement"`
}

func (r UpdateLoadingUnloadingPlanInput) PlanIdentifier() string {
	if r.PlanCode != "" {
		return r.PlanCode
	}
	return r.PlanNumber
}

type LoadingUnloadingPlanResponse struct {
	Header            *LoadingUnloadingPlan        `json:"header"`
	Details           []LoadingUnloadingPlanDetail `json:"details"`
	DetailsEquipement []PostEquipmentPlan          `json:"detailsEquipement,omitempty"`
}

type DetailOpResponse struct {
	LoadingUnloadingPlan
	Details           []LoadingUnloadingPlanDetail `json:"details"`
	DetailsEquipment  []PostEquipmentPlan          `json:"detailsEquipment"`
	DetailsEquipement []PostEquipmentPlan          `json:"detailsEquipement,omitempty"`
}

type DetailDeterminationResponse struct {
	LoadingUnloadingDetermination
	Details           []LoadingUnloadingDeterminationDetail `json:"details"`
	DetailsEquipment  []PostEquipmentDetermination          `json:"detailsEquipment"`
	DetailsEquipement []PostEquipmentDetermination          `json:"detailsEquipement,omitempty"`
}

type OpsPlanAuthLocation struct {
	BranchName   string `json:"branch_name"`
	TerminalName string `json:"terminal_name"`
}
