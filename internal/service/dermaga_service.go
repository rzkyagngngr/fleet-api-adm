package service

import (
	"errors"
	"gin-boilerplate/internal/model/entity"
	"gin-boilerplate/internal/repository"
)

var (
	ErrUnauthorized = errors.New("unauthorized: record does not belong to your branch or terminal")
	ErrNotFound     = errors.New("record not found")
)

type DermagaService interface {
	Create(dermaga *entity.Dermaga) error
	FindAll(kdCabang uint, kdTerminal uint, limit int, offset int) ([]entity.Dermaga, int64, error)
	Update(id uint, kdCabang uint, kdTerminal uint, dermaga *entity.Dermaga) error
	Delete(id uint, kdCabang uint, kdTerminal uint) error
	FindByID(id uint) (*entity.Dermaga, error)
}

type dermagaService struct {
	dermagaRepo repository.DermagaRepository
}

func NewDermagaService(dermagaRepo repository.DermagaRepository) DermagaService {
	return &dermagaService{dermagaRepo: dermagaRepo}
}

func (s *dermagaService) Create(dermaga *entity.Dermaga) error {
	return s.dermagaRepo.Create(dermaga)
}

func (s *dermagaService) FindAll(kdCabang uint, kdTerminal uint, limit int, offset int) ([]entity.Dermaga, int64, error) {
	return s.dermagaRepo.FindAll(kdCabang, kdTerminal, limit, offset)
}

func (s *dermagaService) Update(id uint, kdCabang uint, kdTerminal uint, dermaga *entity.Dermaga) error {
	existing, err := s.dermagaRepo.FindByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrNotFound
	}

	if existing.KdCabang != kdCabang || existing.KdTerminal != kdTerminal {
		return ErrUnauthorized
	}

	return s.dermagaRepo.Update(id, dermaga)
}

func (s *dermagaService) Delete(id uint, kdCabang uint, kdTerminal uint) error {
	existing, err := s.dermagaRepo.FindByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrNotFound
	}

	if existing.KdCabang != kdCabang || existing.KdTerminal != kdTerminal {
		return ErrUnauthorized
	}

	return s.dermagaRepo.Delete(id)
}

func (s *dermagaService) FindByID(id uint) (*entity.Dermaga, error) {
	return s.dermagaRepo.FindByID(id)
}
