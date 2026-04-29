package equipment

import (
	"context"
	"fmt"
	"omniport-api/internal/helper"
	"strings"
	"time"

	"gorm.io/gorm"
)

type EquipmentService interface {
	Search(ctx context.Context, query helper.PaginationQuery) ([]Equipment, helper.PaginationMeta, error)
	ListCustomerOptions(ctx context.Context, branchCode int, terminalCode int, q string, limit int) ([]CustomerOption, error)
	ListEquipmentGroupOptions(ctx context.Context, q string, limit int) ([]EquipmentGroupOption, error)
	Create(ctx context.Context, equipment *Equipment) error
	Update(ctx context.Context, id uint64, equipment *Equipment) error
	Delete(ctx context.Context, id uint64) error
	FindByID(ctx context.Context, id uint64) (*Equipment, error)
	GetAuthLocation(ctx context.Context, userID uint64) (*EquipmentAuthLocation, error)
}

type equipmentService struct {
	db *gorm.DB
}

type EquipmentAuthLocation struct {
	BranchName   string
	TerminalName string
}

func NewEquipmentService(db *gorm.DB) EquipmentService {
	return &equipmentService{db: db}
}

func (s *equipmentService) Search(ctx context.Context, query helper.PaginationQuery) ([]Equipment, helper.PaginationMeta, error) {
	config := helper.NativePaginationConfig{
		TableName: "adm.posm_equipments",
		SelectColumns: []string{
			"id",
			"branch_code",
			"branch_name",
			"terminal_code",
			"terminal_name",
			"equipment_code",
			"equipment_name",
			"equipment_group",
			"equipment_type",
			"capacity",
			"minimal_load_capacity",
			"max_load_capacity",
			"ownership_status",
			"owner_name",
			"owner_code",
			"start_date",
			"end_date",
			"equipment_condition",
			"status",
			"creation_date",
			"creation_by",
			"last_updated_date",
			"last_updated_by",
		},
		SearchColumns: []string{
			"branch_name",
			"terminal_name",
			"equipment_code",
			"equipment_name",
			"equipment_group",
			"equipment_type",
			"capacity",
			"minimal_load_capacity",
			"max_load_capacity",
			"ownership_status",
			"owner_name",
			"owner_code",
			"equipment_condition",
		},
		FilterableColumns: map[string]string{
			"branch_code":           "branch_code",
			"branch_name":           "branch_name",
			"terminal_code":         "terminal_code",
			"terminal_name":         "terminal_name",
			"equipment_code":        "equipment_code",
			"equipment_name":        "equipment_name",
			"equipment_group":       "equipment_group",
			"equipment_type":        "equipment_type",
			"capacity":              "capacity",
			"minimal_load_capacity": "minimal_load_capacity",
			"max_load_capacity":     "max_load_capacity",
			"ownership_status":      "ownership_status",
			"owner_name":            "owner_name",
			"owner_code":            "owner_code",
			"equipment_condition":   "equipment_condition",
			"status":                "status",
			"start_date":            "start_date",
			"end_date":              "end_date",
		},
		SortableColumns: map[string]string{
			"id":             "id",
			"branch_code":    "branch_code",
			"terminal_code":  "terminal_code",
			"equipment_code": "equipment_code",
			"equipment_name": "equipment_name",
			"equipment_type": "equipment_type",
			"last_updated":   "last_updated_date",
		},
		DefaultSortBy:    "equipment_code",
		DefaultSortOrder: "ASC",
		MaxLimit:         100,
		MaxDownloadLimit: 1000,
	}

	var equipments []Equipment
	meta, err := helper.GetDynamicPaginatedNativeData(s.db.WithContext(ctx), config, query, &equipments)
	return equipments, meta, err
}

func (s *equipmentService) ListCustomerOptions(ctx context.Context, branchCode int, terminalCode int, q string, limit int) ([]CustomerOption, error) {
	const baseQuery = `
		SELECT
			id AS customer_id,
			COALESCE(customer_code, '') AS customer_code,
			COALESCE(customer_name, '') AS customer_name,
			COALESCE(phone_number, '') AS phone_number,
			COALESCE(customer_code, '') AS owner_code,
			COALESCE(customer_name, '') AS owner_name,
			COALESCE(customer_code, '') AS value,
			TRIM(COALESCE(customer_code, '') || ' - ' || COALESCE(customer_name, '')) AS label
		FROM adm.posm_customers
		WHERE branch_code = ?
		  AND terminal_code = ?
		  AND COALESCE(status, 0) = 1
	`

	search := strings.TrimSpace(q)
	effectiveLimit := limit
	if effectiveLimit <= 0 {
		effectiveLimit = 20
	}
	if effectiveLimit > 50 {
		effectiveLimit = 50
	}

	query := baseQuery
	args := []interface{}{branchCode, terminalCode}
	if search != "" {
		query += `
		  AND (
			UPPER(COALESCE(customer_code, '')) LIKE UPPER(?)
			OR UPPER(COALESCE(customer_name, '')) LIKE UPPER(?)
			OR UPPER(COALESCE(phone_number, '')) LIKE UPPER(?)
		  )
		`
		keyword := "%" + search + "%"
		args = append(args, keyword, keyword, keyword)
	}

	query += `
		ORDER BY customer_name ASC, customer_code ASC
		LIMIT ?
	`
	args = append(args, effectiveLimit)

	var options []CustomerOption
	if err := s.db.WithContext(ctx).Raw(query, args...).Scan(&options).Error; err != nil {
		return nil, err
	}

	return options, nil
}

