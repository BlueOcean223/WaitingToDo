package security

import (
	"back/utils/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"sync"
	"time"
)

// 请求频率限制
var (
	ipRequestCount = make(map[string][]time.Time)
	ipMutex        = sync.RWMutex{}
	maxRequests    = 100 // 每分钟最大请求数
	timeWindow     = time.Minute
)

// SecurityMiddleware 安全中间件
func SecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查可疑路径
		path := c.Request.URL.Path
		if isSuspiciousPath(path) {
			logger.Warn("可疑路径访问",
				logger.String("ip", c.ClientIP()),
				logger.String("path", path),
				logger.String("user_agent", c.GetHeader("User-Agent")),
				logger.String("method", c.Request.Method))
			c.JSON(http.StatusNotFound, gin.H{"error": "路径不存在"})
			c.Abort()
			return
		}

		// 检查XSS攻击
		if containsXSS(c.Request.URL.RawQuery) {
			logger.Warn("XSS攻击尝试",
				logger.String("ip", c.ClientIP()),
				logger.String("query", c.Request.URL.RawQuery),
				logger.String("path", path),
				logger.String("method", c.Request.Method))
			c.JSON(http.StatusBadRequest, gin.H{"error": "请求被拒绝"})
			c.Abort()
			return
		}

		// 设置安全头
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		c.Next()
	}
}

// RateLimitMiddleware 频率限制中间件
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		ipMutex.Lock()
		now := time.Now()

		// 清理过期记录
		if requests, exists := ipRequestCount[ip]; exists {
			var validRequests []time.Time
			for _, reqTime := range requests {
				if now.Sub(reqTime) < timeWindow {
					validRequests = append(validRequests, reqTime)
				}
			}
			ipRequestCount[ip] = validRequests
		}

		// 检查请求频率
		if len(ipRequestCount[ip]) >= maxRequests {
			ipMutex.Unlock()
			logger.Warn("频率限制触发",
				logger.String("ip", ip),
				logger.Int("request_count", len(ipRequestCount[ip])),
				logger.String("path", c.Request.URL.Path),
				logger.String("method", c.Request.Method))
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "请求过于频繁"})
			c.Abort()
			return
		}

		// 记录当前请求
		ipRequestCount[ip] = append(ipRequestCount[ip], now)
		ipMutex.Unlock()

		c.Next()
	}
}

// 检查可疑路径
func isSuspiciousPath(path string) bool {
	suspiciousPaths := []string{
		"/v1/", "/v2/", "/v3/", "/api/v",
		"/admin", "/wp-admin", "/phpmyadmin",
		"/jobs/", "/chart/", "/templates/",
		"/site/content_store", "/observables",
		"/.env", "/config", "/.git",
	}

	for _, suspicious := range suspiciousPaths {
		if strings.Contains(path, suspicious) {
			return true
		}
	}
	return false
}

// 检查XSS攻击
func containsXSS(query string) bool {
	xssPatterns := []string{
		"<script", "</script>", "javascript:",
		"alert(", "document.domain", "eval(",
		"onload=", "onerror=", "onclick=",
	}

	queryLower := strings.ToLower(query)
	for _, pattern := range xssPatterns {
		if strings.Contains(queryLower, pattern) {
			return true
		}
	}
	return false
}
