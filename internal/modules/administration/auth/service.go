package auth

import (
	"context"
	"errors"
	"time"

	"omniport-api/internal/helper"
	"omniport-api/internal/modules/administration/user"

	"gorm.io/gorm"
)

type AuthService interface {
	Register(ctx context.Context, req *UserRegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error)
	ChangeTerminal(ctx context.Context, userID uint64, req *ChangeTerminalRequest) (*AuthResponse, error)
}

type authService struct {
	userRepo user.UserRepository
	jwtUtil  *helper.JWTUtil
}

func NewAuthService(userRepo user.UserRepository, jwtUtil *helper.JWTUtil) AuthService {
	return &authService{userRepo: userRepo, jwtUtil: jwtUtil}
}

func (s *authService) Register(ctx context.Context, req *UserRegisterRequest) (*AuthResponse, error) {
	if _, err := s.userRepo.FindByEmail(ctx, req.Email); err == nil {
		return nil, errors.New("email already registered")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if _, err := s.userRepo.FindByEmployeeID(ctx, req.EmployeeID); err == nil {
		return nil, errors.New("employee ID already registered")
	}

	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u := &user.User{
		EmployeeID:       req.EmployeeID,
		FullName:         req.FullName,
		Email:            req.Email,
		PasswordHash:     hashedPassword,
		JobTitle:         req.JobTitle,
		PhoneNumber:      req.PhoneNumber,
		SubUnitName:      req.SubUnitName,
		BranchCode:       req.BranchCode,
		BranchName:       req.BranchName,
		TerminalCode:     req.TerminalCode,
		CompanyCode:      req.CompanyCode,
		PersonnelArea:    req.PersonnelArea,
		PersonnelSubArea: req.PersonnelSubArea,
		Status:           "1",
		CreationDate:     time.Now(),
		CreationBy:       "SYSTEM",
	}

	if err := s.userRepo.Create(ctx, u); err != nil {
		return nil, err
	}

	token, err := s.jwtUtil.GenerateToken(u.ID, u.Email, u.EmployeeID, u.FullName, u.BranchCode, u.TerminalCode)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token, User: user.ToResponse(u), Menus: []user.MenuAccessNode{}}, nil
}

func buildTreeRecursive(parentID *int64, raw []user.MenuAccessRow) []user.MenuAccessNode {
	var nodes []user.MenuAccessNode
	for _, row := range raw {
		isChild := false
		if parentID == nil {
			// Handle cases where postgres DB stores root parent_id as nil OR 0
			if row.ParentMenuID == nil || *row.ParentMenuID == 0 {
				isChild = true
			}
		} else {
			// Find actual child instances
			if row.ParentMenuID != nil && *row.ParentMenuID != 0 && *parentID == *row.ParentMenuID {
				isChild = true
			}
		}
		if isChild {
			id := row.MenuID
			node := user.MenuAccessNode{MenuAccessRow: row, Children: buildTreeRecursive(&id, raw)}
			nodes = append(nodes, node)
		}
	}
	if nodes == nil {
		nodes = []user.MenuAccessNode{}
	}
	return nodes
}

func (s *authService) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	u, err := s.userRepo.FindByEmployeeID(ctx, req.EmployeeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid employee ID or password")
		}
		return nil, err
	}

	if !helper.CheckPasswordHash(req.Password, u.PasswordHash) {
		return nil, errors.New("invalid employee ID or password")
	}

	now := time.Now()
	u.LastLoginAt = &now

	token, err := s.jwtUtil.GenerateToken(u.ID, u.Email, u.EmployeeID, u.FullName, u.BranchCode, u.TerminalCode)
	if err != nil {
		return nil, err
	}

	rawMenus, err := s.userRepo.GetUserMenusByRole(ctx, u.RoleID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token, User: user.ToResponse(u), Menus: buildTreeRecursive(nil, rawMenus)}, nil
}

func (s *authService) ChangeTerminal(ctx context.Context, userID uint64, req *ChangeTerminalRequest) (*AuthResponse, error) {
	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	u.BranchCode = req.BranchCode
	u.TerminalCode = req.TerminalCode

	token, err := s.jwtUtil.GenerateToken(u.ID, u.Email, u.EmployeeID, u.FullName, u.BranchCode, u.TerminalCode)
	if err != nil {
		return nil, err
	}

	rawMenus, err := s.userRepo.GetUserMenusByRole(ctx, u.RoleID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token, User: user.ToResponse(u), Menus: buildTreeRecursive(nil, rawMenus)}, nil
}
