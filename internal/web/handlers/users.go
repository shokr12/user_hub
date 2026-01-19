package handlers

import (
	"net/http"
	"strconv"

	"userHub/internal/domain"
	"userHub/internal/web/dto"
	"userHub/pkg/validator"

	"github.com/gin-gonic/gin"
)

// UserHandler holds dependencies for user-related HTTP handlers
type UserHandler struct {
	userService domain.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(svc domain.UserService) *UserHandler {
	return &UserHandler{userService: svc}
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, domain.CodeValidation, "invalid JSON body", nil)
		return
	}
	if err := validator.Validate(req); err != nil {
		Fail(c, http.StatusBadRequest, domain.CodeValidation, "validation failed", validator.ErrorMap(err))
		return
	}

	created, err := h.userService.Create(c.Request.Context(), &domain.User{
		Name:   req.Name,
		Email:  req.Email,
		Gender: req.Gender,
	})
	if err != nil {
		FailFromError(c, err)
		return
	}

	Success(c, http.StatusCreated, dto.UserResponse{
		ID:     created.ID,
		Name:   created.Name,
		Email:  created.Email,
		Gender: created.Gender,
	})
}

// GetUser handles GET /users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, domain.CodeValidation, "invalid user ID", nil)
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		FailFromError(c, err)
		return
	}

	Success(c, http.StatusOK, dto.UserResponse{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Gender: user.Gender,
	})
}

// UpdateUser handles PUT /users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, domain.CodeValidation, "invalid user ID", nil)
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, domain.CodeValidation, "invalid JSON body", nil)
		return
	}
	if err := validator.Validate(req); err != nil {
		Fail(c, http.StatusBadRequest, domain.CodeValidation, "validation failed", validator.ErrorMap(err))
		return
	}

	// Fetch existing user to apply partial updates.
	existing, err := h.userService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		FailFromError(c, err)
		return
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Gender != nil {
		existing.Gender = *req.Gender
	}

	updated, err := h.userService.Update(c.Request.Context(), existing)
	if err != nil {
		FailFromError(c, err)
		return
	}

	Success(c, http.StatusOK, dto.UserResponse{
		ID:     updated.ID,
		Name:   updated.Name,
		Email:  updated.Email,
		Gender: updated.Gender,
	})
}

// DeleteUser handles DELETE /users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, domain.CodeValidation, "invalid user ID", nil)
		return
	}

	if err := h.userService.Delete(c.Request.Context(), uint(id)); err != nil {
		FailFromError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ListUsers handles GET /users
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	q := c.Query("q")

	users, total, err := h.userService.List(c.Request.Context(), page, limit, q)
	if err != nil {
		FailFromError(c, err)
		return
	}

	resp := dto.ListUsersResponse{}
	resp.Data = make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		resp.Data = append(resp.Data, dto.UserResponse{ID: u.ID, Name: u.Name, Email: u.Email, Gender: u.Gender})
	}
	resp.Meta.Page = page
	resp.Meta.Limit = limit
	resp.Meta.Total = total

	Success(c, http.StatusOK, resp)
}
