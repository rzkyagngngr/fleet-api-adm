package repository

import (
	"context"
	"gin-boilerplate/internal/model/entity"

	"gorm.io/gorm"
)

type DermagaRepository interface {
	Create(ctx context.Context, dermaga *entity.Dermaga) error
	FindAll(ctx context.Context, kdCabang uint, kdTerminal uint, limit int, offset int) ([]entity.Dermaga, int64, error)
	Update(ctx context.Context, id uint, dermaga *entity.Dermaga) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*entity.Dermaga, error)
}

type dermagaRepository struct {
	db *gorm.DB
}

func NewDermagaRepository(db *gorm.DB) DermagaRepository {
	return &dermagaRepository{db: db}
}

func (r *dermagaRepository) Create(ctx context.Context, dermaga *entity.Dermaga) error {
	return r.db.WithContext(ctx).Create(dermaga).Error
}

func (r *dermagaRepository) FindAll(ctx context.Context, kdCabang uint, kdTerminal uint, limit int, offset int) ([]entity.Dermaga, int64, error) {
	var dermagas []entity.Dermaga
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Dermaga{}).Where("kd_cabang = ? AND kd_terminal = ?", kdCabang, kdTerminal)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Limit(limit).Offset(offset).Find(&dermagas).Error
	if err != nil {
		return nil, 0, err
	}
	return dermagas, total, nil
}

func (r *dermagaRepository) Update(ctx context.Context, id uint, dermaga *entity.Dermaga) error {
	err := r.db.WithContext(ctx).Where("id = ?", id).Updates(dermaga).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *dermagaRepository) Delete(ctx context.Context, id uint) error {
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.Dermaga{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *dermagaRepository) FindByID(ctx context.Context, id uint) (*entity.Dermaga, error) {
	var dermaga entity.Dermaga
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&dermaga).Error
	if err != nil {
		return nil, err
	}
	return &dermaga, nil
}
