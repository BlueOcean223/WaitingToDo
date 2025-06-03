package service

import (
	"back/configs"
	"back/models"
	"back/models/dto"
	"back/models/vo"
	"back/repository"
	"back/utils"
	"context"
	"encoding/json"
	"fmt"
	"gopkg.in/gomail.v2"
	"time"
)

type AuthService struct {
	authRepository *repository.AuthRepository
}

func NewAuthService(authRepository *repository.AuthRepository) *AuthService {
	return &AuthService{
		authRepository: authRepository,
	}
}

// CheckUser 登录校验
// 返回用户信息及JWT令牌
func (s *AuthService) CheckUser(email, password string) (dto.UserDto, string, error) {
	user, err := s.authRepository.SelectUserByEmail(email)
	// 查询用户异常
	if err != nil {
		return dto.UserDto{}, "", err
	}
	// 该邮箱未注册
	if user == (models.User{}) {
		return dto.UserDto{}, "", nil
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

	userDto := dto.UserDto{
		Id:          user.Id,
		Email:       user.Email,
		Name:        user.Name,
		Pic:         user.Pic,
		Description: user.Description,
	}

	// 将用户信息写入redis
	userInfo, err := json.Marshal(userDto)
	if err != nil {
		return dto.UserDto{}, "", err
	}

	emailKey := utils.UserInfoKey + email
	idKey := fmt.Sprintf(utils.UserInfoKey+"%d", user.Id)
	redisClient := configs.RedisClient
	redisClient.Set(context.Background(), emailKey, userInfo, 24*time.Hour)
	redisClient.Set(context.Background(), idKey, userInfo, 24*time.Hour)

	return userDto, token, nil
}

// Register 注册
func (s *AuthService) Register(userVo vo.UserVo) error {
	// 验证邮箱是否已经存在，不允许已存在的邮箱重复注册
	user, err := s.authRepository.SelectUserByEmail(userVo.Email)
	if err != nil {
		return err
	}
	// 邮箱已经被注册
	if user != (models.User{}) {
		return utils.NewMyError("该邮箱已注册")
	}

	// 校验验证码
	err = s.CheckCaptcha(userVo.Email, userVo.Captcha)
	if err != nil {
		return err
	}

	// 密码加密
	hashPassword, err := utils.HashPassword(userVo.Password)
	if err != nil {
		return err
	}
	// 插入数据库
	user = models.User{
		Email:    userVo.Email,
		Name:     userVo.Name,
		Password: hashPassword,
	}
	return s.authRepository.InsertUser(user, nil)
}

// ForgetPassword 忘记密码
func (s *AuthService) ForgetPassword(userVo vo.UserVo) error {
	// 检查数据库中是否有该用户
	user, err := s.authRepository.SelectUserByEmail(userVo.Email)
	// 查询发生异常
	if err != nil {
		return err
	}

	// 用户不存在
	if user == (models.User{}) {
		return utils.NewMyError("该邮箱尚未注册，请先注册！")
	}

	// 校验验证码
	err = s.CheckCaptcha(userVo.Email, userVo.Captcha)
	if err != nil {
		return err
	}

	// 更新密码
	hashPassword, err := utils.HashPassword(userVo.Password)
	if err != nil {
		return err
	}
	user.Password = hashPassword
	return s.authRepository.UpdateUser(user, nil)
}

// Captcha 获取验证码
func (s *AuthService) Captcha(to []string) error {
	// 生成验证码
	code := utils.GenerateCaptcha()
	// 封装邮件
	mail := models.Mail{
		To:      to,
		Subject: "您的验证码",
		Body: fmt.Sprintf(`
		<h2>欢迎使用WaitingToDo</h2>
		<p>您的验证码是: <strong>%s</strong></p>
		<p>验证码5分钟内有效，请勿泄露给他人</p>
		`, code),
	}
	// 发送验证码
	err := s.Send163Mail(mail)
	if err != nil {
		return err
	}
	// 存入redis
	redisClient := configs.RedisClient
	ctx := context.Background()
	err = redisClient.Set(ctx, utils.CaptchaKey+mail.To[0], code, 5*time.Minute).Err()
	return err
}

// Send163Mail 通过163邮箱发送邮件
func (s *AuthService) Send163Mail(mail models.Mail) error {
	config := configs.AppConfigs.MailConfig

	m := gomail.NewMessage()
	m.SetHeader("From", config.From)
	m.SetHeader("To", mail.To...)
	m.SetHeader("Subject", mail.Subject)
	m.SetBody("text/html", mail.Body)

	d := gomail.NewDialer(
		config.SMTPHost,
		config.SMTPPort,
		config.From,
		config.Password,
	)

	return d.DialAndSend(m)
}

// CheckCaptcha 校验验证码
func (s *AuthService) CheckCaptcha(email, captcha string) error {
	// 从redis获取验证码
	redisClient := configs.RedisClient
	ctx := context.Background()
	val, err := redisClient.Get(ctx, utils.CaptchaKey+email).Result()
	if err != nil {
		return err
	}
	// 校验验证码
	if !utils.VerifyCaptcha(val, captcha) {
		return utils.NewMyError("验证码错误")
	}
	// 验证码正确，删除redis中的验证码
	return redisClient.Del(ctx, utils.CaptchaKey+email).Err()
}
