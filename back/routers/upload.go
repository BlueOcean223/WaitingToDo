package routers

import (
	"back/configs"
	"back/controllers"
	"back/repository"
	"back/service"
	"github.com/gin-gonic/gin"
)

func SetLoadRoutes(r *gin.Engine) {
	// 初始化依赖
	authRepository := repository.NewAuthRepository(configs.MysqlDb)
	imageRepository := repository.NewImageRepository(configs.MysqlDb)
	fileRepository := repository.NewFileRepository(configs.MysqlDb)

	uploadService := service.NewUploadService(authRepository, imageRepository, fileRepository)
	uploadController := controllers.NewUploadController(uploadService)

	load := r.Group("/upload")
	{
		// 上传图片
		load.POST("/img", uploadController.UploadImg)
		// 上传文件
		load.POST(":id/file", uploadController.UploadFile)
		// 删除文件
		load.DELETE("/deletefile", uploadController.DeleteFile)
	}
}
