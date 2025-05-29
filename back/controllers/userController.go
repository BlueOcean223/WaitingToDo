package controllers

import (
	"back/models"
	"back/models/vo"
	"back/service"
	"back/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserController struct {
	authService *service.AuthService
	userService *service.UserService
}

func NewUserController(authService *service.AuthService,
	userService *service.UserService) *UserController {
	return &UserController{
		authService: authService,
		userService: userService,
	}
}

// CheckCaptcha /user/checkCaptcha 检查验证码
func (s *UserController) CheckCaptcha(c *gin.Context) {
	var userVo vo.UserVo
	if err := c.ShouldBindJSON(&userVo); err != nil {
		c.JSON(http.StatusBadRequest, models.Fail("", "参数错误", nil))
		return
	}

	// 从token获取用户email
	email, err := utils.GetUserFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.Fail("", "令牌无效", nil))
		return
	}

	err = s.authService.CheckCaptcha(email, userVo.Captcha)
	if err != nil {
		if utils.IsMyError(err) {
			// 自定义错误
			c.JSON(http.StatusOK, models.Fail("", err.Error(), nil))
		} else {
			// 系统错误
			c.JSON(http.StatusInternalServerError, models.Fail("", err.Error(), nil))
		}
		return
	}

	c.JSON(http.StatusOK, models.Success("", "验证成功", nil))
}

// Reset /user/reset 重置密码
func (s *UserController) Reset(c *gin.Context) {
	var userVo vo.UserVo
	if err := c.ShouldBindJSON(&userVo); err != nil {
		c.JSON(http.StatusBadRequest, models.Fail("", "参数错误", nil))
		return
	}

	// 从token获取用户email
	email, err := utils.GetUserFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.Fail("", "令牌无效", nil))
		return
	}

	err = s.userService.ResetPassword(email, userVo.Password)
	if err != nil {
		if utils.IsMyError(err) {
			// 自定义错误
			c.JSON(http.StatusOK, models.Fail("", err.Error(), nil))
		} else {
			// 系统错误
			c.JSON(http.StatusInternalServerError, models.Fail("", err.Error(), nil))
		}
		return
	}

	c.JSON(http.StatusOK, models.Success("", "重置密码成功！", nil))

}

// UpdateUserInfo /user/update 更新用户信息
func (s *UserController) UpdateUserInfo(c *gin.Context) {
	var userVo vo.UserVo
	if err := c.ShouldBindJSON(&userVo); err != nil {
		c.JSON(http.StatusBadRequest, models.Fail("", "参数错误", nil))
		return
	}

	// 从token获取用户email
	email, err := utils.GetUserFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.Fail("", "令牌无效", nil))
		return
	}

	userVo.Email = email
	err = s.userService.UpdateUserInfo(userVo)
	if err != nil {
		if utils.IsMyError(err) {
			// 自定义错误
			c.JSON(http.StatusOK, models.Fail("", err.Error(), nil))
		} else {
			// 系统错误
			c.JSON(http.StatusInternalServerError, models.Fail("", err.Error(), nil))
		}
		return
	}

	c.JSON(http.StatusOK, models.Success("", "更新个人信息成功！", nil))
}

// GetUserInfo /user/info 获取用户信息
func (s *UserController) GetUserInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Fail("", "参数错误", nil))
		return
	}
	user, err := s.userService.GetUserInfo(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Fail("", "系统错误", nil))
		return
	}

	c.JSON(http.StatusOK, models.Success("", "查询成功", user))
}
