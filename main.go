package main

import (
	"go-api/controllers"
	"go-api/middleware"
	"go-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 初始化
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	db.AutoMigrate(&models.User{})

	// 初始化Gin
	r := gin.Default()

	// 初始化控制器
	authController := controllers.NewAuthController(db)

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}
		protected := api.Group("/protected")
		protected.Use(middleware.JwtAuthMiddleware())
		{
			protected.GET("/test", func(c *gin.Context) {
				userID := c.MustGet("user_id").(uint)
				c.JSON(http.StatusOK, gin.H{
					"message": "Authenticated", "user_id": userID,
				})
			})
		}
	}
	r.Run(":8080")
}
