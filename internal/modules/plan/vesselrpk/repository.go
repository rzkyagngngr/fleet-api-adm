package vesselrpk

import (
	"context"
	"gorm.io/gorm"
)

type VesselRpkRepository interface {
	Create(ctx context.Context, v *VesselRpk) error
	GetByID(ctx context.Context, id uint64) (*VesselRpk, error)
	List(ctx context.Context, branchCode, terminalCode int64, offset, limit int, search string) ([]VesselRpk, int64, error)
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
		Preload("Op").
		Preload("Op.OpDetail").
		First(&v, id).Error
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *vesselRpkRepository) List(ctx context.Context, branchCode, terminalCode int64, offset, limit int, search string) ([]VesselRpk, int64, error) {
	var list []VesselRpk
	var total int64

	query := r.db.WithContext(ctx).Model(&VesselRpk{})

	// Multi-tenancy
	if branchCode > 0 {
		query = query.Where("branch_code = ?", branchCode)
	}
	if terminalCode > 0 {
		query = query.Where("terminal_code = ?", terminalCode)
	}

	// Advanced Search (JSONB & Related Fields)
	if search != "" {
		s := "%" + search + "%"
		query = query.Where("(no_pkk ILIKE ? OR no_ppk ILIKE ? OR no_rkbm ILIKE ?)", s, s, s)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Preload nested relations
	err = query.Preload("Op").Preload("Op.OpDetail").
		Offset(offset).Limit(limit).
		Order("creation_date DESC").Find(&list).Error
	
	return list, total, err
}

func (r *vesselRpkRepository) Update(ctx context.Context, id uint64, v *VesselRpk) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Update Header (Omit relation first)
		if err := tx.Model(&VesselRpk{}).Where("id = ?", id).Omit("Op").Updates(v).Error; err != nil {
			return err
		}

		// 2. Smart-Sync Op & OpDetail
		if v.Op != nil {
			var existingOp VesselRpkOp
			err := tx.Where("vessel_rpk_id = ?", id).First(&existingOp).Error
			
			if err == gorm.ErrRecordNotFound {
				// Create new if not exists
				v.Op.VesselRpkID = id
				if err := tx.Create(v.Op).Error; err != nil {
					return err
				}
			} else if err == nil {
				// Update existing Op
				v.Op.ID = existingOp.ID
				v.Op.VesselRpkID = id
				if err := tx.Save(v.Op).Error; err != nil {
					return err
				}

				// Sync Details (Manual Sync to preserve IDs)
				var incomingDetailIDs []uint64
				for _, detail := range v.Op.OpDetail {
					detail.VesselRpkOpID = existingOp.ID
					if detail.ID > 0 {
						incomingDetailIDs = append(incomingDetailIDs, detail.ID)
						tx.Save(&detail)
					} else {
						tx.Create(&detail)
						incomingDetailIDs = append(incomingDetailIDs, detail.ID)
					}
				}
				
				// Delete removed details
				if len(incomingDetailIDs) > 0 {
					tx.Where("vessel_rpk_op_id = ? AND id NOT IN ?", existingOp.ID, incomingDetailIDs).Delete(&VesselRpkOpDetail{})
				} else {
					tx.Where("vessel_rpk_op_id = ?", existingOp.ID).Delete(&VesselRpkOpDetail{})
				}
			}
		}

		return nil
	})
}

func (r *vesselRpkRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&VesselRpk{}, id).Error
}
