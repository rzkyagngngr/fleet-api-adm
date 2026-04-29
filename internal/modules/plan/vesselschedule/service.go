package vesselschedule

import (
	"context"
	"omniport-api/internal/helper"
	"time"

	"gorm.io/gorm"
)

type VesselScheduleService interface {
	Search(ctx context.Context, query helper.PaginationQuery) ([]VesselSchedule, helper.PaginationMeta, error)
	Create(ctx context.Context, schedule *VesselSchedule) error
	Update(ctx context.Context, id uint64, schedule *VesselSchedule) error
	Delete(ctx context.Context, id uint64) error
	FindByID(ctx context.Context, id uint64) (*VesselSchedule, error)
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

func (s *vesselScheduleService) Search(ctx context.Context, query helper.PaginationQuery) ([]VesselSchedule, helper.PaginationMeta, error) {
	config := helper.NativePaginationConfig{
		TableName: "plan.post_vessel_schedules",
		SelectColumns: []string{
			"id", "branch_code", "terminal_code", "branch_name", "terminal_name",
			"vessel_name", "vessel_code", "vessel_type", "voyage_number", "voyage_type", "grt", "loa",
			"schedule_code", "pkk_number",

			"agency_name", "port_agent", "emergency_contact", "origin_port_code",
			"origin_port_name", "destination_port_code", "destination_port_name",
			"discharge_port_code", "discharge_port_name", "assigned_berth_name", "dock_id",
			"dock_code", "dock_name", "berth_code", "berth_name", "berth_position",
			"position_range", "eta", "etb", "etc", "etd", "status", "creation_date",
			"creation_by", "last_updated_date", "last_updated_by",
		},
		SearchColumns: []string{
			"branch_name", "terminal_name", "vessel_name", "vessel_code", "vessel_type",
			"voyage_number", "schedule_code", "pkk_number", "voyage_type", "agency_name", "port_agent", "origin_port_code", "origin_port_name",

			"destination_port_code", "destination_port_name", "discharge_port_code",
			"discharge_port_name", "assigned_berth_name", "dock_code", "dock_name",
			"berth_code", "berth_name", "berth_position", "position_range",
		},
		FilterableColumns: map[string]string{
			"id":                    "id",
			"branch_code":           "branch_code",
			"terminal_code":         "terminal_code",
			"vessel_code":           "vessel_code",
			"vessel_type":           "vessel_type",
			"voyage_number":         "voyage_number",
			"voyage_type":           "voyage_type",
			"schedule_code":         "schedule_code",
			"pkk_number":            "pkk_number",

			"origin_port_code":      "origin_port_code",
			"destination_port_code": "destination_port_code",
			"discharge_port_code":   "discharge_port_code",
			"dock_id":               "dock_id",
			"dock_code":             "dock_code",
			"berth_code":            "berth_code",
			"status":                "status",
			"eta":                   "eta",
			"etb":                   "etb",
			"etc":                   "etc",
			"etd":                   "etd",
		},
		SortableColumns: map[string]string{
			"id":            "id",
			"vessel_name":   "vessel_name",
			"vessel_code":   "vessel_code",
			"voyage_number": "voyage_number",

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
	return rows, meta, err
}

func (s *vesselScheduleService) Create(ctx context.Context, schedule *VesselSchedule) error {
	now := time.Now()
	schedule.CreationDate = &now
	schedule.LastUpdatedDate = &now
	return s.db.WithContext(ctx).Table(schedule.TableName()).Create(schedule).Error
}

func (s *vesselScheduleService) Update(ctx context.Context, id uint64, schedule *VesselSchedule) error {
	now := time.Now()
	schedule.LastUpdatedDate = &now

	result := s.db.WithContext(ctx).
		Table(schedule.TableName()).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"branch_code":           schedule.BranchCode,
			"terminal_code":         schedule.TerminalCode,
			"branch_name":           schedule.BranchName,
			"terminal_name":         schedule.TerminalName,
			"vessel_name":           schedule.VesselName,
			"vessel_code":           schedule.VesselCode,
			"vessel_type":           schedule.VesselType,
			"voyage_number":         schedule.VoyageNumber,
			"voyage_type":           schedule.VoyageType,
			"schedule_code":         schedule.ScheduleCode,
			"pkk_number":            schedule.PKKNumber,

			"grt":                   schedule.GRT,
			"loa":                   schedule.LOA,
			"agency_name":           schedule.AgencyName,
			"port_agent":            schedule.PortAgent,
			"emergency_contact":     schedule.EmergencyContact,
			"origin_port_code":      schedule.OriginPortCode,
			"origin_port_name":      schedule.OriginPortName,
			"destination_port_code": schedule.DestinationPortCode,
			"destination_port_name": schedule.DestinationPortName,
			"discharge_port_code":   schedule.DischargePortCode,
			"discharge_port_name":   schedule.DischargePortName,
			"assigned_berth_name":   schedule.AssignedBerthName,
			"dock_id":               schedule.DockID,
			"dock_code":             schedule.DockCode,
			"dock_name":             schedule.DockName,
			"berth_code":            schedule.BerthCode,
			"berth_name":            schedule.BerthName,
			"berth_position":        schedule.BerthPosition,
			"position_range":        schedule.PositionRange,
			"eta":                   schedule.ETA,
			"etb":                   schedule.ETB,
			"etc":                   schedule.ETC,
			"etd":                   schedule.ETD,
			"status":                schedule.Status,
			"last_updated_date":     schedule.LastUpdatedDate,
			"last_updated_by":       schedule.LastUpdatedBy,
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
