package middleware

import "github.com/gin-gonic/gin"

// добавляет заголовки для Cross-Origin Resource Sharing (CORS)
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {

		// разрешаем запросы с любых доменов в dev
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		// разрешенные http методы и заголовки
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Content-Length")

		// разрешаем отправку куки
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// если это preflight запрос, завершаем его сразу - браузер проверяет разрешения
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
