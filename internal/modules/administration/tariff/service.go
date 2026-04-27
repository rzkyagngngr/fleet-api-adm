package tariff

import (
	"context"
	"omniport-api/internal/helper"
	"time"

	"gorm.io/gorm"
)

type TariffServiceInterface interface {
	Search(ctx context.Context, query helper.PaginationQuery) ([]Tariff, helper.PaginationMeta, error)
	SearchStatusZero(ctx context.Context, query helper.PaginationQuery) ([]Tariff, helper.PaginationMeta, error)
	Create(ctx context.Context, tariff *Tariff) error
	Update(ctx context.Context, id uint64, tariff *Tariff) error
	UpdateStatus(ctx context.Context, id uint64, status *int, updatedBy *string) error
	Delete(ctx context.Context, id uint64) error
	FindByID(ctx context.Context, id uint64) (*Tariff, error)
	GetAuthLocation(ctx context.Context, userID uint64) (*TariffAuthLocation, error)
}

type service struct {
	db *gorm.DB
}

type TariffAuthLocation struct {
	BranchName   string
	TerminalName string
}

func NewTariffService(db *gorm.DB) TariffServiceInterface {
	return &service{db: db}
}

func (s *service) Search(ctx context.Context, query helper.PaginationQuery) ([]Tariff, helper.PaginationMeta, error) {
	return s.searchWithFixedFilters(ctx, query, nil)
}

func (s *service) SearchStatusZero(ctx context.Context, query helper.PaginationQuery) ([]Tariff, helper.PaginationMeta, error) {
	fixedFilters := map[string]string{
		"status": "0",
	}
	return s.searchWithFixedFilters(ctx, query, fixedFilters)
}

func (s *service) searchWithFixedFilters(ctx context.Context, query helper.PaginationQuery, fixedFilters map[string]string) ([]Tariff, helper.PaginationMeta, error) {
	if query.Filters == nil {
		query.Filters = map[string]string{}
	}
	for key, value := range fixedFilters {
		query.Filters[key] = value
	}

	config := helper.NativePaginationConfig{
		TableName: "adm.posm_tariffs",
		SelectColumns: []string{
			"id",
			"branch_code",
			"branch_name",
			"terminal_code",
			"terminal_name",
			"name_tariff",
			"description",
			"status",
			"agreement_number",
			"start_date",
			"end_date",
			"creation_date",
			"creation_by",
			"last_updated_date",
			"last_updated_by",
		},
		SearchColumns: []string{
			"branch_name",
			"terminal_name",
			"name_tariff",
			"description",
			"agreement_number",
		},
		FilterableColumns: map[string]string{
			"id":               "id",
			"branch_code":      "branch_code",
			"branch_name":      "branch_name",
			"terminal_code":    "terminal_code",
			"terminal_name":    "terminal_name",
			"name_tariff":      "name_tariff",
			"agreement_number": "agreement_number",
			"status":           "status",
			"start_date":       "start_date",
			"end_date":         "end_date",
		},
		SortableColumns: map[string]string{
			"id":            "id",
			"branch_code":   "branch_code",
			"terminal_code": "terminal_code",
			"name_tariff":   "name_tariff",
			"agreement_no":  "agreement_number",
			"status":        "status",
			"start_date":    "start_date",
			"end_date":      "end_date",
			"creation_date": "creation_date",
			"last_updated":  "last_updated_date",
		},
		DefaultSortBy:    "id",
		DefaultSortOrder: "DESC",
		MaxLimit:         100,
		MaxDownloadLimit: 1000,
	}

	var tariffs []Tariff
	meta, err := helper.GetDynamicPaginatedNativeData(s.db.WithContext(ctx), config, query, &tariffs)
	return tariffs, meta, err
}

func (s *service) Create(ctx context.Context, tariff *Tariff) error {
	now := time.Now()
	tariff.CreationDate = &now
	tariff.LastUpdatedDate = &now

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(tariff.TableName()).Omit("Details").Create(tariff).Error; err != nil {
			return err
		}

		if len(tariff.Details) > 0 {
			for i := range tariff.Details {
				tariff.Details[i].IDTariff = tariff.ID
			}
			if err := tx.Table((TariffService{}).TableName()).Create(&tariff.Details).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *service) Update(ctx context.Context, id uint64, tariff *Tariff) error {
	now := time.Now()
	tariff.ID = id
	tariff.LastUpdatedDate = &now

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Table(tariff.TableName()).
			Where("id = ?", id).
			Updates(map[string]interface{}{
				"branch_code":       tariff.BranchCode,
				"branch_name":       tariff.BranchName,
				"terminal_code":     tariff.TerminalCode,
				"terminal_name":     tariff.TerminalName,
				"name_tariff":       tariff.NameTariff,
				"description":       tariff.Description,
				"status":            tariff.Status,
				"agreement_number":  tariff.AgreementNumber,
				"start_date":        tariff.StartDate,
				"end_date":          tariff.EndDate,
				"last_updated_date": tariff.LastUpdatedDate,
				"last_updated_by":   tariff.LastUpdatedBy,
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		if err := tx.Table((TariffService{}).TableName()).Where("id_tariff = ?", id).Delete(&TariffService{}).Error; err != nil {
			return err
		}

		if len(tariff.Details) > 0 {
			for i := range tariff.Details {
				tariff.Details[i].IDTariff = id
			}
			if err := tx.Table((TariffService{}).TableName()).Create(&tariff.Details).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *service) UpdateStatus(ctx context.Context, id uint64, status *int, updatedBy *string) error {
	now := time.Now()

	result := s.db.WithContext(ctx).
		Table((Tariff{}).TableName()).
		Where("id = ?", id).
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

func (s *service) Delete(ctx context.Context, id uint64) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table((TariffService{}).TableName()).Where("id_tariff = ?", id).Delete(&TariffService{}).Error; err != nil {
			return err
		}

		result := tx.Table((Tariff{}).TableName()).Where("id = ?", id).Delete(&Tariff{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return nil
	})
}

func (s *service) FindByID(ctx context.Context, id uint64) (*Tariff, error) {
	var tariff Tariff
	result := s.db.WithContext(ctx).Table((Tariff{}).TableName()).Where("id = ?", id).First(&tariff)
	if result.Error != nil {
		return nil, result.Error
	}

	var details []TariffService
	if err := s.db.WithContext(ctx).
		Table((TariffService{}).TableName()).
		Where("id_tariff = ?", id).
		Order("sequence_no ASC").
		Find(&details).Error; err != nil {
		return nil, err
	}
	tariff.Details = details

	return &tariff, nil
}

func (s *service) GetAuthLocation(ctx context.Context, userID uint64) (*TariffAuthLocation, error) {
	const userQuery = `
		SELECT branch_name, terminal_name
		FROM posm_users
		WHERE id = ?
		LIMIT 1
	`

	var result TariffAuthLocation
	if err := s.db.WithContext(ctx).Raw(userQuery, userID).Scan(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}
