package router

import (
	"pplx2api/middleware"
	"pplx2api/service"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Apply middleware
	r.Use(middleware.CORSMiddleware())

	// Frontend UI
	r.Static("/static", "./web")
	r.GET("/", func(c *gin.Context) {
		c.File("./web/index.html")
	})

	// Health check endpoint
	r.GET("/health", service.HealthCheckHandler)

	// API endpoints
	api := r.Group("/v1")
	api.Use(middleware.AuthMiddleware())
	api.POST("/chat/completions", service.ChatCompletionsHandler)
	api.GET("/models", service.ModelsHandler)

	// Admin config endpoints
	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware())
	admin.GET("/config", service.AdminConfigGetHandler)
	admin.POST("/config", service.AdminConfigUpdateHandler)
	// HuggingFace compatible routes
	hfRouter := r.Group("/hf")
	hfRouter.Use(middleware.AuthMiddleware())
	{
		v1Router := hfRouter.Group("/v1")
		{
			v1Router.POST("/chat/completions", service.ChatCompletionsHandler)
			v1Router.GET("/models", service.ModelsHandler)
		}
	}
}
