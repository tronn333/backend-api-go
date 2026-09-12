package services

import (
	"backend-api-go/internal/apperr"
	"backend-api-go/internal/models"
	"backend-api-go/internal/repositories"
	"context"
	"fmt"
)

type UserService interface {
	GetProfile(ctx context.Context, userID int) (*models.User, error)
	UpdateProfile(ctx context.Context, userID int, req *models.UpdateProfileRequest) (*models.User, error)
	DeleteAccount(ctx context.Context, userID int) error
}

type userService struct {
	userRepo repositories.UserRepo
}

func NewUserService(userRepo repositories.UserRepo) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetProfile(ctx context.Context, userID int) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, apperr.ErrUserNotFound
	}
	return user, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID int, req *models.UpdateProfileRequest) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, apperr.ErrUserNotFound
	}

	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}
	return user, nil
}

func (s *userService) DeleteAccount(ctx context.Context, userID int) error {
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}
	return nil
}
