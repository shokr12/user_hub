package dto

// Transport-layer DTOs for HTTP requests/responses.

type CreateUserRequest struct {
    Name   string `json:"name" validate:"required,min=2,max=50"`
    Email  string `json:"email" validate:"required,email"`
    Gender string `json:"gender" validate:"required,gender"`
}

type UpdateUserRequest struct {
    Name   *string `json:"name" validate:"omitempty,min=2,max=50"`
    Gender *string `json:"gender" validate:"omitempty,gender"`
}

type UserResponse struct {
    ID     uint   `json:"id"`
    Name   string `json:"name"`
    Email  string `json:"email"`
    Gender string `json:"gender"`
}

type ListUsersResponse struct {
    Data []UserResponse `json:"data"`
    Meta struct {
        Page  int   `json:"page"`
        Limit int   `json:"limit"`
        Total int64 `json:"total"`
    } `json:"meta"`
}
