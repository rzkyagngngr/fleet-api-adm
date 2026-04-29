package lookup

import (
	"context"
	"strings"

	"omniport-api/internal/modules/administration/dock"
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
	ListDockOptions(ctx context.Context, branchCode int, terminalCode int, q string, limit int) ([]dock.Dock, error)
	ListVesselOptions(ctx context.Context, branchCode int, terminalCode int, q string, limit int) ([]VesselOption, error)
	ListPortOptions(ctx context.Context, q string, limit int) ([]PortOption, error)
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

type VesselOption struct {
	ID                    uint64 `json:"id"`
	VesselCode            string `json:"vessel_code"`
	VesselName            string `json:"vessel_name"`
	VesselType            string `json:"vessel_type"`
	VesselCallSign        string `json:"vessel_call_sign"`
	VesselImo             string `json:"vessel_imo"`
	VesselGrt             string `json:"vessel_grt"`
	VesselLoa             string `json:"vessel_loa"`
	VesselOwnerName       string `json:"vessel_owner_name"`
	VesselShippingRoute   string `json:"vessel_shipping_route"`
	VesselFlag            string `json:"vessel_flag"`
	VesselCountry         string `json:"vessel_country"`
	VesselHatchNumber     int    `json:"vessel_hatch_number"`
	VesselHatchType       string `json:"vessel_hatch_type"`
	VesselOwnershipStatus string `json:"vessel_ownership_status"`
	VesselOperationStatus string `json:"vessel_operation_status"`
	Status                string `json:"status"`
	PortCode              int64  `json:"port_code"`
	BranchCode            int64  `json:"branch_code"`
	TerminalCode          int64  `json:"terminal_code"`
}

type PortOption struct {
	ID          uint64 `json:"id"`
	PortCode    string `json:"port_code"`
	PortName    string `json:"port_name"`
	PortCity    string `json:"port_city"`
	CountryCode string `json:"country_code"`
	Status      string `json:"status"`
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

func (s *lookupService) ListDockOptions(ctx context.Context, branchCode int, terminalCode int, q string, limit int) ([]dock.Dock, error) {
	const baseQuery = `
		SELECT
			id,
			branch_code,
			branch_name,
			terminal_code,
			terminal_name,
			dock_code,
			dock_name,
			dock_type,
			dock_length_m,
			dock_width_m,
			dock_capacity_ton,
			code_inaportnet,
			location_name_inaportnet,
			status
		FROM adm.posm_docks d
		WHERE d.branch_code = ?
		  AND d.terminal_code = ?
		  AND COALESCE(d.status, 0) = 1
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
			UPPER(COALESCE(d.dock_code, '')) LIKE UPPER(?)
			OR UPPER(COALESCE(d.dock_name, '')) LIKE UPPER(?)
			OR UPPER(COALESCE(d.dock_type, '')) LIKE UPPER(?)
			OR EXISTS (
				SELECT 1
				FROM adm.posm_docks_d dd
				WHERE dd.dock_id = d.id
				  AND COALESCE(dd.status, 0) = 1
				  AND (
					UPPER(COALESCE(dd.berth_code, '')) LIKE UPPER(?)
					OR UPPER(COALESCE(dd.berth_name, '')) LIKE UPPER(?)
				  )
			)
		  )
		`
		keyword := "%" + search + "%"
		args = append(args, keyword, keyword, keyword, keyword, keyword)
	}

	query += `
		ORDER BY d.dock_name ASC, d.dock_code ASC
		LIMIT ?
	`
	args = append(args, effectiveLimit)

	var headers []dock.Dock
	if err := s.db.WithContext(ctx).Raw(query, args...).Scan(&headers).Error; err != nil {
		return nil, err
	}
	if len(headers) == 0 {
		return headers, nil
	}

	dockIDs := make([]uint64, 0, len(headers))
	dockByID := make(map[uint64]*dock.Dock, len(headers))
	for i := range headers {
		dockIDs = append(dockIDs, headers[i].ID)
		dockByID[headers[i].ID] = &headers[i]
	}

	var details []dock.DockDetail
	if err := s.db.WithContext(ctx).
		Table((dock.DockDetail{}).TableName()).
		Where("dock_id IN ?", dockIDs).
		Where("COALESCE(status, 0) = 1").
		Order("dock_id ASC, berth_name ASC, berth_code ASC").
		Find(&details).Error; err != nil {
		return nil, err
	}

	for _, detail := range details {
		if header := dockByID[detail.DockID]; header != nil {
			header.Details = append(header.Details, detail)
		}
	}

	return headers, nil
}

func (s *lookupService) ListVesselOptions(ctx context.Context, branchCode int, terminalCode int, q string, limit int) ([]VesselOption, error) {
	const baseQuery = `
		SELECT
			id,
			COALESCE(vessel_code, '') AS vessel_code,
			COALESCE(vessel_name, '') AS vessel_name,
			COALESCE(vessel_type, '') AS vessel_type,
			COALESCE(vessel_call_sign, '') AS vessel_call_sign,
			COALESCE(vessel_imo, '') AS vessel_imo,
			COALESCE(vessel_grt, '') AS vessel_grt,
			COALESCE(vessel_loa, '') AS vessel_loa,
			COALESCE(vessel_owner_name, '') AS vessel_owner_name,
			COALESCE(vessel_shipping_route, '') AS vessel_shipping_route,
			COALESCE(vessel_flag, '') AS vessel_flag,
			COALESCE(vessel_country, '') AS vessel_country,
			vessel_hatch_number,
			COALESCE(vessel_hatch_type, '') AS vessel_hatch_type,
			COALESCE(vessel_ownership_status, '') AS vessel_ownership_status,
			COALESCE(vessel_operation_status, '') AS vessel_operation_status,
			COALESCE(status, '') AS status,
			port_code,
			branch_code,
			terminal_code
		FROM posm_vessel
		WHERE branch_code = ?
		  AND terminal_code = ?
		  AND UPPER(COALESCE(status, 'ACTIVE')) = 'ACTIVE'
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
			UPPER(COALESCE(vessel_code, '')) LIKE UPPER(?)
			OR UPPER(COALESCE(vessel_name, '')) LIKE UPPER(?)
			OR UPPER(COALESCE(vessel_type, '')) LIKE UPPER(?)
			OR UPPER(COALESCE(vessel_call_sign, '')) LIKE UPPER(?)
			OR UPPER(COALESCE(vessel_imo, '')) LIKE UPPER(?)
		  )
		`
		keyword := "%" + search + "%"
		args = append(args, keyword, keyword, keyword, keyword, keyword)
	}

	query += `
		ORDER BY vessel_name ASC, vessel_code ASC
		LIMIT ?
	`
	args = append(args, effectiveLimit)

	var options []VesselOption
	if err := s.db.WithContext(ctx).Raw(query, args...).Scan(&options).Error; err != nil {
		return nil, err
	}

	return options, nil
}

