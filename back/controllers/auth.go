package controllers

import (
	"back/models"
	"back/models/dto"
	"back/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Login /auth/login
func Login(c *gin.Context) {
	email := c.Query("email")
	password := c.Query("password")
	// 校验参数
	if email == "" || password == "" {
		c.JSON(http.StatusBadRequest, models.Fail("", "邮箱或密码为空", nil))
		return
	}

	userDto, token, err := service.CheckUser(email, password)
	// 错误处理
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Fail("", err.Error(), nil))
		return
	}
	// 密码错误
	if userDto == (dto.UserDto{}) {
		c.JSON(http.StatusUnauthorized, models.Fail("", "邮箱或密码错误", nil))
		return
	}
	// 登陆成功
	c.JSON(http.StatusOK, models.Success(token, "登陆成功", userDto))

}
