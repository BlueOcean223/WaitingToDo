package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

// GetFileMD5 计算文件MD5值
func GetFileMD5(file io.Reader) (string, error) {
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	md5Sum := hash.Sum(nil)
	return hex.EncodeToString(md5Sum), nil
}
