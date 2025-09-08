package api

import (
	"go-api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

type AddInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 所有用户
func (u UserController) Index(c *gin.Context) {
	var users []models.Users

	models.DB.Find(&users)

	c.JSON(http.StatusOK, gin.H{
		"result": users,
	})
}

// 添加用户
func (u UserController) Add(c *gin.Context) {
	var input AddInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	t := time.Unix(int64(time.Now().Unix()), 0)
	curentTime := t.Format("2006-01-02 15:04:05")
	inputUser := models.Users{
		Username: input.Username,
		Email:    input.Email,
		AddTime:  curentTime,
	}
	// 加密
	if err := inputUser.HashPassword(input.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	// 创建用户
	if err := models.DB.Create(&inputUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or email already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registration successful"})

}
