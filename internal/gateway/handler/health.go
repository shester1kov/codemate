package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shester1kov/codemate/internal/gateway/dto"
	"go.uber.org/zap"
)

type HealthHandler struct {
	logger *zap.Logger
}

func NewHealthHandler(logger *zap.Logger) *HealthHandler {
	return &HealthHandler{
		logger: logger,
	}
}

func (h *HealthHandler) Check(c *gin.Context) {

	c.JSON(http.StatusOK, dto.HealthResponse{
		Status:  "ok",
		Version: "1.0.0",
	})
}

// ready - проверяет готовность сервиса к работе (подключения к бд, quadrant и т.д.)
func (h *HealthHandler) Ready(c *gin.Context) {
	// todo - бд, quadrant, ollama, пока просто ok

	c.JSON(http.StatusOK, dto.HealthResponse{
		Status:  "ok",
		Version: "1.0.0",
	})
}
