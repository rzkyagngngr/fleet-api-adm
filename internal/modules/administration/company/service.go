package company

import (
	"context"
	"errors"
	"time"

	"omniport-api/internal/helper"
)

type CompanyService interface {
	CreateCompany(ctx context.Context, req *CompanyRequest, createdBy string) error
	UpdateCompany(ctx context.Context, id uint64, req *CompanyRequest, updatedBy string) error
	DeleteCompany(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*CompanyResponse, error)
	SearchCompanies(ctx context.Context, query helper.PaginationQuery) ([]CompanyResponse, helper.PaginationMeta, error)
}

type companyService struct {
	repo CompanyRepository
}

func NewCompanyService(repo CompanyRepository) CompanyService {
	return &companyService{repo: repo}
}

func (s *companyService) CreateCompany(ctx context.Context, req *CompanyRequest, createdBy string) error {
	now := time.Now()
	c := &Company{
		CompanyCode:  req.CompanyCode,
		CompanyName:  req.CompanyName,
		Npwp:         req.Npwp,
		Address:      req.Address,
		Email:        req.Email,
		PhoneNumber:  req.PhoneNumber,
		BusinessType: req.BusinessType,
		Status:       req.Status,
		CreatedBy:    createdBy,
		CreatedDate:  &now,
		ProgramName:  "OMNIPORT_ADM",
	}

	return s.repo.Create(ctx, c)
}

func (s *companyService) UpdateCompany(ctx context.Context, id uint64, req *CompanyRequest, updatedBy string) error {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("company not found")
	}

	c.CompanyCode = req.CompanyCode
	c.CompanyName = req.CompanyName
	c.Npwp = req.Npwp
	c.Address = req.Address
	c.Email = req.Email
	c.PhoneNumber = req.PhoneNumber
	c.BusinessType = req.BusinessType
	c.Status = req.Status

	now := time.Now()
	c.LastUpdatedDate = &now
	c.LastUpdatedBy = updatedBy

	return s.repo.Update(ctx, id, c)
}

func (s *companyService) DeleteCompany(ctx context.Context, id uint64) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("company not found")
	}

	return s.repo.Delete(ctx, id)
}

func (s *companyService) GetByID(ctx context.Context, id uint64) (*CompanyResponse, error) {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := ToResponse(c)
	return &res, nil
}

func (s *companyService) SearchCompanies(ctx context.Context, query helper.PaginationQuery) ([]CompanyResponse, helper.PaginationMeta, error) {
	rows, meta, err := s.repo.Search(ctx, query)
	if err != nil {
		return nil, meta, err
	}

	var res []CompanyResponse
	for _, r := range rows {
		res = append(res, ToResponse(&r))
	}
	return res, meta, nil
}
