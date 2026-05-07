package op

import (
	"context"
	"errors"
	"fmt"
	"omniport-api/internal/helper"
	"strings"
	"time"
)

const programName = "ADM_SERVICE"

type OpsPlanService interface {
	SearchReady(ctx context.Context, query helper.PaginationQuery) ([]ReadyOpsPlanResponse, helper.PaginationMeta, error)
	GetDataRequest(ctx context.Context, ppkNumber, activityCode string) ([]ReadyOpDetailResponse, error)
	GetDataOp(ctx context.Context, branchCode, terminalCode int, input GetDataOpInput) ([]GetDataOpResponse, error)
	GetDetailOp(ctx context.Context, branchCode, terminalCode int, planCode string) ([]DetailOpResponse, error)
	GetDetailDetermination(ctx context.Context, branchCode, terminalCode int, input GetDetailDeterminationInput) ([]DetailDeterminationResponse, error)
	GetDataVesselSchedule(ctx context.Context, ppkNumber, vesselCode string) ([]RawJSONResponse, error)
	GetDataVesel(ctx context.Context, vesselCode string) ([]RawJSONResponse, error)
	Create(ctx context.Context, input *CreateLoadingUnloadingPlanInput, branchCode, terminalCode int, branchName, terminalName, createdBy string) (*LoadingUnloadingPlanResponse, error)
	CreateDetermination(ctx context.Context, input *CreateLoadingUnloadingDeterminationInput, branchCode, terminalCode int, branchName, terminalName, createdBy string) (*LoadingUnloadingDeterminationResponse, error)
	Update(ctx context.Context, input *UpdateLoadingUnloadingPlanInput, branchCode, terminalCode int, updatedBy string) (*LoadingUnloadingPlanResponse, error)
	UpdateDeterminedPlan(ctx context.Context, input *UpdateLoadingUnloadingPlanInput, branchCode, terminalCode int, updatedBy string) (*LoadingUnloadingPlanResponse, error)
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

func (s *opsPlanService) GetDetailOp(ctx context.Context, branchCode, terminalCode int, planCode string) ([]DetailOpResponse, error) {
	if branchCode == 0 || terminalCode == 0 {
		return nil, errors.New("branch_code and terminal_code are required")
	}
	planCode = strings.TrimSpace(planCode)
	if planCode == "" {
		return nil, errors.New("plan_code is required")
	}

	header, details, detailsEquipment, err := s.repo.GetDetailOp(ctx, branchCode, terminalCode, planCode)
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

func (s *opsPlanService) GetDetailDetermination(ctx context.Context, branchCode, terminalCode int, input GetDetailDeterminationInput) ([]DetailDeterminationResponse, error) {
	if branchCode == 0 || terminalCode == 0 {
		return nil, errors.New("branch_code and terminal_code are required")
	}
	input.DeterminationCode = strings.TrimSpace(input.DeterminationCode)
	if input.DeterminationCode == "" {
		return nil, errors.New("determination_code is required")
	}

	headers, details, detailsEquipment, err := s.repo.GetDetailDetermination(ctx, branchCode, terminalCode, input)
	if err != nil {
		return nil, err
	}

	detailsByCode := make(map[string][]LoadingUnloadingDeterminationDetail, len(headers))
	for _, detail := range details {
		detailsByCode[detail.DeterminationCode] = append(detailsByCode[detail.DeterminationCode], detail)
	}
	equipmentByCode := make(map[string][]PostEquipmentDetermination, len(headers))
	for _, equipment := range detailsEquipment {
		equipmentByCode[equipment.DeterminationCode] = append(equipmentByCode[equipment.DeterminationCode], equipment)
	}

	responses := make([]DetailDeterminationResponse, 0, len(headers))
	for _, header := range headers {
		equipment := equipmentByCode[header.DeterminationCode]
		responses = append(responses, DetailDeterminationResponse{
			LoadingUnloadingDetermination: header,
			Details:                       detailsByCode[header.DeterminationCode],
			DetailsEquipment:              equipment,
			DetailsEquipement:             equipment,
		})
	}

	return responses, nil
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

func (s *opsPlanService) CreateDetermination(
	ctx context.Context,
	input *CreateLoadingUnloadingDeterminationInput,
	branchCode, terminalCode int,
	branchName, terminalName, createdBy string,
) (*LoadingUnloadingDeterminationResponse, error) {
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
	input.RequestCode = strings.TrimSpace(input.RequestCode)
	planCode := strings.TrimSpace(input.PlanIdentifier())
	input.ActivityCode = strings.TrimSpace(input.ActivityCode)
	if planCode == "" {
		return nil, errors.New("plan_code is required")
	}
	if len(input.Details) == 0 && (input.PPKNumber == "" || input.RequestCode == "" || input.ActivityCode == "" || input.PlanDate.IsZero()) {
		builds, err := s.repo.BuildDeterminationsFromPlan(ctx, branchCode, terminalCode, planCode, createdBy)
		if err != nil {
			return nil, err
		}
		for i := range builds {
			header := builds[i].Header
			if strings.TrimSpace(input.DeterminationCode) != "" && len(builds) == 1 {
				header.DeterminationCode = strings.TrimSpace(input.DeterminationCode)
			}
			if !input.DeterminationDate.IsZero() {
				header.DeterminationDate = input.DeterminationDate
			}
			if strings.TrimSpace(input.Remarks) != "" {
				header.Remarks = input.Remarks
			}
			if strings.TrimSpace(input.Status) != "" {
				header.Status = input.Status
			}
			if input.ActivityStatus != nil {
				header.ActivityStatus = input.ActivityStatus
			}
			if input.TruckSequence != nil {
				header.TruckSequence = *input.TruckSequence
			}
		}
		if err := s.repo.CreateDeterminations(ctx, builds); err != nil {
			return nil, err
		}

		headers := make([]LoadingUnloadingDetermination, 0, len(builds))
		details := make([]LoadingUnloadingDeterminationDetail, 0)
		equipmentDeterminations := make([]PostEquipmentDetermination, 0)
		for _, build := range builds {
			headers = append(headers, *build.Header)
			details = append(details, build.Details...)
			equipmentDeterminations = append(equipmentDeterminations, build.EquipmentDeterminations...)
		}

		return &LoadingUnloadingDeterminationResponse{
			Header:            builds[0].Header,
			Headers:           headers,
			Details:           details,
			DetailsEquipment:  equipmentDeterminations,
			DetailsEquipement: equipmentDeterminations,
		}, nil
	}
	if input.PPKNumber == "" {
		return nil, errors.New("ppk_number is required")
	}
	if input.RequestCode == "" {
		return nil, errors.New("request_code is required")
	}
	if input.ActivityCode == "" {
		return nil, errors.New("activity_code is required")
	}
	if input.PlanDate.IsZero() {
		return nil, errors.New("plan_date is required")
	}
	if input.DeterminationDate.IsZero() {
		return nil, errors.New("determination_date is required")
	}
	if len(input.Details) == 0 {
		return nil, errors.New("at least one detail is required")
	}

	now := time.Now()
	truckSequence := 0
	if input.TruckSequence != nil {
		truckSequence = *input.TruckSequence
	}
	header := &LoadingUnloadingDetermination{
		BranchCode:        branchCode,
		TerminalCode:      terminalCode,
		BranchName:        branchName,
		TerminalName:      terminalName,
		VesselCode:        input.VesselCode,
		VesselName:        input.VesselName,
		VesselType:        input.VesselType,
		GRT:               input.GRT,
		LOA:               input.LOA,
		VoyageType:        firstNonEmpty(input.VoyageType, input.ShippingType),
		AgentName:         input.AgentName,
		PPKNumber:         input.PPKNumber,
		RequestCode:       input.RequestCode,
		PlanCode:          planCode,
		PlanDate:          input.PlanDate,
		DeterminationCode: strings.TrimSpace(input.DeterminationCode),
		DeterminationDate: input.DeterminationDate,
		ETA:               input.ETA,
		ETD:               input.ETD,
		TGH:               input.TGH,
		TSD:               input.TSD,
		PBMCode:           input.PBMCode,
		PBMName:           input.PBMName,
		ActivityCode:      input.ActivityCode,
		ActivityName:      input.ActivityName,
		Remarks:           input.Remarks,
		Status:            input.Status,
		ProgramName:       programName,
		Cycle:             input.Cycle,
		ActivityStatus:    input.ActivityStatus,
		TruckSequence:     truckSequence,
		CreationDate:      now,
		CreationBy:        createdBy,
		LastUpdatedDate:   &now,
		LastUpdatedBy:     createdBy,
	}

	details, err := buildLoadingUnloadingDeterminationDetails(input.Details, header, createdBy, now)
	if err != nil {
		return nil, err
	}
	equipmentDeterminations, err := buildPostEquipmentDeterminations(input.EquipmentInputs(), header, createdBy, now)
	if err != nil {
		return nil, err
	}
	if err := s.repo.CreateDeterminations(ctx, []determinationBuild{
		{
			Header:                  header,
			Details:                 details,
			EquipmentDeterminations: equipmentDeterminations,
		},
	}); err != nil {
		return nil, err
	}

	return &LoadingUnloadingDeterminationResponse{
		Header:            header,
		Details:           details,
		DetailsEquipment:  equipmentDeterminations,
		DetailsEquipement: equipmentDeterminations,
	}, nil
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
	planCode := strings.TrimSpace(input.PlanIdentifier())
	if planCode == "" {
		return nil, errors.New("plan_code is required")
	}
	input.PlanCode = planCode
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
			PlanNumber:   planCode,
		}, updatedBy, now)
	}
	detailsEquipement := []PostEquipmentPlan(nil)
	if input.DetailsEquipement != nil {
		var err error
		detailsEquipement, err = buildPostEquipmentPlans(input.DetailsEquipement, &LoadingUnloadingPlan{
			BranchCode:   branchCode,
			TerminalCode: terminalCode,
			PlanNumber:   planCode,
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

func (s *opsPlanService) UpdateDeterminedPlan(
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
	planCode := strings.TrimSpace(input.PlanIdentifier())
	if planCode == "" {
		return nil, errors.New("plan_code is required")
	}
	input.PlanCode = planCode
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
			PlanNumber:   planCode,
		}, updatedBy, now)
	}

	detailsEquipement := []PostEquipmentPlan(nil)
	if input.DetailsEquipement != nil {
		var err error
		detailsEquipement, err = buildPostEquipmentPlans(input.DetailsEquipement, &LoadingUnloadingPlan{
			BranchCode:   branchCode,
			TerminalCode: terminalCode,
			PlanNumber:   planCode,
		}, updatedBy, now)
		if err != nil {
			return nil, err
		}
	}

	header, finalDetails, finalDetailsEquipement, err := s.repo.UpdateDeterminedPlan(
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

func buildLoadingUnloadingDeterminationDetails(inputs []CreateLoadingUnloadingDeterminationDInput, header *LoadingUnloadingDetermination, createdBy string, now time.Time) ([]LoadingUnloadingDeterminationDetail, error) {
	details := make([]LoadingUnloadingDeterminationDetail, 0, len(inputs))
	for i, input := range inputs {
		sequenceNo := i + 1
		if input.SequenceNo != nil && *input.SequenceNo > 0 {
			sequenceNo = *input.SequenceNo
		}
		details = append(details, LoadingUnloadingDeterminationDetail{
			BranchCode:        header.BranchCode,
			TerminalCode:      header.TerminalCode,
			DeterminationCode: header.DeterminationCode,
			RequestCode:       header.RequestCode,
			WorkOrderCode:     strings.TrimSpace(input.WorkOrderCode),
			SequenceNo:        sequenceNo,
			ActivityDate:      input.ActivityDate,
			Stowage:           input.Stowage,
			CargoCode:         input.CargoCode,
			CargoName:         input.CargoName,
			TotalQuantity:     input.TotalQuantity,
			CargoUnit:         input.CargoUnit,
			CargoPackaging:    input.CargoPackaging,
			DayNo:             input.DayNo,
			DockCode:          input.DockCode,
			DockName:          input.DockName,
			BerthCode:         input.BerthCode,
			BerthName:         input.BerthName,
			Shift1:            input.Shift1,
			Shift2:            input.Shift2,
			Shift3:            input.Shift3,
			PBMCode:           input.PBMCode,
			PBMName:           input.PBMName,
			ConsigneeCode:     input.ConsigneeCode,
			ConsigneeName:     input.ConsigneeName,
			TruckCount:        input.TruckCount,
			TruckCapacity:     input.TruckCapacity,
			GangCount:         intPtr(0),
			Attribute1:        firstNonEmpty(input.Attribute1, input.Attrib1),
			Attribute2:        firstNonEmpty(input.Attribute2, input.Attrib2),
			Attribute3:        firstNonEmpty(input.Attribute3, input.Attrib3),
			Value1:            firstFloat(input.Value1, input.Val1),
			Value2:            firstFloat(input.Value2, input.Val2),
			Value3:            firstFloat(input.Value3, input.Val3),
			Status:            input.Status,
			CreationDate:      now,
			CreationBy:        createdBy,
			ProgramName:       programName,
			CargoNature:       input.CargoNature,
			CargoNatureDesc:   input.CargoNatureDesc,
			RequestDetailID:   nil,
		})
	}
	return details, nil
}

func buildPostEquipmentDeterminations(inputs []CreatePostEquipmentDeterminationInput, header *LoadingUnloadingDetermination, createdBy string, now time.Time) ([]PostEquipmentDetermination, error) {
	equipmentDeterminations := make([]PostEquipmentDetermination, 0, len(inputs))
	for i, input := range inputs {
		equipmentCode := strings.TrimSpace(input.EquipmentCode)
		if equipmentCode == "" {
			return nil, errors.New("detailsEquipment.equipment_code is required")
		}

		sequenceNo := i + 1
		if input.SequenceNo != nil && *input.SequenceNo > 0 {
			sequenceNo = *input.SequenceNo
		}
		equipmentDeterminations = append(equipmentDeterminations, PostEquipmentDetermination{
			BranchCode:        header.BranchCode,
			TerminalCode:      header.TerminalCode,
			RequestCode:       header.RequestCode,
			DeterminationCode: header.DeterminationCode,
			SequenceNo:        sequenceNo,
			EquipmentCode:     equipmentCode,
			EquipmentName:     input.EquipmentName,
			UnitCode:          input.UnitCode,
			PBMCode:           input.PBMCode,
			PBMName:           input.PBMName,
			ConsigneeCode:     input.ConsigneeCode,
			ConsigneeName:     input.ConsigneeName,
			Remarks:           firstNonEmpty(input.Remarks, input.Description),
			EquipmentGroup:    input.EquipmentGroup,
			UnitTon:           input.UnitTon,
			Attribute1:        firstNonEmpty(input.Attribute1, input.Attr1),
			Attribute2:        firstNonEmpty(input.Attribute2, input.Attr2),
			Attribute3:        firstNonEmpty(input.Attribute3, input.Attr3),
			Value1:            firstFloat(input.Value1, input.Val1),
			Value2:            firstFloat(input.Value2, input.Val2),
			Value3:            firstFloat(input.Value3, input.Val3),
			CreationDate:      now,
			CreationBy:        createdBy,
			ProgramName:       programName,
			RequestDetailID:   input.RequestDetailID,
			DayNo:             input.DayNo,
			ActivityDate:      input.ActivityDate,
			Stowage:           input.Stowage,
		})
	}
	return equipmentDeterminations, nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func firstFloat(values ...*float64) *float64 {
	for _, value := range values {
		if value != nil {
			return value
		}
	}
	return nil
}

func floatPtrToString(value *float64) string {
	if value == nil {
		return ""
	}
	formatted := fmt.Sprintf("%.2f", *value)
	formatted = strings.TrimRight(formatted, "0")
	return strings.TrimRight(formatted, ".")
}

func intPtrToString(value *int) string {
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%d", *value)
}

func int64ToString(value int64) string {
	if value == 0 {
		return ""
	}
	return fmt.Sprintf("%d", value)
}
