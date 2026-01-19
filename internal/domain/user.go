package domain

import "context"

// User is the core domain entity.
// Keep domain types free from transport concerns (no JSON/validation tags).
type User struct {
    ID     uint
    Name   string
    Email  string
    Gender string
}

// UserRepository is the persistence contract.
type UserRepository interface {
    Create(ctx context.Context, user *User) (*User, error)
    GetByID(ctx context.Context, id uint) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
    Update(ctx context.Context, user *User) (*User, error)
    Delete(ctx context.Context, id uint) error
    List(ctx context.Context, page, limit int, q string) (users []*User, total int64, err error)
}

// UserService is the business logic contract.
type UserService interface {
    Create(ctx context.Context, user *User) (*User, error)
    GetByID(ctx context.Context, id uint) (*User, error)
    Update(ctx context.Context, user *User) (*User, error)
    Delete(ctx context.Context, id uint) error
    List(ctx context.Context, page, limit int, q string) (users []*User, total int64, err error)
}
