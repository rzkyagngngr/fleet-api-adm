package user

import (
	"omniport-api/internal/helper"
	"time"
)

type UserResponse struct {
	ID              uint64     `json:"id"`
	EmployeeID      string     `json:"employee_id"`
	FullName        string     `json:"full_name"`
	JobTitle        string     `json:"job_title"`
	Email           string     `json:"email"`
	PhoneNumber     string     `json:"phone_number"`
	SubUnitName     string     `json:"sub_unit_name"`
	Status          string     `json:"status"`
	BranchCode      *int64     `json:"branch_code"`
	BranchName      string     `json:"branch_name"`
	TerminalCode    *int64     `json:"terminal_code"`
	TerminalName    string     `json:"terminal_name"`
	RoleID          *int64     `json:"role_id"`
	RoleDescription string     `json:"role_description"`
	CompanyCode     string     `json:"company_code"`
	CompanyName     string     `json:"company_name"`
	Superuser       bool       `json:"superuser"`
	CreationDate    time.Time  `json:"creation_date"`
	LastLoginAt     *time.Time `json:"last_login_at"`
}

type UserStatsResponse struct {
	TotalUsers     int64 `json:"total_users"`
	ActiveNow      int64 `json:"active_now"`
	AdminCount     int64 `json:"admin_count"`
	TerminalAccess int64 `json:"terminal_access"`
}

type MenuAccessRow struct {
	RolesID      int64   `gorm:"column:roles_id" json:"roles_id"`
	MenuID       int64   `gorm:"column:menu_id" json:"menu_id"`
	MenuCode     string  `gorm:"column:menu_code" json:"menu_code"`
	MenuIcon     *string `gorm:"column:menu_icon" json:"menu_icon"`
	MenuText     string  `gorm:"column:menu_text" json:"menu_text"`
	MenuURL      *string `gorm:"column:menu_url" json:"menu_url"`
	View         int     `gorm:"column:view" json:"view"`
	Insert       int     `gorm:"column:insert" json:"insert"`
	Update       int     `gorm:"column:update" json:"update"`
	Delete       int     `gorm:"column:delete" json:"delete"`
	MenuLevel    int     `gorm:"column:menu_level" json:"menu_level"`
	ParentMenuID *int64  `gorm:"column:parent_menu_id" json:"parent_menu_id"`
}

type MenuAccessNode struct {
	MenuAccessRow
	Children []MenuAccessNode `json:"children"`
}

type UserRequest struct {
	EmployeeID       string `json:"employee_id" binding:"required"`
	FullName         string `json:"full_name" binding:"required"`
	Email            string `json:"email" binding:"required,email"`
	JobTitle         string `json:"job_title"`
	PhoneNumber      string `json:"phone_number"`
	SubUnitName      string `json:"sub_unit_name"`
	BranchCode       *int64 `json:"branch_code"`
	BranchName       string `json:"branch_name"`
	TerminalCode     *int64 `json:"terminal_code"`
	TerminalName     string `json:"terminal_name"`
	CompanyCode      string `json:"company_code"`
	ProfitCenter     string `json:"profit_center"`
	PersonnelArea    string `json:"personnel_area"`
	PersonnelSubArea string `json:"personnel_sub_area"`
	RoleID           *int64 `json:"role_id"`
	Password         string `json:"password"`
	Status           string `json:"status"`
	Superuser        *bool  `json:"superuser"`
	AccessStatus     *int64 `json:"access_status"`
}

type SearchUsersRequest struct {
	Page     int                  `json:"page"`
	Limit    int                  `json:"limit"`
	Search   string               `json:"search"`
	Filters  map[string]string    `json:"filters"`
	Sort     helper.SortQuery     `json:"sort"`
	Download helper.DownloadQuery `json:"download"`
}

func (r SearchUsersRequest) ToPaginationQuery() helper.PaginationQuery {
	return helper.PaginationQuery{
		Page:     r.Page,
		Limit:    r.Limit,
		Search:   r.Search,
		Filters:  r.Filters,
		Sort:     r.Sort,
		Download: r.Download,
	}
}