func (s *lookupService) ListPortOptions(ctx context.Context, q string, limit int) ([]PortOption, error) {
	const baseQuery = `
		SELECT
			id,
			COALESCE(port_code, '') AS port_code,
			COALESCE(port_name, '') AS port_name,
			COALESCE(port_city, '') AS port_city,
			COALESCE(country_code, '') AS country_code,
			TRIM(COALESCE(status, '')) AS status
		FROM adm.posm_port
		WHERE TRIM(COALESCE(status, 'A')) = 'A'
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
	args := make([]interface{}, 0, 5)
	if search != "" {
		query += `
		  AND (
			UPPER(COALESCE(port_code, '')) LIKE UPPER(?)
			OR UPPER(COALESCE(port_name, '')) LIKE UPPER(?)
			OR UPPER(COALESCE(port_city, '')) LIKE UPPER(?)
			OR UPPER(COALESCE(country_code, '')) LIKE UPPER(?)
		  )
		`
		keyword := "%" + search + "%"
		args = append(args, keyword, keyword, keyword, keyword)
	}

	query += `
		ORDER BY port_name ASC, port_code ASC
		LIMIT ?
	`
	args = append(args, effectiveLimit)

	var options []PortOption
	if err := s.db.WithContext(ctx).Raw(query, args...).Scan(&options).Error; err != nil {
		return nil, err
	}

	return options, nil
}
