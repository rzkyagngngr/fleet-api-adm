package warehouse

import (
	"context"
	"omniport-api/internal/helper"
	"time"

	"gorm.io/gorm"
)

type WarehouseService interface {
	Search(ctx context.Context, query helper.PaginationQuery) ([]Warehouse, helper.PaginationMeta, error)
	Create(ctx context.Context, warehouse *Warehouse) error
	Update(ctx context.Context, id uint64, warehouse *Warehouse) error
	Delete(ctx context.Context, id uint64) error
	FindByID(ctx context.Context, id uint64) (*Warehouse, error)
	GetAuthLocation(ctx context.Context, userID uint64) (*WarehouseAuthLocation, error)
}

type warehouseService struct {
	db *gorm.DB
}

type WarehouseAuthLocation struct {
	BranchName   string
	TerminalName string
}

func NewWarehouseService(db *gorm.DB) WarehouseService {
	return &warehouseService{db: db}
}

func (s *warehouseService) Search(ctx context.Context, query helper.PaginationQuery) ([]Warehouse, helper.PaginationMeta, error) {
	config := helper.NativePaginationConfig{
		TableName: "adm.posm_warehouses",
		SelectColumns: []string{
			"id",
			"branch_code",
			"branch_name",
			"terminal_code",
			"terminal_name",
			"warehouse_code",
			"warehouse_name",
			"warehouse_type",
			"warehouse_capacity",
			"status",
			"creation_date",
			"creation_by",
			"last_updated_date",
			"last_updated_by",
		},
		SearchColumns: []string{
			"branch_name",
			"terminal_name",
			"warehouse_code",
			"warehouse_name",
			"warehouse_type",
			"warehouse_capacity",
		},
		FilterableColumns: map[string]string{
			"branch_code":        "branch_code",
			"branch_name":        "branch_name",
			"terminal_code":      "terminal_code",
			"terminal_name":      "terminal_name",
			"warehouse_code":     "warehouse_code",
			"warehouse_name":     "warehouse_name",
			"warehouse_type":     "warehouse_type",
			"warehouse_capacity": "warehouse_capacity",
			"status":             "status",
		},
		SortableColumns: map[string]string{
			"id":             "id",
			"branch_code":    "branch_code",
			"terminal_code":  "terminal_code",
			"warehouse_code": "warehouse_code",
			"warehouse_name": "warehouse_name",
			"last_updated":   "last_updated_date",
		},
		DefaultSortBy:    "warehouse_code",
		DefaultSortOrder: "ASC",
		MaxLimit:         100,
		MaxDownloadLimit: 1000,
	}

	var warehouses []Warehouse
	meta, err := helper.GetDynamicPaginatedNativeData(s.db.WithContext(ctx), config, query, &warehouses)
	return warehouses, meta, err
}

func (s *warehouseService) Create(ctx context.Context, warehouse *Warehouse) error {
	now := time.Now()
	warehouse.CreationDate = &now
	warehouse.LastUpdatedDate = &now

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(warehouse.TableName()).Omit("Details").Create(warehouse).Error; err != nil {
			return err
		}

		if len(warehouse.Details) > 0 {
			for i := range warehouse.Details {
				warehouse.Details[i].WarehouseID = warehouse.ID
				warehouse.Details[i].CreationDate = &now
				warehouse.Details[i].CreationBy = warehouse.CreationBy
				warehouse.Details[i].LastUpdatedDate = &now
				warehouse.Details[i].LastUpdatedBy = warehouse.LastUpdatedBy
			}
			if err := tx.Table((WarehouseDetail{}).TableName()).Create(&warehouse.Details).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *warehouseService) Update(ctx context.Context, id uint64, warehouse *Warehouse) error {
	now := time.Now()
	warehouse.ID = id
	warehouse.LastUpdatedDate = &now

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Table(warehouse.TableName()).
			Where("id = ?", id).
			Updates(map[string]interface{}{
				"branch_code":        warehouse.BranchCode,
				"branch_name":        warehouse.BranchName,
				"terminal_code":      warehouse.TerminalCode,
				"terminal_name":      warehouse.TerminalName,
				"warehouse_code":     warehouse.WarehouseCode,
				"warehouse_name":     warehouse.WarehouseName,
				"warehouse_type":     warehouse.WarehouseType,
				"warehouse_capacity": warehouse.WarehouseCapacity,
				"status":             warehouse.Status,
				"last_updated_date":  warehouse.LastUpdatedDate,
				"last_updated_by":    warehouse.LastUpdatedBy,
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		if err := tx.Table((WarehouseDetail{}).TableName()).Where("warehouse_id = ?", id).Delete(&WarehouseDetail{}).Error; err != nil {
			return err
		}

		if len(warehouse.Details) > 0 {
			for i := range warehouse.Details {
				warehouse.Details[i].WarehouseID = id
				warehouse.Details[i].CreationDate = &now
				warehouse.Details[i].CreationBy = warehouse.LastUpdatedBy
				warehouse.Details[i].LastUpdatedDate = &now
				warehouse.Details[i].LastUpdatedBy = warehouse.LastUpdatedBy
			}
			if err := tx.Table((WarehouseDetail{}).TableName()).Create(&warehouse.Details).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *warehouseService) Delete(ctx context.Context, id uint64) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table((WarehouseDetail{}).TableName()).Where("warehouse_id = ?", id).Delete(&WarehouseDetail{}).Error; err != nil {
			return err
		}

		result := tx.Table((Warehouse{}).TableName()).Where("id = ?", id).Delete(&Warehouse{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}

func (s *warehouseService) FindByID(ctx context.Context, id uint64) (*Warehouse, error) {
	var warehouse Warehouse
	result := s.db.WithContext(ctx).
		Table((Warehouse{}).TableName()).
		Where("id = ?", id).
		First(&warehouse)
	if result.Error != nil {
		return nil, result.Error
	}

	var details []WarehouseDetail
	if err := s.db.WithContext(ctx).
		Table((WarehouseDetail{}).TableName()).
		Where("warehouse_id = ?", id).
		Order("id ASC").
		Find(&details).Error; err != nil {
		return nil, err
	}
	warehouse.Details = details

	return &warehouse, nil
}

func (s *warehouseService) GetAuthLocation(ctx context.Context, userID uint64) (*WarehouseAuthLocation, error) {
	const userQuery = `
		SELECT branch_name, terminal_name
		FROM posm_users
		WHERE id = ?
		LIMIT 1
	`

	var result WarehouseAuthLocation
	if err := s.db.WithContext(ctx).Raw(userQuery, userID).Scan(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}
