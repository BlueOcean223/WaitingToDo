package handlers

import (
	"backend/internal/models/dto"
	"backend/internal/models/vo"
	"backend/internal/services"
	"backend/pkg/myError"
	"backend/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login /auth/login 登录
func (s *AuthHandler) Login(c *gin.Context) {
	var userVo vo.UserVo
	if err := c.ShouldBindJSON(&userVo); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail("", "参数错误", nil))
		return
	}
	email := userVo.Email
	password := userVo.Password

	// 校验参数
	if email == "" || password == "" {
		c.JSON(http.StatusBadRequest, response.Fail("", "邮箱或密码为空", nil))
		return
	}

	userDto, token, err := s.authService.CheckUser(email, password)
	// 错误处理
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Fail("", err.Error(), nil))
		return
	}
	// 邮箱或密码错误
	if userDto == (dto.UserDto{}) {
		c.JSON(http.StatusOK, response.Fail("", "邮箱或密码错误", nil))
		return
	}
	// 登陆成功
	c.JSON(http.StatusOK, response.Success(token, "登陆成功", userDto))

}

// Register /auth/register 注册
func (s *AuthHandler) Register(c *gin.Context) {
	var userVo vo.UserVo
	if err := c.ShouldBindJSON(&userVo); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail("", "参数错误", nil))
		return
	}

	err := s.authService.Register(userVo)
	if err != nil {
		if myError.IsMyError(err) {
			// 自定义错误，验证码错误等
			c.JSON(http.StatusOK, response.Fail("", err.Error(), nil))
		} else {
			// 系统错误
			c.JSON(http.StatusInternalServerError, response.Fail("", err.Error(), nil))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success("", "注册成功", nil))
}

// Forget /auth/forget 重置密码
func (s *AuthHandler) Forget(c *gin.Context) {
	var userVo vo.UserVo
	if err := c.ShouldBindJSON(&userVo); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail("", "参数错误", nil))
		return
	}

	err := s.authService.ForgetPassword(userVo)
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

	c.JSON(http.StatusOK, response.Success("", "重置密码成功", nil))
}

// Captcha /auth/captcha 发送验证码
func (s *AuthHandler) Captcha(c *gin.Context) {
	email := c.Query("email")
	to := []string{email}
	err := s.authService.Captcha(to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Fail("", "发送验证码失败", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success("", "发送验证码成功", nil))
}
