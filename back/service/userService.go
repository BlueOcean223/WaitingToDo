package service

import (
	"back/models"
	"back/models/vo"
	"back/repository"
	"back/utils"
)

type UserService struct {
	authRepository *repository.AuthRepository
}

func NewUserService(authRepository *repository.AuthRepository) *UserService {
	return &UserService{authRepository: authRepository}
}

// ResetPassword 重置密码
func (s *UserService) ResetPassword(email, password string) error {
	// 检查数据库是否有该邮箱的用户
	user, err := s.authRepository.SelectUserByEmail(email)
	// 查询异常
	if err != nil {
		return err
	}
	// 用户不存在
	if user == (models.User{}) {
		return utils.NewMyError("用户不存在")
	}
	// 更新密码
	hashPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	user.Password = hashPassword
	return s.authRepository.UpdateUser(user)
}

// UpdateUserInfo 更新用户信息
func (s *UserService) UpdateUserInfo(userVo vo.UserVo) error {
	user, err := s.authRepository.SelectUserByEmail(userVo.Email)
	if err != nil {
		return err
	}
	if user == (models.User{}) {
		return utils.NewMyError("用户不存在")
	}

	// 更新个人信息
	user.Name = userVo.Name
	user.Description = userVo.Description
	return s.authRepository.UpdateUser(user)
}
