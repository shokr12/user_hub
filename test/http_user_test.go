package http_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	apphttp "userHub/internal/web"
	"userHub/internal/service"
	"userHub/internal/store/memory"
)

func TestCreateUser(t *testing.T) {
	// Use an in-memory store for fast + reliable unit testing
	userStore := memory.NewUserStore()
	userService := service.NewUserService(userStore)

	// Build router (this should accept the service OR build handlers using it internally)
	r := apphttp.SetupRouter(userService)

	body := `{"name":"John Doe","email":"john@example.com","gender":"male"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Creating a resource should return 201
	assert.Equal(t, http.StatusCreated, w.Code)

	// Gin usually returns "application/json; charset=utf-8"
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	// Optional: assert response body contains created user fields
	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	// These keys depend on your response format; adjust if your API wraps data differently
	// Example if you return {"id":"...","name":"...","email":"...","gender":"..."}
	assert.Equal(t, "John Doe", resp["name"])
	assert.Equal(t, "john@example.com", resp["email"])
	assert.Equal(t, "male", resp["gender"])
}
