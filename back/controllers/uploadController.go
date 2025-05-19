package controllers

import (
	"back/models"
	"back/service"
	"back/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UploadImg /upload/img 上传图片
func UploadImg(c *gin.Context) {
	// 接收前端传递的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Fail("", err.Error(), nil))
	}

	email, err := utils.GetUserFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.Fail("", err.Error(), nil))
	}

	// 上传文件，并获取url
	url, err := service.UploadImg(email, file)
	if err != nil {
		c.JSON(http.StatusOK, models.Fail("", err.Error(), nil))
	}

	c.JSON(http.StatusOK, models.Success("", "上传成功", url))
}
