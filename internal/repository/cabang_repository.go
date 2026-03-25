package repository

import (
	"gin-boilerplate/internal/model/entity"

	"gorm.io/gorm"
)

type CabangRepository interface {
	FindByKdCabangAndKdterminal(kdCabang string, kdTerminal string) (*entity.Cabang, error)
}

type cabangRepository struct {
	db *gorm.DB
}

func NewCabangRepository(db *gorm.DB) CabangRepository {
	return &cabangRepository{db: db}
}

func (r *cabangRepository) FindByKdCabangAndKdterminal(kdCabang string, kdTerminal string) (*entity.Cabang, error) {
	var cabang entity.Cabang
	err := r.db.Where("kd_cabang = ? AND kd_terminal = ?", kdCabang, kdTerminal).First(&cabang).Error
	if err != nil {
		return nil, err
	}
	return &cabang, nil
}
