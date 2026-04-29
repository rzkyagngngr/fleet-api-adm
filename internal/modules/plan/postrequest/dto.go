package postrequest

import (
	"omniport-api/internal/helper"
	"time"
)

// ─────────────────────────────────────────────────────────────
// REQUEST DTOs
// ─────────────────────────────────────────────────────────────

// PostRequestDetailInput represents one row of the manifest table.
type PostRequestDetailInput struct {
	SequenceNumber      *int       `json:"sequence_number"`
	StackingType        string     `json:"stacking_type"`
	CargoCode           string     `json:"cargo_code"          binding:"required"`
	CargoName           string     `json:"cargo_name"          binding:"required"`
	CargoUnit           string     `json:"cargo_unit"`
	Total               *float64   `json:"total"`
	QuantityMT          *float64   `json:"quantity_mt"`
	Quantity            *float64   `json:"quantity"`
	CargoNature         string     `json:"cargo_nature"`
	CargoNatureDesc     string     `json:"cargo_nature_desc"`
	CargoPackaging      string     `json:"cargo_packaging"`
	Stowage             string     `json:"stowage"`
	PlannedDate         *time.Time `json:"planned_date"`
	WarehouseID         string     `json:"warehouse_id"`
	BLAWBNumber         string     `json:"bl_awb_number"`
	BLAWBDate           *time.Time `json:"bl_awb_date"`
	Description         string     `json:"description"`
	PackageCount        *float64   `json:"package_count"`
	OriginPortCode      string     `json:"origin_port_code"`
	DestinationPortCode string     `json:"destination_port_code"`
	OriginPortName      string     `json:"origin_port_name"`
	DestinationPortName string     `json:"destination_port_name"`
	StorageReference    string     `json:"storage_reference"`
	StorageStackDate    *time.Time `json:"storage_stack_date"`
	WarehouseDetailID   string     `json:"warehouse_detail_id"`
	WarehouseDetailName string     `json:"warehouse_detail_name"`
	WarehouseName       string     `json:"warehouse_name"`
	ConsigneeCode       string     `json:"consignee_code"`
	ConsigneeName       string     `json:"consignee_name"`
}

// CreatePostRequestInput is the main payload for creating a new cargo request.
type CreatePostRequestInput struct {
	PPKNumber     string     `json:"ppk_number"`
	VesselCode    string     `json:"vessel_code"     binding:"required"`
	VesselName    string     `json:"vessel_name"     binding:"required"`
	VesselType    string     `json:"vessel_type"`
	VoyageType    string     `json:"voyage_type"`
	AgentName     string     `json:"agent_name"`
	RequestDate   time.Time  `json:"request_date"    binding:"required"`
	PBMCode       string     `json:"pbm_code"`
	PBMName       string     `json:"pbm_name"`
	NoBC11        string     `json:"no_bc11"`
	DateBC11      *time.Time `json:"date_bc11"`
	Description   string     `json:"description"`
	RefNumber     string     `json:"ref_number"`
	RefDate       *time.Time `json:"ref_date"`
	Ref1          string     `json:"ref1"`
	Ref2          string     `json:"ref2"`
	Val1          *float64   `json:"val1"`
	Val2          *float64   `json:"val2"`
	TotalManifest *float64   `json:"total_manifest"`
	BillableCode  string     `json:"billable_code"`
	BillableName  string     `json:"billable_name"`
	VesselCodeDst string     `json:"vessel_code_dst"`
	VesselNameDst string     `json:"vessel_name_dst"`
	ActivityCode  string     `json:"activity_code"`
	ActivityName  string     `json:"activity_name"`
	ToPPKNumber   string     `json:"to_ppk_number"`

	Details []PostRequestDetailInput `json:"details" binding:"required,min=1,dive"`
}

// UpdatePostRequestInput is the payload for updating an existing request.
// Identical fields, but all optional at the header level.
type UpdatePostRequestInput struct {
	PPKNumber     string     `json:"ppk_number"`
	VesselCode    string     `json:"vessel_code"`
	VesselName    string     `json:"vessel_name"`
	VesselType    string     `json:"vessel_type"`
	VoyageType    string     `json:"voyage_type"`
	AgentName     string     `json:"agent_name"`
	RequestDate   *time.Time `json:"request_date"`
	PBMCode       string     `json:"pbm_code"`
	PBMName       string     `json:"pbm_name"`
	NoBC11        string     `json:"no_bc11"`
	DateBC11      *time.Time `json:"date_bc11"`
	Description   string     `json:"description"`
	RefNumber     string     `json:"ref_number"`
	RefDate       *time.Time `json:"ref_date"`
	Ref1          string     `json:"ref1"`
	Ref2          string     `json:"ref2"`
	Val1          *float64   `json:"val1"`
	Val2          *float64   `json:"val2"`
	TotalManifest *float64   `json:"total_manifest"`
	BillableCode  string     `json:"billable_code"`
	BillableName  string     `json:"billable_name"`
	VesselCodeDst string     `json:"vessel_code_dst"`
	VesselNameDst string     `json:"vessel_name_dst"`
	ActivityCode  string     `json:"activity_code"`
	ActivityName  string     `json:"activity_name"`
	ToPPKNumber   string     `json:"to_ppk_number"`

	// Replaces all details when provided
	Details []PostRequestDetailInput `json:"details"`
}

