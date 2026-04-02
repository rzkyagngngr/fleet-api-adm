package access

import "context"

type AccessService interface {
	GetRoleAccess(ctx context.Context, roleID uint64) ([]Access, error)
	UpdateRoleAccess(ctx context.Context, roleID uint64, accessList []Access) error
}

type accessService struct{ accessRepo AccessRepository }

func NewAccessService(accessRepo AccessRepository) AccessService {
	return &accessService{accessRepo: accessRepo}
}
func (s *accessService) GetRoleAccess(ctx context.Context, roleID uint64) ([]Access, error) {
	return s.accessRepo.FindByRoleID(ctx, roleID)
}
func (s *accessService) UpdateRoleAccess(ctx context.Context, roleID uint64, accessList []Access) error {
	if err := s.accessRepo.DeleteByRoleID(ctx, roleID); err != nil {
		return err
	}
	if len(accessList) == 0 {
		return nil
	}
	return s.accessRepo.BulkCreate(ctx, accessList)
}
