package dock

import (
	"context"
	"omniport-api/internal/helper"
	"time"

	"gorm.io/gorm"
)

type DockService interface {
	Search(ctx context.Context, query helper.PaginationQuery) ([]Dock, helper.PaginationMeta, error)
	Create(ctx context.Context, dock *Dock) error
	Update(ctx context.Context, id uint64, dock *Dock) error
	Delete(ctx context.Context, id uint64) error
	FindByID(ctx context.Context, id uint64) (*Dock, error)
	GetAuthLocation(ctx context.Context, userID uint64) (*DockAuthLocation, error)
}

type dockService struct {
	db *gorm.DB
}

type DockAuthLocation struct {
	BranchName   string
	TerminalName string
}

func NewDockService(db *gorm.DB) DockService {
	return &dockService{db: db}
}

func (s *dockService) Search(ctx context.Context, query helper.PaginationQuery) ([]Dock, helper.PaginationMeta, error) {
	config := helper.NativePaginationConfig{
		TableName: "adm.posm_docks",
		SelectColumns: []string{
			"id",
			"branch_code",
			"branch_name",
			"terminal_code",
			"terminal_name",
			"dock_code",
			"dock_name",
			"dock_type",
			"dock_length_m",
			"dock_width_m",
			"dock_capacity_ton",
			"code_inaportnet",
			"location_name_inaportnet",
			"status",
			"creation_date",
			"creation_by",
			"last_updated_date",
			"last_updated_by",
		},
		SearchColumns: []string{
			"branch_name",
			"terminal_name",
			"dock_code",
			"dock_name",
			"dock_type",
			"code_inaportnet",
			"location_name_inaportnet",
		},
		FilterableColumns: map[string]string{
			"branch_code":              "branch_code",
			"branch_name":              "branch_name",
			"terminal_code":            "terminal_code",
			"terminal_name":            "terminal_name",
			"dock_code":                "dock_code",
			"dock_name":                "dock_name",
			"dock_type":                "dock_type",
			"status":                   "status",
			"dock_length_m":            "dock_length_m",
			"dock_width_m":             "dock_width_m",
			"dock_capacity_ton":        "dock_capacity_ton",
			"code_inaportnet":          "code_inaportnet",
			"location_name_inaportnet": "location_name_inaportnet",
		},
		SortableColumns: map[string]string{
			"id":                       "id",
			"branch_code":              "branch_code",
			"terminal_code":            "terminal_code",
			"dock_code":                "dock_code",
			"dock_name":                "dock_name",
			"code_inaportnet":          "code_inaportnet",
			"location_name_inaportnet": "location_name_inaportnet",
			"last_updated":             "last_updated_date",
		},
		DefaultSortBy:    "dock_code",
		DefaultSortOrder: "ASC",
		MaxLimit:         100,
		MaxDownloadLimit: 1000,
	}

	var docks []Dock
	meta, err := helper.GetDynamicPaginatedNativeData(s.db.WithContext(ctx), config, query, &docks)
	return docks, meta, err
}

func (s *dockService) Create(ctx context.Context, dock *Dock) error {
	now := time.Now()
	dock.CreationDate = &now
	dock.LastUpdatedDate = &now

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(dock.TableName()).Omit("Details").Create(dock).Error; err != nil {
			return err
		}

		if len(dock.Details) > 0 {
			for i := range dock.Details {
				dock.Details[i].DockID = dock.ID
				dock.Details[i].CreationDate = &now
				dock.Details[i].CreationBy = dock.CreationBy
				dock.Details[i].LastUpdatedDate = &now
				dock.Details[i].LastUpdatedBy = dock.LastUpdatedBy
			}
			if err := tx.Table((DockDetail{}).TableName()).Create(&dock.Details).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *dockService) Update(ctx context.Context, id uint64, dock *Dock) error {
	now := time.Now()
	dock.ID = id
	dock.LastUpdatedDate = &now

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Table(dock.TableName()).
			Where("id = ?", id).
			Updates(map[string]interface{}{
				"branch_code":              dock.BranchCode,
				"branch_name":              dock.BranchName,
				"terminal_code":            dock.TerminalCode,
				"terminal_name":            dock.TerminalName,
				"dock_code":                dock.DockCode,
				"dock_name":                dock.DockName,
				"dock_type":                dock.DockType,
				"dock_length_m":            dock.DockLengthM,
				"dock_width_m":             dock.DockWidthM,
				"dock_capacity_ton":        dock.DockCapacityTon,
				"code_inaportnet":          dock.CodeInaportnet,
				"location_name_inaportnet": dock.LocationNameIna,
				"status":                   dock.Status,
				"last_updated_date":        dock.LastUpdatedDate,
				"last_updated_by":          dock.LastUpdatedBy,
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		if err := tx.Table((DockDetail{}).TableName()).Where("dock_id = ?", id).Delete(&DockDetail{}).Error; err != nil {
			return err
		}

		if len(dock.Details) > 0 {
			for i := range dock.Details {
				dock.Details[i].DockID = id
				dock.Details[i].CreationDate = &now
				dock.Details[i].CreationBy = dock.LastUpdatedBy
				dock.Details[i].LastUpdatedDate = &now
				dock.Details[i].LastUpdatedBy = dock.LastUpdatedBy
			}
			if err := tx.Table((DockDetail{}).TableName()).Create(&dock.Details).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *dockService) Delete(ctx context.Context, id uint64) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table((DockDetail{}).TableName()).Where("dock_id = ?", id).Delete(&DockDetail{}).Error; err != nil {
			return err
		}

		result := tx.Table((Dock{}).TableName()).Where("id = ?", id).Delete(&Dock{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}

func (s *dockService) FindByID(ctx context.Context, id uint64) (*Dock, error) {
	var dock Dock
	result := s.db.WithContext(ctx).
		Table((Dock{}).TableName()).
		Where("id = ?", id).
		First(&dock)
	if result.Error != nil {
		return nil, result.Error
	}

	var details []DockDetail
	if err := s.db.WithContext(ctx).
		Table((DockDetail{}).TableName()).
		Where("dock_id = ?", id).
		Order("id ASC").
		Find(&details).Error; err != nil {
		return nil, err
	}
	dock.Details = details

	return &dock, nil
}

func (s *dockService) GetAuthLocation(ctx context.Context, userID uint64) (*DockAuthLocation, error) {
	const userQuery = `
		SELECT branch_name, terminal_name
		FROM posm_users
		WHERE id = ?
		LIMIT 1
	`

	var result DockAuthLocation
	if err := s.db.WithContext(ctx).Raw(userQuery, userID).Scan(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}
