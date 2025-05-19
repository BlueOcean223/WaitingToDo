package routers

import (
	"back/controllers"
	"github.com/gin-gonic/gin"
)

func SetLoadRoutes(r *gin.Engine) {
	load := r.Group("/upload")
	{
		// 上传图片
		load.POST("/img", controllers.UploadImg)
	}
}
