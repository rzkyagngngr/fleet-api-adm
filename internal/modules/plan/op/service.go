package op

import (
	"context"
	"errors"
	"omniport-api/internal/helper"
	"strings"
	"time"
)

const programName = "ADM_SERVICE"

type OpsPlanService interface {
	SearchReady(ctx context.Context, query helper.PaginationQuery) ([]ReadyOpsPlanResponse, helper.PaginationMeta, error)
	GetDataRequest(ctx context.Context, ppkNumber, activityCode string) ([]ReadyOpDetailResponse, error)
	GetDataOp(ctx context.Context, branchCode, terminalCode int, input GetDataOpInput) ([]GetDataOpResponse, error)
	GetDetailOp(ctx context.Context, branchCode, terminalCode int, planNumber string) ([]DetailOpResponse, error)
	GetDataVesselSchedule(ctx context.Context, ppkNumber, vesselCode string) ([]RawJSONResponse, error)
	GetDataVesel(ctx context.Context, vesselCode string) ([]RawJSONResponse, error)
	Create(ctx context.Context, input *CreateLoadingUnloadingPlanInput, branchCode, terminalCode int, branchName, terminalName, createdBy string) (*LoadingUnloadingPlanResponse, error)
	Update(ctx context.Context, input *UpdateLoadingUnloadingPlanInput, branchCode, terminalCode int, updatedBy string) (*LoadingUnloadingPlanResponse, error)
	GetAuthLocation(ctx context.Context, userID uint64) (*OpsPlanAuthLocation, error)
}

type opsPlanService struct {
	repo OpsPlanRepository
}

func NewOpsPlanService(repo OpsPlanRepository) OpsPlanService {
	return &opsPlanService{repo: repo}
}

func (s *opsPlanService) SearchReady(ctx context.Context, query helper.PaginationQuery) ([]ReadyOpsPlanResponse, helper.PaginationMeta, error) {
	return s.repo.SearchReady(ctx, query)
}

func (s *opsPlanService) GetDataRequest(ctx context.Context, ppkNumber, activityCode string) ([]ReadyOpDetailResponse, error) {
	return s.repo.GetDataRequest(ctx, ppkNumber, activityCode)
}

func (s *opsPlanService) GetDataOp(ctx context.Context, branchCode, terminalCode int, input GetDataOpInput) ([]GetDataOpResponse, error) {
	if branchCode == 0 || terminalCode == 0 {
		return nil, errors.New("branch_code and terminal_code are required")
	}
	return s.repo.GetDataOp(ctx, branchCode, terminalCode, input)
}

func (s *opsPlanService) GetDetailOp(ctx context.Context, branchCode, terminalCode int, planNumber string) ([]DetailOpResponse, error) {
	if branchCode == 0 || terminalCode == 0 {
		return nil, errors.New("branch_code and terminal_code are required")
	}
	planNumber = strings.TrimSpace(planNumber)
	if planNumber == "" {
		return nil, errors.New("plan_number is required")
	}

	header, details, detailsEquipment, err := s.repo.GetDetailOp(ctx, branchCode, terminalCode, planNumber)
	if err != nil {
		return nil, err
	}

	return []DetailOpResponse{
		{
			LoadingUnloadingPlan: *header,
			Details:              details,
			DetailsEquipment:     detailsEquipment,
			DetailsEquipement:    detailsEquipment,
		},
	}, nil
}

func (s *opsPlanService) GetDataVesselSchedule(ctx context.Context, ppkNumber, vesselCode string) ([]RawJSONResponse, error) {
	return s.repo.GetDataVesselSchedule(ctx, ppkNumber, vesselCode)
}

func (s *opsPlanService) GetDataVesel(ctx context.Context, vesselCode string) ([]RawJSONResponse, error) {
	return s.repo.GetDataVesel(ctx, vesselCode)
}

func (s *opsPlanService) GetAuthLocation(ctx context.Context, userID uint64) (*OpsPlanAuthLocation, error) {
	return s.repo.GetAuthLocation(ctx, userID)
}

