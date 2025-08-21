package handlers

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sahmaragaev/lunaria-backend/internal/database/mongodb"
	"github.com/sahmaragaev/lunaria-backend/internal/database/postgres"
	"github.com/sahmaragaev/lunaria-backend/internal/response"
)

type HealthHandler struct {
	PostgresDB *postgres.PostgresDB
	MongoDB    *mongodb.MongoDB
}

func NewHealthHandler(pg *postgres.PostgresDB, mg *mongodb.MongoDB) *HealthHandler {
	return &HealthHandler{
		PostgresDB: pg,
		MongoDB:    mg,
	}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	status := gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"services":  gin.H{},
	}

	if err := h.PostgresDB.DB.Ping(); err != nil {
		status["services"].(gin.H)["postgres"] = "unhealthy"
		status["status"] = "degraded"
	} else {
		status["services"].(gin.H)["postgres"] = "healthy"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := h.MongoDB.Client.Ping(ctx, nil); err != nil {
		status["services"].(gin.H)["mongodb"] = "unhealthy"
		status["status"] = "degraded"
	} else {
		status["services"].(gin.H)["mongodb"] = "healthy"
	}

	if status["status"] == "healthy" {
		response.Success(c, status, "OK")
	} else {
		response.Error(c, 503, nil, status)
	}
}

func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	h.HealthCheck(c)
}

func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	response.Success(c, gin.H{
		"status":    "alive",
		"timestamp": time.Now().UTC(),
	}, "OK")
}
