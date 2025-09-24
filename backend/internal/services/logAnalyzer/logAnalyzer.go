package logAnalyzer

import (
	"backend/internal/configs"
	"backend/pkg/logger"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

// LogEntry 表示JSON格式的日志条目
type LogEntry struct {
	Level  string `json:"level"`
	Time   string `json:"time"`
	Caller string `json:"caller"`
	Msg    string `json:"msg"`
}

// LogAnalyzer 日志分析器
type LogAnalyzer struct {
	llm    *openai.LLM
	config configs.LogAnalyzerConfig
}

// NewLogAnalyzer 创建新的日志分析器
func NewLogAnalyzer() *LogAnalyzer {
	logAnalyzer := &LogAnalyzer{
		llm:    nil,
		config: configs.AppConfigs.LogAnalyzerConfig,
	}

	if !logAnalyzer.config.Enabled {
		return logAnalyzer
	}

	// 创建LLM客户端
	llm, err := openai.New(
		openai.WithModel(logAnalyzer.config.Model),
		openai.WithBaseURL(logAnalyzer.config.URL),
		openai.WithToken(logAnalyzer.config.APIKey),
	)
	if err != nil {
		logger.Error("创建LLM客户端失败", logger.Err(err))
		return logAnalyzer
	}

	logAnalyzer.llm = llm
	return logAnalyzer
}

// getReportPath 获取报告存储路径
func getReportPath() string {
	reportPath := filepath.Dir(configs.AppConfigs.LogConfig.Filename)
	if err := os.MkdirAll(reportPath, 0755); err != nil {
		logger.Error("创建分析报告目录失败", logger.Err(err))
		return filepath.Join("logs", "analysis_reports")
	}
	return reportPath
}

// AnalyzeLogs 分析日志并保存报告
func (la *LogAnalyzer) AnalyzeLogs() error {
	if !la.config.Enabled {
		return fmt.Errorf("LLM分析未启用")
	}

	// 从配置获取分析时间范围
	analysisHours := la.config.AnalysisHours
	if analysisHours <= 0 {
		analysisHours = 24 // 默认24小时
	}

	// 获取指定时间范围的日志内容
	logContent, err := la.getLogContentSince(time.Duration(analysisHours) * time.Hour)
	if err != nil {
		return fmt.Errorf("获取日志内容失败: %w", err)
	}

	// 从配置获取最大日志大小限制
	maxLogSize := la.config.MaxLogSize
	if maxLogSize <= 0 {
		maxLogSize = 100000 // 默认100KB
	}

	// 限制日志内容大小
	if len(logContent) > maxLogSize {
		logContent = logContent[:maxLogSize] + "\n... [日志内容过大，已截断]"
	}

	// 从配置获取分析级别
	promptType := ParsePromptType(la.config.AnalysisLevel)

	// 构建提示词
	prompt := la.buildPrompt(logContent, promptType)

	// 调用LLM进行分析
	result, err := la.callLLM(prompt)
	if err != nil {
		return fmt.Errorf("LLM分析失败: %w", err)
	}

	// 清理JSON结果（移除可能的markdown包装）
	cleanedResult := la.CleanJSONResult(result)

	// 验证JSON格式
	if err = la.validateJSON(cleanedResult); err != nil {
		logger.Error("LLM返回的JSON格式无效", logger.Err(err))
		// 如果JSON无效，仍然保存原始结果，但添加错误标记
		cleanedResult = fmt.Sprintf(`{"error": "JSON格式无效", "raw_result": %q, "validation_error": %q, "timestamp": "%s"}`,
			cleanedResult, err.Error(), time.Now().Format(time.RFC3339))
	}

	// 保存LLM生成的JSON内容
	_, err = la.saveJSONReport(cleanedResult)
	if err != nil {
		return fmt.Errorf("保存分析报告失败: %w", err)
	}

	// 如果是detailed模式，额外生成Markdown文档
	if promptType == PromptTypeDetailed {
		err = la.SaveMarkdownReport(cleanedResult)
		if err != nil {
			logger.Error("保存Markdown报告失败", logger.Err(err))
			// 不返回错误，因为JSON报告已经成功保存
		}
	}

	return nil
}

// validateJSON 验证JSON格式并检查必要字段
func (la *LogAnalyzer) validateJSON(jsonStr string) error {
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return fmt.Errorf("JSON解析失败: %w", err)
	}

	// 检查必要字段
	requiredFields := []string{"timestamp"}
	for _, field := range requiredFields {
		if _, exists := result[field]; !exists {
			return fmt.Errorf("缺少必要字段: %s", field)
		}
	}

	return nil
}

