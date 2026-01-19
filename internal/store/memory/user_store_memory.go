package memory

import (
    "context"
    "strings"
    "sync"

    "userHub/internal/domain"
)

// userStore is an in-memory implementation of domain.UserRepository.
// Useful for tests and local development.
type userStore struct {
    mu     sync.RWMutex
    nextID uint
    users  map[uint]domain.User
}

func NewUserStore() domain.UserRepository {
    return &userStore{
        nextID: 1,
        users:  make(map[uint]domain.User),
    }
}

func (s *userStore) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    // Enforce unique email
    for _, u := range s.users {
        if strings.EqualFold(u.Email, user.Email) {
            return nil, domain.NewConflict("email already exists")
        }
    }

    u := *user
    u.ID = s.nextID
    s.nextID++
    s.users[u.ID] = u
    return &u, nil
}

func (s *userStore) GetByID(ctx context.Context, id uint) (*domain.User, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    u, ok := s.users[id]
    if !ok {
        return nil, domain.NewNotFound("user not found")
    }
    return &u, nil
}

func (s *userStore) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    for _, u := range s.users {
        if strings.EqualFold(u.Email, email) {
            uu := u
            return &uu, nil
        }
    }
    return nil, domain.NewNotFound("user not found")
}

func (s *userStore) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    if _, ok := s.users[user.ID]; !ok {
        return nil, domain.NewNotFound("user not found")
    }

    // If email is being updated in the future, enforce uniqueness here.
    u := *user
    s.users[u.ID] = u
    return &u, nil
}

func (s *userStore) Delete(ctx context.Context, id uint) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    if _, ok := s.users[id]; !ok {
        return domain.NewNotFound("user not found")
    }
    delete(s.users, id)
    return nil
}

func (s *userStore) List(ctx context.Context, page, limit int, q string) ([]*domain.User, int64, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    if page < 1 {
        page = 1
    }
    if limit < 1 {
        limit = 10
    }
    if limit > 100 {
        limit = 100
    }

    qq := strings.ToLower(strings.TrimSpace(q))

    // Filter
    filtered := make([]domain.User, 0, len(s.users))
    for _, u := range s.users {
        if qq == "" {
            filtered = append(filtered, u)
            continue
        }
        if strings.Contains(strings.ToLower(u.Name), qq) || strings.Contains(strings.ToLower(u.Email), qq) {
            filtered = append(filtered, u)
        }
    }

    total := int64(len(filtered))

    // Paginate
    start := (page - 1) * limit
    if start >= len(filtered) {
        return []*domain.User{}, total, nil
    }
    end := start + limit
    if end > len(filtered) {
        end = len(filtered)
    }

    out := make([]*domain.User, 0, end-start)
    for i := start; i < end; i++ {
        u := filtered[i]
        out = append(out, &u)
    }
    return out, total, nil
}
