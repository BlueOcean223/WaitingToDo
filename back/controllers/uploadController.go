package controllers

import (
	"back/models"
	"back/service"
	"back/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UploadController struct {
	uploadService *service.UploadService
}

func NewUploadController(uploadService *service.UploadService) *UploadController {
	return &UploadController{
		uploadService: uploadService,
	}
}

// UploadImg /upload/img 上传图片
func (s *UploadController) UploadImg(c *gin.Context) {
	// 接收前端传递的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Fail("", err.Error(), nil))
		return
	}

	email, err := utils.GetUserFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.Fail("", err.Error(), nil))
		return
	}

	// 上传文件，并获取url
	url, err := s.uploadService.UploadImg(email, file)
	if err != nil {
		c.JSON(http.StatusOK, models.Fail("", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, models.Success("", "上传成功", url))
}
