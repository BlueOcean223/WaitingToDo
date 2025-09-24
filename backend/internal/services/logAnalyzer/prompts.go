package logAnalyzer

import (
	"strings"
	"time"

	"github.com/tmc/langchaingo/prompts"
)

// 日志分析提示词模板
const (
	// 简化分析提示词 - 专注于实用性
	SimplifiedAnalysisPrompt = `你是一位专业的系统日志分析师。请分析以下日志内容，快速识别关键信息并提供实用建议。

## 日志内容
{{.LogContent}}

## 分析要求
请仔细分析日志内容，重点关注：
1. 系统运行状态
2. 错误和异常情况
3. 需要关注的问题
4. 实用的改进建议

请直接返回JSON格式结果，不要使用markdown代码块包装：

{
  "timestamp": "{{.AnalysisTime}}",
  "analysis_period": "过去24小时",
  "system_status": "系统状态(normal/warning/critical)",
  "summary": "系统运行状况简要总结(100字以内)",
  "log_statistics": {
    "total_entries": 日志条目总数,
    "error_count": 错误日志数量,
    "warning_count": 警告日志数量,
    "info_count": 信息日志数量
  },
  "key_findings": [
    "发现的关键问题1",
    "发现的关键问题2",
    "发现的关键问题3"
  ],
  "recommendations": [
    "具体可行的建议1",
    "具体可行的建议2",
    "具体可行的建议3"
  ]
}

重要提醒：
- 直接返回JSON内容，不要包装
- 确保JSON格式正确，所有字符串用双引号
- key_findings和recommendations数组最多包含5个项目
- 如果没有相关信息，使用空数组[]`

	// 详细分析提示词
	DetailedAnalysisPrompt = `你是一位专业的软件系统日志分析专家。请仔细分析以下过去24小时的系统日志，并提供专业的分析报告。

## 日志信息
- 分析时间范围：过去24小时
- 日志来源：{{.LogSource}}
- 分析时间：{{.AnalysisTime}}

## 日志内容
{{.LogContent}}

## 分析要求
请分析日志并返回严格的JSON格式结果，包含以下字段：

{
  "timestamp": "分析时间戳(ISO 8601格式)",
  "analysis_time": "{{.AnalysisTime}}",
  "log_source": "{{.LogSource}}",
  "system_status": "系统状态(normal/warning/critical)",
  "summary": "整体状况概述(200字以内)",
  "statistics": {
    "total_lines": 日志总行数,
    "error_count": 错误日志数量,
    "warn_count": 警告日志数量,
    "info_count": 信息日志数量,
    "debug_count": 调试日志数量
  },
  "error_analysis": {
    "critical_errors": ["关键错误列表"],
    "error_patterns": ["错误模式分析"],
    "root_causes": ["根因分析"]
  },
  "performance_analysis": {
    "bottlenecks": ["性能瓶颈"],
    "response_time_issues": ["响应时间问题"],
    "resource_usage": ["资源使用情况"]
  },
  "security_analysis": {
    "suspicious_activities": ["可疑活动"],
    "auth_failures": ["认证失败"],
    "security_risks": ["安全风险"]
  },
  "business_analysis": {
    "api_issues": ["API问题"],
    "business_errors": ["业务错误"],
    "workflow_problems": ["工作流问题"]
  },
  "key_issues": ["最重要的3-5个问题"],
  "recommendations": ["具体的优化建议"],
  "trends_and_predictions": ["趋势分析和预测"],
  "detailed_report": "详细的分析报告(Markdown格式，包含具体例子和数据)"
}

## 重要说明
1. 必须返回有效的JSON格式，不要包含任何其他文本，直接返回json内容，不要包装
2. 所有字符串值都要用双引号包围
3. 数组和对象要正确格式化
4. 如果某个字段没有相关信息，使用空数组[]或空字符串""
5. detailed_report字段可以包含完整的Markdown格式分析报告`

	// 快速分析提示词
	QuickAnalysisPrompt = `作为系统日志分析师，请快速分析以下日志内容，重点关注错误和异常：

日志内容：
{{.LogContent}}

请返回JSON格式结果：
{
  "timestamp": "分析时间戳(ISO 8601格式)",
  "analysis_type": "quick",
  "summary": "快速分析摘要",
  "statistics": {
    "error_count": 错误数量,
    "warn_count": 警告数量,
    "critical_issues": 严重问题数量
  },
  "top_issues": ["最严重的3个问题"],
  "recommendations": ["3-5个关键建议"],
  "system_status": "系统状态(normal/warning/critical)"
}

必须返回有效的JSON格式，不要包含其他文本。直接返回json内容，不要包装。`

	// 趋势分析提示词
	TrendAnalysisPrompt = `作为系统运维专家，请分析以下日志中的趋势和模式：

日志内容：
{{.LogContent}}

请返回JSON格式结果：
{
  "timestamp": "分析时间戳(ISO 8601格式)",
  "analysis_type": "trend",
  "summary": "趋势分析摘要",
  "error_trends": {
    "frequency_change": "错误频率变化趋势",
    "pattern_analysis": ["错误模式分析"],
    "time_distribution": "时间分布特征"
  },
  "performance_trends": {
    "response_time": "响应时间趋势",
    "resource_usage": "资源使用趋势",
    "bottlenecks": ["性能瓶颈变化"]
  },
  "user_behavior": {
    "access_patterns": ["用户访问模式"],
    "usage_trends": ["使用趋势"]
  },
  "system_load": {
    "load_changes": "系统负载变化",
    "peak_times": ["高峰时段"],
    "capacity_analysis": "容量分析"
  },
  "predictions": ["未来预测和建议"],
  "recommendations": ["基于趋势的优化建议"]
}

必须返回有效的JSON格式，不要包含其他文本。直接返回json内容，不要包装。`
)

