package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sahmaragaev/lunaria-backend/internal/errors"
)

type Response struct {
	Status    int    `json:"status"`
	Success   bool   `json:"success"`
	Message   string `json:"message,omitempty"`
	Data      any    `json:"data,omitempty"`
	ErrorCode string `json:"error_code,omitempty"`
	Error     string `json:"error,omitempty"`
	Details   any    `json:"details,omitempty"`
}

func Success(c *gin.Context, data any, message string) {
	resp := Response{
		Status:  http.StatusOK,
		Success: true,
		Data:    data,
		Message: message,
	}

	c.JSON(resp.Status, resp)
}

func Error(c *gin.Context, status int, err error, details any) {
	var errorCode string
	var errorMessage string

	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			errorCode = string(appErr.Code)
			errorMessage = appErr.Error()
		} else {
			errorCode = string(errors.ErrCodeInternalError)
			errorMessage = err.Error()
		}
	} else {
		errorCode = string(errors.ErrCodeInternalError)
		errorMessage = "An error occurred"
	}

	resp := Response{
		Status:    status,
		Success:   false,
		ErrorCode: errorCode,
		Error:     errorMessage,
		Details:   details,
	}

	c.JSON(resp.Status, resp)
}

func BadRequest(c *gin.Context, err error, details any) {
	Error(c, http.StatusBadRequest, err, details)
}

func InternalServerError(c *gin.Context, err error, details any) {
	Error(c, http.StatusInternalServerError, err, details)
}

func Unauthorized(c *gin.Context, err error, details any) {
	Error(c, http.StatusUnauthorized, err, details)
}

func Created(c *gin.Context, data any, message string) {
	resp := Response{
		Status:  http.StatusCreated,
		Success: true,
		Data:    data,
		Message: message,
	}

	c.JSON(resp.Status, resp)
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Forbidden(c *gin.Context, err error, details any) {
	Error(c, http.StatusForbidden, err, details)
}

func NotFound(c *gin.Context, err error, details any) {
	Error(c, http.StatusNotFound, err, details)
}
