package api

import (
	"go-api/models"
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

// 注册
func (ac AuthController) Register(c *gin.Context) {
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
	c.JSON(http.StatusCreated, gin.H{"message": "Registration successful"})
}

// 添加用户
// func (u UserController) Add(c *gin.Context) {
// 	var input AddInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	t := time.Unix(int64(time.Now().Unix()), 0)
// 	curentTime := t.Format("2006-01-02 15:04:05")
// 	inputUser := models.Users{
// 		Username: input.Username,
// 		Email:    input.Email,
// 		AddTime:  curentTime,
// 	}
// 	// 加密
// 	if err := inputUser.HashPassword(input.Password); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
// 		return
// 	}
// 	// 创建用户
// 	if err := models.DB.Create(&inputUser).Error; err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or email already exists"})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{"message": "Registration successful"})

// }
