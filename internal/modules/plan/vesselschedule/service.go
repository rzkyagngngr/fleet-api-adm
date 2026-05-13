package vesselschedule

import (
	"context"
	"fmt"
	"log/slog"
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
	InitChatGroup(ctx context.Context, scheduleCode string, actor string) (*VesselSchedule, error)
	Delete(ctx context.Context, id uint64) error
	FindByID(ctx context.Context, id uint64) (*VesselSchedule, error)
	FindByScheduleCode(ctx context.Context, scheduleCode string) (*VesselScheduleDetailResponse, error)
	FindByTopicID(ctx context.Context, topicID int64) (*VesselSchedule, error)
	GetAuthLocation(ctx context.Context, userID uint64) (*VesselScheduleAuthLocation, error)
}

type vesselScheduleService struct {
	db                   *gorm.DB
	authDB               *gorm.DB
	chatInit             ScheduleChatInitializer
	telegramParentChatID int64
}

type VesselScheduleAuthLocation struct {
	BranchName   string
	TerminalName string
}

type ChatInitSettings struct {
	Initializer          ScheduleChatInitializer
	TelegramParentChatID int64
}

func NewVesselScheduleService(db *gorm.DB, opts ...interface{}) VesselScheduleService {
	locationDB := db
	var chatInit ScheduleChatInitializer
	var telegramParentChatID int64

	for _, opt := range opts {
		switch v := opt.(type) {
		case *gorm.DB:
			if v != nil {
				locationDB = v
			}
		case ScheduleChatInitializer:
			chatInit = v
		case ChatInitSettings:
			chatInit = v.Initializer
			telegramParentChatID = v.TelegramParentChatID
		}
	}

	return &vesselScheduleService{
		db:                   db,
		authDB:               locationDB,
		chatInit:             chatInit,
		telegramParentChatID: telegramParentChatID,
	}
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
			"start_berth_position", "end_berth_position", "eta", "etb", "etc", "etd",
			"telegram_topic_id", "telegram_topic_name",
			"status", "creation_date",
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
			"branch_name":              "branch_name",
			"terminal_name":            "terminal_name",
			"schedule_code":            "schedule_code",
			"vessel_name":              "vessel_name",
			"vessel_code":              "vessel_code",
			"vessel_type":              "vessel_type",
			"voyage_number":            "voyage_number",
			"pkk_number":               "pkk_number",
			"ppk_number":               "ppk_number",
			"voyage_type":              "voyage_type",
			"agency_name":              "agency_name",
			"port_agent":               "port_agent",
			"origin_port_code":         "origin_port_code",
			"origin_port_name":         "origin_port_name",
			"destination_port_code":    "destination_port_code",
			"destination_port_name":    "destination_port_name",
			"discharge_port_code":      "discharge_port_code",
			"discharge_port_name":      "discharge_port_name",
			"assigned_berth_name":      "assigned_berth_name",
			"dock_id":                  "dock_id",
			"dock_code":                "dock_code",
			"dock_name":                "dock_name",
			"berth_code":               "berth_code",
			"berth_name":               "berth_name",
			"code_inaportnet":          "code_inaportnet",
			"location_name_inaportnet": "location_name_inaportnet",
			"start_berth_position":     "start_berth_position",
			"end_berth_position":       "end_berth_position",
			"eta":                      "eta",
			"etb":                      "etb",
			"etc":                      "etc",
			"etd":                      "etd",
		},
		SortableColumns: map[string]string{
			"id":            "id",
			"schedule_code": "schedule_code",
		},
		DefaultSortBy:    "id",
		DefaultSortOrder: "DESC",
	}

	// Handle semantic status mapping
	if statusVal, ok := query.Filters["status"]; ok {
		switch strings.ToLower(statusVal) {
		case "active":
			query.Filters["status"] = "active_mapping"
		case "suspended":
			query.Filters["status"] = "2"
		}
	}

	var rows []VesselSchedule
	db := s.db.WithContext(ctx)
	if query.Filters["status"] == "active_mapping" {
		delete(query.Filters, "status")
		db = db.Where("status IN ?", []int{0, 1})
	}

	meta, err := helper.GetDynamicPaginatedNativeData(db, config, query, &rows)
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

		query := `
			INSERT INTO plan.post_vessel_schedules (
				branch_code, terminal_code, branch_name, terminal_name, schedule_code, 
				vessel_name, vessel_code, vessel_type, vessel_hatch_number, voyage_number, 
				pkk_number, ppk_number, voyage_type, grt, loa, agency_name, port_agent, 
				emergency_contact, origin_port_code, origin_port_name, destination_port_code, 
				destination_port_name, discharge_port_code, discharge_port_name, 
				assigned_berth_name, dock_id, dock_code, dock_name, berth_code, berth_name, 
				berth_latitude, berth_longitude, code_inaportnet, location_name_inaportnet, 
				start_berth_position, end_berth_position, eta, etb, etc, etd, status, 
				creation_date, creation_by, last_updated_date, last_updated_by
			) VALUES (
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
			) RETURNING id`

		return tx.Raw(query,
			schedule.BranchCode, schedule.TerminalCode, schedule.BranchName, schedule.TerminalName, schedule.ScheduleCode,
			schedule.VesselName, schedule.VesselCode, schedule.VesselType, schedule.VesselHatchNumber, schedule.VoyageNumber,
			schedule.PKKNumber, schedule.PPKNumber, schedule.VoyageType, schedule.GRT, schedule.LOA, schedule.AgencyName, schedule.PortAgent,
			schedule.EmergencyContact, schedule.OriginPortCode, schedule.OriginPortName, schedule.DestinationPortCode,
			schedule.DestinationPortName, schedule.DischargePortCode, schedule.DischargePortName,
			schedule.AssignedBerthName, schedule.DockID, schedule.DockCode, schedule.DockName, schedule.BerthCode, schedule.BerthName,
			schedule.BerthLatitude, schedule.BerthLongitude, schedule.CodeInaportnet, schedule.LocationNameInaportnet,
			schedule.StartBerthPosition, schedule.EndBerthPosition, schedule.ETA, schedule.ETB, schedule.ETC, schedule.ETD, schedule.Status,
			schedule.CreationDate, schedule.CreationBy, schedule.LastUpdatedDate, schedule.LastUpdatedBy,
		).Scan(&schedule.ID).Error
	})
}

