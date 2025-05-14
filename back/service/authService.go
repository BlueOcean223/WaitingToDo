package service

import (
	"back/configs"
	"back/models"
	"back/models/dto"
	"back/models/vo"
	"back/repository"
	"back/utils"
	"context"
	"errors"
	"fmt"
	"gopkg.in/gomail.v2"
	"time"
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

// Register 注册
func Register(userVo vo.UserVo) error {
	// 验证邮箱是否已经存在，不允许已存在的邮箱重复注册
	_, err := repository.SelectUserByEmail(userVo.Email)
	if err == nil {
		return errors.New("该邮箱已注册")
	}

	// 校验验证码
	redisClient := configs.RedisClient
	ctx := context.Background()
	val, err := redisClient.Get(ctx, utils.CaptchaKey+userVo.Email).Result()
	if err != nil {
		return err
	}
	if !utils.VerifyCaptcha(val, userVo.Captcha) {
		return errors.New("验证码错误")
	}
	// 密码加密
	hashPassword, err := utils.HashPassword(userVo.Password)
	if err != nil {
		return err
	}
	// 插入数据库
	user := models.User{
		Email:    userVo.Email,
		Name:     userVo.Name,
		Password: hashPassword,
	}
	return repository.InsertUser(user)
}

// ForgetPassword 忘记密码
func ForgetPassword(userVo vo.UserVo) error {
	// 检查数据库中是否有该用户
	user, err := repository.SelectUserByEmail(userVo.Email)
	if err != nil {
		return errors.New("该邮箱尚未注册，请先注册！")
	}
	// 校验验证码
	redisClient := configs.RedisClient
	ctx := context.Background()
	val, err := redisClient.Get(ctx, utils.CaptchaKey+userVo.Email).Result()
	if err != nil {
		return err
	}
	if !utils.VerifyCaptcha(val, userVo.Captcha) {
		return errors.New("验证码错误")
	}
	// 更新密码
	hashPassword, err := utils.HashPassword(userVo.Password)
	if err != nil {
		return err
	}
	user.Password = hashPassword
	return repository.UpdateUser(user)
}

// Captcha 获取验证码
func Captcha(to []string) error {
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
	err := Send163Mail(mail)
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
func Send163Mail(mail models.Mail) error {
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