// PromptType 定义提示词类型
type PromptType string

const (
	PromptTypeSimplified PromptType = "simplified" // 简化分析
	PromptTypeDetailed   PromptType = "detailed"   // 详细分析
	PromptTypeQuick      PromptType = "quick"      // 快速分析
	PromptTypeTrend      PromptType = "trend"      // 趋势分析
)

// ParsePromptType 将字符串转换为PromptType
func ParsePromptType(s string) PromptType {
	switch strings.ToLower(s) {
	case "simplified", "simple":
		return PromptTypeSimplified
	case "detailed", "comprehensive":
		return PromptTypeDetailed
	case "quick":
		return PromptTypeQuick
	case "trend":
		return PromptTypeTrend
	default:
		return PromptTypeSimplified // 默认使用简化分析
	}
}

// PromptConfig 提示词配置
type PromptConfig struct {
	Type       PromptType             `json:"type"`
	Parameters map[string]any         `json:"parameters"`
	Template   prompts.PromptTemplate `json:"-"` // langchaingo 提示词模板
}

// GetPromptConfig 获取默认提示词配置
func GetPromptConfig(promptType PromptType) *PromptConfig {
	config := &PromptConfig{
		Type:       promptType,
		Parameters: make(map[string]any),
	}

	// 设置默认参数和创建对应的提示词模板
	switch promptType {
	case PromptTypeSimplified:
		config.Parameters["AnalysisTime"] = time.Now().Format(time.RFC3339)
		config.Template = prompts.NewPromptTemplate(
			SimplifiedAnalysisPrompt,
			[]string{"LogContent", "AnalysisTime"},
		)
	case PromptTypeDetailed:
		config.Parameters["LogSource"] = "系统应用日志"
		config.Parameters["AnalysisTime"] = time.Now().Format("2006-01-02 15:04:05")
		config.Template = prompts.NewPromptTemplate(
			DetailedAnalysisPrompt,
			[]string{"LogContent", "LogSource", "AnalysisTime"},
		)
	case PromptTypeQuick:
		config.Template = prompts.NewPromptTemplate(
			QuickAnalysisPrompt,
			[]string{"LogContent"},
		)
	case PromptTypeTrend:
		config.Parameters["AnalysisTime"] = time.Now().Format("2006-01-02 15:04:05")
		config.Template = prompts.NewPromptTemplate(
			TrendAnalysisPrompt,
			[]string{"LogContent"},
		)
	}

	return config
}

// BuildPrompt 构建完整的提示词
func BuildPrompt(config *PromptConfig, logContent string) string {
	// 添加日志内容到参数
	config.Parameters["LogContent"] = logContent

	// 使用 langchaingo 的 PromptTemplate 格式化提示词
	if config.Template.Template != "" {
		result, err := config.Template.Format(config.Parameters)
		if err != nil {
			// 如果格式化失败，返回错误信息或使用原始模板
			return "Error formatting prompt: " + err.Error()
		}
		return result
	}

	// 如果没有模板，返回空字符串
	return ""
}
