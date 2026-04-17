package menu

import (
	"context"

	"omniport-api/internal/helper"

	"gorm.io/gorm"
)

type MenuService interface {
	Create(ctx context.Context, menu *Menu) error
	FindAll(ctx context.Context) ([]Menu, error)
	Search(ctx context.Context, query helper.PaginationQuery) ([]Menu, helper.PaginationMeta, error)
	FindByID(ctx context.Context, id uint64) (*Menu, error)
	Update(ctx context.Context, id uint64, menu *Menu) error
	Delete(ctx context.Context, id uint64) error
}

type menuService struct{ db *gorm.DB }

func NewMenuService(db *gorm.DB) MenuService { return &menuService{db: db} }

func (s *menuService) Create(ctx context.Context, m *Menu) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		const insertMenuQuery = `
			INSERT INTO adm.posm_menus (
				menu_code,
				menu_text,
				menu_desc,
				menu_url,
				menu_level,
				menu_order,
				parent_menu_id,
				menu_icon,
				application_id,
				menu_header_id,
				menu_status,
				creation_by,
				creation_date,
				last_updated_by,
				last_updated_date
			) VALUES (
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
			)
			RETURNING id
		`

		if err := tx.Raw(
			insertMenuQuery,
			m.MenuCode,
			m.MenuText,
			m.MenuDesc,
			m.MenuURL,
			m.MenuLevel,
			m.MenuOrder,
			m.ParentMenuID,
			m.MenuIcon,
			m.ApplicationID,
			m.MenuHeaderID,
			m.MenuStatus,
			m.CreationBy,
			m.CreationDate,
			m.LastUpdatedBy,
			m.LastUpdatedDate,
		).Scan(&m.ID).Error; err != nil {
			return err
		}

		// Automate access record creation for all roles
		var roles []uint64
		if err := tx.Raw("SELECT hak_akses_id FROM adm.posm_roles").Scan(&roles).Error; err != nil {
			return err
		}

		for _, roleID := range roles {
			const insertAccessQuery = `
				INSERT INTO adm.posm_access (
					roles_id,
					menu_id,
					menu_text,
					menu_url,
					status,
					application_id,
					parent_menu_id,
					can_insert,
					can_update,
					can_delete,
					menu_order,
					menu_icon
				) VALUES (
					?, ?, ?, ?, 0, ?, ?, 0, 0, 0, ?, ?
				)
			`
			if err := tx.Exec(
				insertAccessQuery,
				roleID,
				m.ID,
				m.MenuText,
				m.MenuURL,
				m.ApplicationID,
				m.ParentMenuID,
				m.MenuOrder,
				m.MenuIcon,
			).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *menuService) FindAll(ctx context.Context) ([]Menu, error) {
	const query = `
		SELECT
			id,
			menu_code,
			menu_text,
			menu_desc,
			menu_url,
			menu_level,
			menu_order,
			parent_menu_id,
			menu_icon,
			application_id,
			menu_header_id,
			menu_status,
			creation_by,
			creation_date,
			last_updated_by,
			last_updated_date
		FROM adm.posm_menus
		ORDER BY menu_level ASC, id ASC
	`

	var menus []Menu
	err := s.db.WithContext(ctx).Raw(query).Scan(&menus).Error
	return menus, err
}

