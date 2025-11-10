package main

import (
	"log"
	"os"
	"time"
	"todo-api/internal/database"
	"todo-api/internal/handlers"
	"todo-api/internal/middleware"
	"todo-api/internal/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if  err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	database.ConnectDB()
	
	database.Migrate(&models.User{}, &models.Todo{})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"}, // Add your frontend URLs
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": "Ok",
			"message": "API is Running",
		})
	})

	v1 := router.Group("/api/v1")
	{
		v1.GET("/", func(ctx *gin.Context) {
			ctx.JSON(200,gin.H{
				"message": "Todo API v1",
				"version": "1.0.0",
			})
		})
		
	// API v1 routes group
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication needed)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		// Protected routes (authentication required)
		todos := v1.Group("/todos")
		todos.Use(middleware.AuthMiddleware())  // THIS LINE IS CRITICAL!
		{
			todos.POST("", handlers.CreateTodo)
			todos.GET("", handlers.GetTodos)
			todos.GET("/stats", handlers.GetStats)
			todos.GET("/categories", handlers.GetCategories)
			todos.GET("/:id", handlers.GetTodo)
			todos.PUT("/:id", handlers.UpdateTodo)
			todos.DELETE("/:id", handlers.DeleteTodo)
			todos.PATCH("/:id/toggle", handlers.ToggleComplete)
			todos.POST("/bulk-delete", handlers.BulkDelete)
			todos.POST("/bulk-complete", handlers.BulkComplete)
		}

		// Profile route
		v1.GET("/profile", middleware.AuthMiddleware(), handlers.GetProfile)
	}
	}

	log.Printf("Server starting on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}