package service

import (
	"back/models/dto"
	"back/repository"
	"back/utils"
)

// CheckUser 登录校验
// 返回用户信息及JWT令牌
func CheckUser(email, password string) (dto.UserDto, string, error) {
	user, err := repository.SelectUserByEmail(email)
	// 查询用户异常
	if err != nil {
		return dto.UserDto{}, "", err
	}

	// 校验密码
	if !utils.CheckPasswordHash(password, user.Password) {
		// 密码错误
		return dto.UserDto{}, "", nil
	}

	// 校验成功，下发令牌
	token, e := utils.GenerateToken(email)
	if e != nil {
		return dto.UserDto{}, "", e
	}

	return dto.UserDto{
		Email:       user.Email,
		Name:        user.Name,
		Pic:         user.Pic,
		Description: user.Description,
	}, token, nil
}
