package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/shester1kov/codemate/internal/gateway/dto"
	"go.uber.org/zap"
)

// middleware для восстановления после паники и логгирования ошибок
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			// recover возвращает значение, переданное в panic()
			if err := recover(); err != nil {

				// получаем stack trace
				stack := debug.Stack()

				// логируем критическую ошибку
				logger.Error("Panic recovered",
					zap.String("error", fmt.Sprintf("%v", err)),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
					zap.String("stack", string(stack)),
				)

				c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error: "Internal Server Error",
				})

				c.Abort()
			}

		}()

		c.Next()
	}
}
