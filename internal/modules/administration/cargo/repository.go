package cargo

import (
	"context"
	"omniport-api/internal/helper"

	"gorm.io/gorm"
)

type CargoRepository interface {
	Create(ctx context.Context, cargo *Cargo) error
	FindByID(ctx context.Context, id uint64) (*Cargo, error)
	Update(ctx context.Context, id uint64, cargo *Cargo) error
	Delete(ctx context.Context, id uint64) error
	Search(ctx context.Context, param helper.PaginationQuery) ([]Cargo, helper.PaginationMeta, error)
	GetStats(ctx context.Context) (*CargoStatsResponse, error)
}

type cargoRepository struct{ db *gorm.DB }

func NewCargoRepository(db *gorm.DB) CargoRepository { return &cargoRepository{db: db} }

func (r *cargoRepository) Create(ctx context.Context, cargo *Cargo) error {
	return r.db.WithContext(ctx).Create(cargo).Error
}

func (r *cargoRepository) FindByID(ctx context.Context, id uint64) (*Cargo, error) {
	var c Cargo
	err := r.db.WithContext(ctx).First(&c, id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *cargoRepository) Update(ctx context.Context, id uint64, c *Cargo) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Updates(c).Error
}

func (r *cargoRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&Cargo{}, id).Error
}

func (r *cargoRepository) Search(ctx context.Context, param helper.PaginationQuery) ([]Cargo, helper.PaginationMeta, error) {
	config := helper.NativePaginationConfig{
		TableName: `(SELECT *, CASE WHEN cargo_imdg_code > 0 THEN '1' ELSE '0' END as is_dangerous_goods FROM posm_cargos) as t`,
		SelectColumns: []string{
			"id", "branch_code", "terminal_code", "cargo_code", "cargo_name",
			"cargo_group", "cargo_commodity", "cargo_characteristic",
			"cargo_imdg_code", "cargo_imdg_description", "is_active",
			"cargo_product_name", "created_date",
		},
		SearchColumns: []string{
			"cargo_code", "cargo_name", "cargo_group", "cargo_product_name",
		},
		FilterableColumns: map[string]string{
			"cargo_code":         "cargo_code",
			"cargo_name":         "cargo_name",
			"cargo_group":        "cargo_group",
			"is_active":          "is_active",
			"is_dangerous_goods": "is_dangerous_goods",
		},
		SortableColumns: map[string]string{
			"id":         "id",
			"item_code":  "cargo_code",
			"item_name":  "cargo_name",
			"created_date": "created_date",
		},
		DefaultSortBy:    "cargo_code",
		DefaultSortOrder: "ASC",
		MaxLimit:         100,
	}

	var rows []Cargo
	meta, err := helper.GetDynamicPaginatedNativeData(r.db.WithContext(ctx), config, param, &rows)
	return rows, meta, err
}

func (r *cargoRepository) GetStats(ctx context.Context) (*CargoStatsResponse, error) {
	var stats CargoStatsResponse

	// Total Cargo Masters
	if err := r.db.WithContext(ctx).Model(&Cargo{}).Count(&stats.TotalCargoMasters).Error; err != nil {
		return nil, err
	}

	// Active Commodities (is_active = '1')
	if err := r.db.WithContext(ctx).Model(&Cargo{}).Where("is_active = ?", "1").Count(&stats.ActiveCommodities).Error; err != nil {
		return nil, err
	}

	// Hazmat Registry (cargo_imdg_code IS NOT NULL AND > 0)
	if err := r.db.WithContext(ctx).Model(&Cargo{}).Where("cargo_imdg_code > ?", 0).Count(&stats.HazmatRegistry).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}
