package service

import (
	"context"
	"errors"
	"gin-boilerplate/internal/model/entity"
	"gin-boilerplate/internal/repository"
)

var (
	ErrUnauthorized = errors.New("unauthorized: record does not belong to your branch or terminal")
	ErrNotFound     = errors.New("record not found")
)

type DermagaService interface {
	Create(ctx context.Context, dermaga *entity.Dermaga) error
	FindAll(ctx context.Context, kdCabang uint, kdTerminal uint, limit int, offset int) ([]entity.Dermaga, int64, error)
	Update(ctx context.Context, id uint, kdCabang uint, kdTerminal uint, dermaga *entity.Dermaga) error
	Delete(ctx context.Context, id uint, kdCabang uint, kdTerminal uint) error
	FindByID(ctx context.Context, id uint) (*entity.Dermaga, error)
}

type dermagaService struct {
	dermagaRepo repository.DermagaRepository
}

func NewDermagaService(dermagaRepo repository.DermagaRepository) DermagaService {
	return &dermagaService{dermagaRepo: dermagaRepo}
}

func (s *dermagaService) Create(ctx context.Context, dermaga *entity.Dermaga) error {
	return s.dermagaRepo.Create(ctx, dermaga)
}

func (s *dermagaService) FindAll(ctx context.Context, kdCabang uint, kdTerminal uint, limit int, offset int) ([]entity.Dermaga, int64, error) {
	return s.dermagaRepo.FindAll(ctx, kdCabang, kdTerminal, limit, offset)
}

func (s *dermagaService) Update(ctx context.Context, id uint, kdCabang uint, kdTerminal uint, dermaga *entity.Dermaga) error {
	existing, err := s.dermagaRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrNotFound
	}

	if existing.KdCabang != kdCabang || existing.KdTerminal != kdTerminal {
		return ErrUnauthorized
	}

	return s.dermagaRepo.Update(ctx, id, dermaga)
}

func (s *dermagaService) Delete(ctx context.Context, id uint, kdCabang uint, kdTerminal uint) error {
	existing, err := s.dermagaRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrNotFound
	}

	if existing.KdCabang != kdCabang || existing.KdTerminal != kdTerminal {
		return ErrUnauthorized
	}

	return s.dermagaRepo.Delete(ctx, id)
}

func (s *dermagaService) FindByID(ctx context.Context, id uint) (*entity.Dermaga, error) {
	return s.dermagaRepo.FindByID(ctx, id)
}