// SearchPostRequestInput is the body for the search/list endpoint.
type SearchPostRequestInput struct {
	Page     int                  `json:"page"`
	Limit    int                  `json:"limit"`
	Search   string               `json:"search"`
	Filters  map[string]string    `json:"filters"`
	Sort     helper.SortQuery     `json:"sort"`
	Download helper.DownloadQuery `json:"download"`
}

func (r SearchPostRequestInput) ToPaginationQuery() helper.PaginationQuery {
	return helper.PaginationQuery{
		Page:     r.Page,
		Limit:    r.Limit,
		Search:   r.Search,
		Filters:  r.Filters,
		Sort:     r.Sort,
		Download: r.Download,
	}
}

// ─────────────────────────────────────────────────────────────
// RESPONSE DTOs
// ─────────────────────────────────────────────────────────────

// PostRequestDetailResponse is the read-model for one manifest line.
type PostRequestDetailResponse struct {
	ID                  int64      `json:"id"`
	BranchCode          int        `json:"branch_code"`
	TerminalCode        int        `json:"terminal_code"`
	BranchName          string     `json:"branch_name"`
	TerminalName        string     `json:"terminal_name"`
	RequestCode         string     `json:"request_code"`
	SequenceNumber      *int       `json:"sequence_number"`
	StackingType        string     `json:"stacking_type"`
	CargoCode           string     `json:"cargo_code"`
	CargoName           string     `json:"cargo_name"`
	CargoUnit           string     `json:"cargo_unit"`
	Total               *float64   `json:"total"`
	QuantityMT          *float64   `json:"quantity_mt"`
	Quantity            *float64   `json:"quantity"`
	CargoNature         string     `json:"cargo_nature"`
	CargoNatureDesc     string     `json:"cargo_nature_desc"`
	CargoPackaging      string     `json:"cargo_packaging"`
	Stowage             string     `json:"stowage"`
	PlannedDate         *time.Time `json:"planned_date"`
	WarehouseID         string     `json:"warehouse_id"`
	BLAWBNumber         string     `json:"bl_awb_number"`
	BLAWBDate           *time.Time `json:"bl_awb_date"`
	Description         string     `json:"description"`
	PackageCount        *float64   `json:"package_count"`
	OriginPortCode      string     `json:"origin_port_code"`
	DestinationPortCode string     `json:"destination_port_code"`
	OriginPortName      string     `json:"origin_port_name"`
	DestinationPortName string     `json:"destination_port_name"`
	StorageReference    string     `json:"storage_reference"`
	StorageStackDate    *time.Time `json:"storage_stack_date"`
	WarehouseDetailID   string     `json:"warehouse_detail_id"`
	WarehouseDetailName string     `json:"warehouse_detail_name"`
	WarehouseName       string     `json:"warehouse_name"`
	ConsigneeCode       string     `json:"consignee_code"`
	ConsigneeName       string     `json:"consignee_name"`
	CreationDate        *time.Time `json:"creation_date"`
	CreationBy          string     `json:"creation_by"`
	LastUpdatedDate     *time.Time `json:"last_updated_date"`
	LastUpdatedBy       string     `json:"last_updated_by"`
}

