package vesselschedule

import (
	"context"
	"fmt"
	"omniport-api/internal/helper"
	"strings"
	"time"

	"gorm.io/gorm"
)

type VesselScheduleService interface {
	Search(ctx context.Context, query helper.PaginationQuery) ([]VesselScheduleSearchResponse, helper.PaginationMeta, error)
	Create(ctx context.Context, schedule *VesselSchedule) error
	Update(ctx context.Context, scheduleCode string, schedule *VesselSchedule) error
	UpdateStatus(ctx context.Context, scheduleCode string, status int, updatedBy string) error
	Delete(ctx context.Context, id uint64) error
	FindByID(ctx context.Context, id uint64) (*VesselSchedule, error)
	FindByScheduleCode(ctx context.Context, scheduleCode string) (*VesselScheduleDetailResponse, error)
	GetAuthLocation(ctx context.Context, userID uint64) (*VesselScheduleAuthLocation, error)
}

type vesselScheduleService struct {
	db     *gorm.DB
	authDB *gorm.DB
}

type VesselScheduleAuthLocation struct {
	BranchName   string
	TerminalName string
}

func NewVesselScheduleService(db *gorm.DB, authDB ...*gorm.DB) VesselScheduleService {
	locationDB := db
	if len(authDB) > 0 && authDB[0] != nil {
		locationDB = authDB[0]
	}
	return &vesselScheduleService{db: db, authDB: locationDB}
}

func (s *vesselScheduleService) Search(ctx context.Context, query helper.PaginationQuery) ([]VesselScheduleSearchResponse, helper.PaginationMeta, error) {
	config := helper.NativePaginationConfig{
		TableName: "plan.post_vessel_schedules",
		SelectColumns: []string{
			"id", "branch_code", "terminal_code", "branch_name", "terminal_name",
			"schedule_code", "vessel_name", "vessel_code", "vessel_type", "voyage_number",
			"vessel_hatch_number", "pkk_number", "ppk_number", "voyage_type", "grt", "loa",
			"agency_name", "port_agent", "emergency_contact", "origin_port_code",
			"origin_port_name", "destination_port_code", "destination_port_name",
			"discharge_port_code", "discharge_port_name", "assigned_berth_name", "dock_id",
			"dock_code", "dock_name", "berth_code", "berth_name", "berth_latitude",
			"berth_longitude", "code_inaportnet", "location_name_inaportnet",
			"start_berth_position", "end_berth_position", "eta", "etb", "etc", "etd", "status", "creation_date",
			"creation_by", "last_updated_date", "last_updated_by",
		},
		SearchColumns: []string{
			"branch_name", "terminal_name", "schedule_code", "vessel_name", "vessel_code",
			"vessel_type", "voyage_number", "pkk_number", "ppk_number", "voyage_type", "agency_name",
			"port_agent", "origin_port_code", "origin_port_name",
			"destination_port_code", "destination_port_name", "discharge_port_code",
			"discharge_port_name", "assigned_berth_name", "dock_code", "dock_name",
			"berth_code", "berth_name", "code_inaportnet", "location_name_inaportnet",
			"start_berth_position", "end_berth_position",
		},
		FilterableColumns: map[string]string{
			"id":                       "id",
			"branch_code":              "branch_code",
			"terminal_code":            "terminal_code",
			"schedule_code":            "schedule_code",
			"vessel_code":              "vessel_code",
			"vessel_type":              "vessel_type",
			"voyage_number":            "voyage_number",
			"pkk_number":               "pkk_number",
			"ppk_number":               "ppk_number",
			"voyage_type":              "voyage_type",
			"origin_port_code":         "origin_port_code",
			"destination_port_code":    "destination_port_code",
			"discharge_port_code":      "discharge_port_code",
			"dock_id":                  "dock_id",
			"dock_code":                "dock_code",
			"berth_code":               "berth_code",
			"code_inaportnet":          "code_inaportnet",
			"location_name_inaportnet": "location_name_inaportnet",
			"start_berth_position":     "start_berth_position",
			"end_berth_position":       "end_berth_position",
			"status":                   "status",
			"eta":                      "eta",
			"etb":                      "etb",
			"etc":                      "etc",
			"etd":                      "etd",
		},
		SortableColumns: map[string]string{
			"id":            "id",
			"schedule_code": "schedule_code",
			"vessel_name":   "vessel_name",
			"vessel_code":   "vessel_code",
			"voyage_number": "voyage_number",
			"pkk_number":    "pkk_number",
			"ppk_number":    "ppk_number",
			"voyage_type":   "voyage_type",
			"agency_name":   "agency_name",
			"eta":           "eta",
			"etb":           "etb",
			"etc":           "etc",
			"etd":           "etd",
			"status":        "status",
			"creation_date": "creation_date",
			"last_updated":  "last_updated_date",
		},
		DefaultSortBy:    "id",
		DefaultSortOrder: "DESC",
		MaxLimit:         100,
		MaxDownloadLimit: 1000,
	}

	var rows []VesselSchedule
	meta, err := helper.GetDynamicPaginatedNativeData(s.db.WithContext(ctx), config, query, &rows)
	if err != nil {
		return nil, meta, err
	}

	res, err := s.withPlans(ctx, rows)
	return res, meta, err
}

