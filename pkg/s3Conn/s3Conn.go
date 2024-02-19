package s3Conn

import (
	"context"
	"log"

	configApp "github.com/NatthawutSK/NoTeams-Backend/config"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func S3Connect(cfg configApp.IS3Config) *s3.Client {
	s3Cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.S3AccessKey(), cfg.S3SecretKey(), cfg.S3Session())),
		config.WithRegion(cfg.S3Region()),
	)
	if err != nil {
		log.Fatal("configuration S3 error, " + err.Error())
	}

	client := s3.NewFromConfig(s3Cfg)

	return client
}
