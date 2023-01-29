package common

import (
	"bytes"
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

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

	cfg, err := config.LoadDefaultConfig(context.Background())
	client := s3.NewFromConfig(cfg)

	s3key := strings.Join([]string{time.Now().Format("2006/01/02/"), *input.ID, *input.FileExtension}, "")
	reader := bytes.NewReader(*input.RawBytes)

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
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
