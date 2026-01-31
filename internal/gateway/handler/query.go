package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shester1kov/codemate/internal/gateway/dto"
	"go.uber.org/zap"
)

// обрабатывает запросы к кодовой базе
type QueryHandler struct {
	logger *zap.Logger
}

func NewQueryHandler(logger *zap.Logger) *QueryHandler {
	return &QueryHandler{
		logger: logger,
	}
}

// обрабатывает post запрос
func (h *QueryHandler) Query(c *gin.Context) {
	var req dto.QueryRequest
	// автоматическая проверка требуемых полей
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request",
			zap.Error(err),
		)

		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: fmt.Sprintf("Invalid request: %s", err.Error()),
		})
		return
	}

	// дополнительная проверка
	if req.Query == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Question cannot be empty",
		})
		return

	}

	h.logger.Info("Received query",
		zap.String("query", req.Query),
		zap.Int("max_results", req.MaxResults),
	)

	// заглушка - здесь будет логика обработки запроса к кодовой базе

	response := dto.QueryResponse{
		Answer: "THIS IS A PLACEHOLDER RESPONSE, RAG PIPELINE WILL BE IMPLEMENTED LATER",
		Sources: []dto.Source{
			{
				FilePath: "example/main.go",
				Name:     "main",
				Type:     "function",
				Score:    0.95,
			},
		},
	}
	c.JSON(http.StatusOK, response)
}
