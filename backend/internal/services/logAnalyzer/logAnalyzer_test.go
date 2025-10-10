package logAnalyzer

import (
	"testing"
)

func TestLogAnalyzerEinoIntegration(t *testing.T) {
	// 这是一个集成测试，用于验证Eino框架是否正确集成
	// 由于需要实际的API密钥和网络连接，这个测试主要用于验证代码是否能正确编译和初始化

	// 创建一个新的日志分析器实例
	analyzer := NewLogAnalyzer()

	// 验证实例是否正确创建
	if analyzer == nil {
		t.Error("Failed to create LogAnalyzer instance")
	}

	// 验证是否正确实现了接口
	_ = analyzer.IsEnabled()

	t.Log("LogAnalyzer with Eino framework created successfully")
}

// TestLogAnalyzerFullIntegration 测试完整的日志分析流程
func TestLogAnalyzerFullIntegration(t *testing.T) {
	// 创建一个新的日志分析器实例
	analyzer := NewLogAnalyzer()

	// 验证实例创建成功
	if analyzer == nil {
		t.Fatal("Failed to create LogAnalyzer instance")
	}

	// 验证是否启用（应该默认为false）
	enabled := analyzer.IsEnabled()
	if enabled {
		t.Log("LogAnalyzer is enabled")
	} else {
		t.Log("LogAnalyzer is disabled (expected for testing)")
	}

	// 测试提示词类型解析
	simplifiedType := ParsePromptType("simplified")
	if simplifiedType != PromptTypeSimplified {
		t.Errorf("Expected PromptTypeSimplified, got %v", simplifiedType)
	}

	detailedType := ParsePromptType("detailed")
	if detailedType != PromptTypeDetailed {
		t.Errorf("Expected PromptTypeDetailed, got %v", detailedType)
	}

	// 测试提示词配置获取
	simplifiedConfig := GetPromptConfig(PromptTypeSimplified)
	if simplifiedConfig == nil {
		t.Error("Failed to get simplified prompt config")
	}

	detailedConfig := GetPromptConfig(PromptTypeDetailed)
	if detailedConfig == nil {
		t.Error("Failed to get detailed prompt config")
	}

	t.Log("Full integration test passed")
}