func (s *opsPlanService) Create(
	ctx context.Context,
	input *CreateLoadingUnloadingPlanInput,
	branchCode, terminalCode int,
	branchName, terminalName, createdBy string,
) (*LoadingUnloadingPlanResponse, error) {
	if input == nil {
		return nil, errors.New("payload is required")
	}
	if branchCode == 0 || terminalCode == 0 {
		return nil, errors.New("branch_code and terminal_code are required")
	}
	if branchName == "" {
		branchName = input.BranchName
	}
	if terminalName == "" {
		terminalName = input.TerminalName
	}
	if createdBy == "" {
		createdBy = "SYSTEM"
	}

	input.PPKNumber = strings.TrimSpace(input.PPKNumber)
	input.ActivityCode = strings.TrimSpace(input.ActivityCode)
	if input.PPKNumber == "" {
		return nil, errors.New("ppk_number is required")
	}
	if input.ActivityCode == "" {
		return nil, errors.New("activity_code is required")
	}
	if input.ActivityCode != "BONGKAR" && input.ActivityCode != "MUAT" {
		return nil, errors.New("activity_code must be BONGKAR or MUAT")
	}
	if input.PlanDate.IsZero() {
		return nil, errors.New("plan_date is required")
	}
	if len(input.Details) == 0 {
		return nil, errors.New("at least one detail is required")
	}

	now := time.Now()
	header := &LoadingUnloadingPlan{
		BranchCode:        branchCode,
		TerminalCode:      terminalCode,
		BranchName:        branchName,
		TerminalName:      terminalName,
		VesselCode:        input.VesselCode,
		VesselName:        input.VesselName,
		VesselType:        input.VesselType,
		GRT:               input.GRT,
		LOA:               input.LOA,
		ShippingType:      input.ShippingType,
		AgentName:         input.AgentName,
		PPKNumber:         input.PPKNumber,
		PlanDate:          input.PlanDate,
		ETA:               input.ETA,
		ETD:               input.ETD,
		BilledTo:          input.BilledTo,
		AssignedTo:        input.AssignedTo,
		ActivityCode:      input.ActivityCode,
		ActivityName:      input.ActivityName,
		Remarks:           input.Remarks,
		Status:            input.Status,
		Cycle:             input.Cycle,
		TotalDays:         input.TotalDays,
		TotalShifts:       input.TotalShifts,
		ActivityStartDate: input.ActivityStartDate,
		ActivityEndDate:   input.ActivityEndDate,
		VesselFacing:      input.VesselFacing,
		MooringLimit:      input.MooringLimit,
		BT:                input.BT,
		CreationDate:      now,
		CreationBy:        createdBy,
		LastUpdatedDate:   &now,
		LastUpdatedBy:     createdBy,
		ProgramName:       programName,
	}

	details := buildLoadingUnloadingPlanDetails(input.Details, header, createdBy, now)
	detailsEquipement, err := buildPostEquipmentPlans(input.DetailsEquipement, header, createdBy, now)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, header, details, detailsEquipement); err != nil {
		return nil, err
	}

	return &LoadingUnloadingPlanResponse{Header: header, Details: details, DetailsEquipement: detailsEquipement}, nil
}

func (s *opsPlanService) Update(
	ctx context.Context,
	input *UpdateLoadingUnloadingPlanInput,
	branchCode, terminalCode int,
	updatedBy string,
) (*LoadingUnloadingPlanResponse, error) {
	if input == nil {
		return nil, errors.New("payload is required")
	}
	if branchCode == 0 || terminalCode == 0 {
		return nil, errors.New("branch_code and terminal_code are required")
	}
	input.PlanNumber = strings.TrimSpace(input.PlanNumber)
	if input.PlanNumber == "" {
		return nil, errors.New("plan_number is required")
	}
	if updatedBy == "" {
		updatedBy = "SYSTEM"
	}

	now := time.Now()
	replaceDetails := input.Details != nil
	replaceEquipmentPlans := input.DetailsEquipement != nil || replaceDetails
	details := []LoadingUnloadingPlanDetail(nil)
	if replaceDetails {
		details = buildLoadingUnloadingPlanDetails(input.Details, &LoadingUnloadingPlan{
			BranchCode:   branchCode,
			TerminalCode: terminalCode,
			PlanNumber:   input.PlanNumber,
		}, updatedBy, now)
	}
	detailsEquipement := []PostEquipmentPlan(nil)
	if input.DetailsEquipement != nil {
		var err error
		detailsEquipement, err = buildPostEquipmentPlans(input.DetailsEquipement, &LoadingUnloadingPlan{
			BranchCode:   branchCode,
			TerminalCode: terminalCode,
			PlanNumber:   input.PlanNumber,
		}, updatedBy, now)
		if err != nil {
			return nil, err
		}
	}

	header, finalDetails, finalDetailsEquipement, err := s.repo.Update(
		ctx,
		branchCode,
		terminalCode,
		input,
		details,
		replaceDetails,
		detailsEquipement,
		replaceEquipmentPlans,
		updatedBy,
	)
	if err != nil {
		return nil, err
	}

	return &LoadingUnloadingPlanResponse{Header: header, Details: finalDetails, DetailsEquipement: finalDetailsEquipement}, nil
}

