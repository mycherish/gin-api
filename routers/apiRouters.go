package routers

import (
	"go-api/controllers/api"

	"github.com/gin-gonic/gin"
)

func ApiRoutersInst(r *gin.Engine) {
	apiRouter := r.Group("/api")

	user := apiRouter.Group("/user")
	user.GET("/", api.UserController{}.Index)
	user.POST("/add", api.UserController{}.Add)

	auth := apiRouter.Group("/auth")
	authController := &api.AuthController{}
	auth.POST("/register", authController.Register)
	auth.POST("/login", authController.Login)
}
