package handlers

import (
	"net/http"

	"userHub/internal/domain"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, status int, data any) {
	c.JSON(status, data)
}

func Fail(c *gin.Context, status int, code domain.ErrorCode, message string, fields map[string]string) {
	payload := gin.H{"error": gin.H{"code": code, "message": message}}
	if len(fields) > 0 {
		payload["error"].(gin.H)["fields"] = fields
	}
	c.JSON(status, payload)
}

func FailFromError(c *gin.Context, err error) {
	if err == nil {
		Fail(c, http.StatusInternalServerError, domain.CodeInternal, "unknown error", nil)
		return
	}

	if ae, ok := err.(*domain.AppError); ok {
		switch ae.Code {
		case domain.CodeValidation:
			Fail(c, http.StatusBadRequest, ae.Code, ae.Message, ae.Fields)
			return
		case domain.CodeNotFound:
			Fail(c, http.StatusNotFound, ae.Code, ae.Message, nil)
			return
		case domain.CodeConflict:
			Fail(c, http.StatusConflict, ae.Code, ae.Message, nil)
			return
		default:
			Fail(c, http.StatusInternalServerError, domain.CodeInternal, ae.Message, nil)
			return
		}
	}

	// Fallback
	Fail(c, http.StatusInternalServerError, domain.CodeInternal, err.Error(), nil)
}