func (s *vesselScheduleService) withPlans(ctx context.Context, rows []VesselSchedule) ([]VesselScheduleSearchResponse, error) {
	res := make([]VesselScheduleSearchResponse, len(rows))
	ppkSet := make(map[string]struct{})

	for i, row := range rows {
		res[i] = VesselScheduleSearchResponse{
			VesselSchedule: row,
			Plans:          []VesselSchedulePlanResponse{},
		}
		if row.PPKNumber == nil {
			continue
		}
		ppkNumber := strings.TrimSpace(*row.PPKNumber)
		if ppkNumber == "" {
			continue
		}
		ppkSet[ppkNumber] = struct{}{}
	}

	if len(ppkSet) == 0 {
		return res, nil
	}

	ppkNumbers := make([]string, 0, len(ppkSet))
	for ppkNumber := range ppkSet {
		ppkNumbers = append(ppkNumbers, ppkNumber)
	}

	var plans []VesselSchedulePlanResponse
	if err := s.db.WithContext(ctx).
		Table("plan.post_vessel_plan").
		Select("ppk_number, plan_code, plan_date, activity_code, activity_name, activity_start_date, activity_end_date").
		Where("ppk_number IN ?", ppkNumbers).
		Order("ppk_number ASC, plan_date DESC, id DESC").
		Find(&plans).Error; err != nil {
		return nil, err
	}

	plansByPPK := make(map[string][]VesselSchedulePlanResponse)
	for _, plan := range plans {
		ppkNumber := strings.TrimSpace(plan.PPKNumber)
		if ppkNumber == "" {
			continue
		}
		plansByPPK[ppkNumber] = append(plansByPPK[ppkNumber], plan)
	}

	for i := range res {
		if res[i].PPKNumber == nil {
			continue
		}
		ppkNumber := strings.TrimSpace(*res[i].PPKNumber)
		if grouped, ok := plansByPPK[ppkNumber]; ok {
			res[i].Plans = grouped
		}
	}

	return res, nil
}

func (s *vesselScheduleService) Create(ctx context.Context, schedule *VesselSchedule) error {
	now := time.Now()
	schedule.ID = 0
	schedule.CreationDate = &now
	schedule.LastUpdatedDate = &now

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("LOCK TABLE plan.post_vessel_schedules IN EXCLUSIVE MODE").Error; err != nil {
			return err
		}

		scheduleCode, err := s.nextScheduleCode(tx, now, schedule)
		if err != nil {
			return err
		}
		schedule.ScheduleCode = &scheduleCode

		return tx.Table(schedule.TableName()).Omit("id").Create(schedule).Error
	})
}

