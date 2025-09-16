package api

import (
	"go-api/models"
	"go-api/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct{}

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

// 注册
func (ac *AuthController) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t := time.Unix(int64(time.Now().Unix()), 0)
	curentTime := t.Format("2006-01-02 15:04:05")
	user := models.Users{
		Username: input.Username,
		Password: input.Password,
		Email:    input.Email,
		AddTime:  curentTime,
	}
	if err := user.HashPassword(input.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	if err := models.DB.Create(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful",
		"ID":      user.ID,
	})
}

// 登陆
func (ac *AuthController) Login(c *gin.Context) {
	var input LoginInput
	// 绑定输入数据
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查询用户
	var user models.Users
	if err := models.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
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
