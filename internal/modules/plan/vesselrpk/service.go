package vesselrpk

import (
	"context"
	"errors"
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
	if s == nil || s.repo == nil {
		return nil, helper.PaginationMeta{}, errors.New("vessel rpk repository is not initialized")
	}
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit
	list, total, err := s.repo.List(ctx, branchCode, terminalCode, offset, limit, search, filters)
	if err != nil {
		return nil, helper.PaginationMeta{}, err
	}

	var res []VesselRpkResponse
	for _, v := range list {
		mapped := s.mapEntityToResponse(&v)
		if mapped != nil {
			res = append(res, *mapped)
		}
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
		OpsPlanCode:            input.OpsPlanCode,
		ActivityCode:           input.ActivityCode,
		Payload:                input.Payload,
		BranchCode:             input.BranchCode,
		TerminalCode:           input.TerminalCode,
	}

	for _, opInput := range input.Ops {
		op := VesselRpkOp{
			ID:               opInput.ID,
			Pbm:              opInput.Pbm,
			Emkl:             opInput.Emkl,
			Shipper:          opInput.Shipper,
			StartDischarging: opInput.StartDischarging,
			EndDischarging:   opInput.EndDischarging,
			StartActivityDate: opInput.StartActivityDate,
			EndActivityDate:   opInput.EndActivityDate,
		}

		for _, d := range opInput.OpDetail {
			op.OpDetail = append(op.OpDetail, VesselRpkOpDetail{
				ID:                d.ID,
				RkbmMuatNumber:    d.RkbmMuatNumber,
				RkbmBongkarNumber: d.RkbmBongkarNumber,
				Loading:           d.Loading,
				Discharging:       d.Discharging,
				Commodity:         d.Commodity,
			})
		}
		v.Ops = append(v.Ops, op)
	}

	return v
}

func (s *vesselRpkService) mapEntityToResponse(v *VesselRpk) *VesselRpkResponse {
	if v == nil {
		return nil
	}

	res := &VesselRpkResponse{
		ID:                     v.ID,
		NoPkk:                  v.NoPkk,
		NoPpk:                  v.NoPpk,
		VesselName:             v.VesselName,
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
		OpsPlanCode:            v.OpsPlanCode,
		ActivityCode:           v.ActivityCode,
		Payload:                v.Payload,
		BranchCode:             v.BranchCode,
		TerminalCode:           v.TerminalCode,
		CreationDate:           v.CreationDate,
		CreationBy:             v.CreationBy,
		LastUpdatedDate:        v.LastUpdatedDate,
		LastUpdatedBy:          v.LastUpdatedBy,
	}

	for _, opEntity := range v.Ops {
		opRes := VesselRpkOpResponse{
			ID:               opEntity.ID,
			Pbm:              opEntity.Pbm,
			Emkl:             opEntity.Emkl,
			Shipper:          opEntity.Shipper,
			StartDischarging: opEntity.StartDischarging,
			EndDischarging:   opEntity.EndDischarging,
			StartActivityDate: opEntity.StartActivityDate,
			EndActivityDate:   opEntity.EndActivityDate,
		}

		for _, d := range opEntity.OpDetail {
			opRes.OpDetail = append(opRes.OpDetail, VesselRpkOpDetailResponse{
				ID:                d.ID,
				RkbmMuatNumber:    d.RkbmMuatNumber,
				RkbmBongkarNumber: d.RkbmBongkarNumber,
				Loading:           d.Loading,
				Discharging:       d.Discharging,
				Commodity:         d.Commodity,
			})
		}
		res.Ops = append(res.Ops, opRes)
	}

	return res
}
