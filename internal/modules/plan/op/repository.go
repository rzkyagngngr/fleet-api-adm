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
	GetDetailOp(ctx context.Context, branchCode, terminalCode int, planNumber string) (*LoadingUnloadingPlan, []LoadingUnloadingPlanDetail, []PostEquipmentPlan, error)
	GetDataVesselSchedule(ctx context.Context, ppkNumber, vesselCode string) ([]RawJSONResponse, error)
	GetDataVesel(ctx context.Context, vesselCode string) ([]RawJSONResponse, error)
	Create(ctx context.Context, header *LoadingUnloadingPlan, details []LoadingUnloadingPlanDetail, equipmentPlans []PostEquipmentPlan) error
	Update(ctx context.Context, branchCode, terminalCode int, input *UpdateLoadingUnloadingPlanInput, details []LoadingUnloadingPlanDetail, replaceDetails bool, equipmentPlans []PostEquipmentPlan, replaceEquipmentPlans bool, updatedBy string) (*LoadingUnloadingPlan, []LoadingUnloadingPlanDetail, []PostEquipmentPlan, error)
	GetAuthLocation(ctx context.Context, userID uint64) (*OpsPlanAuthLocation, error)
}

type opsPlanRepository struct {
	db    *gorm.DB
	admDB *gorm.DB
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
		"COALESCE(a.status, 0) <> 2",
		"a.branch_code = ?",
		"a.terminal_code = ?",
	}
	args := []interface{}{branchCode, terminalCode}

	if value := strings.TrimSpace(input.PPKNumber); value != "" {
		whereParts = append(whereParts, "a.ppk_number = ?")
		args = append(args, value)
	}
	if value := strings.TrimSpace(input.PlanNumber); value != "" {
		whereParts = append(whereParts, "a.plan_number = ?")
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

func (r *opsPlanRepository) GetDetailOp(ctx context.Context, branchCode, terminalCode int, planNumber string) (*LoadingUnloadingPlan, []LoadingUnloadingPlanDetail, []PostEquipmentPlan, error) {
	var header LoadingUnloadingPlan
	if err := r.db.WithContext(ctx).
		Where("branch_code = ? AND terminal_code = ? AND plan_number = ?", branchCode, terminalCode, planNumber).
		First(&header).Error; err != nil {
		return nil, nil, nil, err
	}

	var details []LoadingUnloadingPlanDetail
	if err := r.db.WithContext(ctx).
		Where("branch_code = ? AND terminal_code = ? AND plan_number = ?", branchCode, terminalCode, planNumber).
		Order("sequence_no ASC").
		Find(&details).Error; err != nil {
		return nil, nil, nil, err
	}

	var detailsEquipment []PostEquipmentPlan
	if err := r.db.WithContext(ctx).
		Where("branch_code = ? AND terminal_code = ? AND plan_number = ?", branchCode, terminalCode, planNumber).
		Order("sequence_no ASC").
		Find(&detailsEquipment).Error; err != nil {
		return nil, nil, nil, err
	}

	return &header, details, detailsEquipment, nil
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
		if err := tx.Exec("LOCK TABLE plan.post_loading_unloading_plans IN EXCLUSIVE MODE").Error; err != nil {
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

func (r *opsPlanRepository) Update(ctx context.Context, branchCode, terminalCode int, input *UpdateLoadingUnloadingPlanInput, details []LoadingUnloadingPlanDetail, replaceDetails bool, equipmentPlans []PostEquipmentPlan, replaceEquipmentPlans bool, updatedBy string) (*LoadingUnloadingPlan, []LoadingUnloadingPlanDetail, []PostEquipmentPlan, error) {
	var finalHeader LoadingUnloadingPlan
	var finalDetails []LoadingUnloadingPlanDetail
	var finalEquipmentPlans []PostEquipmentPlan

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where(
			"branch_code = ? AND terminal_code = ? AND plan_number = ?",
			branchCode,
			terminalCode,
			input.PlanNumber,
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
				"branch_code = ? AND terminal_code = ? AND plan_number = ?",
				branchCode,
				terminalCode,
				input.PlanNumber,
			).Delete(&PostEquipmentPlan{}).Error; err != nil {
				return err
			}
			equipmentPlansDeleted = true

			if err := tx.Where(
				"branch_code = ? AND terminal_code = ? AND plan_number = ?",
				branchCode,
				terminalCode,
				input.PlanNumber,
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
				details[i].PlanNumber = input.PlanNumber
				details[i].PlanDetailCode = fmt.Sprintf("%s%06d", detailPrefix, nextDetailSequence+i)
			}
			if len(details) > 0 {
				if err := tx.CreateInBatches(details, 100).Error; err != nil {
					return err
				}
			}
		}

		if err := tx.Where(
			"branch_code = ? AND terminal_code = ? AND plan_number = ?",
			branchCode,
			terminalCode,
			input.PlanNumber,
		).Order("sequence_no ASC").Find(&finalDetails).Error; err != nil {
			return err
		}

		if replaceEquipmentPlans {
			if !equipmentPlansDeleted {
				if err := tx.Where(
					"branch_code = ? AND terminal_code = ? AND plan_number = ?",
					branchCode,
					terminalCode,
					input.PlanNumber,
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
			"branch_code = ? AND terminal_code = ? AND plan_number = ?",
			branchCode,
			terminalCode,
			input.PlanNumber,
		).Order("sequence_no ASC").Find(&finalEquipmentPlans).Error
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
		SELECT COALESCE(MAX(CAST(SUBSTRING(plan_number FROM ? FOR 6) AS INTEGER)), 0)
		FROM plan.post_loading_unloading_plans
		WHERE plan_number LIKE ?
			AND SUBSTRING(plan_number FROM ? FOR 6) ~ '^[0-9]{6}$'
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
		FROM plan.post_loading_unloading_plans_d
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
			NULL::numeric AS target_performance,
			NULL::numeric AS target_productivity,
			0::numeric AS grt,
			0::numeric AS loa,
			a.vessel_type,
			a.pbm_code,
			a.pbm_name,
			a.activity_name,
			a.activity_code,
			STRING_AGG(COALESCE(b.cargo_name, ''), ';' ORDER BY b.sequence_number) AS cargo_name_list,
			STRING_AGG(COALESCE(b.total::text, ''), ';' ORDER BY b.sequence_number) AS total_list,
			STRING_AGG(COALESCE(b.cargo_unit, ''), ';' ORDER BY b.sequence_number) AS cargo_unit_list,
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
			b.stowage,
			b.consignee_code,
			b.consignee_name
	`, detailTable)
}

func getDataOpQuery(whereClause string) string {
	return fmt.Sprintf(`
		SELECT
			a.branch_code,
			a.terminal_code,
			a.branch_name,
			a.terminal_name,
			a.plan_number,
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
			(
				ARRAY_AGG(NULLIF(TRIM(b.berth_name), '') ORDER BY b.sequence_no)
				FILTER (WHERE NULLIF(TRIM(b.berth_name), '') IS NOT NULL)
			)[1] AS berth_name,
			COALESCE((
				ARRAY_AGG(DISTINCT NULLIF(TRIM(b.pbm_code), ''))
				FILTER (WHERE NULLIF(TRIM(b.pbm_code), '') IS NOT NULL)
			)[1], '') AS pbm_code,
			COALESCE((
				ARRAY_AGG(DISTINCT NULLIF(TRIM(b.pbm_name), ''))
				FILTER (WHERE NULLIF(TRIM(b.pbm_name), '') IS NOT NULL)
			)[1], '') AS pbm_name,
			a.status
		FROM plan.post_loading_unloading_plans a
		LEFT JOIN plan.post_loading_unloading_plans_d b
			ON a.plan_number = b.plan_number
			AND a.branch_code = b.branch_code
			AND a.terminal_code = b.terminal_code
		WHERE %s
		GROUP BY
			a.branch_code,
			a.terminal_code,
			a.branch_name,
			a.terminal_name,
			a.plan_number,
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
			a.status
		ORDER BY a.plan_date DESC, a.plan_number DESC
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
