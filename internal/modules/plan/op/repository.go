package op

import (
	"context"
	"fmt"
	"math"
	"omniport-api/internal/helper"
	"strings"
	"time"

	"gorm.io/gorm"
)

type OpsPlanRepository interface {
	SearchReady(ctx context.Context, param helper.PaginationQuery) ([]ReadyOpsPlanResponse, helper.PaginationMeta, error)
	GetDataRequest(ctx context.Context, ppkNumber, activityCode string) ([]ReadyOpDetailResponse, error)
	GetDataOp(ctx context.Context, branchCode, terminalCode int, input GetDataOpInput) ([]GetDataOpResponse, error)
	GetDetailOp(ctx context.Context, branchCode, terminalCode int, planCode string) (*LoadingUnloadingPlan, []LoadingUnloadingPlanDetail, []PostEquipmentPlan, error)
	GetDetailDetermination(ctx context.Context, branchCode, terminalCode int, input GetDetailDeterminationInput) ([]LoadingUnloadingDetermination, []LoadingUnloadingDeterminationDetail, []PostEquipmentDetermination, error)
	GetDataVesselSchedule(ctx context.Context, ppkNumber, vesselCode string) ([]RawJSONResponse, error)
	GetDataVesel(ctx context.Context, vesselCode string) ([]RawJSONResponse, error)
	Create(ctx context.Context, header *LoadingUnloadingPlan, details []LoadingUnloadingPlanDetail, equipmentPlans []PostEquipmentPlan) error
	BuildDeterminationsFromPlan(ctx context.Context, branchCode, terminalCode int, planCode, createdBy string) ([]determinationBuild, error)
	CreateDeterminations(ctx context.Context, builds []determinationBuild) error
	Update(ctx context.Context, branchCode, terminalCode int, input *UpdateLoadingUnloadingPlanInput, details []LoadingUnloadingPlanDetail, replaceDetails bool, equipmentPlans []PostEquipmentPlan, replaceEquipmentPlans bool, updatedBy string) (*LoadingUnloadingPlan, []LoadingUnloadingPlanDetail, []PostEquipmentPlan, error)
	UpdateDeterminedPlan(ctx context.Context, branchCode, terminalCode int, input *UpdateLoadingUnloadingPlanInput, details []LoadingUnloadingPlanDetail, replaceDetails bool, equipmentPlans []PostEquipmentPlan, replaceEquipmentPlans bool, updatedBy string) (*LoadingUnloadingPlan, []LoadingUnloadingPlanDetail, []PostEquipmentPlan, error)
	GetAuthLocation(ctx context.Context, userID uint64) (*OpsPlanAuthLocation, error)
}

type opsPlanRepository struct {
	db    *gorm.DB
	admDB *gorm.DB
}

type determinationBuild struct {
	Header                  *LoadingUnloadingDetermination
	Details                 []LoadingUnloadingDeterminationDetail
	EquipmentDeterminations []PostEquipmentDetermination
}

func NewOpsPlanRepository(db *gorm.DB, admDB ...*gorm.DB) OpsPlanRepository {
	vesselDB := db
	if len(admDB) > 0 && admDB[0] != nil {
		vesselDB = admDB[0]
	}
	return &opsPlanRepository{db: db, admDB: vesselDB}
}

func (r *opsPlanRepository) SearchReady(ctx context.Context, param helper.PaginationQuery) ([]ReadyOpsPlanResponse, helper.PaginationMeta, error) {
	meta := helper.PaginationMeta{Page: 1, Limit: 10, MaxDownloadLimit: 1000}

	detailTable, err := r.resolvePostRequestDetailTable(ctx)
	if err != nil {
		return nil, meta, err
	}

	page := maxInt(param.Page, 1)
	limit := maxInt(param.Limit, 1)
	if limit > 200 {
		limit = 200
	}

	baseQuery := readyOpsPlanBaseQuery(detailTable)
	whereParts := []string{"1=1"}
	args := make([]interface{}, 0)

	search := strings.TrimSpace(param.Search)
	if search != "" {
		whereParts = append(whereParts, `(
			UPPER(COALESCE(request_code, '')) LIKE UPPER(?) OR
			UPPER(COALESCE(ppk_number, '')) LIKE UPPER(?) OR
			UPPER(COALESCE(vessel_name, '')) LIKE UPPER(?) OR
			UPPER(COALESCE(pbm_name, '')) LIKE UPPER(?) OR
			UPPER(COALESCE(agent_name, '')) LIKE UPPER(?)
		)`)
		like := "%" + search + "%"
		args = append(args, like, like, like, like, like)
	}

	filterable := map[string]string{
		"branch_code":   "branch_code",
		"terminal_code": "terminal_code",
		"request_code":  "request_code",
		"ppk_number":    "ppk_number",
		"vessel_code":   "vessel_code",
		"vessel_name":   "vessel_name",
		"pbm_code":      "pbm_code",
		"activity_code": "activity_code",
	}
	for key, value := range param.Filters {
		column, ok := filterable[key]
		filterValue := strings.TrimSpace(value)
		if !ok || filterValue == "" {
			continue
		}
		whereParts = append(whereParts, fmt.Sprintf("UPPER(CAST(%s AS TEXT)) LIKE UPPER(?)", column))
		args = append(args, "%"+filterValue+"%")
	}

	filterQuery := " WHERE " + strings.Join(whereParts, " AND ")
	countQuery := "WITH ready_ops_plan AS (" + baseQuery + ") SELECT COUNT(1) FROM ready_ops_plan" + filterQuery
	if err := r.db.WithContext(ctx).Raw(countQuery, args...).Scan(&meta.TotalItems).Error; err != nil {
		return nil, meta, err
	}

	sortable := map[string]string{
		"request_date":  "request_date",
		"request_code":  "request_code",
		"ppk_number":    "ppk_number",
		"vessel_name":   "vessel_name",
		"pbm_name":      "pbm_name",
		"activity_code": "activity_code",
		"total":         "total",
	}
	sortColumn := sortable[param.Sort.By]
	if sortColumn == "" {
		sortColumn = "request_date"
	}
	sortOrder := strings.ToUpper(strings.TrimSpace(param.Sort.Order))
	if sortOrder != "ASC" {
		sortOrder = "DESC"
	}

	offset := (page - 1) * limit
	effectiveLimit := limit
	if param.Download.IsDownload {
		switch strings.ToLower(strings.TrimSpace(param.Download.Type)) {
		case "range":
			rangeStart := maxInt(param.Download.RangeStart, 1)
			rangeEnd := maxInt(param.Download.RangeEnd, rangeStart)
			offset = (rangeStart - 1) * limit
			effectiveLimit = (rangeEnd - rangeStart + 1) * limit
		default:
			offset = 0
			effectiveLimit = int(meta.TotalItems)
		}
		if effectiveLimit > meta.MaxDownloadLimit {
			effectiveLimit = meta.MaxDownloadLimit
		}
	}

	dataArgs := append([]interface{}{}, args...)
	dataQuery := fmt.Sprintf(
		"WITH ready_ops_plan AS (%s) SELECT * FROM ready_ops_plan%s ORDER BY %s %s LIMIT ? OFFSET ?",
		baseQuery,
		filterQuery,
		sortColumn,
		sortOrder,
	)
	dataArgs = append(dataArgs, effectiveLimit, offset)

	var rows []ReadyOpsPlanResponse
	if err := r.db.WithContext(ctx).Raw(dataQuery, dataArgs...).Scan(&rows).Error; err != nil {
		return nil, meta, err
	}

	meta.Page = page
	meta.Limit = limit
	meta.TotalPages = int(math.Ceil(float64(meta.TotalItems) / float64(limit)))
	return rows, meta, nil
}

