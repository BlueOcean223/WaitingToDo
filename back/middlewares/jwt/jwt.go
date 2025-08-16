package jwt

import (
	"back/utils/logger"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

// JWTKey JWT密钥
var JWTKey = []byte("WaitingToDo")

// CustomClaims 自定义JWTClaims
type CustomClaims struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(name string, duration time.Duration) (string, error) {
	claims := CustomClaims{
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)), // 过期时间
			Issuer:    "WaitingToDo",                                // 令牌签发者
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTKey)
}

// ParseToken 解析并验证JWT令牌
func ParseToken(tokenString string) (*CustomClaims, error) {
	if tokenString == "" {
		return nil, errors.New("token不能为空")
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 校验签名方法是否正确
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
		}
		return JWTKey, nil
	})

	if err != nil {
		return nil, err
	}
	// 防止token为空，触发panic
	if token == nil {
		return nil, errors.New("无效的token")
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("无效的token")
	}
}

// JWTAuthMiddleware JWT认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Header获取JWT
		tokenString := c.GetHeader("Authorization")
		// 校验Token格式
		if !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token格式错误"})
			c.Abort()
			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供Token"})
			c.Abort()
			return
		}

		// 解析Token
		claims, err := ParseToken(tokenString)
		if err != nil {
			// 区分过期错误和其它错误
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token已过期"})
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "无效Token"})
			}
			return
		}

		// 检查token剩余时间(12小时)
		remaining := time.Until(claims.ExpiresAt.Time)
		refresh := 12 * time.Hour

		// 令牌即将过期时生成新令牌
		if remaining < refresh {
			// 新令牌有效期为三天
			newToken, err := GenerateToken(claims.Name, 72*time.Hour)
			if err != nil {
				logger.Error("令牌刷新失败",
					logger.String("user", claims.Name),
					logger.Err(err),
					logger.String("ip", c.ClientIP()),
					logger.String("path", c.Request.URL.Path))
			} else {
				// 设置新令牌到响应头
				c.Header("New-Access-Token", newToken)
				logger.Info("令牌刷新成功",
					logger.String("user", claims.Name),
					logger.String("ip", c.ClientIP()),
					logger.String("path", c.Request.URL.Path))
			}
		}

		// 将用户信息存入上下文
		c.Set("user", claims.Name)
		// 放行
		c.Next()
	}
}
