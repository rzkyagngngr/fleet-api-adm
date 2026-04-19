package user

import "time"

type User struct {
	ID                     uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	AccessID               *int64     `gorm:"column:access_id" json:"access_id"`
	RoleID                 *int64     `gorm:"column:role_id" json:"role_id"`
	ApplicationID          *int64     `gorm:"column:application_id" json:"application_id"`
	UserID                 *int64     `gorm:"column:user_id" json:"user_id"`
	EmployeeID             string     `gorm:"column:employee_id" json:"employee_id"`
	FullName               string     `gorm:"column:full_name" json:"full_name"`
	JobTitle               string     `gorm:"column:job_title" json:"job_title"`
	PasswordHash           string     `gorm:"column:password_hash" json:"-"`
	Email                  string     `gorm:"column:email" json:"email"`
	PhoneNumber            string     `gorm:"column:phone_number" json:"phone_number"`
	SubUnitName            string     `gorm:"column:sub_unit_name" json:"sub_unit_name"`
	Status                 string     `gorm:"column:status" json:"status"`
	RoleDescription        string     `gorm:"column:role_description" json:"role_description"`
	ApplicationDescription string     `gorm:"column:application_description" json:"application_description"`
	BranchCode             *int64     `gorm:"column:branch_code" json:"branch_code"`
	BranchName             string     `gorm:"column:branch_name" json:"branch_name"`
	TerminalCode           *int64     `gorm:"column:terminal_code" json:"terminal_code"`
	ProfitCenter           string     `gorm:"column:profit_center" json:"profit_center"`
	ApplicationURL         string     `gorm:"column:application_url" json:"application_url"`
	AccessStatus           *int64     `gorm:"column:access_status" json:"access_status"`
	CompanyCode            string     `gorm:"column:company_code" json:"company_code"`
	AccessUpdatedAt        *time.Time `gorm:"column:access_updated_at" json:"access_updated_at"`
	LastLoginAt            *time.Time `gorm:"column:last_login_at" json:"last_login_at"`
	PersonnelArea          string     `gorm:"column:personnel_area" json:"personnel_area"`
	PersonnelSubArea       string     `gorm:"column:personnel_sub_area" json:"personnel_sub_area"`
	Superuser              bool       `gorm:"column:superuser;default:false" json:"superuser"`
	CreationDate           time.Time  `gorm:"column:creation_date;default:CURRENT_TIMESTAMP" json:"creation_date"`
	CreationBy             string     `gorm:"column:creation_by" json:"creation_by"`
	LastUpdatedDate        *time.Time `gorm:"column:last_updated_date" json:"last_updated_date"`
	LastUpdatedBy          string     `gorm:"column:last_updated_by" json:"last_updated_by"`
	TerminalName           string     `gorm:"column:terminal_name" json:"terminal_name"`
}

func (User) TableName() string { return "posm_users" }

func ToResponse(u *User) UserResponse {
	return UserResponse{
		ID:           u.ID,
		EmployeeID:   u.EmployeeID,
		FullName:     u.FullName,
		JobTitle:     u.JobTitle,
		Email:        u.Email,
		PhoneNumber:  u.PhoneNumber,
		SubUnitName:  u.SubUnitName,
		Status:       u.Status,
		BranchCode:   u.BranchCode,
		BranchName:   u.BranchName,
		TerminalCode: u.TerminalCode,
		TerminalName: u.TerminalName,
		RoleID:       u.RoleID,
		Superuser:    u.Superuser,
		CreationDate: u.CreationDate,
		LastLoginAt:  u.LastLoginAt,
	}
}
