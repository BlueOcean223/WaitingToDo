package service

import (
	"back/configs"
	"back/middlewares/jwt"
	"back/models"
	"back/models/dto"
	"back/models/vo"
	"back/repository"
	"back/utils/captcha"
	"back/utils/hashPassword"
	"back/utils/logger"
	"back/utils/myError"
	"back/utils/redisContent"
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
		logger.Error("查询用户失败",
			logger.String("email", email),
			logger.Err(err))
		return dto.UserDto{}, "", err
	}
	// 该邮箱未注册
	if user == (models.User{}) {
		logger.Warn("用户不存在",
			logger.String("email", email))
		return dto.UserDto{}, "", nil
	}

	// 校验密码
	if !hashPassword.CheckPasswordHash(password, user.Password) {
		// 密码错误
		logger.Warn("密码验证失败",
			logger.String("email", email),
			logger.String("user_id", fmt.Sprintf("%d", user.Id)))
		return dto.UserDto{}, "", nil
	}

	// 校验成功，下发令牌
	token, e := jwt.GenerateToken(email, 24*time.Hour)
	if e != nil {
		logger.Error("JWT令牌生成失败",
			logger.String("email", email),
			logger.Err(e))
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
		logger.Error("用户信息JSON序列化失败",
			logger.String("email", email),
			logger.Err(err))
		return dto.UserDto{}, "", err
	}

	emailKey := redisContent.UserInfoKey + email
	idKey := fmt.Sprintf(redisContent.UserInfoKey+"%d", user.Id)
	redisClient := configs.RedisClient
	err = redisClient.Set(context.Background(), emailKey, userInfo, 24*time.Hour).Err()
	if err != nil {
		logger.Error("Redis写入用户信息失败(email key)",
			logger.String("email", email),
			logger.Err(err))
		return dto.UserDto{}, "", err
	}

	err = redisClient.Set(context.Background(), idKey, userInfo, 24*time.Hour).Err()
	if err != nil {
		logger.Error("Redis写入用户信息失败(id key)",
			logger.String("email", email),
			logger.String("user_id", fmt.Sprintf("%d", user.Id)),
			logger.Err(err))
		return dto.UserDto{}, "", err
	}

	return userDto, token, nil
}

// Register 注册
func (s *AuthService) Register(userVo vo.UserVo) error {
	// 验证邮箱是否已经存在，不允许已存在的邮箱重复注册
	user, err := s.authRepository.SelectUserByEmail(userVo.Email)
	if err != nil {
		logger.Error("查询用户邮箱失败",
			logger.String("email", userVo.Email),
			logger.Err(err))
		return err
	}
	// 邮箱已经被注册
	if user != (models.User{}) {
		logger.Warn("邮箱已存在",
			logger.String("email", userVo.Email),
			logger.String("name", userVo.Name))
		return myError.NewMyError("该邮箱已注册")
	}

	// 校验验证码
	err = s.CheckCaptcha(userVo.Email, userVo.Captcha)
	if err != nil {
		return err
	}

	// 密码加密
	hash, err := hashPassword.HashPassword(userVo.Password)
	if err != nil {
		logger.Error("密码加密失败",
			logger.String("email", userVo.Email),
			logger.Err(err))
		return err
	}

	// 插入数据库
	user = models.User{
		Email:    userVo.Email,
		Name:     userVo.Name,
		Password: hash,
	}

	err = s.authRepository.InsertUser(user, nil)
	if err != nil {
		logger.Error("用户注册失败",
			logger.String("email", userVo.Email),
			logger.String("name", userVo.Name),
			logger.Err(err))
		return err
	}

	return nil
}

// ForgetPassword 忘记密码
func (s *AuthService) ForgetPassword(userVo vo.UserVo) error {
	// 检查数据库中是否有该用户
	user, err := s.authRepository.SelectUserByEmail(userVo.Email)
	// 查询发生异常
	if err != nil {
		logger.Error("查询用户失败",
			logger.String("email", userVo.Email),
			logger.Err(err))
		return err
	}

	// 用户不存在
	if user == (models.User{}) {
		logger.Warn("用户不存在，忘记密码失败",
			logger.String("email", userVo.Email))
		return myError.NewMyError("该邮箱尚未注册，请先注册！")
	}

	// 校验验证码
	err = s.CheckCaptcha(userVo.Email, userVo.Captcha)
	if err != nil {
		return err
	}

	// 更新密码
	hash, err := hashPassword.HashPassword(userVo.Password)
	if err != nil {
		logger.Error("密码加密失败",
			logger.String("email", userVo.Email),
			logger.Err(err))
		return err
	}

	user.Password = hash

	err = s.authRepository.UpdateUser(user, nil)
	if err != nil {
		logger.Error("密码更新失败",
			logger.String("email", userVo.Email),
			logger.String("user_id", fmt.Sprintf("%d", user.Id)),
			logger.Err(err))
		return err
	}

	return nil
}

// Captcha 获取验证码
func (s *AuthService) Captcha(to []string) error {
	if len(to) == 0 {
		logger.Error("验证码发送失败：收件人列表为空")
		return myError.NewMyError("收件人列表不能为空")
	}

	// 生成验证码
	code := captcha.GenerateCaptcha()

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
		logger.Error("验证码邮件发送失败",
			logger.String("email", to[0]),
			logger.Err(err))
		return err
	}

	// 存入redis
	redisClient := configs.RedisClient
	ctx := context.Background()
	captchaKey := redisContent.CaptchaKey + mail.To[0]
	err = redisClient.Set(ctx, captchaKey, code, 5*time.Minute).Err()
	if err != nil {
		logger.Error("验证码Redis存储失败",
			logger.String("email", to[0]),
			logger.Err(err))
		return err
	}

	return nil
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
func (s *AuthService) CheckCaptcha(email, captchaCode string) error {
	// 从redis获取验证码
	redisClient := configs.RedisClient
	ctx := context.Background()
	captchaKey := redisContent.CaptchaKey + email
	val, err := redisClient.Get(ctx, captchaKey).Result()
	if err != nil {
		logger.Error("从Redis获取验证码失败",
			logger.String("email", email),
			logger.Err(err))
		return err
	}

	// 校验验证码
	if !captcha.VerifyCaptcha(val, captchaCode) {
		logger.Warn("验证码校验失败",
			logger.String("email", email),
			logger.String("input_code", captchaCode))
		return myError.NewMyError("验证码错误")
	}

	// 验证码正确，删除redis中的验证码
	err = redisClient.Del(ctx, captchaKey).Err()
	if err != nil {
		logger.Error("删除Redis验证码失败",
			logger.String("email", email),
			logger.Err(err))
		return err
	}

	return nil
}
