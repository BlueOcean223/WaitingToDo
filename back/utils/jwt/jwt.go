package jwt

import (
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
func GenerateToken(name string) (string, error) {
	claims := CustomClaims{
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 过期时间
			Issuer:    "WaitingToDo",                                      // 令牌签发者
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTKey)
}

// ParseToken 解析并验证JWT令牌
func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JWTKey, nil
	})

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

// JWTAuthMiddleware JWT认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 放行/auth/*路由
		if strings.Contains(c.Request.URL.Path, "/auth") {
			c.Next()
			return
		}

		// 从Header获取JWT
		tokenString := c.GetHeader("Authorization")
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供Token"})
			c.Abort()
			return
		}

		// 解析Token
		_, err := ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的Token"})
			c.Abort()
			return
		}

		// 放行
		c.Next()
	}
}

// GetUserFromToken 根据Token获取用户名
func GetUserFromToken(c *gin.Context) (string, error) {
	// 获取token
	tokenString := c.GetHeader("Authorization")
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	// 解析token
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Name, nil
}
