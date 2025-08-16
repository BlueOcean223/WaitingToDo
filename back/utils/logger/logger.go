package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var Logger *zap.Logger

// LogConfig 日志配置结构
type LogConfig struct {
	Level      string `yaml:"level"`       // 日志级别: debug, info, warn, error
	Filename   string `yaml:"filename"`    // 日志文件名
	MaxSize    int    `yaml:"max_size"`    // 单个日志文件最大大小(MB)
	MaxBackups int    `yaml:"max_backups"` // 保留的日志文件数量
	MaxAge     int    `yaml:"max_age"`     // 日志文件保留天数
	Compress   bool   `yaml:"compress"`    // 是否压缩
	Console    bool   `yaml:"console"`     // 是否输出到控制台
}

// InitLogger 初始化日志系统
func InitLogger(config LogConfig) error {
	// 设置默认配置
	defaultConfig := LogConfig{
		Level:      "info",
		Filename:   "logs/app.log",
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   true,
		Console:    true,
	}

	// 如果传入的配置有效，使用传入的配置
	if config.Level != "" {
		defaultConfig = config
	}

	// 解析日志级别
	level := parseLogLevel(defaultConfig.Level)

	// 创建编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 创建核心配置
	var cores []zapcore.Core

	// 文件输出
	if defaultConfig.Filename != "" {
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   defaultConfig.Filename,
			MaxSize:    defaultConfig.MaxSize,
			MaxBackups: defaultConfig.MaxBackups,
			MaxAge:     defaultConfig.MaxAge,
			Compress:   defaultConfig.Compress,
		})
		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			fileWriter,
			level,
		)
		cores = append(cores, fileCore)
	}

	// 控制台输出
	if defaultConfig.Console {
		consoleCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			level,
		)
		cores = append(cores, consoleCore)
	}

	// 创建日志器
	core := zapcore.NewTee(cores...)
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return nil
}

// parseLogLevel 解析日志级别
func parseLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// String 创建字符串字段
func String(key, val string) zap.Field {
	return zap.String(key, val)
}

// Int 创建整数字段
func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

// Err 创建错误字段
func Err(err error) zap.Field {
	return zap.Error(err)
}

// Debug 调试日志
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

// Debugf 格式化调试日志
func Debugf(format string, args ...interface{}) {
	Logger.Debug(fmt.Sprintf(format, args...))
}

// Info 信息日志
func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

// Infof 格式化信息日志
func Infof(format string, args ...interface{}) {
	Logger.Info(fmt.Sprintf(format, args...))
}

// Warn 警告日志
func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

// Warnf 格式化警告日志
func Warnf(format string, args ...interface{}) {
	Logger.Warn(fmt.Sprintf(format, args...))
}

// Error 错误日志
func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

// Errorf 格式化错误日志
func Errorf(format string, args ...interface{}) {
	Logger.Error(fmt.Sprintf(format, args...))
}

// Fatal 致命错误日志
func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

// Fatalf 格式化致命错误日志
func Fatalf(format string, args ...interface{}) {
	Logger.Fatal(fmt.Sprintf(format, args...))
}

// With 创建带字段的日志器
func With(fields ...zap.Field) *zap.Logger {
	return Logger.With(fields...)
}

// Sync 同步日志
func Sync() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}

// Close 关闭日志系统
func Close() {
	Sync()
}

// GetLogger 获取日志器实例
func GetLogger() *zap.Logger {
	return Logger
}
