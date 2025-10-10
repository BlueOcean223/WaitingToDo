package logAnalyzer

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"github.com/stretchr/testify/assert"
)

// TestParsePromptType tests the ParsePromptType function
func TestParsePromptType(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected PromptType
	}{
		{"Simplified", "simplified", PromptTypeSimplified},
		{"Simple", "simple", PromptTypeSimplified},
		{"Detailed", "detailed", PromptTypeDetailed},
		{"Comprehensive", "comprehensive", PromptTypeDetailed},
		{"Quick", "quick", PromptTypeQuick},
		{"Trend", "trend", PromptTypeTrend},
		{"Default", "unknown", PromptTypeSimplified},
		{"Empty", "", PromptTypeSimplified},
		{"Uppercase", "SIMPLIFIED", PromptTypeSimplified},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, ParsePromptType(tc.input))
		})
	}
}

// TestGetPromptConfig tests the GetPromptConfig function
func TestGetPromptConfig(t *testing.T) {
	testCases := []struct {
		name         string
		promptType   PromptType
		expectedType PromptType
		expectedKeys []string
	}{
		{
			name:         "Simplified",
			promptType:   PromptTypeSimplified,
			expectedType: PromptTypeSimplified,
			expectedKeys: []string{"AnalysisTime"},
		},
		{
			name:         "Detailed",
			promptType:   PromptTypeDetailed,
			expectedType: PromptTypeDetailed,
			expectedKeys: []string{"LogSource", "AnalysisTime"},
		},
		{
			name:         "Quick",
			promptType:   PromptTypeQuick,
			expectedType: PromptTypeQuick,
			expectedKeys: nil,
		},
		{
			name:         "Trend",
			promptType:   PromptTypeTrend,
			expectedType: PromptTypeTrend,
			expectedKeys: []string{"AnalysisTime"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := GetPromptConfig(tc.promptType)

			assert.NotNil(t, config)
			assert.Equal(t, tc.expectedType, config.Type)
			assert.NotNil(t, config.Parameters)
			assert.NotNil(t, config.Template)

			if len(tc.expectedKeys) > 0 {
				for _, key := range tc.expectedKeys {
					assert.Contains(t, config.Parameters, key)
				}
			} else {
				assert.Empty(t, config.Parameters)
			}
		})
	}
}

// TestBuildPrompt tests the BuildPrompt function
func TestBuildPrompt(t *testing.T) {
	logContent := "this is a test log"

	testCases := []struct {
		name           string
		promptType     PromptType
		expectedPrompt string
	}{
		{
			name:           "Simplified",
			promptType:     PromptTypeSimplified,
			expectedPrompt: SimplifiedAnalysisPrompt,
		},
		{
			name:           "Detailed",
			promptType:     PromptTypeDetailed,
			expectedPrompt: DetailedAnalysisPrompt,
		},
		{
			name:           "Quick",
			promptType:     PromptTypeQuick,
			expectedPrompt: QuickAnalysisPrompt,
		},
		{
			name:           "Trend",
			promptType:     PromptTypeTrend,
			expectedPrompt: TrendAnalysisPrompt,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := GetPromptConfig(tc.promptType)
			assert.NotNil(t, config)

			// Mock time-dependent parameters for consistent output
			if _, ok := config.Parameters["AnalysisTime"]; ok {
				config.Parameters["AnalysisTime"] = "2023-10-27T10:00:00Z"
			}
			if _, ok := config.Parameters["LogSource"]; ok {
				config.Parameters["LogSource"] = "test-source"
			}

			prompt := BuildPrompt(config, logContent)
			assert.NotEmpty(t, prompt)

			// Basic check if the log content is included
			assert.Contains(t, prompt, logContent)
		})
	}
}

// TestBuildPrompt_ErrorHandling tests error handling in BuildPrompt
func TestBuildPrompt_ErrorHandling(t *testing.T) {
	config := &PromptConfig{
		Type:       PromptTypeSimplified,
		Parameters: make(map[string]any),
		Template: prompt.FromMessages(schema.GoTemplate,
			&schema.Message{
				Role:    schema.User,
				Content: SimplifiedAnalysisPrompt,
			},
		),
	}

	// Missing "AnalysisTime" parameter should cause a formatting error
	prompt := BuildPrompt(config, "test log")
	assert.True(t, strings.HasPrefix(prompt, "Error formatting prompt:"))
}

// TestGetPromptConfig_TimeFormat tests the time format in GetPromptConfig
func TestGetPromptConfig_TimeFormat(t *testing.T) {
	// Test for Simplified
	configSimplified := GetPromptConfig(PromptTypeSimplified)
	_, errSimplified := time.Parse(time.RFC3339, fmt.Sprintf("%v", configSimplified.Parameters["AnalysisTime"]))
	assert.NoError(t, errSimplified, "Simplified prompt should have RFC3339 time format")

	// Test for Detailed
	configDetailed := GetPromptConfig(PromptTypeDetailed)
	_, errDetailed := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%v", configDetailed.Parameters["AnalysisTime"]))
	assert.NoError(t, errDetailed, "Detailed prompt should have '2006-01-02 15:04:05' time format")

	// Test for Trend
	configTrend := GetPromptConfig(PromptTypeTrend)
	_, errTrend := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%v", configTrend.Parameters["AnalysisTime"]))
	assert.NoError(t, errTrend, "Trend prompt should have '2006-01-02 15:04:05' time format")
}
