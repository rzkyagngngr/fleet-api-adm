package lookup

import (
	"context"
	"strings"

	"omniport-api/internal/modules/administration/equipment"

	"gorm.io/gorm"
)

type LookupService interface {
	ListEquipmentGroupOptions(ctx context.Context, q string, limit int) ([]equipment.EquipmentGroupOption, error)
	ListCustomerOptions(ctx context.Context, branchCode int, terminalCode int, q string, limit int) ([]equipment.CustomerOption, error)
	ListEquipmentOptions(ctx context.Context, branchCode int, terminalCode int, q string, limit int) ([]EquipmentOption, error)
	ListCargoPackageOptions(ctx context.Context, q string, limit int) ([]CargoPackageOption, error)
	ListCargoUnitOptions(ctx context.Context, q string, limit int) ([]CargoUnitOption, error)
	ListBillingServiceOptions(ctx context.Context, branchCode int, terminalCode int, q string, limit int) ([]BillingServiceOption, error)
	ListCargoOptions(ctx context.Context, q string, limit int) ([]CargoOption, error)
}

type lookupService struct {
	equipmentService equipment.EquipmentService
	db               *gorm.DB
}

type SearchOptionRequest struct {
	Q     string `json:"q"`
	Limit int    `json:"limit"`
}

type CargoPackageOption struct {
	IDRefKey   string `json:"id_ref_key"`
	KetRefData string `json:"ket_ref_data"`
}

type CargoUnitOption struct {
	IDRefKey   string `json:"id_ref_key"`
	KetRefData string `json:"ket_ref_data"`
}

type BillingServiceOption struct {
	IDRefKey   string `json:"id_ref_key"`
	KetRefData string `json:"ket_ref_data"`
}

type CargoOption struct {
	CargoCode string `json:"cargo_code"`
	CargoName string `json:"cargo_name"`
}

type EquipmentOption struct {
	EquipmentCode string `json:"equipment_code"`
	EquipmentName string `json:"equipment_name"`
}

func NewLookupService(db *gorm.DB, equipmentService equipment.EquipmentService) LookupService {
	return &lookupService{
		equipmentService: equipmentService,
		db:               db,
	}
}

func (s *lookupService) ListEquipmentGroupOptions(ctx context.Context, q string, limit int) ([]equipment.EquipmentGroupOption, error) {
	return s.equipmentService.ListEquipmentGroupOptions(ctx, q, limit)
}

func (s *lookupService) ListCustomerOptions(ctx context.Context, branchCode int, terminalCode int, q string, limit int) ([]equipment.CustomerOption, error) {
	return s.equipmentService.ListCustomerOptions(ctx, branchCode, terminalCode, q, limit)
}

func (s *lookupService) ListEquipmentOptions(ctx context.Context, branchCode int, terminalCode int, q string, limit int) ([]EquipmentOption, error) {
	const baseQuery = `
		SELECT
			COALESCE(equipment_code, '') AS equipment_code,
			COALESCE(equipment_name, '') AS equipment_name
		FROM adm.posm_equipments
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
			UPPER(COALESCE(equipment_code, '')) LIKE UPPER(?)
			OR UPPER(COALESCE(equipment_name, '')) LIKE UPPER(?)
		  )
		`
		keyword := "%" + search + "%"
		args = append(args, keyword, keyword)
	}

	query += `
		ORDER BY equipment_name ASC, equipment_code ASC
		LIMIT ?
	`
	args = append(args, effectiveLimit)

	var options []EquipmentOption
	if err := s.db.WithContext(ctx).Raw(query, args...).Scan(&options).Error; err != nil {
		return nil, err
	}

	return options, nil
}

