package branch

import (
	"context"
	"errors"
	"time"

	"omniport-api/internal/helper"
)

type BranchService interface {
	Create(ctx context.Context, req *BranchRequest, companyCode, companyName, userEmp string) error
	Update(ctx context.Context, id uint64, req *BranchRequest, userEmp string) error
	Delete(ctx context.Context, id uint64) error
	Search(ctx context.Context, param helper.PaginationQuery) ([]Branch, helper.PaginationMeta, error)
	GetByID(ctx context.Context, id uint64) (*Branch, error)
	GetStats(ctx context.Context) (*BranchStats, error)
}

type branchService struct {
	repo BranchRepository
}

func NewBranchService(repo BranchRepository) BranchService {
	return &branchService{repo: repo}
}

func (s *branchService) Create(ctx context.Context, req *BranchRequest, companyCode, companyName, userEmp string) error {
	if _, err := s.repo.GetByCode(ctx, req.BranchCode); err == nil {
		return errors.New("branch code already exists")
	}

	now := time.Now()
	branch := &Branch{
		BranchCode:      req.BranchCode,
		BranchName:      req.BranchName,
		CompanyCode:     companyCode,
		CompanyName:     companyName,
		KdPort:          req.KdPort,
		Address:         req.Address,
		Status:          req.Status,
		CreatedBy:       userEmp,
		CreatedDate:     &now,
		ProgramName:     "OMNIPORT_ADM",
	}

	return s.repo.Create(ctx, branch)
}

func (s *branchService) Update(ctx context.Context, id uint64, req *BranchRequest, userEmp string) error {
	b, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("branch not found")
	}

	b.BranchName = req.BranchName
	b.KdPort = req.KdPort
	b.Address = req.Address
	b.Status = req.Status
	b.LastUpdatedBy = userEmp
	now := time.Now()
	b.LastUpdatedDate = &now

	return s.repo.Update(ctx, id, b)
}

func (s *branchService) Delete(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}

func (s *branchService) Search(ctx context.Context, param helper.PaginationQuery) ([]Branch, helper.PaginationMeta, error) {
	return s.repo.Search(ctx, param)
}

func (s *branchService) GetByID(ctx context.Context, id uint64) (*Branch, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *branchService) GetStats(ctx context.Context) (*BranchStats, error) {
	return s.repo.GetStats(ctx)
}
