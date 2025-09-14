package controllers

import (
	"go-api/models"
	"go-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

// NewAuthController creates a new instance of AuthController with the provided database connection.
func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{
		DB: db,
	}
}

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

// LoginInput represents the input data for user login.
type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register registers a new user.
// func (ac *AuthController) Register(c *gin.Context) {
// 	var input RegisterInput
// 	// Bind the input data
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	user := models.User{
// 		Username: input.Username,
// 		Email:    input.Email,
// 	}
// 	// Hash the password before saving
// 	if err := user.HashPassword(input.Password); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
// 		return
// 	}
// 	if err := ac.DB.Create(&user).Error; err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or email already exists"})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{"message": "Registration successful"})
// }

// 登陆
func (ac *AuthController) Login(c *gin.Context) {
	var input LoginInput
	// 绑定输入数据
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查询用户
	var user models.User
	if err := ac.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	}

	// Check password
	if err := user.CheckPassword(input.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateToken(uint(user.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}
