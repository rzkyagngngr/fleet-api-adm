package role

import "context"

type RoleService interface {
	Create(ctx context.Context, role *Role) error
	FindAll(ctx context.Context) ([]Role, error)
	FindByID(ctx context.Context, id uint64) (*Role, error)
	Update(ctx context.Context, id uint64, role *Role) error
	Delete(ctx context.Context, id uint64) error
}

type roleService struct{ roleRepo RoleRepository }

func NewRoleService(roleRepo RoleRepository) RoleService { return &roleService{roleRepo: roleRepo} }
func (s *roleService) Create(ctx context.Context, role *Role) error {
	return s.roleRepo.Create(ctx, role)
}
func (s *roleService) FindAll(ctx context.Context) ([]Role, error) { return s.roleRepo.FindAll(ctx) }
func (s *roleService) FindByID(ctx context.Context, id uint64) (*Role, error) {
	return s.roleRepo.FindByID(ctx, id)
}
func (s *roleService) Update(ctx context.Context, id uint64, role *Role) error {
	return s.roleRepo.Update(ctx, id, role)
}
func (s *roleService) Delete(ctx context.Context, id uint64) error { return s.roleRepo.Delete(ctx, id) }