func (s *vesselScheduleService) Update(ctx context.Context, scheduleCode string, schedule *VesselSchedule) error {
	now := time.Now()
	query := `
		UPDATE plan.post_vessel_schedules SET 
			branch_code = ?, terminal_code = ?, vessel_name = ?, vessel_code = ?, 
			vessel_type = ?, vessel_hatch_number = ?, voyage_number = ?, pkk_number = ?, 
			ppk_number = ?, voyage_type = ?, grt = ?, loa = ?, agency_name = ?, 
			port_agent = ?, emergency_contact = ?, origin_port_code = ?, origin_port_name = ?, 
			destination_port_code = ?, destination_port_name = ?, discharge_port_code = ?, 
			discharge_port_name = ?, assigned_berth_name = ?, dock_id = ?, dock_code = ?, 
			dock_name = ?, berth_code = ?, berth_name = ?, berth_latitude = ?, 
			berth_longitude = ?, code_inaportnet = ?, location_name_inaportnet = ?, 
			start_berth_position = ?, end_berth_position = ?, eta = ?, etb = ?, 
			etc = ?, etd = ?, status = ?, last_updated_date = ?, last_updated_by = ? 
		WHERE schedule_code = ?`

	result := s.db.WithContext(ctx).Exec(query,
		schedule.BranchCode, schedule.TerminalCode, schedule.VesselName, schedule.VesselCode,
		schedule.VesselType, schedule.VesselHatchNumber, schedule.VoyageNumber, schedule.PKKNumber,
		schedule.PPKNumber, schedule.VoyageType, schedule.GRT, schedule.LOA, schedule.AgencyName,
		schedule.PortAgent, schedule.EmergencyContact, schedule.OriginPortCode, schedule.OriginPortName,
		schedule.DestinationPortCode, schedule.DestinationPortName, schedule.DischargePortCode,
		schedule.DischargePortName, schedule.AssignedBerthName, schedule.DockID, schedule.DockCode,
		schedule.DockName, schedule.BerthCode, schedule.BerthName, schedule.BerthLatitude,
		schedule.BerthLongitude, schedule.CodeInaportnet, schedule.LocationNameInaportnet,
		schedule.StartBerthPosition, schedule.EndBerthPosition, schedule.ETA, schedule.ETB,
		schedule.ETC, schedule.ETD, schedule.Status, &now, schedule.LastUpdatedBy,
		scheduleCode,
	)

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
	query := `UPDATE plan.post_vessel_schedules SET status = ?, last_updated_date = ?, last_updated_by = ? WHERE schedule_code = ?`
	result := s.db.WithContext(ctx).Exec(query, status, &now, updatedBy, scheduleCode)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s *vesselScheduleService) Delete(ctx context.Context, id uint64) error {
	query := `DELETE FROM plan.post_vessel_schedules WHERE id = ?`
	result := s.db.WithContext(ctx).Exec(query, id)
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
	query := `SELECT * FROM plan.post_vessel_schedules WHERE id = ? LIMIT 1`
	if err := s.db.WithContext(ctx).Raw(query, id).Scan(&row).Error; err != nil {
		return nil, err
	}
	if row.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &row, nil
}

