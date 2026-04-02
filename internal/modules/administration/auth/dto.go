package auth

import "omniport-api/internal/modules/administration/user"

type UserRegisterRequest struct {
	EmployeeID       string `json:"employee_id" binding:"required"`
	FullName         string `json:"full_name" binding:"required"`
	Email            string `json:"email" binding:"required,email"`
	Password         string `json:"password" binding:"required,min=6"`
	JobTitle         string `json:"job_title"`
	PhoneNumber      string `json:"phone_number"`
	SubUnitName      string `json:"sub_unit_name"`
	BranchCode       *int64 `json:"branch_code"`
	BranchName       string `json:"branch_name"`
	TerminalCode     *int64 `json:"terminal_code"`
	CompanyCode      string `json:"company_code"`
	PersonnelArea    string `json:"personnel_area"`
	PersonnelSubArea string `json:"personnel_sub_area"`
}

type LoginRequest struct {
	EmployeeID string `json:"employee_id" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type ChangeTerminalRequest struct {
	BranchCode   *int64 `json:"branch_code" binding:"required"`
	TerminalCode *int64 `json:"terminal_code" binding:"required"`
}

type AuthResponse struct {
	Token string                `json:"token"`
	User  user.UserResponse     `json:"user"`
	Menus []user.MenuAccessNode `json:"menus"`
}
