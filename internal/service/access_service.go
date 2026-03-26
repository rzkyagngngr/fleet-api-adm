package service

import (
	"gin-boilerplate/internal/model/entity"
	"gin-boilerplate/internal/repository"
)

type AccessService interface {
	GetRoleAccess(roleID uint64) ([]entity.Access, error)
	UpdateRoleAccess(roleID uint64, accessList []entity.Access) error
	InitializeRoleAccess(roleID uint64) error
	InitializeMenuAccessForAllRoles(menu *entity.Menu) error
}

type accessService struct {
	accessRepo repository.AccessRepository
	menuRepo   repository.MenuRepository
}

func NewAccessService(accessRepo repository.AccessRepository, menuRepo repository.MenuRepository) AccessService {
	return &accessService{
		accessRepo: accessRepo,
		menuRepo:   menuRepo,
	}
}

func (s *accessService) GetRoleAccess(roleID uint64) ([]entity.Access, error) {
	return s.accessRepo.FindByRoleID(roleID)
}

func (s *accessService) UpdateRoleAccess(roleID uint64, accessList []entity.Access) error {
	// Simple implementation: delete existing and bulk create
	// Or individually update. For bulk, deleting and re-inserting is often cleaner for matrix data.
	err := s.accessRepo.DeleteByRoleID(roleID)
	if err != nil {
		return err
	}
	return s.accessRepo.BulkCreate(accessList)
}

func (s *accessService) InitializeRoleAccess(roleID uint64) error {
	menus, err := s.menuRepo.FindAll()
	if err != nil {
		return err
	}

	var accessList []entity.Access
	for _, menu := range menus {
		appID := int64(0)
		if menu.ApplicationID != nil {
			appID = int64(*menu.ApplicationID)
		}
		
		parentID := int64(0)
		if menu.ParentMenuID != nil {
			parentID = int64(*menu.ParentMenuID)
		}

		zero := int16(0)
		
		access := entity.Access{
			RolesID:       ptrInt64(int64(roleID)),
			MenuID:        ptrInt64(int64(menu.ID)),
			MenuText:      menu.MenuText,
			MenuUrl:       menu.MenuUrl,
			Status:        &menu.MenuStatus,
			ApplicationID: &appID,
			ParentMenuID:  &parentID,
			CanInsert:     &zero,
			CanUpdate:     &zero,
			CanDelete:     &zero,
			MenuIcon:      menu.MenuIcon,
		}
		accessList = append(accessList, access)
	}

	if len(accessList) > 0 {
		return s.accessRepo.BulkCreate(accessList)
	}
	return nil
}

func (s *accessService) InitializeMenuAccessForAllRoles(menu *entity.Menu) error {
	roles, err := s.accessRepo.FindAllRoles() // I need to add this to repo
	if err != nil {
		return err
	}

	var accessList []entity.Access
	for _, role := range roles {
		appID := int64(0)
		if menu.ApplicationID != nil {
			appID = int64(*menu.ApplicationID)
		}
		
		parentID := int64(0)
		if menu.ParentMenuID != nil {
			parentID = int64(*menu.ParentMenuID)
		}

		zero := int16(0)
		
		access := entity.Access{
			RolesID:       ptrInt64(int64(role.HakAksesID)),
			MenuID:        ptrInt64(int64(menu.ID)),
			MenuText:      menu.MenuText,
			MenuUrl:       menu.MenuUrl,
			Status:        &menu.MenuStatus,
			ApplicationID: &appID,
			ParentMenuID:  &parentID,
			CanInsert:     &zero,
			CanUpdate:     &zero,
			CanDelete:     &zero,
			MenuIcon:      menu.MenuIcon,
		}
		accessList = append(accessList, access)
	}

	if len(accessList) > 0 {
		return s.accessRepo.BulkCreate(accessList)
	}
	return nil
}

func ptrInt64(i int64) *int64 {
	return &i
}
