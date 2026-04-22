package terminal

import (
	"context"
	"errors"
	"time"

	"omniport-api/internal/helper"
	"omniport-api/internal/modules/administration/branch"
)

type TerminalService interface {
	Create(ctx context.Context, req *TerminalRequest, companyCode, companyName, userEmp string) error
	Update(ctx context.Context, id uint64, req *TerminalRequest, userEmp string) error
	Delete(ctx context.Context, id uint64) error
	Search(ctx context.Context, param helper.PaginationQuery) ([]Terminal, helper.PaginationMeta, error)
	GetByID(ctx context.Context, id uint64) (*Terminal, error)
	GetStats(ctx context.Context, companyCode string) (*TerminalStats, error)
}

type terminalService struct {
	repo       TerminalRepository
	branchRepo branch.BranchRepository
}

func NewTerminalService(repo TerminalRepository, branchRepo branch.BranchRepository) TerminalService {
	return &terminalService{
		repo:       repo,
		branchRepo: branchRepo,
	}
}

func (s *terminalService) Create(ctx context.Context, req *TerminalRequest, companyCode, companyName, userEmp string) error {
	if _, err := s.repo.GetByCode(ctx, req.TerminalCode); err == nil {
		return errors.New("terminal code already exists")
	}

	// Prefer company info from request if available (passed from frontend based on branch)
	finalCompCode := companyCode
	if req.CompanyCode != "" {
		finalCompCode = req.CompanyCode
	}
	finalCompName := companyName
	if req.CompanyName != "" {
		finalCompName = req.CompanyName
	}

	// Smart lookup: resolve branch name
	branchName := ""
	if b, err := s.branchRepo.GetByCode(ctx, req.BranchCode); err == nil {
		branchName = b.BranchName
	}

	now := time.Now()
	terminal := &Terminal{
		BranchCode:    req.BranchCode,
		BranchName:    branchName,
		TerminalCode:  req.TerminalCode,
		TerminalName:  req.TerminalName,
		GoLiveDate:    req.GoLiveDate,
		IsGoLive:      req.IsGoLive,
		ProfitCenter:  req.ProfitCenter,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		Status:        req.Status,
		VersionCode:   req.VersionCode,
		VersionName:   req.VersionName,
		DocumentCode:  req.DocumentCode,
		CompanyCode:   finalCompCode,
		CompanyName:   finalCompName,
		VesselVersion: req.VesselVersion,
		LogoURL:       req.LogoURL,
		LogoMiniURL:   req.LogoMiniURL,
		Address:       req.Address,
		CompanyType:   req.CompanyType,
		PortCode:      req.PortCode,
		CreatedBy:     userEmp,
		CreatedDate:   &now,
		ProgramName:   "OMNIPORT_ADM",
	}

	return s.repo.Create(ctx, terminal)
}

func (s *terminalService) Update(ctx context.Context, id uint64, req *TerminalRequest, userEmp string) error {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("terminal not found")
	}

	// Update branch name and company info if branch code changed
	if t.BranchCode != req.BranchCode {
		if b, err := s.branchRepo.GetByCode(ctx, req.BranchCode); err == nil {
			t.BranchName = b.BranchName
			t.CompanyCode = b.CompanyCode
			t.CompanyName = b.CompanyName
		}
	} else {
		// Even if branch didn't change, update company info if provided in request
		if req.CompanyCode != "" {
			t.CompanyCode = req.CompanyCode
		}
		if req.CompanyName != "" {
			t.CompanyName = req.CompanyName
		}
	}

	t.BranchCode = req.BranchCode
	t.TerminalName = req.TerminalName
	t.GoLiveDate = req.GoLiveDate
	t.IsGoLive = req.IsGoLive
	t.ProfitCenter = req.ProfitCenter
	t.Latitude = req.Latitude
	t.Longitude = req.Longitude
	t.Status = req.Status
	t.VersionCode = req.VersionCode
	t.VersionName = req.VersionName
	t.DocumentCode = req.DocumentCode
	t.VesselVersion = req.VesselVersion
	t.LogoURL = req.LogoURL
	t.LogoMiniURL = req.LogoMiniURL
	t.Address = req.Address
	t.CompanyType = req.CompanyType
	t.PortCode = req.PortCode
	t.LastUpdatedBy = userEmp
	now := time.Now()
	t.LastUpdatedDate = &now

	return s.repo.Update(ctx, id, t)
}

func (s *terminalService) Delete(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}

func (s *terminalService) Search(ctx context.Context, param helper.PaginationQuery) ([]Terminal, helper.PaginationMeta, error) {
	return s.repo.Search(ctx, param)
}

func (s *terminalService) GetByID(ctx context.Context, id uint64) (*Terminal, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *terminalService) GetStats(ctx context.Context, companyCode string) (*TerminalStats, error) {
	return s.repo.GetStats(ctx, companyCode)
}