func (s *vesselScheduleService) Update(ctx context.Context, scheduleCode string, schedule *VesselSchedule) error {
	now := time.Now()
	schedule.LastUpdatedDate = &now

	result := s.db.WithContext(ctx).
		Table(schedule.TableName()).
		Where("schedule_code = ?", scheduleCode).
		Updates(map[string]interface{}{
			"branch_code":              schedule.BranchCode,
			"terminal_code":            schedule.TerminalCode,
			"branch_name":              schedule.BranchName,
			"terminal_name":            schedule.TerminalName,
			"vessel_name":              schedule.VesselName,
			"vessel_code":              schedule.VesselCode,
			"vessel_type":              schedule.VesselType,
			"vessel_hatch_number":      schedule.VesselHatchNumber,
			"voyage_number":            schedule.VoyageNumber,
			"pkk_number":               schedule.PKKNumber,
			"ppk_number":               schedule.PPKNumber,
			"voyage_type":              schedule.VoyageType,
			"grt":                      schedule.GRT,
			"loa":                      schedule.LOA,
			"agency_name":              schedule.AgencyName,
			"port_agent":               schedule.PortAgent,
			"emergency_contact":        schedule.EmergencyContact,
			"origin_port_code":         schedule.OriginPortCode,
			"origin_port_name":         schedule.OriginPortName,
			"destination_port_code":    schedule.DestinationPortCode,
			"destination_port_name":    schedule.DestinationPortName,
			"discharge_port_code":      schedule.DischargePortCode,
			"discharge_port_name":      schedule.DischargePortName,
			"assigned_berth_name":      schedule.AssignedBerthName,
			"dock_id":                  schedule.DockID,
			"dock_code":                schedule.DockCode,
			"dock_name":                schedule.DockName,
			"berth_code":               schedule.BerthCode,
			"berth_name":               schedule.BerthName,
			"berth_latitude":           schedule.BerthLatitude,
			"berth_longitude":          schedule.BerthLongitude,
			"code_inaportnet":          schedule.CodeInaportnet,
			"location_name_inaportnet": schedule.LocationNameInaportnet,
			"start_berth_position":     schedule.StartBerthPosition,
			"end_berth_position":       schedule.EndBerthPosition,
			"eta":                      schedule.ETA,
			"etb":                      schedule.ETB,
			"etc":                      schedule.ETC,
			"etd":                      schedule.ETD,
			"status":                   schedule.Status,
			"last_updated_date":        schedule.LastUpdatedDate,
			"last_updated_by":          schedule.LastUpdatedBy,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s *vesselScheduleService) UpdateStatus(ctx context.Context, scheduleCode string, status int, updatedBy string) error {
	now := time.Now()
	result := s.db.WithContext(ctx).
		Table((VesselSchedule{}).TableName()).
		Where("schedule_code = ?", scheduleCode).
		Updates(map[string]interface{}{
			"status":            status,
			"last_updated_date": &now,
			"last_updated_by":   updatedBy,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s *vesselScheduleService) Delete(ctx context.Context, id uint64) error {
	result := s.db.WithContext(ctx).Table((VesselSchedule{}).TableName()).Where("id = ?", id).Delete(&VesselSchedule{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s *vesselScheduleService) FindByID(ctx context.Context, id uint64) (*VesselSchedule, error) {
	var row VesselSchedule
	result := s.db.WithContext(ctx).Table((VesselSchedule{}).TableName()).Where("id = ?", id).First(&row)
	if result.Error != nil {
		return nil, result.Error
	}
	return &row, nil
}

func (s *vesselScheduleService) FindByScheduleCode(ctx context.Context, scheduleCode string) (*VesselScheduleDetailResponse, error) {
	var row VesselSchedule
	result := s.db.WithContext(ctx).
		Table((VesselSchedule{}).TableName()).
		Where("schedule_code = ?", scheduleCode).
		First(&row)
	if result.Error != nil {
		return nil, result.Error
	}

	res := &VesselScheduleDetailResponse{
		VesselSchedule: row,
		HatchDetails:   []interface{}{},
	}

	vCode := ""
	if row.VesselCode != nil {
		vCode = *row.VesselCode
	}
	bCode := 0
	if row.BranchCode != nil {
		bCode = *row.BranchCode
	}
	tCode := 0
	if row.TerminalCode != nil {
		tCode = *row.TerminalCode
	}

	if vCode != "" {
		var vessel map[string]interface{}
		if err := s.authDB.WithContext(ctx).
			Table("adm.posm_vessel").
			Where("vessel_code = ? AND branch_code = ? AND terminal_code = ?", vCode, bCode, tCode).
			Limit(1).
			Scan(&vessel).Error; err == nil && len(vessel) > 0 {

			res.Vessel = vessel

			var hatches []map[string]interface{}
			if err := s.authDB.WithContext(ctx).
				Table("adm.posm_vessel_d").
				Where("vessel_code = ? AND branch_code = ? AND terminal_code = ?", vCode, bCode, tCode).
				Order("hatch_code ASC").
				Find(&hatches).Error; err == nil {

				res.HatchDetails = make([]interface{}, len(hatches))
				for i, h := range hatches {
					res.HatchDetails[i] = h
				}
			}
		}
	}

	return res, nil
}

func (s *vesselScheduleService) nextScheduleCode(tx *gorm.DB, now time.Time, schedule *VesselSchedule) (string, error) {
	if schedule.BranchCode == nil {
		return "", fmt.Errorf("branch code is required to generate schedule code")
	}
	if schedule.TerminalCode == nil {
		return "", fmt.Errorf("terminal code is required to generate schedule code")
	}

	period := now.Format("200601")
	prefix := fmt.Sprintf("VS%d%d%s", *schedule.BranchCode, *schedule.TerminalCode, period)
	pattern := fmt.Sprintf("^%s[0-9]{6}$", prefix)

	var lastSequence int
	if err := tx.Raw(`
		SELECT COALESCE(MAX(CAST(SUBSTRING(schedule_code FROM ? FOR 6) AS INTEGER)), 0)
		FROM plan.post_vessel_schedules
		WHERE schedule_code ~ ?
	`, len(prefix)+1, pattern).Scan(&lastSequence).Error; err != nil {
		return "", err
	}

	nextSequence := lastSequence + 1
	if nextSequence > 999999 {
		return "", fmt.Errorf("schedule code sequence for period %s is exhausted", period)
	}

	return fmt.Sprintf("%s%06d", prefix, nextSequence), nil
}

func (s *vesselScheduleService) GetAuthLocation(ctx context.Context, userID uint64) (*VesselScheduleAuthLocation, error) {
	const userQuery = `
		SELECT branch_name, terminal_name
		FROM posm_users
		WHERE id = ?
		LIMIT 1
	`

	var result VesselScheduleAuthLocation
	if err := s.authDB.WithContext(ctx).Raw(userQuery, userID).Scan(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}
