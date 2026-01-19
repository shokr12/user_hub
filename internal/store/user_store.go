package store

import (
	"context"

	"userHub/internal/domain"

	"gorm.io/gorm"
)

// userStore implements domain.UserRepository
type userStore struct {
	db *gorm.DB
}

// NewUserStore creates a new UserRepository backed by GORM
func NewUserStore(db *gorm.DB) domain.UserRepository {
	return &userStore{db: db}
}

func (s *userStore) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
    if err := s.db.WithContext(ctx).Create(user).Error; err != nil {
        return nil, err
    }
    return user, nil
}

func (s *userStore) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User
    if err := s.db.WithContext(ctx).First(&user, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, domain.NewNotFound("user not found")
        }
        return nil, err
	}
	return &user, nil
}

func (s *userStore) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
    if err := s.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, domain.NewNotFound("user not found")
        }
        return nil, err
	}
	return &user, nil
}

func (s *userStore) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
    if err := s.db.WithContext(ctx).Save(user).Error; err != nil {
        return nil, err
    }
    return user, nil
}

func (s *userStore) Delete(ctx context.Context, id uint) error {
    res := s.db.WithContext(ctx).Delete(&domain.User{}, id)
    if res.Error != nil {
        return res.Error
    }
    if res.RowsAffected == 0 {
        return domain.NewNotFound("user not found")
    }
    return nil
}

func (s *userStore) List(ctx context.Context, page, limit int, q string) ([]*domain.User, int64, error) {
	var users []*domain.User
	var total int64

	query := s.db.WithContext(ctx).Model(&domain.User{})

	if q != "" {
		query = query.Where("name LIKE ? OR email LIKE ?", "%"+q+"%", "%"+q+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
