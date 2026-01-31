package router

import (
	"github.com/gin-gonic/gin"
	"github.com/shester1kov/codemate/internal/gateway/handler"
	"github.com/shester1kov/codemate/internal/gateway/middleware"
	"go.uber.org/zap"
)

// инициализация роутера Gin с маршрутами и middleware
func Setup(logger *zap.Logger, serverMode string) *gin.Engine {
	// debug - подробные логи, release - минимальные логи
	if serverMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// создаем роутер
	r := gin.New()

	// подключаем middleware, порядок важен

	// recovery должен быть первым, чтобы ловить паники в других middleware и хендлерах
	r.Use(middleware.Recovery(logger))
	// logger логгирует все запросы
	r.Use(middleware.Logger(logger))
	// cors разрешает кросс-доменные запросы
	r.Use(middleware.CORS())

	healthHandler := handler.NewHealthHandler(logger)
	queryHandler := handler.NewQueryHandler(logger)

	// healthcheck маршруты не входят в /api
	r.GET("/health", healthHandler.Check)
	r.GET("/ready", healthHandler.Ready)

	v1 := r.Group("/api/v1")
	{
		v1.POST("/query", queryHandler.Query)

		// TODO
		// v1.POST("index", indexHandler.Index)
		// v1.GET("status/:id", indexHandler.Status)
	}

	return r

}
