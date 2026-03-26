package service

import (
	"errors"
	"gin-boilerplate/internal/model/entity"
	"time"

	"gin-boilerplate/internal/repository"
	"gin-boilerplate/pkg/utils"

	"gin-boilerplate/internal/model/dto"

	"gorm.io/gorm"
)

type AuthService interface {
	Register(req *dto.UserRegisterRequest) (*dto.AuthResponse, error)
	Login(req *dto.LoginRequest) (*dto.AuthResponse, error)
}

type authService struct {
	userRepo repository.UserRepository
	jwtUtil  *utils.JWTUtil
}

func NewAuthService(userRepo repository.UserRepository, jwtUtil *utils.JWTUtil) AuthService {
	return &authService{
		userRepo: userRepo,
		jwtUtil:  jwtUtil,
	}
}

func (s *authService) Register(req *dto.UserRegisterRequest) (*dto.AuthResponse, error) {
	// Check if email already exists
	_, err := s.userRepo.FindByEmail(req.Email)
	if err == nil {
		return nil, errors.New("email already registered")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Check if employee ID already exists
	_, err = s.userRepo.FindByEmployeeID(req.EmployeeID)
	if err == nil {
		return nil, errors.New("employee ID already registered")
	}

	// Hash password using SHA-512 (as requested)
	hashedPassword := utils.HashSHA512(req.Password)

	user := &entity.User{
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
		Status:           "1", // Default active
		CreationDate:     time.Now(),
		CreationBy:       "SYSTEM", // Or from context if available
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := s.jwtUtil.GenerateToken(user.ID, user.Email, user.EmployeeID, user.FullName, user.BranchCode)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token: token,
		User:  user.ToResponse(),
		Menus: []dto.MenuAccessNode{},
	}, nil
}

// buildTreeRecursive helps to convert flat menu rows into hierarchical nodes
func buildTreeRecursive(parentID *int64, raw []dto.MenuAccessRow) []dto.MenuAccessNode {
	var nodes []dto.MenuAccessNode
	for _, row := range raw {
		isChild := false
		if parentID == nil && row.ParentMenuID == nil {
			isChild = true
		} else if parentID != nil && row.ParentMenuID != nil && *parentID == *row.ParentMenuID {
			isChild = true
		}

		if isChild {
			menuID := row.MenuID
			node := dto.MenuAccessNode{
				MenuAccessRow: row,
				Children:      buildTreeRecursive(&menuID, raw),
			}
			nodes = append(nodes, node)
		}
	}
	if nodes == nil {
		nodes = []dto.MenuAccessNode{}
	}
	return nodes
}

func (s *authService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.FindByEmployeeID(req.EmployeeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid employee ID or password")
		}
		return nil, err
	}

	// Compare password using SHA-512 comparison
	inputHash := utils.HashSHA512(req.Password)
	if user.PasswordHash != inputHash {
		return nil, errors.New("invalid employee ID or password")
	}

	// Update last login
	now := time.Now()
	user.LastLoginAt = &now
	// We might want to save this in the repo, but for now we proceed

	// Generate JWT token
	token, err := s.jwtUtil.GenerateToken(user.ID, user.Email, user.EmployeeID, user.FullName, user.BranchCode)
	if err != nil {
		return nil, err
	}

	// Fetch Menus from vw_access_login
	rawMenus, err := s.userRepo.GetUserMenusByRole(user.RoleID)
	if err != nil {
		return nil, err
	}

	// Build hierarchy tree
	menuTree := buildTreeRecursive(nil, rawMenus)

	return &dto.AuthResponse{
		Token: token,
		User:  user.ToResponse(),
		Menus: menuTree,
	}, nil
}
