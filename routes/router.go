package routes

import (
	"Gin_webDemo/controller"
	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.POST("/api/register", controller.Register)
	return r
}
