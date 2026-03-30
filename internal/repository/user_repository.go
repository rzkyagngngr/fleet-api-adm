package repository

import (
	"context"
	"gin-boilerplate/internal/model/dto"
	"gin-boilerplate/internal/model/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByID(ctx context.Context, id uint64) (*entity.User, error)
	FindByEmployeeID(ctx context.Context, employeeID string) (*entity.User, error)
	GetUserMenusByRole(ctx context.Context, roleID *int) ([]dto.MenuAccessRow, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id uint64) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmployeeID(ctx context.Context, employeeID string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("employee_id = ?", employeeID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserMenusByRole(ctx context.Context, roleID *int) ([]dto.MenuAccessRow, error) {
	var menus []dto.MenuAccessRow
	if roleID == nil {
		return menus, nil
	}
	err := r.db.WithContext(ctx).Raw("SELECT roles_id, menu_id, menu_code, menu_icon, menu_text, menu_url, view, insert, update, delete, menu_level, parent_menu_id FROM vw_access_login WHERE roles_id = ?", *roleID).Scan(&menus).Error
	if err != nil {
		return nil, err
	}
	return menus, nil
}
