package repository

import (
	"gin-boilerplate/internal/model/entity"

	"gorm.io/gorm"
)

type DermagaRepository interface {
	Create(dermaga *entity.Dermaga) error
	FindAll(kdCabang uint, kdTerminal uint, limit int, offset int) ([]entity.Dermaga, int64, error)
	Update(id uint, dermaga *entity.Dermaga) error
	Delete(id uint) error
	FindByID(id uint) (*entity.Dermaga, error)
}

type dermagaRepository struct {
	db *gorm.DB
}

func NewDermagaRepository(db *gorm.DB) DermagaRepository {
	return &dermagaRepository{db: db}
}

func (r *dermagaRepository) Create(dermaga *entity.Dermaga) error {
	return r.db.Create(dermaga).Error
}

func (r *dermagaRepository) FindAll(kdCabang uint, kdTerminal uint, limit int, offset int) ([]entity.Dermaga, int64, error) {
	var dermagas []entity.Dermaga
	var total int64

	query := r.db.Model(&entity.Dermaga{}).Where("kd_cabang = ? AND kd_terminal = ?", kdCabang, kdTerminal)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Limit(limit).Offset(offset).Find(&dermagas).Error
	if err != nil {
		return nil, 0, err
	}
	return dermagas, total, nil
}

func (r *dermagaRepository) Update(id uint, dermaga *entity.Dermaga) error {
	err := r.db.Where("id = ?", id).Updates(dermaga).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *dermagaRepository) Delete(id uint) error {
	err := r.db.Where("id = ?", id).Delete(&entity.Dermaga{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *dermagaRepository) FindByID(id uint) (*entity.Dermaga, error) {
	var dermaga entity.Dermaga
	err := r.db.Where("id = ?", id).First(&dermaga).Error
	if err != nil {
		return nil, err
	}
	return &dermaga, nil
}
