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
	RefreshToken(ctx context.Context, userID uint64) (*AuthResponse, error)
}

type authService struct {
	userRepo user.UserRepository
	db       *gorm.DB
	jwtUtil  *helper.JWTUtil
}

func NewAuthService(userRepo user.UserRepository, db *gorm.DB, jwtUtil *helper.JWTUtil) AuthService {
	return &authService{userRepo: userRepo, db: db, jwtUtil: jwtUtil}
}

// --- Helpers ---

func buildTreeRecursive(parentID *int64, raw []user.MenuAccessRow) []user.MenuAccessNode {
	var nodes []user.MenuAccessNode
	for _, row := range raw {
		isChild := false
		if parentID == nil {
			if row.ParentMenuID == nil || *row.ParentMenuID == 0 {
				isChild = true
			}
		} else {
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

// lookupBranchName fetches branch_name from posm_branches by code.
func (s *authService) lookupBranchName(ctx context.Context, code string) string {
	var name string
	s.db.WithContext(ctx).
		Table("adm.posm_branches").
		Select("branch_name").
		Where("branch_code = ?", code).
		Limit(1).
		Scan(&name)
	return name
}

// lookupTerminalName fetches terminal_name from posm_terminals by code.
func (s *authService) lookupTerminalName(ctx context.Context, code string) string {
	var name string
	s.db.WithContext(ctx).
		Table("adm.posm_terminals").
		Select("terminal_name").
		Where("terminal_code = ?", code).
		Limit(1).
		Scan(&name)
	return name
}

// buildToken wraps jwtUtil.GenerateToken with the full user context.
func (s *authService) buildToken(u *user.User) (string, error) {
	return s.jwtUtil.GenerateToken(
		u.ID, u.Email, u.EmployeeID, u.FullName,
		u.BranchCode, u.BranchName,
		u.TerminalCode, u.TerminalName,
		u.CompanyCode, u.CompanyName,
	)
}

// --- Service Implementations ---

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

	token, err := s.buildToken(u)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token, User: user.ToResponse(u), Menus: []user.MenuAccessNode{}}, nil
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

	// Re-fetch with full relations (Branches, Terminals)
	uFull, err := s.userRepo.FindByID(ctx, u.ID)
	if err != nil {
		return nil, err
	}
	u = uFull

	// GUARD: User must have at least one terminal, one branch, and a valid company code.
	// Without these, the system cannot establish a valid operational context.
	if len(u.Terminals) == 0 || len(u.Branches) == 0 || u.CompanyCode == "" {
		return nil, errors.New("[ACCESS_DENIED] Akun Anda belum memiliki akses terminal atau cabang yang valid. Silakan hubungi administrator untuk mengatur akses Anda.")
	}

	// Determine active branch:
	// - If user has a saved preference (branch_code in posm_users from a prior ChangeTerminal), use it.
	// - Otherwise, default to the first entry from posm_user_branches (M2M junction).
	// Name is always resolved from the preloaded posm_branches record, not from posm_users.
	activeBranchCode := u.BranchCode
	if activeBranchCode == "" && len(u.Branches) > 0 {
		activeBranchCode = u.Branches[0].BranchCode
	}
	activeBranchName := ""
	if activeBranchCode != "" {
		for _, b := range u.Branches {
			if b.BranchCode == activeBranchCode {
				activeBranchName = b.BranchName
				break
			}
		}
		if activeBranchName == "" {
			activeBranchName = s.lookupBranchName(ctx, activeBranchCode)
		}
	}
	u.BranchCode = activeBranchCode
	u.BranchName = activeBranchName

	// Determine active terminal:
	// - If user has a saved preference (terminal_code in posm_users from a prior ChangeTerminal), use it.
	// - Otherwise, default to the first entry from posm_user_terminals (M2M junction).
	// Name is always resolved from the preloaded posm_terminals record, not from posm_users.
	activeTerminalCode := u.TerminalCode
	if activeTerminalCode == "" && len(u.Terminals) > 0 {
		activeTerminalCode = u.Terminals[0].TerminalCode
	}
	activeTerminalName := ""
	if activeTerminalCode != "" {
		for _, t := range u.Terminals {
			if t.TerminalCode == activeTerminalCode {
				activeTerminalName = t.TerminalName
				break
			}
		}
		if activeTerminalName == "" {
			activeTerminalName = s.lookupTerminalName(ctx, activeTerminalCode)
		}
	}
	u.TerminalCode = activeTerminalCode
	u.TerminalName = activeTerminalName

	now := time.Now()
	u.LastLoginAt = &now

	token, err := s.buildToken(u)
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
	// 1. Load current user (needed for token generation: email, employeeID, etc.)
	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 2. Resolve names from master data tables
	branchName := s.lookupBranchName(ctx, req.BranchCode)
	terminalName := s.lookupTerminalName(ctx, req.TerminalCode)

	// 3. Persist new context to DB using a dedicated, lightweight update
	if err := s.userRepo.UpdateTerminalContext(
		ctx, u.ID,
		req.BranchCode, branchName,
		req.TerminalCode, terminalName,
	); err != nil {
		return nil, errors.New("failed to update terminal context")
	}

	// 4. Apply to in-memory user object for token & response generation
	u.BranchCode = req.BranchCode
	u.BranchName = branchName
	u.TerminalCode = req.TerminalCode
	u.TerminalName = terminalName

	// 5. Issue new JWT with the full, updated identity claims
	token, err := s.buildToken(u)
	if err != nil {
		return nil, err
	}

	// 6. Fetch menus
	rawMenus, err := s.userRepo.GetUserMenusByRole(ctx, u.RoleID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token, User: user.ToResponse(u), Menus: buildTreeRecursive(nil, rawMenus)}, nil
}

func (s *authService) RefreshToken(ctx context.Context, userID uint64) (*AuthResponse, error) {
	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Ensure names are populated even on refresh
	if u.BranchCode != "" && u.BranchName == "" {
		u.BranchName = s.lookupBranchName(ctx, u.BranchCode)
	}
	if u.TerminalCode != "" && u.TerminalName == "" {
		u.TerminalName = s.lookupTerminalName(ctx, u.TerminalCode)
	}

	token, err := s.buildToken(u)
	if err != nil {
		return nil, err
	}

	rawMenus, err := s.userRepo.GetUserMenusByRole(ctx, u.RoleID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token, User: user.ToResponse(u), Menus: buildTreeRecursive(nil, rawMenus)}, nil
}
