package company

import (
	"context"
	"omniport-api/internal/helper"

	"gorm.io/gorm"
)

type CompanyRepository interface {
	Create(ctx context.Context, c *Company) error
	Update(ctx context.Context, id uint64, c *Company) error
	Delete(ctx context.Context, id uint64) error
	FindByID(ctx context.Context, id uint64) (*Company, error)
	Search(ctx context.Context, param helper.PaginationQuery) ([]Company, helper.PaginationMeta, error)
}

type companyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) CompanyRepository {
	return &companyRepository{db: db}
}

func (r *companyRepository) Create(ctx context.Context, c *Company) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *companyRepository) Update(ctx context.Context, id uint64, c *Company) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Updates(c).Error
}

func (r *companyRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&Company{}).Error
}

func (r *companyRepository) FindByID(ctx context.Context, id uint64) (*Company, error) {
	var c Company
	err := r.db.WithContext(ctx).First(&c, id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *companyRepository) Search(ctx context.Context, param helper.PaginationQuery) ([]Company, helper.PaginationMeta, error) {
	config := helper.NativePaginationConfig{
		TableName: "adm.posm_companies",
		SelectColumns: []string{
			"id", "company_code", "company_name", "npwp", "address",
			"email", "phone_number", "business_type", "status",
			"created_by", "created_date", "last_updated_by", "last_updated_date", "program_name",
		},
		SearchColumns: []string{
			"company_code", "company_name", "npwp", "email",
		},
		FilterableColumns: map[string]string{
			"status":        "status",
			"business_type": "business_type",
		},
		SortableColumns: map[string]string{
			"id":           "id",
			"company_code": "company_code",
			"company_name": "company_name",
		},
		DefaultSortBy:    "company_name",
		DefaultSortOrder: "ASC",
		MaxLimit:         100,
		MaxDownloadLimit: 1000,
	}

	var rows []Company
	meta, err := helper.GetDynamicPaginatedNativeData(r.db.WithContext(ctx), config, param, &rows)
	return rows, meta, err
}
