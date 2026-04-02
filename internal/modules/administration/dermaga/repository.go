package dermaga

import (
	"context"

	"gorm.io/gorm"
)

type DermagaRepository interface {
	Create(ctx context.Context, dermaga *Dermaga) error
	FindAll(ctx context.Context, kdCabang uint, kdTerminal uint, limit int, offset int) ([]Dermaga, int64, error)
	Update(ctx context.Context, id uint, dermaga *Dermaga) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*Dermaga, error)
}

type dermagaRepository struct{ db *gorm.DB }

func NewDermagaRepository(db *gorm.DB) DermagaRepository { return &dermagaRepository{db: db} }
func (r *dermagaRepository) Create(ctx context.Context, d *Dermaga) error {
	return r.db.WithContext(ctx).Create(d).Error
}
func (r *dermagaRepository) FindAll(ctx context.Context, kdCabang uint, kdTerminal uint, limit int, offset int) ([]Dermaga, int64, error) {
	var rows []Dermaga
	var total int64
	q := r.db.WithContext(ctx).Model(&Dermaga{}).Where("kd_cabang = ? AND kd_terminal = ?", kdCabang, kdTerminal)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Limit(limit).Offset(offset).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}
func (r *dermagaRepository) Update(ctx context.Context, id uint, d *Dermaga) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Updates(d).Error
}
func (r *dermagaRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&Dermaga{}).Error
}
func (r *dermagaRepository) FindByID(ctx context.Context, id uint) (*Dermaga, error) {
	var d Dermaga
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&d).Error; err != nil {
		return nil, err
	}
	return &d, nil
}