// getLogContentSince 获取指定时间范围内的日志内容
func (la *LogAnalyzer) getLogContentSince(duration time.Duration) (string, error) {
	logPath := configs.AppConfigs.LogConfig.Filename

	file, err := os.Open(logPath)
	if err != nil {
		return "", fmt.Errorf("打开日志文件失败: %w", err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			logger.Error("关闭日志文件失败", logger.Err(err))
		}
	}(file)

	var content strings.Builder
	scanner := bufio.NewScanner(file)

	cutoffTime := time.Now().Add(-duration)
	lineCount := 0
	errorCount := 0
	warnCount := 0
	infoCount := 0
	debugCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		// 跳过空行
		if strings.TrimSpace(line) == "" {
			continue
		}

		// 尝试解析JSON格式的日志
		var logEntry LogEntry
		if err := json.Unmarshal([]byte(line), &logEntry); err == nil {
			// 解析时间戳（RFC3339格式：2025-09-19T21:27:45.273+0800）
			if logTime, err := time.Parse(time.RFC3339, logEntry.Time); err == nil {
				if logTime.Before(cutoffTime) {
					continue
				}
			}

			// 统计日志级别
			switch strings.ToLower(logEntry.Level) {
			case "error":
				errorCount++
			case "warn", "warning":
				warnCount++
			case "info":
				infoCount++
			case "debug":
				debugCount++
			}
		} else {
			// 如果不是JSON格式，尝试传统格式解析（向后兼容）
			// 检查时间戳（假设日志格式为 [时间戳] 日志内容）
			if strings.Contains(line, "[") && strings.Contains(line, "]") {
				timestampStart := strings.Index(line, "[") + 1
				timestampEnd := strings.Index(line, "]")
				if timestampEnd > timestampStart {
					timestampStr := line[timestampStart:timestampEnd]
					if logTime, err := time.Parse("2006-01-02 15:04:05", timestampStr); err == nil {
						if logTime.Before(cutoffTime) {
							continue
						}
					}
				}
			}

			// 统计日志级别
			if strings.Contains(line, "ERROR") {
				errorCount++
			} else if strings.Contains(line, "WARN") {
				warnCount++
			} else if strings.Contains(line, "INFO") {
				infoCount++
			} else if strings.Contains(line, "DEBUG") {
				debugCount++
			}
		}

		content.WriteString(line)
		content.WriteString("\n")

		// 限制行数避免内存问题
		if lineCount > 10000 {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("读取日志文件失败: %w", err)
	}

	logger.Info("日志统计", logger.String("error_count", fmt.Sprintf("%d", errorCount)),
		logger.String("warn_count", fmt.Sprintf("%d", warnCount)),
		logger.String("info_count", fmt.Sprintf("%d", infoCount)),
		logger.String("debug_count", fmt.Sprintf("%d", debugCount)))

	return content.String(), nil
}

// buildPrompt 构建指定类型的提示词
func (la *LogAnalyzer) buildPrompt(logContent string, promptType PromptType) string {
	// 使用新的提示词系统
	config := GetPromptConfig(promptType)
	return BuildPrompt(config, logContent)
}

// callLLM 调用LLM进行分析
func (la *LogAnalyzer) callLLM(prompt string) (string, error) {
	if la.llm == nil {
		return "LLM客户端未初始化", nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(la.config.Timeout)*time.Second)
	defer cancel()

	result, err := llms.GenerateFromSinglePrompt(ctx, la.llm, prompt)
	if err != nil {
		return "", fmt.Errorf("LLM生成失败: %w", err)
	}

	return result, nil
}

