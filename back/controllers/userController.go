package controllers

import (
	"back/models"
	"back/models/vo"
	"back/service"
	"back/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CheckCaptcha /user/checkCaptcha 检查验证码
func CheckCaptcha(c *gin.Context) {
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

	err = service.CheckCaptcha(email, userVo.Captcha)
	if err != nil {
		if utils.IsMyError(err) {
			// 自定义错误
			c.JSON(http.StatusOK, models.Fail("", err.Error(), nil))
		} else {
			// 系统错误
			c.JSON(http.StatusInternalServerError, models.Fail("", err.Error(), nil))
		}
	}

	c.JSON(http.StatusOK, models.Success("", "验证成功", nil))
}

// Reset /user/reset 重置密码
func Reset(c *gin.Context) {
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

	err = service.ResetPassword(email, userVo.Password)
	if err != nil {
		if utils.IsMyError(err) {
			// 自定义错误
			c.JSON(http.StatusOK, models.Fail("", err.Error(), nil))
		} else {
			// 系统错误
			c.JSON(http.StatusInternalServerError, models.Fail("", err.Error(), nil))
		}
	}

	c.JSON(http.StatusOK, models.Success("", "重置密码成功！", nil))

}

// UpdateUserInfo /user/update 更新用户信息
func UpdateUserInfo(c *gin.Context) {
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
	err = service.UpdateUserInfo(userVo)
	if err != nil {
		if utils.IsMyError(err) {
			// 自定义错误
			c.JSON(http.StatusOK, models.Fail("", err.Error(), nil))
		} else {
			// 系统错误
			c.JSON(http.StatusInternalServerError, models.Fail("", err.Error(), nil))
		}
	}

	c.JSON(http.StatusOK, models.Success("", "更新个人信息成功！", nil))
}
