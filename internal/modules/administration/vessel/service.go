package vessel

import (
	"context"
	"errors"
	"time"

	"omniport-api/internal/helper"
)

type VesselService interface {
	CreateVessel(ctx context.Context, req *VesselRequest, createdBy string) error
	UpdateVessel(ctx context.Context, id uint64, req *VesselRequest, updatedBy string) error
	DeleteVessel(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*VesselResponse, error)
	SearchVessels(ctx context.Context, query helper.PaginationQuery) ([]VesselResponse, helper.PaginationMeta, error)
}

type vesselService struct {
	repo VesselRepository
}

func NewVesselService(repo VesselRepository) VesselService {
	return &vesselService{repo: repo}
}

func (s *vesselService) CreateVessel(ctx context.Context, req *VesselRequest, createdBy string) error {
	v := &Vessel{
		VesselCode:            req.VesselCode,
		VesselName:            req.VesselName,
		VesselType:            req.VesselType,
		VesselCallSign:        req.VesselCallSign,
		VesselImo:             req.VesselImo,
		VesselGrt:             req.VesselGrt,
		VesselLoa:             req.VesselLoa,
		VesselOwnerName:       req.VesselOwnerName,
		VesselShippingRoute:   req.VesselShippingRoute,
		VesselFlag:            req.VesselFlag,
		VesselCountry:         req.VesselCountry,
		VesselYearMade:        req.VesselYearMade,
		VesselHatchNumber:     req.VesselHatchNumber,
		VesselHatchType:       req.VesselHatchType,
		VesselOwnershipStatus: req.VesselOwnershipStatus,
		VesselOperationStatus: req.VesselOperationStatus,
		Status:                req.Status,
		Remark:                req.Remark,
		PortCode:              req.PortCode,
		TerminalCode:          req.TerminalCode,
		CreationDate:          time.Now(),
		CreationBy:            createdBy,
	}

	return s.repo.Create(ctx, v)
}

func (s *vesselService) UpdateVessel(ctx context.Context, id uint64, req *VesselRequest, updatedBy string) error {
	v, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("vessel not found")
	}

	v.VesselCode = req.VesselCode
	v.VesselName = req.VesselName
	v.VesselType = req.VesselType
	v.VesselCallSign = req.VesselCallSign
	v.VesselImo = req.VesselImo
	v.VesselGrt = req.VesselGrt
	v.VesselLoa = req.VesselLoa
	v.VesselOwnerName = req.VesselOwnerName
	v.VesselShippingRoute = req.VesselShippingRoute
	v.VesselFlag = req.VesselFlag
	v.VesselCountry = req.VesselCountry
	v.VesselYearMade = req.VesselYearMade
	v.VesselHatchNumber = req.VesselHatchNumber
	v.VesselHatchType = req.VesselHatchType
	v.VesselOwnershipStatus = req.VesselOwnershipStatus
	v.VesselOperationStatus = req.VesselOperationStatus
	v.Status = req.Status
	v.Remark = req.Remark
	v.PortCode = req.PortCode
	v.TerminalCode = req.TerminalCode

	now := time.Now()
	v.LastUpdatedDate = &now
	v.LastUpdatedBy = updatedBy

	return s.repo.Update(ctx, id, v)
}

func (s *vesselService) DeleteVessel(ctx context.Context, id uint64) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("vessel not found")
	}

	return s.repo.Delete(ctx, id)
}

func (s *vesselService) GetByID(ctx context.Context, id uint64) (*VesselResponse, error) {
	v, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := ToResponse(v)
	return &res, nil
}

func (s *vesselService) SearchVessels(ctx context.Context, query helper.PaginationQuery) ([]VesselResponse, helper.PaginationMeta, error) {
	rows, meta, err := s.repo.Search(ctx, query)
	if err != nil {
		return nil, meta, err
	}

	var res []VesselResponse
	for _, r := range rows {
		res = append(res, ToResponse(&r))
	}
	return res, meta, nil
}