func (s *vesselScheduleService) FindByScheduleCode(ctx context.Context, scheduleCode string) (*VesselScheduleDetailResponse, error) {
	var row VesselSchedule
	query := `SELECT * FROM plan.post_vessel_schedules WHERE schedule_code = ? LIMIT 1`
	if err := s.db.WithContext(ctx).Raw(query, scheduleCode).Scan(&row).Error; err != nil {
		return nil, err
	}
	if row.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	res := &VesselScheduleDetailResponse{
		VesselSchedule: row,
		HatchDetails:   []interface{}{},
	}

	if row.VesselCode != nil && *row.VesselCode != "" {
		var vessel map[string]interface{}
		vesselQuery := `SELECT * FROM adm.posm_vessel WHERE vessel_code = ? AND branch_code = ? AND terminal_code = ? LIMIT 1`
		if err := s.authDB.WithContext(ctx).Raw(vesselQuery, *row.VesselCode, row.BranchCode, row.TerminalCode).Scan(&vessel).Error; err == nil && len(vessel) > 0 {
			res.Vessel = vessel
			var hatches []map[string]interface{}
			hatchQuery := `SELECT * FROM adm.posm_vessel_d WHERE vessel_code = ? AND branch_code = ? AND terminal_code = ? ORDER BY hatch_code ASC`
			if err := s.authDB.WithContext(ctx).Raw(hatchQuery, *row.VesselCode, row.BranchCode, row.TerminalCode).Scan(&hatches).Error; err == nil {
				res.HatchDetails = make([]interface{}, len(hatches))
				for i, h := range hatches {
					res.HatchDetails[i] = h
				}
			}
		}
	}

	return res, nil
}

func (s *vesselScheduleService) FindByTopicID(ctx context.Context, topicID int64) (*VesselSchedule, error) {
	var row VesselSchedule
	sTopicID := fmt.Sprintf("%d", topicID)
	query := `SELECT * FROM plan.post_vessel_schedules WHERE telegram_topic_id = ? LIMIT 1`
	if err := s.db.WithContext(ctx).Raw(query, sTopicID).Scan(&row).Error; err != nil {
		return nil, err
	}
	if row.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &row, nil
}

