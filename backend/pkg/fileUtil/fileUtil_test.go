package fileUtil

import (
	"os"
	"strings"
	"testing"
)

func TestGetFileMD5(t *testing.T) {
	// 创建临时测试文件
	testContent := "This is a test file for MD5 calculation"
	tempFile, err := os.CreateTemp("", "test_md5_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() {
		err = tempFile.Close()
		if err != nil {
			t.Fatalf("Failed to close temp file: %v", err)
		}
		err = os.Remove(tempFile.Name()) // 清理临时文件
		if err != nil {
			t.Fatalf("Failed to remove temp file: %v", err)
		}
	}()

	// 写入测试内容
	if _, err := tempFile.WriteString(testContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// 重置文件指针到开头
	if _, err := tempFile.Seek(0, 0); err != nil {
		t.Fatalf("Failed to seek temp file: %v", err)
	}

	// 测试 GetFileMD5 函数
	md5Hash, err := GetFileMD5(tempFile)
	if err != nil {
		t.Fatalf("GetFileMD5 failed: %v", err)
	}

	// 验证 MD5 值不为空且格式正确
	if md5Hash == "" {
		t.Error("MD5 hash should not be empty")
	}

	// MD5 哈希值应该是 32 位十六进制字符串
	if len(md5Hash) != 32 {
		t.Errorf("Expected MD5 hash length 32, got %d", len(md5Hash))
	}

	// 验证是否为有效的十六进制字符串
	for _, char := range md5Hash {
		if !strings.ContainsRune("0123456789abcdef", char) {
			t.Errorf("MD5 hash contains invalid character: %c", char)
			break
		}
	}

	t.Logf("MD5 hash: %s", md5Hash)
}

// TestGetFileMD5WithStringReader 测试使用字符串读取器
func TestGetFileMD5WithStringReader(t *testing.T) {
	testContent := "Hello, World!"
	reader := strings.NewReader(testContent)

	md5Hash, err := GetFileMD5(reader)
	if err != nil {
		t.Fatalf("GetFileMD5 failed: %v", err)
	}

	// 预期的 MD5 值（"Hello, World!" 的 MD5）
	expected := "65a8e27d8879283831b664bd8b7f0ad4"
	if md5Hash != expected {
		t.Errorf("Expected MD5 %s, got %s", expected, md5Hash)
	}

	t.Logf("MD5 hash: %s", md5Hash)
}
