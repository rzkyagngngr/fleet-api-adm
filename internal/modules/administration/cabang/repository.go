package cabang

import (
	"context"

	"gorm.io/gorm"
)

type CabangRepository interface {
	FindByKdCabangAndKdterminal(ctx context.Context, kdCabang string, kdTerminal string) (*Cabang, error)
}

type cabangRepository struct{ db *gorm.DB }

func NewCabangRepository(db *gorm.DB) CabangRepository { return &cabangRepository{db: db} }
func (r *cabangRepository) FindByKdCabangAndKdterminal(ctx context.Context, kdCabang string, kdTerminal string) (*Cabang, error) {
	var c Cabang
	err := r.db.WithContext(ctx).Where("kd_cabang = ? AND kd_terminal = ?", kdCabang, kdTerminal).First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}
