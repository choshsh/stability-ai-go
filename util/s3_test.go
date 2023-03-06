package util

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	TestRegion     = "ap-northeast-2"
	TestS3Endpoint = "http://localhost:9000"
	TestS3Bucket   = "test"
	TestS3Cred     = "minio1234"
)

func TestCreateBucket(t *testing.T) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:       "aws",
			URL:               TestS3Endpoint,
			SigningRegion:     TestRegion,
			HostnameImmutable: true,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(TestS3Cred, TestS3Cred, "")),
	)
	assert.NoError(t, err)

	client := s3.NewFromConfig(cfg)
	_, err = client.CreateBucket(context.Background(), &s3.CreateBucketInput{
		Bucket: aws.String(TestS3Bucket),
	})
	assert.NoError(t, err)
}
