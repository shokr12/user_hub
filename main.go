package main

import (
	"log"

	"userHub/internal/config"
	apphttp "userHub/internal/web"
	"userHub/internal/service"
	"userHub/internal/store"
	"userHub/pkg/validator"
)

func main() {
	validator.Init()

	// Initialize database connection
	db := config.InitDB()

	// Setup repository and service layers
	userRepo := store.NewUserStore(db)
	userService := service.NewUserService(userRepo)

	// Setup router with all dependencies
	router := apphttp.SetupRouter(userService)

	// Start the server
	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
