package branch

import (
	"context"
	"omniport-api/internal/helper"

	"gorm.io/gorm"
)

type BranchRepository interface {
	Create(ctx context.Context, branch *Branch) error
	GetByID(ctx context.Context, id uint64) (*Branch, error)
	GetByCode(ctx context.Context, code string) (*Branch, error)
	Update(ctx context.Context, id uint64, branch *Branch) error
	Delete(ctx context.Context, id uint64) error
	Search(ctx context.Context, param helper.PaginationQuery) ([]Branch, helper.PaginationMeta, error)
	GetStats(ctx context.Context, companyCode string) (*BranchStats, error)
}

type branchRepository struct {
	db *gorm.DB
}

func NewBranchRepository(db *gorm.DB) BranchRepository {
	return &branchRepository{db: db}
}

func (r *branchRepository) Create(ctx context.Context, branch *Branch) error {
	return r.db.WithContext(ctx).Create(branch).Error
}

func (r *branchRepository) GetByID(ctx context.Context, id uint64) (*Branch, error) {
	var b Branch
	err := r.db.WithContext(ctx).First(&b, id).Error
	return &b, err
}

func (r *branchRepository) GetByCode(ctx context.Context, code string) (*Branch, error) {
	var b Branch
	err := r.db.WithContext(ctx).Where("branch_code = ?", code).First(&b).Error
	return &b, err
}

func (r *branchRepository) Update(ctx context.Context, id uint64, branch *Branch) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Updates(branch).Error
}

func (r *branchRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&Branch{}).Error
}

func (r *branchRepository) Search(ctx context.Context, param helper.PaginationQuery) ([]Branch, helper.PaginationMeta, error) {
	config := helper.NativePaginationConfig{
		TableName: "adm.posm_branches",
		SelectColumns: []string{
			"id", "branch_code", "branch_name", "company_code", "company_name",
			"kd_port", "regional_area", "profit_center", "status", "created_by", "created_date",
		},
		SearchColumns: []string{"branch_code::text", "branch_name", "company_name", "kd_port"},
		FilterableColumns: map[string]string{
			"status":       "status",
			"company_code": "company_code",
		},
		SortableColumns: map[string]string{
			"id":          "id",
			"branch_code": "branch_code",
			"branch_name": "branch_name",
		},
		DefaultSortBy:    "id",
		DefaultSortOrder: "DESC",
	}
	var rows []Branch
	meta, err := helper.GetDynamicPaginatedNativeData(r.db.WithContext(ctx), config, param, &rows)
	return rows, meta, err
}

func (r *branchRepository) GetStats(ctx context.Context, companyCode string) (*BranchStats, error) {
	var stats BranchStats
	db := r.db.WithContext(ctx).Model(&Branch{})
	if companyCode != "" {
		db = db.Where("company_code = ?", companyCode)
	}
	db.Count(&stats.TotalBranches)
	
	dbActive := r.db.WithContext(ctx).Model(&Branch{}).Where("status = ?", "1")
	if companyCode != "" {
		dbActive = dbActive.Where("company_code = ?", companyCode)
	}
	dbActive.Count(&stats.ActiveBranches)
	return &stats, nil
}
