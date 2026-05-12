package vesselrpk

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type VesselRpkRepository interface {
	Create(ctx context.Context, v *VesselRpk) error
	GetByID(ctx context.Context, id uint64) (*VesselRpk, error)
	List(ctx context.Context, branchCode, terminalCode int64, offset, limit int, search string, filters map[string]interface{}) ([]VesselRpk, int64, error)
	Update(ctx context.Context, id uint64, v *VesselRpk) error
	Delete(ctx context.Context, id uint64) error
}

type vesselRpkRepository struct {
	db *gorm.DB
}

func NewVesselRpkRepository(db *gorm.DB) VesselRpkRepository {
	return &vesselRpkRepository{db: db}
}

func (r *vesselRpkRepository) Create(ctx context.Context, v *VesselRpk) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Create(v).Error
	})
}

func (r *vesselRpkRepository) GetByID(ctx context.Context, id uint64) (*VesselRpk, error) {
	var v VesselRpk
	err := r.db.WithContext(ctx).
		Preload("Ops").
		Preload("Ops.OpDetail").
		First(&v, id).Error
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *vesselRpkRepository) List(ctx context.Context, branchCode, terminalCode int64, offset, limit int, search string, filters map[string]interface{}) ([]VesselRpk, int64, error) {
	if r == nil || r.db == nil {
		return nil, 0, errors.New("vessel rpk database is not initialized")
	}

	var list []VesselRpk
	var total int64

	query := r.db.WithContext(ctx).Model(&VesselRpk{}).
		Select("plan.post_vessel_rpk.*, vp.vessel_name").
		Joins("LEFT JOIN plan.post_vessel_plan vp ON plan.post_vessel_rpk.ops_plan_code = vp.plan_code")

	// Multi-tenancy
	if branchCode > 0 {
		query = query.Where("plan.post_vessel_rpk.branch_code = ?", branchCode)
	}
	if terminalCode > 0 {
		query = query.Where("plan.post_vessel_rpk.terminal_code = ?", terminalCode)
	}

	// Advanced Search (JSONB & Related Fields)
	if search != "" {
		s := "%" + search + "%"
		query = query.Where("(plan.post_vessel_rpk.no_pkk ILIKE ? OR plan.post_vessel_rpk.no_ppk ILIKE ? OR plan.post_vessel_rpk.no_rkbm ILIKE ? OR vp.vessel_name ILIKE ?)", s, s, s, s)
	}

	// Dynamic Filters from Frontend
	if filters != nil {
		allowedCols := []string{"no_pkk", "rpk_type", "distribution", "creation_by", "start_mooring"}
		for _, col := range allowedCols {
			if val, ok := filters[col]; ok && val != "" {
				query = query.Where("plan.post_vessel_rpk."+col+" ILIKE ?", "%"+fmt.Sprintf("%v", val)+"%")
			}
		}

		if branchCode == 0 {
			if val, ok := filters["branch_code"]; ok && val != "" && val != "0" {
				query = query.Where("plan.post_vessel_rpk.branch_code = ?", val)
			}
		}
		if terminalCode == 0 {
			if val, ok := filters["terminal_code"]; ok && val != "" && val != "0" {
				query = query.Where("plan.post_vessel_rpk.terminal_code = ?", val)
			}
		}
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Ops").Preload("Ops.OpDetail").
		Offset(offset).Limit(limit).
		Order("plan.post_vessel_rpk.creation_date DESC").Find(&list).Error

	return list, total, err
}

func (r *vesselRpkRepository) Update(ctx context.Context, id uint64, v *VesselRpk) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Update Header (Omit relation first)
		if err := tx.Model(&VesselRpk{}).Where("id = ?", id).Omit("Ops").Updates(v).Error; err != nil {
			return err
		}

		// 2. Smart-Sync Ops & OpDetails
		var incomingOpIDs []uint64
		for _, op := range v.Ops {
			op.VesselRpkID = id
			if op.ID > 0 {
				incomingOpIDs = append(incomingOpIDs, op.ID)
				if err := tx.Save(&op).Error; err != nil {
					return err
				}
			} else {
				if err := tx.Create(&op).Error; err != nil {
					return err
				}
				incomingOpIDs = append(incomingOpIDs, op.ID)
			}

			// Sync Details for this Op
			var incomingDetailIDs []uint64
			for _, detail := range op.OpDetail {
				detail.VesselRpkOpID = op.ID
				if detail.ID > 0 {
					incomingDetailIDs = append(incomingDetailIDs, detail.ID)
					tx.Save(&detail)
				} else {
					tx.Create(&detail)
					incomingDetailIDs = append(incomingDetailIDs, detail.ID)
				}
			}

			// Delete removed details for this specific Op
			if len(incomingDetailIDs) > 0 {
				tx.Where("vessel_rpk_op_id = ? AND id NOT IN ?", op.ID, incomingDetailIDs).Delete(&VesselRpkOpDetail{})
			} else {
				tx.Where("vessel_rpk_op_id = ?", op.ID).Delete(&VesselRpkOpDetail{})
			}
		}

		// 3. Delete removed Ops (and cascading details if DB allows, else manual delete)
		if len(incomingOpIDs) > 0 {
			// First delete details of ops that will be removed
			tx.Where("vessel_rpk_op_id IN (SELECT id FROM plan.post_vessel_rpk_op WHERE vessel_rpk_id = ? AND id NOT IN ?)", id, incomingOpIDs).Delete(&VesselRpkOpDetail{})
			// Then delete the ops
			tx.Where("vessel_rpk_id = ? AND id NOT IN ?", id, incomingOpIDs).Delete(&VesselRpkOp{})
		} else {
			tx.Where("vessel_rpk_op_id IN (SELECT id FROM plan.post_vessel_rpk_op WHERE vessel_rpk_id = ?)", id).Delete(&VesselRpkOpDetail{})
			tx.Where("vessel_rpk_id = ?", id).Delete(&VesselRpkOp{})
		}

		return nil
	})
}

func (r *vesselRpkRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Manual cascade delete
		tx.Where("vessel_rpk_op_id IN (SELECT id FROM plan.post_vessel_rpk_op WHERE vessel_rpk_id = ?)", id).Delete(&VesselRpkOpDetail{})
		tx.Where("vessel_rpk_id = ?", id).Delete(&VesselRpkOp{})
		return tx.Delete(&VesselRpk{}, id).Error
	})
}
