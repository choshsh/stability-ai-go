package util

import (
	"bytes"
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
	"time"
)

var client *s3.Client

func init() {
	var cfg aws.Config
	var err error

	switch os.Getenv("GO_ENV") {
	case "production":
		cfg, err = config.LoadDefaultConfig(context.Background())
		if err != nil {
			panic(err)
		}

	default:
		const (
			TestRegion     = "ap-northeast-2"
			TestS3Endpoint = "http://localhost:9000"
			TestS3Cred     = "minio1234"
		)

		customResolver := aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					PartitionID:       "aws",
					URL:               TestS3Endpoint,
					SigningRegion:     TestRegion,
					HostnameImmutable: true,
				}, nil
			},
		)

		cfg, err = config.LoadDefaultConfig(
			context.Background(),
			config.WithEndpointResolverWithOptions(customResolver),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(TestS3Cred, TestS3Cred, "")),
		)
		if err != nil {
			panic(err)
		}
	}

	client = s3.NewFromConfig(cfg)
}

type S3UploadInput struct {
	Bucket        *string
	ID            *string
	RawBytes      *[]byte
	ContentType   *string
	FileExtension *string
}

func (si *S3UploadInput) Validate() bool {
	if si.Bucket == nil || si.ID == nil || si.RawBytes == nil ||
		si.ContentType == nil || si.FileExtension == nil {
		return false
	}
	return true
}

func S3Upload(input S3UploadInput) (string, error) {
	if !input.Validate() {
		return "", errors.New("invalid input for s3 upload")
	}

	s3key := strings.Join([]string{time.Now().Format("2006/01/02/"), *input.ID, *input.FileExtension}, "")
	reader := bytes.NewReader(*input.RawBytes)

	_, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      input.Bucket,
		Key:         aws.String(s3key),
		Body:        reader,
		ContentType: input.ContentType,
	})

	if err != nil {
		log.Error().Err(err).Msg("Upload to S3")
		return "", err
	}
	return s3key, nil
}
