package web

import (
	"time"

	"userHub/internal/domain"
	"userHub/internal/web/handlers"
	
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures the gin engine with production-ready middleware and routes.
func SetupRouter(userService domain.UserService) *gin.Engine {
	r := gin.New()

	// Standard production middleware
	r.Use(HeadersMiddleware())
	r.Use(LoggerMiddleware())
	r.Use(gin.Recovery())

	// Professional CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Adjust for production environments
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// Health check endpoint
	r.GET("/health", handlers.HealthCheck)

	// Create handler with dependencies
	userHandler := handlers.NewUserHandler(userService)

	// API Versioning Group
	v1 := r.Group("/api/v1")
	{
		v1.POST("/users", userHandler.CreateUser)
		v1.GET("/users/:id", userHandler.GetUser)
		v1.PUT("/users/:id", userHandler.UpdateUser)
		v1.DELETE("/users/:id", userHandler.DeleteUser)
		v1.GET("/users", userHandler.ListUsers)
	}

	return r
}