// PostRequestResponse is the read-model for a cargo service request header.
type PostRequestResponse struct {
	ID              int64                       `json:"id"`
	BranchCode      *int                        `json:"branch_code"`
	TerminalCode    *int                        `json:"terminal_code"`
	BranchName      string                      `json:"branch_name"`
	TerminalName    string                      `json:"terminal_name"`
	PPKNumber       string                      `json:"ppk_number"`
	VesselCode      string                      `json:"vessel_code"`
	VesselName      string                      `json:"vessel_name"`
	VesselType      string                      `json:"vessel_type"`
	VoyageType      string                      `json:"voyage_type"`
	AgentName       string                      `json:"agent_name"`
	RequestCode     string                      `json:"request_code"`
	RequestDate     time.Time                   `json:"request_date"`
	PBMCode         string                      `json:"pbm_code"`
	PBMName         string                      `json:"pbm_name"`
	NoBC11          string                      `json:"no_bc11"`
	DateBC11        *time.Time                  `json:"date_bc11"`
	Description     string                      `json:"description"`
	Status          *int                        `json:"status"`
	PlanStatus      *int                        `json:"plan_status"`
	RefNumber       string                      `json:"ref_number"`
	RefDate         *time.Time                  `json:"ref_date"`
	Ref1            string                      `json:"ref1"`
	Ref2            string                      `json:"ref2"`
	Val1            *float64                    `json:"val1"`
	Val2            *float64                    `json:"val2"`
	TotalManifest   *float64                    `json:"total_manifest"`
	BillableCode    string                      `json:"billable_code"`
	BillableName    string                      `json:"billable_name"`
	VesselCodeDst   string                      `json:"vessel_code_dst"`
	VesselNameDst   string                      `json:"vessel_name_dst"`
	ActivityCode    string                      `json:"activity_code"`
	ActivityName    string                      `json:"activity_name"`
	ToPPKNumber     string                      `json:"to_ppk_number"`
	ApprovalDate    *time.Time                  `json:"approval_date"`
	CreationDate    *time.Time                  `json:"creation_date"`
	CreationBy      string                      `json:"creation_by"`
	LastUpdatedDate *time.Time                  `json:"last_updated_date"`
	LastUpdatedBy   string                      `json:"last_updated_by"`
	Details         []PostRequestDetailResponse `json:"details"`
}

// PostRequestStatsResponse holds aggregated dashboard counts.
type PostRequestStatsResponse struct {
	Total    int64 `json:"total"`
	Pending  int64 `json:"pending"`
	Approved int64 `json:"approved"`
	Rejected int64 `json:"rejected"`
}

// ─────────────────────────────────────────────────────────────
// MAPPERS
// ─────────────────────────────────────────────────────────────

func detailToResponse(d PostRequestDetail) PostRequestDetailResponse {
	return PostRequestDetailResponse{
		ID:                  d.ID,
		BranchCode:          d.BranchCode,
		TerminalCode:        d.TerminalCode,
		BranchName:          d.BranchName,
		TerminalName:        d.TerminalName,
		RequestCode:         d.RequestCode,
		SequenceNumber:      d.SequenceNumber,
		StackingType:        d.StackingType,
		CargoCode:           d.CargoCode,
		CargoName:           d.CargoName,
		CargoUnit:           d.CargoUnit,
		Total:               d.Total,
		QuantityMT:          d.QuantityMT,
		Quantity:            d.Quantity,
		CargoNature:         d.CargoNature,
		CargoNatureDesc:     d.CargoNatureDesc,
		CargoPackaging:      d.CargoPackaging,
		Stowage:             d.Stowage,
		PlannedDate:         d.PlannedDate,
		WarehouseID:         d.WarehouseID,
		BLAWBNumber:         d.BLAWBNumber,
		BLAWBDate:           d.BLAWBDate,
		Description:         d.Description,
		PackageCount:        d.PackageCount,
		OriginPortCode:      d.OriginPortCode,
		DestinationPortCode: d.DestinationPortCode,
		OriginPortName:      d.OriginPortName,
		DestinationPortName: d.DestinationPortName,
		StorageReference:    d.StorageReference,
		StorageStackDate:    d.StorageStackDate,
		WarehouseDetailID:   d.WarehouseDetailID,
		WarehouseDetailName: d.WarehouseDetailName,
		WarehouseName:       d.WarehouseName,
		ConsigneeCode:       d.ConsigneeCode,
		ConsigneeName:       d.ConsigneeName,
		CreationDate:        d.CreationDate,
		CreationBy:          d.CreationBy,
		LastUpdatedDate:     d.LastUpdatedDate,
		LastUpdatedBy:       d.LastUpdatedBy,
	}
}

func (h *PostRequest) ToResponse(details []PostRequestDetail) PostRequestResponse {
	res := PostRequestResponse{
		ID:              h.ID,
		BranchCode:      h.BranchCode,
		TerminalCode:    h.TerminalCode,
		BranchName:      h.BranchName,
		TerminalName:    h.TerminalName,
		PPKNumber:       h.PPKNumber,
		VesselCode:      h.VesselCode,
		VesselName:      h.VesselName,
		VesselType:      h.VesselType,
		VoyageType:      h.VoyageType,
		AgentName:       h.AgentName,
		RequestCode:     h.RequestCode,
		RequestDate:     h.RequestDate,
		PBMCode:         h.PBMCode,
		PBMName:         h.PBMName,
		NoBC11:          h.NoBC11,
		DateBC11:        h.DateBC11,
		Description:     h.Description,
		Status:          h.Status,
		PlanStatus:      h.PlanStatus,
		RefNumber:       h.RefNumber,
		RefDate:         h.RefDate,
		Ref1:            h.Ref1,
		Ref2:            h.Ref2,
		Val1:            h.Val1,
		Val2:            h.Val2,
		TotalManifest:   h.TotalManifest,
		BillableCode:    h.BillableCode,
		BillableName:    h.BillableName,
		VesselCodeDst:   h.VesselCodeDst,
		VesselNameDst:   h.VesselNameDst,
		ActivityCode:    h.ActivityCode,
		ActivityName:    h.ActivityName,
		ToPPKNumber:     h.ToPPKNumber,
		ApprovalDate:    h.ApprovalDate,
		CreationDate:    h.CreationDate,
		CreationBy:      h.CreationBy,
		LastUpdatedDate: h.LastUpdatedDate,
		LastUpdatedBy:   h.LastUpdatedBy,
	}
	for _, d := range details {
		res.Details = append(res.Details, detailToResponse(d))
	}
	return res
}

