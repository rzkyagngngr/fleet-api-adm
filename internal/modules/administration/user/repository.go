package user

import (
	"context"
	"omniport-api/internal/helper"
	"omniport-api/internal/modules/administration/branch"
	"omniport-api/internal/modules/administration/terminal"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id uint64) (*User, error)
	FindByEmployeeID(ctx context.Context, employeeID string) (*User, error)
	GetUserMenusByRole(ctx context.Context, roleID *int64) ([]MenuAccessRow, error)
	FindAll(ctx context.Context, limit int, offset int) ([]User, int64, error)
	Update(ctx context.Context, id uint64, user *User) error
	Delete(ctx context.Context, id uint64) error
	Search(ctx context.Context, param helper.PaginationQuery) ([]User, helper.PaginationMeta, error)
	GetStats(ctx context.Context) (*UserStatsResponse, error)
	GetUserLocations(ctx context.Context, userID uint64) ([]UserLocationBranch, []UserLocationTerminal, error)
}

type userRepository struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) UserRepository { return &userRepository{db: db} }


func (r *userRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}
func (r *userRepository) FindByID(ctx context.Context, id uint64) (*User, error) {
	var u User
	err := r.db.WithContext(ctx).
		Table("adm.posm_users").
		Select("adm.posm_users.*, c.company_name").
		Joins("LEFT JOIN adm.posm_companies c ON adm.posm_users.company_code = c.company_code").
		Preload("Branches").
		Preload("Terminals").
		Where("adm.posm_users.id = ?", id).
		First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}
func (r *userRepository) FindByEmployeeID(ctx context.Context, employeeID string) (*User, error) {
	var u User
	err := r.db.WithContext(ctx).Where("employee_id = ?", employeeID).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}
func (r *userRepository) GetUserMenusByRole(ctx context.Context, roleID *int64) ([]MenuAccessRow, error) {
	var menus []MenuAccessRow
	if roleID == nil {
		return menus, nil
	}
	err := r.db.WithContext(ctx).Raw("SELECT roles_id, menu_id, menu_code, menu_icon, menu_text, menu_url, view, insert, update, delete, menu_level, parent_menu_id FROM vw_access_login WHERE roles_id = ?", *roleID).Scan(&menus).Error
	if err != nil {
		return nil, err
	}
	return menus, nil
}

func (r *userRepository) Create(ctx context.Context, u *User) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Branches", "Terminals").Create(u).Error; err != nil {
			return err
		}
		return r.syncAssociations(tx, u.ID, u.Branches, u.Terminals)
	})
}

func (r *userRepository) Update(ctx context.Context, id uint64, u *User) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{
			"employee_id": u.EmployeeID, "full_name": u.FullName, "email": u.Email,
			"job_title": u.JobTitle, "phone_number": u.PhoneNumber, "sub_unit_name": u.SubUnitName,
			"company_code": u.CompanyCode, "role_id": u.RoleID, "status": u.Status,
			"superuser": u.Superuser, "last_updated_date": u.LastUpdatedDate, "last_updated_by": u.LastUpdatedBy,
			"branch_code": "", "branch_name": "", "terminal_code": "", "terminal_name": "",
			"application_id": u.ApplicationID, "access_status": u.AccessStatus,
		}
		if u.PasswordHash != "" { updates["password_hash"] = u.PasswordHash }

		if err := tx.Model(&User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return err
		}
		return r.syncAssociations(tx, id, u.Branches, u.Terminals)
	})
}

func (r *userRepository) syncAssociations(tx *gorm.DB, userID uint64, branches []branch.Branch, terminals []terminal.Terminal) error {
	// Sync Branches
	if err := tx.Exec("DELETE FROM adm.posm_user_branches WHERE user_id = ?", userID).Error; err != nil {
		return err
	}
	for _, b := range branches {
		if err := tx.Exec("INSERT INTO adm.posm_user_branches (user_id, branch_code) VALUES (?, ?)", userID, b.BranchCode).Error; err != nil {
			return err
		}
	}

	// Sync Terminals
	if err := tx.Exec("DELETE FROM adm.posm_user_terminals WHERE user_id = ?", userID).Error; err != nil {
		return err
	}
	for _, t := range terminals {
		if err := tx.Exec("INSERT INTO adm.posm_user_terminals (user_id, terminal_code) VALUES (?, ?)", userID, t.TerminalCode).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&User{}).Error
}