func (s *lookupService) ListCargoPackageOptions(ctx context.Context, q string, limit int) ([]CargoPackageOption, error) {
	const baseQuery = `
		SELECT
			COALESCE(id_ref_key, '') AS id_ref_key,
			COALESCE(ket_ref_data, '') AS ket_ref_data
		FROM adm.posm_reference_d
		WHERE id_ref_file = 'JNSKEMASAN'
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
	args := make([]interface{}, 0, 3)
	if search != "" {
		query += `
		  AND (
			UPPER(COALESCE(id_ref_key, '')) LIKE UPPER(?)
			OR UPPER(COALESCE(ket_ref_data, '')) LIKE UPPER(?)
		  )
		`
		keyword := "%" + search + "%"
		args = append(args, keyword, keyword)
	}

	query += `
		ORDER BY ket_ref_data ASC, id_ref_key ASC
		LIMIT ?
	`
	args = append(args, effectiveLimit)

	var options []CargoPackageOption
	if err := s.db.WithContext(ctx).Raw(query, args...).Scan(&options).Error; err != nil {
		return nil, err
	}

	return options, nil
}

func (s *lookupService) ListCargoUnitOptions(ctx context.Context, q string, limit int) ([]CargoUnitOption, error) {
	const baseQuery = `
		SELECT
			COALESCE(id_ref_key, '') AS id_ref_key,
			COALESCE(ket_ref_data, '') AS ket_ref_data
		FROM adm.posm_reference_d
		WHERE id_ref_file = 'SATDEFINPUT'
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
	args := make([]interface{}, 0, 3)
	if search != "" {
		query += `
		  AND (
			UPPER(COALESCE(id_ref_key, '')) LIKE UPPER(?)
			OR UPPER(COALESCE(ket_ref_data, '')) LIKE UPPER(?)
		  )
		`
		keyword := "%" + search + "%"
		args = append(args, keyword, keyword)
	}

	query += `
		ORDER BY ket_ref_data ASC, id_ref_key ASC
		LIMIT ?
	`
	args = append(args, effectiveLimit)

	var options []CargoUnitOption
	if err := s.db.WithContext(ctx).Raw(query, args...).Scan(&options).Error; err != nil {
		return nil, err
	}

	return options, nil
}

func (s *lookupService) ListBillingServiceOptions(ctx context.Context, branchCode int, terminalCode int, q string, limit int) ([]BillingServiceOption, error) {
	const baseQuery = `
		SELECT
			COALESCE(id_ref_key, '') AS id_ref_key,
			COALESCE(ket_ref_data, '') AS ket_ref_data
		FROM adm.posm_reference_d
		WHERE id_ref_file = 'JASA_TAGIHAN'
		  AND branch_code = ?
		  AND terminal_code = ?
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
	args := []interface{}{branchCode, terminalCode}
	if search != "" {
		query += `
		  AND (
			UPPER(COALESCE(id_ref_key, '')) LIKE UPPER(?)
			OR UPPER(COALESCE(ket_ref_data, '')) LIKE UPPER(?)
		  )
		`
		keyword := "%" + search + "%"
		args = append(args, keyword, keyword)
	}

	query += `
		ORDER BY ket_ref_data ASC, id_ref_key ASC
		LIMIT ?
	`
	args = append(args, effectiveLimit)

	var options []BillingServiceOption
	if err := s.db.WithContext(ctx).Raw(query, args...).Scan(&options).Error; err != nil {
		return nil, err
	}

	return options, nil
}

func (s *lookupService) ListCargoOptions(ctx context.Context, q string, limit int) ([]CargoOption, error) {
	const baseQuery = `
		SELECT
			COALESCE(cargo_code, '') AS cargo_code,
			COALESCE(cargo_name, '') AS cargo_name
		FROM posm_cargos
		WHERE COALESCE(is_active, 'Y') = 'Y'
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
	args := make([]interface{}, 0, 3)
	if search != "" {
		query += `
		  AND (
			UPPER(COALESCE(cargo_code, '')) LIKE UPPER(?)
			OR UPPER(COALESCE(cargo_name, '')) LIKE UPPER(?)
		  )
		`
		keyword := "%" + search + "%"
		args = append(args, keyword, keyword)
	}

	query += `
		ORDER BY cargo_name ASC, cargo_code ASC
		LIMIT ?
	`
	args = append(args, effectiveLimit)

	var options []CargoOption
	if err := s.db.WithContext(ctx).Raw(query, args...).Scan(&options).Error; err != nil {
		return nil, err
	}

	return options, nil
}
