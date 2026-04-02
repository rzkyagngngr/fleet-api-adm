package dermaga

import (
	"context"
	"errors"
)

var (
	ErrUnauthorized = errors.New("unauthorized: record does not belong to your branch or terminal")
	ErrNotFound     = errors.New("record not found")
)

type DermagaService interface {
	Create(ctx context.Context, dermaga *Dermaga) error
	FindAll(ctx context.Context, kdCabang uint, kdTerminal uint, limit int, offset int) ([]Dermaga, int64, error)
	Update(ctx context.Context, id uint, kdCabang uint, kdTerminal uint, dermaga *Dermaga) error
	Delete(ctx context.Context, id uint, kdCabang uint, kdTerminal uint) error
	FindByID(ctx context.Context, id uint) (*Dermaga, error)
}

type dermagaService struct{ dermagaRepo DermagaRepository }

func NewDermagaService(repo DermagaRepository) DermagaService {
	return &dermagaService{dermagaRepo: repo}
}
func (s *dermagaService) Create(ctx context.Context, d *Dermaga) error {
	return s.dermagaRepo.Create(ctx, d)
}
func (s *dermagaService) FindAll(ctx context.Context, kdCabang uint, kdTerminal uint, limit int, offset int) ([]Dermaga, int64, error) {
	return s.dermagaRepo.FindAll(ctx, kdCabang, kdTerminal, limit, offset)
}
func (s *dermagaService) Update(ctx context.Context, id uint, kdCabang uint, kdTerminal uint, d *Dermaga) error {
	ex, err := s.dermagaRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if ex == nil {
		return ErrNotFound
	}
	if ex.KdCabang != kdCabang || ex.KdTerminal != kdTerminal {
		return ErrUnauthorized
	}
	return s.dermagaRepo.Update(ctx, id, d)
}
func (s *dermagaService) Delete(ctx context.Context, id uint, kdCabang uint, kdTerminal uint) error {
	ex, err := s.dermagaRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if ex == nil {
		return ErrNotFound
	}
	if ex.KdCabang != kdCabang || ex.KdTerminal != kdTerminal {
		return ErrUnauthorized
	}
	return s.dermagaRepo.Delete(ctx, id)
}
func (s *dermagaService) FindByID(ctx context.Context, id uint) (*Dermaga, error) {
	return s.dermagaRepo.FindByID(ctx, id)
}
