package repositories

import (
	"context"
	"errors"
	error2 "user-service/common/error"
	errConstant "user-service/constants/error"
	"user-service/domain/dto"
	"user-service/domain/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

type IUserRepository interface {
	Register(context.Context, *dto.RegisterRequest) (*models.User, error)
	Update(context.Context, *dto.UpdateRequest, string) (*models.User, error)
	FindByUsername(context.Context, string) (*models.User, error)
	FindByEmail(context.Context, string) (*models.User, error)
	FindByUUID(context.Context, string) (*models.User, error)
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Register(ctx context.Context, request *dto.RegisterRequest) (*models.User, error) {
	user := models.User{
		UUID:        uuid.New(),
		Name:        request.Name,
		Username:    request.Username,
		Password:    request.Password,
		PhoneNumber: request.PhoneNumber,
		Email:       request.Email,
		RoleID:      request.RoleID,
	}

	err := r.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return nil, error2.WrapError(errConstant.ErrSQLError)
	}

	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, request *dto.UpdateRequest, uuid string) (*models.User, error) {
	user := models.User{
		Name:        request.Name,
		Username:    request.Username,
		Password:    *request.Password,
		PhoneNumber: request.PhoneNumber,
		Email:       request.Email,
	}

	err := r.db.WithContext(ctx).
		Where("uuid = ?", uuid).
		Updates(&user).Error

	if err != nil {
		return nil, error2.WrapError(errConstant.ErrSQLError)
	}

	return &user, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Preload("Role").
		Where("username = ?", username).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, error2.WrapError(errConstant.ErrUserNotFound)
		}
		return nil, error2.WrapError(errConstant.ErrSQLError)
	}

	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Preload("Role").
		Where("email = ?", email).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, error2.WrapError(errConstant.ErrUserNotFound)
		}
		return nil, error2.WrapError(errConstant.ErrSQLError)
	}

	return &user, nil
}

func (r *UserRepository) FindByUUID(ctx context.Context, uuid string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Preload("Role").
		Where("uuid = ?", uuid).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, error2.WrapError(errConstant.ErrUserNotFound)
		}
		return nil, error2.WrapError(errConstant.ErrSQLError)
	}

	return &user, nil
}
