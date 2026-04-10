package pelabuhan

import (
	"context"
	"omniport-api/internal/helper"
	"time"

	"gorm.io/gorm"
)

type PortService interface {
	Search(ctx context.Context, query helper.PaginationQuery) ([]Port, helper.PaginationMeta, error)
	Create(ctx context.Context, port *Port) error
	Update(ctx context.Context, id uint64, port *Port) error
	Delete(ctx context.Context, id uint64) error
	FindByID(ctx context.Context, id uint64) (*Port, error)
}

type portService struct {
	db *gorm.DB
}

func NewPortService(db *gorm.DB) PortService {
	return &portService{db: db}
}

func (s *portService) Search(ctx context.Context, query helper.PaginationQuery) ([]Port, helper.PaginationMeta, error) {
	config := helper.NativePaginationConfig{
		TableName: "adm.posm_port",
		SelectColumns: []string{
			"id",
			"port_code",
			"port_name",
			"port_city",
			"country_code",
			"created_by",
			"created_date",
			"last_updated_date",
			"last_updated_by",
			"program_name",
			"status",
		},
		SearchColumns: []string{
			"port_code",
			"port_name",
			"port_city",
		},
		FilterableColumns: map[string]string{
			"port_code":    "port_code",
			"port_name":    "port_name",
			"port_city":    "port_city",
			"country_code": "country_code",
			"status":       "status",
		},
		SortableColumns: map[string]string{
			"id":           "id",
			"port_code":    "port_code",
			"port_name":    "port_name",
			"last_updated": "last_updated_date",
		},
		DefaultSortBy:    "port_code",
		DefaultSortOrder: "ASC",
		MaxLimit:         100,
		MaxDownloadLimit: 1000,
	}

	var ports []Port
	meta, err := helper.GetDynamicPaginatedNativeData(s.db.WithContext(ctx), config, query, &ports)
	return ports, meta, err
}

func (s *portService) Create(ctx context.Context, p *Port) error {
	const query = `
		INSERT INTO adm.posm_port (
			port_code,
			port_name,
			port_city,
			country_code,
			created_by,
			created_date,
			last_updated_date,
			last_updated_by,
			program_name,
			status
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id
	`

	now := time.Now()
	p.CreatedDate = &now
	p.LastUpdatedDate = now
	p.ProgramName = "Master Pelabuhan"

	return s.db.WithContext(ctx).Raw(
		query,
		p.PortCode,
		p.PortName,
		p.PortCity,
		p.CountryCode,
		p.CreatedBy,
		p.CreatedDate,
		p.LastUpdatedDate,
		p.LastUpdatedBy,
		p.ProgramName,
		p.Status,
	).Scan(&p.ID).Error
}

func (s *portService) Update(ctx context.Context, id uint64, p *Port) error {
	const query = `
		UPDATE adm.posm_port
		SET
			port_code = ?,
			port_name = ?,
			port_city = ?,
			country_code = ?,
			last_updated_date = ?,
			last_updated_by = ?,
			program_name = ?,
			status = ?
		WHERE id = ?
	`

	now := time.Now()
	p.LastUpdatedDate = now
	p.ProgramName = "Master Pelabuhan"

	result := s.db.WithContext(ctx).Exec(
		query,
		p.PortCode,
		p.PortName,
		p.PortCity,
		p.CountryCode,
		p.LastUpdatedDate,
		p.LastUpdatedBy,
		p.ProgramName,
		p.Status,
		id,
	)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s *portService) Delete(ctx context.Context, id uint64) error {
	const query = `DELETE FROM adm.posm_port WHERE id = ?`
	result := s.db.WithContext(ctx).Exec(query, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s *portService) FindByID(ctx context.Context, id uint64) (*Port, error) {
	const query = `
		SELECT
			id,
			port_code,
			port_name,
			port_city,
			country_code,
			created_by,
			created_date,
			last_updated_date,
			last_updated_by,
			program_name,
			status
		FROM adm.posm_port
		WHERE id = ?
		LIMIT 1
	`

	var port Port
	result := s.db.WithContext(ctx).Raw(query, id).Scan(&port)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &port, nil
}