// saveJSONReport 保存JSON格式的分析报告，返回文件路径
func (la *LogAnalyzer) saveJSONReport(jsonResult string) (string, error) {
	// 从配置获取报告路径
	reportPath := la.config.ReportPath
	if reportPath == "" {
		reportPath = getReportPath() // 使用默认路径
	}

	// 确保报告目录存在
	if err := os.MkdirAll(reportPath, 0755); err != nil {
		return "", fmt.Errorf("创建报告目录失败: %w", err)
	}

	// 生成报告文件名
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("log_analysis_%s.json", timestamp)
	reportFile := filepath.Join(reportPath, filename)

	// 格式化JSON（如果有效的话）
	var formattedJSON []byte
	var temp interface{}
	if err := json.Unmarshal([]byte(jsonResult), &temp); err == nil {
		// JSON有效，格式化输出
		formattedJSON, _ = json.MarshalIndent(temp, "", "  ")
	} else {
		// JSON无效，直接保存原始内容
		formattedJSON = []byte(jsonResult)
	}

	// 写入文件
	err := os.WriteFile(reportFile, formattedJSON, 0644)
	if err != nil {
		return "", fmt.Errorf("写入报告文件失败: %w", err)
	}

	logger.GetLogger().Info("日志分析报告已保存", logger.String("file", reportFile))
	return reportFile, nil
}

// CleanJSONResult 清理LLM返回的JSON内容，移除可能的markdown代码块包装
func (la *LogAnalyzer) CleanJSONResult(result string) string {
	// 去除首尾空白
	result = strings.TrimSpace(result)

	// 移除markdown代码块标记
	if strings.HasPrefix(result, "```json") {
		result = strings.TrimPrefix(result, "```json")
	}
	if strings.HasPrefix(result, "```") {
		result = strings.TrimPrefix(result, "```")
	}
	if strings.HasSuffix(result, "```") {
		result = strings.TrimSuffix(result, "```")
	}

	// 再次去除空白
	result = strings.TrimSpace(result)

	// 如果结果不是以{开头，尝试找到第一个{
	if !strings.HasPrefix(result, "{") {
		if idx := strings.Index(result, "{"); idx != -1 {
			result = result[idx:]
		}
	}

	// 如果结果不是以}结尾，尝试找到最后一个}
	if !strings.HasSuffix(result, "}") {
		if idx := strings.LastIndex(result, "}"); idx != -1 {
			result = result[:idx+1]
		}
	}

	return result
}

// SaveMarkdownReport 从JSON结果中提取detailed_report字段并保存为Markdown文档
func (la *LogAnalyzer) SaveMarkdownReport(jsonResult string) error {
	// 解析JSON
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonResult), &result); err != nil {
		return fmt.Errorf("解析JSON失败: %w", err)
	}

	// 提取detailed_report字段
	detailedReport, exists := result["detailed_report"]
	if !exists {
		return fmt.Errorf("JSON中未找到detailed_report字段")
	}

	// 转换为字符串
	markdownContent, ok := detailedReport.(string)
	if !ok {
		return fmt.Errorf("detailed_report字段不是字符串类型")
	}

	// 从配置获取报告路径
	reportPath := la.config.ReportPath
	if reportPath == "" {
		reportPath = getReportPath() // 使用默认路径
	}

	// 确保报告目录存在
	if err := os.MkdirAll(reportPath, 0755); err != nil {
		return fmt.Errorf("创建报告目录失败: %w", err)
	}

	// 生成Markdown文件名
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("log_analysis_detailed_%s.md", timestamp)
	reportFile := filepath.Join(reportPath, filename)

	// 写入Markdown文件
	err := os.WriteFile(reportFile, []byte(markdownContent), 0644)
	if err != nil {
		return fmt.Errorf("写入Markdown报告文件失败: %w", err)
	}

	logger.GetLogger().Info("详细日志分析Markdown报告已保存", logger.String("file", reportFile))
	return nil
}

// IsEnabled 检查是否启用
func (la *LogAnalyzer) IsEnabled() bool {
	return la.config.Enabled
}
