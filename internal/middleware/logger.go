package middleware

import (
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggerMiddleware() gin.HandlerFunc {
	logger, _ := zap.NewProduction()
	return ginzap.Ginzap(logger, time.RFC3339, true)
}