// PostVesselScheduleResponse represents the read-model for a vessel schedule.
type PostVesselScheduleResponse struct {
	ID                  int64      `json:"id"`
	BranchCode          *int       `json:"branch_code"`
	TerminalCode        *int       `json:"terminal_code"`
	BranchName          string     `json:"branch_name"`
	TerminalName        string     `json:"terminal_name"`
	ScheduleCode        string     `json:"schedule_code"`
	PKKNumber           string     `json:"pkk_number"`
	VesselName          string     `json:"vessel_name"`
	VesselCode          string     `json:"vessel_code"`
	VesselType          string     `json:"vessel_type"`
	VoyageNumber        string     `json:"voyage_number"`
	VoyageType          string     `json:"voyage_type"`

	GRT                 *int       `json:"grt"`
	LOA                 *float64   `json:"loa"`
	AgencyName          string     `json:"agency_name"`
	PortAgent           string     `json:"port_agent"`
	EmergencyContact    string     `json:"emergency_contact"`
	OriginPortCode      string     `json:"origin_port_code"`
	OriginPortName      string     `json:"origin_port_name"`
	DestinationPortCode string     `json:"destination_port_code"`
	DestinationPortName string     `json:"destination_port_name"`
	DischargePortCode   string     `json:"discharge_port_code"`
	DischargePortName   string     `json:"discharge_port_name"`
	AssignedBerthName   string     `json:"assigned_berth_name"`
	DockID              *int       `json:"dock_id"`
	DockCode            string     `json:"dock_code"`
	DockName            string     `json:"dock_name"`
	BerthCode           string     `json:"berth_code"`
	BerthName           string     `json:"berth_name"`
	BerthPosition       string     `json:"berth_position"`
	PositionRange       string     `json:"position_range"`
	ETA                 *time.Time `json:"eta"`
	ETB                 *time.Time `json:"etb"`
	ETC                 *time.Time `json:"etc"`
	ETD                 *time.Time `json:"etd"`
	Status              *int       `json:"status"`
	CreationDate        *time.Time `json:"creation_date"`
	CreationBy          string     `json:"creation_by"`
	LastUpdatedDate     *time.Time `json:"last_updated_date"`
	LastUpdatedBy       string     `json:"last_updated_by"`
}

func (h *PostVesselSchedule) ToResponse() PostVesselScheduleResponse {
	return PostVesselScheduleResponse{
		ID:                  h.ID,
		BranchCode:          h.BranchCode,
		TerminalCode:        h.TerminalCode,
		BranchName:          h.BranchName,
		TerminalName:        h.TerminalName,
		ScheduleCode:        h.ScheduleCode,
		PKKNumber:           h.PKKNumber,
		VesselName:          h.VesselName,
		VesselCode:          h.VesselCode,
		VesselType:          h.VesselType,
		VoyageNumber:        h.VoyageNumber,
		VoyageType:          h.VoyageType,

		GRT:                 h.GRT,
		LOA:                 h.LOA,
		AgencyName:          h.AgencyName,
		PortAgent:           h.PortAgent,
		EmergencyContact:    h.EmergencyContact,
		OriginPortCode:      h.OriginPortCode,
		OriginPortName:      h.OriginPortName,
		DestinationPortCode: h.DestinationPortCode,
		DestinationPortName: h.DestinationPortName,
		DischargePortCode:   h.DischargePortCode,
		DischargePortName:   h.DischargePortName,
		AssignedBerthName:   h.AssignedBerthName,
		DockID:              h.DockID,
		DockCode:            h.DockCode,
		DockName:            h.DockName,
		BerthCode:           h.BerthCode,
		BerthName:           h.BerthName,
		BerthPosition:       h.BerthPosition,
		PositionRange:       h.PositionRange,
		ETA:                 h.ETA,
		ETB:                 h.ETB,
		ETC:                 h.ETC,
		ETD:                 h.ETD,
		Status:              h.Status,
		CreationDate:        h.CreationDate,
		CreationBy:          h.CreationBy,
		LastUpdatedDate:     h.LastUpdatedDate,
		LastUpdatedBy:       h.LastUpdatedBy,
	}
}

