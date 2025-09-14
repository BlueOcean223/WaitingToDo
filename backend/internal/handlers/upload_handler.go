package handlers

import (
	"backend/internal/services"
	"backend/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UploadHandler struct {
	uploadService services.UploadService
}

func NewUploadHandler(uploadService services.UploadService) *UploadHandler {
	return &UploadHandler{
		uploadService: uploadService,
	}
}

// UploadImg /upload/img 上传图片
func (s *UploadHandler) UploadImg(c *gin.Context) {
	// 接收前端传递的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Fail("", err.Error(), nil))
		return
	}

	email := c.GetString("user")

	// 上传文件，并获取url
	url, err := s.uploadService.UploadImg(email, file)
	if err != nil {
		c.JSON(http.StatusOK, response.Fail("", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.Success("", "上传成功", url))
}

// UploadFile 上传文件
func (s *UploadHandler) UploadFile(c *gin.Context) {
	taskId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Fail("", err.Error(), nil))
		return
	}
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Fail("", err.Error(), nil))
		return
	}
	files := form.File["files"]
	err = s.uploadService.UploadFile(taskId, files)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Fail("", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.Success("", "上传成功", nil))
}

// DeleteFile 删除文件
func (s *UploadHandler) DeleteFile(c *gin.Context) {
	type temp struct {
		Ids []int `json:"delete_ids"`
	}
	var ids temp
	if err := c.ShouldBindJSON(&ids); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail("", err.Error(), nil))
		return
	}

	err := s.uploadService.DeleteFile(ids.Ids)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Fail("", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.Success("", "删除成功", nil))
}