func (r *userRepository) FindAll(ctx context.Context, limit int, offset int) ([]User, int64, error) {
	var rows []User
	var total int64
	q := r.db.WithContext(ctx).Model(&User{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Limit(limit).Offset(offset).Order("id DESC").Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *userRepository) Search(ctx context.Context, param helper.PaginationQuery) ([]User, helper.PaginationMeta, error) {
	config := helper.NativePaginationConfig{
		TableName: "adm.posm_users u LEFT JOIN adm.posm_roles r ON u.role_id = r.hak_akses_id LEFT JOIN adm.posm_companies c ON u.company_code = c.company_code",
		SelectColumns: []string{
			"u.id", "u.access_id", "u.role_id", "u.application_id", "u.user_id", "u.employee_id",
			"u.full_name", "u.job_title", "u.email", "u.phone_number", "u.sub_unit_name",
			"u.status", "u.branch_code", "u.branch_name", "u.terminal_code", "u.terminal_name",
			"u.profit_center", "u.access_status", "u.company_code", "u.superuser",
			"u.creation_date", "u.last_login_at", "r.hak_akses_nama as role_description", "c.company_name",
		},
		SearchColumns: []string{"u.employee_id", "u.full_name", "u.email", "u.job_title", "u.phone_number"},
		FilterableColumns: map[string]string{
			"employee_id": "u.employee_id", "full_name": "u.full_name", "status": "u.status",
			"role_id": "u.role_id", "branch_code": "u.branch_code", "terminal_code": "u.terminal_code",
			"company_code": "u.company_code", "superuser": "u.superuser",
		},
		SortableColumns: map[string]string{"id": "u.id", "employee_id": "u.employee_id", "full_name": "u.full_name", "creation_date": "u.creation_date"},
		DefaultSortBy: "id", DefaultSortOrder: "DESC", MaxLimit: 100, MaxDownloadLimit: 1000,
	}

	var rows []User
	meta, err := helper.GetDynamicPaginatedNativeData(r.db.WithContext(ctx), config, param, &rows)
	if err != nil || len(rows) == 0 {
		return rows, meta, err
	}

	// Optimized: Collect all IDs and preload in one go instead of N+1
	userIDs := make([]uint64, len(rows))
	for i := range rows { userIDs[i] = rows[i].ID }

	var preloadedUsers []User
	r.db.WithContext(ctx).Preload("Branches").Preload("Terminals").Where("id IN ?", userIDs).Find(&preloadedUsers)

	// Map back to original order from native pagination
	userMap := make(map[uint64]User)
	for _, u := range preloadedUsers { userMap[u.ID] = u }
	for i := range rows {
		if preloaded, ok := userMap[rows[i].ID]; ok {
			rows[i].Branches = preloaded.Branches
			rows[i].Terminals = preloaded.Terminals
		}
	}

	return rows, meta, err
}

func (r *userRepository) GetStats(ctx context.Context) (*UserStatsResponse, error) {
	var stats UserStatsResponse

	// Total Users
	if err := r.db.WithContext(ctx).Model(&User{}).Count(&stats.TotalUsers).Error; err != nil {
		return nil, err
	}

	// Active Now (status = '1')
	if err := r.db.WithContext(ctx).Model(&User{}).Where("status = ?", "1").Count(&stats.ActiveNow).Error; err != nil {
		return nil, err
	}

	// Admin Count (superuser = true)
	if err := r.db.WithContext(ctx).Model(&User{}).Where("superuser = ?", true).Count(&stats.AdminCount).Error; err != nil {
		return nil, err
	}

	// Terminal Access (terminal_code IS NOT NULL)
	if err := r.db.WithContext(ctx).Model(&User{}).Where("terminal_code IS NOT NULL").Count(&stats.TerminalAccess).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

func (r *userRepository) GetUserLocations(ctx context.Context, userID uint64) ([]UserLocationBranch, []UserLocationTerminal, error) {
	var branches []UserLocationBranch
	var terminals []UserLocationTerminal

	err := r.db.WithContext(ctx).Raw(`
		SELECT ub.branch_code, b.branch_name 
		FROM adm.posm_user_branches ub
		JOIN adm.posm_branches b ON ub.branch_code = b.branch_code
		WHERE ub.user_id = ?
	`, userID).Scan(&branches).Error
	if err != nil {
		return nil, nil, err
	}

	err = r.db.WithContext(ctx).Raw(`
		SELECT t.terminal_code, t.terminal_name, t.branch_code
		FROM adm.posm_user_terminals ut
		JOIN adm.posm_terminals t ON (ut.terminal_code = t.terminal_code OR ut.terminal_code = CAST(t.id AS TEXT))
		WHERE ut.user_id = ?
	`, userID).Scan(&terminals).Error
	if err != nil {
		return nil, nil, err
	}

	return branches, terminals, nil
}
