package routers

import (
	"go-api/controllers/api"

	"github.com/gin-gonic/gin"
)

func ApiRoutersInst(r *gin.Engine) {
	apiRouter := r.Group("/api")
	apiRouter.GET("/users", api.UserController{}.Index)
	apiRouter.POST("/users/add", api.UserController{}.Add)
}
