package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sahmaragaev/lunaria-backend/internal/errors"
	"github.com/sahmaragaev/lunaria-backend/internal/response"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		if err, ok := recovered.(string); ok {
			response.InternalServerError(c, errors.NewAppError(errors.ErrCodeInternalError, err, nil), nil)
		}
		c.AbortWithStatus(500)
	})
}
