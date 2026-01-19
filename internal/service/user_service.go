package service

import (
	"context"

	"userHub/internal/domain"
)

// userService implements domain.UserService
type userService struct {
	repo domain.UserRepository
}

// NewUserService creates a new UserService
func NewUserService(repo domain.UserRepository) domain.UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
    // Enforce uniqueness at the service layer too (works for DB + memory).
    if user == nil {
        return nil, domain.NewInternal("user is nil")
    }

    if existing, err := s.repo.GetByEmail(ctx, user.Email); err == nil && existing != nil {
        return nil, domain.NewConflict("email already exists")
    }

    return s.repo.Create(ctx, user)
}

func (s *userService) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *userService) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
    return s.repo.Update(ctx, user)
}

func (s *userService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *userService) List(ctx context.Context, page, limit int, q string) ([]*domain.User, int64, error) {
	return s.repo.List(ctx, page, limit, q)
}
