package configs

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioConfig struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
}

var MinioClient *minio.Client

func InitMinioClient() error {
	var err error
	// 获取配置信息
	endpoint := AppConfigs.MinioConfig.Endpoint
	accessKeyID := AppConfigs.MinioConfig.AccessKeyID
	secretAccessKey := AppConfigs.MinioConfig.SecretAccessKey

	// 初始化minio客户端
	MinioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		return err
	}

	return nil
}