func buildLoadingUnloadingPlanDetails(inputs []CreateLoadingUnloadingPlanDInput, header *LoadingUnloadingPlan, createdBy string, now time.Time) []LoadingUnloadingPlanDetail {
	details := make([]LoadingUnloadingPlanDetail, 0, len(inputs))
	for i, input := range inputs {
		sequenceNo := i + 1
		if input.SequenceNo != nil && *input.SequenceNo > 0 {
			sequenceNo = *input.SequenceNo
		}
		details = append(details, LoadingUnloadingPlanDetail{
			BranchCode:      header.BranchCode,
			TerminalCode:    header.TerminalCode,
			PlanNumber:      header.PlanNumber,
			SequenceNo:      sequenceNo,
			ActivityDate:    input.ActivityDate,
			Stowage:         input.Stowage,
			CargoCode:       input.CargoCode,
			CargoName:       input.CargoName,
			TotalQuantity:   input.TotalQuantity,
			PlannedQuantity: input.PlannedQuantity,
			CargoUnit:       input.CargoUnit,
			CargoPackaging:  input.CargoPackaging,
			DayNo:           input.DayNo,
			BerthCode:       input.BerthCode,
			BerthName:       input.BerthName,
			DockCode:        input.DockCode,
			DockName:        input.DockName,
			Shift1:          input.Shift1,
			Shift2:          input.Shift2,
			Shift3:          input.Shift3,
			PBMCode:         input.PBMCode,
			PBMName:         input.PBMName,
			ConsigneeCode:   input.ConsigneeCode,
			ConsigneeName:   input.ConsigneeName,
			TruckCount:      input.TruckCount,
			TruckCapacity:   input.TruckCapacity,
			GangCount:       input.GangCount,
			EquipmentCode:   input.EquipmentCode,
			EquipmentName:   input.EquipmentName,
			EquipmentGroup:  input.EquipmentGroup,
			Attrib1:         input.Attrib1,
			Attrib2:         input.Attrib2,
			Attrib3:         input.Attrib3,
			Val1:            input.Val1,
			Val2:            input.Val2,
			Val3:            input.Val3,
			Status:          input.Status,
			CargoNature:     input.CargoNature,
			CargoNatureDesc: input.CargoNatureDesc,
			PlanDetailCode:  input.PlanDetailCode,
			FromDockCode:    input.FromDockCode,
			FromBerthCode:   input.FromBerthCode,
			ProgramName:     programName,
			CreationDate:    now,
			CreationBy:      createdBy,
			LastUpdatedDate: &now,
			LastUpdatedBy:   createdBy,
		})
	}
	return details
}

func buildPostEquipmentPlans(inputs []CreatePostEquipmentPlanInput, header *LoadingUnloadingPlan, createdBy string, now time.Time) ([]PostEquipmentPlan, error) {
	equipmentPlans := make([]PostEquipmentPlan, 0, len(inputs))
	for i, input := range inputs {
		equipmentCode := strings.TrimSpace(input.EquipmentCode)
		if equipmentCode == "" {
			return nil, errors.New("detailsEquipement.equipment_code is required")
		}

		sequenceNo := i + 1
		if input.SequenceNo != nil && *input.SequenceNo > 0 {
			sequenceNo = *input.SequenceNo
		}

		equipmentPlans = append(equipmentPlans, PostEquipmentPlan{
			BranchCode:      header.BranchCode,
			TerminalCode:    header.TerminalCode,
			PlanNumber:      header.PlanNumber,
			SequenceNo:      sequenceNo,
			EquipmentCode:   equipmentCode,
			EquipmentName:   input.EquipmentName,
			UnitCode:        input.UnitCode,
			PBMCode:         input.PBMCode,
			PBMName:         input.PBMName,
			ConsigneeCode:   input.ConsigneeCode,
			ConsigneeName:   input.ConsigneeName,
			Description:     input.Description,
			EquipmentGroup:  input.EquipmentGroup,
			UnitTon:         input.UnitTon,
			Attr1:           input.Attr1,
			Attr2:           input.Attr2,
			Attr3:           input.Attr3,
			Value1:          input.Value1,
			Value2:          input.Value2,
			Value3:          input.Value3,
			DayNo:           input.DayNo,
			ActivityDate:    input.ActivityDate,
			Stowage:         input.Stowage,
			Quantity:        input.Quantity,
			CreationDate:    now,
			CreationBy:      createdBy,
			LastUpdatedDate: &now,
			LastUpdatedBy:   createdBy,
			ProgramName:     programName,
		})
	}
	return equipmentPlans, nil
}
