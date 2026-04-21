package cargo

import (
	"context"
	"errors"
	"omniport-api/internal/helper"
	"time"
)

type CargoService interface {
	CreateCargo(ctx context.Context, req *CargoRequest, adminEmp string) error
	GetByID(ctx context.Context, id uint64) (*CargoResponse, error)
	UpdateCargo(ctx context.Context, id uint64, req *CargoRequest, adminEmp string) error
	DeleteCargo(ctx context.Context, id uint64) error
	Search(ctx context.Context, query helper.PaginationQuery) ([]CargoResponse, helper.PaginationMeta, error)
	GetStats(ctx context.Context) (*CargoStatsResponse, error)
}

type cargoService struct{ repo CargoRepository }

func NewCargoService(repo CargoRepository) CargoService { return &cargoService{repo: repo} }

func (s *cargoService) CreateCargo(ctx context.Context, req *CargoRequest, adminEmp string) error {
	isActive := "N"
	if req.Status == 1 || req.IsActive == "1" || req.IsActive == "Y" {
		isActive = "Y"
	}

	// Mapping & Validation
	cargoCode := req.CargoCode
	if cargoCode == "" {
		cargoCode = req.ItemCode
	}
	if cargoCode == "" {
		return errors.New("cargo identity code is required (cargo_code or item_code)")
	}

	cargoName := req.CargoName
	if cargoName == "" {
		cargoName = req.ItemName
	}
	if cargoName == "" {
		return errors.New("cargo designation name is required (cargo_name or item_name)")
	}

	now := time.Now()
	c := &Cargo{
		BranchCode:            req.BranchCode,
		TerminalCode:          req.TerminalCode,
		CargoCode:             cargoCode,
		CargoSitcCode:         req.CargoSitcCode,
		CargoHsHarmonizedCode: req.CargoHsHarmonizedCode,
		CargoName:             cargoName,
		CargoGroup:            req.CargoGroup,
		CargoCommodity:        req.CargoCommodity,
		CargoCharacteristic:   req.CargoCharacteristic,
		CargoImdgCode:         req.CargoImdgCode,
		CargoImdgDescription:  req.CargoImdgDescription,
		CargoPackaging1:       req.CargoPackaging1,
		CargoConversion1:      req.CargoConversion1,
		CargoDimension1:       req.CargoDimension1,
		CargoUnit1:            req.CargoUnit1,
		CargoPackaging2:       req.CargoPackaging2,
		CargoConversion2:      req.CargoConversion2,
		CargoDimension2:       req.CargoDimension2,
		CargoUnit2:            req.CargoUnit2,
		CargoPackaging3:       req.CargoPackaging3,
		CargoConversion3:      req.CargoConversion3,
		CargoDimension3:       req.CargoDimension3,
		CargoUnit3:            req.CargoUnit3,
		CargoMooringType:      req.CargoMooringType,
		CargoNotes:            req.CargoNotes,
		CargoCommodityGroup:   req.CargoCommodityGroup,
		CargoCommodityType:    req.CargoCommodityType,
		IsActive:              isActive,
		CargoDocument:         req.CargoDocument,
		HsCode:                req.HsCode,
		HsDescription:         req.HsDescription,
		CargoProductName:      req.CargoProductName,
		CreatedBy:             adminEmp,
		CreatedDate:           &now,
		LastUpdatedBy:         adminEmp,
		LastUpdatedDate:       now,
		ProgramName:           "ADM_SERVICE",
	}

	// Fallback mappings for older field keys if they were the only ones provided
	if c.CargoGroup == "" { c.CargoGroup = req.Category }
	if c.CargoUnit1 == "" { c.CargoUnit1 = req.UOM }
	if c.CargoCharacteristic == "" { c.CargoCharacteristic = req.StorageType }
	if c.CargoProductName == "" { c.CargoProductName = req.Brand }
	if c.CargoNotes == "" { c.CargoNotes = req.Remark }

	return s.repo.Create(ctx, c)
}

func (s *cargoService) GetByID(ctx context.Context, id uint64) (*CargoResponse, error) {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := c.ToResponse()
	return &res, nil
}

func (s *cargoService) UpdateCargo(ctx context.Context, id uint64, req *CargoRequest, adminEmp string) error {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("cargo record not found")
	}

	isActive := "N"
	if req.Status == 1 || req.IsActive == "1" || req.IsActive == "Y" {
		isActive = "Y"
	}

	// Mapping updates & Validation
	if req.CargoCode != "" {
		c.CargoCode = req.CargoCode
	} else if req.ItemCode != "" {
		c.CargoCode = req.ItemCode
	}
	if c.CargoCode == "" {
		return errors.New("cargo identity code is required")
	}

	if req.CargoName != "" {
		c.CargoName = req.CargoName
	} else if req.ItemName != "" {
		c.CargoName = req.ItemName
	}
	if c.CargoName == "" {
		return errors.New("cargo designation name is required")
	}
	
	c.BranchCode = req.BranchCode
	c.TerminalCode = req.TerminalCode
	c.CargoSitcCode = req.CargoSitcCode
	c.CargoHsHarmonizedCode = req.CargoHsHarmonizedCode
	c.CargoGroup = req.CargoGroup
	if c.CargoGroup == "" { c.CargoGroup = req.Category }
	
	c.CargoCommodity = req.CargoCommodity
	c.CargoCharacteristic = req.CargoCharacteristic
	if c.CargoCharacteristic == "" { c.CargoCharacteristic = req.StorageType }
	
	c.CargoImdgCode = req.CargoImdgCode
	c.CargoImdgDescription = req.CargoImdgDescription
	c.CargoPackaging1 = req.CargoPackaging1
	c.CargoConversion1 = req.CargoConversion1
	c.CargoDimension1 = req.CargoDimension1
	c.CargoUnit1 = req.CargoUnit1
	if c.CargoUnit1 == "" { c.CargoUnit1 = req.UOM }
	
	c.CargoPackaging2 = req.CargoPackaging2
	c.CargoConversion2 = req.CargoConversion2
	c.CargoDimension2 = req.CargoDimension2
	c.CargoUnit2 = req.CargoUnit2
	c.CargoPackaging3 = req.CargoPackaging3
	c.CargoConversion3 = req.CargoConversion3
	c.CargoDimension3 = req.CargoDimension3
	c.CargoUnit3 = req.CargoUnit3
	c.CargoMooringType = req.CargoMooringType
	c.CargoNotes = req.CargoNotes
	if c.CargoNotes == "" { c.CargoNotes = req.Remark }
	
	c.CargoCommodityGroup = req.CargoCommodityGroup
	c.CargoCommodityType = req.CargoCommodityType
	c.IsActive = isActive
	c.CargoDocument = req.CargoDocument
	c.HsCode = req.HsCode
	c.HsDescription = req.HsDescription
	c.CargoProductName = req.CargoProductName
	if c.CargoProductName == "" { c.CargoProductName = req.Brand }
	
	c.LastUpdatedBy = adminEmp
	c.LastUpdatedDate = time.Now()

	return s.repo.Update(ctx, id, c)
}

func (s *cargoService) DeleteCargo(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}

func (s *cargoService) Search(ctx context.Context, query helper.PaginationQuery) ([]CargoResponse, helper.PaginationMeta, error) {
	rows, meta, err := s.repo.Search(ctx, query)
	if err != nil {
		return nil, meta, err
	}

	var res []CargoResponse
	for _, r := range rows {
		res = append(res, r.ToResponse())
	}
	return res, meta, nil
}

func (s *cargoService) GetStats(ctx context.Context) (*CargoStatsResponse, error) {
	return s.repo.GetStats(ctx)
}
