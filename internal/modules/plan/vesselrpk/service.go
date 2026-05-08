package vesselrpk

import (
	"context"
	"omniport-api/internal/helper"
)

type VesselRpkService interface {
	Create(ctx context.Context, input CreateVesselRpkInput, branchCode, terminalCode int64, userID string) (*VesselRpkResponse, error)
	GetByID(ctx context.Context, id uint64) (*VesselRpkResponse, error)
	List(ctx context.Context, branchCode, terminalCode int64, page, limit int, search string, filters map[string]interface{}) ([]VesselRpkResponse, helper.PaginationMeta, error)
	Update(ctx context.Context, id uint64, input CreateVesselRpkInput, userID string) error
	Delete(ctx context.Context, id uint64) error
}

type vesselRpkService struct {
	repo VesselRpkRepository
}

func NewVesselRpkService(repo VesselRpkRepository) VesselRpkService {
	return &vesselRpkService{repo: repo}
}

func (s *vesselRpkService) Create(ctx context.Context, input CreateVesselRpkInput, branchCode, terminalCode int64, userID string) (*VesselRpkResponse, error) {
	v := s.mapInputToEntity(input)
	v.BranchCode = branchCode
	v.TerminalCode = terminalCode
	v.CreationBy = userID

	if err := s.repo.Create(ctx, v); err != nil {
		return nil, err
	}

	return s.mapEntityToResponse(v), nil
}

func (s *vesselRpkService) GetByID(ctx context.Context, id uint64) (*VesselRpkResponse, error) {
	v, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.mapEntityToResponse(v), nil
}

func (s *vesselRpkService) List(ctx context.Context, branchCode, terminalCode int64, page, limit int, search string, filters map[string]interface{}) ([]VesselRpkResponse, helper.PaginationMeta, error) {
	offset := (page - 1) * limit
	list, total, err := s.repo.List(ctx, branchCode, terminalCode, offset, limit, search, filters)
	if err != nil {
		return nil, helper.PaginationMeta{}, err
	}

	var res []VesselRpkResponse
	for _, v := range list {
		res = append(res, *s.mapEntityToResponse(&v))
	}

	return res, helper.NewPaginationMeta(total, page, limit), nil
}

func (s *vesselRpkService) Update(ctx context.Context, id uint64, input CreateVesselRpkInput, userID string) error {
	v := s.mapInputToEntity(input)
	v.LastUpdatedBy = userID
	return s.repo.Update(ctx, id, v)
}

func (s *vesselRpkService) Delete(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}

// ─────────────────────────────────────────────────────────────
// MAPPERS
// ─────────────────────────────────────────────────────────────

func (s *vesselRpkService) mapInputToEntity(input CreateVesselRpkInput) *VesselRpk {
	v := &VesselRpk{
		NoPkk:                  input.NoPkk,
		NoPpk:                  input.NoPpk,
		LocationCodeInaportnet: input.LocationCodeInaportnet,
		RpkType:                input.RpkType,
		BerthPosition:          input.BerthPosition,
		VesselPosition:         input.VesselPosition,
		StartMeter:             input.StartMeter,
		EndMeter:               input.EndMeter,
		StartMooring:           input.StartMooring,
		EndMooring:             input.EndMooring,
		RampDoor:               input.RampDoor,
		Distribution:           input.Distribution,
		Packaging:              input.Packaging,
		NoRkbm:                 input.NoRkbm,
		Reason:                 input.Reason,
		Notes:                  input.Notes,
		Payload:                input.Payload,
		BranchCode:             input.BranchCode,
		TerminalCode:           input.TerminalCode,
	}

	if input.Op != nil {
		op := &VesselRpkOp{
			ID:               input.Op.ID,
			Pbm:              input.Op.Pbm,
			Emkl:             input.Op.Emkl,
			Shipper:          input.Op.Shipper,
			StartDischarging: input.Op.StartDischarging,
			EndDischarging:   input.Op.EndDischarging,
		}

		for _, d := range input.Op.OpDetail {
			op.OpDetail = append(op.OpDetail, VesselRpkOpDetail{
				ID:                d.ID,
				RkbmMuatNumber:    d.RkbmMuatNumber,
				RkbmBongkarNumber: d.RkbmBongkarNumber,
				Loading:           d.Loading,
				Discharging:       d.Discharging,
				Commodity:         d.Commodity,
			})
		}
		v.Op = op
	}

	return v
}

func (s *vesselRpkService) mapEntityToResponse(v *VesselRpk) *VesselRpkResponse {
	res := &VesselRpkResponse{
		ID:                     v.ID,
		NoPkk:                  v.NoPkk,
		NoPpk:                  v.NoPpk,
		LocationCodeInaportnet: v.LocationCodeInaportnet,
		RpkType:                v.RpkType,
		BerthPosition:          v.BerthPosition,
		VesselPosition:         v.VesselPosition,
		StartMeter:             v.StartMeter,
		EndMeter:               v.EndMeter,
		StartMooring:           v.StartMooring,
		EndMooring:             v.EndMooring,
		RampDoor:               v.RampDoor,
		Distribution:           v.Distribution,
		Packaging:              v.Packaging,
		NoRkbm:                 v.NoRkbm,
		Reason:                 v.Reason,
		Notes:                  v.Notes,
		Payload:                v.Payload,
		BranchCode:             v.BranchCode,
		TerminalCode:           v.TerminalCode,
		CreationDate:           v.CreationDate,
		CreationBy:             v.CreationBy,
		LastUpdatedDate:        v.LastUpdatedDate,
		LastUpdatedBy:          v.LastUpdatedBy,
	}

	if v.Op != nil {
		opRes := &VesselRpkOpResponse{
			ID:               v.Op.ID,
			Pbm:              v.Op.Pbm,
			Emkl:             v.Op.Emkl,
			Shipper:          v.Op.Shipper,
			StartDischarging: v.Op.StartDischarging,
			EndDischarging:   v.Op.EndDischarging,
		}

		for _, d := range v.Op.OpDetail {
			opRes.OpDetail = append(opRes.OpDetail, VesselRpkOpDetailResponse{
				ID:                d.ID,
				RkbmMuatNumber:    d.RkbmMuatNumber,
				RkbmBongkarNumber: d.RkbmBongkarNumber,
				Loading:           d.Loading,
				Discharging:       d.Discharging,
				Commodity:         d.Commodity,
			})
		}
		res.Op = opRes
	}

	return res
}