func (s *equipmentService) ListEquipmentGroupOptions(ctx context.Context, q string, limit int) ([]EquipmentGroupOption, error) {
	const baseQuery = `
		SELECT DISTINCT
			COALESCE(id_ref_key, '') AS value,
			COALESCE(id_ref_key, '') AS label
		FROM posm_reference_d
		WHERE id_ref_file = 'GROUPALAT'
		  AND COALESCE(kd_aktif, 'A') = 'A'
	`

	search := strings.TrimSpace(q)
	effectiveLimit := limit
	if effectiveLimit <= 0 {
		effectiveLimit = 20
	}
	if effectiveLimit > 50 {
		effectiveLimit = 50
	}

	query := baseQuery
	args := make([]interface{}, 0, 2)
	if search != "" {
		query += `
		  AND UPPER(COALESCE(id_ref_key, '')) LIKE UPPER(?)
		`
		args = append(args, "%"+search+"%")
	}

	query += `
		ORDER BY label ASC
		LIMIT ?
	`
	args = append(args, effectiveLimit)

	var options []EquipmentGroupOption
	if err := s.db.WithContext(ctx).Raw(query, args...).Scan(&options).Error; err != nil {
		return nil, err
	}

	return options, nil
}

func (s *equipmentService) Create(ctx context.Context, e *Equipment) error {
	const insertQuery = `
		INSERT INTO adm.posm_equipments (
			branch_code,
			branch_name,
			terminal_code,
			terminal_name,
			equipment_code,
			equipment_name,
			equipment_group,
			equipment_type,
			capacity,
			minimal_load_capacity,
			max_load_capacity,
			ownership_status,
			owner_name,
			owner_code,
			start_date,
			end_date,
			equipment_condition,
			status,
			creation_date,
			creation_by,
			last_updated_date,
			last_updated_by
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id
	`

	const updateCodeQuery = `
		UPDATE adm.posm_equipments
		SET equipment_code = ?
		WHERE id = ?
	`

	now := time.Now()
	e.CreationDate = &now
	e.LastUpdatedDate = &now

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Raw(
			insertQuery,
			e.BranchCode,
			e.BranchName,
			e.TerminalCode,
			e.TerminalName,
			nil,
			e.EquipmentName,
			e.EquipmentGroup,
			e.EquipmentType,
			e.Capacity,
			e.MinimalLoadCapacity,
			e.MaxLoadCapacity,
			e.OwnershipStatus,
			e.OwnerName,
			e.OwnerCode,
			e.StartDate,
			e.EndDate,
			e.EquipmentCondition,
			e.Status,
			e.CreationDate,
			e.CreationBy,
			e.LastUpdatedDate,
			e.LastUpdatedBy,
		).Scan(&e.ID).Error; err != nil {
			return err
		}

		equipmentCode := fmt.Sprintf("EQUI%06d", e.ID)
		if err := tx.Exec(updateCodeQuery, equipmentCode, e.ID).Error; err != nil {
			return err
		}

		e.EquipmentCode = &equipmentCode
		return nil
	})
}

func (s *equipmentService) Update(ctx context.Context, id uint64, e *Equipment) error {
	const query = `
		UPDATE adm.posm_equipments
		SET
			branch_code = ?,
			branch_name = ?,
			terminal_code = ?,
			terminal_name = ?,
			equipment_code = ?,
			equipment_name = ?,
			equipment_group = ?,
			equipment_type = ?,
			capacity = ?,
			minimal_load_capacity = ?,
			max_load_capacity = ?,
			ownership_status = ?,
			owner_name = ?,
			owner_code = ?,
			start_date = ?,
			end_date = ?,
			equipment_condition = ?,
			status = ?,
			last_updated_date = ?,
			last_updated_by = ?
		WHERE id = ?
	`

	now := time.Now()
	e.LastUpdatedDate = &now

	result := s.db.WithContext(ctx).Exec(
		query,
		e.BranchCode,
		e.BranchName,
		e.TerminalCode,
		e.TerminalName,
		e.EquipmentCode,
		e.EquipmentName,
		e.EquipmentGroup,
		e.EquipmentType,
		e.Capacity,
		e.MinimalLoadCapacity,
		e.MaxLoadCapacity,
		e.OwnershipStatus,
		e.OwnerName,
		e.OwnerCode,
		e.StartDate,
		e.EndDate,
		e.EquipmentCondition,
		e.Status,
		e.LastUpdatedDate,
		e.LastUpdatedBy,
		id,
	)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s *equipmentService) Delete(ctx context.Context, id uint64) error {
	const query = `DELETE FROM adm.posm_equipments WHERE id = ?`
	result := s.db.WithContext(ctx).Exec(query, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s *equipmentService) FindByID(ctx context.Context, id uint64) (*Equipment, error) {
	const query = `
		SELECT
			id,
			branch_code,
			branch_name,
			terminal_code,
			terminal_name,
			equipment_code,
			equipment_name,
			equipment_group,
			equipment_type,
			capacity,
			minimal_load_capacity,
			max_load_capacity,
			ownership_status,
			owner_name,
			owner_code,
			start_date,
			end_date,
			equipment_condition,
			status,
			creation_date,
			creation_by,
			last_updated_date,
			last_updated_by
		FROM adm.posm_equipments
		WHERE id = ?
		LIMIT 1
	`

	var equipment Equipment
	result := s.db.WithContext(ctx).Raw(query, id).Scan(&equipment)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &equipment, nil
}

func (s *equipmentService) GetAuthLocation(ctx context.Context, userID uint64) (*EquipmentAuthLocation, error) {
	const userQuery = `
		SELECT branch_name, terminal_name
		FROM posm_users
		WHERE id = ?
		LIMIT 1
	`

	var result EquipmentAuthLocation
	if err := s.db.WithContext(ctx).Raw(userQuery, userID).Scan(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}
