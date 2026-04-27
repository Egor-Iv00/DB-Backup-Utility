package Cloud

import (
	"context"
	"dbtool/DBinterface"
	"fmt"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
)

type CloudConfig struct {
	IsUse      bool
	endpoint   string
	accessKey  string
	secretKey  string
	useSSL     string
	bucketName string
}

var GlobalMinio *minio.Client

func InitCloud(CloudConf *CloudConfig, GlobalViper *viper.Viper) error {
	CloudConf.endpoint = GlobalViper.GetString("endpoint")
	CloudConf.accessKey = GlobalViper.GetString("accesskey")
	CloudConf.secretKey = GlobalViper.GetString("secretkey")
	CloudConf.bucketName = GlobalViper.GetString("bucketname")

	MClient, err := minio.New(CloudConf.endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(CloudConf.accessKey, CloudConf.secretKey, ""),
		Secure: true,
	})

	if err != nil {
		return fmt.Errorf("cloud connect error: %w", err)
	}
	GlobalMinio = MClient
	return nil
}

func ConnectToCloud(CloudConf *CloudConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if _, err := GlobalMinio.BucketExists(ctx, CloudConf.bucketName); err != nil {
		return fmt.Errorf("connect cloud error: %w", err)
	}
	return nil
}

func BackupCloud(CloudConf CloudConfig, DBconf DBinterface.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	_, err := GlobalMinio.FPutObject(ctx, CloudConf.bucketName, filepath.Base(DBconf.FilePath), DBconf.FilePath, minio.PutObjectOptions{})
	return err
}

func RestoreCloud(CloudConf CloudConfig, DBconf DBinterface.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err := GlobalMinio.FGetObject(ctx, CloudConf.bucketName, filepath.Base(DBconf.FilePath), DBconf.FilePath, minio.GetObjectOptions{})
	return err
}
