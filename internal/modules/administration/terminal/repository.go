package terminal

import (
	"context"
	"omniport-api/internal/helper"

	"gorm.io/gorm"
)

type TerminalRepository interface {
	Create(ctx context.Context, terminal *Terminal) error
	GetByID(ctx context.Context, id uint64) (*Terminal, error)
	GetByCode(ctx context.Context, code int64) (*Terminal, error)
	Update(ctx context.Context, id uint64, terminal *Terminal) error
	Delete(ctx context.Context, id uint64) error
	Search(ctx context.Context, param helper.PaginationQuery) ([]Terminal, helper.PaginationMeta, error)
	GetStats(ctx context.Context) (*TerminalStats, error)
}

type terminalRepository struct {
	db *gorm.DB
}

func NewTerminalRepository(db *gorm.DB) TerminalRepository {
	return &terminalRepository{db: db}
}

func (r *terminalRepository) Create(ctx context.Context, terminal *Terminal) error {
	return r.db.WithContext(ctx).Create(terminal).Error
}

func (r *terminalRepository) GetByID(ctx context.Context, id uint64) (*Terminal, error) {
	var t Terminal
	err := r.db.WithContext(ctx).First(&t, id).Error
	return &t, err
}

func (r *terminalRepository) GetByCode(ctx context.Context, code int64) (*Terminal, error) {
	var t Terminal
	err := r.db.WithContext(ctx).Where("terminal_code = ?", code).First(&t).Error
	return &t, err
}

func (r *terminalRepository) Update(ctx context.Context, id uint64, terminal *Terminal) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Updates(terminal).Error
}

func (r *terminalRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&Terminal{}).Error
}

func (r *terminalRepository) Search(ctx context.Context, param helper.PaginationQuery) ([]Terminal, helper.PaginationMeta, error) {
	config := helper.NativePaginationConfig{
		TableName: "adm.posm_terminals",
		SelectColumns: []string{
			"id", "branch_code", "branch_name", "terminal_code", "terminal_name",
			"status", "profit_center", "port_code", "go_live_date", "is_go_live",
			"company_name", "created_by", "created_date",
		},
		SearchColumns: []string{"terminal_code::text", "terminal_name", "branch_name", "company_name"},
		FilterableColumns: map[string]string{
			"status":      "status",
			"branch_code": "branch_code",
			"is_go_live":  "is_go_live",
		},
		SortableColumns: map[string]string{
			"id":            "id",
			"terminal_code": "terminal_code",
			"terminal_name": "terminal_name",
		},
		DefaultSortBy:    "id",
		DefaultSortOrder: "DESC",
	}
	var rows []Terminal
	meta, err := helper.GetDynamicPaginatedNativeData(r.db.WithContext(ctx), config, param, &rows)
	return rows, meta, err
}

func (r *terminalRepository) GetStats(ctx context.Context) (*TerminalStats, error) {
	var stats TerminalStats
	r.db.WithContext(ctx).Model(&Terminal{}).Count(&stats.TotalTerminals)
	r.db.WithContext(ctx).Model(&Terminal{}).Where("is_go_live = ?", "1").Count(&stats.GoLiveTerminals)
	return &stats, nil
}
