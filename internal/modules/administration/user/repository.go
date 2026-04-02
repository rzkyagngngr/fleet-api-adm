package user

import (
	"context"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id uint64) (*User, error)
	FindByEmployeeID(ctx context.Context, employeeID string) (*User, error)
	GetUserMenusByRole(ctx context.Context, roleID *int) ([]MenuAccessRow, error)
}

type userRepository struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) UserRepository { return &userRepository{db: db} }

func (r *userRepository) Create(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}
func (r *userRepository) FindByID(ctx context.Context, id uint64) (*User, error) {
	var u User
	err := r.db.WithContext(ctx).First(&u, id).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}
func (r *userRepository) FindByEmployeeID(ctx context.Context, employeeID string) (*User, error) {
	var u User
	err := r.db.WithContext(ctx).Where("employee_id = ?", employeeID).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}
func (r *userRepository) GetUserMenusByRole(ctx context.Context, roleID *int) ([]MenuAccessRow, error) {
	var menus []MenuAccessRow
	if roleID == nil {
		return menus, nil
	}
	err := r.db.WithContext(ctx).Raw("SELECT roles_id, menu_id, menu_code, menu_icon, menu_text, menu_url, view, insert, update, delete, menu_level, parent_menu_id FROM vw_access_login WHERE roles_id = ?", *roleID).Scan(&menus).Error
	if err != nil {
		return nil, err
	}
	return menus, nil
}