func (s *vesselScheduleService) GetAuthLocation(ctx context.Context, userID uint64) (*VesselScheduleAuthLocation, error) {
	var result VesselScheduleAuthLocation
	query := `SELECT branch_name, terminal_name FROM adm.posm_users WHERE id = ? LIMIT 1`
	if err := s.authDB.WithContext(ctx).Raw(query, userID).Scan(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *vesselScheduleService) InitChatGroup(ctx context.Context, scheduleCode string, actor string) (*VesselSchedule, error) {
	trimmedCode := strings.TrimSpace(scheduleCode)
	if trimmedCode == "" {
		return nil, fmt.Errorf("schedule_code is required")
	}

	var current VesselSchedule
	findQuery := `SELECT * FROM plan.post_vessel_schedules WHERE schedule_code = ? LIMIT 1`
	if err := s.db.WithContext(ctx).Raw(findQuery, trimmedCode).Scan(&current).Error; err != nil {
		return nil, err
	}

	if err := s.ensureScheduleChatInitialized(ctx, &current, actor); err != nil {
		return nil, err
	}

	if err := s.db.WithContext(ctx).Raw(findQuery, trimmedCode).Scan(&current).Error; err != nil {
		return nil, err
	}
	return &current, nil
}

func (s *vesselScheduleService) ensureScheduleChatInitialized(ctx context.Context, current *VesselSchedule, actor string) error {
	if s.chatInit == nil || current == nil || current.ID == 0 {
		return nil
	}
	if s.telegramParentChatID == 0 {
		return fmt.Errorf("TELEGRAM_PARENT_CHAT_ID is not configured")
	}
	if current.TelegramTopicID != nil && *current.TelegramTopicID != "" {
		return nil
	}

	topicID, topicName, err := s.chatInit.InitScheduleChat(ctx, ChatInitRequest{
		ScheduleID:     current.ID,
		TelegramChatID: s.telegramParentChatID,
		TopicName:      buildScheduleTopicName(current),
		Actor:          actor,
	})
	if err != nil {
		return err
	}

	if topicID != nil && *topicID != 0 {
		now := time.Now()
		sTopicID := fmt.Sprintf("%d", *topicID)
		updateQuery := `UPDATE plan.post_vessel_schedules SET telegram_topic_id = ?, telegram_topic_name = ?, last_updated_date = ?, last_updated_by = ? WHERE id = ?`
		return s.db.WithContext(ctx).Exec(updateQuery, sTopicID, topicName, &now, actor, current.ID).Error
	}

	return nil
}

func buildScheduleTopicName(s *VesselSchedule) string {
	vessel := "VESSEL"
	if s.VesselName != nil && strings.TrimSpace(*s.VesselName) != "" {
		vessel = strings.TrimSpace(*s.VesselName)
	}
	eta := time.Now().Format("02 Jan")
	if s.ETA != nil {
		eta = s.ETA.Format("02 Jan")
	}
	return fmt.Sprintf("%s - %s", eta, vessel)
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

func actorFromSchedule(schedule *VesselSchedule) string {
	if schedule == nil {
		return "SYSTEM"
	}
	if schedule.LastUpdatedBy != nil && strings.TrimSpace(*schedule.LastUpdatedBy) != "" {
		return strings.TrimSpace(*schedule.LastUpdatedBy)
	}
	if schedule.CreationBy != nil && strings.TrimSpace(*schedule.CreationBy) != "" {
		return strings.TrimSpace(*schedule.CreationBy)
	}
	return "SYSTEM"
}

func scheduleIDValue(schedule *VesselSchedule) uint64 {
	if schedule == nil {
		return 0
	}
	return schedule.ID
}

func scheduleCodeValue(schedule *VesselSchedule) string {
	if schedule == nil || schedule.ScheduleCode == nil {
		return ""
	}
	return strings.TrimSpace(*schedule.ScheduleCode)
}

func (s *vesselScheduleService) tryEnsureScheduleChatInitialized(ctx context.Context, schedule *VesselSchedule, actor string) {
	if err := s.ensureScheduleChatInitialized(ctx, schedule, actor); err != nil {
		slog.Warn("failed to initialize telegram schedule chat", "error", err, "schedule_id", scheduleIDValue(schedule), "schedule_code", scheduleCodeValue(schedule))
	}
}
