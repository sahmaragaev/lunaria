package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

func CORSMiddleware() gin.HandlerFunc {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://lunaria.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	return func(ctx *gin.Context) {
		c.HandlerFunc(ctx.Writer, ctx.Request)
		ctx.Next()
	}
}
