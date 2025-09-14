package handlers

import (
	"backend/internal/models/vo"
	"backend/internal/services"
	"backend/pkg/myError"
	"backend/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserHandler struct {
	authService services.AuthService
	userService services.UserService
}

func NewUserHandler(authService services.AuthService,
	userService services.UserService) *UserHandler {
	return &UserHandler{
		authService: authService,
		userService: userService,
	}
}

// CheckCaptcha /user/checkCaptcha 检查验证码
func (s *UserHandler) CheckCaptcha(c *gin.Context) {
	var userVo vo.UserVo
	if err := c.ShouldBindJSON(&userVo); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail("", "参数错误", nil))
		return
	}

	// 从token获取用户email
	email := c.GetString("user")

	err := s.authService.CheckCaptcha(email, userVo.Captcha)
	if err != nil {
		if myError.IsMyError(err) {
			// 自定义错误
			c.JSON(http.StatusOK, response.Fail("", err.Error(), nil))
		} else {
			// 系统错误
			c.JSON(http.StatusInternalServerError, response.Fail("", err.Error(), nil))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success("", "验证成功", nil))
}

// Reset /user/reset 重置密码
func (s *UserHandler) Reset(c *gin.Context) {
	var userVo vo.UserVo
	if err := c.ShouldBindJSON(&userVo); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail("", "参数错误", nil))
		return
	}

	// 从token获取用户email
	email := c.GetString("user")

	err := s.userService.ResetPassword(email, userVo.Password)
	if err != nil {
		if myError.IsMyError(err) {
			// 自定义错误
			c.JSON(http.StatusOK, response.Fail("", err.Error(), nil))
		} else {
			// 系统错误
			c.JSON(http.StatusInternalServerError, response.Fail("", err.Error(), nil))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success("", "重置密码成功！", nil))

}

// UpdateUserInfo /user/update 更新用户信息
func (s *UserHandler) UpdateUserInfo(c *gin.Context) {
	var userVo vo.UserVo
	if err := c.ShouldBindJSON(&userVo); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail("", "参数错误", nil))
		return
	}

	// 从token获取用户email
	email := c.GetString("user")

	userVo.Email = email
	err := s.userService.UpdateUserInfo(userVo)
	if err != nil {
		if myError.IsMyError(err) {
			// 自定义错误
			c.JSON(http.StatusOK, response.Fail("", err.Error(), nil))
		} else {
			// 系统错误
			c.JSON(http.StatusInternalServerError, response.Fail("", err.Error(), nil))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success("", "更新个人信息成功！", nil))
}

// GetUserInfo /user/info 获取用户信息
func (s *UserHandler) GetUserInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Fail("", "参数错误", nil))
		return
	}
	user, err := s.userService.GetUserInfo(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Fail("", "系统错误", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success("", "查询成功", user))
}
