package controllers

import (
	"back/models"
	"back/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

	email := c.GetString("user")

	// 上传文件，并获取url
	url, err := s.uploadService.UploadImg(email, file)
	if err != nil {
		c.JSON(http.StatusOK, models.Fail("", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, models.Success("", "上传成功", url))
}

// UploadFile 上传文件
func (s *UploadController) UploadFile(c *gin.Context) {
	taskId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Fail("", err.Error(), nil))
		return
	}
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Fail("", err.Error(), nil))
		return
	}
	files := form.File["files"]
	err = s.uploadService.UploadFile(taskId, files)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Fail("", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, models.Success("", "上传成功", nil))
}

// DeleteFile 删除文件
func (s *UploadController) DeleteFile(c *gin.Context) {
	type temp struct {
		Ids []int `json:"delete_ids"`
	}
	var ids temp
	if err := c.ShouldBindJSON(&ids); err != nil {
		c.JSON(http.StatusBadRequest, models.Fail("", err.Error(), nil))
		return
	}

	err := s.uploadService.DeleteFile(ids.Ids)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Fail("", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, models.Success("", "删除成功", nil))
}
