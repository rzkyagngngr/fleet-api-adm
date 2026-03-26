package dto

import (
	"time"
)

type UserRegisterRequest struct {
	EmployeeID             string `json:"employee_id" binding:"required"`
	FullName               string `json:"full_name" binding:"required"`
	Email                  string `json:"email" binding:"required,email"`
	Password               string `json:"password" binding:"required,min=6"`
	JobTitle               string `json:"job_title"`
	PhoneNumber            string `json:"phone_number"`
	SubUnitName            string `json:"sub_unit_name"`
	BranchCode             *int64 `json:"branch_code"`
	BranchName             string `json:"branch_name"`
	TerminalCode           *int64 `json:"terminal_code"`
	CompanyCode            string `json:"company_code"`
	PersonnelArea          string `json:"personnel_area"`
	PersonnelSubArea       string `json:"personnel_sub_area"`
}

type LoginRequest struct {
	EmployeeID string `json:"employee_id" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string           `json:"token"`
	User  UserResponse     `json:"user"`
	Menus []MenuAccessNode `json:"menus"`
}

type UserResponse struct {
	ID                     uint64     `json:"id"`
	EmployeeID             string     `json:"employee_id"`
	FullName               string     `json:"full_name"`
	JobTitle               string     `json:"job_title"`
	Email                  string     `json:"email"`
	PhoneNumber            string     `json:"phone_number"`
	SubUnitName            string     `json:"sub_unit_name"`
	Status                 string     `json:"status"`
	BranchCode             *int64     `json:"branch_code"`
	BranchName             string     `json:"branch_name"`
	TerminalCode           *int64     `json:"terminal_code"`
	Superuser              bool       `json:"superuser"`
	CreationDate           time.Time  `json:"creation_date"`
	LastLoginAt            *time.Time `json:"last_login_at"`
}

type MenuAccessRow struct {
	RolesID        int64   `gorm:"column:roles_id" json:"roles_id"`
	MenuID         int64   `gorm:"column:menu_id" json:"menu_id"`
	MenuCode       string  `gorm:"column:menu_code" json:"menu_code"`
	MenuIcon       *string `gorm:"column:menu_icon" json:"menu_icon"`
	MenuText       string  `gorm:"column:menu_text" json:"menu_text"`
	MenuUrl        *string `gorm:"column:menu_url" json:"menu_url"`
	View           int     `gorm:"column:view" json:"view"`
	Insert         int     `gorm:"column:insert" json:"insert"`
	Update         int     `gorm:"column:update" json:"update"`
	Delete         int     `gorm:"column:delete" json:"delete"`
	MenuLevel      int     `gorm:"column:menu_level" json:"menu_level"`
	ParentMenuID   *int64  `gorm:"column:parent_menu_id" json:"parent_menu_id"`
}

type MenuAccessNode struct {
	MenuAccessRow
	Children []MenuAccessNode `json:"children"`
}