func (s *menuService) Search(ctx context.Context, query helper.PaginationQuery) ([]Menu, helper.PaginationMeta, error) {
	config := helper.NativePaginationConfig{
		TableName: "adm.posm_menus",
		SelectColumns: []string{
			"id",
			"menu_code",
			"menu_text",
			"menu_desc",
			"menu_url",
			"menu_level",
			"menu_order",
			"parent_menu_id",
			"menu_icon",
			"application_id",
			"menu_header_id",
			"menu_status",
			"creation_by",
			"creation_date",
			"last_updated_by",
			"last_updated_date",
		},
		SearchColumns: []string{
			"menu_code",
			"menu_text",
			"menu_desc",
			"menu_url",
			"menu_icon",
		},
		FilterableColumns: map[string]string{
			"menu_code":      "menu_code",
			"menu_text":      "menu_text",
			"menu_desc":      "menu_desc",
			"menu_url":       "menu_url",
			"menu_level":     "menu_level",
			"menu_order":     "menu_order",
			"menu_status":    "menu_status",
			"parent_menu_id": "parent_menu_id",
			"menu_icon":      "menu_icon",
		},
		SortableColumns: map[string]string{
			"id":          "id",
			"menu_code":   "menu_code",
			"menu_text":   "menu_text",
			"menu_level":  "menu_level",
			"menu_order":  "menu_order",
			"menu_status": "menu_status",
		},
		DefaultSortBy:    "menu_level",
		DefaultSortOrder: "ASC",
		MaxLimit:         100,
		MaxDownloadLimit: 1000,
	}

	var menus []Menu
	meta, err := helper.GetDynamicPaginatedNativeData(s.db.WithContext(ctx), config, query, &menus)
	return menus, meta, err
}

func (s *menuService) FindByID(ctx context.Context, id uint64) (*Menu, error) {
	const query = `
		SELECT
			id,
			menu_code,
			menu_text,
			menu_desc,
			menu_url,
			menu_level,
			menu_order,
			parent_menu_id,
			menu_icon,
			application_id,
			menu_header_id,
			menu_status,
			creation_by,
			creation_date,
			last_updated_by,
			last_updated_date
		FROM adm.posm_menus
		WHERE id = ?
		LIMIT 1
	`

	var menu Menu
	result := s.db.WithContext(ctx).Raw(query, id).Scan(&menu)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &menu, nil
}

func (s *menuService) Update(ctx context.Context, id uint64, m *Menu) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		const query = `
			UPDATE adm.posm_menus
			SET
				menu_code = ?,
				menu_text = ?,
				menu_desc = ?,
				menu_url = ?,
				menu_level = ?,
				menu_order = ?,
				parent_menu_id = ?,
				menu_icon = ?,
				application_id = ?,
				menu_header_id = ?,
				menu_status = ?,
				creation_by = ?,
				creation_date = ?,
				last_updated_by = ?,
				last_updated_date = ?
			WHERE id = ?
		`

		result := tx.Exec(
			query,
			m.MenuCode,
			m.MenuText,
			m.MenuDesc,
			m.MenuURL,
			m.MenuLevel,
			m.MenuOrder,
			m.ParentMenuID,
			m.MenuIcon,
			m.ApplicationID,
			m.MenuHeaderID,
			m.MenuStatus,
			m.CreationBy,
			m.CreationDate,
			m.LastUpdatedBy,
			m.LastUpdatedDate,
			id,
		)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		// Sync missing access records for all roles
		syncAccessQuery := `
			INSERT INTO adm.posm_access (
				roles_id, menu_id, menu_text, menu_url, status, 
				application_id, parent_menu_id, can_insert, can_update, can_delete, 
				menu_order, menu_icon
			)
			SELECT 
				r.hak_akses_id, ?, ?, ?, 0, 
				?, ?, 0, 0, 0, ?, ?
			FROM adm.posm_roles r
			WHERE NOT EXISTS (
				SELECT 1 FROM adm.posm_access a 
				WHERE a.menu_id = ? AND a.roles_id = r.hak_akses_id
			)
		`
		if err := tx.Exec(
			syncAccessQuery,
			id, m.MenuText, m.MenuURL,
			m.ApplicationID, m.ParentMenuID, m.MenuOrder, m.MenuIcon,
			id,
		).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *menuService) Delete(ctx context.Context, id uint64) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Clean up associated access records first
		if err := tx.Exec("DELETE FROM adm.posm_access WHERE menu_id = ?", id).Error; err != nil {
			return err
		}

		const query = `DELETE FROM adm.posm_menus WHERE id = ?`
		result := tx.Exec(query, id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}
