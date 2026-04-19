package user

import (
	"context"
	"errors"
	"time"

	"omniport-api/internal/helper"
)

type UserService interface {
	GetProfile(ctx context.Context, userID uint64) (*UserResponse, error)
	FindAll(ctx context.Context, page int, size int) ([]UserResponse, int64, error)
	GetByID(ctx context.Context, id uint64) (*UserResponse, error)
	CreateUser(ctx context.Context, req *UserRequest, adminEmp string) error
	UpdateUser(ctx context.Context, id uint64, req *UserRequest, adminEmp string) error
	DeleteUser(ctx context.Context, id uint64) error
	Search(ctx context.Context, query helper.PaginationQuery) ([]UserResponse, helper.PaginationMeta, error)
}

type userService struct{ userRepo UserRepository }

func NewUserService(userRepo UserRepository) UserService { return &userService{userRepo: userRepo} }

func (s *userService) GetProfile(ctx context.Context, userID uint64) (*UserResponse, error) {
	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	res := ToResponse(u)
	return &res, nil
}

func (s *userService) FindAll(ctx context.Context, page int, size int) ([]UserResponse, int64, error) {
	limit := size
	offset := (page - 1) * size
	rows, total, err := s.userRepo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var res []UserResponse
	for _, r := range rows {
		res = append(res, ToResponse(&r))
	}
	return res, total, nil
}

func (s *userService) GetByID(ctx context.Context, id uint64) (*UserResponse, error) {
	u, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := ToResponse(u)
	return &res, nil
}

func (s *userService) CreateUser(ctx context.Context, req *UserRequest, adminEmp string) error {
	if _, err := s.userRepo.FindByEmployeeID(ctx, req.EmployeeID); err == nil {
		return errors.New("employee ID already exists")
	}

	if _, err := s.userRepo.FindByEmail(ctx, req.Email); err == nil {
		return errors.New("email already exists")
	}

	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}

	superuser := false
	if req.Superuser != nil {
		superuser = *req.Superuser
	}

	u := &User{
		EmployeeID:       req.EmployeeID,
		FullName:         req.FullName,
		Email:            req.Email,
		JobTitle:         req.JobTitle,
		PhoneNumber:      req.PhoneNumber,
		SubUnitName:      req.SubUnitName,
		BranchCode:       req.BranchCode,
		BranchName:       req.BranchName,
		TerminalCode:     req.TerminalCode,
		CompanyCode:      req.CompanyCode,
		ProfitCenter:     req.ProfitCenter,
		PersonnelArea:    req.PersonnelArea,
		PersonnelSubArea: req.PersonnelSubArea,
		RoleID:           req.RoleID,
		AccessStatus:     req.AccessStatus,
		PasswordHash:     hashedPassword,
		Status:           req.Status,
		Superuser:        superuser,
		CreationDate:     time.Now(),
		CreationBy:       adminEmp,
	}

	return s.userRepo.Create(ctx, u)
}

func (s *userService) UpdateUser(ctx context.Context, id uint64, req *UserRequest, adminEmp string) error {
	u, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("user not found")
	}

	if req.Password != "" {
		hashedPassword, err := helper.HashPassword(req.Password)
		if err != nil {
			return errors.New("failed to hash password")
		}
		u.PasswordHash = hashedPassword
	}

	u.EmployeeID = req.EmployeeID
	u.FullName = req.FullName
	u.Email = req.Email
	u.JobTitle = req.JobTitle
	u.PhoneNumber = req.PhoneNumber
	u.SubUnitName = req.SubUnitName
	u.BranchCode = req.BranchCode
	u.BranchName = req.BranchName
	u.TerminalCode = req.TerminalCode
	u.CompanyCode = req.CompanyCode
	u.ProfitCenter = req.ProfitCenter
	u.PersonnelArea = req.PersonnelArea
	u.PersonnelSubArea = req.PersonnelSubArea
	u.RoleID = req.RoleID
	u.Status = req.Status
	u.AccessStatus = req.AccessStatus
	if req.Superuser != nil {
		u.Superuser = *req.Superuser
	}

	now := time.Now()
	u.LastUpdatedDate = &now
	u.LastUpdatedBy = adminEmp

	return s.userRepo.Update(ctx, id, u)
}

func (s *userService) DeleteUser(ctx context.Context, id uint64) error {
	_, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("user not found")
	}

	return s.userRepo.Delete(ctx, id)
}

func (s *userService) Search(ctx context.Context, query helper.PaginationQuery) ([]UserResponse, helper.PaginationMeta, error) {
	rows, meta, err := s.userRepo.Search(ctx, query)
	if err != nil {
		return nil, meta, err
	}

	var res []UserResponse
	for _, r := range rows {
		res = append(res, ToResponse(&r))
	}
	return res, meta, nil
}

