package repository

import (
	"gin-boilerplate/internal/model/dto"
	"gin-boilerplate/internal/model/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
	FindByID(id uint64) (*entity.User, error)
	FindByEmployeeID(employeeID string) (*entity.User, error)
	GetUserMenusByRole(roleID *int) ([]dto.MenuAccessRow, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uint64) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmployeeID(employeeID string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("employee_id = ?", employeeID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserMenusByRole(roleID *int) ([]dto.MenuAccessRow, error) {
	var menus []dto.MenuAccessRow
	if roleID == nil {
		return menus, nil
	}
	err := r.db.Raw("SELECT roles_id, menu_id, menu_code, menu_icon, menu_text, menu_url, view, insert, update, delete, menu_level, parent_menu_id FROM vw_access_login WHERE roles_id = ?", *roleID).Scan(&menus).Error
	if err != nil {
		return nil, err
	}
	return menus, nil
}
