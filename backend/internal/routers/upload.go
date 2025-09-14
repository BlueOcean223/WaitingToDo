package routers

import (
	"backend/internal/configs"
	"backend/internal/handlers"
	"backend/internal/repositories"
	"backend/internal/services"
	"github.com/gin-gonic/gin"
)

func SetLoadRoutes(r *gin.RouterGroup) {
	// 初始化依赖
	authRepository := repository.NewAuthRepository(configs.MysqlDb)
	imageRepository := repository.NewImageRepository(configs.MysqlDb)
	fileRepository := repository.NewFileRepository(configs.MysqlDb)

	uploadService := services.NewUploadService(authRepository, imageRepository, fileRepository)
	uploadController := handlers.NewUploadHandler(uploadService)

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