func (r *opsPlanRepository) GetDataRequest(ctx context.Context, ppkNumber, activityCode string) ([]ReadyOpDetailResponse, error) {
	detailTable, err := r.resolvePostRequestDetailTable(ctx)
	if err != nil {
		return nil, err
	}

	var rows []ReadyOpDetailResponse
	if err := r.db.WithContext(ctx).Raw(getDataRequestQuery(detailTable), ppkNumber, activityCode).Scan(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *opsPlanRepository) GetDataOp(ctx context.Context, branchCode, terminalCode int, input GetDataOpInput) ([]GetDataOpResponse, error) {
	whereParts := []string{
		"a.branch_code = ?",
		"a.terminal_code = ?",
	}
	args := []interface{}{branchCode, terminalCode}

	if value := strings.TrimSpace(input.PPKNumber); value != "" {
		whereParts = append(whereParts, "a.ppk_number = ?")
		args = append(args, value)
	}
	if value := strings.TrimSpace(input.PlanIdentifier()); value != "" {
		whereParts = append(whereParts, "a.plan_code = ?")
		args = append(args, value)
	}
	if value := strings.TrimSpace(input.ActivityCode); value != "" {
		whereParts = append(whereParts, "a.activity_code = ?")
		args = append(args, value)
	}

	var rows []GetDataOpResponse
	if err := r.db.WithContext(ctx).Raw(getDataOpQuery(strings.Join(whereParts, " AND ")), args...).Scan(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *opsPlanRepository) GetDetailOp(ctx context.Context, branchCode, terminalCode int, planCode string) (*LoadingUnloadingPlan, []LoadingUnloadingPlanDetail, []PostEquipmentPlan, error) {
	var header LoadingUnloadingPlan
	if err := r.db.WithContext(ctx).
		Table("plan.post_vessel_plan a").
		Select("a.*, COALESCE(rpk.id, 0) AS vessel_rpk_id").
		Joins("LEFT JOIN plan.post_vessel_rpk rpk ON a.plan_code = rpk.ops_plan_code AND a.branch_code = rpk.branch_code AND a.terminal_code = rpk.terminal_code").
		Where("a.branch_code = ? AND a.terminal_code = ? AND a.plan_code = ?", branchCode, terminalCode, planCode).
		First(&header).Error; err != nil {
		return nil, nil, nil, err
	}

	var details []LoadingUnloadingPlanDetail
	if err := r.db.WithContext(ctx).
		Where("branch_code = ? AND terminal_code = ? AND plan_code = ?", branchCode, terminalCode, planCode).
		Order("sequence_no ASC").
		Find(&details).Error; err != nil {
		return nil, nil, nil, err
	}

	var detailsEquipment []PostEquipmentPlan
	if err := r.db.WithContext(ctx).
		Where("branch_code = ? AND terminal_code = ? AND plan_code = ?", branchCode, terminalCode, planCode).
		Order("sequence_no ASC").
		Find(&detailsEquipment).Error; err != nil {
		return nil, nil, nil, err
	}

	return &header, details, detailsEquipment, nil
}

func (r *opsPlanRepository) GetDetailDetermination(ctx context.Context, branchCode, terminalCode int, input GetDetailDeterminationInput) ([]LoadingUnloadingDetermination, []LoadingUnloadingDeterminationDetail, []PostEquipmentDetermination, error) {
	determinationCode := strings.TrimSpace(input.DeterminationCode)
	workOrderCode := strings.TrimSpace(input.WorkOrderCode)
	if determinationCode == "" && workOrderCode == "" {
		return nil, nil, nil, fmt.Errorf("confirmed_plan_code or work_order_code is required")
	}

	var headers []LoadingUnloadingDetermination
	headerQuery := r.db.WithContext(ctx).
		Where("branch_code = ? AND terminal_code = ?", branchCode, terminalCode)
	if determinationCode != "" {
		headerQuery = headerQuery.Where("confirmed_plan_code = ?", determinationCode)
	}
	if workOrderCode != "" {
		headerQuery = headerQuery.Where(`
			confirmed_plan_code IN (
				SELECT confirmed_plan_code
				FROM plan.post_vessel_confirmed_plan_d
				WHERE branch_code = ?
					AND terminal_code = ?
					AND work_order_code = ?
			)
		`, branchCode, terminalCode, workOrderCode)
	}
	if err := headerQuery.Order("confirmed_plan_code ASC").Find(&headers).Error; err != nil {
		return nil, nil, nil, err
	}
	if len(headers) == 0 {
		return nil, nil, nil, gorm.ErrRecordNotFound
	}

	determinationCodes := make([]string, 0, len(headers))
	for _, header := range headers {
		determinationCodes = append(determinationCodes, header.DeterminationCode)
	}

	var details []LoadingUnloadingDeterminationDetail
	if err := r.db.WithContext(ctx).
		Where("branch_code = ? AND terminal_code = ? AND confirmed_plan_code IN ?", branchCode, terminalCode, determinationCodes).
		Order("confirmed_plan_code ASC, sequence_no ASC").
		Find(&details).Error; err != nil {
		return nil, nil, nil, err
	}

	var detailsEquipment []PostEquipmentDetermination
	if err := r.db.WithContext(ctx).
		Where("branch_code = ? AND terminal_code = ? AND confirmed_plan_code IN ?", branchCode, terminalCode, determinationCodes).
		Order("confirmed_plan_code ASC, sequence_no ASC").
		Find(&detailsEquipment).Error; err != nil {
		return nil, nil, nil, err
	}

	return headers, details, detailsEquipment, nil
}

type determinationRequestHeader struct {
	BranchName   string `gorm:"column:branch_name"`
	TerminalName string `gorm:"column:terminal_name"`
	RequestCode  string `gorm:"column:request_code"`
	VesselCode   string `gorm:"column:vessel_code"`
	VesselName   string `gorm:"column:vessel_name"`
	VesselType   string `gorm:"column:vessel_type"`
	VoyageType   string `gorm:"column:voyage_type"`
	AgentName    string `gorm:"column:agent_name"`
	PBMCode      string `gorm:"column:pbm_code"`
	PBMName      string `gorm:"column:pbm_name"`
	ActivityName string `gorm:"column:activity_name"`
}

type determinationRequestDetail struct {
	ID              int64    `gorm:"column:id"`
	SequenceNumber  *int     `gorm:"column:sequence_number"`
	CargoCode       string   `gorm:"column:cargo_code"`
	CargoName       string   `gorm:"column:cargo_name"`
	CargoUnit       string   `gorm:"column:cargo_unit"`
	Total           *float64 `gorm:"column:total"`
	CargoNature     string   `gorm:"column:cargo_nature"`
	CargoNatureDesc string   `gorm:"column:cargo_nature_desc"`
	CargoPackaging  string   `gorm:"column:cargo_packaging"`
	Stowage         string   `gorm:"column:stowage"`
	ConsigneeCode   string   `gorm:"column:consignee_code"`
	ConsigneeName   string   `gorm:"column:consignee_name"`
}

func (r *opsPlanRepository) findDeterminationRequestHeader(ctx context.Context, branchCode, terminalCode int, ppkNumber, activityCode string) (determinationRequestHeader, error) {
	headers, err := r.findDeterminationRequestHeaders(ctx, branchCode, terminalCode, ppkNumber, activityCode)
	if err != nil {
		return determinationRequestHeader{}, err
	}
	return headers[0], nil
}

func (r *opsPlanRepository) findDeterminationRequestHeaders(ctx context.Context, branchCode, terminalCode int, ppkNumber, activityCode string) ([]determinationRequestHeader, error) {
	var headers []determinationRequestHeader
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			branch_name,
			terminal_name,
			request_code,
			vessel_code,
			vessel_name,
			vessel_type,
			voyage_type,
			agent_name,
			pbm_code,
			pbm_name,
			activity_name
		FROM plan.post_requests
		WHERE branch_code = ?
			AND terminal_code = ?
			AND ppk_number = ?
			AND activity_code = ?
			AND status = 1
		ORDER BY request_date DESC, id DESC
	`, branchCode, terminalCode, ppkNumber, activityCode).Scan(&headers).Error
	if err != nil {
		return nil, err
	}
	if len(headers) == 0 {
		return nil, fmt.Errorf("approved post_request not found for ppk_number %s and activity_code %s", ppkNumber, activityCode)
	}
	return headers, nil
}

func (r *opsPlanRepository) findDeterminationRequestDetails(ctx context.Context, branchCode, terminalCode int, requestCode string) ([]determinationRequestDetail, error) {
	detailTable, err := r.resolvePostRequestDetailTable(ctx)
	if err != nil {
		return nil, err
	}

	var details []determinationRequestDetail
	if err := r.db.WithContext(ctx).Raw(fmt.Sprintf(`
		SELECT
			id,
			sequence_number,
			cargo_code,
			cargo_name,
			cargo_unit,
			total,
			cargo_nature,
			cargo_nature_desc,
			cargo_packaging,
			stowage,
			consignee_code,
			consignee_name
		FROM %s
		WHERE branch_code = ?
			AND terminal_code = ?
			AND request_code = ?
		ORDER BY sequence_number ASC, id ASC
	`, detailTable), branchCode, terminalCode, requestCode).Scan(&details).Error; err != nil {
		return nil, err
	}
	return details, nil
}

func (r *opsPlanRepository) GetDataVesselSchedule(ctx context.Context, ppkNumber, vesselCode string) ([]RawJSONResponse, error) {
	tableName, err := r.resolveExistingTable(ctx, []string{"post.vessel_schedule", "plan.post_vessel_schedules"})
	if err != nil {
		return nil, err
	}

	var rows []RawJSONResponse
	if err := r.db.WithContext(ctx).Raw(
		getDataVesselScheduleQuery(tableName),
		ppkNumber,
		ppkNumber,
		ppkNumber,
		ppkNumber,
		vesselCode,
		vesselCode,
	).Scan(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *opsPlanRepository) GetDataVesel(ctx context.Context, vesselCode string) ([]RawJSONResponse, error) {
	vesselTable, err := r.resolveExistingTableWithDB(ctx, r.admDB, []string{"adm.posm_vessel", "posm_vessel"})
	if err != nil {
		return nil, err
	}
	vesselDetailTable, err := r.resolveExistingTableWithDB(ctx, r.admDB, []string{"adm.posm_vessel_d", "posm_vessel_d"})
	if err != nil {
		return nil, err
	}

	var rows []RawJSONResponse
	if err := r.admDB.WithContext(ctx).Raw(
		getDataVeselQuery(vesselTable, vesselDetailTable),
		vesselCode,
	).Scan(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *opsPlanRepository) Create(ctx context.Context, header *LoadingUnloadingPlan, details []LoadingUnloadingPlanDetail, equipmentPlans []PostEquipmentPlan) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("LOCK TABLE plan.post_vessel_plan IN EXCLUSIVE MODE").Error; err != nil {
			return err
		}
		planNumber, err := r.nextPlanNumber(tx, header.BranchCode, header.TerminalCode, header.PlanDate)
		if err != nil {
			return err
		}
		header.PlanNumber = planNumber

		if err := tx.Create(header).Error; err != nil {
			return err
		}

		nextDetailSequence, err := r.nextPlanDetailSequence(tx, header.BranchCode, header.TerminalCode, header.PlanDate)
		if err != nil {
			return err
		}
		detailPrefix := fmt.Sprintf("OPD%d%d%s", header.BranchCode, header.TerminalCode, header.PlanDate.Format("200601"))
		for i := range details {
			details[i].BranchCode = header.BranchCode
			details[i].TerminalCode = header.TerminalCode
			details[i].PlanNumber = header.PlanNumber
			details[i].PlanDetailCode = fmt.Sprintf("%s%06d", detailPrefix, nextDetailSequence+i)
		}
		if len(details) > 0 {
			if err := tx.CreateInBatches(details, 100).Error; err != nil {
				return err
			}
		}

		if len(equipmentPlans) == 0 {
			equipmentPlans = buildPostEquipmentPlansFromDetails(details)
		}
		if err := preparePostEquipmentPlans(header, details, equipmentPlans); err != nil {
			return err
		}
		if len(equipmentPlans) > 0 {
			if err := tx.CreateInBatches(equipmentPlans, 100).Error; err != nil {
				return err
			}
		}

		requestTable, err := r.resolveExistingTableWithDB(ctx, tx, []string{"plan.post_requests", "plan.post_request"})
		if err != nil {
			return err
		}
		statusColumn, err := r.resolveExistingColumn(ctx, tx, requestTable, []string{"plan_status", "status_plan"})
		if err != nil {
			return err
		}

		result := tx.Exec(
			fmt.Sprintf(
				"UPDATE %s SET %s = 1, last_updated_date = CURRENT_TIMESTAMP, last_updated_by = ? WHERE ppk_number = ? AND activity_code = ? AND status = 1",
				requestTable,
				statusColumn,
			),
			header.CreationBy,
			header.PPKNumber,
			header.ActivityCode,
		)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("approved post_request not found for ppk_number %s and activity_code %s", header.PPKNumber, header.ActivityCode)
		}

		return nil
	})
}

func (r *opsPlanRepository) BuildDeterminationsFromPlan(ctx context.Context, branchCode, terminalCode int, planCode, createdBy string) ([]determinationBuild, error) {
	var planHeader LoadingUnloadingPlan
	if err := r.db.WithContext(ctx).
		Where("branch_code = ? AND terminal_code = ? AND plan_code = ?", branchCode, terminalCode, planCode).
		First(&planHeader).Error; err != nil {
		return nil, err
	}

	requestHeaders, err := r.findDeterminationRequestHeaders(ctx, branchCode, terminalCode, planHeader.PPKNumber, planHeader.ActivityCode)
	if err != nil {
		return nil, err
	}

	var planDetails []LoadingUnloadingPlanDetail
	if err := r.db.WithContext(ctx).
		Where("branch_code = ? AND terminal_code = ? AND plan_code = ?", branchCode, terminalCode, planCode).
		Order("sequence_no ASC").
		Find(&planDetails).Error; err != nil {
		return nil, err
	}
	if len(planDetails) == 0 {
		return nil, fmt.Errorf("plan detail not found for plan_code %s", planCode)
	}

	var equipmentPlans []PostEquipmentPlan
	if err := r.db.WithContext(ctx).
		Where("branch_code = ? AND terminal_code = ? AND plan_code = ?", branchCode, terminalCode, planCode).
		Order("sequence_no ASC").
		Find(&equipmentPlans).Error; err != nil {
		return nil, err
	}

	now := time.Now()
	if createdBy == "" {
		createdBy = "SYSTEM"
	}

	builds := make([]determinationBuild, 0, len(requestHeaders))
	for _, requestHeader := range requestHeaders {
		header := &LoadingUnloadingDetermination{
			BranchCode:        branchCode,
			TerminalCode:      terminalCode,
			BranchName:        firstNonEmpty(planHeader.BranchName, requestHeader.BranchName),
			TerminalName:      firstNonEmpty(planHeader.TerminalName, requestHeader.TerminalName),
			VesselCode:        firstNonEmpty(planHeader.VesselCode, requestHeader.VesselCode),
			VesselName:        firstNonEmpty(planHeader.VesselName, requestHeader.VesselName),
			VesselType:        firstNonEmpty(planHeader.VesselType, requestHeader.VesselType),
			GRT:               floatPtrToString(planHeader.GRT),
			LOA:               floatPtrToString(planHeader.LOA),
			VoyageType:        firstNonEmpty(requestHeader.VoyageType, planHeader.ShippingType),
			AgentName:         firstNonEmpty(planHeader.AgentName, requestHeader.AgentName),
			PPKNumber:         planHeader.PPKNumber,
			RequestCode:       requestHeader.RequestCode,
			PlanCode:          planHeader.PlanNumber,
			PlanDate:          planHeader.PlanDate,
			DeterminationDate: now,
			ETA:               planHeader.ETA,
			ETD:               planHeader.ETD,
			PBMCode:           requestHeader.PBMCode,
			PBMName:           requestHeader.PBMName,
			ActivityCode:      planHeader.ActivityCode,
			ActivityName:      firstNonEmpty(planHeader.ActivityName, requestHeader.ActivityName),
			Remarks:           planHeader.Remarks,
			Status:            intPtrToString(planHeader.Status),
			ProgramName:       programName,
			Cycle:             planHeader.Cycle,
			CreationDate:      now,
			CreationBy:        createdBy,
			LastUpdatedDate:   &now,
			LastUpdatedBy:     createdBy,
		}

		selectedPlanDetails := filterPlanDetailsForRequest(planDetails, requestHeader, len(requestHeaders))
		if len(selectedPlanDetails) == 0 {
			continue
		}

		requestDetails, err := r.findDeterminationRequestDetails(ctx, branchCode, terminalCode, requestHeader.RequestCode)
		if err != nil {
			return nil, err
		}
		requestDetailsBySequence := make(map[int]determinationRequestDetail, len(requestDetails))
		for _, detail := range requestDetails {
			if detail.SequenceNumber != nil {
				requestDetailsBySequence[*detail.SequenceNumber] = detail
			}
		}

		details := make([]LoadingUnloadingDeterminationDetail, 0, len(selectedPlanDetails))
		for i, planDetail := range selectedPlanDetails {
			requestDetail, ok := requestDetailsBySequence[planDetail.SequenceNo]
			if !ok && i < len(requestDetails) {
				requestDetail = requestDetails[i]
			}

			details = append(details, LoadingUnloadingDeterminationDetail{
				BranchCode:      branchCode,
				TerminalCode:    terminalCode,
				RequestCode:     requestHeader.RequestCode,
				SequenceNo:      planDetail.SequenceNo,
				ActivityDate:    planDetail.ActivityDate,
				Stowage:         planDetail.Stowage,
				CargoCode:       planDetail.CargoCode,
				CargoName:       planDetail.CargoName,
				TotalQuantity:   firstFloat(planDetail.PlannedQuantity, planDetail.TotalQuantity),
				CargoUnit:       planDetail.CargoUnit,
				CargoPackaging:  planDetail.CargoPackaging,
				DayNo:           planDetail.DayNo,
				DockCode:        planDetail.DockCode,
				DockName:        planDetail.DockName,
				BerthCode:       planDetail.BerthCode,
				BerthName:       planDetail.BerthName,
				Shift1:          planDetail.Shift1,
				Shift2:          planDetail.Shift2,
				Shift3:          planDetail.Shift3,
				PBMCode:         firstNonEmpty(planDetail.PBMCode, requestHeader.PBMCode),
				PBMName:         firstNonEmpty(planDetail.PBMName, requestHeader.PBMName),
				ConsigneeCode:   firstNonEmpty(planDetail.ConsigneeCode, requestDetail.ConsigneeCode),
				ConsigneeName:   firstNonEmpty(planDetail.ConsigneeName, requestDetail.ConsigneeName),
				TruckCount:      planDetail.TruckCount,
				TruckCapacity:   planDetail.TruckCapacity,
				GangCount:       intPtr(0),
				Attribute1:      planDetail.Attrib1,
				Attribute2:      planDetail.Attrib2,
				Attribute3:      planDetail.Attrib3,
				Value1:          planDetail.Val1,
				Value2:          planDetail.Val2,
				Value3:          planDetail.Val3,
				Status:          intPtrToString(planDetail.Status),
				CreationDate:    now,
				CreationBy:      createdBy,
				ProgramName:     programName,
				CargoNature:     firstNonEmpty(planDetail.CargoNature, requestDetail.CargoNature),
				CargoNatureDesc: firstNonEmpty(planDetail.CargoNatureDesc, requestDetail.CargoNatureDesc),
				RequestDetailID: nil,
			})
		}

		selectedEquipmentPlans := filterEquipmentPlansForRequest(equipmentPlans, requestHeader, len(requestHeaders))
		equipmentDeterminations := make([]PostEquipmentDetermination, 0, len(selectedEquipmentPlans))
		for _, equipmentPlan := range selectedEquipmentPlans {
			equipmentDeterminations = append(equipmentDeterminations, PostEquipmentDetermination{
				BranchCode:     branchCode,
				TerminalCode:   terminalCode,
				RequestCode:    requestHeader.RequestCode,
				SequenceNo:     equipmentPlan.SequenceNo,
				EquipmentCode:  equipmentPlan.EquipmentCode,
				EquipmentName:  equipmentPlan.EquipmentName,
				UnitCode:       equipmentPlan.UnitCode,
				PBMCode:        firstNonEmpty(equipmentPlan.PBMCode, requestHeader.PBMCode),
				PBMName:        firstNonEmpty(equipmentPlan.PBMName, requestHeader.PBMName),
				ConsigneeCode:  equipmentPlan.ConsigneeCode,
				ConsigneeName:  equipmentPlan.ConsigneeName,
				Remarks:        equipmentPlan.Description,
				EquipmentGroup: equipmentPlan.EquipmentGroup,
				UnitTon:        equipmentPlan.UnitTon,
				Attribute1:     equipmentPlan.Attr1,
				Attribute2:     equipmentPlan.Attr2,
				Attribute3:     equipmentPlan.Attr3,
				Value1:         equipmentPlan.Value1,
				Value2:         equipmentPlan.Value2,
				Value3:         equipmentPlan.Value3,
				CreationDate:   now,
				CreationBy:     createdBy,
				ProgramName:    programName,
				DayNo:          equipmentPlan.DayNo,
				ActivityDate:   equipmentPlan.ActivityDate,
				Stowage:        equipmentPlan.Stowage,
			})
		}

		builds = append(builds, determinationBuild{
			Header:                  header,
			Details:                 details,
			EquipmentDeterminations: equipmentDeterminations,
		})
	}

	if len(builds) == 0 {
		return nil, fmt.Errorf("approved request not found for plan_code %s", planCode)
	}

	return builds, nil
}

func filterPlanDetailsForRequest(details []LoadingUnloadingPlanDetail, requestHeader determinationRequestHeader, requestCount int) []LoadingUnloadingPlanDetail {
	pbmCode := normalizeCode(requestHeader.PBMCode)
	if pbmCode == "" {
		return details
	}

	selected := make([]LoadingUnloadingPlanDetail, 0, len(details))
	hasPBMOnPlan := false
	for _, detail := range details {
		detailPBMCode := normalizeCode(detail.PBMCode)
		if detailPBMCode != "" {
			hasPBMOnPlan = true
		}
		if detailPBMCode == pbmCode {
			selected = append(selected, detail)
		}
	}
	if len(selected) == 0 && (!hasPBMOnPlan || requestCount == 1) {
		return details
	}
	return selected
}

func filterEquipmentPlansForRequest(equipmentPlans []PostEquipmentPlan, requestHeader determinationRequestHeader, requestCount int) []PostEquipmentPlan {
	pbmCode := normalizeCode(requestHeader.PBMCode)
	if pbmCode == "" {
		return equipmentPlans
	}

	selected := make([]PostEquipmentPlan, 0, len(equipmentPlans))
	hasPBMOnPlan := false
	for _, equipmentPlan := range equipmentPlans {
		equipmentPBMCode := normalizeCode(equipmentPlan.PBMCode)
		if equipmentPBMCode != "" {
			hasPBMOnPlan = true
		}
		if equipmentPBMCode == pbmCode {
			selected = append(selected, equipmentPlan)
		}
	}
	if len(selected) == 0 && (!hasPBMOnPlan || requestCount == 1) {
		return equipmentPlans
	}
	return selected
}

func normalizeCode(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}

type planPBMWorkOrderKey struct {
	BranchCode   int
	TerminalCode int
	PlanCode     string
	PBMCode      string
}

type planPBMDeterminationCodes struct {
	DeterminationCode string
	WorkOrderCode     string
}

func (r *opsPlanRepository) CreateDeterminations(ctx context.Context, builds []determinationBuild) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("LOCK TABLE plan.post_vessel_confirmed_plan IN EXCLUSIVE MODE").Error; err != nil {
			return err
		}

		workOrderSequences := make(map[string]int)
		workOrderCodesByPBM := make(map[string]string)
		planCodesByPBM := make(map[planPBMWorkOrderKey]planPBMDeterminationCodes)
		for buildIndex := range builds {
			header := builds[buildIndex].Header
			if header == nil {
				return fmt.Errorf("determination header is required")
			}
			if strings.TrimSpace(header.DeterminationCode) == "" {
				determinationCode, err := r.nextDeterminationCode(tx, header.BranchCode, header.TerminalCode, header.DeterminationDate)
				if err != nil {
					return err
				}
				header.DeterminationCode = determinationCode
			}

			if err := tx.Create(header).Error; err != nil {
				return err
			}

			sequenceKey := workOrderSequenceKey(header.BranchCode, header.TerminalCode, header.DeterminationDate)
			nextWorkOrderSequence, ok := workOrderSequences[sequenceKey]
			if !ok {
				var err error
				nextWorkOrderSequence, err = r.nextWorkOrderSequence(tx, header.BranchCode, header.TerminalCode, header.DeterminationDate)
				if err != nil {
					return err
				}
			}
			workOrderPrefix := fmt.Sprintf("SPMK%d%d%s", header.BranchCode, header.TerminalCode, header.DeterminationDate.Format("200601"))

			for i := range builds[buildIndex].Details {
				details := builds[buildIndex].Details
				details[i].BranchCode = header.BranchCode
				details[i].TerminalCode = header.TerminalCode
				details[i].DeterminationCode = header.DeterminationCode
				details[i].RequestCode = header.RequestCode
				if strings.TrimSpace(details[i].WorkOrderCode) == "" {
					pbmCode := normalizeCode(firstNonEmpty(details[i].PBMCode, header.PBMCode))
					workOrderKey := workOrderSequenceKey(header.BranchCode, header.TerminalCode, header.DeterminationDate) + "|" + pbmCode
					workOrderCode, ok := workOrderCodesByPBM[workOrderKey]
					if !ok {
						workOrderCode = fmt.Sprintf("%s%06d", workOrderPrefix, nextWorkOrderSequence)
						workOrderCodesByPBM[workOrderKey] = workOrderCode
						nextWorkOrderSequence++
					}
					details[i].WorkOrderCode = workOrderCode
				}
				planCodesByPBM[planPBMWorkOrderKey{
					BranchCode:   header.BranchCode,
					TerminalCode: header.TerminalCode,
					PlanCode:     header.PlanCode,
					PBMCode:      normalizeCode(firstNonEmpty(details[i].PBMCode, header.PBMCode)),
				}] = planPBMDeterminationCodes{
					DeterminationCode: header.DeterminationCode,
					WorkOrderCode:     details[i].WorkOrderCode,
				}
				if strings.TrimSpace(details[i].CreationBy) == "" {
					details[i].CreationBy = header.CreationBy
				}
				if strings.TrimSpace(details[i].ProgramName) == "" {
					details[i].ProgramName = header.ProgramName
				}
				if details[i].CreationDate.IsZero() {
					details[i].CreationDate = header.CreationDate
				}
				builds[buildIndex].Details = details
			}
			workOrderSequences[sequenceKey] = nextWorkOrderSequence

			if len(builds[buildIndex].Details) > 0 {
				if err := tx.CreateInBatches(builds[buildIndex].Details, 100).Error; err != nil {
					return err
				}
			}

			for i := range builds[buildIndex].EquipmentDeterminations {
				equipmentDeterminations := builds[buildIndex].EquipmentDeterminations
				equipmentDeterminations[i].BranchCode = header.BranchCode
				if equipmentDeterminations[i].TerminalCode == 0 {
					equipmentDeterminations[i].TerminalCode = header.TerminalCode
				}
				equipmentDeterminations[i].DeterminationCode = header.DeterminationCode
				equipmentDeterminations[i].RequestCode = header.RequestCode
				if strings.TrimSpace(equipmentDeterminations[i].CreationBy) == "" {
					equipmentDeterminations[i].CreationBy = header.CreationBy
				}
				if strings.TrimSpace(equipmentDeterminations[i].ProgramName) == "" {
					equipmentDeterminations[i].ProgramName = header.ProgramName
				}
				if equipmentDeterminations[i].CreationDate.IsZero() {
					equipmentDeterminations[i].CreationDate = header.CreationDate
				}
				builds[buildIndex].EquipmentDeterminations = equipmentDeterminations
			}
			if len(builds[buildIndex].EquipmentDeterminations) > 0 {
				if err := tx.CreateInBatches(builds[buildIndex].EquipmentDeterminations, 100).Error; err != nil {
					return err
				}
			}

			if err := tx.Model(&LoadingUnloadingPlan{}).
				Where("branch_code = ? AND terminal_code = ? AND plan_code = ?", header.BranchCode, header.TerminalCode, header.PlanCode).
				Updates(map[string]interface{}{
					"status":            1,
					"last_updated_date": time.Now(),
					"last_updated_by":   header.CreationBy,
				}).Error; err != nil {
				return err
			}
		}

		for key, codes := range planCodesByPBM {
			if err := tx.Model(&LoadingUnloadingPlanDetail{}).
				Where(
					"branch_code = ? AND terminal_code = ? AND plan_code = ? AND UPPER(TRIM(COALESCE(pbm_code, ''))) = ?",
					key.BranchCode,
					key.TerminalCode,
					key.PlanCode,
					key.PBMCode,
				).
				Updates(map[string]interface{}{
					"confirmed_plan_code": codes.DeterminationCode,
					"work_order_code":     codes.WorkOrderCode,
				}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *opsPlanRepository) Update(ctx context.Context, branchCode, terminalCode int, input *UpdateLoadingUnloadingPlanInput, details []LoadingUnloadingPlanDetail, replaceDetails bool, equipmentPlans []PostEquipmentPlan, replaceEquipmentPlans bool, updatedBy string) (*LoadingUnloadingPlan, []LoadingUnloadingPlanDetail, []PostEquipmentPlan, error) {
	var finalHeader LoadingUnloadingPlan
	var finalDetails []LoadingUnloadingPlanDetail
	var finalEquipmentPlans []PostEquipmentPlan

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		planCode := input.PlanIdentifier()
		if err := tx.Where(
			"branch_code = ? AND terminal_code = ? AND plan_code = ?",
			branchCode,
			terminalCode,
			planCode,
		).First(&finalHeader).Error; err != nil {
			return err
		}

		now := time.Now()
		updates := map[string]interface{}{
			"last_updated_date": now,
			"last_updated_by":   updatedBy,
		}
		if input.TotalDays != nil {
			updates["total_days"] = *input.TotalDays
			finalHeader.TotalDays = input.TotalDays
		}
		if input.TotalShifts != nil {
			updates["total_shifts"] = *input.TotalShifts
			finalHeader.TotalShifts = input.TotalShifts
		}
		if input.ActivityStartDate != nil {
			updates["activity_start_date"] = *input.ActivityStartDate
			finalHeader.ActivityStartDate = input.ActivityStartDate
		}
		if input.ActivityEndDate != nil {
			updates["activity_end_date"] = *input.ActivityEndDate
			finalHeader.ActivityEndDate = input.ActivityEndDate
		}

		if err := tx.Model(&LoadingUnloadingPlan{}).
			Where("id = ?", finalHeader.ID).
			Updates(updates).Error; err != nil {
			return err
		}
		finalHeader.LastUpdatedDate = &now
		finalHeader.LastUpdatedBy = updatedBy

		equipmentPlansDeleted := false
		if replaceDetails {
			if err := tx.Where(
				"branch_code = ? AND terminal_code = ? AND plan_code = ?",
				branchCode,
				terminalCode,
				planCode,
			).Delete(&PostEquipmentPlan{}).Error; err != nil {
				return err
			}
			equipmentPlansDeleted = true

			if err := tx.Where(
				"branch_code = ? AND terminal_code = ? AND plan_code = ?",
				branchCode,
				terminalCode,
				planCode,
			).Delete(&LoadingUnloadingPlanDetail{}).Error; err != nil {
				return err
			}

			nextDetailSequence, err := r.nextPlanDetailSequence(tx, branchCode, terminalCode, finalHeader.PlanDate)
			if err != nil {
				return err
			}
			detailPrefix := fmt.Sprintf("OPD%d%d%s", branchCode, terminalCode, finalHeader.PlanDate.Format("200601"))
			for i := range details {
				details[i].BranchCode = branchCode
				details[i].TerminalCode = terminalCode
				details[i].PlanNumber = planCode
				details[i].PlanDetailCode = fmt.Sprintf("%s%06d", detailPrefix, nextDetailSequence+i)
			}
			if len(details) > 0 {
				if err := tx.CreateInBatches(details, 100).Error; err != nil {
					return err
				}
			}
		}

		if err := tx.Where(
			"branch_code = ? AND terminal_code = ? AND plan_code = ?",
			branchCode,
			terminalCode,
			planCode,
		).Order("sequence_no ASC").Find(&finalDetails).Error; err != nil {
			return err
		}

		if replaceEquipmentPlans {
			if !equipmentPlansDeleted {
				if err := tx.Where(
					"branch_code = ? AND terminal_code = ? AND plan_code = ?",
					branchCode,
					terminalCode,
					planCode,
				).Delete(&PostEquipmentPlan{}).Error; err != nil {
					return err
				}
			}
			if len(equipmentPlans) == 0 && replaceDetails {
				equipmentPlans = buildPostEquipmentPlansFromDetails(finalDetails)
			}
			if err := preparePostEquipmentPlans(&finalHeader, finalDetails, equipmentPlans); err != nil {
				return err
			}
			if len(equipmentPlans) > 0 {
				if err := tx.CreateInBatches(equipmentPlans, 100).Error; err != nil {
					return err
				}
			}
		}

		return tx.Where(
			"branch_code = ? AND terminal_code = ? AND plan_code = ?",
			branchCode,
			terminalCode,
			planCode,
		).Order("sequence_no ASC").Find(&finalEquipmentPlans).Error
	})
	if err != nil {
		return nil, nil, nil, err
	}

	return &finalHeader, finalDetails, finalEquipmentPlans, nil
}

func (r *opsPlanRepository) UpdateDeterminedPlan(ctx context.Context, branchCode, terminalCode int, input *UpdateLoadingUnloadingPlanInput, details []LoadingUnloadingPlanDetail, replaceDetails bool, equipmentPlans []PostEquipmentPlan, replaceEquipmentPlans bool, updatedBy string) (*LoadingUnloadingPlan, []LoadingUnloadingPlanDetail, []PostEquipmentPlan, error) {
	var finalHeader LoadingUnloadingPlan
	var finalDetails []LoadingUnloadingPlanDetail
	var finalEquipmentPlans []PostEquipmentPlan

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		planCode := input.PlanIdentifier()
		if err := tx.Where(
			"branch_code = ? AND terminal_code = ? AND plan_code = ? AND status IN ?",
			branchCode,
			terminalCode,
			planCode,
			[]int{1, 2},
		).First(&finalHeader).Error; err != nil {
			return err
		}

		var determinations []LoadingUnloadingDetermination
		if err := tx.Where(
			"branch_code = ? AND terminal_code = ? AND plan_code = ?",
			branchCode,
			terminalCode,
			planCode,
		).Order("confirmed_plan_code ASC").Find(&determinations).Error; err != nil {
			return err
		}
		if len(determinations) == 0 {
			return fmt.Errorf("determination not found for plan_code %s", planCode)
		}

		workOrderByPBM, err := r.findExistingWorkOrderByPBM(tx, branchCode, terminalCode, planCode)
		if err != nil {
			return err
		}

		now := time.Now()
		updates := map[string]interface{}{
			"last_updated_date": now,
			"last_updated_by":   updatedBy,
		}
		if input.TotalDays != nil {
			updates["total_days"] = *input.TotalDays
			finalHeader.TotalDays = input.TotalDays
		}
		if input.TotalShifts != nil {
			updates["total_shifts"] = *input.TotalShifts
			finalHeader.TotalShifts = input.TotalShifts
		}
		if input.ActivityStartDate != nil {
			updates["activity_start_date"] = *input.ActivityStartDate
			finalHeader.ActivityStartDate = input.ActivityStartDate
		}
		if input.ActivityEndDate != nil {
			updates["activity_end_date"] = *input.ActivityEndDate
			finalHeader.ActivityEndDate = input.ActivityEndDate
		}

		if err := tx.Model(&LoadingUnloadingPlan{}).
			Where("id = ?", finalHeader.ID).
			Updates(updates).Error; err != nil {
			return err
		}
		finalHeader.LastUpdatedDate = &now
		finalHeader.LastUpdatedBy = updatedBy

		equipmentPlansDeleted := false
		if replaceDetails {
			if err := tx.Where(
				"branch_code = ? AND terminal_code = ? AND plan_code = ?",
				branchCode,
				terminalCode,
				planCode,
			).Delete(&PostEquipmentPlan{}).Error; err != nil {
				return err
			}
			equipmentPlansDeleted = true

			if err := tx.Where(
				"branch_code = ? AND terminal_code = ? AND plan_code = ?",
				branchCode,
				terminalCode,
				planCode,
			).Delete(&LoadingUnloadingPlanDetail{}).Error; err != nil {
				return err
			}

			nextDetailSequence, err := r.nextPlanDetailSequence(tx, branchCode, terminalCode, finalHeader.PlanDate)
			if err != nil {
				return err
			}
			detailPrefix := fmt.Sprintf("OPD%d%d%s", branchCode, terminalCode, finalHeader.PlanDate.Format("200601"))
			for i := range details {
				details[i].BranchCode = branchCode
				details[i].TerminalCode = terminalCode
				details[i].PlanNumber = planCode
				details[i].PlanDetailCode = fmt.Sprintf("%s%06d", detailPrefix, nextDetailSequence+i)

				pbmCode := normalizeCode(details[i].PBMCode)
				workOrderCode := workOrderByPBM[pbmCode]
				if strings.TrimSpace(workOrderCode) == "" {
					return fmt.Errorf("existing work_order_code not found for pbm_code %s", details[i].PBMCode)
				}
				details[i].WorkOrderCode = workOrderCode
			}
			if len(details) > 0 {
				if err := tx.CreateInBatches(details, 100).Error; err != nil {
					return err
				}
			}
		}

		if err := tx.Where(
			"branch_code = ? AND terminal_code = ? AND plan_code = ?",
			branchCode,
			terminalCode,
			planCode,
		).Order("sequence_no ASC").Find(&finalDetails).Error; err != nil {
			return err
		}

		if replaceEquipmentPlans {
			if !equipmentPlansDeleted {
				if err := tx.Where(
					"branch_code = ? AND terminal_code = ? AND plan_code = ?",
					branchCode,
					terminalCode,
					planCode,
				).Delete(&PostEquipmentPlan{}).Error; err != nil {
					return err
				}
			}
			if len(equipmentPlans) == 0 && replaceDetails {
				equipmentPlans = buildPostEquipmentPlansFromDetails(finalDetails)
			}
			if err := preparePostEquipmentPlans(&finalHeader, finalDetails, equipmentPlans); err != nil {
				return err
			}
			if len(equipmentPlans) > 0 {
				if err := tx.CreateInBatches(equipmentPlans, 100).Error; err != nil {
					return err
				}
			}
		}

		if err := tx.Where(
			"branch_code = ? AND terminal_code = ? AND plan_code = ?",
			branchCode,
			terminalCode,
			planCode,
		).Order("sequence_no ASC").Find(&finalEquipmentPlans).Error; err != nil {
			return err
		}

		determinationCodes := make([]string, 0, len(determinations))
		for _, determination := range determinations {
			determinationCodes = append(determinationCodes, determination.DeterminationCode)
			if err := tx.Model(&LoadingUnloadingDetermination{}).
				Where("id = ?", determination.ID).
				Updates(map[string]interface{}{
					"plan_date":         finalHeader.PlanDate,
					"eta":               finalHeader.ETA,
					"etd":               finalHeader.ETD,
					"last_updated_date": now,
					"last_updated_by":   updatedBy,
				}).Error; err != nil {
				return err
			}
		}

		if replaceDetails {
			if err := tx.Where(
				"branch_code = ? AND terminal_code = ? AND confirmed_plan_code IN ?",
				branchCode,
				terminalCode,
				determinationCodes,
			).Delete(&LoadingUnloadingDeterminationDetail{}).Error; err != nil {
				return err
			}

			for _, determination := range determinations {
				selectedDetails := filterPlanDetailsForRequest(finalDetails, determinationRequestHeader{
					RequestCode: determination.RequestCode,
					PBMCode:     determination.PBMCode,
					PBMName:     determination.PBMName,
				}, len(determinations))
				determinationDetails := buildDeterminationDetailsFromPlanDetails(selectedDetails, determination, updatedBy, now)
				if len(determinationDetails) > 0 {
					if err := tx.CreateInBatches(determinationDetails, 100).Error; err != nil {
						return err
					}
				}
			}
		}

		if replaceEquipmentPlans {
			if err := tx.Where(
				"branch_code = ? AND terminal_code = ? AND confirmed_plan_code IN ?",
				branchCode,
				terminalCode,
				determinationCodes,
			).Delete(&PostEquipmentDetermination{}).Error; err != nil {
				return err
			}

			for _, determination := range determinations {
				selectedEquipmentPlans := filterEquipmentPlansForRequest(finalEquipmentPlans, determinationRequestHeader{
					RequestCode: determination.RequestCode,
					PBMCode:     determination.PBMCode,
					PBMName:     determination.PBMName,
				}, len(determinations))
				equipmentDeterminations := buildEquipmentDeterminationsFromPlans(selectedEquipmentPlans, determination, updatedBy, now)
				if len(equipmentDeterminations) > 0 {
					if err := tx.CreateInBatches(equipmentDeterminations, 100).Error; err != nil {
						return err
					}
				}
			}
		}

		return nil
	})
	if err != nil {
		return nil, nil, nil, err
	}

	return &finalHeader, finalDetails, finalEquipmentPlans, nil
}

func (r *opsPlanRepository) nextPlanNumber(tx *gorm.DB, branchCode, terminalCode int, planDate time.Time) (string, error) {
	prefix := fmt.Sprintf("OP%d%d%s", branchCode, terminalCode, planDate.Format("200601"))
	startPosition := len(prefix) + 1

	var lastSequence int
	if err := tx.Raw(`
		SELECT COALESCE(MAX(CAST(SUBSTRING(plan_code FROM ? FOR 6) AS INTEGER)), 0)
		FROM plan.post_vessel_plan
		WHERE plan_code LIKE ?
			AND SUBSTRING(plan_code FROM ? FOR 6) ~ '^[0-9]{6}$'
	`, startPosition, prefix+"%", startPosition).Scan(&lastSequence).Error; err != nil {
		return "", err
	}

	nextSequence := lastSequence + 1
	if nextSequence > 999999 {
		return "", fmt.Errorf("plan number sequence is exhausted for %s", prefix)
	}

	return fmt.Sprintf("%s%06d", prefix, nextSequence), nil
}

func (r *opsPlanRepository) nextPlanDetailSequence(tx *gorm.DB, branchCode, terminalCode int, planDate time.Time) (int, error) {
	prefix := fmt.Sprintf("OPD%d%d%s", branchCode, terminalCode, planDate.Format("200601"))
	startPosition := len(prefix) + 1

	var lastSequence int
	if err := tx.Raw(`
		SELECT COALESCE(MAX(CAST(SUBSTRING(plan_detail_code FROM ? FOR 6) AS INTEGER)), 0)
		FROM plan.post_vessel_plan_d
		WHERE plan_detail_code LIKE ?
			AND SUBSTRING(plan_detail_code FROM ? FOR 6) ~ '^[0-9]{6}$'
	`, startPosition, prefix+"%", startPosition).Scan(&lastSequence).Error; err != nil {
		return 0, err
	}

	nextSequence := lastSequence + 1
	if nextSequence > 999999 {
		return 0, fmt.Errorf("plan detail code sequence is exhausted for %s", prefix)
	}

	return nextSequence, nil
}

func (r *opsPlanRepository) nextDeterminationCode(tx *gorm.DB, branchCode, terminalCode int, determinationDate time.Time) (string, error) {
	prefix := fmt.Sprintf("PNTP%d%d%s", branchCode, terminalCode, determinationDate.Format("200601"))
	startPosition := len(prefix) + 1

	var lastSequence int
	if err := tx.Raw(`
		SELECT COALESCE(MAX(CAST(SUBSTRING(confirmed_plan_code FROM ? FOR 6) AS INTEGER)), 0)
		FROM plan.post_vessel_confirmed_plan
		WHERE confirmed_plan_code LIKE ?
			AND SUBSTRING(confirmed_plan_code FROM ? FOR 6) ~ '^[0-9]{6}$'
	`, startPosition, prefix+"%", startPosition).Scan(&lastSequence).Error; err != nil {
		return "", err
	}

	nextSequence := lastSequence + 1
	if nextSequence > 999999 {
		return "", fmt.Errorf("determination code sequence is exhausted for %s", prefix)
	}

	return fmt.Sprintf("%s%06d", prefix, nextSequence), nil
}

func (r *opsPlanRepository) nextWorkOrderSequence(tx *gorm.DB, branchCode, terminalCode int, determinationDate time.Time) (int, error) {
	prefix := fmt.Sprintf("SPMK%d%d%s", branchCode, terminalCode, determinationDate.Format("200601"))
	startPosition := len(prefix) + 1

	var lastSequence int
	if err := tx.Raw(`
		SELECT COALESCE(MAX(CAST(SUBSTRING(work_order_code FROM ? FOR 6) AS INTEGER)), 0)
		FROM plan.post_vessel_confirmed_plan_d
		WHERE work_order_code LIKE ?
			AND SUBSTRING(work_order_code FROM ? FOR 6) ~ '^[0-9]{6}$'
	`, startPosition, prefix+"%", startPosition).Scan(&lastSequence).Error; err != nil {
		return 0, err
	}

	nextSequence := lastSequence + 1
	if nextSequence > 999999 {
		return 0, fmt.Errorf("work order code sequence is exhausted for %s", prefix)
	}

	return nextSequence, nil
}

func workOrderSequenceKey(branchCode, terminalCode int, determinationDate time.Time) string {
	return fmt.Sprintf("%d|%d|%s", branchCode, terminalCode, determinationDate.Format("200601"))
}

type workOrderPBMRow struct {
	PBMCode       string `gorm:"column:pbm_code"`
	WorkOrderCode string `gorm:"column:work_order_code"`
}

func (r *opsPlanRepository) findExistingWorkOrderByPBM(tx *gorm.DB, branchCode, terminalCode int, planCode string) (map[string]string, error) {
	var rows []workOrderPBMRow
	if err := tx.Raw(`
		SELECT pbm_code, work_order_code
		FROM plan.post_vessel_plan_d
		WHERE branch_code = ?
			AND terminal_code = ?
			AND plan_code = ?
			AND NULLIF(TRIM(COALESCE(work_order_code, '')), '') IS NOT NULL
		ORDER BY sequence_no ASC
	`, branchCode, terminalCode, planCode).Scan(&rows).Error; err != nil {
		return nil, err
	}

	workOrderByPBM := make(map[string]string, len(rows))
	for _, row := range rows {
		pbmCode := normalizeCode(row.PBMCode)
		if pbmCode == "" {
			continue
		}
		if strings.TrimSpace(workOrderByPBM[pbmCode]) == "" {
			workOrderByPBM[pbmCode] = row.WorkOrderCode
		}
	}
	return workOrderByPBM, nil
}

func buildDeterminationDetailsFromPlanDetails(planDetails []LoadingUnloadingPlanDetail, determination LoadingUnloadingDetermination, updatedBy string, now time.Time) []LoadingUnloadingDeterminationDetail {
	details := make([]LoadingUnloadingDeterminationDetail, 0, len(planDetails))
	for _, planDetail := range planDetails {
		details = append(details, LoadingUnloadingDeterminationDetail{
			BranchCode:        determination.BranchCode,
			TerminalCode:      determination.TerminalCode,
			DeterminationCode: determination.DeterminationCode,
			RequestCode:       determination.RequestCode,
			WorkOrderCode:     planDetail.WorkOrderCode,
			SequenceNo:        planDetail.SequenceNo,
			ActivityDate:      planDetail.ActivityDate,
			Stowage:           planDetail.Stowage,
			CargoCode:         planDetail.CargoCode,
			CargoName:         planDetail.CargoName,
			TotalQuantity:     firstFloat(planDetail.PlannedQuantity, planDetail.TotalQuantity),
			CargoUnit:         planDetail.CargoUnit,
			CargoPackaging:    planDetail.CargoPackaging,
			DayNo:             planDetail.DayNo,
			DockCode:          planDetail.DockCode,
			DockName:          planDetail.DockName,
			BerthCode:         planDetail.BerthCode,
			BerthName:         planDetail.BerthName,
			Shift1:            planDetail.Shift1,
			Shift2:            planDetail.Shift2,
			Shift3:            planDetail.Shift3,
			PBMCode:           firstNonEmpty(planDetail.PBMCode, determination.PBMCode),
			PBMName:           firstNonEmpty(planDetail.PBMName, determination.PBMName),
			ConsigneeCode:     planDetail.ConsigneeCode,
			ConsigneeName:     planDetail.ConsigneeName,
			TruckCount:        planDetail.TruckCount,
			TruckCapacity:     planDetail.TruckCapacity,
			GangCount:         intPtr(0),
			Attribute1:        planDetail.Attrib1,
			Attribute2:        planDetail.Attrib2,
			Attribute3:        planDetail.Attrib3,
			Value1:            planDetail.Val1,
			Value2:            planDetail.Val2,
			Value3:            planDetail.Val3,
			Status:            intPtrToString(planDetail.Status),
			CreationDate:      now,
			CreationBy:        updatedBy,
			ProgramName:       programName,
			CargoNature:       planDetail.CargoNature,
			CargoNatureDesc:   planDetail.CargoNatureDesc,
			RequestDetailID:   nil,
		})
	}
	return details
}

func buildEquipmentDeterminationsFromPlans(equipmentPlans []PostEquipmentPlan, determination LoadingUnloadingDetermination, updatedBy string, now time.Time) []PostEquipmentDetermination {
	equipmentDeterminations := make([]PostEquipmentDetermination, 0, len(equipmentPlans))
	for _, equipmentPlan := range equipmentPlans {
		equipmentDeterminations = append(equipmentDeterminations, PostEquipmentDetermination{
			BranchCode:        determination.BranchCode,
			TerminalCode:      determination.TerminalCode,
			RequestCode:       determination.RequestCode,
			DeterminationCode: determination.DeterminationCode,
			SequenceNo:        equipmentPlan.SequenceNo,
			EquipmentCode:     equipmentPlan.EquipmentCode,
			EquipmentName:     equipmentPlan.EquipmentName,
			UnitCode:          equipmentPlan.UnitCode,
			PBMCode:           firstNonEmpty(equipmentPlan.PBMCode, determination.PBMCode),
			PBMName:           firstNonEmpty(equipmentPlan.PBMName, determination.PBMName),
			ConsigneeCode:     equipmentPlan.ConsigneeCode,
			ConsigneeName:     equipmentPlan.ConsigneeName,
			Remarks:           equipmentPlan.Description,
			EquipmentGroup:    equipmentPlan.EquipmentGroup,
			UnitTon:           equipmentPlan.UnitTon,
			Attribute1:        equipmentPlan.Attr1,
			Attribute2:        equipmentPlan.Attr2,
			Attribute3:        equipmentPlan.Attr3,
			Value1:            equipmentPlan.Value1,
			Value2:            equipmentPlan.Value2,
			Value3:            equipmentPlan.Value3,
			CreationDate:      now,
			CreationBy:        updatedBy,
			ProgramName:       programName,
			DayNo:             equipmentPlan.DayNo,
			ActivityDate:      equipmentPlan.ActivityDate,
			Stowage:           equipmentPlan.Stowage,
		})
	}
	return equipmentDeterminations
}

func buildPostEquipmentPlansFromDetails(details []LoadingUnloadingPlanDetail) []PostEquipmentPlan {
	equipmentPlans := make([]PostEquipmentPlan, 0, len(details))
	for _, detail := range details {
		if strings.TrimSpace(detail.EquipmentCode) == "" {
			continue
		}

		equipmentPlans = append(equipmentPlans, PostEquipmentPlan{
			BranchCode:      detail.BranchCode,
			TerminalCode:    detail.TerminalCode,
			PlanNumber:      detail.PlanNumber,
			SequenceNo:      detail.SequenceNo,
			EquipmentCode:   detail.EquipmentCode,
			EquipmentName:   detail.EquipmentName,
			UnitCode:        detail.CargoUnit,
			PBMCode:         detail.PBMCode,
			PBMName:         detail.PBMName,
			ConsigneeCode:   detail.ConsigneeCode,
			ConsigneeName:   detail.ConsigneeName,
			Description:     detail.CargoName,
			EquipmentGroup:  detail.EquipmentGroup,
			Attr1:           detail.Attrib1,
			Attr2:           detail.Attrib2,
			Attr3:           detail.Attrib3,
			Value1:          detail.Val1,
			Value2:          detail.Val2,
			Value3:          detail.Val3,
			HeaderID:        detail.ID,
			DayNo:           detail.DayNo,
			ActivityDate:    detail.ActivityDate,
			Stowage:         detail.Stowage,
			Quantity:        detail.PlannedQuantity,
			CreationDate:    detail.CreationDate,
			CreationBy:      detail.CreationBy,
			LastUpdatedDate: detail.LastUpdatedDate,
			LastUpdatedBy:   detail.LastUpdatedBy,
			ProgramName:     detail.ProgramName,
		})
	}
	return equipmentPlans
}

func preparePostEquipmentPlans(header *LoadingUnloadingPlan, details []LoadingUnloadingPlanDetail, equipmentPlans []PostEquipmentPlan) error {
	detailIDs := make(map[string]int64, len(details))
	for _, detail := range details {
		detailIDs[postEquipmentPlanDetailKey(detail.DayNo, detail.Stowage, detail.ActivityDate)] = detail.ID
	}

	for i := range equipmentPlans {
		equipmentPlans[i].BranchCode = header.BranchCode
		equipmentPlans[i].TerminalCode = header.TerminalCode
		equipmentPlans[i].PlanNumber = header.PlanNumber
		if equipmentPlans[i].ProgramName == "" {
			equipmentPlans[i].ProgramName = header.ProgramName
		}
		if equipmentPlans[i].CreationDate.IsZero() {
			equipmentPlans[i].CreationDate = header.CreationDate
		}
		if equipmentPlans[i].CreationBy == "" {
			equipmentPlans[i].CreationBy = header.CreationBy
		}
		if equipmentPlans[i].LastUpdatedDate == nil {
			equipmentPlans[i].LastUpdatedDate = header.LastUpdatedDate
		}
		if equipmentPlans[i].LastUpdatedBy == "" {
			equipmentPlans[i].LastUpdatedBy = header.LastUpdatedBy
		}
		if equipmentPlans[i].HeaderID != 0 {
			continue
		}

		key := postEquipmentPlanDetailKey(equipmentPlans[i].DayNo, equipmentPlans[i].Stowage, equipmentPlans[i].ActivityDate)
		headerID, ok := detailIDs[key]
		if !ok {
			return fmt.Errorf(
				"loading/unloading detail not found for detailsEquipement[%d] day_no=%s stowage=%s activity_date=%s",
				i,
				formatNullableInt(equipmentPlans[i].DayNo),
				equipmentPlans[i].Stowage,
				formatNullableDate(equipmentPlans[i].ActivityDate),
			)
		}
		equipmentPlans[i].HeaderID = headerID
	}

	return nil
}

func postEquipmentPlanDetailKey(dayNo *int, stowage string, activityDate *time.Time) string {
	return fmt.Sprintf("%s|%s|%s", formatNullableInt(dayNo), strings.ToUpper(strings.TrimSpace(stowage)), formatNullableDate(activityDate))
}

func formatNullableInt(value *int) string {
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%d", *value)
}

func intPtr(value int) *int {
	return &value
}

func formatNullableDate(value *time.Time) string {
	if value == nil {
		return ""
	}
	return value.Format("2006-01-02")
}

func (r *opsPlanRepository) GetAuthLocation(ctx context.Context, userID uint64) (*OpsPlanAuthLocation, error) {
	userTable, err := r.resolveExistingTableWithDB(ctx, r.admDB, []string{"adm.posm_users", "posm_users"})
	if err != nil {
		return nil, err
	}

	var result OpsPlanAuthLocation
	if err := r.admDB.WithContext(ctx).Raw(
		fmt.Sprintf("SELECT branch_name, terminal_name FROM %s WHERE id = ? LIMIT 1", userTable),
		userID,
	).Scan(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *opsPlanRepository) resolvePostRequestDetailTable(ctx context.Context) (string, error) {
	return r.resolveExistingTable(ctx, []string{"plan.post_request_d", "plan.post_requests_d"})
}

func (r *opsPlanRepository) resolveExistingTable(ctx context.Context, candidates []string) (string, error) {
	return r.resolveExistingTableWithDB(ctx, r.db, candidates)
}

func (r *opsPlanRepository) resolveExistingTableWithDB(ctx context.Context, db *gorm.DB, candidates []string) (string, error) {
	for _, tableName := range candidates {
		var exists bool
		if err := db.WithContext(ctx).Raw("SELECT to_regclass(?) IS NOT NULL", tableName).Scan(&exists).Error; err != nil {
			if isPermissionDenied(err) {
				continue
			}
			return "", err
		}
		if exists {
			return tableName, nil
		}
	}
	return "", fmt.Errorf("table not found: %s", strings.Join(candidates, ", "))
}

func (r *opsPlanRepository) resolveExistingColumn(ctx context.Context, db *gorm.DB, tableName string, candidates []string) (string, error) {
	parts := strings.Split(tableName, ".")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid table name: %s", tableName)
	}
	for _, columnName := range candidates {
		var exists bool
		err := db.WithContext(ctx).Raw(`
			SELECT EXISTS (
				SELECT 1
				FROM information_schema.columns
				WHERE table_schema = ?
					AND table_name = ?
					AND column_name = ?
			)
		`, parts[0], parts[1], columnName).Scan(&exists).Error
		if err != nil {
			return "", err
		}
		if exists {
			return columnName, nil
		}
	}
	return "", fmt.Errorf("column not found on %s: %s", tableName, strings.Join(candidates, ", "))
}

func isPermissionDenied(err error) bool {
	return err != nil && strings.Contains(err.Error(), "SQLSTATE 42501")
}

func readyOpsPlanBaseQuery(detailTable string) string {
	return fmt.Sprintf(`
		SELECT
			a.branch_code,
			a.terminal_code,
			a.terminal_name,
			a.request_date,
			a.request_code,
			a.ppk_number,
			a.vessel_code,
			a.vessel_name,
			a.agent_name,
			a.voyage_type,
			0::numeric AS grt,
			0::numeric AS loa,
			a.vessel_type,
			a.pbm_code,
			a.pbm_name,
			a.activity_name,
			a.activity_code,
			STRING_AGG(COALESCE(b.cargo_name, ''), ';' ORDER BY b.sequence_number, b.id) AS cargo_name_list,
			STRING_AGG(COALESCE(b.total::text, ''), ';' ORDER BY b.sequence_number, b.id) AS total_list,
			STRING_AGG(COALESCE(b.cargo_unit, ''), ';' ORDER BY b.sequence_number, b.id) AS cargo_unit_list,
			SUM(COALESCE(b.total, 0)) AS total
		FROM plan.post_requests a
		JOIN %s b
			ON a.branch_code = b.branch_code
			AND a.terminal_code = b.terminal_code
			AND a.request_code = b.request_code
		WHERE a.plan_status = 0
			AND a.status = 1
			AND a.activity_code IN ('BONGKAR', 'MUAT')
		GROUP BY
			a.branch_code,
			a.terminal_code,
			a.terminal_name,
			a.request_date,
			a.request_code,
			a.ppk_number,
			a.vessel_code,
			a.vessel_name,
			a.agent_name,
			a.voyage_type,
			a.vessel_type,
			a.pbm_code,
			a.pbm_name,
			a.activity_name,
			a.activity_code
	`, detailTable)
}

func getDataRequestQuery(detailTable string) string {
	return fmt.Sprintf(`
		SELECT
			a.branch_code,
			a.terminal_code,
			a.ppk_number,
			a.activity_code,
			a.pbm_code,
			a.pbm_name,
			b.cargo_code,
			b.cargo_name,
			b.cargo_unit,
			b.cargo_nature,
			b.cargo_packaging,
			b.stowage_code,
			b.stowage,
			b.consignee_code,
			b.consignee_name,
			SUM(COALESCE(b.total, 0)) AS total
		FROM plan.post_requests a
		JOIN %s b
			ON a.branch_code = b.branch_code
			AND a.terminal_code = b.terminal_code
			AND a.request_code = b.request_code
		WHERE a.status = 1
			AND a.activity_code IN ('BONGKAR', 'MUAT')
			AND a.ppk_number = ?
			AND a.activity_code = ?
		GROUP BY
			a.branch_code,
			a.terminal_code,
			a.ppk_number,
			a.activity_code,
			a.pbm_code,
			a.pbm_name,
			b.cargo_code,
			b.cargo_name,
			b.cargo_unit,
			b.cargo_nature,
			b.cargo_packaging,
			b.stowage_code,
			b.stowage,
			b.consignee_code,
			b.consignee_name
	`, detailTable)
}

func getDataOpQuery(whereClause string) string {
	return fmt.Sprintf(`
		WITH selected_plans AS (
			SELECT a.*
			FROM plan.post_vessel_plan a
			WHERE %s
		),
		detail_rows AS (
			SELECT
				b.branch_code,
				b.terminal_code,
				b.plan_code,
				b.sequence_no,
				NULLIF(TRIM(b.berth_name), '') AS berth_name,
				NULLIF(TRIM(b.pbm_code), '') AS pbm_code,
				NULLIF(TRIM(b.pbm_name), '') AS pbm_name,
				NULLIF(TRIM(b.confirmed_plan_code), '') AS confirmed_plan_code,
				NULLIF(TRIM(b.work_order_code), '') AS work_order_code
			FROM plan.post_vessel_plan_d b
			JOIN selected_plans a
				ON a.plan_code = b.plan_code
				AND a.branch_code = b.branch_code
				AND a.terminal_code = b.terminal_code
		),
		berth_agg AS (
			SELECT
				branch_code,
				terminal_code,
				plan_code,
				(ARRAY_AGG(berth_name ORDER BY sequence_no) FILTER (WHERE berth_name IS NOT NULL))[1] AS berth_name
			FROM detail_rows
			GROUP BY branch_code, terminal_code, plan_code
		),
		detail_by_pbm AS (
			SELECT
				branch_code,
				terminal_code,
				plan_code,
				pbm_code,
				(ARRAY_AGG(pbm_name ORDER BY sequence_no) FILTER (WHERE pbm_name IS NOT NULL))[1] AS pbm_name,
				(ARRAY_AGG(confirmed_plan_code ORDER BY sequence_no) FILTER (WHERE confirmed_plan_code IS NOT NULL))[1] AS confirmed_plan_code,
				(ARRAY_AGG(work_order_code ORDER BY sequence_no) FILTER (WHERE work_order_code IS NOT NULL))[1] AS work_order_code,
				MIN(sequence_no) AS sort_no
			FROM detail_rows
			WHERE pbm_code IS NOT NULL
			GROUP BY branch_code, terminal_code, plan_code, pbm_code
		),
		pbm_agg AS (
			SELECT
				branch_code,
				terminal_code,
				plan_code,
				STRING_AGG(pbm_code, '; ' ORDER BY sort_no, pbm_code) AS pbm_code,
				STRING_AGG(COALESCE(pbm_name, ''), '; ' ORDER BY sort_no, pbm_code) AS pbm_name,
				STRING_AGG(COALESCE(confirmed_plan_code, ''), '; ' ORDER BY sort_no, pbm_code) AS confirmed_plan_code,
				STRING_AGG(COALESCE(work_order_code, ''), '; ' ORDER BY sort_no, pbm_code) AS work_order_code
			FROM detail_by_pbm
			GROUP BY branch_code, terminal_code, plan_code
		)
		SELECT
			a.branch_code,
			a.terminal_code,
			a.branch_name,
			a.terminal_name,
			a.plan_code,
			a.plan_date,
			a.eta,
			a.ppk_number,
			a.activity_code,
			a.activity_name,
			a.vessel_type,
			a.vessel_code,
			a.vessel_name,
			a.grt,
			a.loa,
			a.shipping_type,
			COALESCE(ba.berth_name, '') AS berth_name,
			COALESCE(pa.pbm_code, '') AS pbm_code,
			COALESCE(pa.pbm_name, '') AS pbm_name,
			COALESCE(pa.confirmed_plan_code, '') AS confirmed_plan_code,
			COALESCE(pa.work_order_code, '') AS work_order_code,
			a.status,
			COALESCE(rpk.id, 0) AS vessel_rpk_id
		FROM selected_plans a
		LEFT JOIN plan.post_vessel_rpk rpk
			ON a.plan_code = rpk.ops_plan_code
			AND a.branch_code = rpk.branch_code
			AND a.terminal_code = rpk.terminal_code
		LEFT JOIN berth_agg ba
			ON a.plan_code = ba.plan_code
			AND a.branch_code = ba.branch_code
			AND a.terminal_code = ba.terminal_code
		LEFT JOIN pbm_agg pa
			ON a.plan_code = pa.plan_code
			AND a.branch_code = pa.branch_code
			AND a.terminal_code = pa.terminal_code
		ORDER BY a.plan_date DESC, a.plan_code DESC
	`, whereClause)
}

func getDataVesselScheduleQuery(tableName string) string {
	return fmt.Sprintf(`
		WITH selected AS (
			SELECT
				to_jsonb(vs) AS data,
				1 AS priority,
				COALESCE(vs.eta, vs.etb, vs.creation_date) AS sort_date
			FROM %s vs
			WHERE NULLIF(?, '') IS NOT NULL
				AND vs.pkk_number = ?

			UNION ALL

			SELECT
				to_jsonb(vs) AS data,
				2 AS priority,
				COALESCE(vs.eta, vs.etb, vs.creation_date) AS sort_date
			FROM %s vs
			WHERE NOT EXISTS (
					SELECT 1
					FROM %s found
					WHERE NULLIF(?, '') IS NOT NULL
						AND found.pkk_number = ?
				)
				AND NULLIF(?, '') IS NOT NULL
				AND vs.vessel_code = ?
		)
		SELECT data
		FROM selected
		ORDER BY priority, sort_date DESC NULLS LAST
	`, tableName, tableName, tableName)
}

func getDataVeselQuery(vesselTable, vesselDetailTable string) string {
	return fmt.Sprintf(`
		WITH selected AS (
			SELECT
				to_jsonb(v) || jsonb_build_object(
					'details',
					COALESCE((
						SELECT jsonb_agg(to_jsonb(d))
						FROM %s d
						WHERE d.vessel_code = v.vessel_code
					), '[]'::jsonb)
				) AS data
			FROM %s v
			WHERE v.vessel_code = ?
		)
		SELECT data
		FROM selected
	`, vesselDetailTable, vesselTable)
}

func maxInt(value int, fallback int) int {
	if value < fallback {
		return fallback
	}
	return value
}
