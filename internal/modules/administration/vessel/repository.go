package vessel

import (
	"context"
	"omniport-api/internal/helper"

	"gorm.io/gorm"
)

type VesselRepository interface {
	Create(ctx context.Context, v *Vessel) error
	Update(ctx context.Context, id uint64, v *Vessel) error
	Delete(ctx context.Context, id uint64) error
	FindByID(ctx context.Context, id uint64) (*Vessel, error)
	Search(ctx context.Context, param helper.PaginationQuery) ([]Vessel, helper.PaginationMeta, error)
	GetStats(ctx context.Context) (*VesselStatsResponse, error)
}

type vesselRepository struct {
	db *gorm.DB
}

func NewVesselRepository(db *gorm.DB) VesselRepository {
	return &vesselRepository{db: db}
}

func (r *vesselRepository) Create(ctx context.Context, v *Vessel) error {
	return r.db.WithContext(ctx).Create(v).Error
}

func (r *vesselRepository) Update(ctx context.Context, id uint64, v *Vessel) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Updates(v).Error
}

func (r *vesselRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&Vessel{}).Error
}

func (r *vesselRepository) FindByID(ctx context.Context, id uint64) (*Vessel, error) {
	var v Vessel
	err := r.db.WithContext(ctx).First(&v, id).Error
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *vesselRepository) Search(ctx context.Context, param helper.PaginationQuery) ([]Vessel, helper.PaginationMeta, error) {
	config := helper.NativePaginationConfig{
		TableName: "posm_vessel",
		SelectColumns: []string{
			"id", "vessel_code", "vessel_name", "vessel_type", "vessel_call_sign",
			"vessel_imo", "vessel_grt", "vessel_loa", "vessel_owner_name",
			"vessel_shipping_route", "vessel_flag", "vessel_country", "vessel_year_made",
			"vessel_hatch_number", "vessel_hatch_type", "vessel_ownership_status",
			"vessel_operation_status", "status", "remark", "port_code", "terminal_code",
			"creation_date",
		},
		SearchColumns: []string{
			"vessel_code", "vessel_name", "vessel_imo", "vessel_call_sign", "vessel_owner_name",
		},
		FilterableColumns: map[string]string{
			"vessel_type":  "vessel_type",
			"vessel_flag":  "vessel_flag",
			"status":       "status",
			"port_code":     "port_code",
			"terminal_code": "terminal_code",
		},
		SortableColumns: map[string]string{
			"id":          "id",
			"vessel_code": "vessel_code",
			"vessel_name": "vessel_name",
			"vessel_grt":  "vessel_grt",
			"vessel_loa":  "vessel_loa",
		},
		DefaultSortBy:    "vessel_name",
		DefaultSortOrder: "ASC",
		MaxLimit:         100,
		MaxDownloadLimit: 1000,
	}

	var rows []Vessel
	meta, err := helper.GetDynamicPaginatedNativeData(r.db.WithContext(ctx), config, param, &rows)
	return rows, meta, err
}

func (r *vesselRepository) GetStats(ctx context.Context) (*VesselStatsResponse, error) {
	var stats VesselStatsResponse

	// Total Fleet
	if err := r.db.WithContext(ctx).Model(&Vessel{}).Count(&stats.TotalFleet).Error; err != nil {
		return nil, err
	}

	// Active Vessels
	if err := r.db.WithContext(ctx).Model(&Vessel{}).Where("status IN ?", []string{"ACTIVE", "IN TRANSIT"}).Count(&stats.ActiveVessels).Error; err != nil {
		return nil, err
	}

	// Maintenance
	if err := r.db.WithContext(ctx).Model(&Vessel{}).Where("status = ? OR vessel_operation_status = ?", "MAINTENANCE", "MAINTENANCE").Count(&stats.Maintenance).Error; err != nil {
		return nil, err
	}

	// Deactivated
	if err := r.db.WithContext(ctx).Model(&Vessel{}).Where("status = ?", "INACTIVE").Count(&stats.Deactivated).Error; err != nil {
		return nil, err
	}

	// Counts by Type
	r.db.WithContext(ctx).Model(&Vessel{}).Where("vessel_type = ?", "GENERAL").Count(&stats.CargoCount)
	r.db.WithContext(ctx).Model(&Vessel{}).Where("vessel_type = ?", "TANKER").Count(&stats.TankerCount)
	r.db.WithContext(ctx).Model(&Vessel{}).Where("vessel_type = ?", "CONTAINER").Count(&stats.ContainerCount)
	r.db.WithContext(ctx).Model(&Vessel{}).Where("vessel_type NOT IN ?", []string{"GENERAL", "TANKER", "CONTAINER"}).Count(&stats.OtherCount)

	return &stats, nil
}
